package game

import (
	"ants/internal/user"

	pkg "github.com/gregmus2/ants-pkg"
)

type Ant struct {
	ID      int
	Pos     *pkg.Pos
	User    *user.User
	IsDead  bool
	PosDiff *pkg.Pos // last move
}

type Anthills struct {
	m  map[*user.User]map[*pkg.Pos]*Anthill
	id int
}

type Anthill struct {
	ID       int
	Pos      *pkg.Pos
	User     *user.User
	BirthPos *pkg.Pos
}

type Ants struct {
	m  []*Ant
	id int
}

func NewAnts(cap int) *Ants {
	return &Ants{
		m:  make([]*Ant, 0, cap),
		id: 0,
	}
}

func (a *Ants) ID() int {
	a.id++

	return a.id
}

func NewAnthills() *Anthills {
	return &Anthills{
		m:  make(map[*user.User]map[*pkg.Pos]*Anthill),
		id: 0,
	}
}

func (a *Anthills) ID() int {
	a.id++

	return a.id
}

func (a Anthills) ByUser(user *user.User) map[*pkg.Pos]*Anthill {
	return a.m[user]
}

func (a Anthills) ByPos(pos *pkg.Pos) *Anthill {
	for _, anthills := range a.m {
		if _, exist := anthills[pos]; exist {
			return anthills[pos]
		}
	}

	return nil
}

func (a Anthills) FirstByUser(user *user.User) *Anthill {
	for _, anthill := range a.m[user] {
		return anthill
	}

	return nil
}

func (a Anthills) DeleteByPos(pos *pkg.Pos) *Anthill {
	var obj *Anthill
	for i, positions := range a.m {
		for p := range positions {
			if p.X == pos.X && p.Y == pos.Y {
				obj = a.m[i][p]
				delete(a.m[i], p)

				return obj
			}
		}
	}

	return nil
}

func (a Anthills) Add(user *user.User, pos *pkg.Pos, anthill *Anthill) {
	if _, ok := a.m[user]; !ok {
		a.m[user] = make(map[*pkg.Pos]*Anthill)
	}

	a.m[user][pos] = anthill
}

func (a Ants) Living() []*Ant {
	l := make([]*Ant, 0, len(a.m))
	for _, ant := range a.m {
		if !ant.IsDead {
			l = append(l, ant)
		}
	}

	return l
}
