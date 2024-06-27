package main

import (
	"fmt"
	"iniparser/cmd/pkg"
)

func main() {
	p := parser.NewParser()
	p.LoadFromString("[fruits]\nPort=550")
	fmt.Println()
}
