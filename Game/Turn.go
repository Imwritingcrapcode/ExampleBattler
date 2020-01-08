package Game

import (
	. "../Abstract"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func Turn2Channels(girl1, girl2 CharInt, turnnum int, p1, p2 *ClientChannels) bool {
	player := girl1.(*Girl)
	opp := girl2.(*Girl)
	if turnnum > 20 || turnnum < 1 {
		log.Fatal("[INGAME] what's this turn number for God's sake " + strconv.Itoa(turnnum))
	}
	//Check availability
	var Clock *Clock
	availableSkills := player.CheckAvailableSkills(turnnum)
	playerState := GetGameStateChannels(player, opp, turnnum)
	opponentState := GetGameStateChannelsOpp(player, opp, turnnum)
	playerInputChannel := p1.Input
	playerGivenUp := p1.HasGivenUp
	oppGivenUp := p2.HasGivenUp
	Clock = p1.Clock
	if player.HasEffect(ControlledByStT) {
		playerInputChannel = p2.Input
		Clock = p2.Clock
		playerGivenUp = p2.HasGivenUp
		oppGivenUp = p1.HasGivenUp
	}

	//Decrease CDs
	player.DecreaseCooldowns()

	//Use Skill/s
	movesAvailable := 1
	if player.HasEffect(SpedUp) {
		movesAvailable = 2
	}
	//sending it
	p1.Send(playerState)
	p2.Send(opponentState)
	awaitingInput := true

	//start ticking
	go Clock.StartTicking()
	for ; movesAvailable > 0; movesAvailable-- {
		//CHOOSE YOUR SKILL, MY GOOD FRIEND
		for used := false; !used; {
			select {
			case response, stillOpen := <-playerInputChannel:
				if !stillOpen {
					return !stillOpen
				}
				var num int
				var present bool
				if awaitingInput {
					decision := strings.ToUpper(response)
					num, present = SkillMap[decision]

				} else {
					num = rand.Intn(len(player.Skills))
					for !availableSkills[num] {
						num = rand.Intn(len(player.Skills))
					}
					present = true
				}
				if present && availableSkills[num] {
					player.Skills[num].Use(player, opp, turnnum)
					player.Skills[num].CurrCD = player.Skills[num].CD
					player.LastUsed = num
					used = true
					availableSkills = player.CheckAvailableSkills(turnnum)

					//sending info part! gathering info:
					playerState = GetGameStateChannels(player, opp, turnnum)
					opponentState = GetGameStateChannelsOpp(player, opp, turnnum)
					//setting the animation instruction
					playerState.Instruction = "Animation:" + MapSkill[num]
					opponentState.Instruction = "Animation:" + MapSkill[num]
					//sending it

					p1.Send(playerState)
					p2.Send(opponentState)

				} else if present {
					var state1 GameState
					if player.HasEffect(ControlledByStT) {
						state1 = *p2.LastThing.Copy()
						state1.Instruction += "Input:Can't use that skill. Choose another one!"
						p2.Send(&state1)
					} else {
						state1 = *p1.LastThing.Copy()
						state1.Instruction += "Input:Can't use that skill. Choose another one!"
						p1.Send(&state1)
					}

				} else {
					var state1 GameState
					if player.HasEffect(ControlledByStT) {
						state1 = *p2.LastThing.Copy()
						state1.Instruction += "Input:Error - invalid input."
						p2.Send(&state1)
					} else {
						state1 = *p1.LastThing.Copy()
						state1.Instruction += "Input:Error - invalid input."
						p1.Send(&state1)
					}
				}
			case msg := <-playerGivenUp:
				if msg {
					return msg
				}
				//skipped a turn due to timeout
				awaitingInput = false
				myNewChannel := make(chan string, 2)
				playerInputChannel = myNewChannel
				for i := 0; i < movesAvailable; i++ {
					playerInputChannel <- "Q"
				}

			case <-oppGivenUp:
				return true
			}
		}
	}

	//stop ticking
	if Clock.State() {
		Clock.Stop()
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
	if opp.HasEffect(EuphoricHeal) {
		opp.GetEffect(EuphoricHeal).Remove(opp, player, turnnum)
	} else if player.HasEffect(EuphoricHeal) {
		player.GetEffect(EuphoricHeal).Remove(player, opp, turnnum)
	}
	return false
}

func TurnChannels(girl1, girl2 CharInt, turnnum int, p1 *ClientChannels, botInput chan string,
	botClock *Clock, isBotsTurn bool) bool {
	player := girl1.(*Girl)
	opp := girl2.(*Girl)
	if turnnum > 20 || turnnum < 1 {
		log.Fatal("[INGAME] what's this turn number for God's sake " + strconv.Itoa(turnnum))
	}
	//Check availability
	availableSkills := player.CheckAvailableSkills(turnnum)
	var Clock *Clock
	if !isBotsTurn {
		Clock = p1.Clock
		playerState := GetGameStateChannels(player, opp, turnnum)
		playerInputChannel := p1.Input

		if player.HasEffect(ControlledByStT) {
			playerInputChannel = botInput
		}
		//Decrease CDs
		player.DecreaseCooldowns()
		//Use Skill/s
		movesAvailable := 1
		if player.HasEffect(SpedUp) {
			movesAvailable = 2
		}
		//sending it
		p1.Send(playerState)
		awaitingInput := true

		//start ticking
		go Clock.StartTicking()
		for ; movesAvailable > 0; movesAvailable-- {
			//CHOOSE YOUR SKILL, MY GOOD FRIEND
			for used := false; !used; {
				select {
				case response, stillOpen := <-playerInputChannel:
					if !stillOpen {
						return !stillOpen
					}
					var num int
					var present bool
					if awaitingInput {
						decision := strings.ToUpper(response)
						num, present = SkillMap[decision]

					} else {
						num = rand.Intn(len(player.Skills))
						for !availableSkills[num] {
							num = rand.Intn(len(player.Skills))
						}
						present = true
					}
					if present && availableSkills[num] {
						player.Skills[num].Use(player, opp, turnnum)
						player.Skills[num].CurrCD = player.Skills[num].CD
						player.LastUsed = num
						used = true
						availableSkills = player.CheckAvailableSkills(turnnum)

						//sending info part! gathering info:
						playerState = GetGameStateChannels(player, opp, turnnum)
						//setting the animation instruction
						playerState.Instruction = "Animation:" + MapSkill[num]
						//sending it
						p1.Send(playerState)

					} else if present {
						if !player.HasEffect(ControlledByStT) {
							state1 := p1.LastThing.Copy()
							state1.Instruction = "Input:Can't use that skill. Choose another one!"
							p1.Send(state1)
							p1.Clock.TellTheTime()
						}

					} else {
						if !player.HasEffect(ControlledByStT) {
							state1 := p1.LastThing.Copy()
							state1.Instruction = "Input:Error - invalid input."
							p1.Send(state1)
							p1.Clock.TellTheTime()
						}
					}
				case msg := <-p1.HasGivenUp:
					if msg {
						return msg
					}
					//skipped a turn due to timeout
					awaitingInput = false
					myNewChannel := make(chan string, 2)
					playerInputChannel = myNewChannel
					for i := 0; i < movesAvailable; i++ {
						playerInputChannel <- "Q"
					}
				}
			}
		}
	} else {
		Clock = botClock
		oppState := GetGameStateChannelsOpp(player, opp, turnnum)
		playerInputChannel := botInput
		if player.HasEffect(ControlledByStT) {
			playerInputChannel = p1.Input

		}
		//Decrease CDs
		player.DecreaseCooldowns()
		//Use Skill/s
		movesAvailable := 1
		if player.HasEffect(SpedUp) {
			movesAvailable = 2
		}
		//sending it
		p1.Send(oppState)
		awaitingInput := true

		//start ticking
		go Clock.StartTickingBot()
		//log.Println("starttickingbot", turnnum, isBotsTurn)
		//somehow when battling structure/speed player
		// against story bot when story ulted this got called
		for ; movesAvailable > 0; movesAvailable-- {

			//CHOOSE YOUR SKILL, MY GOOD HUMAN FRIEND
			for used := false; !used; {

				select {
				case response, stillOpen := <-playerInputChannel:
					if !stillOpen {
						return !stillOpen
					}
					var num int
					var present bool
					if awaitingInput {
						decision := strings.ToUpper(response)
						num, present = SkillMap[decision]

					} else {
						num = rand.Intn(len(player.Skills))
						for !availableSkills[num] {
							num = rand.Intn(len(player.Skills))
						}
						present = true
					}
					if present && availableSkills[num] {
						player.Skills[num].Use(player, opp, turnnum)
						player.Skills[num].CurrCD = player.Skills[num].CD
						player.LastUsed = num
						used = true

						//sending info part! gathering info:
						oppState = GetGameStateChannelsOpp(player, opp, turnnum)
						//setting the animation instruction
						oppState.Instruction = "Animation:" + MapSkill[num]
						//sending it
						p1.Send(oppState)

					} else if present {
						if player.HasEffect(ControlledByStT) {
							state1 := p1.LastThing.Copy()
							state1.Instruction = "Input:Can't use that skill. Choose another one!"
							p1.Send(state1)
							p1.Clock.TellTheTime()
						}

					} else {
						if player.HasEffect(ControlledByStT) {
							state1 := p1.LastThing.Copy()
							state1.Instruction = "Input:Error - invalid input."
							p1.Send(state1)
							p1.Clock.TellTheTime()
						}
					}
				case msg := <-p1.HasGivenUp:
					if msg {
						return msg
					}
					//BOT skipped a turn due to timeout
					awaitingInput = false
					myNewChannel := make(chan string, 2)
					playerInputChannel = myNewChannel
					for i := 0; i < movesAvailable; i++ {
						playerInputChannel <- "Q"
					}
				}
			}
		}
	}

	if Clock.State() {
		Clock.Stop()
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
	if opp.HasEffect(EuphoricHeal) {
		opp.GetEffect(EuphoricHeal).Remove(opp, player, turnnum)
	} else if player.HasEffect(EuphoricHeal) {
		player.GetEffect(EuphoricHeal).Remove(player, opp, turnnum)
	}
	return false
}

func TurnApply(girl1, girl2 CharInt, turnnum int, strat []int) {
	player := girl1.(*Girl)
	opp := girl2.(*Girl)
	if turnnum > 20 || turnnum < 1 {
		panic("what's this turn number for God's sake " + strconv.Itoa(turnnum))
	}
	//Check availability
	skillsAvailable := player.CheckAvailableSkills(turnnum)
	//Decrease CDs
	player.DecreaseCooldowns()

	movesAvailable := 1
	if player.HasEffect(SpedUp) {
		movesAvailable = 2
	}
	for ; movesAvailable > 0; movesAvailable-- {
		//GetGameState, send to player
		//CHOOSE YOUR SKILL, MY GOOD FRIEND
		for used := false; !used; {
			var num int
			if player.HasEffect(SpedUp) && movesAvailable == 1 {
				num = strat[1]
			} else {
				num = strat[0]
			}
			if skillsAvailable[num] {
				player.Skills[num].Use(player, opp, turnnum)
				player.Skills[num].CurrCD = player.Skills[num].CD
				player.LastUsed = num
				used = true
				skillsAvailable = player.CheckAvailableSkills(turnnum)
			} else {
				//fmt.Println("You can't use that skill. Choose another one!\n")
				return
			}
		}
	}
	//}
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
	if opp.HasEffect(EuphoricHeal) {
		opp.GetEffect(EuphoricHeal).Remove(opp, player, turnnum)
	} else if player.HasEffect(EuphoricHeal) {
		player.GetEffect(EuphoricHeal).Remove(player, opp, turnnum)
	}

}

func TurnKeyboard(girl1, girl2 CharInt, turnnum int) {
	player := girl1.(*Girl)
	opp := girl2.(*Girl)
	if turnnum > 20 || turnnum < 1 {
		panic("what's this turn number for God's sake " + strconv.Itoa(turnnum))
	}
	//Check availability
	skillsAvailable := player.CheckAvailableSkills(turnnum)
	//Decrease CDs
	player.DecreaseCooldowns()

	//Use Skills
	movesAvailable := 1
	if player.HasEffect(SpedUp) {
		movesAvailable = 2
	}
	for ; movesAvailable > 0; movesAvailable-- {
		//GetGameState, send to player
		output := GetGameState(player, opp, turnnum, skillsAvailable)
		//CHOOSE YOUR SKILL, MY GOOD FRIEND
		for used := false; !used; {
			input := GetMoveFromKeyboard(output)
			use := input
			num := SkillMap[use]
			_, present := SkillMap[use]
			if present && skillsAvailable[num] {
				/*if player.CheckifAppropriate(player, opp, turnnum, num) {
					fmt.Println("good choice!\n")
				} else {
					fmt.Println("bad move, mate!\n")
				}*/
				player.Skills[num].Use(player, opp, turnnum)
				player.Skills[num].CurrCD = player.Skills[num].CD
				player.LastUsed = num
				used = true
				skillsAvailable = player.CheckAvailableSkills(turnnum)
			} else if present {
				output = "You can't use that skill. Choose another one!\n"
			} else {
				output = "Invalid input. Please choose from Q, W, E or R.\n"
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
	if opp.HasEffect(EuphoricHeal) {
		opp.GetEffect(EuphoricHeal).Remove(opp, player, turnnum)
	} else if player.HasEffect(EuphoricHeal) {
		player.GetEffect(EuphoricHeal).Remove(player, opp, turnnum)
	}
}
