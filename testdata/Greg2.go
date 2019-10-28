package main

import (
	"github.com/gregmus2/ants-pkg"
	"math/rand"
	"time"
)

type greg2 string

var r2 *rand.Rand

func init() {
	r2 = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (g greg2) Start(anthill pkg.Pos, birth pkg.Pos) {

}

func (g greg2) Do(fields [9]pkg.FieldType) (field uint8, action pkg.Action) {
	field = uint8(r2.Intn(9))
	action = pkg.ResolveAction(fields[field])

	return
}

var Greg2 greg2
