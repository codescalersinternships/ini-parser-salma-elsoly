package main

import(
	"iniparser/cmd/pkg"
	"fmt"
)

func  main(){
	p:=parser.NewParser()
	p.LoadFromString("[fruits]\nPort=550")
	fmt.Println()
}