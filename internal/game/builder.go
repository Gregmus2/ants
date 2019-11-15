package game

import (
	"ants/internal/user"
	"errors"
	"log"
	"math"
	"math/rand"

	pkg "github.com/gregmus2/ants-pkg"
)

type matchState struct {
	areaSize int
	ants     []*Ant
	players  []*user.User
	area     Area
	anthills Anthills
}

func newMatchState(areaSize int, players []*user.User) (*matchState, error) {
	if len(players) != 2 && len(players) != 4 {
		return nil, errors.New("wrong number of players")
	}

	return &matchState{areaSize: areaSize, players: players}, nil
}

func buildAnts(state *matchState) {
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
			{&pkg.Pos{quartSize, halfSize}, &pkg.Pos{quartSize + 1, halfSize}},
			{&pkg.Pos{state.areaSize - quartSize, halfSize}, &pkg.Pos{state.areaSize - quartSize - 1, halfSize}},
		}
	case 4:
		octoSize := int(math.Round(float64(state.areaSize / 8)))
		lastOctoPiece := state.areaSize - octoSize
		positions = [][2]*pkg.Pos{
			{{octoSize, octoSize}, {octoSize + 1, octoSize + 1}},
			{{lastOctoPiece, octoSize}, {lastOctoPiece - 1, octoSize + 1}},
			{{octoSize, lastOctoPiece}, {octoSize + 1, lastOctoPiece - 1}},
			{{lastOctoPiece, lastOctoPiece}, {lastOctoPiece - 1, lastOctoPiece - 1}},
		}
	default:
		log.Fatal("wrong number of players")
	}

	state.anthills = make(Anthills)
	for i := 0; i < len(state.players); i++ {
		state.area[positions[i][0].X][positions[i][0].Y] = CreateAnthill(state.players[i])
		state.anthills.Add(state.players[i], positions[i][0], &Anthill{
			Pos:      positions[i][0],
			User:     state.players[i],
			BirthPos: positions[i][1],
		})
	}

	state.ants = make([]*Ant, 0, len(state.players))
	for _, anthills := range state.anthills {
		for _, anthill := range anthills {
			ant := &Ant{
				Pos:    anthill.BirthPos,
				User:   anthill.User,
				IsDead: false,
			}

			state.ants = append(state.ants, ant)
			state.area[anthill.BirthPos.X][anthill.BirthPos.Y] = CreateAnt(ant)
		}
	}
}

func buildArea(state *matchState) {
	state.area = make([][]*Object, state.areaSize)
	lastTile := state.areaSize - 1
	for x := 0; x < state.areaSize; x++ {
		state.area[x] = make([]*Object, state.areaSize)
		for y := 0; y < state.areaSize; y++ {
			// edges
			if x == 0 || x == lastTile || y == 0 || y == lastTile {
				state.area[x][y] = CreateWall()
			} else {
				state.area[x][y] = CreateEmptyObject()
			}
		}
	}
}

func buildFood(state *matchState, percentFrom float32, percentTo float32, min int, isUniformDistribution bool) {
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
		if state.area[x][y].Type != pkg.AntField {
			state.area[x][y] = CreateFood()
		}
	}
}

// htodo symmetrically distribution
func foodUniformDistribution(state *matchState, foodCount int) {
	var xPartSize int
	var yPartSize int
	var offsets [][2]int
	halfSize := int(math.Round(float64(state.areaSize / 2)))
	antsCount := len(state.ants)
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
			if state.area[x][y].Type != pkg.AntField {
				state.area[x][y] = CreateFood()
			}
		}
	}
}
