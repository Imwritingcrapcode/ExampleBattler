package Abstract

import (
	"strconv"
	"strings"
)

const MAXTURNNUMFORWINSTARTTEST = 20

func TestStrat(char1, char2 CharInt, turnnum int, strat, response []int, num int, verd int) (*Graph, []int, int) {
	girl1 := char1.(*Girl)
	girl2 := char2.(*Girl)
	var player, opp *Girl
	var turnplayer int
	if turnnum%2 != 0 {
		turnplayer = 0
		player = girl1
		opp = girl2
	} else {
		turnplayer = 1
		player = girl2
		opp = girl1

	}
	//fmt.Println(girl1.Name, girl1.CurrHP, girl2.Name, girl2.CurrHP, turnnum, "player: ", player.Name)
	if !girl1.IsAlive() || !girl2.IsAlive() || turnnum > MAXTURNNUMFORWINSTARTTEST {
		if girl1.CurrHP < girl2.CurrHP && girl2.Number == num {
			//fmt.Println("I WON")
			//fmt.Println(response, 1)
			return nil, response, 1
		} else if girl2.CurrHP < girl1.CurrHP && girl1.Number == num {
			//fmt.Println("I WON")
			//fmt.Println(response, 1)
			return nil, response, 1
		} else if girl2.CurrHP == girl1.CurrHP {
			//fmt.Println("DRAW")
			//fmt.Println(response, 1)

			return nil, response, 1 //draw = victory
		} else {
			//fmt.Println("I LOST")
			//return []int{}, 0
			return nil, []int{}, 0
		}
	}
	//Check availability
	skillsAvailable := player.CheckAvailableSkills(turnnum)
	//Decrease CDs
	player.DecreaseCooldowns()

	finalverd := 1
	//Use Skill/s
	if player.HasEffect(SpedUp) {
		//fmt.Println("spedupbranch")
		//SPEDUPBRANCH
		if player.Number == num && !player.HasEffect(ControlledByStT) {
			//fmt.Println("player mode")
			//Player mode!
			use := strat[GetTurnNum(turnnum)-1]
			if use != 3 && skillsAvailable[use] {
				girlsArray := make([]CharInt, 5)
				opponentsArray := make([]CharInt, 5)
				verdicts := make([]int, 5)
				i := -1
				toUse := make([]int, 2)
				toUse[1] = use
				lostStrats := 0
				unplayable := 0
				for ; toUse != nil; toUse = NextSpedUpStrat(toUse, use) {
					i += 1
					//fmt.Println("trying out: ", ToStringStrat(toUse))
					if skillsAvailable[toUse[0]] &&
						skillsAvailable[toUse[1]] &&
						player.CheckifAppropriate(player, opp, turnnum, toUse[0]) &&
						player.CheckifAppropriate(player, opp, turnnum, toUse[1]) {
						//fmt.Println("OK strat")
						//before using ANY strat, copy everyone.
						//make a copy
						copyPlayer := player.Copy()
						copyOpp := opp.Copy()
						girlsArray[i] = copyPlayer
						opponentsArray[i] = copyOpp
						//use the skills
						for _, skill := range toUse {
							copyPlayer.(*Girl).Skills[skill].Use(copyPlayer.(*Girl), copyOpp.(*Girl), turnnum)
							copyPlayer.(*Girl).Skills[skill].CurrCD = copyPlayer.(*Girl).Skills[skill].CD
							copyPlayer.(*Girl).LastUsed = skill
						}
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
						//send deeper?..
						if turnplayer == 0 {
							//player, opp
							_, response, verdicts[i] = TestStrat(copyPlayer, copyOpp, turnnum+1, strat, response, num, verd)
						} else {
							//opp, player
							_, response, verdicts[i] = TestStrat(copyOpp, copyPlayer, turnnum+1, strat, response, num, verd)
						}
						if verdicts[i] == 0 {
							lostStrats += 1
							//fmt.Println(ToStringStrat(toUse), "has lost", lostStrats)
						} else if verdicts[i] == -10 {
							//fmt.Println(ToStringStrat(toUse), "is unplayable")
							unplayable += 1
							//fmt.Println("invalid branch", ToStringStrat(toUse))
						} else if verdicts[i] == 1 {
							//fmt.Println(ToStringStrat(toUse), "has won")
							//fmt.Println("winstrat")
							break
						}
					} else {
						unplayable += 1
						//fmt.Println("failed", unplayable)
						if (use == 0 && unplayable == 4) ||
							(use == 1 && unplayable == 5) ||
							(use == 2 && unplayable == 3) {
							//not a viable strat, checked everything I could and it  was all wrong.
							//fmt.Println("-10 all strats are WRONG")
							//fmt.Println("response", response, -10)
							return nil, response, -10
						}
					}

				}
				//fmt.Println("lost", lostStrats, "i", i)
				if lostStrats-1 == i || (lostStrats+unplayable-1 == i && lostStrats > 0) {
					//fmt.Println("loststrats==i player SU")
					finalverd = 0
				}
				//fmt.Println("couldn't play", unplayable, "i", i)
				if unplayable-1 == i {
					//fmt.Println("unplayable!!")
					finalverd = -10
				}
			} else {
				//Ulti? How?
				//fmt.Println("-10 ulti")
				//fmt.Println("ULTI???")
				return nil, response, -10
			}
		} else {
			//NPC Mode
			//fmt.Println("NPC SpedUP")
			STRATSMAX := 7
			girlsArray := make([]CharInt, STRATSMAX)
			opponentsArray := make([]CharInt, STRATSMAX)
			verdicts := make([]int, STRATSMAX)
			// start going through the strats
			i := -1
			toUse := make([]int, 2)
			//fixing "no turns"
			allInappr := true
			for ; toUse != nil; toUse = NextSpedUpStrat(toUse, -1) {
				if skillsAvailable[toUse[0]] &&
					skillsAvailable[toUse[1]] &&
					(player.CheckifAppropriate(player, opp, turnnum, toUse[0]) &&
						player.CheckifAppropriate(player, opp, turnnum, toUse[1])) {
					allInappr = false
					break
				}
			}
			//fmt.Println("All inappr?", allInappr)
			toUse = make([]int, 2)
			lostStrats := 0
			unplayable := 0
			for ; toUse != nil; toUse = NextSpedUpStrat(toUse, -1) {
				i += 1
				//fmt.Println("curr strat,", ToStringStrat(toUse))
				if skillsAvailable[toUse[0]] &&
					skillsAvailable[toUse[1]] &&
					(player.CheckifAppropriate(player, opp, turnnum, toUse[0]) &&
						player.CheckifAppropriate(player, opp, turnnum, toUse[1]) || allInappr || player.HasEffect(ControlledByStT)) {
					//make a copy
					//fmt.Println("OK strat")
					copyPlayer := player.Copy()
					copyOpp := opp.Copy()
					girlsArray[i] = copyPlayer
					opponentsArray[i] = copyOpp
					for _, skill := range toUse {
						//use the skill
						copyPlayer.(*Girl).Skills[skill].Use(copyPlayer.(*Girl), copyOpp.(*Girl), turnnum)
						copyPlayer.(*Girl).Skills[skill].CurrCD = copyPlayer.(*Girl).Skills[skill].CD
						copyPlayer.(*Girl).LastUsed = skill
					}
					//decrease effs
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
					if turnplayer == 0 {
						//player, opp
						_, response, verdicts[i] = TestStrat(copyPlayer, copyOpp, turnnum+1, strat, response, num, verd)
					} else {
						//opp, player
						_, response, verdicts[i] = TestStrat(copyOpp, copyPlayer, turnnum+1, strat, response, num, verd)
					}
					if verdicts[i] == 0 && (player.Number == num || !player.HasEffect(ControlledByStT)) {
						finalverd = 0
						//fmt.Println("lost on", ToStringStrat(toUse))
						if player.Number != num {
							response = append(toUse, response...)
						}
						break
					} else if verdicts[i] == 0 && player.Number != num && player.HasEffect(ControlledByStT) {
						lostStrats += 1
						if player.Number != num && lostStrats == 1 {
							response = append(toUse, response...)
						}
					} else if verdicts[i] == -10 && player.HasEffect(ControlledByStT) {
						//fmt.Println("-10")
						unplayable += 1
					} else if verdicts[i] == -10 && !player.HasEffect(ControlledByStT) {
						//fmt.Println("-10")
						finalverd = -10
						break
					} else if verdicts[i] == 1 && player.HasEffect(ControlledByStT) && player.Number != num {
						//fmt.Println("controlled and won")
						break
					}
				} else {
					//fmt.Println("can't use curr strat")
					finalverd = -1
					unplayable += 1
					if unplayable == STRATSMAX {
						//fmt.Println(response, -10)
						return nil, response, -10
					}
				}
			}
			if (lostStrats-1 == i || lostStrats+unplayable-1 == i) && finalverd != -10 {
				//fmt.Println("loststrats==i npc SU")
				finalverd = 0
			}
			if unplayable-1 == i {
				//fmt.Println("unplayable!!")
				finalverd = -10
			}
		}
	} else {
		//fmt.Println("normal branch")
		//NORMAL BRANCH
		verdicts := make([]int, len(skillsAvailable))
		if player.Number == num && !player.HasEffect(ControlledByStT) {
			use := strat[GetTurnNum(turnnum)-1]
			if skillsAvailable[use] && (player.CheckifAppropriate(player, opp, turnnum, use)) {
				//fmt.Println("can use", ToStringStrat([]int{use}))
				//copy not needed
				//skill use
				player.Skills[use].Use(player, opp, turnnum)
				player.Skills[use].CurrCD = player.Skills[use].CD
				player.LastUsed = use
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

				if opp.HasEffect(EuphoricHeal) {
					opp.GetEffect(EuphoricHeal).Remove(opp, player, turnnum)
					if opp.GetEffect(EuphoricHeal).Duration == 1 {
						player.RemoveEffect(EuphoricHeal)
					}
				}
				if player.HasEffect(EuphoricHeal) {
					player.GetEffect(EuphoricHeal).Remove(player, opp, turnnum)
					if player.GetEffect(EuphoricHeal).Duration == 1 {
						player.RemoveEffect(EuphoricHeal)
					}
				}
				if turnplayer == 0 {
					//player, opp
					_, response, finalverd = TestStrat(player, opp, turnnum+1, strat, response, num, verd)
				} else {
					//opp, player
					_, response, finalverd = TestStrat(opp, player, turnnum+1, strat, response, num, verd)
				}
			} else {
				//fmt.Println("w r o n g:", player.CheckifAppropriate(player, opp, turnnum, use), skillsAvailable[use], ToStringStrat([]int{use}))
				finalverd = -10
			}
		} else {
			//NPC BRANCH
			//fmt.Println("NPC normal")
			girlsArray := make([]CharInt, len(skillsAvailable))
			opponentsArray := make([]CharInt, len(skillsAvailable))
			allInappr := true
			for i, isAvailable := range skillsAvailable {
				if isAvailable && (player.CheckifAppropriate(player, opp, turnnum, i)) {
					allInappr = false
					break
				}
			}
			lostTechnology := 0
			allBroken := 0
			for i, isAvailable := range skillsAvailable {
				if isAvailable && (player.CheckifAppropriate(player, opp, turnnum, i) || allInappr || player.HasEffect(ControlledByStT)) {
					//fmt.Println(ToStringStrat([]int{i}), "is ok")
					//make a copy
					copyPlayer := player.Copy()
					copyOpp := opp.Copy()
					girlsArray[i] = copyPlayer
					opponentsArray[i] = copyOpp
					//use the skill
					copyPlayer.(*Girl).Skills[i].Use(copyPlayer.(*Girl), copyOpp.(*Girl), turnnum)
					copyPlayer.(*Girl).Skills[i].CurrCD = copyPlayer.(*Girl).Skills[i].CD
					copyPlayer.(*Girl).LastUsed = i
					//decrease effs
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
					if turnplayer == 0 {
						//player, opp
						_, response, verdicts[i] = TestStrat(copyPlayer, copyOpp, turnnum+1, strat, response, num, verd)
					} else {
						//opp, player
						_, response, verdicts[i] = TestStrat(copyOpp, copyPlayer, turnnum+1, strat, response, num, verd)
					}
					if verdicts[i] == 0 && (!player.HasEffect(ControlledByStT) || player.Number == num) {
						//fmt.Println(MapSkill[i], ToStringStrat(response), turnnum, "lost")
						if player.Number != num {
							response = append([]int{i}, response...)
						}
						finalverd = 0
						break
					} else if verdicts[i] == 0 {
						lostTechnology += 1
						if player.Number != num && lostTechnology == 1 {
							response = append([]int{i}, response...)
						}
					} else if verdicts[i] == -10 && player.Number != num && !player.HasEffect(ControlledByStT) {
						//fmt.Println("minus tennn")
						finalverd = -10
						break
					} else if verdicts[i] == -10 && (player.Number == num || player.HasEffect(ControlledByStT)) {
						allBroken += 1
					} else if player.HasEffect(ControlledByStT) && verdicts[i] == 1 && player.Number != num {
						//fmt.Println("owo story win")
						break
					}
				} else {
					//fmt.Println("-1 :(")
					allBroken += 1
					finalverd = -1
				}
			}
			if lostTechnology == len(skillsAvailable) || lostTechnology+allBroken == len(skillsAvailable) {
				//fmt.Println("rosuto tekki")
				finalverd = 0
			}
			if allBroken == len(skillsAvailable) {
				//fmt.Println("nO wAy OuT")
				finalverd = -10
			}
		}
	}
	//fmt.Println("response", response, "final", finalverd)
	return nil, response, finalverd
}

func NextStrat(prev []int) []int {
	if len(prev) > 10 {
		panic("strat is too long: " + strconv.Itoa(len(prev)))
	}
	curr := make([]int, len(prev))
	needAdd := 1
	for i := len(prev) - 1; i >= 0; i-- {
		curr[i] = (prev[i] + needAdd) % 4
		needAdd = (prev[i] + needAdd) / 4
		if needAdd > 0 && i == 0 {
			return nil
		}
	}
	return curr
}

func ToStringStrat(strat []int) string {
	s := ""
	for i := 0; i < len(strat); i++ {
		s += MapSkill[strat[i]]
	}
	return s
}

func NextSpedUpStrat(current []int, mustHave int) []int {
	if current[0] >= 2 && current[1] >= 1 {
		return nil
	}
	res := make([]int, 2)
	if mustHave != -1 {
		for ; res != nil; res = NextSpedUpStrat(res, -1) {
			if (res[0] == mustHave || res[1] == mustHave) &&
				!(res[0] == 2 && res[1] == 0) &&
				(res[0] == current[0] && res[1] > current[1] || res[0] > current[0]) {
				return res
			}
		}
	} else {
		needadd := 1
		for i := len(current) - 1; i >= 0; i-- {
			res[i] = (current[i] + needadd) % 3
			needadd = (current[i] + needadd) / 3
			if current[i] > 2 {
				panic("why are you making me use ulti when it's forbidden")
			}
		}
		if res[0] == 2 && res[1] == 0 {
			return NextSpedUpStrat(res, -1)
		}
	}
	return res
}

func FromStringStrat(str string) []int {
	s := make([]int, len(str))
	str = strings.ToUpper(str)
	for i := 0; i < len(str); i++ {
		s[i] = SkillMap[string(str[i])]
	}
	return s
}
