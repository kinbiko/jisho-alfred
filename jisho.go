package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
)

type app struct {
	wf *aw.Workflow
}

func main() {
	a := &app{wf: aw.New()}
	a.wf.Run(a.run)
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

// should take the format of:
// [ {data.slug}: {data.japanese[].reading} -- {data.senses[].english_definitions[]} ]
type resp struct {
	Data []*Datum `json:"data"`
}

type entry struct {
	slug     string
	readings []string
	defns    []string
}

func valid(query string) bool {
	sanitized := strings.TrimSpace(query)
	for _, s := range []string{"", `"`, "'"} {
		if s == sanitized {
			return false
		}
	}
	return true
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
	log.Println(string(bodyBytes))

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

// Must be created in Alfred as:
// Script filter
// Language: External Script
// Script file: ~/go/bin/jisho
// Don't tick "Alfred filters results"
//
// There's probably a better way of doing this...
func (a *app) run() {
	query := strings.Join(os.Args[1:], " ")
	if !valid(query) {
		it := a.wf.NewItem(fmt.Sprintf("No match for '%s'", query))
		it.Arg(query)
		it.Valid(true)
		it.Icon(aw.IconNote)
		a.wf.SendFeedback()
		return
	}

	for _, entry := range getResults(query) {
		it := a.wf.NewItem(fmt.Sprintf("%s (%s)", entry.slug, strings.Join(entry.readings, " ")))
		it.Arg(query)
		it.Subtitle(strings.Join(entry.defns, ", "))
		it.Icon(aw.IconNote)
	}

	a.wf.SendFeedback()
}
