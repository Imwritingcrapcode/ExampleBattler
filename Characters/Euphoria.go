package Characters

import (
	. "../Abstract"
	"strconv"
)

type Euphoria Girl

func (self *Euphoria) Init() {
	self.Name = "Euphoria"
	self.Number = 9
	self.MaxHP = 112
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //3, main color - pink(orange)
		Red:    0,
		Orange: 0,
		Yellow: 0,
		Green:  0,
		Cyan:   0,
		Blue:   0,
		Violet: 0,
		Pink:   0,
		Gray:   0,
		Brown:  0,
		Black:  0,
		White:  0}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		1, 0, 0,
		Orange,
		"High Spirits",
		EuphoriaQ,
		"rgb(255,173,135)",
	}
	self.Skills[1] = &Skill{
		false,
		2, 0, 3,
		Orange,
		"Unstudied",
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
		3, 0, 9,
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
	if player.HasEffect(EuphoricSource) {
		player.GetEffect(EuphoricSource).State += 9
	} else {
		eff := player.CreateEff(EuphoricSource, opp, 21, 9)
		player.AddEffect(eff)
	}
	Heal(player, 9)
	Heal(opp, 9)

}

func EuphoriaW(player, opp *Girl, turn int) {
	if turn < opp.Skills[3].StrT {
		opp.Skills[3].StrT += 2
		if player.HasEffect(EuphoricSource) {
			player.GetEffect(EuphoricSource).State += 16
		} else {
			eff := player.CreateEff(EuphoricSource, opp, 21, 16)
			player.AddEffect(eff)
		}
	}
}

func EuphoriaE(player, opp *Girl, turn int) {
	//fmt.Println("E, dmg with an addition")
	var DMG, ADDITIONAL int
	DMG = 5
	if player.HasEffect(EuphoricSource) {
		ADDITIONAL = player.GetEffect(EuphoricSource).State
	} else {
		ADDITIONAL = 0
	}
	Damage(player, opp, DMG+ADDITIONAL, false, player.Skills[2].Colour)
}

func EuphoriaUlti(player, opp *Girl, turn int) {
	if player.HasEffect(EuphoricSource) && !player.HasEffect(EuphoricHeal) {
		eff := player.CreateEff(EuphoricHeal, opp, 4, 0)
		player.AddEffect(eff)
	} /* else if player.HasEffect(EuphoricSource) && player.HasEffect(EuphoricHeal) {
		player.GetEffect(EuphoricSource).State *= 2
		eff := player.CreateEff(EuphoricHeal, opp, 3, 0)
		player.AddEffect(eff)
	}*/
}
