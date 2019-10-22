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
	users                       []*global.User
	ants                        []*global.Ant
	anthills                    map[*global.User][]global.Anthill
	area                        global.Area
	queueAtTheCemetery          []*global.Ant
	queueAtTheMaternityHospital []*global.User
	stat                        *MatchStat
	s                           global.Storage
}

const matchesCollection string = "matches"

func CreateMatch(users []*global.User, ants []*global.Ant, anthills map[*global.User][]global.Anthill, area global.Area, s global.Storage) *Match {
	match := &Match{
		users:                       users,
		ants:                        ants,
		area:                        area,
		anthills:                    anthills,
		queueAtTheCemetery:          make([]*global.Ant, 10),
		queueAtTheMaternityHospital: make([]*global.User, 10),
		s:                           s,
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
	matchPartSizeFloat := float64(global.Config.MatchPartSize)
	round := 1
	part := 1
	states := make([][][]string, 0, global.Config.MatchPartSize)
	// todo give position of ants by start
	// todo add anthills. All ants born in anthills
	for g.stat.CountLiving() > 1 && part < global.Config.MatchPartsLimit {
		for i := 0; i < len(g.ants); i++ {
			ant := g.ants[i]
			if ant.IsDead == true {
				continue
			}

			fieldTypes := g.area.TypesSlice(ant)
			// todo give round to 'Do' function
			field, action := g.ants[i].User.Algorithm().Do(fieldTypes)
			pos := g.area.RelativePosition(ant.Pos, field)
			g.do(ant, pos, action)
		}

		for i := 0; i < len(g.queueAtTheCemetery); i++ {
			g.queueAtTheCemetery[i].IsDead = true
		}
		g.queueAtTheCemetery = make([]*global.Ant, 10)

		latecomers := make([]*global.User, 10)
		for _, user := range g.queueAtTheMaternityHospital {
			ok := g.giveBirth(user)
			if !ok {
				latecomers = append(latecomers, user)
			}
		}
		g.queueAtTheMaternityHospital = latecomers

		states = append(states, g.area.ToColorSlice())
		if math.Mod(float64(round), matchPartSizeFloat) == 0 {
			g.saveRound(name, part, states)
			states = make([][][]string, 0, global.Config.MatchPartSize)
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
	result := make([][][]string, 0, global.Config.MatchPartSize)
	buf := &bytes.Buffer{}
	rawData, err := g.s.Get(matchesCollection, name+part)
	if err != nil {
		log.Fatal(err)
	}

	if len(rawData) == 0 {
		return make([][][]string, 0)
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
	// todo check if the target has already dead
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

		g.queueAtTheMaternityHospital = append(g.queueAtTheMaternityHospital, ant.User)

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

func (g *Match) giveBirth(user *global.User) bool {
	for _, anthill := range g.anthills[user] {
		if g.area[anthill.BirthPos.X()][anthill.BirthPos.Y()].Type != pkg.EmptyField {
			continue
		}

		baby := &global.Ant{
			Pos:    anthill.BirthPos,
			User:   user,
			IsDead: false,
		}
		g.area[anthill.BirthPos.X()][anthill.BirthPos.Y()] = global.CreateAnt(baby)
		g.ants = append(g.ants, baby)
		g.stat.ants[user]++
	}

	return false
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
