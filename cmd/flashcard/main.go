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

func getFileName(deck string) string {
	return strings.ToLower(deck) + ".json"
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Flashcard training program.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2026\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "To use the environment variable: 'export MEMORIZER_FILE=name.json'\n")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	addCard := flag.Bool("addCard", false, "Add a new card")
	deckName := flag.String("deck", "default", "Deck name")
	term := flag.String("term", "", "Card term")
	solution := flag.String("solution", "", "Card solution")
	review := flag.Bool("review", false, "Start review session")
	listCards := flag.Bool("list", false, "List all cards")
	editCard := flag.Bool("edit", false, "Edit a card by ID")
	deleteCard := flag.Bool("delete", false, "Delete a card by ID")
	cardID := flag.String("id", "", "Card ID")
	flag.Parse()

	// Mode counter to prevent using conflicting flags
	modeCount := 0
	if *addCard {
		modeCount++
	}
	if *listCards {
		modeCount++
	}
	if *editCard {
		modeCount++
	}
	if *deleteCard {
		modeCount++
	}
	if *review {
		modeCount++
	}

	if modeCount > 1 {
		fmt.Fprintln(os.Stderr, "please use only onde mode at a time")
		os.Exit(1)
	}

	if *deckName == "" {
		fmt.Fprintln(os.Stderr, "deck name cannot be empty")
		os.Exit(1)
	}

	fileName := getFileName(*deckName)

	// Check if the user defined the ENV VAR for a custom file name
	if env := os.Getenv("MEMORIZER_FILE"); env != "" {
		fileName = env
	}

	// Use LoadOrCreateDeck helper function to read a Deck from file or create a new one
	d, err := LoadOrCreateDeck(fileName, *deckName, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// List all the cards
	if *listCards {
		if len(d.Cards) == 0 {
			fmt.Println("No cards in deck.")
			return
		}

		fmt.Printf("Deck: %s\n\n", d.Name)

		for i, c := range d.Cards {
			fmt.Printf("%d.\n", i+1)
			fmt.Printf("  ID:				%s\n", c.ID)
			fmt.Printf("  Term:			%s\n", c.Term)
			fmt.Printf("  Solution:	%s\n", c.Solution)
			fmt.Printf("  DueDate:	%s\n", c.DueDate.Format("2006-01-02"))
			fmt.Printf("  Interval:	%d days\n", c.Interval)
			fmt.Printf("  Reps:			%d\n", c.Repetitions)
			fmt.Printf("  Lapses:		%d\n", c.Lapses)
			fmt.Println()
		}

		return
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

	// Edit edits a card by ID
	if *editCard {
		if *cardID == "" {
			fmt.Fprintln(os.Stderr, "card ID required")
			os.Exit(1)
		}

		card := d.FindByID(*cardID)
		if card == nil {
			fmt.Fprintln(os.Stderr, "card not found")
			os.Exit(1)
		}

		if *term != "" {
			card.Term = *term
		}

		if *solution != "" {
			card.Solution = *solution
		}

		SaveOrExit(d, fileName)
		fmt.Println("Card updated successfully")
		return
	}

	// Delete a card by ID
	if *deleteCard {
		if *cardID == "" {
			fmt.Fprintln(os.Stderr, "card ID required")
			os.Exit(1)
		}

		if !d.DeleteCardByID(*cardID) {
			fmt.Fprintln(os.Stderr, "card not found")
			os.Exit(1)
		}

		SaveOrExit(d, fileName)
		fmt.Println("Card deleted successfully")
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
			fmt.Printf("No deck found. Creating new deck: %s\n", name)
			return flashcard.NewDeck(name, now), nil
		}
		return nil, err
	}

	fmt.Printf("Using deck: %s\n", d.Name)
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
	queue := d.Due(now)

	if len(queue) == 0 {
		fmt.Println("No cards due")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	total := 0
	againCount := 0

	for len(queue) > 0 {
		card := queue[0]
		queue = queue[1:]

		total++

		fmt.Printf("\n--- Card %d ---\n", total)
		fmt.Println("Term:", card.Term)

		// wait for Enter
		fmt.Print("(Press enter to show solution)")
		reader.ReadString('\n')

		fmt.Println("Solution:", card.Solution)

		for {
			fmt.Print("\n0 Again | 1 Hard | 2 Good | 3 Easy | q Quit\n> ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "q" {
				fmt.Println("\nSession aborted.")
				SaveOrExit(d, filename)
				return
			}

			rating, err := strconv.Atoi(input)
			if err != nil || rating < 0 || rating > 3 {
				fmt.Println("Invalid input, try again")
				continue
			}

			if err := card.Review(rating, now); err != nil {
				fmt.Println(err)
				continue
			}

			// if again, card goes back in queue
			if rating == 0 {
				queue = append(queue, card)
				againCount++
			}

			break
		}
	}

	SaveOrExit(d, filename)

	fmt.Println("\nReview session complete.")
	fmt.Printf("Cards reviewed: %d\n", total)
	fmt.Printf("Again count: %d\n", againCount)
}
