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
