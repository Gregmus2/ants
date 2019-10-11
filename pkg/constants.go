package pkg

type FieldType uint8

const EmptyField FieldType = 0
const FoodField FieldType = 1
const AllyField FieldType = 2
const EnemyField FieldType = 3
const WallField FieldType = 4
const AntField FieldType = 5

type Action uint8

const NoAction Action = 0
const AttackAction Action = 1
const EatAction Action = 2
const MoveAction Action = 3
const DieAction Action = 4
