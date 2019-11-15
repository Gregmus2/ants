package game

import (
	"ants/internal/user"

	pkg "github.com/gregmus2/ants-pkg"
)

type Ant struct {
	Pos    *pkg.Pos
	User   *user.User
	IsDead bool
}

type Anthills map[*user.User]map[*pkg.Pos]*Anthill

type Anthill struct {
	Pos      *pkg.Pos
	User     *user.User
	BirthPos *pkg.Pos
}

type Ants []*Ant

func (ah Anthills) ByUser(user *user.User) map[*pkg.Pos]*Anthill {
	return ah[user]
}

func (ah Anthills) ByPos(pos *pkg.Pos) *Anthill {
	for _, anthills := range ah {
		if _, exist := anthills[pos]; exist {
			return anthills[pos]
		}
	}

	return nil
}

func (ah Anthills) FirstByUser(user *user.User) *Anthill {
	for _, anthill := range ah[user] {
		return anthill
	}

	return nil
}

func (ah Anthills) DeleteByPos(pos *pkg.Pos) *Anthill {
	var obj *Anthill
	for i, anthills := range ah {
		if _, exist := anthills[pos]; exist {
			obj = ah[i][pos]
			delete(ah[i], pos)

			return obj
		}
	}

	return nil
}

func (ah Anthills) Add(user *user.User, pos *pkg.Pos, anthill *Anthill) {
	if _, ok := ah[user]; !ok {
		ah[user] = make(map[*pkg.Pos]*Anthill)
	}

	ah[user][pos] = anthill
}

func (ants Ants) Living() []*Ant {
	l := make([]*Ant, 0, len(ants))
	for _, ant := range ants {
		if !ant.IsDead {
			l = append(l, ant)
		}
	}

	return l
}
