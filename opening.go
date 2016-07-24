package opening

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"github.com/notnil/chess"
)

type Opening struct {
	code  string
	title string
	pgn   string
}

type Directory struct {
	root *node
}

func NewDirectory() *Directory {
	return buildDirectory(nil)
}

func (d *Directory) Find(g *chess.Game) *Opening {
	var o *Opening
	d.findOpening(d.root, g.Moves(), o)
	return o
}

func (d *Directory) findOpening(n *node, moves []*chess.Move, o *Opening) {
	if len(moves) == 0 {
		return
	}
	if n.opening != nil {
		o = n.opening
		log.Println(n.opening.title)
	}
	m := moves[0].String()
	c, ok := n.children[m]
	if !ok {
		return
	}
	d.findOpening(c, moves[1:len(moves)], o)
}

func buildDirectory(f func(o *Opening) bool) *Directory {
	d := &Directory{
		root: &node{
			children: map[string]*node{},
			pos:      chess.NewGame().Position(),
			label:    label(),
		},
	}
	r := csv.NewReader(bytes.NewBuffer(csvData))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i, row := range records {
		if i == 0 {
			continue
		}
		o := &Opening{code: row[0], title: row[1], pgn: row[2]}
		if f == nil || f(o) {
			d.insert(o)
		}
	}
	return d
}
func (d *Directory) insert(o *Opening) error {
	b := bytes.NewBufferString(o.pgn)
	pgn, err := chess.PGN(b)
	if err != nil {
		return err
	}
	g := chess.NewGame(pgn)
	posList := g.Positions()
	moves := g.Moves()
	n := d.root
	d.ins(n, o, posList[1:len(posList)], moves)
	return nil
}

func (d *Directory) ins(n *node, o *Opening, posList []*chess.Position, moves []*chess.Move) {
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
	d.ins(child, o, posList[1:len(posList)], moves[1:len(moves)])
}

type node struct {
	parent   *node
	children map[string]*node
	opening  *Opening
	pos      *chess.Position
	label    string
}

func (d *Directory) draw(w io.Writer) error {
	s := "digraph g {\n"
	ch := make(chan *node)
	go func() {
		d.nodes(d.root, ch)
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

func (d *Directory) nodes(root *node, ch chan *node) {
	ch <- root
	for _, c := range root.children {
		d.nodes(c, ch)
	}
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
