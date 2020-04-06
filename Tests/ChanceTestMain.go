package main

import (
	. "../Abstract"
	. "../Game"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	MAXTEST := 100000
	test, err := os.OpenFile("C:/Users/~C-o-L/GoglandProjects/ExampleBattler/TestFiles/ChanceTest.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer test.Close()
	testingHeader := "Testing " + strconv.Itoa(MAXTEST) + " time(s).\n"
	fmt.Println(testingHeader)
	test.WriteString(testingHeader)
	var g1, g2 CharInt
	var totalWins, totalDiff, totalLen int64
	var WinRates map[int]float64
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	WinRates = make(map[int]float64, len(ReleasedCharacters))
	for _, INIT1 := range ReleasedCharacters {
		for _, INIT2 := range ReleasedCharacters {
			if INIT1 != INIT2 {
				totalDiff = 0
				totalWins = 0
				totalLen = 0
				currLen := 20
				for i := 0; i < MAXTEST; i++ {
					InitAsNumber(gi1, INIT1)
					InitAsNumber(gi2, INIT2)
					for i := 1; i < 21; i++ {
						if i%2 == 1 {
							TurnChance(gi1, gi2, i)
						} else {
							TurnChance(gi2, gi1, i)
						}
						if !gi1.IsAlive() || !gi2.IsAlive() {
							currLen = i
							break

						}
					}
					totalLen += int64(GetTurnNum(currLen))
					totalDiff += int64(gi1.CurrHP - gi2.CurrHP)
					if gi1.CurrHP >= gi2.CurrHP {
						totalWins += 1
					}
				}
				WinRates[gi2.Number] += float64(int64(MAXTEST) - totalWins)
				WinRates[gi1.Number] += float64(totalWins)
				result := gi1.Name + " \t- " + gi2.Name + ", wins with " +
					strconv.FormatFloat(float64(totalWins)/float64(MAXTEST)*100.0, 'f', 3, 64) +
					"% chance, " + strconv.Itoa(int(math.Round(float64(totalDiff)/float64(MAXTEST)))) +
					" average hp difference, and " + strconv.Itoa(int(math.Round(float64(totalLen)/float64(MAXTEST)))) + " turns average game length.\n"
				test.WriteString(result)
				fmt.Println(result)
			}
		}
	}
	keys := make([]int, len(WinRates))
	i := 0
	for k := range WinRates {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for _, key := range keys {
		wins := WinRates[key]
		WinRates[key] = wins * 100.0 / (float64(MAXTEST * (len(ReleasedCharacters) - 1) * 2))
		strin := ReleasedCharactersNames[key] + "'s average winrate is: " + strconv.FormatFloat(WinRates[key], 'f', 3, 64) + "%.\n"
		fmt.Println(strin)
		test.WriteString(strin)
	}

}
