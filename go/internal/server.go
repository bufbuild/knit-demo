// Copyright 2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Serve handles serving the given handler via HTTP using the given listener. The
// server will support H2C (HTTP/2 over plaintext) and will log a single line of
// output for each HTTP request that briefly describes the call and its status.
func Serve(ctx context.Context, listener net.Listener, handler http.Handler) error {
	loggingHandler := http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		// Instead of 404'ing on home page, redirect to repo README.
		if req.URL.Path == "/" && req.Method == http.MethodGet {
			// TODO: eventually this should redirect to a proper website for Knit
			http.Redirect(respWriter, req, "https://github.com/bufbuild/knit-demo/blob/main/README.md", http.StatusFound)
			return
		}

		start := time.Now()
		intercepted, respWriter := intercept(respWriter)
		handler.ServeHTTP(respWriter, req)
		logRequest(req, intercepted.status, intercepted.size, time.Since(start))
	})
	log.Printf("Listening on %s for HTTP requests...\n", listener.Addr().String())
	svr := http.Server{
		Handler:           h2c.NewHandler(loggingHandler, &http2.Server{}),
		ReadHeaderTimeout: 20 * time.Second,
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var alreadyShutdown atomic.Bool
	go func() {
		<-ctx.Done()
		if !alreadyShutdown.Load() {
			// ctx is already cancelled, so we need one with more time
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_ = svr.Shutdown(timeoutCtx) //nolint:contextcheck
		}
	}()

	err := svr.Serve(listener)
	alreadyShutdown.Store(true)
	return err
}

func intercept(w http.ResponseWriter) (*interceptWriter, http.ResponseWriter) {
	intercepted := &interceptWriter{w: w, status: "200"}
	if f, ok := w.(http.Flusher); ok {
		// make sure we return flusher if input writer can flush
		return intercepted, writerAndFlusher{ResponseWriter: intercepted, Flusher: f}
	}
	return intercepted, intercepted
}

type interceptWriter struct {
	w            http.ResponseWriter
	alreadyWrote bool
	status       string
	size         int
}

func (i *interceptWriter) Header() http.Header {
	return i.w.Header()
}

func (i *interceptWriter) Write(bytes []byte) (int, error) {
	i.alreadyWrote = true
	n, err := i.w.Write(bytes)
	i.size += n
	if err != nil && !strings.Contains(i.status, "(") {
		i.status = "499 (" + i.status + ")"
	}
	return n, err
}

func (i *interceptWriter) WriteHeader(statusCode int) {
	if !i.alreadyWrote {
		i.status = strconv.Itoa(statusCode)
	}
	i.w.WriteHeader(statusCode)
}

type writerAndFlusher struct {
	http.ResponseWriter
	http.Flusher
}

func logRequest(r *http.Request, status string, bodySize int, latency time.Duration) {
	log.Printf("%s %s  %s  %s  %db  %v\n", r.Method, r.RequestURI, r.RemoteAddr, status, bodySize, latency)
}
