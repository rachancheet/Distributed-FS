package main

import "fmt"

func main() {
	fs := NewFileserver("raa")
	comms := NewComms("raxx")
	comms.EstablishComms()
	go comms.listenloop(fs)
	// time.Sleep(10 * time.Second)
	for {
		fmt.Print("Enter File to get : ")
		text := ""
		fmt.Scan(&text)
		fmt.Println("Asking bhai for", text)
		comms.askBhai(Instruct{Instruction: 0, Asker_addr: comms.myAddr(), Asked_file: text})
		// fmt.Println(text)
	}
}
