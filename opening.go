// Package opening implements chess opening determination and exploration.
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

// A Opening represents a specific sequence of moves from the staring position.
type Opening struct {
	code  string
	title string
	pgn   string
	game  *chess.Game
}

// Code returns the Encyclopaedia of Chess Openings (ECO) code.
func (o *Opening) Code() string {
	return o.code
}

// Title returns the Encyclopaedia of Chess Openings (ECO) title of the opening.
func (o *Opening) Title() string {
	return o.title
}

// PGN returns the opening in PGN format.
func (o *Opening) PGN() string {
	return o.pgn
}

// Game returns the opening as a game.
func (o *Opening) Game() *chess.Game {
	if o.game == nil {
		pgn, _ := chess.PGN(bytes.NewBufferString(o.pgn))
		o.game = chess.NewGame(pgn)
	}
	return o.game
}

// Find returns the most specific opening for the list of moves.  If no opening is found, Find returns nil.
func Find(moves []*chess.Move) *Opening {
	for n := dir.followPath(dir.root, moves); n != nil; n = n.parent {
		if n.opening != nil {
			return n.opening
		}
	}
	return nil
}

// Possible returns the possible openings after the moves given.  If moves is empty or nil all openings are returned.
func Possible(moves []*chess.Move) []*Opening {
	n := dir.followPath(dir.root, moves)
	openings := []*Opening{}
	for _, n := range dir.nodeList(n) {
		if n.opening != nil {
			openings = append(openings, n.opening)
		}
	}
	return openings
}

var startingPosition *chess.Position
var dir *directory

func init() {
	startingPosition = &chess.Position{}
	if err := startingPosition.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")); err != nil {
		panic(err)
	}
	dir = buildDirectory(nil)
}

type directory struct {
	root *node
}

func (d *directory) followPath(n *node, moves []*chess.Move) *node {
	if len(moves) == 0 {
		return n
	}
	c, ok := n.children[moves[0].String()]
	if !ok {
		return n
	}
	return d.followPath(c, moves[1:len(moves)])
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
	for _, n := range d.nodeList(d.root) {
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

func (d *directory) nodeList(root *node) []*node {
	ch := make(chan *node)
	go func() {
		dir.nodes(root, ch)
		close(ch)
	}()
	nodes := []*node{}
	for n := range ch {
		nodes = append(nodes, n)
	}
	return nodes
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
