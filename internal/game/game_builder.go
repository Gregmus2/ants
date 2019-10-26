package game

import (
	"ants/internal/global"
	"errors"
	"log"
	"math"

	pkg "github.com/gregmus2/ants-pkg"
)

type MatchBuilder struct {
	name     string
	areaSize int
	ants     []*global.Ant
	players  []*global.User
	area     global.Area
	anthills map[*global.User][]global.Anthill
}

func NewMatchBuilder(name string, areaSize int, players []*global.User) (*MatchBuilder, error) {
	if len(players) != 2 && len(players) != 4 {
		return nil, errors.New("wrong number of players")
	}

	return &MatchBuilder{name: name, areaSize: areaSize, players: players}, nil
}

func (gb *MatchBuilder) BuildAnts() {
	if gb.area == nil {
		log.Fatal("builder must have area before build ants")
	}

	// [players][position, birth position]
	var positions [][2][2]uint
	quartSize := uint(math.Round(float64(gb.areaSize / 4)))
	halfSize := uint(math.Round(float64(gb.areaSize / 2)))

	switch len(gb.players) {
	case 2:
		positions = [][2][2]uint{
			{{quartSize, halfSize}, {quartSize + 1, halfSize}},
			{{uint(gb.areaSize) - quartSize, halfSize}, {uint(gb.areaSize) - quartSize - 1, halfSize}},
		}
		break
	case 4:
		octoSize := uint(math.Round(float64(gb.areaSize / 8)))
		lastOctoPiece := uint(gb.areaSize) - octoSize
		positions = [][2][2]uint{
			{{octoSize, octoSize}, {octoSize + 1, octoSize + 1}},
			{{lastOctoPiece, octoSize}, {lastOctoPiece - 1, octoSize + 1}},
			{{octoSize, lastOctoPiece}, {octoSize + 1, lastOctoPiece - 1}},
			{{lastOctoPiece, lastOctoPiece}, {lastOctoPiece - 1, lastOctoPiece - 1}},
		}
		break
	default:
		log.Fatal("wrong number of players")
	}

	gb.anthills = make(map[*global.User][]global.Anthill)
	for i := 0; i < len(gb.players); i++ {
		gb.area[positions[i][0][0]][positions[i][0][1]] = global.CreateAnthill(gb.players[i])
		gb.anthills[gb.players[i]] = append(gb.anthills[gb.players[i]], global.Anthill{
			Pos:      positions[i][0],
			User:     gb.players[i],
			BirthPos: positions[i][1],
		})
	}

	gb.ants = make([]*global.Ant, 0, len(gb.players))
	for _, anthill := range gb.anthills {
		ant := &global.Ant{
			Pos:    anthill[0].BirthPos,
			User:   anthill[0].User,
			IsDead: false,
		}
		gb.ants = append(gb.ants, ant)
		gb.area[anthill[0].BirthPos.X()][anthill[0].BirthPos.Y()] = global.CreateAnt(ant)
	}
}

func (gb *MatchBuilder) BuildArea() {
	gb.area = make([][]*global.Object, gb.areaSize)
	lastTile := gb.areaSize - 1
	for x := 0; x < gb.areaSize; x++ {
		gb.area[x] = make([]*global.Object, gb.areaSize)
		for y := 0; y < gb.areaSize; y++ {
			// edges
			if x == 0 || x == lastTile || y == 0 || y == lastTile {
				gb.area[x][y] = global.CreateWall()
			} else {
				gb.area[x][y] = global.CreateEmptyObject()
			}
		}
	}
}

func (gb *MatchBuilder) BuildFood(percentFrom float32, percentTo float32, min int, isUniformDistribution bool) {
	if gb.area == nil || gb.ants == nil {
		log.Fatal("builder must have ants and area before build food")
	}

	randomPercent := global.Random.Float32()*(percentTo-percentFrom) + percentFrom
	foodCount := int(float32(gb.areaSize*gb.areaSize) * randomPercent)
	if foodCount < min {
		foodCount = min
	}

	if isUniformDistribution {
		gb.foodUniformDistribution(foodCount)
		return
	}

	size := gb.areaSize - 2
	for i := 0; i < foodCount; i++ {
		x := global.Random.Intn(size) + 1
		y := global.Random.Intn(size) + 1
		if gb.area[x][y].Type != pkg.AntField {
			gb.area[x][y] = global.CreateFood()
		}
	}
}

func (gb *MatchBuilder) BuildMatch(s global.Storage) *Match {
	if gb.players == nil || gb.ants == nil {
		log.Fatal("builder must have at least players and ants")
	}

	return CreateMatch(gb.name, gb.players, gb.ants, gb.anthills, gb.area, s)
}

func (gb *MatchBuilder) foodUniformDistribution(foodCount int) {
	var xPartSize int
	var yPartSize int
	var offsets [][2]int
	halfSize := int(math.Round(float64(gb.areaSize / 2)))
	antsCount := len(gb.ants)
	switch antsCount {
	case 2:
		// -1 because of walls
		xPartSize = halfSize - 1
		yPartSize = gb.areaSize - 2
		offsets = [][2]int{{1, 1}, {halfSize, 1}}
		break
	case 4:
		xPartSize = halfSize - 1
		yPartSize = halfSize - 1
		offsets = [][2]int{{1, 1}, {halfSize, 1}, {1, halfSize}, {halfSize, halfSize}}
		break
	default:
		log.Fatal("wrong number of ants")
	}

	for i := 0; i < foodCount; i += antsCount {
		if foodCount-i < antsCount {
			break
		}

		for j := 0; j < antsCount; j++ {
			x := global.Random.Intn(xPartSize) + offsets[j][0]
			y := global.Random.Intn(yPartSize) + offsets[j][1]
			if gb.area[x][y].Type != pkg.AntField {
				gb.area[x][y] = global.CreateFood()
			}
		}
	}
}
