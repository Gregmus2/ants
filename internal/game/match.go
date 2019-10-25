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
	killed map[*global.User]float64
}

type Match struct {
	name                        string
	users                       []*global.User
	ants                        []*global.Ant
	anthills                    map[*global.User][]global.Anthill
	area                        global.Area
	queueAtTheMaternityHospital []*global.User
	stat                        *MatchStat
	s                           global.Storage
	round                       int
	part                        int
	states                      [][][]string
}

const matchesCollection string = "matches"

func CreateMatch(name string, users []*global.User, ants []*global.Ant, anthills map[*global.User][]global.Anthill, area global.Area, s global.Storage) *Match {
	match := &Match{
		name:                        name,
		users:                       users,
		ants:                        ants,
		area:                        area,
		anthills:                    anthills,
		queueAtTheMaternityHospital: make([]*global.User, 0, 10),
		s:                           s,
		stat: &MatchStat{
			ants:   make(map[*global.User]uint),
			dead:   make(map[*global.User]uint),
			killed: make(map[*global.User]float64),
		},
		round:  1,
		part:   1,
		states: make([][][]string, 0, global.Config.MatchPartSize),
	}

	for _, user := range users {
		match.stat.ants[user] = 1
		match.stat.dead[user] = 0
		match.stat.killed[user] = 0
	}

	s.CreateCollectionIfNotExist(matchesCollection)

	return match
}

func (g *Match) Run() {
	// todo give position of ants by start
	for g.isOver() {
		actions := g.collectActions()
		g.play(actions)
		g.birthStep()

		g.switchRound()
	}
}

func (g *Match) isOver() bool {
	return g.stat.CountLiving() > 1 && g.part < global.Config.MatchPartsLimit
}

func (g *Match) switchRound() {
	g.states = append(g.states, g.area.ToColorSlice())
	matchPartSizeFloat := float64(global.Config.MatchPartSize)
	if math.Mod(float64(g.round), matchPartSizeFloat) == 0 {
		g.savePart()
		g.states = make([][][]string, 0, global.Config.MatchPartSize)
		g.part++
	}
	g.round++
}

func (g *Match) collectActions() map[pkg.Action]map[global.Pos]global.Ants {
	actions := make(map[pkg.Action]map[global.Pos]global.Ants)
	for _, ant := range g.ants {
		// todo can I remove dead ants from g.ants?
		if ant.IsDead == true {
			continue
		}

		fieldTypes := g.area.NearestArea(ant)
		// todo provide round to 'Do' function
		field, action := ant.User.Algorithm().Do(fieldTypes)
		pos := g.area.RelativePosition(ant.Pos, field)
		if _, ok := actions[action]; !ok {
			actions[action] = make(map[global.Pos]global.Ants)
		}

		actions[action][pos] = append(actions[action][pos], ant)
	}

	return actions
}

// todo write this logic to instruction
func (g *Match) play(actions map[pkg.Action]map[global.Pos]global.Ants) {
	if fields, ok := actions[pkg.AttackAction]; ok {
		for targetPos, ants := range fields {
			target := g.area.ByPos(targetPos)
			if target.Type != pkg.AntField {
				continue
			}

			if target.Ant.IsDead {
				// todo remove after tests
				panic("BUG: attempt to attack ant, which has already dead")
			}

			target.Ant.IsDead = true
			g.area[targetPos.X()][targetPos.Y()] = global.CreateEmptyObject()
			g.stat.Kill(ants, target.Ant.User)
		}
	}

	if fields, ok := actions[pkg.DieAction]; ok {
		for _, ants := range fields {
			for _, ant := range ants {
				if ant.IsDead {
					continue
				}

				ant.IsDead = true
				g.area[ant.Pos.X()][ant.Pos.Y()] = global.CreateEmptyObject()
				g.stat.Die(ant.User)
			}
		}
	}

	if fields, ok := actions[pkg.EatAction]; ok {
		for targetPos, ants := range fields {
			target := g.area.ByPos(targetPos)
			if target.Type != pkg.FoodField {
				continue
			}

			if len(ants.Living()) > 1 {
				continue
			}

			ant := ants[0]
			g.queueAtTheMaternityHospital = append(g.queueAtTheMaternityHospital, ant.User)
			g.area[targetPos.X()][targetPos.Y()] = global.CreateEmptyObject()
		}
	}

	if fields, ok := actions[pkg.MoveAction]; ok {
		for targetPos, ants := range fields {
			target := g.area.ByPos(targetPos)
			if target.Type != pkg.EmptyField {
				continue
			}

			if len(ants.Living()) > 1 {
				continue
			}

			// todo in that case ant can move to field, where another ant was in that round
			ant := ants[0]
			g.area[targetPos.X()][targetPos.Y()] = g.area[ant.Pos.X()][ant.Pos.Y()]
			g.area[ant.Pos.X()][ant.Pos.Y()] = global.CreateEmptyObject()
			ant.Pos = targetPos
		}
	}
}

func (g *Match) birthStep() {
	latecomers := make([]*global.User, 0, 10)
	for _, user := range g.queueAtTheMaternityHospital {
		ok := g.giveBirth(user)
		if !ok {
			latecomers = append(latecomers, user)
		}
	}
	g.queueAtTheMaternityHospital = latecomers
}

func (g *Match) savePart() {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(g.states)
	if err != nil {
		log.Fatal(err)
	}

	err = g.s.Put(matchesCollection, g.name+strconv.Itoa(g.part), buf.Bytes())
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

		return true
	}

	return false
}

func (s *MatchStat) Kill(killers []*global.Ant, victim *global.User) {
	s.Die(victim)

	piece := math.Round(float64(1/len(killers)*100)) / 100
	for _, killer := range killers {
		s.killed[killer.User] += piece
	}
}

func (s *MatchStat) Die(who *global.User) {
	s.ants[who]--
	s.dead[who]++
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
