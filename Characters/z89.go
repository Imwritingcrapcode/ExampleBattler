package Characters

import (
	. "../Abstract"
	"strconv"
)

type Z89 Girl

func (self *Z89) Init() {
	self.Name = "z89"
	self.Number = 8
	self.MaxHP = 116
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //3, main color - blue green black?
		Red: 0,
		Orange: 0,
		Yellow: 0,
		Green: 0,
		Cyan: 0,
		Blue: 0,
		Violet: 0,
		Pink: 0,
		Gray: 0,
		Brown: 0,
		Black: 0,
		White: 0}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		1, 0, 0,
		Black,
		"Z89 lower hp",
		Z89Q,
		"rgb(0, 0, 0)",
	}
	self.Skills[1] = &Skill{
		false,
		0, 0, 0,
		Cyan,
		"Z89 eff removal and damage 15 no def",
		Z89W,
		"rgb(0, 255, 255)",
	}
	self.Skills[2] = &Skill{
		false,
		0, 0, 0,
		Green,
		"Green Sphere",
		Z89E,
		"rgb(0, 255, 0)",
	}
	self.Skills[3] = &Skill{
		true,
		2, 0, 11,
		Blue,
		"Z89 R",
		Z89Ulti,
		"rgb(0, 0, 255)",
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
		return true
	case 2:
		return true
	case 3:
		return true
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}

}

func Z89Q(player, opp *Girl, turn int) {
	opp.MaxHP -= 12
	if opp.CurrHP > opp.MaxHP {
		opp.CurrHP = opp.MaxHP
	}
}

func Z89W(player, opp *Girl, turn int) {
	DMG := 0
	for i := 0; i < TOTALEFFECTS; i++ {
		if opp.HasEffect(EffectID(i)) {
			if opp.GetEffect(EffectID(i)).Duration <= 2 {
				opp.RemoveEffect(EffectID(i))
				Damage(player, opp, 15, true, player.Skills[1].Colour)
				DMG += 15
			}
		}
	}
	if DMG == 0 {
		Damage(player,opp, 5, true, player.Skills[1].Colour)
	}
}

func Z89E(player, opp *Girl, turn int) {
	var CURR, MAX, BASEDMG int
	CURR = opp.CurrHP
	MAX = opp.MaxHP
	BASEDMG = 20
	DMG := BASEDMG - (MAX - CURR)
	if DMG > 0 {
		Damage(player, opp, DMG, false, player.Skills[2].Colour)
	}
}

func Z89Ulti(player, opp *Girl, turn int) {
	DMG := 50
	MAXHP := opp.MaxHP
	THRESHOLD := 70
	Damage(player, opp, DMG - (MAXHP - THRESHOLD), false, player.Skills[3].Colour)
}
