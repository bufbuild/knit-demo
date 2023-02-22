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

import { createClient, Query } from "@bufbuild/knit";
import type { FilmService } from "./gen/buf/knit/demo/swapi/film/v1/film_knit.js";
import type { PersonService } from "./gen/buf/knit/demo/swapi/person/v1/person_knit.js";
import type { PlanetService } from "./gen/buf/knit/demo/swapi/planet/v1/planet_knit.js";
import type { SpeciesService } from "./gen/buf/knit/demo/swapi/species/v1/species_knit.js";
import type { StarshipService } from "./gen/buf/knit/demo/swapi/starship/v1/starship_knit.js";
import type {
  VehicleService,
  Vehicle,
} from "./gen/buf/knit/demo/swapi/vehicle/v1/vehicle_knit.js";
import "./gen/buf/knit/demo/swapi/relations/v1/relations_knit.js";

type Schema = FilmService &
  PersonService &
  PlanetService &
  SpeciesService &
  StarshipService &
  VehicleService;

const client = createClient<Schema>({
  baseUrl: "http://127.0.0.1:30480/",
  credentials: "include",
});

const vehicleQuery = {
  name: {},
  class: {},
  manufacturers: {},
  model: {},
} satisfies Query<Vehicle>;

const res = await client.fetch({
  "buf.knit.demo.swapi.film.v1.FilmService": {
    getFilms: {
      $: { ids: ["1"] },
      films: {
        title: {},
        episodeNumber: {},
        director: {},
        releaseDate: {},
        characters: {
          $: { limit: 20 },
          name: {},
          birthYear: {},
          species: {
            $: { limit: 2 },
            name: {},
            classification: {},
            homeworld: {
              name: {},
            },
          },
          homeworld: {},
          vehicles: {
            $: { limit: 5 },
            ...vehicleQuery,
          },
          starships: {
            $: { limit: 2 },
            pilots: {
              $: { limit: 2 },
              name: {},
              species: {
                $: { limit: 1 },
                name: {},
              },
              homeworld: {
                name: {},
              },
            },
          },
        },
      },
    },
  },
});

console.log(JSON.stringify(res, null, 2));
