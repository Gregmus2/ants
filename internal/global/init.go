package global

import (
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var Random *rand.Rand
var Config *ConfigType

func init() {
	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/ants/.env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Dotenv loaded successfully")

	areaSize, err := strconv.ParseInt(os.Getenv("AREA_SIZE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	matchPartsLimit, err := strconv.ParseInt(os.Getenv("MATCH_PARTS_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	matchPartSize, err := strconv.ParseInt(os.Getenv("MATCH_PART_SIZE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	Config = &ConfigType{
		AreaSize:        int(areaSize),
		MatchPartsLimit: int(matchPartsLimit),
		MatchPartSize:   int(matchPartSize),
	}

	Random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
