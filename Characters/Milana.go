package Characters

import (
	. "../Abstract"
	"strconv"
)

type Milana Girl

func (self *Milana) Init() {
	self.Number = 51
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 114
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //+1, main color - white(green)
		Red:    0,
		Orange: -1,
		Yellow: 1,
		Green:  1,
		Cyan:   1,
		Blue:   0,
		Violet: -1,
		Pink:   0,
		Gray:   0,
		Brown:  0,
		Black:  -2,
		White:  2}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		0, 0, 0,
		Green,
		"Royal Move",
		MilanaQ,
		"rgb(49,255,185)",
	}
	self.Skills[1] = &Skill{
		false,
		0, 0, 0,
		White,
		"Composure",
		MilanaW,
		"rgb(232,255,243)",
	}
	self.Skills[2] = &Skill{
		false,
		2, 0, 0,
		White,
		"Mint Mist",
		MilanaE,
		"rgb(232,255,243)",
	}
	self.Skills[3] = &Skill{
		true,
		0, 0, 15,
		Cyan,
		"Pride",
		MilanaUlti,
		"rgb(115,255,240)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = MlnCheck

}

func MlnCheck(player, opp *Girl, turn, skill int) bool {
	switch skill {
	case 0:
		return true
	case 1:
		return player.CurrHP < player.MaxHP && player.HasEffect(StolenHP) && player.GetEffect(StolenHP).State > 0
	case 2:
		return true
	case 3:
		return player.HasEffect(StolenHP) && player.GetEffect(StolenHP).State > 0
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}

}

func MilanaQ(player, opp *Girl, turn int) {
	var dmg int
	DMGLOW := 12
	DMGHIGH := 20
	if player.HasEffect(Unseen) {
		dmg = Damage(player, opp, DMGHIGH, false, player.Skills[0].Colour)
	} else {
		dmg = Damage(player, opp, DMGLOW, false, player.Skills[0].Colour)
	}
	if player.HasEffect(StolenHP) {
		player.GetEffect(StolenHP).State += dmg
	} else {
		eff := player.CreateEff(StolenHP, opp, 21, dmg)
		player.AddEffect(eff)
	}
}
func MilanaW(player, opp *Girl, turn int) {
	var COST, UCOST, MAXH, MAXUH int
	COST = 6
	MAXH = 20
	UCOST = 10
	MAXUH = 30
	if !player.HasEffect(StolenHP) {
		eff := player.CreateEff(StolenHP, opp, 21, 0)
		player.AddEffect(eff)
	}
	if !player.HasEffect(Unseen) {
		//enough2use
		if player.GetEffect(StolenHP).State >= COST {
			player.GetEffect(StolenHP).State -= COST
			Heal(player, MAXH)
			//not enough 2 use
		} else {
			player.GetEffect(StolenHP).State -= Heal(player, player.GetEffect(StolenHP).State)

		}
		//Unseen
	} else {
		//enough2use
		if player.GetEffect(StolenHP).State >= UCOST {
			player.GetEffect(StolenHP).State -= UCOST
			Heal(player, MAXUH)
			//JUST A BIT
		} else {
			Heal(player, player.GetEffect(StolenHP).State)
			player.GetEffect(StolenHP).State -= Heal(player, player.GetEffect(StolenHP).State)

		}
	}
}
func MilanaE(player, opp *Girl, turn int) {
	eff := player.CreateEff(Unseen, opp, 3, 0)
	player.AddEffect(eff)
}
func MilanaUlti(player, opp *Girl, turn int) {
	if player.HasEffect(StolenHP) {
		Damage(player, opp, player.GetEffect(StolenHP).State, false, player.Skills[3].Colour)
		player.GetEffect(StolenHP).State = 0
	}
}
