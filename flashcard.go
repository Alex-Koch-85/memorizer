package flashcard

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
)

const (
	AgainPenalty = 0.2
	HardPenalty = 0.15
	EasyBonus = 1.3
	MinEase = 1.3
)

// Card struct represents a flashcard
type Card struct {
	ID					string		`json:"id"`
	Term				string		`json:"term"`
	Solution		string		`json:"solution"`

	Interval		int 			`json:"interval"`			// days until next repetition
	Repetitions	int				`json:"repetitions"`	// consecutive correct answers
	EaseFactor	float64		`json:"ease_factor"`	// how easy was the card (Anki-lite version, no full SM-2)

	DueDate			time.Time	`json:"due_date"`			// next due date
	LastReview	time.Time	`json:"last_review"`	// last learned

	Lapses			int 			`json:"lapses"`				// how many times wrong
	CreatedAt		time.Time	`json:"created_at"`
}

type Deck struct {
	Name			string		`json:"name"`
	Cards			[]Card		`json:"cards"`
	CreatedAt	time.Time	`json:"created_at"`
}

// NewCard function generates a new card (constructor method)
func NewCard(term, solution string, now time.Time) Card {
	return Card{
		ID:						generateID(),
		Term: 				term,
		Solution: 		solution,

		Repetitions: 	0,
		Interval: 		0,

		EaseFactor: 	2.5,	// standard from SM-2

		DueDate: 			now,
		LastReview: 	time.Time{},

		Lapses: 			0,
		CreatedAt: 		now,
	}
}

// NewDeck function generates a new deck (constructor method)
func NewDeck(name string, now time.Time) *Deck {
	return &Deck{
		Name: 			name,
		Cards: 			[]Card{},
		CreatedAt: 	now,
	}
}

// generateID function generates an ID for NewCard function
func generateID() string {
	return uuid.NewString()
}

// NewCard method adds a card to a Deck
func (d *Deck) AddNewCard(term, solution string, now time.Time) *Card {
	c := NewCard(term, solution, now)
	d.Cards = append(d.Cards, c)
	
	return &d.Cards[len(d.Cards)-1] 
}

// UpdateCard method updates a card via ID
func (d *Deck) UpdateCard(updated Card) bool {
	for i := range d.Cards {
		if d.Cards[i].ID == updated.ID {
			d.Cards[i] = updated
			return true
		}
	}

	return false
}

// Due method gets due cards for review
func (d *Deck) Due(now time.Time) []*Card {
	var result []*Card
	
	for i := range d.Cards {
		if !d.Cards[i].DueDate.After(now) {
			result = append(result, &d.Cards[i])
		}
	}
	
	// Sort the slice to return the oldest card first
	sort.Slice(result, func(i, j int) bool {
		return result[i].DueDate.Before(result[j].DueDate)
	})

	return result
}

// FindByID method gets the card by ID
func (d *Deck) FindByID(id string) *Card {
	for i := range d.Cards {
		if d.Cards[i].ID == id {
			return &d.Cards[i]
		}
	}

	return nil
}

// Review method implements Anki-lite algorithm for reviewing cards
func (c *Card) Review(rating int, now time.Time) error {
	if rating < 0 || rating > 3 {
		return fmt.Errorf("invalid rating: %d", rating)
	}
	
	switch rating {
		// rating 0: again; Show card again
	case 0:
		c.Repetitions = 0
		c.Interval = 1
		c.Lapses++
		c.EaseFactor -= AgainPenalty
	
		// rating 1: hard
	case 1:
		c.Repetitions++
		c.Interval = max(1, int(float64(c.Interval)*1.2))
		c.EaseFactor -= HardPenalty

	// rating 2: good
	case 2:
		c.Repetitions++
		if c.Repetitions == 1 {
			c.Interval = 1
		} else if c.Repetitions == 2 {
			c.Interval = 3
		} else {
			c.Interval = int(float64(c.Interval) * c.EaseFactor)
		}

	// rating 3: easy
	case 3:
		c.Repetitions++

		base := c.Interval
		if base == 0 {
			base = 1
		}

		c.Interval = int(float64(base) * c.EaseFactor * EasyBonus)
		c.EaseFactor += 0.15
	}

	// minimum EaseFactor
	if c.EaseFactor < MinEase {
		c.EaseFactor = MinEase
	}

	if c.Interval < 1 {
		c.Interval = 1
	}

	c.LastReview = now
	c.DueDate = now.AddDate(0, 0, c.Interval)

	return nil
}

// Save method encodes the Deck as JSON and saves it using 
// the provided filename
func (d *Deck) Save(filename string) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadDeck function opens the provided file name, decodes
// the JSON and parses it into a Deck
func LoadDeck(filename string) (*Deck, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var d Deck
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return &d, nil
}
