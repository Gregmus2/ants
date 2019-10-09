package pkg

type FieldType uint8

const EmptyField FieldType = 0
const FoodField FieldType = 1
const AllyField FieldType = 2
const EnemyField FieldType = 3
const WallField FieldType = 4
const AntField FieldType = 5

const AttackAction uint8 = 0
const EatAction uint8 = 1
const MoveAction uint8 = 2
const DieAction uint8 = 3
