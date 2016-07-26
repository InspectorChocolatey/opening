package opening

import (
	"os"
	"os/exec"
	"testing"

	"github.com/notnil/chess"
)

func TestFind(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("e4"); err != nil {
		t.Fatal(err)
	}
	if err := g.MoveStr("d5"); err != nil {
		t.Fatal(err)
	}
	o := Find(g.Moves())
	expected := "Center Counter; Scandanavian; B01"
	if o == nil || o.Title() != expected {
		t.Fatalf("expected to find opening %s but got %s", expected, o.title)
	}
}

func TestPossible(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("g3"); err != nil {
		t.Fatal(err)
	}
	openings := Possible(g.Moves())
	actual := len(openings)
	if actual != 4 {
		t.Fatalf("expected %d possible openings but got %d", 4, actual)
	}
}

func TestDrawOpening(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	fun := func(o *Opening) bool {
		return o.code == "A00"
	}
	d := buildDirectory(fun)
	f, err := os.Create("test.dot")
	if err != nil {
		t.Fatal(err)
	}
	if err := d.draw(f); err != nil {
		t.Fatal(err)
	}
	// dot -Tpng mygraph.dot -o mygraph.png
	if err := exec.Command("dot", "-Tpng", "test.dot", "-o", "test.png").Run(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildDirectory(nil)
	}
}
