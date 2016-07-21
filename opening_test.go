package opening

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestOpening(t *testing.T) {
	if err := buildTree(nil); err != nil {
		t.Fatal(err)
	}
	ch := make(chan *node)
	go func() {
		tree.nodes(tree.root, ch)
		close(ch)
	}()
	count := 0
	for range ch {
		count++
	}
	log.Println(count)
}

func TestDrawOpening(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	fun := func(o *Opening) bool {
		return o.code == "A00"
	}
	if err := buildTree(fun); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("test.dot")
	if err != nil {
		t.Fatal(err)
	}
	if err := tree.draw(f); err != nil {
		t.Fatal(err)
	}
	// dot -Tpng mygraph.dot -o mygraph.png
	if err := exec.Command("dot", "-Tpng", "test.dot", "-o", "test.png").Run(); err != nil {
		t.Fatal(err)
	}
}
