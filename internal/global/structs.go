package global

import pkg "github.com/gregmus2/ants-pkg"

type Ant struct {
	Pos    pkg.Pos
	User   *User
	IsDead bool
}

type Anthill struct {
	Pos      pkg.Pos
	User     *User
	BirthPos pkg.Pos
}

type Ants []*Ant

func (ants Ants) Living() []*Ant {
	l := make([]*Ant, 0, len(ants))
	for _, ant := range ants {
		if !ant.IsDead {
			l = append(l, ant)
		}
	}

	return l
}
