# opening
[![GoDoc](https://godoc.org/github.com/notnil/opening?status.svg)](https://godoc.org/github.com/notnil/opening)
[![Build Status](https://drone.io/github.com/notnil/opening/status.png)](https://drone.io/github.com/notnil/opening/latest)
[![Go Report Card](https://goreportcard.com/badge/notnil/opening)](https://goreportcard.com/report/notnil/opening)

## Datasource

The [Encyclopaedia of Chess Openings](https://en.wikipedia.org/wiki/Encyclopaedia_of_Chess_Openings) (ECO) functions as the datasource for this package.  A consise list of openings with PGNs can be found [here](http://www.webcitation.org/query?url=http://www.geocities.com/siliconvalley/lab/7378/eco.htm&date=2010-02-20+10:14:24).

## Example

```go   
package main

import (
    "fmt"

    "github.com/notnil/chess"
    "github.com/notnil/opening"
)

func main(){
    g := chess.NewGame()
    g.MoveStr("e4")
    g.MoveStr("e6")
    
    o := opening.Find(g.Moves())
    fmt.Println(o.Title()) // French Defense
}
```