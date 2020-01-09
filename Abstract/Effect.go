package Abstract

import (
	"strconv"
)

type Effect struct {
	//+condition
	ID       EffectID
	Duration int
	State    int
	Type     EffectType
	wielder  *Girl
	Activate func(self, opp *Girl, turn int)
	Remove   func(self, opp *Girl, turn int)
}

func (self *Girl) CreateEff(ID EffectID, opp *Girl, duration int, state int) *Effect {
	s := *new(Effect)
	switch ID {
	//Basic
	case DmgMul:
		s = Effect{
			ID,
			duration,
			state,
			Basic,
			self,
			dummy,
			dummy,
		}

	case DmgAdd:
		s = Effect{
			ID,
			duration,
			state,
			Basic,
			self,
			dummy,
			dummy,
		}
		//Prohibiting
	case CantHeal:
		s = Effect{
			ID,
			duration,
			state,
			Prohibiting,
			self,
			dummy,
			dummy,
		}
	case CantUse:
		s = Effect{
			ID,
			duration,
			state,
			Prohibiting,
			self,
			dummy,
			dummy,
		}
		//Control
	case ControlledByStT:
		s = Effect{
			ID,
			duration,
			state,
			Control,
			self,
			dummy,
			dummy,
		}
		//Debuff
	case TurnThreshold:
		s = Effect{
			ID,
			duration,
			state,
			Debuff,
			self,
			dummy,
			dummy,
		}
	case TurnReduc:
		s = Effect{
			ID,
			duration,
			state,
			Debuff,
			self,
			dummy,
			dummy,
		}
		//Buff
	case AtkReduc:
		s = Effect{
			ID,
			duration,
			state,
			Buff,
			self,
			dummy,
			dummy,
		}
		//State
	case Unseen:
		s = Effect{
			ID,
			duration,
			state,
			State,
			self,
			allow,
			dummy,
		}
	case SpedUp:
		s = Effect{
			ID,
			duration,
			state,
			State,
			self,
			dummy,
			dummy,
		}
	case DelayedHeal:
		s = Effect{
			ID,
			duration,
			state,
			State,
			self,
			dummy,
			delayed,
		}
	case Invulnerable:
		s = Effect{
			ID,
			duration,
			state,
			State,
			self,
			dummy,
			dummy,
		}
	case EuphoricHeal:
		s = Effect{
			ID,
			duration,
			state,
			State,
			self,
			dummy,
			euphoricHeal,
		}
		//Numerical
	case GreenToken:
		s = Effect{
			ID,
			duration,
			state,
			Numerical,
			self,
			dummy,
			dummy,
		}
	case BlackToken:
		s = Effect{
			ID,
			duration,
			state,
			Numerical,
			self,
			dummy,
			dummy,
		}
	case StolenHP:
		s = Effect{
			ID,
			duration,
			state,
			Numerical,
			self,
			dummy,
			dummy,
		}
	case BoostShock:
		s = Effect{
			ID,
			duration,
			state,
			Numerical,
			self,
			dummy,
			dummy,
		}
	case BoostLayers:
		s = Effect{
			ID,
			duration,
			state,
			Numerical,
			self,
			dummy,
			dummy,
		}
	case EuphoricSource:
		s = Effect{
			ID,
			duration,
			state,
			Numerical,
			self,
			dummy,
			dummy,
		}

	default:
		panic("*notices ID*... " + strconv.Itoa(int(ID)) + " - OwO what's this?")
	}
	return &s
}

//Equals does not take into consideration the wielder
func (s *Effect) Equals(other *Effect) bool {
	return s.ID == other.ID &&
		s.Duration == other.Duration &&
		s.Type == other.Type &&
		s.State == other.State

}

//DeepEqual takes into consideration the wielder, do not use inside of Girl.Equals
func (s *Effect) DeepEqual(other *Effect) bool {
	return s.ID == other.ID &&
		s.Duration == other.Duration &&
		s.Type == other.Type &&
		s.State == other.State &&
		s.wielder.Equals(other.wielder)
}

func (s *Effect) Copy(newWielder *Girl) *Effect {
	r := Effect{
		ID:       s.ID,
		Duration: s.Duration,
		Type:     s.Type,
		wielder:  newWielder,
		Activate: s.Activate,
		State:    s.State,
		Remove:   s.Remove,
	}
	return &r
}

func (s *Effect) DecreaseDuration() bool {
	if s.Duration--; s.Duration > 0 {
		return false
	}
	return true
}
