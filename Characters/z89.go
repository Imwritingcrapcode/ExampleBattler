package Characters

import (
	. "../Abstract"
	"strconv"
)

type Z89 Girl

func (self *Z89) Init() {
	self.Number = 8
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 111
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //4, main color - blue green black?
		Red: -1,
		Orange: 0,
		Yellow: -1,
		Green: 2,
		Cyan: 3,
		Blue: 2,
		Violet: -1,
		Pink: -1,
		Gray: 0,
		Brown: 0,
		Black: 2,
		White: -1}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		1, 0, 0,
		Black,
		"Scarcity",
		Z89Q,
		"rgb(9, 15, 26)",
	}
	self.Skills[1] = &Skill{
		false,
		2, 0, 3,
		Cyan,
		"Indifference",
		Z89W,
		"rgb(165, 252, 252)",
	}
	self.Skills[2] = &Skill{
		false,
		0, 0, 0,
		Green,
		"Green Sphere",
		Z89E,
		"rgb(16, 143, 38)",
	}
	self.Skills[3] = &Skill{
		true,
		20, 0, 17,
		Blue,
		"Despondency",
		Z89Ulti,
		"rgb(22, 63, 145)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = Z89Check

}

func Z89Check(player, opp *Girl, turn, skill int) bool {
	switch skill {
	case 0:
		return true
	case 1:
		return turn <= opp.Skills[3].StrT
	case 2:
		return true
	case 3:
		return true
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}

}

func Z89Q(player, opp *Girl, turn int) {
	Damage(player, opp, 12, false, player.Skills[0].Colour)
	if opp.CurrHP < opp.MaxHP {
		opp.MaxHP = opp.CurrHP
	}
}

func Z89W(player, opp *Girl, turn int) {
	if turn <= opp.Skills[3].StrT {
		if (opp.Skills[3].StrT + 2) <= 19 {
			opp.Skills[3].StrT += 2
		}
	}
}

func Z89E(player, opp *Girl, turn int) {
	var CURR, MAX, BASEDMG int
	CURR = opp.CurrHP
	MAX = opp.MaxHP
	BASEDMG = 15
	DMG := BASEDMG - (MAX - CURR)
	if DMG > 0 {
		Damage(player, opp, DMG, false, player.Skills[2].Colour)
	}
}

func Z89Ulti(player, opp *Girl, turn int) {
	DMG := 40
	MAXHP := opp.MaxHP
	THRESHOLD := 70
	Damage(player, opp, DMG-(MAXHP-THRESHOLD), false, player.Skills[3].Colour)
}
