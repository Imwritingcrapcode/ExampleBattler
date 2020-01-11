package Abstract

import (
	"reflect"
)

type Girl struct {
	Name                  string
	CurrHP, MaxHP, Number int
	Defenses              map[Colour]int
	LastUsed              int
	LastDmgTaken          int
	Skills                []*Skill
	Effects               EffectSet
	CheckifAppropriate    func(player, opp *Girl, turn, skill int) bool
}

func (self *Girl) Equals(other *Girl) bool {
	res := self.Name == other.Name &&
		self.Number == other.Number &&
		self.MaxHP == other.MaxHP &&
		self.CurrHP == other.CurrHP &&
		self.LastDmgTaken == other.LastDmgTaken
	reflect.DeepEqual(self.Defenses, other.Defenses)
	if len(self.Skills) != len(other.Skills) || !res {
		return false
	}
	for i := range self.Skills {
		res = res && self.Skills[i].Equals(other.Skills[i])
	}
	for i := 0; i < TOTALEFFECTS; i++ {
		if self.HasEffect(EffectID(i)) && other.HasEffect(EffectID(i)) {
			res = res && self.GetEffect(EffectID(i)).Equals(other.GetEffect(EffectID(i)))
		} else {
			return false
		}
	}
	res = res && (self.Skills[self.LastUsed].Equals(other.Skills[other.LastUsed]))
	return res
}

func (self *Girl) Copy() CharInt {
	var cope Girl
	cope = *new(Girl)
	targetMap := make(map[Colour]int)
	for key, value := range self.Defenses {
		targetMap[key] = value
	}
	newEffs := *new(EffectSet)
	newEffs.Init(TOTALEFFECTS)
	for i := 0; i < TOTALEFFECTS; i++ {
		if self.Effects.Contains(EffectID(i)) {
			oldEff := self.GetEffect(EffectID(i))
			newEff := oldEff.Copy(&cope)
			newEffs.Set(newEff)
		}
	}
	newSkills := make([]*Skill, len(self.Skills))
	for i := range self.Skills {
		newSkills[i] = self.Skills[i].Copy()
	}
	cope = Girl{
		Name:               self.Name,
		Number:             self.Number,
		CurrHP:             self.CurrHP,
		MaxHP:              self.MaxHP,
		Defenses:           targetMap,
		Skills:             newSkills,
		Effects:            newEffs,
		LastUsed:           self.LastUsed,
		LastDmgTaken:       self.LastDmgTaken,
		CheckifAppropriate: self.CheckifAppropriate,
	}

	return &cope
}

func (self *Girl) AddEffect(effect *Effect) {
	self.Effects.Set(effect)
}
func (self *Girl) RemoveEffect(ID EffectID) {
	self.Effects.Remove(ID)
}
func (self *Girl) HasEffect(ID EffectID) bool {
	return self.Effects.Contains(ID)
}
func (self *Girl) GetEffect(ID EffectID) *Effect {
	return self.Effects.Get(ID)
}

func (self *Girl) ModifyDef(colour Colour, amnt int) {
	if self.Defenses[colour]+amnt <= 5 && self.Defenses[colour]+amnt >= -5 {
		self.Defenses[colour] += amnt

	} else if amnt > 0 {
		self.Defenses[colour] = 5
	} else {
		self.Defenses[colour] = -5
	}
}

//Skills of the girls
func (self *Girl) SkillQ(opp *Girl, turn int) {
	self.Skills[0].Use(self, opp, turn)
}

func (self *Girl) SkillW(opp *Girl, turn int) {
	self.Skills[1].Use(self, opp, turn)
}
func (self *Girl) SkillE(opp *Girl, turn int) {
	self.Skills[2].Use(self, opp, turn)
}
func (self *Girl) SkillUlti(opp *Girl, turn int) {
	self.Skills[3].Use(self, opp, turn)
}

func (self *Girl) IsAlive() bool {
	return self.CurrHP > 0
}

func (self *Girl) CanHeal() bool {
	return !self.HasEffect(CantHeal)
}

func (self *Girl) TurnEnd(opp *Girl) {

}

func (self *Girl) Init() {

}
func (self *Girl) DecreaseCooldowns() {
	for _, k := range self.Skills {
		if k.CurrCD > 0 {
			k.CurrCD--
		}
	}
}
func (player *Girl) DecreaseEffects(opp *Girl, turn int) {
	for i := 0; i < TOTALEFFECTS; i++ {
		if player.HasEffect(EffectID(i)) {
			if player.GetEffect(EffectID(i)).DecreaseDuration() {
				//fmt.Println("REMOVED", EffectNames[EffectID(i)])
				player.GetEffect(EffectID(i)).Remove(player, opp, turn)
				player.RemoveEffect(EffectID(i))

			}
		}
	}
}

func (player *Girl) CheckAvailableSkills(turnnum int) []bool {
	skillsAvailable := make([]bool, len(player.Skills))
	for i, skill := range player.Skills {
		if skill.CurrCD == 0 &&
			skill.StrT <= turnnum &&
			!(player.HasEffect(CantUse) && Colour(player.GetEffect(CantUse).State) == skill.Colour) &&
			!(player.HasEffect(SpedUp) && skill.IsUlti) {
			skillsAvailable[i] = true

		} else {
			skillsAvailable[i] = false
		}
	}
	return skillsAvailable
}

func (girl *Girl) SkillsStringList() []string {
	var res []string
	for _, skill := range girl.Skills {
		res = append(res, skill.Name)
	}
	return res
}

func (girl *Girl) SkillColoursToString() []string {
	var res []string
	for _, skill := range girl.Skills {
		res = append(res, ColoursToString[skill.Colour])
	}
	return res
}

func (girl *Girl) SkillColourCodesToString() []string {
	var res []string
	for _, skill := range girl.Skills {
		res = append(res, skill.ColourCode)
	}
	return res
}