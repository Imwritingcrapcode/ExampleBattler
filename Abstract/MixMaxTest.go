package Abstract

func MiniMax(girl1, girl2 CharInt, turnNum, depth int, isMaximizerTurn bool, bestStrat []int) (int, int, []int) {
	var player, opp *Girl
	var bestValue, endedAt, bestMove int
	endedAt = 21
	if isMaximizerTurn {
		player = girl1.(*Girl)
		opp = girl2.(*Girl)
	} else {
		player = girl2.(*Girl)
		opp = girl1.(*Girl)
	}
	//fmt.Println("is the turn of", player.Name, "strat is", bestStrat, "depth", depth, "turnNum", turnNum, player.CurrHP, "&", opp.CurrHP)
	if turnNum > 20 || !player.IsAlive() || !opp.IsAlive() || depth <= 0 {
		//fmt.Println("the game has ended, diff", girl1.(*Girl).CurrHP-girl2.(*Girl).CurrHP, turnNum > 20, !player.IsAlive(), !opp.IsAlive(), depth <= 0)
		return girl1.(*Girl).CurrHP - girl2.(*Girl).CurrHP, turnNum, bestStrat
	}
	if isMaximizerTurn && !player.HasEffect(ControlledByStT) || !isMaximizerTurn && player.HasEffect(ControlledByStT) {
		bestValue = -1000
	} else {
		bestValue = 1000

	}
	//Check availability
	skillsAvailable := player.CheckAvailableSkills(turnNum)
	//fmt.Println("availableskills", skillsAvailable)
	//Decrease CDs
	player.DecreaseCooldowns()
	if player.HasEffect(SpedUp) {
		//fmt.Println("spedupfml")
		STRATSMAX := 7
		girlsArray := make([]CharInt, STRATSMAX)
		opponentsArray := make([]CharInt, STRATSMAX)
		movesArray := make([][]int, STRATSMAX)
		valArray := make([]int, STRATSMAX)
		lengthArray := make([]int, STRATSMAX)
		currStrat := -1
		toUse := make([]int, 2)
		for ; toUse != nil; toUse = NextSpedUpStrat(toUse, -1) {
			currStrat += 1
			if skillsAvailable[toUse[0]] &&
				skillsAvailable[toUse[1]] {
				copyPlayer := player.Copy()
				copyOpp := opp.Copy()
				girlsArray[currStrat] = copyPlayer
				opponentsArray[currStrat] = copyOpp
				for _, skill := range toUse {
					//use the skill
					copyPlayer.(*Girl).Skills[skill].Use(copyPlayer.(*Girl), copyOpp.(*Girl), turnNum)
					copyPlayer.(*Girl).Skills[skill].CurrCD = copyPlayer.(*Girl).Skills[skill].CD
					copyPlayer.(*Girl).LastUsed = skill
					//decrease effects
					copyPlayer.(*Girl).DecreaseEffects(copyOpp.(*Girl), turnNum)
					if copyOpp.(*Girl).HasEffect(DelayedHeal) {
						copyOpp.(*Girl).GetEffect(DelayedHeal).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnNum)
						copyOpp.(*Girl).RemoveEffect(DelayedHeal)
					}
					if copyPlayer.(*Girl).HasEffect(EuphoricHeal) {
						copyPlayer.(*Girl).GetEffect(EuphoricHeal).Remove(copyPlayer.(*Girl), copyOpp.(*Girl), turnNum)
						if copyPlayer.(*Girl).GetEffect(EuphoricHeal).Duration == 1 {
							copyPlayer.(*Girl).RemoveEffect(EuphoricHeal)
						}
					}
					if copyOpp.(*Girl).HasEffect(EuphoricHeal) {
						//fmt.Println("THIS IS TURN NUMBER,", turnNum, "HP", copyOpp.(*Girl).CurrHP)
						copyOpp.(*Girl).GetEffect(EuphoricHeal).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnNum)
						if copyOpp.(*Girl).GetEffect(EuphoricHeal).Duration == 1 {
							copyOpp.(*Girl).RemoveEffect(EuphoricHeal)
						}
						//fmt.Println("THIS IS TURN NUMBER,", turnNum, "HP", copyOpp.(*Girl).CurrHP)
					}
					if copyPlayer.(*Girl).HasEffect(CantHeal) && copyPlayer.(*Girl).GetEffect(CantHeal).Duration == 1 {
						copyPlayer.(*Girl).GetEffect(CantHeal).Remove(copyPlayer.(*Girl), copyOpp.(*Girl), turnNum)
						copyPlayer.(*Girl).RemoveEffect(CantHeal)
					}
					if copyOpp.(*Girl).HasEffect(TurnThreshold) && copyOpp.(*Girl).GetEffect(TurnThreshold).Duration == 1 {
						copyOpp.(*Girl).GetEffect(TurnThreshold).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnNum)
						copyOpp.(*Girl).RemoveEffect(TurnThreshold)
					}

					//send deeper
					if isMaximizerTurn {
						valArray[currStrat], lengthArray[currStrat], movesArray[currStrat] = MiniMax(copyPlayer, copyOpp, turnNum+1, depth-1, !isMaximizerTurn, append(bestStrat, toUse...))
					} else {
						valArray[currStrat], lengthArray[currStrat], movesArray[currStrat] = MiniMax(copyOpp, copyPlayer, turnNum+1, depth, !isMaximizerTurn, bestStrat)

					}
				}

			}

		} //check for all of the strats you've just played.
		//fmt.Println("checking the strats")
		for i := 0; i <= currStrat; i++ {
			if girlsArray[i] != nil {
				if (isMaximizerTurn && !player.HasEffect(ControlledByStT)) ||
					(!isMaximizerTurn && player.HasEffect(ControlledByStT)) {
					//search for max
					if lengthArray[i] < endedAt && valArray[i] > bestValue && valArray[i] > 0 ||
						valArray[i] > bestValue {
						bestValue = valArray[i]
						endedAt = lengthArray[i]
						bestMove = i
						//fmt.Println("best move is", ToStringStrat([]int{i}))
					}

				} else {
					//search for min
					if lengthArray[i] < endedAt && valArray[i] < bestValue && valArray[i] < 0 ||
						valArray[i] < bestValue {
						bestValue = valArray[i]
						endedAt = lengthArray[i]
						bestMove = i
						//fmt.Println("best move is", ToStringStrat([]int{i}))
					}
				}
			}
		}
		//if isMaximizerTurn {
		//	bestStrat = append(movesArray[bestMove], bestStrat...)
		//}
		//fmt.Println("bst diff", bestValue, "bestStrat", movesArray[bestMove])
		return bestValue, endedAt, movesArray[bestMove]
	} else {
		//fmt.Println("notspedupphew")
		//Normal branch
		STRATSMAX := len(skillsAvailable)
		girlsArray := make([]CharInt, STRATSMAX)
		opponentsArray := make([]CharInt, STRATSMAX)
		movesArray := make([][]int, STRATSMAX)
		valArray := make([]int, STRATSMAX)
		lengthArray := make([]int, STRATSMAX)
		currStrat := -1
		for i, isAvailable := range skillsAvailable {
			//fmt.Println("checking", ToStringStrat([]int{i}), "for", player.Name, "turn", turnNum, isAvailable)
			currStrat++
			if isAvailable {
				//make a copy
				copyPlayer := player.Copy()
				copyOpp := opp.Copy()
				girlsArray[i] = copyPlayer
				opponentsArray[i] = copyOpp
				//use the skill
				copyPlayer.(*Girl).Skills[i].Use(copyPlayer.(*Girl), copyOpp.(*Girl), turnNum)
				copyPlayer.(*Girl).Skills[i].CurrCD = copyPlayer.(*Girl).Skills[i].CD
				copyPlayer.(*Girl).LastUsed = i
				//decrease effs
				copyPlayer.(*Girl).DecreaseEffects(copyOpp.(*Girl), turnNum)
				if copyOpp.(*Girl).HasEffect(DelayedHeal) {
					//fmt.Println("THIS IS TURN NUMBER,", turnNum, "HP", copyOpp.(*Girl).CurrHP)
					copyOpp.(*Girl).GetEffect(DelayedHeal).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnNum)
					copyOpp.(*Girl).RemoveEffect(DelayedHeal)
					//fmt.Println("THIS IS TURN NUMBER,", turnNum, "HP", copyOpp.(*Girl).CurrHP)
				}
				if copyPlayer.(*Girl).HasEffect(EuphoricHeal) {
					copyPlayer.(*Girl).GetEffect(EuphoricHeal).Remove(copyPlayer.(*Girl), copyOpp.(*Girl), turnNum)
					if copyPlayer.(*Girl).GetEffect(EuphoricHeal).Duration == 1 {
						copyPlayer.(*Girl).RemoveEffect(EuphoricHeal)
					}
				}
				if copyOpp.(*Girl).HasEffect(EuphoricHeal) {
					//fmt.Println("THIS IS TURN NUMBER,", turnNum, "HP", copyOpp.(*Girl).CurrHP)
					copyOpp.(*Girl).GetEffect(EuphoricHeal).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnNum)
					if copyOpp.(*Girl).GetEffect(EuphoricHeal).Duration == 1 {
						copyOpp.(*Girl).RemoveEffect(EuphoricHeal)
					}
					//fmt.Println("THIS IS TURN NUMBER,", turnNum, "HP", copyOpp.(*Girl).CurrHP)
				}
				if copyPlayer.(*Girl).HasEffect(CantHeal) && copyPlayer.(*Girl).GetEffect(CantHeal).Duration == 1 {
					copyPlayer.(*Girl).GetEffect(CantHeal).Remove(copyPlayer.(*Girl), copyOpp.(*Girl), turnNum)
					copyPlayer.(*Girl).RemoveEffect(CantHeal)
				}
				if copyOpp.(*Girl).HasEffect(TurnThreshold) && copyOpp.(*Girl).GetEffect(TurnThreshold).Duration == 1 {
					copyOpp.(*Girl).GetEffect(TurnThreshold).Remove(copyOpp.(*Girl), copyPlayer.(*Girl), turnNum)
					copyOpp.(*Girl).RemoveEffect(TurnThreshold)
				}

				//send deeper
				if isMaximizerTurn {
					valArray[currStrat], lengthArray[currStrat], movesArray[currStrat] = MiniMax(copyPlayer, copyOpp, turnNum+1, depth-1, !isMaximizerTurn, append(bestStrat, []int{i, -1}...))
				} else {
					valArray[currStrat], lengthArray[currStrat], movesArray[currStrat] = MiniMax(copyOpp, copyPlayer, turnNum+1, depth, !isMaximizerTurn, bestStrat)

				}
			}
		}
		//fmt.Println("checking the strats")
		for i := 0; i <= currStrat; i++ {
			if girlsArray[i] != nil {
				if (isMaximizerTurn && !player.HasEffect(ControlledByStT)) ||
					(!isMaximizerTurn && player.HasEffect(ControlledByStT)) {
					//search for max
					if lengthArray[i] < endedAt && valArray[i] > bestValue && valArray[i] > 0 ||
						valArray[i] > bestValue {
						bestValue = valArray[i]
						endedAt = lengthArray[i]
						bestMove = i
						//fmt.Println("this is its' len:", lengthArray[i])
						//fmt.Println("best max move for", player.Name, "on the turn", turnNum, "is", movesArray[bestMove])
					}

				} else {
					//search for min
					if lengthArray[i] < endedAt && valArray[i] < bestValue && valArray[i] < 0 ||
						valArray[i] < bestValue {
						bestValue = valArray[i]
						endedAt = lengthArray[i]
						bestMove = i
						//fmt.Println("this is its' len:", lengthArray[i])
						//fmt.Println("best min move for", player.Name, "on the turn", turnNum, "is", movesArray[bestMove])
					}
				}
			}
		}
		//bestStrat = append(bestStrat, movesArray[bestMove]...)

		//fmt.Println("bst diff", bestValue, "bestStrat", movesArray[bestMove])
		return bestValue, endedAt, movesArray[bestMove]
	}
}
