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

package swapi

// Use "go generate" to re-create the static snapshot of swapi data.
//go:generate go run ../cmd/gendata swapi ./swapi_data.gen.go
//go:generate gofmt -w -s ./swapi_data.gen.go

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/film/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/film/v1/filmv1connect"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/person/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/person/v1/personv1connect"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/planet/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/planet/v1/planetv1connect"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/relations/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/relations/v1/relationsv1connect"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/species/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/species/v1/speciesv1connect"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/starship/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/starship/v1/starshipv1connect"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/vehicle/v1"
	"github.com/bufbuild/knit-demo/go/gen/buf/knit/demo/swapi/vehicle/v1/vehiclev1connect"
	"github.com/peterhellberg/swapi"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler implements the Star Wars API.
//
// For performance, it uses an in-memory snapshot of the data (e.g. instead
// of sending queries to swapi.dev). Run "go generate" for this package to
// re-generate the snapshot.
type Handler struct {
	filmv1connect.UnimplementedFilmServiceHandler
	personv1connect.UnimplementedPersonServiceHandler
	starshipv1connect.UnimplementedStarshipServiceHandler
	vehiclev1connect.UnimplementedVehicleServiceHandler
	speciesv1connect.UnimplementedSpeciesServiceHandler
	planetv1connect.UnimplementedPlanetServiceHandler
	relationsv1connect.UnimplementedFilmResolverServiceHandler
	relationsv1connect.UnimplementedPersonResolverServiceHandler
	relationsv1connect.UnimplementedPlanetResolverServiceHandler
	relationsv1connect.UnimplementedSpeciesResolverServiceHandler
	relationsv1connect.UnimplementedStarshipResolverServiceHandler
	relationsv1connect.UnimplementedVehicleResolverServiceHandler
}

// NewHandler returns a new handler that serves the Star Wars API.
func NewHandler() *Handler {
	return &Handler{}
}

// GetFilms implements the GetFilms RPC of the FilmService.
func (h *Handler) GetFilms(_ context.Context, req *connect.Request[filmv1.GetFilmsRequest]) (*connect.Response[filmv1.GetFilmsResponse], error) {
	films, err := getAll(req.Msg.Ids, allFilms)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&filmv1.GetFilmsResponse{
		Films: transform(films, transformFilm),
	}), nil
}

// ListFilms implements the ListFilms RPC of the FilmService.
func (h *Handler) ListFilms(_ context.Context, req *connect.Request[filmv1.ListFilmsRequest]) (*connect.Response[filmv1.ListFilmsResponse], error) {
	films, nextPageToken := paginate(allFilms, int(req.Msg.PageSize), req.Msg.PageToken)
	return connect.NewResponse(
		&filmv1.ListFilmsResponse{
			Films:         transform(films, transformFilm),
			NextPageToken: nextPageToken,
		},
	), nil
}

// GetPeople implements the GetPeople RPC of the PersonService.
func (h *Handler) GetPeople(_ context.Context, req *connect.Request[personv1.GetPeopleRequest]) (*connect.Response[personv1.GetPeopleResponse], error) {
	people, err := getAll(req.Msg.Ids, allPeople)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&personv1.GetPeopleResponse{
		People: transform(people, transformPerson),
	}), nil
}

// ListPeople implements the ListPeople RPC of the PersonService.
func (h *Handler) ListPeople(_ context.Context, req *connect.Request[personv1.ListPeopleRequest]) (*connect.Response[personv1.ListPeopleResponse], error) {
	people, nextPageToken := paginate(allPeople, int(req.Msg.PageSize), req.Msg.PageToken)
	return connect.NewResponse(
		&personv1.ListPeopleResponse{
			People:        transform(people, transformPerson),
			NextPageToken: nextPageToken,
		},
	), nil
}

// GetStarships implements the GetStarships RPC of the StarshipService.
func (h *Handler) GetStarships(_ context.Context, req *connect.Request[starshipv1.GetStarshipsRequest]) (*connect.Response[starshipv1.GetStarshipsResponse], error) {
	starships, err := getAll(req.Msg.Ids, allStarships)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&starshipv1.GetStarshipsResponse{
		Starships: transform(starships, transformStarship),
	}), nil
}

// ListStarships implements the ListStarships RPC of the StarshipService.
func (h *Handler) ListStarships(_ context.Context, req *connect.Request[starshipv1.ListStarshipsRequest]) (*connect.Response[starshipv1.ListStarshipsResponse], error) {
	starships, nextPageToken := paginate(allStarships, int(req.Msg.PageSize), req.Msg.PageToken)
	return connect.NewResponse(
		&starshipv1.ListStarshipsResponse{
			Starships:     transform(starships, transformStarship),
			NextPageToken: nextPageToken,
		},
	), nil
}

// GetVehicles implements the GetVehicles RPC of the VehicleService.
func (h *Handler) GetVehicles(_ context.Context, req *connect.Request[vehiclev1.GetVehiclesRequest]) (*connect.Response[vehiclev1.GetVehiclesResponse], error) {
	vehicles, err := getAll(req.Msg.Ids, allVehicles)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&vehiclev1.GetVehiclesResponse{
		Vehicles: transform(vehicles, transformVehicle),
	}), nil
}

// ListVehicles implements the ListVehicles RPC of the VehicleService.
func (h *Handler) ListVehicles(_ context.Context, req *connect.Request[vehiclev1.ListVehiclesRequest]) (*connect.Response[vehiclev1.ListVehiclesResponse], error) {
	vehicles, nextPageToken := paginate(allVehicles, int(req.Msg.PageSize), req.Msg.PageToken)
	return connect.NewResponse(
		&vehiclev1.ListVehiclesResponse{
			Vehicles:      transform(vehicles, transformVehicle),
			NextPageToken: nextPageToken,
		},
	), nil
}

// GetSpecies implements the GetSpecies RPC of the SpeciesService.
func (h *Handler) GetSpecies(_ context.Context, req *connect.Request[speciesv1.GetSpeciesRequest]) (*connect.Response[speciesv1.GetSpeciesResponse], error) {
	species, err := getAll(req.Msg.Ids, allSpecies)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&speciesv1.GetSpeciesResponse{
		Species: transform(species, transformSpecies),
	}), nil
}

// ListSpecies implements the ListSpecies RPC of the SpeciesService.
func (h *Handler) ListSpecies(_ context.Context, req *connect.Request[speciesv1.ListSpeciesRequest]) (*connect.Response[speciesv1.ListSpeciesResponse], error) {
	species, nextPageToken := paginate(allSpecies, int(req.Msg.PageSize), req.Msg.PageToken)
	return connect.NewResponse(
		&speciesv1.ListSpeciesResponse{
			Species:       transform(species, transformSpecies),
			NextPageToken: nextPageToken,
		},
	), nil
}

// GetPlanets implements the GetPlanets RPC of the PlanetService.
func (h *Handler) GetPlanets(_ context.Context, req *connect.Request[planetv1.GetPlanetsRequest]) (*connect.Response[planetv1.GetPlanetsResponse], error) {
	planets, err := getAll(req.Msg.Ids, allPlanets)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&planetv1.GetPlanetsResponse{
		Planets: transform(planets, transformPlanet),
	}), nil
}

// ListPlanets implements the ListPlanets RPC of the PlanetService.
func (h *Handler) ListPlanets(_ context.Context, req *connect.Request[planetv1.ListPlanetsRequest]) (*connect.Response[planetv1.ListPlanetsResponse], error) {
	planets, nextPageToken := paginate(allPlanets, int(req.Msg.PageSize), req.Msg.PageToken)
	return connect.NewResponse(
		&planetv1.ListPlanetsResponse{
			Planets:       transform(planets, transformPlanet),
			NextPageToken: nextPageToken,
		},
	), nil
}

// GetFilmCharacters implements the GetFilmCharacters RPC of the PersonResolverService.
func (h *Handler) GetFilmCharacters(ctx context.Context, req *connect.Request[relationsv1.GetFilmRelationsRequest]) (*connect.Response[relationsv1.GetCharactersResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(film *filmv1.Film) []string { return film.CharacterIds },
		func(ctx context.Context, ids []string) (*connect.Response[personv1.GetPeopleResponse], error) {
			return h.GetPeople(ctx, connect.NewRequest(&personv1.GetPeopleRequest{Ids: ids}))
		},
		func(msg *personv1.GetPeopleResponse) []*personv1.Person { return msg.People },
		func(values []*personv1.Person, result *relationsv1.GetCharactersResponse_Result) {
			result.Characters = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetCharactersResponse{Values: wrappers}), nil
}

// GetFilmPlanets implements the GetFilmPlanets RPC of the PlanetResolverService.
func (h *Handler) GetFilmPlanets(ctx context.Context, req *connect.Request[relationsv1.GetFilmRelationsRequest]) (*connect.Response[relationsv1.GetPlanetsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(film *filmv1.Film) []string { return film.PlanetIds },
		func(ctx context.Context, ids []string) (*connect.Response[planetv1.GetPlanetsResponse], error) {
			return h.GetPlanets(ctx, connect.NewRequest(&planetv1.GetPlanetsRequest{Ids: ids}))
		},
		func(msg *planetv1.GetPlanetsResponse) []*planetv1.Planet { return msg.Planets },
		func(values []*planetv1.Planet, result *relationsv1.GetPlanetsResponse_Result) {
			result.Planets = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetPlanetsResponse{Values: wrappers}), nil
}

// GetFilmSpecies implements the GetFilmSpecies RPC of the SpeciesResolverService.
func (h *Handler) GetFilmSpecies(ctx context.Context, req *connect.Request[relationsv1.GetFilmRelationsRequest]) (*connect.Response[relationsv1.GetSpeciesResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(film *filmv1.Film) []string { return film.SpeciesIds },
		func(ctx context.Context, ids []string) (*connect.Response[speciesv1.GetSpeciesResponse], error) {
			return h.GetSpecies(ctx, connect.NewRequest(&speciesv1.GetSpeciesRequest{Ids: ids}))
		},
		func(msg *speciesv1.GetSpeciesResponse) []*speciesv1.Species { return msg.Species },
		func(values []*speciesv1.Species, result *relationsv1.GetSpeciesResponse_Result) {
			result.Species = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetSpeciesResponse{Values: wrappers}), nil
}

// GetFilmStarships implements the GetFilmStarships RPC of the StarshipResolverService.
func (h *Handler) GetFilmStarships(ctx context.Context, req *connect.Request[relationsv1.GetFilmRelationsRequest]) (*connect.Response[relationsv1.GetStarshipsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(film *filmv1.Film) []string { return film.StarshipIds },
		func(ctx context.Context, ids []string) (*connect.Response[starshipv1.GetStarshipsResponse], error) {
			return h.GetStarships(ctx, connect.NewRequest(&starshipv1.GetStarshipsRequest{Ids: ids}))
		},
		func(msg *starshipv1.GetStarshipsResponse) []*starshipv1.Starship { return msg.Starships },
		func(values []*starshipv1.Starship, result *relationsv1.GetStarshipsResponse_Result) {
			result.Starships = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetStarshipsResponse{Values: wrappers}), nil
}

// GetFilmVehicles implements the GetFilmVehicles RPC of the VehicleResolverService.
func (h *Handler) GetFilmVehicles(ctx context.Context, req *connect.Request[relationsv1.GetFilmRelationsRequest]) (*connect.Response[relationsv1.GetVehiclesResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(film *filmv1.Film) []string { return film.VehicleIds },
		func(ctx context.Context, ids []string) (*connect.Response[vehiclev1.GetVehiclesResponse], error) {
			return h.GetVehicles(ctx, connect.NewRequest(&vehiclev1.GetVehiclesRequest{Ids: ids}))
		},
		func(msg *vehiclev1.GetVehiclesResponse) []*vehiclev1.Vehicle { return msg.Vehicles },
		func(values []*vehiclev1.Vehicle, result *relationsv1.GetVehiclesResponse_Result) {
			result.Vehicles = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetVehiclesResponse{Values: wrappers}), nil
}

// GetPersonHomeworld implements the GetPersonHomeworld RPC of the PlanetResolverService.
func (h *Handler) GetPersonHomeworld(ctx context.Context, req *connect.Request[relationsv1.GetPersonRelationRequest]) (*connect.Response[relationsv1.GetHomeworldResponse], error) {
	wrappers, err := resolve1to1Batch(
		ctx,
		req.Msg.Bases,
		func(person *personv1.Person) string { return person.HomeworldId },
		func(ctx context.Context, ids []string) (*connect.Response[planetv1.GetPlanetsResponse], error) {
			return h.GetPlanets(ctx, connect.NewRequest(&planetv1.GetPlanetsRequest{Ids: ids}))
		},
		func(msg *planetv1.GetPlanetsResponse) []*planetv1.Planet { return msg.Planets },
		func(value *planetv1.Planet, result *relationsv1.GetHomeworldResponse_Result) {
			result.Homeworld = value
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetHomeworldResponse{Values: wrappers}), nil
}

// GetPersonFilms implements the GetPersonFilms RPC of the FilmResolverService.
func (h *Handler) GetPersonFilms(ctx context.Context, req *connect.Request[relationsv1.GetPersonRelationsRequest]) (*connect.Response[relationsv1.GetFilmsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(person *personv1.Person) []string { return person.FilmIds },
		func(ctx context.Context, ids []string) (*connect.Response[filmv1.GetFilmsResponse], error) {
			return h.GetFilms(ctx, connect.NewRequest(&filmv1.GetFilmsRequest{Ids: ids}))
		},
		func(msg *filmv1.GetFilmsResponse) []*filmv1.Film { return msg.Films },
		func(values []*filmv1.Film, result *relationsv1.GetFilmsResponse_Result) { result.Films = values },
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetFilmsResponse{Values: wrappers}), nil
}

// GetPersonSpecies implements the GetPersonSpecies RPC of the SpeciesResolverService.
func (h *Handler) GetPersonSpecies(ctx context.Context, req *connect.Request[relationsv1.GetPersonRelationsRequest]) (*connect.Response[relationsv1.GetSpeciesResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(person *personv1.Person) []string { return person.SpeciesIds },
		func(ctx context.Context, ids []string) (*connect.Response[speciesv1.GetSpeciesResponse], error) {
			return h.GetSpecies(ctx, connect.NewRequest(&speciesv1.GetSpeciesRequest{Ids: ids}))
		},
		func(msg *speciesv1.GetSpeciesResponse) []*speciesv1.Species { return msg.Species },
		func(values []*speciesv1.Species, result *relationsv1.GetSpeciesResponse_Result) {
			result.Species = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetSpeciesResponse{Values: wrappers}), nil
}

// GetPersonStarships implements the GetPersonStarships RPC of the StarshipResolverService.
func (h *Handler) GetPersonStarships(ctx context.Context, req *connect.Request[relationsv1.GetPersonRelationsRequest]) (*connect.Response[relationsv1.GetStarshipsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(person *personv1.Person) []string { return person.StarshipIds },
		func(ctx context.Context, ids []string) (*connect.Response[starshipv1.GetStarshipsResponse], error) {
			return h.GetStarships(ctx, connect.NewRequest(&starshipv1.GetStarshipsRequest{Ids: ids}))
		},
		func(msg *starshipv1.GetStarshipsResponse) []*starshipv1.Starship { return msg.Starships },
		func(values []*starshipv1.Starship, result *relationsv1.GetStarshipsResponse_Result) {
			result.Starships = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetStarshipsResponse{Values: wrappers}), nil
}

// GetPersonVehicles implements the GetPersonVehicles RPC of the VehicleResolverService.
func (h *Handler) GetPersonVehicles(ctx context.Context, req *connect.Request[relationsv1.GetPersonRelationsRequest]) (*connect.Response[relationsv1.GetVehiclesResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(person *personv1.Person) []string { return person.VehicleIds },
		func(ctx context.Context, ids []string) (*connect.Response[vehiclev1.GetVehiclesResponse], error) {
			return h.GetVehicles(ctx, connect.NewRequest(&vehiclev1.GetVehiclesRequest{Ids: ids}))
		},
		func(msg *vehiclev1.GetVehiclesResponse) []*vehiclev1.Vehicle { return msg.Vehicles },
		func(values []*vehiclev1.Vehicle, result *relationsv1.GetVehiclesResponse_Result) {
			result.Vehicles = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetVehiclesResponse{Values: wrappers}), nil
}

// GetPlanetFilms implements the GetPlanetFilms RPC of the FilmResolverService.
func (h *Handler) GetPlanetFilms(ctx context.Context, req *connect.Request[relationsv1.GetPlanetRelationsRequest]) (*connect.Response[relationsv1.GetFilmsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(planet *planetv1.Planet) []string { return planet.FilmIds },
		func(ctx context.Context, ids []string) (*connect.Response[filmv1.GetFilmsResponse], error) {
			return h.GetFilms(ctx, connect.NewRequest(&filmv1.GetFilmsRequest{Ids: ids}))
		},
		func(msg *filmv1.GetFilmsResponse) []*filmv1.Film { return msg.Films },
		func(values []*filmv1.Film, result *relationsv1.GetFilmsResponse_Result) { result.Films = values },
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetFilmsResponse{Values: wrappers}), nil
}

// GetPlanetResidents implements the GetPlanetResidents RPC of the PersonResolverService.
func (h *Handler) GetPlanetResidents(ctx context.Context, req *connect.Request[relationsv1.GetPlanetRelationsRequest]) (*connect.Response[relationsv1.GetResidentsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(planet *planetv1.Planet) []string { return planet.ResidentIds },
		func(ctx context.Context, ids []string) (*connect.Response[personv1.GetPeopleResponse], error) {
			return h.GetPeople(ctx, connect.NewRequest(&personv1.GetPeopleRequest{Ids: ids}))
		},
		func(msg *personv1.GetPeopleResponse) []*personv1.Person { return msg.People },
		func(values []*personv1.Person, result *relationsv1.GetResidentsResponse_Result) {
			result.Residents = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetResidentsResponse{Values: wrappers}), nil
}

// GetSpeciesHomeworld implements the GetSpeciesHomeworld RPC of the PlanetResolverService.
func (h *Handler) GetSpeciesHomeworld(ctx context.Context, req *connect.Request[relationsv1.GetSpeciesRelationRequest]) (*connect.Response[relationsv1.GetHomeworldResponse], error) {
	wrappers, err := resolve1to1Batch(
		ctx,
		req.Msg.Bases,
		func(species *speciesv1.Species) string { return species.HomeworldId },
		func(ctx context.Context, ids []string) (*connect.Response[planetv1.GetPlanetsResponse], error) {
			return h.GetPlanets(ctx, connect.NewRequest(&planetv1.GetPlanetsRequest{Ids: ids}))
		},
		func(msg *planetv1.GetPlanetsResponse) []*planetv1.Planet { return msg.Planets },
		func(value *planetv1.Planet, result *relationsv1.GetHomeworldResponse_Result) {
			result.Homeworld = value
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetHomeworldResponse{Values: wrappers}), nil
}

// GetSpeciesFilms implements the GetSpeciesFilms RPC of the FilmResolverService.
func (h *Handler) GetSpeciesFilms(ctx context.Context, req *connect.Request[relationsv1.GetSpeciesRelationsRequest]) (*connect.Response[relationsv1.GetFilmsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(species *speciesv1.Species) []string { return species.FilmIds },
		func(ctx context.Context, ids []string) (*connect.Response[filmv1.GetFilmsResponse], error) {
			return h.GetFilms(ctx, connect.NewRequest(&filmv1.GetFilmsRequest{Ids: ids}))
		},
		func(msg *filmv1.GetFilmsResponse) []*filmv1.Film { return msg.Films },
		func(values []*filmv1.Film, result *relationsv1.GetFilmsResponse_Result) { result.Films = values },
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetFilmsResponse{Values: wrappers}), nil
}

// GetSpeciesCharacters implements the GetSpeciesCharacters RPC of the PersonResolverService.
func (h *Handler) GetSpeciesCharacters(ctx context.Context, req *connect.Request[relationsv1.GetSpeciesRelationsRequest]) (*connect.Response[relationsv1.GetCharactersResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(species *speciesv1.Species) []string { return species.PeopleIds },
		func(ctx context.Context, ids []string) (*connect.Response[personv1.GetPeopleResponse], error) {
			return h.GetPeople(ctx, connect.NewRequest(&personv1.GetPeopleRequest{Ids: ids}))
		},
		func(msg *personv1.GetPeopleResponse) []*personv1.Person { return msg.People },
		func(values []*personv1.Person, result *relationsv1.GetCharactersResponse_Result) {
			result.Characters = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetCharactersResponse{Values: wrappers}), nil
}

// GetStarshipFilms implements the GetStarshipFilms RPC of the FilmResolverService.
func (h *Handler) GetStarshipFilms(ctx context.Context, req *connect.Request[relationsv1.GetStarshipRelationsRequest]) (*connect.Response[relationsv1.GetFilmsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(starship *starshipv1.Starship) []string { return starship.FilmIds },
		func(ctx context.Context, ids []string) (*connect.Response[filmv1.GetFilmsResponse], error) {
			return h.GetFilms(ctx, connect.NewRequest(&filmv1.GetFilmsRequest{Ids: ids}))
		},
		func(msg *filmv1.GetFilmsResponse) []*filmv1.Film { return msg.Films },
		func(values []*filmv1.Film, result *relationsv1.GetFilmsResponse_Result) { result.Films = values },
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetFilmsResponse{Values: wrappers}), nil
}

// GetStarshipPilots implements the GetStarshipPilots RPC of the PersonResolverService.
func (h *Handler) GetStarshipPilots(ctx context.Context, req *connect.Request[relationsv1.GetStarshipRelationsRequest]) (*connect.Response[relationsv1.GetPilotsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(starship *starshipv1.Starship) []string { return starship.PilotIds },
		func(ctx context.Context, ids []string) (*connect.Response[personv1.GetPeopleResponse], error) {
			return h.GetPeople(ctx, connect.NewRequest(&personv1.GetPeopleRequest{Ids: ids}))
		},
		func(msg *personv1.GetPeopleResponse) []*personv1.Person { return msg.People },
		func(values []*personv1.Person, result *relationsv1.GetPilotsResponse_Result) {
			result.Pilots = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetPilotsResponse{Values: wrappers}), nil
}

// GetVehicleFilms implements the GetVehicleFilms RPC of the FilmResolverService.
func (h *Handler) GetVehicleFilms(ctx context.Context, req *connect.Request[relationsv1.GetVehicleRelationsRequest]) (*connect.Response[relationsv1.GetFilmsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(vehicle *vehiclev1.Vehicle) []string { return vehicle.FilmIds },
		func(ctx context.Context, ids []string) (*connect.Response[filmv1.GetFilmsResponse], error) {
			return h.GetFilms(ctx, connect.NewRequest(&filmv1.GetFilmsRequest{Ids: ids}))
		},
		func(msg *filmv1.GetFilmsResponse) []*filmv1.Film { return msg.Films },
		func(values []*filmv1.Film, result *relationsv1.GetFilmsResponse_Result) { result.Films = values },
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetFilmsResponse{Values: wrappers}), nil
}

// GetVehiclePilots implements the GetVehiclePilots RPC of the PersonResolverService.
func (h *Handler) GetVehiclePilots(ctx context.Context, req *connect.Request[relationsv1.GetVehicleRelationsRequest]) (*connect.Response[relationsv1.GetPilotsResponse], error) {
	wrappers, err := resolveBatch(
		ctx,
		req.Msg.Bases,
		int(req.Msg.Limit),
		func(vehicle *vehiclev1.Vehicle) []string { return vehicle.PilotIds },
		func(ctx context.Context, ids []string) (*connect.Response[personv1.GetPeopleResponse], error) {
			return h.GetPeople(ctx, connect.NewRequest(&personv1.GetPeopleRequest{Ids: ids}))
		},
		func(msg *personv1.GetPeopleResponse) []*personv1.Person { return msg.People },
		func(values []*personv1.Person, result *relationsv1.GetPilotsResponse_Result) {
			result.Pilots = values
		},
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&relationsv1.GetPilotsResponse{Values: wrappers}), nil
}

func getAll[T any](ids []string, items []*T) ([]*T, error) {
	idIndex := make(map[string]int, len(ids))
	for index, id := range ids {
		idIndex[id] = index
	}
	results := make([]*T, len(ids))
	for _, item := range items {
		url := reflect.ValueOf(item).Elem().FieldByName("URL").String()
		id := urlToID(url)
		index, ok := idIndex[id]
		if !ok {
			continue
		}
		results[index] = item
	}
	var missingIDs []string
	for i, item := range results {
		if item == nil {
			missingIDs = append(missingIDs, ids[i])
		}
	}
	if len(missingIDs) > 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown IDs: %v", missingIDs))
	}
	return results, nil
}

func paginate[T any](entities []T, pageSize int, pageToken string) ([]T, string) {
	start := 0
	if pageToken != "" {
		token, err := strconv.Atoi(pageToken)
		if err == nil {
			start = token
		}
	}
	var nextPageToken string
	if pageSize <= 0 || pageSize > len(entities) {
		pageSize = len(entities)
	}
	if len(entities[start:]) != pageSize {
		nextPageToken = strconv.Itoa(len(entities[:start+pageSize]))
	}
	entities = entities[start : start+pageSize]
	return entities, nextPageToken
}

func transform[T, V any](entities []T, transformFn func(T) V) []V {
	results := make([]V, 0, len(entities))
	for _, v := range entities {
		results = append(results, transformFn(v))
	}
	return results
}

func urlToID(url string) string {
	splits := strings.Split(strings.Trim(url, "/"), "/")
	return splits[len(splits)-1]
}

func transformFilm(film *swapi.Film) *filmv1.Film {
	return &filmv1.Film{
		Id:            urlToID(film.URL),
		Title:         film.Title,
		EpisodeNumber: int32(film.EpisodeID),
		OpeningCrawl:  film.OpeningCrawl,
		Director:      film.Director,
		Producer:      film.Producer,
		CharacterIds:  transform(film.CharacterURLs, urlToID),
		PlanetIds:     transform(film.PlanetURLs, urlToID),
		SpeciesIds:    transform(film.SpeciesURLs, urlToID),
		StarshipIds:   transform(film.StarshipURLs, urlToID),
		VehicleIds:    transform(film.VehicleURLs, urlToID),
		Created:       mustTimestamp(film.Created),
		Edited:        mustTimestamp(film.Edited),
	}
}

func transformPerson(person *swapi.Person) *personv1.Person {
	return &personv1.Person{
		Id:          urlToID(person.URL),
		Name:        person.Name,
		Mass:        person.Mass,
		HairColor:   person.HairColor,
		SkinColor:   person.SkinColor,
		EyeColor:    person.EyeColor,
		BirthYear:   person.BirthYear,
		Gender:      person.Gender,
		HomeworldId: urlToID(person.Homeworld),
		FilmIds:     transform(person.FilmURLs, urlToID),
		SpeciesIds:  transform(person.SpeciesURLs, urlToID),
		StarshipIds: transform(person.StarshipURLs, urlToID),
		VehicleIds:  transform(person.VehicleURLs, urlToID),
		Created:     mustTimestamp(person.Created),
		Edited:      mustTimestamp(person.Edited),
	}
}

func transformStarship(starship *swapi.Starship) *starshipv1.Starship {
	return &starshipv1.Starship{
		Id:                   urlToID(starship.URL),
		Name:                 starship.Name,
		Mglt:                 mustInt(starship.MGLT),
		CargoCapacity:        mustFloat(starship.CargoCapacity),
		Consumable:           starship.Consumables,
		CostInCredits:        maybeInt(starship.CostInCredits),
		Crew:                 mustInt(starship.Crew),
		HyperDriveRating:     mustFloat(starship.HyperdriveRating),
		Length:               mustFloat(starship.Length),
		Manufacturers:        transform(strings.Split(starship.Manufacturer, ","), strings.TrimSpace),
		MaxAtmospheringSpeed: maybeInt(starship.MaxAtmospheringSpeed),
		Model:                starship.Model,
		Passengers:           mustInt(starship.Passengers),
		Class:                starship.StarshipClass,
		PilotIds:             transform(starship.PilotURLs, urlToID),
		FilmIds:              transform(starship.FilmURLs, urlToID),
		Created:              mustTimestamp(starship.Created),
		Edited:               mustTimestamp(starship.Edited),
	}
}

func transformVehicle(vehicle *swapi.Vehicle) *vehiclev1.Vehicle {
	return &vehiclev1.Vehicle{
		Id:            urlToID(vehicle.URL),
		Name:          vehicle.Name,
		CargoCapacity: mustFloat(vehicle.CargoCapacity),
		CostInCredits: maybeInt(vehicle.CostInCredits),
		Crew:          mustInt(vehicle.Crew),
		Length:        mustFloat(vehicle.Length),
		Manufacturers: transform(strings.Split(vehicle.Manufacturer, ","), strings.TrimSpace),
		Model:         vehicle.Model,
		Passengers:    mustInt(vehicle.Passengers),
		Class:         vehicle.VehicleClass,
		PilotIds:      transform(vehicle.PilotURLs, urlToID),
		FilmIds:       transform(vehicle.FilmURLs, urlToID),
		Created:       mustTimestamp(vehicle.Created),
		Edited:        mustTimestamp(vehicle.Edited),
	}
}

func transformSpecies(species *swapi.Species) *speciesv1.Species {
	return &speciesv1.Species{
		Id:              urlToID(species.URL),
		AverageHeight:   mustFloat(species.AverageHeight),
		AverageLifespan: mustInt(species.AverageLifespan),
		Classification:  species.Classification,
		Designation:     species.Designation,
		EyeColors:       transform(strings.Split(species.EyeColors, ","), strings.TrimSpace),
		HairColors:      transform(strings.Split(species.HairColors, ","), strings.TrimSpace),
		Language:        species.Language,
		Name:            species.Name,
		HomeworldId:     urlToID(species.Homeworld),
		SkinColors:      transform(strings.Split(species.SkinColors, ","), strings.TrimSpace),
		PeopleIds:       transform(species.PeopleURLs, urlToID),
		FilmIds:         transform(species.FilmURLs, urlToID),
		Created:         mustTimestamp(species.Created),
		Edited:          mustTimestamp(species.Edited),
	}
}

func transformPlanet(planet *swapi.Planet) *planetv1.Planet {
	return &planetv1.Planet{
		Id:             urlToID(planet.URL),
		Climates:       []string{planet.Climate},
		Diameter:       int32(mustInt(planet.Diameter)),
		Gravity:        planet.Gravity,
		Name:           planet.Name,
		OrbitalPeriod:  int32(mustInt(planet.OrbitalPeriod)),
		Population:     mustFloat(planet.Population),
		RotationPeriod: int32(mustInt(planet.RotationPeriod)),
		SurfaceWater:   mustFloat(planet.SurfaceWater),
		Terrains:       []string{planet.Terrain},
		ResidentIds:    transform(planet.ResidentURLs, urlToID),
		FilmIds:        transform(planet.FilmURLs, urlToID),
		Created:        mustTimestamp(planet.Created),
		Edited:         mustTimestamp(planet.Edited),
	}
}

func mustInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func maybeInt(s string) *int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}

func mustFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func mustTimestamp(s string) *timestamppb.Timestamp {
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return timestamppb.New(v)
}

func resolveBatch[E, R, M, W any](
	ctx context.Context,
	entities []E,
	limit int,
	idExtractor func(E) []string,
	invoker func(context.Context, []string) (*connect.Response[M], error),
	resultExtractor func(*M) []R,
	resultStorer func([]R, *W),
) ([]*W, error) {
	idSet := map[string]struct{}{}
	idBatches := make([][]string, len(entities))
	for i, entity := range entities {
		ids := idExtractor(entity)
		if limit > 0 && len(ids) > limit {
			ids = ids[:limit]
		}
		for _, item := range ids {
			idSet[item] = struct{}{}
		}
		idBatches[i] = ids
	}
	idSlice := make([]string, 0, len(idSet))
	for item := range idSet {
		idSlice = append(idSlice, item)
	}

	resp, err := invoker(ctx, idSlice)
	if err != nil {
		return nil, err
	}
	results := resultExtractor(resp.Msg)

	indices := map[string]int{}
	for i, item := range idSlice {
		indices[item] = i
	}
	batchedResults := make([]*W, len(entities))
	for i := range entities {
		ids := idBatches[i]
		batch := make([]R, len(ids))
		for j, item := range ids {
			batch[j] = results[indices[item]]
		}
		var w W
		resultStorer(batch, &w)
		batchedResults[i] = &w
	}
	return batchedResults, nil
}

func resolve1to1Batch[E, R, M, W any](
	ctx context.Context,
	entities []E,
	idExtractor func(E) string,
	invoker func(context.Context, []string) (*connect.Response[M], error),
	resultExtractor func(*M) []R,
	resultStorer func(R, *W),
) ([]*W, error) {
	return resolveBatch(
		ctx,
		entities,
		0,
		func(e E) []string {
			id := idExtractor(e)
			if id == "" {
				return nil
			}
			return []string{id}
		},
		invoker,
		resultExtractor,
		func(r []R, w *W) {
			if len(r) > 0 {
				resultStorer(r[0], w)
			}
		},
	)
}
