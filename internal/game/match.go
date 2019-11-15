package game

import (
	"ants/internal/user"
	"bytes"
	"encoding/gob"
	"log"
	"math"
	"strconv"

	pkg "github.com/gregmus2/ants-pkg"
)

type MatchStat struct {
	ants   map[*user.User]uint
	dead   map[*user.User]uint
	killed map[*user.User]float64
}

type Match struct {
	name     string
	players  []*user.User
	ants     []*Ant
	anthills Anthills
	area     Area
	birthQ   []*user.User
	stat     *MatchStat
	service  *Service
	round    int
	part     int
	states   [][][]string
}

const matchesCollection string = "matches"

func CreateMatch(gameService *Service, state *matchState, name string) *Match {
	if state.players == nil || state.ants == nil {
		log.Fatal("builder must have at least players and ants")
	}

	match := &Match{
		name:     name,
		players:  state.players,
		ants:     state.ants,
		area:     state.area,
		anthills: state.anthills,
		birthQ:   make([]*user.User, 0, 10),
		service:  gameService,
		stat: &MatchStat{
			ants:   make(map[*user.User]uint),
			dead:   make(map[*user.User]uint),
			killed: make(map[*user.User]float64),
		},
		round:  1,
		part:   1,
		states: make([][][]string, 0, gameService.config.Match.PartSize),
	}

	for _, player := range state.players {
		match.stat.ants[player] = 1
		match.stat.dead[player] = 0
		match.stat.killed[player] = 0
	}

	gameService.storage.CreateCollectionIfNotExist(matchesCollection)

	return match
}

func (g *Match) Run() {
	g.start()
	for g.isOver() {
		actions := g.collectActions()
		g.play(actions)
		g.birthStep()

		g.switchRound()
	}
}

func (g *Match) isOver() bool {
	return g.stat.CountLiving() > 1 && g.part < g.service.config.Match.PartsLimit
}

func (g *Match) switchRound() {
	g.states = append(g.states, g.area.ToColorSlice())
	matchPartSizeFloat := float64(g.service.config.Match.PartSize)
	if math.Mod(float64(g.round), matchPartSizeFloat) == 0 {
		g.savePart()
		g.states = make([][][]string, 0, g.service.config.Match.PartSize)
		g.part++
	}
	g.round++
}

func (g *Match) collectActions() map[pkg.Action]map[*pkg.Pos]Ants {
	actions := make(map[pkg.Action]map[*pkg.Pos]Ants)
	for _, ant := range g.ants {
		// todo can I remove dead ants from g.ants?
		if ant.IsDead {
			continue
		}

		fieldTypes := g.area.VisibleArea(ant)
		pos, action := ant.User.Algorithm().Do(fieldTypes, g.round*g.part)
		if pos.X < -1 || pos.X > 1 || pos.Y < -1 || pos.Y > 1 {
			continue
		}

		pos.Add(ant.Pos)
		if pos.X < 0 || pos.Y < 0 {
			continue
		}

		if _, ok := actions[action]; !ok {
			actions[action] = make(map[*pkg.Pos]Ants)
		}

		actions[action][pos] = append(actions[action][pos], ant)
	}

	return actions
}

func (g *Match) start() {
	for _, ant := range g.ants {
		anthill := g.anthills.FirstByUser(ant.User)
		ant.User.Algorithm().Start(*anthill.Pos, *anthill.BirthPos)
	}
}

// todo write this logic to instruction
func (g *Match) play(actions map[pkg.Action]map[*pkg.Pos]Ants) {
	if fields, ok := actions[pkg.AttackAction]; ok {
		g.attackStep(fields)
	}

	if fields, ok := actions[pkg.DieAction]; ok {
		g.suicideStep(fields)
	}

	if fields, ok := actions[pkg.EatAction]; ok {
		g.eatStep(fields)
	}

	if fields, ok := actions[pkg.MoveAction]; ok {
		g.moveStep(fields)
	}
}

// htodo capture anthill
func (g *Match) attackStep(fields map[*pkg.Pos]Ants) {
	for targetPos, ants := range fields {
		target := g.area.ByPos(targetPos)
		switch target.Type {
		case pkg.AntField:
			g.handleAttackAnt(targetPos, target.Ant, ants)
		case pkg.AnthillField:
			g.handleAttackAnthill(targetPos, ants)
		}
	}
}

func (g *Match) handleAttackAnt(targetPos *pkg.Pos, victim *Ant, ants Ants) {
	if victim.IsDead {
		// todo remove after tests
		panic("BUG: attempt to attack ant, which has already dead")
	}

	killers := make([]*user.User, 0, 1)
	bestPower := 0
	for _, ant := range ants {
		power := g.area.CalcAtkPower(victim, ant)
		switch {
		case power < bestPower:
			continue
		case power == bestPower:
			killers = append(killers, ant.User)
		default:
			killers = []*user.User{ant.User}
		}
	}

	if len(killers) <= 0 {
		return
	}

	victim.IsDead = true
	// todo ant would be part of atkPower of another ant in that round
	g.area[targetPos.X][targetPos.Y] = CreateEmptyObject()
	g.stat.Kill(killers, victim.User)
}

func (g *Match) handleAttackAnthill(targetPos *pkg.Pos, ants Ants) {
	users := make(map[*user.User]bool)
	invaders := 0
	for _, ant := range ants {
		if _, exist := users[ant.User]; !exist {
			users[ant.User] = true
			invaders++
		}
	}

	if invaders > 1 {
		return
	}

	g.area[targetPos.X][targetPos.Y] = CreateAnthill(ants[0].User)
	anthill := g.anthills.DeleteByPos(targetPos)
	anthill.User = ants[0].User
	g.anthills.Add(ants[0].User, targetPos, anthill)
}

func (g *Match) suicideStep(fields map[*pkg.Pos]Ants) {
	for _, ants := range fields {
		for _, ant := range ants {
			if ant.IsDead {
				continue
			}

			ant.IsDead = true
			g.area[ant.Pos.X][ant.Pos.Y] = CreateEmptyObject()
			g.stat.Die(ant.User)
		}
	}
}

func (g *Match) eatStep(fields map[*pkg.Pos]Ants) {
	for targetPos, ants := range fields {
		target := g.area.ByPos(targetPos)
		if target.Type != pkg.FoodField {
			continue
		}

		if len(ants.Living()) > 1 {
			continue
		}

		ant := ants[0]
		g.birthQ = append(g.birthQ, ant.User)
		g.area[targetPos.X][targetPos.Y] = CreateEmptyObject()
	}
}

func (g *Match) moveStep(fields map[*pkg.Pos]Ants) {
	for targetPos, ants := range fields {
		target := g.area.ByPos(targetPos)
		if target.Type != pkg.EmptyField {
			continue
		}

		if len(ants.Living()) > 1 {
			continue
		}

		// fixme in that case ant can move to field, where another ant was in that round
		ant := ants[0]
		g.area[targetPos.X][targetPos.Y] = g.area[ant.Pos.X][ant.Pos.Y]
		g.area[ant.Pos.X][ant.Pos.Y] = CreateEmptyObject()
		ant.Pos = targetPos
	}
}

func (g *Match) birthStep() {
	latecomers := make([]*user.User, 0, 10)
	for _, user := range g.birthQ {
		ok := g.giveBirth(user)
		if !ok {
			latecomers = append(latecomers, user)
		}
	}
	g.birthQ = latecomers
}

// todo remove old parts or load on restart
func (g *Match) savePart() {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(g.states)
	if err != nil {
		log.Fatal(err)
	}

	// todo delegate on service
	err = g.service.storage.Put(matchesCollection, g.name+strconv.Itoa(g.part), buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Match) LoadRound(name string, part string) [][][]string {
	result := make([][][]string, 0, g.service.config.Match.PartSize)
	buf := &bytes.Buffer{}
	rawData, err := g.service.storage.Get(matchesCollection, name+part)
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

func (g *Match) giveBirth(user *user.User) bool {
	for _, anthill := range g.anthills.ByUser(user) {
		if g.area[anthill.BirthPos.X][anthill.BirthPos.Y].Type != pkg.EmptyField {
			continue
		}

		baby := &Ant{
			Pos:    anthill.BirthPos,
			User:   user,
			IsDead: false,
		}
		g.area[anthill.BirthPos.X][anthill.BirthPos.Y] = CreateAnt(baby)
		g.ants = append(g.ants, baby)
		g.stat.ants[user]++

		return true
	}

	return false
}

func (s *MatchStat) Kill(killers []*user.User, victim *user.User) {
	s.Die(victim)

	piece := math.Round(float64(1/len(killers)*100)) / 100
	for _, killer := range killers {
		s.killed[killer] += piece
	}
}

func (s *MatchStat) Die(who *user.User) {
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
