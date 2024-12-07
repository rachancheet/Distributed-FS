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
			log.Println("Error Accepting")
			con.Close()
			continue
		}
		var p1 peer
		if err = json.NewDecoder(con).Decode(&p1); err != nil {
			log.Println("Unable to decode Child's msg")
			con.Close()
			continue
		}
		fmt.Println(p1)
		peers = append(peers, p1)
		if len(peers) == 1 {
			sendpeerinfo(con, peer{Name: "testing"})
		} else {
			sendpeerinfo(con, peers[len(peers)-2])
		}
		// os.Exit(1)
		con.Close()

	}
}
func sendpeerinfo(con net.Conn, p peer) {
	if err := json.NewEncoder(con).Encode(p); err != nil {
		log.Println("Error occurred while sending peer info to child")
		return
	}
}
