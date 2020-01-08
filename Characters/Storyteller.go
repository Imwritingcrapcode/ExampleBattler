package Characters

import (
	. "../Abstract"
	"math"
	"strconv"
)

type Storyteller Girl

func (self *Storyteller) Init() {
	self.Number = 1
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 119
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //-4, Main Colour - Orange
		Red:    -1,
		Orange: 1,
		Yellow: 0,
		Green:  -2,
		Cyan:   -1,
		Blue:   1,
		Violet: 1,
		Pink:   0,
		Gray:   -1,
		Brown:  -1,
		Black:  -2,
		White:  1}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		0, 0, 0,
		Orange,
		"Your Number",
		StorytellerQ,
		"rgb(255,69,002)",
	}
	self.Skills[1] = &Skill{
		false,
		1, 0, 2,
		White,
		"Your Colour",
		StorytellerW,
		"rgb(237,235,243)",
	}
	self.Skills[2] = &Skill{
		false,
		0, 0, 0,
		Violet,
		"Your Dream",
		StorytellerE,
		"rgb(104,022,253)",
	}
	self.Skills[3] = &Skill{
		true,
		1, 0, 13,
		Blue,
		"My Story",
		StorytellerUlti,
		"rgb(29,104,255)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = StTCheck

}

func StTCheck(player, opp *Girl, turn, skill int) bool {
	switch skill {
	case 0:
		return true
	case 1:
		switch opp.Name {
		case "Ruby":
			//return opp.Skills[opp.LastUsed].Colour == Red
			return true
		case "Speed":
			return true
		case "Milana":
			//return !(opp.LastUsed == 3)
			return true
		case "Structure":
			//return !(opp.LastUsed == 3)
			return true
		case "Euphoria":
			return true
		default:
			panic("And everybody wants to know who is that girl: " + opp.Name)
		}
	case 2:
		return !(player.HasEffect(CantHeal) || player.CurrHP == player.MaxHP)
	case 3:
		return true
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}
}

func StorytellerQ(player, opp *Girl, turn int) {
	//DMG := 16 + opp.Number%7
	DMG := 10 + opp.Number%7
	Damage(player, opp, DMG, false, player.Skills[0].Colour)
}
func StorytellerW(player, opp *Girl, turn int) {
	DMG := 15
	//DMG := 0
	Damage(player, opp, DMG, false, opp.Skills[opp.LastUsed].Colour)
	if !opp.HasEffect(Unseen) {
		eff := opp.CreateEff(CantUse, player, 1, int(opp.Skills[opp.LastUsed].Colour))
		opp.AddEffect(eff)
	}

}
func StorytellerE(player, opp *Girl, turn int) {
	var SUB int
	if opp.Number > 83 {
		SUB = 83
	} else {
		SUB = opp.Number
	}
	HEAL := int(math.Ceil(float64(player.MaxHP-SUB) / float64(GetTurnNum(turn))))
	if !player.HasEffect(CantHeal) {
		Heal(player, HEAL)
	}
}
func StorytellerUlti(player, opp *Girl, turn int) {
	eff := opp.CreateEff(ControlledByStT, player, 1, 0)
	opp.AddEffect(eff)
}
