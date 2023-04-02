package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

var adjectives = []string{
	"aged",
	"ancient",
	"autumn",
	"billowing",
	"bitter",
	"black",
	"blue",
	"bold",
	"broad",
	"broken",
	"calm",
	"cold",
	"cool",
	"crimson",
	"curly",
	"damp",
	"dark",
	"dawn",
	"delicate",
	"divine",
	"dry",
	"empty",
	"falling",
	"fancy",
	"flat",
	"floral",
	"fragrant",
	"frosty",
	"gentle",
	"green",
	"hidden",
	"holy",
	"icy",
	"jolly",
	"late",
	"lingering",
	"little",
	"lively",
	"long",
	"lucky",
	"misty",
	"morning",
	"muddy",
	"mute",
	"nameless",
	"noisy",
	"odd",
	"old",
	"orange",
	"patient",
	"plain",
	"polished",
	"proud",
	"purple",
	"quiet",
	"rapid",
	"raspy",
	"red",
	"restless",
	"rough",
	"round",
	"royal",
	"shiny",
	"shrill",
	"shy",
	"silent",
	"small",
	"snowy",
	"soft",
	"solitary",
	"sparkling",
	"spring",
	"square",
	"steep",
	"still",
	"summer",
	"super",
	"sweet",
	"throbbing",
	"tight",
	"tiny",
	"twilight",
	"wandering",
	"weathered",
	"wild",
	"winter",
	"wispy",
	"withered",
	"yellow",
	"young",
}

var nouns = []string{
	"art",
	"band",
	"bar",
	"base",
	"bird",
	"block",
	"boat",
	"bonus",
	"bread",
	"breeze",
	"brook",
	"bush",
	"butterfly",
	"cake",
	"cell",
	"cherry",
	"cloud",
	"credit",
	"darkness",
	"dawn",
	"dew",
	"disk",
	"dream",
	"dust",
	"feather",
	"field",
	"fire",
	"firefly",
	"flower",
	"fog",
	"forest",
	"frog",
	"frost",
	"glade",
	"glitter",
	"grass",
	"hall",
	"hat",
	"haze",
	"heart",
	"hill",
	"king",
	"lab",
	"lake",
	"leaf",
	"limit",
	"math",
	"meadow",
	"mode",
	"moon",
	"morning",
	"mountain",
	"mouse",
	"mud",
	"night",
	"paper",
	"pine",
	"poetry",
	"pond",
	"queen",
	"rain",
	"recipe",
	"resonance",
	"rice",
	"river",
	"salad",
	"scene",
	"sea",
	"shadow",
	"shape",
	"silence",
	"sky",
	"smoke",
	"snow",
	"snowflake",
	"sound",
	"star",
	"sun",
	"sunset",
	"surf",
	"term",
	"thunder",
	"tooth",
	"tree",
	"truth",
	"union",
	"unit",
	"violet",
	"voice",
	"water",
	"waterfall",
	"wave",
	"wildflower",
	"wind",
	"wood",
}

func GeneRateUserName(db sqlx.DB) string {
	var name string

	for {
		name = generate()

		fmt.Println(name)

		var id int
		_ = db.Get(&id, "select id from users where username=?", name)

		if id == 0 {
			break
		} else {
			continue
		}

	}

	return name

}

func generate() string {
	a := randomAdjective()
	n := randomNoun()
	i := randomNumber()
	return fmt.Sprintf("%s-%s-%d", a, n, i)
}

func randomAdjective() string {
	index := randomRange(0, len(adjectives)-1)
	return adjectives[index]
}

func randomNoun() string {
	index := randomRange(0, len(nouns)-1)
	return nouns[index]
}

func randomNumber() int {
	return randomRange(10, 99999)
}

func randomRange(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ir := IntRange{min, max}
	return ir.NextRandom(r)
}

// range specification, note that min <= max
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}
