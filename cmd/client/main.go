package main

import (
	"flag"
	"fmt"
	"notey/pkg/note"
	"os"
)

func main() {
	client := note.Client{}

	list := flag.Bool("list", false, "list all of your notes")
	get := flag.Int("get", -1, "get an individual note (--get 1)")
	del := flag.Int("delete", -1, "delete a note (--delete 1)")
	create := flag.Bool("new", false, "create a new note")
	title := flag.String("title", "", "title of your note (used with --new)")
	content := flag.String("content", "", "content of your note (used with --new)")

	flag.Parse()

	if flag.NFlag() > 1 && !*create {
		fmt.Println("You can only do one action at a time")
		flag.Usage()
		os.Exit(1)
	}

	if *list {
		notes, err := client.GetAllNotes()
		if err != nil {
			fmt.Printf("failed to retireve notes, %v\n", err)
			os.Exit(1)
		}
		for _, note := range notes {
			fmt.Print(note.String())
		}
	} else if *get > -1 {
		note, err := client.GetNoteById(*get)
		if err != nil {
			fmt.Println("note not found")
			os.Exit(1)
		}
		fmt.Println(note.String())
	} else if *del > -1 {
		err := client.DeleteNoteById(*del)
		if err != nil {
			fmt.Printf("failed to delete note %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("deleted note %d\n", *del)
	} else if *create {
		if len(*title) < 1 || len(*content) < 1 {
			fmt.Println("You must provide a title and content for your new note.")
			flag.Usage()
			os.Exit(1)
		}
		n, err := client.CreateNote(note.NewNote{Title: *title, Content: *content})
		if err != nil {
			fmt.Printf("failed to create note %v", err)
			os.Exit(1)
		}
		fmt.Println("Created note!")
		fmt.Println(n.String())
	}

}
