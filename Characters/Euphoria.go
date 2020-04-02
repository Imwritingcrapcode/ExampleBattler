package Characters

import (
	. "../Abstract"
	"strconv"
)

type Euphoria Girl

func (self *Euphoria) Init() {
	self.Number = 9
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 117
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //-2, main color - pink(orange)
		Red: 0,
		Orange: 2,
		Yellow: 0,
		Green: 0,
		Cyan: -2,
		Blue: 0,
		Violet: 0,
		Pink: 3,
		Gray: 0,
		Brown: -3,
		Black: 0,
		White: 0}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		1, 0, 0,
		Orange,
		"Ampleness",
		EuphoriaQ,
		"rgb(255,173,135)",
	}
	self.Skills[1] = &Skill{
		false,
		2, 0, 0,
		Orange,
		"Exuberance",
		EuphoriaW,
		"rgb(255,173,135)",
	}
	self.Skills[2] = &Skill{
		false,
		0, 0, 0,
		Pink,
		"Pink Sphere",
		EuphoriaE,
		"rgb(255,135,173)",
	}
	self.Skills[3] = &Skill{
		true,
		20, 0, 7,
		Pink,
		"Euphoria",
		EuphoriaUlti,
		"rgb(255,135,175)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = EupCheck

}

func EupCheck(player, opp *Girl, turn, skill int) bool {
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

func EuphoriaQ(player, opp *Girl, turn int) {
	AMNT := 12
	if player.HasEffect(EuphoricSource) {
		player.GetEffect(EuphoricSource).State += AMNT
	} else {
		eff := player.CreateEff(EuphoricSource, opp, 21, AMNT)
		player.AddEffect(eff)
	}
	player.MaxHP += AMNT
	opp.MaxHP += AMNT
}

func EuphoriaW(player, opp *Girl, turn int) {
	var AMNT int
	if turn <= opp.Skills[3].StrT {
		if (opp.Skills[3].StrT - 2) >= 0 {
			opp.Skills[3].StrT -= 2
		}
		AMNT = 10
	} else {
		AMNT = 20
	}
		if player.HasEffect(EuphoricSource) {
			player.GetEffect(EuphoricSource).State += AMNT
		} else {
			eff := player.CreateEff(EuphoricSource, opp, 21, AMNT)
			player.AddEffect(eff)
		}
		player.MaxHP += AMNT
		opp.MaxHP += AMNT
		Heal(player, AMNT)
}

func EuphoriaE(player, opp *Girl, turn int) {
	DMG := 12
	Damage(player, opp, DMG, false, player.Skills[2].Colour)
	player.MaxHP += DMG
	opp.MaxHP += DMG

}

func EuphoriaUlti(player, opp *Girl, turn int) {
	if player.HasEffect(EuphoricSource) && !player.HasEffect(EuphoricHeal) {
		eff := player.CreateEff(EuphoricHeal, opp, 21, 0)
		player.AddEffect(eff)
	}
}
