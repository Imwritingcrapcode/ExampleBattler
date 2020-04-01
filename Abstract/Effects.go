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
	if self.IsAlive() {
		Heal(self, self.MaxHP)
	}
}

func euphoricHeal(self, opp *Girl, turn int) {
	if self.HasEffect(EuphoricSource) {
		var SAUCE = self.GetEffect(EuphoricSource).State
		if SAUCE > 0 {
			Heal(self, SAUCE)
			Heal(opp, SAUCE)
			var DRAIN int
			DRAIN = 9
			if SAUCE >= DRAIN {
				self.GetEffect(EuphoricSource).State -= DRAIN
			} else {
				self.RemoveEffect(EuphoricSource)
			}
		}
	}
}
