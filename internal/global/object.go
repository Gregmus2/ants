package global

import "ants/pkg"

type Object struct {
	Type  pkg.FieldType
	Color string
	Ant   *Ant
}

func CreateEmptyObject() *Object {
	return &Object{
		Type:  pkg.EmptyField,
		Color: "#FFFFFF",
		Ant:   nil,
	}
}

func CreateWall() *Object {
	return &Object{
		Type:  pkg.WallField,
		Color: "#8A4B1C",
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

func (o *Object) FieldTypeForUser(ant *Ant) pkg.FieldType {
	if o.Type == pkg.AntField {
		if ant.User == o.Ant.User {
			return pkg.AllyField
		} else {
			return pkg.EnemyField
		}
	}

	return o.Type
}
