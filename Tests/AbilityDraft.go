package main

import (
	. "../Abstract"
	. "../Game"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	//INIT1 := 10
	//INIT2 := 51
	INIT1 := ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	INIT2 := ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	for INIT2 == INIT1 {
		INIT2 = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	}
	coin := rand.Intn(2)
	if coin == 0 {
		c := INIT1
		INIT1 = INIT2
		INIT2 = c
	}

	var g1, g2 CharInt
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	InitAsNumber(gi1, INIT1)
	InitAsNumber(gi2, INIT2)
	//empty their skills.
	gi1.Skills = make([]*Skill, 4)
	gi2.Skills = make([]*Skill, 4)
	for i := 0; i < 4; i++ {
		gi1.Skills[i] = &Skill{
			Name: "",
		}
		gi2.Skills[i] = &Skill{
			Name: "",
		}
	}

	//selected 4 characters.
	selected := make([]int, 4)
	skills := make([][]int, 4)

	all := make([]int, len(ReleasedCharacters))
	for i := 0; i < len(all); i++ {
		all[i] = ReleasedCharacters[i]
		//fmt.Println(all[i])
	}
	for i := 0; i < 4; i++ {
		aGirl := rand.Intn(len(all))
		selected[i] = all[aGirl]
		all = append(all[0:aGirl], all[aGirl+1:]...)
		//for v:=0; v<len(all); v++ {
		//	fmt.Println(all[v])
		//}
		skills[i] = make([]int, 4)
		skills[i][0] = 0
		skills[i][1] = 1
		skills[i][2] = 2
		skills[i][3] = 3
	}

	//drafting.
	for i := 0; i < 8; i++ {
		fmt.Println(gi1.Name, gi1.Skills[0].Name, gi1.Skills[1].Name, gi1.Skills[2].Name, gi1.Skills[3].Name)
		fmt.Println(gi2.Name, gi2.Skills[0].Name, gi2.Skills[1].Name, gi2.Skills[2].Name, gi2.Skills[3].Name)

		fmt.Println("available skills")
		for j := 0; j < len(selected); j++ {
			fmt.Println(ReleasedCharactersNames[selected[j]])
			for k := 0; k < len(skills[j]); k++ {
				fmt.Println(MapSkill[skills[j][k]])
			}
		}
		fmt.Println("\ninput girl")
		reader := bufio.NewReader(os.Stdin)
		girl, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		numb, err := strconv.Atoi(girl[:len(girl)-1])
		if err == nil {
			found := false
			index := -1
			for j := 0; j < len(selected); j++ {
				if selected[j] == numb {
					found = true
					index = j
				}
			}
			if !found {
				i--
				fmt.Println("Wrong input! not in selected")
			} else {
				fmt.Println("input skill")
				response, err := reader.ReadString('\n')
				if err != nil {
					panic(err)
				}
				value, present := SkillMap[strings.ToUpper(string(response[0]))]
				if present {
					//check if this skill is not taken yet
					isTaken := true
					skillIndex := -1
					for h := 0; h < len(skills[index]); h++ {
						if skills[index][h] == value {
							skillIndex = h
							isTaken = false
							break
						}
					}
					if isTaken {
						i--
						fmt.Println("It's taken already!")
					} else {
						var temp CharInt
						var turnGirl *Girl
						if i%2 == 0 {
							turnGirl = gi1
						} else {
							turnGirl = gi2
						}
						temp = new(Girl)
						temp1 := temp.(*Girl)
						InitAsNumber(temp1, numb)
						skillToSet := temp1.Skills[value]
						if skillToSet.IsUlti && turnGirl.Skills[3].Name != "" {
							i--
							fmt.Println("You already have an ultimate!!")
						} else if skillToSet.IsUlti {
							turnGirl.Skills[3] = temp1.Skills[value].Copy()
							skills[index] = append(skills[index][0:skillIndex], skills[index][skillIndex+1:]...)
							fmt.Println("Ulti set!")
						} else {
							for m := 0; m < 4; m++ {
								if turnGirl.Skills[m].Name == "" && m != 3 {
									turnGirl.Skills[m] = temp1.Skills[value].Copy()
									skills[index] = append(skills[index][0:skillIndex], skills[index][skillIndex+1:]...)
									fmt.Println(temp1.Skills[value].Name, "set as", MapSkill[m])
									break
								} else if turnGirl.Skills[m].Name == "" && m == 3 {
									i--
									fmt.Println("You now need an ultimate!!")
									break
								}
							}
						}

					}
				} else {
					i--
					fmt.Println("Wrong input on skill!")
				}
			}
		} else {
			fmt.Println(err)
			i--
			fmt.Println("Wrong input! not a number")
		}
	}

	fmt.Println("Time to play")
	//now play.
	/*for i := 1; i < 21; i++ {
		TurnKeyboard(gi1, gi2, i)
		other := gi1
		gi1 = gi2
		gi2 = other
		if !gi1.IsAlive() || !gi2.IsAlive() {
			break

		}
	}
	fmt.Println(g1.(*Girl).Name, g1.(*Girl).CurrHP, g2.(*Girl).Name, g2.(*Girl).CurrHP)*/
	TWOBOTS := false
	DEPTH := 5
	PLAYER := gi1.Number
	for i := 1; i < 21; i++ {
		if (gi1.Number != PLAYER && !gi1.HasEffect(ControlledByStT) ||
			gi1.Number == PLAYER && gi1.HasEffect(ControlledByStT)) || TWOBOTS {
			testfood1 := gi1.Copy()
			testfood2 := gi2.Copy()
			prediction, _, moves := MiniMax(testfood1, testfood2, i, DEPTH, true, []int{})
			use := moves[0:2]
			fmt.Println(GetGameState(gi1, gi2, i, gi1.CheckAvailableSkills(i)))
			if gi1.HasEffect(SpedUp) {
				fmt.Println("The bot has used", ToStringStrat([]int{moves[0], moves[1]})+", and predicted", prediction, "\n")
			} else {
				fmt.Println("the bot has used", ToStringStrat([]int{moves[0]})+", and predicted", prediction, "\n")
				//fmt.Println("the bot has used", ToStringStrat(moves) + ", and predicted", prediction, "\n")
			}
			TurnApply(gi1, gi2, i, use)

		} else {
			TurnKeyboard(gi1, gi2, i)
		}

		other := gi1
		gi1 = gi2
		gi2 = other
		if !gi1.IsAlive() || !gi2.IsAlive() {
			/*close(input)
			close(output)*/
			break

		}
	}
	fmt.Println(g1.(*Girl).Name, g1.(*Girl).CurrHP, g2.(*Girl).Name, g2.(*Girl).CurrHP)

}
