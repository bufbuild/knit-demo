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

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: buf/knit/demo/swapi/species/v1/species.proto

package speciesv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/species/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// SpeciesServiceName is the fully-qualified name of the SpeciesService service.
	SpeciesServiceName = "buf.knit.demo.swapi.species.v1.SpeciesService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// SpeciesServiceGetSpeciesProcedure is the fully-qualified name of the SpeciesService's GetSpecies
	// RPC.
	SpeciesServiceGetSpeciesProcedure = "/buf.knit.demo.swapi.species.v1.SpeciesService/GetSpecies"
	// SpeciesServiceListSpeciesProcedure is the fully-qualified name of the SpeciesService's
	// ListSpecies RPC.
	SpeciesServiceListSpeciesProcedure = "/buf.knit.demo.swapi.species.v1.SpeciesService/ListSpecies"
)

// SpeciesServiceClient is a client for the buf.knit.demo.swapi.species.v1.SpeciesService service.
type SpeciesServiceClient interface {
	GetSpecies(context.Context, *connect_go.Request[v1.GetSpeciesRequest]) (*connect_go.Response[v1.GetSpeciesResponse], error)
	ListSpecies(context.Context, *connect_go.Request[v1.ListSpeciesRequest]) (*connect_go.Response[v1.ListSpeciesResponse], error)
}

// NewSpeciesServiceClient constructs a client for the buf.knit.demo.swapi.species.v1.SpeciesService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewSpeciesServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) SpeciesServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &speciesServiceClient{
		getSpecies: connect_go.NewClient[v1.GetSpeciesRequest, v1.GetSpeciesResponse](
			httpClient,
			baseURL+SpeciesServiceGetSpeciesProcedure,
			opts...,
		),
		listSpecies: connect_go.NewClient[v1.ListSpeciesRequest, v1.ListSpeciesResponse](
			httpClient,
			baseURL+SpeciesServiceListSpeciesProcedure,
			opts...,
		),
	}
}

// speciesServiceClient implements SpeciesServiceClient.
type speciesServiceClient struct {
	getSpecies  *connect_go.Client[v1.GetSpeciesRequest, v1.GetSpeciesResponse]
	listSpecies *connect_go.Client[v1.ListSpeciesRequest, v1.ListSpeciesResponse]
}

// GetSpecies calls buf.knit.demo.swapi.species.v1.SpeciesService.GetSpecies.
func (c *speciesServiceClient) GetSpecies(ctx context.Context, req *connect_go.Request[v1.GetSpeciesRequest]) (*connect_go.Response[v1.GetSpeciesResponse], error) {
	return c.getSpecies.CallUnary(ctx, req)
}

// ListSpecies calls buf.knit.demo.swapi.species.v1.SpeciesService.ListSpecies.
func (c *speciesServiceClient) ListSpecies(ctx context.Context, req *connect_go.Request[v1.ListSpeciesRequest]) (*connect_go.Response[v1.ListSpeciesResponse], error) {
	return c.listSpecies.CallUnary(ctx, req)
}

// SpeciesServiceHandler is an implementation of the buf.knit.demo.swapi.species.v1.SpeciesService
// service.
type SpeciesServiceHandler interface {
	GetSpecies(context.Context, *connect_go.Request[v1.GetSpeciesRequest]) (*connect_go.Response[v1.GetSpeciesResponse], error)
	ListSpecies(context.Context, *connect_go.Request[v1.ListSpeciesRequest]) (*connect_go.Response[v1.ListSpeciesResponse], error)
}

// NewSpeciesServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewSpeciesServiceHandler(svc SpeciesServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(SpeciesServiceGetSpeciesProcedure, connect_go.NewUnaryHandler(
		SpeciesServiceGetSpeciesProcedure,
		svc.GetSpecies,
		opts...,
	))
	mux.Handle(SpeciesServiceListSpeciesProcedure, connect_go.NewUnaryHandler(
		SpeciesServiceListSpeciesProcedure,
		svc.ListSpecies,
		opts...,
	))
	return "/buf.knit.demo.swapi.species.v1.SpeciesService/", mux
}

// UnimplementedSpeciesServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedSpeciesServiceHandler struct{}

func (UnimplementedSpeciesServiceHandler) GetSpecies(context.Context, *connect_go.Request[v1.GetSpeciesRequest]) (*connect_go.Response[v1.GetSpeciesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("buf.knit.demo.swapi.species.v1.SpeciesService.GetSpecies is not implemented"))
}

func (UnimplementedSpeciesServiceHandler) ListSpecies(context.Context, *connect_go.Request[v1.ListSpeciesRequest]) (*connect_go.Response[v1.ListSpeciesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("buf.knit.demo.swapi.species.v1.SpeciesService.ListSpecies is not implemented"))
}
