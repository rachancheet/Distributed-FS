package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("n", "rachan", "Name Identifier")
	foldername := flag.String("f", "Seed_folder", "Name Identifier")
	fs := NewFileserver(*foldername)
	shipaddr := ":8000"
	comms := NewComms(*name, shipaddr)
	comms.EstablishComms()
	go comms.listenloop(fs)
	for {
		fmt.Print("Enter File to get : ")
		text := ""
		fmt.Scan(&text)
		fmt.Println("Asking bhai for", text)
		comms.askBhai(Instruct{Instruction: 0, Asker_addr: comms.myAddr(), Asked_file: text})
	}
}
