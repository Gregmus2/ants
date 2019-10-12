package game

import (
	"ants/internal/global"
	"ants/pkg"
)

type Match struct {
	users              []*global.User
	ants               []*global.Ant
	area               global.Area
	queueAtTheCemetery []*global.Ant
}

func CreateMatch(users []*global.User, ants []*global.Ant, area global.Area) *Match {
	return &Match{
		users:              users,
		ants:               ants,
		area:               area,
		queueAtTheCemetery: make([]*global.Ant, 0),
	}
}

func (g *Match) Run(pipe chan [][]string) {
	for len(g.users) > 1 {
		for i := 0; i < len(g.ants); i++ {
			ant := g.ants[i]
			if ant.IsDead == true {
				continue
			}

			fieldTypes := g.area.TypesSlice(ant)
			field, action := g.ants[i].User.Algorithm().Do(fieldTypes)
			pos := g.area.RelativePosition(ant.Pos, field)
			g.do(ant, pos, action)
		}

		for i := 0; i < len(g.queueAtTheCemetery); i++ {
			g.queueAtTheCemetery[i].IsDead = true
		}
		g.queueAtTheCemetery = make([]*global.Ant, 0)

		pipe <- g.area.ToColorSlice()
	}

	close(pipe)
}

func (g *Match) do(ant *global.Ant, targetPos global.Pos, action pkg.Action) {
	target := g.area.ByPos(targetPos)
	switch action {
	case pkg.AttackAction:
		if target.Type != pkg.AntField {
			break
		}

		g.queueAtTheCemetery = append(g.queueAtTheCemetery, target.Ant)
		g.area[targetPos.X()][targetPos.Y()] = global.CreateEmptyObject()
		break

	case pkg.EatAction:
		if target.Type != pkg.FoodField {
			break
		}

		baby := &global.Ant{
			Pos:    targetPos,
			User:   ant.User,
			IsDead: false,
		}
		g.area[targetPos.X()][targetPos.Y()] = global.CreateAnt(baby)
		g.ants = append(g.ants, baby)

		break

	case pkg.MoveAction:
		if target.Type != pkg.EmptyField {
			break
		}

		g.area[targetPos.X()][targetPos.Y()] = g.area[ant.Pos.X()][ant.Pos.Y()]
		g.area[ant.Pos.X()][ant.Pos.Y()] = global.CreateEmptyObject()
		ant.Pos = targetPos
		break

	case pkg.DieAction:
		g.queueAtTheCemetery = append(g.queueAtTheCemetery, ant)
		g.area[ant.Pos.X()][ant.Pos.Y()] = global.CreateEmptyObject()
		break
	}
}
