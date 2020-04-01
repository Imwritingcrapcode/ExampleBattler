package Abstract

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//This file is for structs and functions we use to support other game objects.

//The one we store the effects in. More O(1)s!
type EffectSet struct {
	dense []*Effect
	max   int
	len   int
}

func (self *EffectSet) Init(highest int) {
	if highest > 0 {
		self.dense = make([]*Effect, highest)
		self.len = 0
		self.max = highest
	} else {
		panic("wrong EffectSet initialization: " + strconv.Itoa(highest))
	}
}

func (self *EffectSet) Set(eff *Effect) {
	if self.len <= self.max {
		if !self.Contains(eff.ID) {
			self.len++
		}
		self.dense[int(eff.ID)] = eff
	} else {
		panic("Set overflow, can't add")
	}
}

func (self *EffectSet) Get(ID EffectID) *Effect {
	if self.Contains(ID) {
		return self.dense[int(ID)]
	} else {
		panic("No such effect in EffectSet: " + strconv.Itoa(int(ID))) //thank you, my precious panic, you've saved my code!
	}
}

func (self *EffectSet) Remove(ID EffectID) {
	if self.Contains(ID) {
		self.len--
		self.dense[int(ID)] = nil
	} else {
		panic("No such effect in EffectSet: " + strconv.Itoa(int(ID)))
	}
}

func (self *EffectSet) Contains(ID EffectID) bool {
	return int(ID) < self.max && self.dense[int(ID)] != nil && self.dense[int(ID)].ID == ID
}

func (self *EffectSet) Size() int {
	return self.len
}

//It's my ... turn if it's the ith turn of the game.
func GetTurnNum(i int) int {
	return (i-1)/2 + 1
}

func Damage(player, opp *Girl, deal int, ignoreDef bool, colour Colour) int {
	if opp.HasEffect(Invulnerable) {
		return 0
	}
	var dmg, atkred, turnthr, turnred, mul, add int
	dmg = deal
	if player.HasEffect(DmgMul) {
		mul = player.GetEffect(DmgMul).State
	} else {
		mul = 1
	}

	if player.HasEffect(DmgAdd) {
		add = player.GetEffect(DmgAdd).State
	} else {
		add = 0
	}

	if opp.HasEffect(AtkReduc) {
		atkred = opp.GetEffect(AtkReduc).State
		opp.RemoveEffect(AtkReduc)
	} else {
		atkred = 0
	}
	if player.HasEffect(TurnReduc) {
		turnred = player.GetEffect(TurnReduc).State
	} else {
		turnred = 0
	}
	if opp.HasEffect(TurnThreshold) {
		turnthr = opp.GetEffect(TurnThreshold).State
	} else {
		turnthr = 0
	}
	/*fmt.Println("DMG", dmg)
	fmt.Println("DMGMUL", mul)
	fmt.Println("DMGADD", add)
	fmt.Println("TURN RED", turnred)
	fmt.Println("ATK RED", atkred)
	fmt.Println("DEF", opp.Defences[colour])*/
	dmg = dmg*mul + add - atkred - turnred
	if !ignoreDef {
		dmg -= opp.Defenses[colour]
	}
	if dmg > turnthr {
		opp.CurrHP -= dmg
	} else {
		dmg = 0
	}
	opp.LastDmgTaken = dmg
	return dmg
}

func Heal(player *Girl, heal int) int {
	if !player.HasEffect(CantHeal) && heal > 0 {
		if player.CurrHP+heal < player.MaxHP {
			player.CurrHP += heal
		} else {
			heal = player.MaxHP - player.CurrHP
			player.CurrHP = player.MaxHP
		}
	} else {
		heal = 0
	}
	return heal
}

func GetMoveFromKeyboard(State string) string {
	output := State
	fmt.Println(output)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	decision := strings.ToUpper(string(response[0]))
	return decision
}

func GetSkillNames(girl *Girl) map[string]string {
	names := make(map[string]string, len(girl.Skills))
	for index, skill := range girl.Skills {
		names[MapSkill[index]] = skill.Name
	}
	return names
}

func GetOppSkillNames(girl *Girl) map[string]string {
	names := make(map[string]string, len(girl.Skills))
	for index, skill := range girl.Skills {
		names["Opp"+MapSkill[index]] = skill.Name
	}
	return names
}

func GetSkillColours(girl *Girl) map[string]string {
	codes := make(map[string]string, len(girl.Skills))
	for index, skill := range girl.Skills {
		codes[MapSkill[index]] = skill.ColourCode
	}
	return codes
}
func GetOppSkillColours(girl *Girl) map[string]string {
	codes := make(map[string]string, len(girl.Skills))
	for index, skill := range girl.Skills {
		codes["Opp"+MapSkill[index]] = skill.ColourCode
	}
	return codes
}

func GetGameState(player, opp *Girl, Turn int, skillsAvailable []bool) string {
	var res strings.Builder
	res.WriteString("Turn: " + strconv.Itoa(Turn) + " (it's your " + strconv.Itoa(GetTurnNum(Turn)) + ")\n")
	res.WriteString("Player: " + player.Name + ";\tOpponent: " + opp.Name + "\n")
	res.WriteString("Your HP: " + strconv.Itoa(player.CurrHP) + "/" + strconv.Itoa(player.MaxHP) + ";\tOpponent's HP: " + strconv.Itoa(opp.CurrHP) + "/" + strconv.Itoa(opp.MaxHP) + "\n")
	res.WriteString("Your Effects:\n")
	ToStringEffects(player, &res)
	res.WriteString("Opponent's Effects:\n")
	ToStringEffects(opp, &res)
	for i, isAvailable := range skillsAvailable {
		res.WriteString("[" + strings.ToUpper(MapSkill[i]) + "] " + player.Skills[i].Name + ", ")
		if isAvailable {
			res.WriteString("is ready!\n")
		} else {
			res.WriteString("is NOT ready!\n")
		}
	}
	if !player.HasEffect(ControlledByStT) {
		res.WriteString("Choose your skill!\n")
	} else {
		res.WriteString("This turn, Storyteller chooses which skill you use!\n")
	}
	return res.String()
}

func GetGameStateChannels(player, opp *Girl, Turn int) *GameState {
	var res GameState
	res.Instruction = ""
	res.EndState = GameContinue
	res.TurnNum = Turn
	if player.HasEffect(ControlledByStT) {
		res.TurnPlayer = 2
		res.SkillState = *DescribeSkillStatesWithSecrets(player, Turn, false)
		res.OppSkillState = *DescribeSkillStatesNoSecrets(player, Turn, true)
		res.SkillNames = GetSkillNames(player)
		res.SkillColours = GetSkillColours(player)
		res.OppSkillNames = GetOppSkillNames(player)
		res.OppSkillColours = GetOppSkillColours(player)
	} else {
		res.TurnPlayer = 1
		res.SkillState = *DescribeSkillStatesWithSecrets(player, Turn, false)
		res.OppSkillState = *DescribeSkillStatesNoSecrets(opp, Turn, true)
		res.SkillNames = GetSkillNames(player)
		res.SkillColours = GetSkillColours(player)
		res.OppSkillNames = GetOppSkillNames(opp)
		res.OppSkillColours = GetOppSkillColours(opp)
	}
	res.PlayerNum = player.Number
	res.OppNum = opp.Number
	res.PlayerName = player.Name
	res.OppName = opp.Name
	res.HP = player.CurrHP
	res.MaxHP = player.MaxHP
	res.OppHP = opp.CurrHP
	res.OppMaxHP = opp.MaxHP
	res.Effects = *DescribeEffects(player)
	res.OppEffects = *DescribeEffects(opp)
	res.Defenses = player.Defenses
	res.OppDefenses = opp.Defenses
	return &res
}

func GetGameStateChannelsOpp(player, opp *Girl, Turn int) *GameState {
	var res GameState
	res.Instruction = ""
	res.EndState = GameContinue
	res.TurnNum = Turn
	if player.HasEffect(ControlledByStT) {
		res.TurnPlayer = 1
		res.SkillState = *DescribeSkillStatesNoSecrets(player, Turn, false)
		res.OppSkillState = *DescribeSkillStatesWithSecrets(player, Turn, true)
		res.SkillNames = GetSkillNames(player)
		res.SkillColours = GetSkillColours(player)
		res.OppSkillNames = GetOppSkillNames(player)
		res.OppSkillColours = GetOppSkillColours(player)
	} else {
		res.TurnPlayer = 2
		res.SkillState = *DescribeSkillStatesWithSecrets(opp, Turn, false)
		res.OppSkillState = *DescribeSkillStatesNoSecrets(player, Turn, true)
		res.SkillNames = GetSkillNames(opp)
		res.SkillColours = GetSkillColours(opp)
		res.OppSkillNames = GetOppSkillNames(player)
		res.OppSkillColours = GetOppSkillColours(player)
	}
	res.PlayerNum = opp.Number
	res.OppNum = player.Number
	res.PlayerName = opp.Name
	res.OppName = player.Name
	res.HP = opp.CurrHP
	res.MaxHP = opp.MaxHP
	res.OppHP = player.CurrHP
	res.OppMaxHP = player.MaxHP
	res.Effects = *DescribeEffects(opp)
	res.OppEffects = *DescribeEffects(player)
	res.Defenses = opp.Defenses
	res.OppDefenses = player.Defenses
	return &res
}

func DescribeSkillStatesWithSecrets(player *Girl, Turn int, add bool) *map[string]int {
	dest := make(map[string]int)
	var line string
	for index, skill := range player.Skills {
		if add {
			line = "Opp" + MapSkill[index]
		} else {
			line = MapSkill[index]
		}
		if skill.StrT > Turn {
			dest[line] = NotUnlockedYet
		} else if player.HasEffect(CantUse) && Colour(player.GetEffect(CantUse).State) == skill.Colour ||
			player.HasEffect(SpedUp) && skill.IsUlti || player.HasEffect(ControlledByStT) {
			dest[line] = DisabledByEffect
		} else if skill.CurrCD > 0 && GetTurnNum(Turn)+skill.CurrCD > 10 {
			dest[line] = Disabled
		} else if skill.CurrCD > 0 {
			dest[line] = skill.CurrCD
		} else {
			dest[line] = 0
		}

	}
	return &dest
}

func DescribeSkillStatesNoSecrets(player *Girl, Turn int, add bool) *map[string]int { //doesn't send you secret data.
	dest := make(map[string]int)
	var line string
	for index, skill := range player.Skills {
		if add {
			line = "Opp" + MapSkill[index]
		} else {
			line = MapSkill[index]
		}
		if skill.StrT > Turn {
			dest[line] = NotUnlockedYet
		} else if player.HasEffect(CantUse) && Colour(player.GetEffect(CantUse).State) == skill.Colour ||
			player.HasEffect(SpedUp) && skill.IsUlti {
			dest[line] = DisabledByEffect
		} else if skill.CurrCD > 0 && GetTurnNum(Turn)+skill.CurrCD > 10 {
			dest[line] = Disabled
		} else if skill.CurrCD > 0 {
			dest[line] = skill.CurrCD
		} else {
			dest[line] = 0
		}

	}
	return &dest
}

func DescribeEffects(player *Girl) *map[int]string {
	desc := make(map[int]string)
	for i := 0; i < TOTALEFFECTS; i++ {
		if player.HasEffect(EffectID(i)) {
			if EffectID(i) == ControlledByStT || EffectID(i) == CantHeal {
				desc[int(EffectID(i))] = ""
			} else if EffectID(i) == CantUse {
				desc[int(EffectID(i))] = ColoursToString[Colour(player.GetEffect(EffectID(i)).State)]
			} else if player.GetEffect(EffectID(i)).Type == State || EffectID(i) == DmgMul {
				desc[int(EffectID(i))] = strconv.Itoa(player.GetEffect(EffectID(i)).Duration)
			} else {
				desc[int(EffectID(i))] = strconv.Itoa(player.GetEffect(EffectID(i)).State)
			}
		}
	}
	return &desc
}

func ToStringEffects(player *Girl, res *strings.Builder) {
	if player.Effects.Size() != 0 {
		j := 0
		for i := 0; i < TOTALEFFECTS; i++ {
			if player.HasEffect(EffectID(i)) {
				j += 1
				if EffectID(i) == ControlledByStT || EffectID(i) == CantHeal {
					res.WriteString(EffectNames[EffectID(i)] + ";")
				} else if EffectID(i) == CantUse {
					res.WriteString(EffectNames[EffectID(i)] + " " + ColoursToString[Colour(player.GetEffect(EffectID(i)).State)] + ";")
				} else if player.GetEffect(EffectID(i)).Type == State || EffectID(i) == DmgMul {
					res.WriteString(EffectNames[EffectID(i)] + " " + strconv.Itoa(player.GetEffect(EffectID(i)).Duration) + ";")
				} else {
					res.WriteString(EffectNames[EffectID(i)] + " " + strconv.Itoa(player.GetEffect(EffectID(i)).State) + ";")
				}
				if j != player.Effects.Size() && j%2 == 0 {
					res.WriteString("\n")
				} else if j != player.Effects.Size() {
					res.WriteString("\t")
				}
			}
		}
	} else {
		res.WriteString("[None so far]")
	}
	res.WriteString("\n")
}

// An Item is something we manage in a Priority queue.
type Item struct {
	UserID   int64 // The value of the item; arbitrary.
	Priority int64 // The Priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The Index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Remove(i int64) bool {
	for k, item := range pq {
		if item.UserID == i {
			pq = append(pq[:k], pq[k+1:]...)
			return true
		}
	}
	return false
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the Priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value int64, priority int64) {
	item.UserID = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

func Contains(a interface{}, e interface{}) int {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return i
		}
	}
	return -1
}
