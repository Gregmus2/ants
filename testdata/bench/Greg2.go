package main

import (
	"github.com/gregmus2/ants-pkg"
	"math/rand"
	"time"
)

type greg2 string

var r2 *rand.Rand

func init() {
	r2 = rand.New(rand.NewSource(time.Now().Unix()))
}

func (g greg2) Start(anthillID int, birthPos pkg.Pos) {}
func (g greg2) OnAntDie(antID int) {}
func (g greg2) OnAnthillDie(anthillID int) {}
func (g greg2) OnAntBirth(antID int, anthillID int) {}
func (g greg2) OnNewAnthill(invaderID int, birthPos pkg.Pos, anthillID int) {}

func (g greg2) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff pkg.Pos) (target *pkg.Pos, action pkg.Action) {
	target = &pkg.Pos{X: r2.Intn(3) - 1, Y: r2.Intn(3) - 1}
	action = pkg.ResolveAction(fields[target.X + 3][target.Y + 3])

	return
}

var Greg2 greg2
