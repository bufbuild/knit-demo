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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/film/v1/filmv1connect"
	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/person/v1/personv1connect"
	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/planet/v1/planetv1connect"
	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/relations/v1/relationsv1connect"
	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/species/v1/speciesv1connect"
	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/starship/v1/starshipv1connect"
	"buf.build/gen/go/bufbuild/knit-demo/bufbuild/connect-go/buf/knit/demo/swapi/vehicle/v1/vehiclev1connect"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/bufbuild/knit-demo/go/internal"
	"github.com/bufbuild/knit-demo/go/internal/swapi"
	"github.com/bufbuild/knit-go"
	"github.com/rs/cors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	bindAddr := flags.String("bind", "127.0.0.1", "The local IP on which to listen for HTTP requests. Use 0.0.0.0 to bind to all interfaces.")
	port := flags.Int("port", 30485, "The port on which to listen for HTTP requests.")
	var serviceNames multiStringFlag
	flags.Var(&serviceNames, "service", "The set of services to implement. If not specified, all services will be implemented.")
	embedGateway := flags.Bool("embed-gateway", false, "If true, the server will embed a Knit gateway and also expose the Knit protocol.")

	_ = flags.Parse(os.Args[1:])

	handler := swapi.NewHandler()

	mux := http.NewServeMux()

	allServices := map[string]struct {
		register func()
	}{
		filmv1connect.FilmServiceName: {
			register: func() {
				mux.Handle(filmv1connect.NewFilmServiceHandler(handler))
			},
		},
		relationsv1connect.FilmResolverServiceName: {
			register: func() {
				mux.Handle(relationsv1connect.NewFilmResolverServiceHandler(handler))
			},
		},
		personv1connect.PersonServiceName: {
			register: func() {
				mux.Handle(personv1connect.NewPersonServiceHandler(handler))
			},
		},
		relationsv1connect.PersonResolverServiceName: {
			register: func() {
				mux.Handle(relationsv1connect.NewPersonResolverServiceHandler(handler))
			},
		},
		planetv1connect.PlanetServiceName: {
			register: func() {
				mux.Handle(planetv1connect.NewPlanetServiceHandler(handler))
			},
		},
		relationsv1connect.PlanetResolverServiceName: {
			register: func() {
				mux.Handle(relationsv1connect.NewPlanetResolverServiceHandler(handler))
			},
		},
		speciesv1connect.SpeciesServiceName: {
			register: func() {
				mux.Handle(speciesv1connect.NewSpeciesServiceHandler(handler))
			},
		},
		relationsv1connect.SpeciesResolverServiceName: {
			register: func() {
				mux.Handle(relationsv1connect.NewSpeciesResolverServiceHandler(handler))
			},
		},
		starshipv1connect.StarshipServiceName: {
			register: func() {
				mux.Handle(starshipv1connect.NewStarshipServiceHandler(handler))
			},
		},
		relationsv1connect.StarshipResolverServiceName: {
			register: func() {
				mux.Handle(relationsv1connect.NewStarshipResolverServiceHandler(handler))
			},
		},
		vehiclev1connect.VehicleServiceName: {
			register: func() {
				mux.Handle(vehiclev1connect.NewVehicleServiceHandler(handler))
			},
		},
		relationsv1connect.VehicleResolverServiceName: {
			register: func() {
				mux.Handle(relationsv1connect.NewVehicleResolverServiceHandler(handler))
			},
		},
	}
	// if none specified, we'll expose all of them
	if len(serviceNames) == 0 {
		for serviceName := range allServices {
			serviceNames = append(serviceNames, serviceName)
		}
		sort.Strings(serviceNames)
	} else if *embedGateway {
		// if embedding gateway, cannot indicate service names; they must all be exposed
		log.Fatalln("cannot use --embed-gateway while only exposing subset of services via --service")
	}

	for _, serviceName := range serviceNames {
		info, ok := allServices[serviceName]
		if !ok {
			log.Fatalf("unknown service %q\n", serviceName)
		}
		log.Printf("registering handler for %q\n", serviceName)
		info.register()
	}

	// support gRPC reflection
	reflector := grpcreflect.NewStaticReflector(serviceNames...)
	// TODO: also include the Knit service once it is public
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddr, *port))
	if err != nil {
		log.Fatalln(err)
	}

	if *embedGateway {
		routeURL := &url.URL{
			Scheme: "http",
			Host:   listener.Addr().String(),
			Path:   "/",
		}
		gateway := &knit.Gateway{
			Client:                   http.DefaultClient,
			Route:                    routeURL,
			MaxParallelismPerRequest: 10, // TODO: flag?
		}
		for _, svc := range serviceNames {
			if err := gateway.AddServiceByName(protoreflect.FullName(svc)); err != nil {
				log.Fatalln(err)
			}
		}
		mux.Handle(gateway.AsHandler())
	}

	if err := internal.Serve(context.Background(), listener, cors.AllowAll().Handler(mux)); err != nil {
		log.Fatalln(err)
	}
}

type multiStringFlag []string

func (m *multiStringFlag) String() string {
	return strings.Join(*m, ", ")
}

func (m *multiStringFlag) Set(s string) error {
	vals := strings.Split(s, ",")
	for _, val := range vals {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		*m = append(*m, val)
	}
	return nil
}
