package opening

import (
	"os"
	"os/exec"
	"testing"
)

// this test is used to visualize the graph of openings using a dot file
func TestDrawOpening(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	fun := func(o *Opening) bool {
		return o.code == "C00"
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
