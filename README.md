# Jisho-Alfred

A program to be used in an [Alfred workflow](https://www.alfredapp.com/) in order to look up [Japanese words/translations with Jisho](http://jisho.org/).

![Demo of the workflow](./demo.gif)

## Install

1. If you don't already have `go` installed, install it with `brew install go`.
1. Run `go install github.com/kinbiko/jisho-alfred@latest`.
1. Grab the `jisho.alfredworkflow` from this repository and open it with Alfred to install.

## Usage

`ji <your search here` or `じ<your search here>` (no space) will first list all results within Alfred, letting you quickly look up pronunciations, kanji, and English meanings.
You can then hit enter to open the selected result in Jisho.org.
