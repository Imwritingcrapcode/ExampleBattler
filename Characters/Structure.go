package Characters

import (
	. "../Abstract"
	"strconv"
)

type Structure Girl

func (self *Structure) Init() {
	self.Number = 119
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 119
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //-4, main colour - gray
		Red:    -2,
		Orange: -1,
		Yellow: -2,
		Green:  0,
		Cyan:   2,
		Blue:   -2,
		Violet: 1,
		Pink:   0,
		Gray:   2,
		Brown:  0,
		Black:  -3,
		White:  1}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		0, 0, 0,
		Cyan,
		"E-Shock",
		StructureQ,
		"rgb(148,255,246)",
	}
	self.Skills[1] = &Skill{
		false,
		0, 0, 0,
		Violet,
		"I Boost",
		StructureW,
		"rgb(90,0,255)",
	}
	self.Skills[2] = &Skill{
		false,
		0, 0, 0,
		Gray,
		"S Layers",
		StructureE,
		"rgb(157,168,168)",
	}
	self.Skills[3] = &Skill{
		true,
		4, 0, 13,
		White,
		"Last Chance",
		StructureUlti,
		"rgb(240,245,252)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = StrCheck

}

func StrCheck(player, opp *Girl, turn, skill int) bool {
	switch skill {
	case 0:
		return true
	case 1:
		return true
	case 2:
		/*for colour, value := range player.Defences {
			if colour != Black {
				res = res && (value < 5)
			}
		}*/
		return true
	case 3:
		return true
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}

}

func StructureQ(player, opp *Girl, turn int) {
	//BASEDMG := 15
	BASEDMG := 10
	if player.HasEffect(BoostShock) {
		Damage(player, opp, BASEDMG+player.GetEffect(BoostShock).State, false, player.Skills[0].Colour)
	} else {
		Damage(player, opp, BASEDMG, false, player.Skills[0].Colour)
	}
}
func StructureW(player, opp *Girl, turn int) {
	var SHOCKB, LAYERSB int
	SHOCKB = 5
	LAYERSB = 5
	if player.HasEffect(BoostShock) {
		player.GetEffect(BoostShock).State += SHOCKB
		if player.GetEffect(BoostShock).State == 15 {
			for i, skill := range player.Skills {
				if skill.Name == "I Boost" {
					player.Skills[i].CD = 10
				}
			}
		}
	} else {
		eff := player.CreateEff(BoostShock, opp, 21, SHOCKB)
		player.AddEffect(eff)
	}
	if player.HasEffect(BoostLayers) {
		player.GetEffect(BoostLayers).State += LAYERSB
	} else {
		eff := player.CreateEff(BoostLayers, opp, 21, LAYERSB)
		player.AddEffect(eff)
	}
}
func StructureE(player, opp *Girl, turn int) {
	STARTINGTHRESHOLD := 5
	for k := range player.Defenses {
		if k != Black {
			player.ModifyDef(k, 1)
		}
	}
	if player.HasEffect(BoostLayers) {
		eff := opp.CreateEff(TurnThreshold, player, 1, player.GetEffect(BoostLayers).State+STARTINGTHRESHOLD)
		opp.AddEffect(eff)
	} else {
		eff := opp.CreateEff(TurnThreshold, player, 1, STARTINGTHRESHOLD)
		opp.AddEffect(eff)
	}
}
func StructureUlti(player, opp *Girl, turn int) {
	eff := player.CreateEff(DelayedHeal, opp, 2, 0)
	player.AddEffect(eff)
}
