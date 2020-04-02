package Characters

import (
	. "../Abstract"
	"strconv"
)

type Ruby Girl

func (self *Ruby) Init() {
	self.Number = 10
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 110
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //+5 Main colour - Red (Yellow?)
		Red:    4,
		Orange: -1,
		Yellow: 2,
		Green:  0,
		Cyan:   1,
		Blue:   -2,
		Violet: 0,
		Pink:   0,
		Gray:   0,
		Brown:  1,
		Black:  0,
		White:  0}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		0, 0, 0,
		Yellow,
		"Dance",
		RubyQ,
		"rgb(255,249,151)",
	}
	self.Skills[1] = &Skill{
		false,
		0, 0, 0,
		Red,
		"Rage",
		RubyW,
		"rgb(255,4,3)",
	}
	self.Skills[2] = &Skill{
		false,
		1, 0, 0,
		Cyan,
		"Stop",
		RubyE,
		"rgb(165,235,240)",
	}
	self.Skills[3] = &Skill{
		true,
		0, 0, 0,
		Red,
		".Execute",
		RubyUlti,
		"rgb(255,4,3)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = RubCheck

}

func RubCheck(player, opp *Girl, turn, skill int) bool {
	switch skill {
	case 0:
		return true
		//return !player.HasEffect(DmgMul) ||
		//			(player.HasEffect(CantUse) && Colour(player.GetEffect(CantUse).State) == Red)
	case 1:
		return true
	case 2:
		return true
	case 3:
		return float64(opp.CurrHP) < float64(opp.MaxHP)*0.1 ||
			(float64(opp.CurrHP) < float64(opp.MaxHP)*0.2 && opp.HasEffect(CantHeal))
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}

}

func RubyQ(player, opp *Girl, turn int) {
	DMGMUL := 2
	DURATION := 3
	eff := player.CreateEff(DmgMul, opp, DURATION, DMGMUL)
	player.AddEffect(eff)
}
func RubyW(player, opp *Girl, turn int) {
	DMG := 24 - 2*GetTurnNum(turn)
	Damage(player, opp, DMG, false, player.Skills[1].Colour)

}
func RubyE(player, opp *Girl, turn int) {
	DURATION := 2
	if !opp.HasEffect(Unseen) {
		eff := opp.CreateEff(CantHeal, player, DURATION, 1)
		opp.AddEffect(eff)
	}
	eff := player.CreateEff(CantHeal, opp, DURATION+1, -1)
	player.AddEffect(eff)
}
func RubyUlti(player, opp *Girl, turn int) {
	THRESHOLD := 10.0 / 100.0
	if player.HasEffect(CantHeal) {
		THRESHOLD = 20.0 / 100.0
	}
	if float64(opp.CurrHP) < float64(opp.MaxHP)*(float64(THRESHOLD)) {
		opp.CurrHP = 0
	}
}
