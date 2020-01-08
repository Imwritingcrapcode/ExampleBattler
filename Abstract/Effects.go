package Abstract

func dummy(self, opp *Girl, turn int) {
}

func allow(self, opp *Girl, turn int) {
	for i := 0; i < TOTALEFFECTS; i++ {
		if self.HasEffect(EffectID(i)) && self.GetEffect(EffectID(i)).Type == Prohibiting {
			self.RemoveEffect(EffectID(i))
		}
	}
}

func delayed(self, opp *Girl, turn int) {
	if self.IsAlive() && self.CanHeal() {
		Heal(self, self.MaxHP)
	}
}

func euphoricHeal(self, opp *Girl, turn int) {
	if self.HasEffect(EuphoricSource) && self.GetEffect(EuphoricSource).State > 0 && self.CanHeal() {
		Heal(self, self.GetEffect(EuphoricSource).State)
		if self.GetEffect(EuphoricSource).State >= 10 {
			self.GetEffect(EuphoricSource).State -= 10
		} else {
			self.GetEffect(EuphoricSource).State = 0
		}
	}
}
