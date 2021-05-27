package main

import (
	"fmt"
	"strings"
)

var NoConsole bool = true

func ReadLine() string {
	var s string
	if NoConsole {
		s = ""
	} else {
		fmt.Scanln(&s)
	}
	return s
}

func Prompt(msg string) string {
	print(msg)
	return ReadLine()
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintError(err error) {
	if err != nil {
		println(err.Error())
	}
}

func PromptOptions(msg string, options map[string]string) string {
	println(msg)
	for o, m := range options {
		println("\t[" + o + "] = " + m)
	}

	var r string = ""
	if NoConsole {
		// Select First key
		for o := range options {
			r = o
			break
		}
	} else {
		r = strings.TrimSpace(strings.ToLower(Prompt("Enter [value] : ")))
	}

	_, ok := options[r]
	if ok {
		return r
	} else {
		println("Sorry, I didn't get that. Please enter the [option] you want ")
		return PromptOptions(msg, options)
	}

}
