package Abstract

type Skill struct {
	IsUlti           bool
	CD, CurrCD, StrT int
	Colour
	Name       string
	Use        func(player, opp *Girl, turn int)
	ColourCode string
}

func (s *Skill) Equals(other *Skill) bool {
	return other.IsUlti == s.IsUlti &&
		other.CD == s.CD &&
		other.CurrCD == s.CurrCD &&
		other.StrT == s.StrT &&
		other.Colour == s.Colour &&
		other.Name == s.Name && other.ColourCode == s.ColourCode
}

func (s *Skill) Copy() *Skill {
	r := Skill{
		IsUlti:     s.IsUlti,
		CD:         s.CD,
		CurrCD:     s.CurrCD,
		StrT:       s.StrT,
		Colour:     s.Colour,
		Name:       s.Name,
		Use:        s.Use,
		ColourCode: s.ColourCode,
	}
	return &r
}
