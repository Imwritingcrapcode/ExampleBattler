package Abstract

import (
	"strconv"
	"strings"
)

const MAXTURN = 6

var Time = 0

var VerticesSeen = 0

var Win1, Win2, Draw Graph

var Combinations = 0

var GraphColours = map[Colour]string{
	Red:    "red",
	Orange: "\"#ff87ad\"",
	Yellow: "\"#71ffd1\"",
	Green:  "green",
	Cyan:   "cyan",
	Blue:   "blue",
	Violet: "blueviolet",
	Pink:   "pink",
	Gray:   "\"#8ca1a7\"",
	Brown:  "brown",
	Black:  "black",
	White:  "lightgray",
}

func (v *Graph) InitWinVert(label string) {
	v.Name = TimeName()
	v.Label = label
	v.Seen = false
}

//give vertex a name in 26-based num of letters.
func TimeName() string {
	var name string
	name = recursive(VerticesSeen)
	VerticesSeen += 1
	return name
}

func recursive(left int) string {
	var time string
	var str []string
	ost := left % 26
	if left >= 26 {
		str = append(str, recursive(left/26))
	}
	str = append(str, string(rune(ost+'a')))
	time = strings.Join(str, "")
	return time
}

//The one we use for tests.
type Graph struct {
	Seen        bool
	Turnnum     int
	Label, Name string
	G1, G2      *Girl
	Edges       []*Graph
}

func (g *Graph) ToString() string {
	var builder strings.Builder
	builder.WriteString("diGraph ")
	builder.WriteString(g.G1.Name)
	builder.WriteString(g.G2.Name)
	builder.WriteString(" {\n")
	builder.WriteString(inner(g))
	builder.WriteString("}")
	return builder.String()

}

func inner(self *Graph) string {
	var builder strings.Builder
	//generating label, writing name & label
	var s []string
	if self.G1 == nil || self.G2 == nil {
		if !self.Seen {
			s = []string{"[label=\"", self.Label, "\", color=", GraphColours[Gray], ",style=filled];\n"}
			self.Label = strings.Join(s, "")
			s = []string{"    ", self.Name, " ", self.Label}
			builder.WriteString(strings.Join(s, ""))
			self.Seen = true
		} else {
			s = []string{"    ", self.Name, " ", self.Label}
			builder.WriteString(strings.Join(s, ""))
		}

	} else {
		hp1 := strconv.Itoa(self.G1.CurrHP)
		hp2 := strconv.Itoa(self.G2.CurrHP)
		//var Q, W, E, R string
		var colour1, colour2, colour3, colour4 Colour
		//Q = self.G1.Skills[0].Name
		colour1 = self.G1.Skills[0].Colour
		//W = self.G1.Skills[1].Name
		colour2 = self.G1.Skills[1].Colour
		//E = self.G1.Skills[2].Name
		colour3 = self.G1.Skills[2].Colour
		//R = self.G1.Skills[3].Name
		colour4 = self.G1.Skills[3].Colour

		if self.Label == "" {
			s = []string{"[label=\"", hp1, " & ", hp2, "\", color=", GraphColours[colour1], ",style=filled];\n"}
			self.Label = strings.Join(s, "")
		}
		s = []string{"    ", self.Name, " ", self.Label}
		builder.WriteString(strings.Join(s, ""))
		//let's do the same with other vertices.
		i := 0
		for i < len(self.Edges) {
			if self.Edges[i] != nil {
				other := self.Edges[i]
				if other.Name == "" {
					//formiruyu imya
					other.Name = TimeName()
				}
				//формирую ребро, записываю
				//определяю скилл этого ребра.
				var letter, colour string
				switch i {
				case 0:
					letter = "Q"
					//skill_name = Q
					colour = GraphColours[colour1]
				case 1:
					letter = "W"
					//skill_name = W
					colour = GraphColours[colour2]
				case 2:
					letter = "E"
					//skill_name = E
					colour = GraphColours[colour3]
				case 3:
					letter = "R"
					//skill_name = R
					colour = GraphColours[colour4]
				}
				var s []string
				if self.Turnnum != 0 && other != &Win1 && other != &Win2 && other != &Draw {
					s = []string{"    ", self.Name, " -> ", other.Name, "[label =\"", letter, "\",color=", colour, "];\n"}
				} else {
					s = []string{"    ", self.Name, " -> ", other.Name, ";\n"}
				}
				builder.WriteString(strings.Join(s, ""))
				//NEXT!
				add := inner(other)
				builder.WriteString(add)
			}
			i += 1
		}
	}
	return builder.String()
}

func TurnGraph(char1, char2 CharInt, turnnum int) *Graph {
	girl1 := char1.(*Girl)
	girl2 := char2.(*Girl)
	var player, opp *Girl
	var current Graph
	edges := make([]*Graph, len(girl1.Skills))
	player = girl1
	opp = girl2

	if !girl1.IsAlive() || !girl2.IsAlive() || turnnum > MAXTURN {
		Combinations += 1
		if girl1.CurrHP < girl2.CurrHP {
			if girl2.Name == Win2.Label {

				return &Win2
			} else {
				return &Win1
			}
		} else if girl2.CurrHP < girl1.CurrHP {
			if girl1.Name == Win1.Label {

				return &Win1
			} else {
				return &Win2
			}
		} else {
			return &Draw
		}
	}

	//Check availability
	skillsAvailable := player.CheckAvailableSkills(turnnum)
	//Decrease CDs
	player.DecreaseCooldowns()
	current = Graph{
		false,
		turnnum,
		"",
		TimeName(),
		player,
		opp,
		edges,
	}
	//Use Skill/s
	if player.HasEffect(SpedUp) {
	} else {
		girlsArray := make([]CharInt, len(skillsAvailable))
		opponentsArray := make([]CharInt, len(skillsAvailable))
		for i, isAvailable := range skillsAvailable {
			if isAvailable && player.CheckifAppropriate(player, opp, turnnum, i) {
				//make a copy
				copyPlayer := player.Copy()
				copyOpp := opp.Copy()
				girlsArray[i] = copyPlayer
				opponentsArray[i] = copyOpp
				//use the skill
				copyPlayer.(*Girl).Skills[i].Use(copyPlayer.(*Girl), copyOpp.(*Girl), turnnum)
				copyPlayer.(*Girl).Skills[i].CurrCD = copyPlayer.(*Girl).Skills[i].CD
				copyPlayer.(*Girl).LastUsed = i

				//decrease effects
				copyPlayer.(*Girl).DecreaseEffects(copyOpp.(*Girl), turnnum)
				if copyOpp.(*Girl).HasEffect(DelayedHeal) {
					copyOpp.(*Girl).GetEffect(DelayedHeal).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnnum)
					copyOpp.(*Girl).RemoveEffect(DelayedHeal)
				}
				if copyPlayer.(*Girl).HasEffect(CantHeal) && copyPlayer.(*Girl).GetEffect(CantHeal).Duration == 1 {
					copyPlayer.(*Girl).GetEffect(CantHeal).Remove(copyPlayer.(*Girl), copyPlayer.(*Girl), turnnum)
					copyPlayer.(*Girl).RemoveEffect(CantHeal)
				}
				if copyPlayer.(*Girl).HasEffect(TurnThreshold) && copyPlayer.(*Girl).GetEffect(TurnThreshold).Duration == 1 {
					copyPlayer.(*Girl).GetEffect(TurnThreshold).Remove(copyPlayer.(*Girl), copyPlayer.(*Girl), turnnum)
					copyPlayer.(*Girl).RemoveEffect(TurnThreshold)
				}
				if copyOpp.(*Girl).HasEffect(EuphoricHeal) {
					copyOpp.(*Girl).GetEffect(EuphoricHeal).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnnum)
					if copyOpp.(*Girl).GetEffect(EuphoricHeal).Duration == 1 {
						copyOpp.(*Girl).RemoveEffect(EuphoricHeal)
					}
				}
				if copyPlayer.(*Girl).HasEffect(EuphoricHeal) {
					copyPlayer.(*Girl).GetEffect(EuphoricHeal).Remove(copyPlayer.(*Girl), copyOpp.(*Girl), turnnum)
					if copyPlayer.(*Girl).GetEffect(EuphoricHeal).Duration == 1 {
						copyPlayer.(*Girl).RemoveEffect(EuphoricHeal)
					}
				}

				//send deeper
				current.Edges[i] = TurnGraph(copyOpp, copyPlayer, turnnum+1)

			}
		}
	}
	return &current
}
