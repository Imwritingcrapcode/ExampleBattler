package main

import (
	. "../Abstract"
	. "../Characters"
	"fmt"
	"os"
)

func main() {
	var g1, g2 CharInt
	var Tree Graph
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	girl1 := (*Storyteller)(gi1)
	girl2 := (*Ruby)(gi2)
	girl1.Init()
	girl2.Init()
	Win1.InitWinVert(girl1.Name)
	Win2.InitWinVert(girl2.Name)
	Draw.InitWinVert("Draw!!")
	Tree = Graph{
		false,
		0,
		"[label=\"" + girl1.Name + ": Game Start!" + "\", color=" + GraphColours[Gray] + ", style=filled];\n",
		TimeName(),
		(*Girl)(girl1),
		(*Girl)(girl2),
		make([]*Graph, len(girl1.Skills)),
	}
	Tree.Edges[0] = TurnGraph(gi1, gi2, 1)
	fmt.Println(Combinations)
	graph, err := os.OpenFile("C:/Users/~C-o-L/GoglandProjects/Battler/TestFiles/graph.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	graph.WriteString(Tree.ToString())

}
