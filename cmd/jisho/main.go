package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("requires exactly 1 parameter")
	}
	entries := getResults(os.Args[1])
	if len(entries) == 0 {
		log.Fatalln("0 results for", os.Args[1])
	}
	fmt.Printf("\"%s(%s)\",%s\n", entries[0].readings[0], entries[0].defns[0], entries[0].slug)
}

type entry struct {
	slug     string
	readings []string
	defns    []string
}

// should take the format of:
// [ {data.slug}: {data.japanese[].reading} -- {data.senses[].english_definitions[]} ]
type resp struct {
	Data []*Datum `json:"data"`
}

func getResults(query string) []*entry {
	res, err := http.Get("http://jisho.org/api/v1/search/words?keyword=" + url.QueryEscape(query))
	if err != nil {
		return nil
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	if res.StatusCode != 200 {
		return nil
	}

	r := &resp{}
	_ = json.Unmarshal(bodyBytes, r)

	entries := []*entry{}

	for _, datum := range r.Data {
		entries = append(entries, &entry{
			slug:     datum.Slug,
			readings: extractReadings(datum),
			defns:    extractDefinitions(datum),
		})
	}
	return entries
}

func extractReadings(datum *Datum) []string {
	readings := []string{}
	for _, r := range datum.Japanese {
		readings = appendIfNew(readings, r.Reading)
	}
	return readings
}

func appendIfNew(existing []string, cand string) []string {
	for _, s := range existing {
		if cand == s {
			return existing
		}
	}
	return append(existing, cand)
}

func extractDefinitions(datum *Datum) []string {
	defns := []string{}
	for _, s := range datum.Senses {
		for _, def := range s.EnglishDefinitions {
			defns = appendIfNew(defns, def)
		}
	}
	return defns
}

type Datum struct {
	Slug     string `json:"slug"`
	Japanese []struct {
		Reading string `json:"reading"`
	} `json:"japanese"`
	Senses []struct {
		EnglishDefinitions []string `json:"english_definitions"`
	} `json:"senses"`
}
