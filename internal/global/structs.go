package global

type Pos [2]uint

func (p *Pos) X() uint {
	return p[0]
}

func (p *Pos) Y() uint {
	return p[1]
}

type Ant struct {
	Pos    Pos
	User   *User
	IsDead bool
}

type Anthill struct {
	Pos      Pos
	User     *User
	BirthPos Pos
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
