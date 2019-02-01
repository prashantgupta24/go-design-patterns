package main

import (
	"fmt"
	"strings"
)

type tt struct {
	s string
	a int
}

// func init() {
// 	fmt.Println("init")
// }

func prettyPrintStructString(structString string) {
	braceReplacer := strings.NewReplacer("{", "\n\n",
		"}", "\n")
	replacedString := braceReplacer.Replace(structString)
	splitStrings := strings.Split(replacedString, ",")

	for _, str := range splitStrings {
		//trimmedString := strings.Trim(str, " ")
		fmt.Printf(str + "\n")
	}
}

func main() {
	t := &tt{}
	t.s = "test"
	t.a = 123

	s := fmt.Sprintf("%#v", t)
	prettyPrintStructString(s)
	rep := strings.NewReplacer("{", "\n\n",
		"}", "\n")
	s1 := rep.Replace(s)
	sp := strings.Split(s1, ",")
	for _, str := range sp {
		st := strings.Trim(str, " ")
		fmt.Printf(st + "\n")
	}

	// s1 := (strings.Replace(s, "{", "\n", -1)).Replace(s1, "}", "\n", -1)
	// s1 = strings.Replace(s1, "}", "\n", -1)

	//fmt.Println(s)
}
