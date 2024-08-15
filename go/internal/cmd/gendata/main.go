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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/peterhellberg/swapi"
)

const (
	filmsURL     = "https://swapi.dev/api/films"
	peopleURL    = "https://swapi.dev/api/people"
	planetsURL   = "https://swapi.dev/api/planets"
	speciesURL   = "https://swapi.dev/api/species"
	starshipsURL = "https://swapi.dev/api/starships"
	vehiclesURL  = "https://swapi.dev/api/vehicles"
)

func main() {
	var pkgName string
	var outputName string
	switch len(os.Args) {
	case 3:
		outputName = os.Args[2]
		fallthrough
	case 2:
		pkgName = os.Args[1]
	default:
		log.Fatalln("usage: gendata package [output_file]")
	}

	if err := generateSwapiData(context.Background(), pkgName, outputName); err != nil {
		log.Fatalln(err)
	}
}

func generateSwapiData(ctx context.Context, pkgName, outputName string) error {
	if outputName != "" {
		out, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := out.Close(); err != nil {
				log.Printf("failed to close output file %q: %v\n", outputName, err)
			}
		}()
		os.Stdout = out
	}

	log.Println("Retrieving films...")
	allFilms, err := readAll[*swapi.Film](ctx, filmsURL)
	if err != nil {
		return err
	}
	log.Println("Retrieving people...")
	allPeople, err := readAll[*swapi.Person](ctx, peopleURL)
	if err != nil {
		return err
	}
	log.Println("Retrieving planets...")
	allPlanets, err := readAll[*swapi.Planet](ctx, planetsURL)
	if err != nil {
		return err
	}
	log.Println("Retrieving species...")
	allSpecies, err := readAll[*swapi.Species](ctx, speciesURL)
	if err != nil {
		return err
	}
	log.Println("Retrieving starships...")
	allStarships, err := readAll[*swapi.Starship](ctx, starshipsURL)
	if err != nil {
		return err
	}
	log.Println("Retrieving vehicles...")
	allVehicles, err := readAll[*swapi.Vehicle](ctx, vehiclesURL)
	if err != nil {
		return err
	}

	// Now print out generated file
	fmt.Println("// THIS IS GENERATED CODE. DO NOT EDIT.")
	fmt.Println("package", pkgName)
	fmt.Println("")
	fmt.Println(`import "github.com/peterhellberg/swapi"`)
	fmt.Println("")
	fmt.Println("var (")
	printItems("allFilms", allFilms)
	printItems("allPeople", allPeople)
	printItems("allPlanets", allPlanets)
	printItems("allSpecies", allSpecies)
	printItems("allStarships", allStarships)
	printItems("allVehicles", allVehicles)
	fmt.Println(")")
	return nil
}

func printItems[T any](varName string, items []*T) {
	var t T
	fmt.Printf("\t%s = []*swapi.%s{\n", varName, reflect.TypeOf(t).Name())
	for _, item := range items {
		fmt.Println("\t\t{")
		v := reflect.ValueOf(item).Elem()
		vt := v.Type()
		for i := range v.NumField() {
			fmt.Printf("\t\t\t%s: ", vt.Field(i).Name)
			printItem(v.Field(i).Interface())
			fmt.Println(",")
		}
		fmt.Println("\t\t},")
	}
	fmt.Println("\t}")
}

func printItem(item any) {
	// NB: We only need to support ints, strings, and string slices for now
	//     since those arethe only actual data types in the swapi structs.
	switch item := item.(type) {
	case int:
		fmt.Printf("%d", item)
	case string:
		fmt.Printf("%q", item)
	case []string:
		fmt.Print("[]string{")
		for i, s := range item {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%q", s)
		}
		fmt.Print("}")
	default:
		panic(fmt.Sprintf("unexpected value type %T", item))
	}
}

type pageResult[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

func readAll[T any](ctx context.Context, start string) ([]T, error) {
	var results []T
	next := start
	for {
		var page pageResult[T]
		if err := getJSON(ctx, next, &page); err != nil {
			return nil, err
		}
		results = append(results, page.Results...)
		if page.Next == nil {
			break
		}
		next = *page.Next
	}
	return results, nil
}

func getJSON(ctx context.Context, url string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unknow status code: %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(&dest)
}
