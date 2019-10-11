package pkg

var fieldTypeToActionMap = map[FieldType]Action{
	EmptyField: MoveAction,
	FoodField:  EatAction,
	EnemyField: AttackAction,
	AllyField:  NoAction,
	WallField:  NoAction,
	AntField:   NoAction,
}

func ResolveAction(fieldType FieldType) Action {
	return fieldTypeToActionMap[fieldType]
}
