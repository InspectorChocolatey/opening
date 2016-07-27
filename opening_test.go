package opening_test

import (
	"fmt"
	"testing"

	"github.com/notnil/chess"
	"github.com/notnil/opening"
)

func ExampleFind() {
	g := chess.NewGame()
	g.MoveStr("e4")
	g.MoveStr("e6")

	// print French Defense
	o := opening.Find(g.Moves())
	fmt.Println(o.Title())
}

func ExamplePossible() {
	g := chess.NewGame()
	g.MoveStr("e4")
	g.MoveStr("d5")

	// print all variantions of the Scandinavian Defense
	for _, o := range opening.Possible(g.Moves()) {
		fmt.Println(o.Title())
	}
}

func TestFind(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("e4"); err != nil {
		t.Fatal(err)
	}
	if err := g.MoveStr("d5"); err != nil {
		t.Fatal(err)
	}
	o := opening.Find(g.Moves())
	expected := "Center Counter; Scandanavian; B01"
	if o == nil || o.Title() != expected {
		t.Fatalf("expected to find opening %s but got %s", expected, o.Title())
	}
}

func TestPossible(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("g3"); err != nil {
		t.Fatal(err)
	}
	openings := opening.Possible(g.Moves())
	actual := len(openings)
	if actual != 4 {
		t.Fatalf("expected %d possible openings but got %d", 4, actual)
	}
}
