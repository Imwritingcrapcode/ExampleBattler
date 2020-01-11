package Abstract

import (
	"math/rand"
	"strconv"
)

func TurnChance(girl1, girl2 CharInt, turnnum int) {
	player := girl1.(*Girl)
	opp := girl2.(*Girl)
	if turnnum > 20 || turnnum < 1 {
		panic("what's this turn number for God's sake " + strconv.Itoa(turnnum))
	}
	//Check availability
	skillsAvailable := player.CheckAvailableSkills(turnnum)
	//Decrease CDs
	player.DecreaseCooldowns()

	//Check if skills are worth using with a simple ai
	movesAvailable := 1
	if player.HasEffect(SpedUp) {
		movesAvailable = 2
	}
	for ; movesAvailable > 0; movesAvailable-- {
		//CHOOSE YOUR SKILL, MY GOOD FRIEND
		for used := false; !used; {
			allInappr := true
			for i := 0; i < len(skillsAvailable); i++ {
				if skillsAvailable[i] && player.CheckifAppropriate(player, opp, turnnum, i) {
					allInappr = false
					break
				}
			}
			//Use Skill/s
			use := 0
			if allInappr {
				//fmt.Println("don't have appr skills")
				for chosen := rand.Intn(len(skillsAvailable)); ; chosen = rand.Intn(len(skillsAvailable)) {
					if skillsAvailable[chosen] {
						use = chosen
						break
					}
				}
			} else {
				//fmt.Println("have appr skills")
				//fixPrint(player, opp, turnnum, skillsAvailable, use)
				for chosen := rand.Intn(len(skillsAvailable)); ; chosen = rand.Intn(len(skillsAvailable)) {
					if skillsAvailable[chosen] && player.CheckifAppropriate(player, opp, turnnum, chosen) {
						use = chosen
						break
					}
				}
			}
			//fixPrint(player, opp, turnnum, skillsAvailable, use)
			if skillsAvailable[use] &&
				(player.CheckifAppropriate(player, opp, turnnum, use) ||
					allInappr ||
					player.HasEffect(ControlledByStT)) {
				//fmt.Println("if")
				player.Skills[use].Use(player, opp, turnnum)
				player.Skills[use].CurrCD = player.Skills[use].CD
				player.LastUsed = use
				used = true
				skillsAvailable = player.CheckAvailableSkills(turnnum)
			} else {
				//fmt.Println("else")
				panic(nil)
			}
		}
	}

	//Remove Effects, decrease duration
	player.DecreaseEffects(opp, turnnum)
	if opp.HasEffect(DelayedHeal) {
		opp.GetEffect(DelayedHeal).Remove(opp, player, turnnum)
		opp.RemoveEffect(DelayedHeal)
	}
	if opp.HasEffect(CantHeal) && opp.GetEffect(CantHeal).Duration == 1 {
		opp.GetEffect(CantHeal).Remove(opp, player, turnnum)
		opp.RemoveEffect(CantHeal)
	}
	if opp.HasEffect(TurnThreshold) && opp.GetEffect(TurnThreshold).Duration == 1 {
		opp.GetEffect(TurnThreshold).Remove(opp, player, turnnum)
		opp.RemoveEffect(TurnThreshold)
	}

}
