package game

import (
	"ants/internal/user"

	pkg "github.com/gregmus2/ants-pkg"
)

type Object struct {
	Type  pkg.FieldType
	Color string
	Ant   *Ant
}

func CreateEmptyObject() *Object {
	return &Object{
		Type:  pkg.EmptyField,
		Color: "",
		Ant:   nil,
	}
}

func CreateWall() *Object {
	return &Object{
		Type:  pkg.WallField,
		Color: "brown",
		Ant:   nil,
	}
}

func CreateAnt(ant *Ant) *Object {
	return &Object{
		Type:  pkg.AntField,
		Color: ant.User.Color,
		Ant:   ant,
	}
}

func CreateFood() *Object {
	return &Object{
		Type:  pkg.FoodField,
		Color: "yellow",
		Ant:   nil,
	}
}

// todo color of anthill must be
func CreateAnthill(u *user.User) *Object {
	return &Object{
		Type:  pkg.AnthillField,
		Color: u.Color,
		Ant:   nil,
	}
}

func (o *Object) FieldTypeForUser(ant *Ant) pkg.FieldType {
	if o.Type == pkg.AntField {
		if ant.User == o.Ant.User {
			return pkg.AllyField
		}

		return pkg.EnemyField
	}

	if o.Type == pkg.AnthillField {
		// change AnthillField and get user from it to compare
		if ant.User.Color == o.Color {
			return pkg.AllyAnthillField
		}

		return pkg.EnemyAnthillField
	}

	return o.Type
}
