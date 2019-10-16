package game

import (
	"ants/internal/global"
	"bytes"
	"encoding/gob"
	pkg "github.com/gregmus2/ants-pkg"
	"log"
	"math"
	"strconv"
)

type MatchStat struct {
	ants   map[*global.User]uint
	dead   map[*global.User]uint
	killed map[*global.User]uint
}

type Match struct {
	users              []*global.User
	ants               []*global.Ant
	area               global.Area
	queueAtTheCemetery []*global.Ant
	stat               *MatchStat
	s                  global.Storage
}

const matchesCollection string = "matches"

func CreateMatch(users []*global.User, ants []*global.Ant, area global.Area, s global.Storage) *Match {
	match := &Match{
		users:              users,
		ants:               ants,
		area:               area,
		queueAtTheCemetery: make([]*global.Ant, 0),
		s:                  s,
		stat: &MatchStat{
			ants:   make(map[*global.User]uint),
			dead:   make(map[*global.User]uint),
			killed: make(map[*global.User]uint),
		},
	}

	for _, user := range users {
		match.stat.ants[user] = 1
		match.stat.dead[user] = 0
		match.stat.killed[user] = 0
	}

	s.CreateCollectionIfNotExist(matchesCollection)

	return match
}

func (g *Match) Run(name string) {
	round := 1
	part := 1
	states := make([][][]string, 0, 100)
	for g.stat.CountLiving() > 1 && part < 100 {
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

		states = append(states, g.area.ToColorSlice())
		if math.Mod(float64(round), 100) == 0 {
			g.saveRound(name, part, states)
			states = make([][][]string, 0, 100)
			part++
		}
		round++
	}
}

func (g *Match) saveRound(name string, part int, states [][][]string) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(states)
	if err != nil {
		log.Fatal(err)
	}

	err = g.s.Put(matchesCollection, name+strconv.Itoa(part), buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Match) LoadRound(name string, part string) [][][]string {
	result := make([][][]string, 0, 100)
	buf := &bytes.Buffer{}
	rawData, err := g.s.Get(matchesCollection, name+part)
	if err != nil {
		log.Fatal(err)
	}

	buf.Write(rawData)
	err = gob.NewDecoder(buf).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
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
		g.stat.Kill(ant.User, target.Ant.User)
		break

	// @todo we need to handle case, when two ants want to eat one food
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
		g.stat.ants[ant.User]++

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
		g.stat.ants[ant.User]--
		g.stat.dead[ant.User]++
		break
	}
}

func (s *MatchStat) Kill(who *global.User, whom *global.User) {
	s.ants[whom]--
	s.dead[whom]++
	s.killed[who]++
}

func (s *MatchStat) CountLiving() int {
	counter := 0
	for _, count := range s.ants {
		if count > 0 {
			counter++
		}
	}

	return counter
}
