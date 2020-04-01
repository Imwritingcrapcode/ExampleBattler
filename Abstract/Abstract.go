package Abstract

//This file is for the interfaces, probably unchangeable constants, global variables.

//Current main interface.
type CharInt interface {
	Init()
	TurnEnd(opp *Girl)
	Copy() CharInt
	Equals(other *Girl) bool
	IsAlive() bool
	SkillQ(opp *Girl, turn int)
	SkillW(opp *Girl, turn int)
	SkillE(opp *Girl, turn int)
	SkillUlti(opp *Girl, turn int)
}

//Types we use
type Colour int
type EffectType int
type EffectID int

//All of the constants for the game.
const TOTALEFFECTS = 19
const TOTALCOLOURS = 12
const MAXTAILSNUMBER = 3

const (
	None Colour = iota
	Red
	Orange
	Yellow
	Green
	Cyan
	Blue
	Violet
	Pink
	Gray
	Brown
	Black
	White
)

const (
	Basic EffectType = iota
	Prohibiting
	Control
	Debuff
	Buff
	State
	Numerical
)

const (
	//Basic
	DmgMul EffectID = iota
	DmgAdd          //Not yet in the game but...
	//Prohibiting
	CantHeal
	CantUse
	//Control
	ControlledByStT
	//Debuff
	TurnReduc
	//Buff
	AtkReduc
	TurnThreshold
	//State
	Unseen
	SpedUp
	DelayedHeal
	Invulnerable
	EuphoricHeal
	//Numerical
	GreenToken
	BlackToken
	StolenHP
	BoostShock
	BoostLayers
	EuphoricSource
)

var ColoursToString = map[Colour]string{
	Red:    "Red",
	Orange: "Orange",
	Yellow: "Yellow",
	Green:  "Green",
	Cyan:   "Cyan",
	Blue:   "Blue",
	Violet: "Violet",
	Pink:   "Pink",
	Gray:   "Gray",
	Brown:  "Brown",
	Black:  "Black",
	White:  "White",
}

//Map skillletters -> numbers
var SkillMap = map[string]int{
	"Q": 0,
	"W": 1,
	"E": 2,
	"R": 3,
}

var MapSkill = map[int]string{
	0: "Q",
	1: "W",
	2: "E",
	3: "R",
}

//Map EffectIDs -> names
var EffectNames = map[EffectID]string{
	//Basic
	DmgMul: "Damage Multiplied:", //state
	DmgAdd: "Additional Damage:", //state
	//Prohibiting
	CantHeal: "Can't Heal!",         //none
	CantUse:  "Can't use a colour:", //colour
	//Control
	ControlledByStT: "Controlled by Storyteller!", //none
	//Debuff
	TurnReduc:     "Damage reduced this turn by:",        //state
	//Buff
	TurnThreshold: "Damage min this turn:",               //state
	AtkReduc:      "Next opp's attack's damage is reduced by:", //state
	//State
	Unseen:       "Unseen for:",              //duration
	SpedUp:       "SpedUp for:",              //duration
	DelayedHeal:  "Delayed Heal for:",        //duration
	Invulnerable: "Invulnerable for:",        //duration
	EuphoricHeal: "Healing from Source for:", //duration
	//Numerical
	GreenToken:     "Green Tokens:",        //state
	BlackToken:     "Black Tokens:",        //state
	StolenHP:       "Stolen HP:",           //state
	BoostShock:     "Boosted Shock by:",    //state
	BoostLayers:    "Boosted Layers by:",   //state
	EuphoricSource: "Euphoric Source has:", //state
}

var Rarities = []string{
	0: "ST",
	1: "AD",
	2: "SP",
	3: "RP",
	4: "LF",
}

type GirlInfo struct {
	Name         string
	Number       int
	Rarity       string
	Tags         []string
	Skills       []string
	SkillColours []string
	SkillColourCodes []string
	Description  string
	MainColour   string
}

var ReleasedCharacters = []int{
	1,
	9,
	10,
	33,
	51,
	119,
}

var ReleasedCharactersNames = map[int]string{
	1:   "Storyteller",
	8:	  "Z89",
	9:	 "Euphoria",
	10:  "Ruby",
	33:  "Speed",
	51:  "Milana",
	119: "Structure",
}

var ReleasedCharactersPacks = map[string][]int{
	"ST": {},
	"AD": {8, 9, 33},
	"SP": {10, 51},
	"RP": {},
	"LF": {1, 119},
}

var BetaCharacters = append(ReleasedCharacters, []int{8, 9, 118}...)

var BetaCharactersNames = func() map[int]string {
	m := make(map[int]string)
	for key, value := range ReleasedCharactersNames {
		m[key] = value
	}
	m[8] = "z89"
	m[9] = "Euphoria"
	m[118] = "Void"
	return m
}()
