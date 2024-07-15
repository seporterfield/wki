package main

import (
	"fmt"
	"os"
)

func run() {
}

const Help = `wki - Wikipedia CLI
usage: wki [-h] [-s SEARCH]

options:
  -h, --help            show this help message
  -s SEARCH, --search SEARCH
			search for topic <SEARCH>`

func printHelpAndExit() {
	fmt.Println(Help)
}

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		run()
	case 1:
		printHelpAndExit()
	default:
		run()
	}
}
