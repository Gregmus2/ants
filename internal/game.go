package internal

import (
	"ants/internal/global"
	"ants/pkg"
	"log"
)

type Game struct {
	users              []*global.User
	ants               []*global.Ant
	area               global.Area
	queueAtTheCemetery []*global.Ant
}

func CreateGame(users []*global.User, ants []*global.Ant, area global.Area) *Game {
	return &Game{
		users:              users,
		ants:               ants,
		area:               area,
		queueAtTheCemetery: make([]*global.Ant, 0),
	}
}

func (g *Game) Run(pipe chan [][]string) {
	for len(g.users) > 1 {
		for i := 0; i < len(g.ants); i++ {
			ant := g.ants[i]
			if ant.IsDead == true {
				continue
			}

			fieldTypes := g.area.TypesSlice(ant)
			field, action := g.ants[i].User.Algorithm().Do(fieldTypes)
			log.Println(ant.User.Name)
			log.Println(len(g.ants))
			log.Println(field, action)
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

func (g *Game) do(ant *global.Ant, targetPos global.Pos, action uint8) {
	target := g.area.ByPos(targetPos)
	switch action {
	case pkg.AttackAction:
		log.Println("attack")
		if target.Type != pkg.AntField {
			break
		}

		g.queueAtTheCemetery = append(g.queueAtTheCemetery, target.Ant)
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
		log.Println("die")
		g.queueAtTheCemetery = append(g.queueAtTheCemetery, ant)
		break
	}
}
