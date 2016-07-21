package opening

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/notnil/chess"
)

type Opening struct {
	code  string
	title string
	pgn   string
}

type node struct {
	parent   *node
	children map[string]*node
	opening  *Opening
	pos      *chess.Position
	label    string
}

type openingTree struct {
	root *node
}

func newOpeningTree() *openingTree {
	return &openingTree{
		root: &node{
			children: map[string]*node{},
			pos:      chess.NewGame().Position(),
			label:    label(),
		},
	}
}

func (t *openingTree) insert(o *Opening) error {
	b := bytes.NewBufferString(o.pgn)
	pgn, err := chess.PGN(b)
	if err != nil {
		return err
	}
	g := chess.NewGame(pgn)
	posList := g.Positions()
	moves := g.Moves()
	n := t.root
	t.ins(n, o, posList[1:len(posList)], moves)
	return nil
}

func (t *openingTree) ins(n *node, o *Opening, posList []*chess.Position, moves []*chess.Move) {
	pos := posList[0]
	move := moves[0]
	moveStr := move.String()
	var child *node
	for mv, c := range n.children {
		if mv == moveStr {
			child = c
			break
		}
	}
	if child == nil {
		child = &node{
			parent:   n,
			children: map[string]*node{},
			pos:      pos,
			label:    label(),
		}
		n.children[moveStr] = child
	}
	if len(posList) == 1 {
		child.opening = o
		return
	}
	t.ins(child, o, posList[1:len(posList)], moves[1:len(moves)])
}

func (t *openingTree) draw(w io.Writer) error {
	s := "digraph g {\n"
	ch := make(chan *node)
	go func() {
		t.nodes(t.root, ch)
		close(ch)
	}()
	for n := range ch {
		title := ""
		if n.opening != nil {
			title = n.opening.title
		}
		s += fmt.Sprintf(`%s [label="%s"];`+"\n", n.label, title)
		for m, c := range n.children {
			s += fmt.Sprintf(`%s -> %s [label="%s"];`+"\n", n.label, c.label, m)
		}
	}
	s += "}"
	_, err := w.Write([]byte(s))
	return err
}

func (t *openingTree) nodes(root *node, ch chan *node) {
	ch <- root
	for _, c := range root.children {
		t.nodes(c, ch)
	}
}

var (
	tree *openingTree
)

func buildTree(f func(o *Opening) bool) error {
	r := csv.NewReader(bytes.NewBuffer(csvData))
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	tree = newOpeningTree()
	for i, row := range records {
		if i == 0 {
			continue
		}
		o := &Opening{code: row[0], title: row[1], pgn: row[2]}
		if f == nil || f(o) {
			tree.insert(o)
		}
	}
	return nil
}

var (
	labelCount = 0
	alphabet   = "abcdefghijklmnopqrstuvwxyz"
)

func label() string {
	s := "a" + fmt.Sprint(labelCount)
	labelCount++
	return s
}
