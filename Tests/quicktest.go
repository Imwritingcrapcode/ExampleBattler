package main

import (
	"math/rand"
	"fmt"
)


var ReleasedCharactersPacks = map[string][]int{
	"ST": {},
	"AD": {8, 9, 33},
	"SP": {10, 51},
	"RP": {},
	"LF": {1, 119},
}

func GenerateStartingPack(currTime int64) (int, int) {
	rand.Seed(currTime)
	var free_girl1, free_girl2, n int
	var Rarity1 string
	ST_CHANCE := 450
	AD_CHANCE := 300
	SP_CHANCE := 150
	RP_CHANCE := 80
	LF_CHANCE := 20
	for free_girl1 == 0 {
		n = rand.Intn(1000)
		if (n < ST_CHANCE) && len(ReleasedCharactersPacks["ST"]) > 0 {
			Rarity1 = "ST"
			free_girl1 = ReleasedCharactersPacks["ST"][rand.Intn(len(ReleasedCharactersPacks["ST"]))]
		} else if (n < ST_CHANCE+AD_CHANCE) && len(ReleasedCharactersPacks["AD"]) > 0 {
			Rarity1 = "AD"
			free_girl1 = ReleasedCharactersPacks["AD"][rand.Intn(len(ReleasedCharactersPacks["AD"]))]
		} else if (n < ST_CHANCE+AD_CHANCE+SP_CHANCE) && len(ReleasedCharactersPacks["SP"]) > 0 {
			Rarity1 = "SP"
			free_girl1 = ReleasedCharactersPacks["SP"][rand.Intn(len(ReleasedCharactersPacks["SP"]))]
		} else if (n < ST_CHANCE+AD_CHANCE+SP_CHANCE+RP_CHANCE) && len(ReleasedCharactersPacks["RP"]) > 0 {
			Rarity1 = "RP"
			free_girl1 = ReleasedCharactersPacks["RP"][rand.Intn(len(ReleasedCharactersPacks["RP"]))]
		} else if len(ReleasedCharactersPacks["LF"]) > 0 {
			Rarity1 = "LF"
			free_girl1 = ReleasedCharactersPacks["LF"][rand.Intn(len(ReleasedCharactersPacks["LF"]))]
		}
	}
	for free_girl2 == 0 || free_girl2 == free_girl1 {
		n = rand.Intn(1000)
		if (n >= 1000-LF_CHANCE) && len(ReleasedCharactersPacks["LF"]) > 0 && Rarity1 == "ST" && len(ReleasedCharactersPacks["LF"]) > 0 {
			free_girl2 = ReleasedCharactersPacks["LF"][rand.Intn(len(ReleasedCharactersPacks["LF"]))]
		} else if (n >= 1000-LF_CHANCE-RP_CHANCE) && (Rarity1 == "ST" || Rarity1 == "AD") && len(ReleasedCharactersPacks["RP"]) > 0 {
			free_girl2 = ReleasedCharactersPacks["RP"][rand.Intn(len(ReleasedCharactersPacks["RP"]))]
		} else if (n >= 1000-LF_CHANCE-RP_CHANCE-SP_CHANCE) && (Rarity1 == "ST" || Rarity1 == "AD" || Rarity1 == "SP") && len(ReleasedCharactersPacks["SP"]) > 0 {
			free_girl2 = ReleasedCharactersPacks["SP"][rand.Intn(len(ReleasedCharactersPacks["SP"]))]
		} else if (n >= 1000-LF_CHANCE-RP_CHANCE-SP_CHANCE-AD_CHANCE) && (Rarity1 != "LF") && len(ReleasedCharactersPacks["AD"]) > 0 {
			free_girl2 = ReleasedCharactersPacks["AD"][rand.Intn(len(ReleasedCharactersPacks["AD"]))]
		} else if len(ReleasedCharactersPacks["ST"]) > 0 {
			free_girl2 = ReleasedCharactersPacks["ST"][rand.Intn(len(ReleasedCharactersPacks["ST"]))]
		} else {
			free_girl2 = ReleasedCharactersPacks["AD"][rand.Intn(len(ReleasedCharactersPacks["AD"]))]
		}
	}
	/*free_girl1 = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	free_girl2 = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	for free_girl2 == free_girl1 {
		free_girl2 = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	}*/

	return free_girl1, free_girl2
}

func main() {
	fmt.Println(GenerateStartingPack(119))
}
