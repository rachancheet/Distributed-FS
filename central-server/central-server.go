package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type peer struct {
	Ip   string
	Port uint16
	Name string
}

var peers []peer

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Somebody already listening")
	}
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatal("Error Accepting")
		}
		var p1 peer
		if err = json.NewDecoder(con).Decode(&p1); err != nil {
			log.Fatal("Unable to decode Child's msg")
		}
		fmt.Println(p1)
		peers = append(peers, p1)
		if len(peers) == 1 {
			send(con, peer{Name: "testing"})
		} else {
			send(con, peers[0])
		}
		// os.Exit(1)

	}
}
func send(con net.Conn, p peer) {
	if err := json.NewEncoder(con).Encode(p); err != nil {
		log.Fatal("Error occurred while sending peer info to child")
	}
}
