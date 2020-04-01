package Characters

import (
	. "../Abstract"
	"log"
	"strconv"
)

func GetGirlRarity(number int) string {
	for key, value := range ReleasedCharactersPacks {
		if Contains(value, number) != -1 {
			return key
		}
	}
	log.Panic("This girl is not released yet: " + strconv.Itoa(number))
	return ""
}

func GetGirlInfo(number int) *GirlInfo {
	var info GirlInfo
	var girl Girl
	switch number {
	case 1:
		(*Storyteller)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Control", "Meta", "Heal", "Colours"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "The one that is the beginning and the end."
		info.MainColour = girl.Skills[0].ColourCode
	case 8:
		(*Z89)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Preventive", "Control", "Nuke"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "The spy girl. Euphoria's twin sister."
		info.MainColour = girl.Skills[3].ColourCode
	case 9:
		(*Euphoria)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Heal", "Durable"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "The happy girl. Z89's twin sister."
		info.MainColour = girl.Skills[0].ColourCode
	case 10:
		(*Ruby)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Aggressive", "Preventive", "Combo"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "A fierce fighter that dances till the end."
		info.MainColour = girl.Skills[0].ColourCode
	case 33:
		(*Speed)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Combo", "Nuke", "Aggressive"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "Asneakyassassininthecornerofyoureye."
		info.MainColour = girl.Skills[0].ColourCode
	case 51:
		(*Milana)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Durable", "Heal", "Nuke", "Preventive"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "An overconfident spoiled princess."
		info.MainColour = girl.Skills[0].ColourCode
	case 119:
		(*Structure)(&girl).Init()
		//Name string
		info.Name = girl.Name
		//Number int
		info.Number = girl.Number
		//Rarity string
		info.Rarity = GetGirlRarity(girl.Number)
		//Tags []string
		info.Tags = []string{"Defensive", "Durable", "Heal"}
		//Skills []string
		info.Skills = girl.SkillsStringList()
		//SkillColours []string
		info.SkillColours = girl.SkillColoursToString()
		//SkillColourCodes []string
		info.SkillColourCodes = girl.SkillColourCodesToString()
		//Description
		info.Description = "An android girl with no eyes."
		info.MainColour = girl.Skills[0].ColourCode
	default:
		log.Println("This girl is not released yet: " + strconv.Itoa(number))

	}
	return &info
}
