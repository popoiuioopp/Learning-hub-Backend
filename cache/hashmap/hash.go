package hashmap

import "fmt"

type Flashcard struct {
	FlashcardID int
	DeckID      int
	Term        string
	Definition  string
}

type Node struct {
	Next      *Node
	Flashcard Flashcard
}

type HashTable struct {
	Nodes []Node
	Head  *Node
	Size  int
}

func (f Flashcard) Status() string {
	return fmt.Sprintf("Flashcard ID %d, Deck ID %d, \nterm : %s\ndefinition : %s\n", f.FlashcardID, f.DeckID, f.Term, f.Definition)
}
