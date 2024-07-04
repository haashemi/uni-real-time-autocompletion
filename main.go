package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Max options to show on a match
const OptionsLimit int = 12

// wordsFile is the embedded 'words.txt' file at the compile-time.
//
//go:embed words.txt
var wordsFile []byte

// words holds the words loaded from the 'wordsFile' variable.
var words []string

// initialize the global variables.
func init() {
	// split the wordsBytes by \n to the words variable.
	words = strings.Split(string(wordsFile), "\n")
}

func main() {
	// An infinite loop.
	for {
		// Ask for the user's input with a 'real-time' autocompletion from the
		// words list, managed by the 'huh' package.
		var inputWord string
		huh.NewInput().
			Title("What word is in your mind? ('q' to exit)").
			Prompt("? ").
			Value(&inputWord).
			Suggestions(words). // Real-time autocompletion on text input.
			Run()

		// huh's Input adds a space at the end of the input, So we trim it.
		inputWord = strings.TrimSpace(inputWord)

		// break the loop if user wants to quit.
		if inputWord == "q" {
			break
		}

		// Do a fuzzy search with the input through the words list.
		matches := fuzzy.RankFind(inputWord, words)
		// Sort the matched results with their ranks.
		sort.Sort(matches)

		// options holds the select options.
		options := make([]huh.Option[string], OptionsLimit)

		// Generate select options to the 'options' variable.
		for i, match := range matches {
			options[i] = huh.NewOption(match.Target, match.Target)

			// break if we added enough options.
			if i+1 == OptionsLimit {
				break
			}
		}

		// Initialize a new select component and ask for the user input.
		var chosenWord string
		huh.NewSelect[string]().
			Title("Closest match to " + inputWord).
			Options(options...).
			Value(&chosenWord).
			Run()

		fmt.Printf("You chose: %s\n\n", chosenWord)
	}
}
