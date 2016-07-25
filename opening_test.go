package opening

import (
	"os"
	"os/exec"
	"testing"

	"github.com/notnil/chess"
)

func TestOpening(t *testing.T) {
	d := buildDirectory(nil)
	g := chess.NewGame()
	if err := g.MoveStr("e4"); err != nil {
		t.Fatal(err)
	}
	if err := g.MoveStr("d5"); err != nil {
		t.Fatal(err)
	}
	o := d.Find(g)
	if o == nil {
		t.Fatal("expected to find scandanavian opening")
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
