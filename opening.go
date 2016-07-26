package opening

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/notnil/chess"
)

var startingPosition *chess.Position
var dir *directory

func init() {
	startingPosition = &chess.Position{}
	if err := startingPosition.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")); err != nil {
		panic(err)
	}
	dir = buildDirectory(nil)
}

type Opening struct {
	code  string
	title string
	pgn   string
}

func Find(g *chess.Game) *Opening {
	return dir.findOpening(dir.root, g.Moves(), nil)
}

type directory struct {
	root *node
}

func (d *directory) findOpening(n *node, moves []*chess.Move, o *Opening) *Opening {
	if n.opening != nil {
		o = n.opening
	}
	if len(moves) == 0 {
		return o
	}
	m := moves[0].String()
	c, ok := n.children[m]
	if !ok {
		return o
	}
	return d.findOpening(c, moves[1:len(moves)], o)
}

func buildDirectory(f func(o *Opening) bool) *directory {
	d := &directory{
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

func (d *directory) insert(o *Opening) error {
	posList := []*chess.Position{startingPosition}
	moves := []*chess.Move{}
	for _, s := range parseMoveList(o.pgn) {
		pos := posList[len(posList)-1]
		m, err := chess.LongAlgebraicNotation{}.Decode(pos, s)
		if err != nil {
			panic(err)
		}
		moves = append(moves, m)
		posList = append(posList, pos.Update(m))
	}
	n := d.root
	d.ins(n, o, posList[1:len(posList)], moves)
	return nil
}

func (d *directory) ins(n *node, o *Opening, posList []*chess.Position, moves []*chess.Move) {
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

func (d *directory) draw(w io.Writer) error {
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

func (d *directory) nodes(root *node, ch chan *node) {
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

// 1.b2b4 e7e5 2.c1b2 f7f6 3.e2e4 f8b4 4.f1c4 b8c6 5.f2f4 d8e7 6.f4f5 g7g6
func parseMoveList(pgn string) []string {
	strs := strings.Split(pgn, " ")
	cp := []string{}
	for _, s := range strs {
		i := strings.Index(s, ".")
		if i == -1 {
			cp = append(cp, s)
		} else {
			cp = append(cp, s[i+1:len(s)])
		}
	}
	return cp
}
