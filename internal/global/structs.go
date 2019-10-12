package global

type Pos [2]uint

func (p *Pos) X() uint {
	return p[0]
}

func (p *Pos) Y() uint {
	return p[1]
}

// @todo we can add memory to ant (like small storage for users uniq by ant)
type Ant struct {
	Pos    Pos
	User   *User
	IsDead bool
}
