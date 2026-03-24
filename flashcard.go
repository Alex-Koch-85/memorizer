package flashcard


// item struct represents a flashcard
type item struct {
	Term			string
	Solution	string
	TimesRepeated	int
	Done			bool
}

// List represents a list of flashcard items
type List []item
