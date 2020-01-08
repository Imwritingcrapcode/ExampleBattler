package Characters

import (
	. "../Abstract"
	"strconv"
)

type Speed Girl

func (self *Speed) Init() {
	self.Number = 33
	self.Name = ReleasedCharactersNames[self.Number]
	self.MaxHP = 113
	self.CurrHP = self.MaxHP
	self.Defenses = map[Colour]int{ //+2, main colour - Green
		Red:    0,
		Orange: 0,
		Yellow: 0,
		Green:  4,
		Cyan:   0,
		Blue:   0,
		Violet: -2,
		Pink:   0,
		Gray:   0,
		Brown:  0,
		Black:  2,
		White:  -2}
	self.Skills = make([]*Skill, 4)
	self.Skills[0] = &Skill{
		false,
		0, 0, 0,
		Green,
		"Run",
		SpeedQ,
		"rgb(14,51,20)",
	}
	self.Skills[1] = &Skill{
		false,
		0, 0, 0,
		Black,
		"Weaken",
		SpeedW,
		"rgb(0,10,0)",
	}
	self.Skills[2] = &Skill{
		false,
		0, 0, 0,
		Green,
		"Speed",
		SpeedE,
		"rgb(14,51,20)",
	}
	self.Skills[3] = &Skill{
		true,
		0, 0, 0,
		Black,
		"Stab",
		SpeedUlti,
		"rgb(0,10,0)",
	}
	self.Effects = EffectSet{}
	self.Effects.Init(TOTALEFFECTS)
	self.CheckifAppropriate = SpdCheck

}

func SpdCheck(player, opp *Girl, turn, skill int) bool {
	switch skill {
	case 0:
		return !(player.HasEffect(GreenToken) && player.GetEffect(GreenToken).State > 4)
	case 1:
		return !(player.HasEffect(BlackToken) && player.GetEffect(BlackToken).State > 4)
	case 2:
		return !((player.HasEffect(GreenToken) && player.GetEffect(GreenToken).State > 4) &&
			(player.HasEffect(BlackToken) && player.GetEffect(BlackToken).State > 4)) &&
			!(player.HasEffect(SpedUp) && player.GetEffect(SpedUp).Duration == 2)
	case 3:
		return player.HasEffect(BlackToken) || player.HasEffect(GreenToken)
	default:
		panic("HOW DID YOU EVEN " + strconv.Itoa(skill))

	}

}

func SpeedQ(player, opp *Girl, turn int) {
	DMGREDATK := 5
	GETGTKNS := 1
	if !player.HasEffect(GreenToken) {
		eff := player.CreateEff(GreenToken, opp, 21, GETGTKNS)
		player.AddEffect(eff)
	} else if player.GetEffect(GreenToken).State < 5 {
		player.GetEffect(GreenToken).State += GETGTKNS
	}
	if opp.HasEffect(AtkReduc) {
		opp.GetEffect(AtkReduc).State += DMGREDATK
	} else {
		eff := opp.CreateEff(AtkReduc, player, 21, DMGREDATK)
		opp.AddEffect(eff)
	}
}

func SpeedW(player, opp *Girl, turn int) {
	GETBTKNS := 1
	COLOUR := Green
	REDUCE := 1
	if !player.HasEffect(BlackToken) {
		eff := player.CreateEff(BlackToken, opp, 21, GETBTKNS)
		player.AddEffect(eff)
	} else if player.GetEffect(BlackToken).State < 5 {
		player.GetEffect(BlackToken).State += GETBTKNS
	}
	opp.ModifyDef(COLOUR, -REDUCE)
}
func SpeedE(player, opp *Girl, turn int) {
	GETGTKNS := 1
	if !player.HasEffect(GreenToken) {
		eff := player.CreateEff(GreenToken, opp, 21, GETGTKNS)
		player.AddEffect(eff)
	} else if player.GetEffect(GreenToken).State < 5 {
		player.GetEffect(GreenToken).State += GETGTKNS
	}
	eff2 := player.CreateEff(SpedUp, opp, 2, 0)
	player.AddEffect(eff2)
}
func SpeedUlti(player, opp *Girl, turn int) {
	DMGPERTKN := 6
	var green, black int
	if player.HasEffect(GreenToken) {
		green = player.GetEffect(GreenToken).State

	} else {
		green = 0
	}
	if player.HasEffect(BlackToken) {
		black = player.GetEffect(BlackToken).State

	} else {
		black = 0
	}
	var dmgG, dmgB, atkred, turnthr, mul, add int
	if player.HasEffect(DmgMul) {
		mul = player.GetEffect(DmgMul).State
	} else {
		mul = 1
	}
	if player.HasEffect(DmgAdd) {
		add = player.GetEffect(DmgAdd).State
	} else {
		add = 0
	}
	if player.HasEffect(AtkReduc) {
		atkred = player.GetEffect(AtkReduc).State
	} else {
		atkred = 0
	}
	if player.HasEffect(TurnThreshold) {
		turnthr = player.GetEffect(TurnThreshold).State
	} else {
		turnthr = 0
	}
	dmgG = DMGPERTKN*green - opp.Defenses[Green]
	dmgB = DMGPERTKN*black - opp.Defenses[Black]
	dmg := (dmgG+dmgB)*mul + add - atkred
	if dmg > 0 && dmg > turnthr && !opp.HasEffect(Invulnerable) && (player.HasEffect(GreenToken) || player.HasEffect(BlackToken)) {
		opp.CurrHP -= dmg
		opp.LastDmgTaken = dmg
	} else {
		opp.LastDmgTaken = 0
	}
}
