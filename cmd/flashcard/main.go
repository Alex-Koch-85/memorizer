package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	flashcard "github.com/Alex-Koch-85/memorizer"
)

var fileName = ".fc.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Flashcard training program.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2026\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	addCard := flag.Bool("addCard", false, "Add a new card")
	term := flag.String("term", "", "Card term")
	solution := flag.String("solution", "", "Card solution")
	review := flag.Bool("review", false, "Start review session")
	flag.Parse()

	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("MEMORIZER_FILE") != "" {
		fileName = os.Getenv("MEMORIZER_FILE")
	}

	// Use LoadOrCreateDeck helper function to read a Deck from file or create a new one
	d, err := LoadOrCreateDeck(fileName, "default", time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Add a card to a deck
	if *addCard {
		if *term == "" || *solution == "" {
			fmt.Fprintln(os.Stderr, "term and solution required")
			os.Exit(1)
		}

		d.AddNewCard(*term, *solution, time.Now())
		SaveOrExit(d, fileName)

		fmt.Println("Card added successfully")
		return
	}

	// Start review session
	if *review {
		RunReview(d, fileName, time.Now())
		return
	}
}

// Helper function to load or create a deck
func LoadOrCreateDeck(filename, name string, now time.Time) (*flashcard.Deck, error) {
	d, err := flashcard.LoadDeck(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return flashcard.NewDeck(name, now), nil
		}
		return nil, err
	}
	return d, nil
}

// Helper function to save a deck or Exit
func SaveOrExit(d *flashcard.Deck, filename string) {
	if err := d.Save(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// RunReview function runs the review loop
func RunReview(d *flashcard.Deck, filename string, now time.Time) {
	due := d.Due(now)

	if len(due) == 0 {
		fmt.Println("No cards due")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for i, card := range due {
		fmt.Printf("\nCard %d/%d\n", i+1, len(due))
		fmt.Println("Term:", card.Term)

		// wait for Enter
		reader.ReadString('\n')

		fmt.Println("Solution:", card.Solution)

		for {
			fmt.Print("\n0 Again | 1 Hard | 2 Good | 3 Easy\n> ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			rating, err := strconv.Atoi(input)
			if err != nil || rating < 0 || rating > 3 {
				fmt.Println("Invalid input, try again")
				continue
			}

			if err := card.Review(rating, now); err != nil {
				fmt.Println(err)
				continue
			}

			break
		}
	}

	SaveOrExit(d, filename)
	fmt.Println("\nReview session complete.")
}
