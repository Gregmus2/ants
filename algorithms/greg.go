package main

import (
	"ants/pkg"
	"math/rand"
	"time"
)

type greg string

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (g greg) Do(fields [9]pkg.FieldType) (field uint8, action uint8) {
	return uint8(r.Intn(9)), pkg.MoveAction
}

var Greg greg
