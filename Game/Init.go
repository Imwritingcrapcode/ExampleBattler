package Game

import (
	. "../Abstract"
	. "../Characters"
	"strconv"
)

func InitAsNumber(girl *Girl, number int) {
	switch number {
	case 1:
		(*Storyteller)(girl).Init()
	case 9:
		(*Euphoria)(girl).Init()
	case 10:
		(*Ruby)(girl).Init()
	case 33:
		(*Speed)(girl).Init()
	case 51:
		(*Milana)(girl).Init()
	case 119:
		(*Structure)(girl).Init()
	default:
		panic("This girl is not released yet: " + strconv.Itoa(number))
	}
}

func InitAsNumberBeta(girl *Girl, number int) {
	switch number {
	case 1:
		(*Storyteller)(girl).Init()
	case 8:
		(*Z89)(girl).Init()
	case 9:
		(*Euphoria)(girl).Init()
	case 10:
		(*Ruby)(girl).Init()
	case 33:
		(*Speed)(girl).Init()
	case 51:
		(*Milana)(girl).Init()
	case 119:
		(*Structure)(girl).Init()
	default:
		panic("Invalid girl number " + strconv.Itoa(number))
	}
}
