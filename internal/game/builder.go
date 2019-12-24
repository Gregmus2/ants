package game

import (
	"ants/internal/user"
	"errors"
	"log"
	"math"
	"math/rand"

	pkg "github.com/gregmus2/ants-pkg"
)

type MatchState struct {
	areaSize int
	ants     *Ants
	players  []*user.User
	area     *Area
	anthills *Anthills
}

func NewMatchState(areaSize int, players []*user.User) (*MatchState, error) {
	if len(players) != 2 && len(players) != 4 {
		return nil, errors.New("wrong number of players")
	}

	return &MatchState{areaSize: areaSize, players: players}, nil
}

func BuildAnts(state *MatchState) {
	if state.area == nil {
		log.Fatal("builder must have area before build ants")
	}

	// [players][position, birthQ position]
	var positions [][2]*pkg.Pos
	quartSize := int(math.Round(float64(state.areaSize / 4)))
	halfSize := int(math.Round(float64(state.areaSize / 2)))

	switch len(state.players) {
	case 2:
		positions = [][2]*pkg.Pos{
			{&pkg.Pos{X: quartSize, Y: halfSize}, &pkg.Pos{X: quartSize + 1, Y: halfSize}},
			{&pkg.Pos{X: state.areaSize - quartSize, Y: halfSize}, &pkg.Pos{X: state.areaSize - quartSize - 1, Y: halfSize}},
		}
	case 4:
		octoSize := int(math.Round(float64(state.areaSize / 8)))
		lastOctoPiece := state.areaSize - octoSize
		positions = [][2]*pkg.Pos{
			{{X: octoSize, Y: octoSize}, {X: octoSize + 1, Y: octoSize + 1}},
			{{X: lastOctoPiece, Y: octoSize}, {X: lastOctoPiece - 1, Y: octoSize + 1}},
			{{X: octoSize, Y: lastOctoPiece}, {X: octoSize + 1, Y: lastOctoPiece - 1}},
			{{X: lastOctoPiece, Y: lastOctoPiece}, {X: lastOctoPiece - 1, Y: lastOctoPiece - 1}},
		}
	default:
		log.Fatal("wrong number of players")
	}

	state.anthills = NewAnthills()
	for i := 0; i < len(state.players); i++ {
		state.area.matrix[positions[i][0].X][positions[i][0].Y] = CreateAnthill(state.players[i])
		state.anthills.Add(state.players[i], positions[i][0], &Anthill{
			ID:       state.anthills.ID(),
			Pos:      positions[i][0],
			User:     state.players[i],
			BirthPos: positions[i][1],
		})
	}

	state.ants = NewAnts(len(state.players))
	for _, anthills := range state.anthills.m {
		for _, anthill := range anthills {
			ant := &Ant{
				ID:      state.ants.ID(),
				Pos:     anthill.BirthPos,
				User:    anthill.User,
				IsDead:  false,
				PosDiff: &pkg.Pos{},
			}

			state.ants.m = append(state.ants.m, ant)
			state.area.matrix[anthill.BirthPos.X][anthill.BirthPos.Y] = CreateAnt(ant)
		}
	}
}

func BuildArea(state *MatchState) {
	state.area = NewArea(state.areaSize, state.areaSize)
}

func BuildFood(state *MatchState, percentFrom float32, percentTo float32, min int, isUniformDistribution bool) {
	if state.area == nil || state.ants == nil {
		log.Fatal("builder must have ants and area before build food")
	}

	randomPercent := rand.Float32()*(percentTo-percentFrom) + percentFrom
	foodCount := int(float32(state.areaSize*state.areaSize) * randomPercent)
	if foodCount < min {
		foodCount = min
	}

	if isUniformDistribution {
		foodUniformDistribution(state, foodCount)
		return
	}

	size := state.areaSize - 2
	for i := 0; i < foodCount; i++ {
		x := rand.Intn(size) + 1
		y := rand.Intn(size) + 1
		if state.area.matrix[x][y].Type != pkg.AntField {
			state.area.matrix[x][y] = CreateFood()
		}
	}
}

// todo symmetrically distribution
func foodUniformDistribution(state *MatchState, foodCount int) {
	var xPartSize int
	var yPartSize int
	var offsets [][2]int
	halfSize := int(math.Round(float64(state.areaSize / 2)))
	antsCount := len(state.ants.m)
	switch antsCount {
	case 2:
		// -1 because of walls
		xPartSize = halfSize - 1
		yPartSize = state.areaSize - 2
		offsets = [][2]int{{1, 1}, {halfSize, 1}}
	case 4:
		xPartSize = halfSize - 1
		yPartSize = halfSize - 1
		offsets = [][2]int{{1, 1}, {halfSize, 1}, {1, halfSize}, {halfSize, halfSize}}
	default:
		log.Fatal("wrong number of ants")
	}

	for i := 0; i < foodCount; i += antsCount {
		if foodCount-i < antsCount {
			break
		}

		for j := 0; j < antsCount; j++ {
			x := rand.Intn(xPartSize) + offsets[j][0]
			y := rand.Intn(yPartSize) + offsets[j][1]
			if state.area.matrix[x][y].Type != pkg.AntField {
				state.area.matrix[x][y] = CreateFood()
			}
		}
	}
}
