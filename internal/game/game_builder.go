package game

import (
	"ants/internal/global"
	"errors"
	pkg "github.com/gregmus2/ants-pkg"
	"log"
	"math"
)

type MatchBuilder struct {
	areaSize int
	ants     []*global.Ant
	players  []*global.User
	area     global.Area
}

func NewMatchBuilder(areaSize int, players []*global.User) (*MatchBuilder, error) {
	if len(players) != 2 || len(players) != 4 {
		return nil, errors.New("wrong number of players")
	}

	return &MatchBuilder{areaSize: areaSize, players: players}, nil
}

func (gb *MatchBuilder) BuildAnts() {
	var positions [][2]uint
	quartSize := uint(math.Round(float64(gb.areaSize / 4)))
	halfSize := uint(math.Round(float64(gb.areaSize / 2)))

	switch len(gb.players) {
	case 2:
		positions = [][2]uint{{quartSize, halfSize}, {uint(gb.areaSize) - quartSize, halfSize}}
		break
	case 4:
		octoSize := uint(math.Round(float64(gb.areaSize / 8)))
		lastOctoPiece := uint(gb.areaSize) - octoSize
		positions = [][2]uint{
			{octoSize, octoSize}, {lastOctoPiece, octoSize},
			{octoSize, lastOctoPiece}, {lastOctoPiece, lastOctoPiece},
		}
		break
	default:
		log.Fatal("wrong number of players")
	}

	gb.ants = make([]*global.Ant, len(gb.players))
	for i := 0; i < len(gb.ants); i++ {
		gb.ants[i] = &global.Ant{
			Pos:    positions[i],
			User:   gb.players[i],
			IsDead: false,
		}
	}
}

func (gb *MatchBuilder) BuildArea() {
	if gb.area == nil || gb.ants == nil {
		log.Fatal("builder must have ants before build area")
	}

	gb.area = make([][]*global.Object, gb.areaSize)
	lastTile := gb.areaSize - 1
	for x := 0; x < gb.areaSize; x++ {
		gb.area[x] = make([]*global.Object, gb.areaSize)
		for y := 0; y < gb.areaSize; y++ {
			if x == 0 || x == lastTile || y == 0 || y == lastTile {
				gb.area[x][y] = global.CreateWall()
			} else {
				gb.area[x][y] = global.CreateEmptyObject()
			}
		}
	}

	for _, ant := range gb.ants {
		gb.area[ant.Pos.X()][ant.Pos.Y()] = global.CreateAnt(ant)
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

func (gb *MatchBuilder) BuildMatch() *Match {
	if gb.players == nil || gb.ants == nil {
		log.Fatal("builder must have at least players and ants")
	}

	return CreateMatch(gb.players, gb.ants, gb.area)
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
		if foodCount-i < 4 {
			break
		}

		for j := 0; j < antsCount; i++ {
			x := global.Random.Intn(xPartSize) + offsets[j][0]
			y := global.Random.Intn(yPartSize) + offsets[j][1]
			if gb.area[x][y].Type != pkg.AntField {
				gb.area[x][y] = global.CreateFood()
			}
		}
	}
}
