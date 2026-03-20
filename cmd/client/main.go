package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	list := flag.Bool("list", false, "list all of your notes")
	get := flag.Int("get", -1, "get an individual note (--get 1)")
	del := flag.Int("delete", -1, "delete a note (--delete 1)")

	flag.Parse()

	if flag.NFlag() > 1 {
		fmt.Println("You can only do one action at a time")
		flag.Usage()
		os.Exit(1)
	}

	if *list {
		fmt.Println("list")
	} else if *get > -1 {
		fmt.Println("get")
	} else if *del > -1 {
		fmt.Println("del")
	}

}
