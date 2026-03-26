package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	pb "notey/pkg/gen/note/v1"
	"os"
	"text/tabwriter"

	"google.golang.org/grpc"
)

type Client struct {
	noteClient pb.NoteServiceClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		noteClient: pb.NewNoteServiceClient(conn),
	}
}

func printNote(note *pb.Note) {
	fmt.Printf("\033[1m%s\033[0m\n", note.Title)
	fmt.Printf("\033[3m%s\033[0m\n", note.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("%s\n", note.Content)
}
func printTableItem(w io.Writer, n *pb.NoteAbbr) {
	fmt.Fprintf(w, "%d\t%s\t%s\n", n.Id, n.Title, n.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
}

func printListTable(w io.Writer, notes []*pb.NoteAbbr) {
	fmt.Fprintf(w, "\033[1mID\tTitle\tCreated At\033[0m\n")
	for _, v := range notes {
		printTableItem(w, v)
	}
}

func (c *Client) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	list := flag.Bool("list", false, "list all of your notes")
	watch := flag.Bool("watch", false, "lists and then watches for new notes")
	get := flag.Int("get", -1, "get an individual note (--get 1)")
	del := flag.Int("delete", -1, "delete a note (--delete 1)")
	create := flag.Bool("new", false, "create a new note")
	title := flag.String("title", "", "title of your note (used with --new)")
	content := flag.String("content", "", "content of your note (used with --new)")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	flag.Parse()

	if flag.NFlag() > 1 && !*create {
		fmt.Println("You can only do one action at a time")
		flag.Usage()
		os.Exit(1)
	}

	if *list {
		res, err := c.noteClient.GetNotes(ctx, &pb.GetNotesRequest{})
		if err != nil {
			fmt.Printf("failed to retireve notes, %v\n", err)
			os.Exit(1)
		}
		printListTable(w, res.Notes)
		w.Flush()
	} else if *watch {
		res, err := c.noteClient.GetNotes(ctx, &pb.GetNotesRequest{})
		if err != nil {
			fmt.Printf("failed to retireve notes, %v\n", err)
			os.Exit(1)
		}
		printListTable(w, res.Notes)
		w.Flush()
		stream, err := c.noteClient.WatchNotes(ctx, &pb.WatchNotesRequest{})
		if err != nil {
			log.Fatalf("failed to open stream watch stream %v", err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("stream error %v", err)
			}
			printTableItem(w, res.Note)
			w.Flush()

		}
	} else if *get > -1 {
		res, err := c.noteClient.GetNote(ctx, &pb.GetNoteRequest{Id: int32(*get)})
		if err != nil {
			fmt.Println("note not found")
			os.Exit(1)
		}
		printNote(res.Note)
	} else if *del > -1 {
		_, err := c.noteClient.DeleteNote(ctx, &pb.DeleteNoteRequest{Id: int32(*del)})
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
		res, err := c.noteClient.CreateNote(ctx, &pb.CreateNoteRequest{Title: *title, Content: *content})
		if err != nil {
			fmt.Printf("failed to create note %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created note!")
		printNote(res.Note)
	}
}
