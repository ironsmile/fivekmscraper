package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ironsmile/logger"
)

type namesArchive struct {
	Male   []string `json:"male"`
	Female []string `json:"female"`
}

func loadNamesDicts(filename string) (map[string]struct{}, map[string]struct{}) {
	male := make(map[string]struct{})
	female := make(map[string]struct{})

	file, err := os.Open(filename)
	if err != nil {
		logger.Errorf("Could not open names archive: %s\n", err)
		return male, female
	}

	defer file.Close()

	var arch namesArchive

	dec := json.NewDecoder(file)
	if err := dec.Decode(&arch); err != nil {
		logger.Errorf("Could not decode names archive: %s\n", err)
		return male, female
	}

	for _, name := range arch.Male {
		name = strings.ToLower(name)
		male[name] = struct{}{}
	}

	for _, name := range arch.Female {
		name = strings.ToLower(name)
		female[name] = struct{}{}
	}

	return male, female
}

func saveNemesDicts(filename string, males, females map[string]struct{}) {
	file, err := os.Create(filename)
	if err != nil {
		logger.Errorf("Could not archive names to disk: %s\n", err)
		return
	}

	defer file.Close()

	enc := json.NewEncoder(file)

	var arch namesArchive

	for k := range males {
		k = strings.ToLower(k)
		arch.Male = append(arch.Male, k)
	}

	for k := range females {
		k = strings.ToLower(k)
		arch.Female = append(arch.Female, k)
	}

	if err := enc.Encode(arch); err != nil {
		logger.Errorf("Could not encode archive: %s\n", err)
	}
}
