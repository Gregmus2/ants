package game

import (
	"ants/internal/global"
	"errors"
	"log"
	"math"
	"math/rand"

	pkg "github.com/gregmus2/ants-pkg"
)

type MatchBuilder struct {
	name     string
	areaSize int
	ants     []*global.Ant
	players  []*global.User
	area     global.Area
	anthills global.Anthills
}

func NewMatchBuilder(name string, areaSize int, players []*global.User) (*MatchBuilder, error) {
	if len(players) != 2 && len(players) != 4 {
		return nil, errors.New("wrong number of players")
	}

	return &MatchBuilder{name: name, areaSize: areaSize, players: players}, nil
}

func (mb *MatchBuilder) BuildAnts() {
	if mb.area == nil {
		log.Fatal("builder must have area before build ants")
	}

	// [players][position, birthQ position]
	var positions [][2][2]uint
	quartSize := uint(math.Round(float64(mb.areaSize / 4)))
	halfSize := uint(math.Round(float64(mb.areaSize / 2)))

	switch len(mb.players) {
	case 2:
		positions = [][2][2]uint{
			{{quartSize, halfSize}, {quartSize + 1, halfSize}},
			{{uint(mb.areaSize) - quartSize, halfSize}, {uint(mb.areaSize) - quartSize - 1, halfSize}},
		}
	case 4:
		octoSize := uint(math.Round(float64(mb.areaSize / 8)))
		lastOctoPiece := uint(mb.areaSize) - octoSize
		positions = [][2][2]uint{
			{{octoSize, octoSize}, {octoSize + 1, octoSize + 1}},
			{{lastOctoPiece, octoSize}, {lastOctoPiece - 1, octoSize + 1}},
			{{octoSize, lastOctoPiece}, {octoSize + 1, lastOctoPiece - 1}},
			{{lastOctoPiece, lastOctoPiece}, {lastOctoPiece - 1, lastOctoPiece - 1}},
		}
	default:
		log.Fatal("wrong number of players")
	}

	mb.anthills = make(global.Anthills)
	for i := 0; i < len(mb.players); i++ {
		mb.area[positions[i][0][0]][positions[i][0][1]] = global.CreateAnthill(mb.players[i])
		mb.anthills.Add(mb.players[i], positions[i][0], &global.Anthill{
			Pos:      positions[i][0],
			User:     mb.players[i],
			BirthPos: positions[i][1],
		})
	}

	mb.ants = make([]*global.Ant, 0, len(mb.players))
	for _, anthills := range mb.anthills {
		for _, anthill := range anthills {
			ant := &global.Ant{
				Pos:    anthill.BirthPos,
				User:   anthill.User,
				IsDead: false,
			}
			mb.ants = append(mb.ants, ant)
			mb.area[anthill.BirthPos.X()][anthill.BirthPos.Y()] = global.CreateAnt(ant)
		}
	}
}

func (mb *MatchBuilder) BuildArea() {
	mb.area = make([][]*global.Object, mb.areaSize)
	lastTile := mb.areaSize - 1
	for x := 0; x < mb.areaSize; x++ {
		mb.area[x] = make([]*global.Object, mb.areaSize)
		for y := 0; y < mb.areaSize; y++ {
			// edges
			if x == 0 || x == lastTile || y == 0 || y == lastTile {
				mb.area[x][y] = global.CreateWall()
			} else {
				mb.area[x][y] = global.CreateEmptyObject()
			}
		}
	}
}

func (mb *MatchBuilder) BuildFood(percentFrom float32, percentTo float32, min int, isUniformDistribution bool) {
	if mb.area == nil || mb.ants == nil {
		log.Fatal("builder must have ants and area before build food")
	}

	randomPercent := rand.Float32()*(percentTo-percentFrom) + percentFrom
	foodCount := int(float32(mb.areaSize*mb.areaSize) * randomPercent)
	if foodCount < min {
		foodCount = min
	}

	if isUniformDistribution {
		mb.foodUniformDistribution(foodCount)
		return
	}

	size := mb.areaSize - 2
	for i := 0; i < foodCount; i++ {
		x := rand.Intn(size) + 1
		y := rand.Intn(size) + 1
		if mb.area[x][y].Type != pkg.AntField {
			mb.area[x][y] = global.CreateFood()
		}
	}
}

func (mb *MatchBuilder) BuildMatch(s global.Storage) *Match {
	if mb.players == nil || mb.ants == nil {
		log.Fatal("builder must have at least players and ants")
	}

	return CreateMatch(mb, s)
}

// htodo symmetrically distribution
func (mb *MatchBuilder) foodUniformDistribution(foodCount int) {
	var xPartSize int
	var yPartSize int
	var offsets [][2]int
	halfSize := int(math.Round(float64(mb.areaSize / 2)))
	antsCount := len(mb.ants)
	switch antsCount {
	case 2:
		// -1 because of walls
		xPartSize = halfSize - 1
		yPartSize = mb.areaSize - 2
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
			if mb.area[x][y].Type != pkg.AntField {
				mb.area[x][y] = global.CreateFood()
			}
		}
	}
}
