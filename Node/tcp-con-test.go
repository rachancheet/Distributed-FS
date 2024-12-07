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
type Comms struct {
	bhai           peer
	name           string
	ip             string
	motherShip     net.Conn
	Listener       net.Listener
	portAddr       uint16
	mothershipaddr string
}

type Instruct struct {
	Instruction uint16
	Asker_addr  string
	Asked_file  string
	File_data   []byte
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
func (C Comms) Callhome() net.Conn {

	motherShip, err := net.Dial("tcp", C.mothershipaddr)
	if err != nil {
		log.Fatal("Mother-ship unreachable")
	}
	return motherShip
}
func Listening() (net.Listener, uint16) {

	// var Listener net.Listener
	// var err error
	// var portAddr uint16

	// //selecting listening port address
	// portAddr = 8000
	// for portAddr < 10000 {
	// 	portAddr++
	// 	Listener, err = net.Listen("tcp", ":"+string(portAddr))
	// 	if err == nil {
	// 		fmt.Printf("Port %d available\n", portAddr)
	// 		break
	// 	}
	// }
	Listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("Error listening to tcp traffic")
	}
	return Listener, uint16(Listener.Addr().(*net.TCPAddr).Port)
}
func RecvLoop(Listener net.Listener) {

}

func (C *Comms) EstablishComms() {

	C.motherShip = C.Callhome()

	C.Listener, C.portAddr = Listening()
	fmt.Print("listener info ", C.portAddr, "\n Done \n")
	// C.Listener.Accept()

	myinfo := peer{GetOutboundIP().String(), C.portAddr, C.name}

	//sending my info
	sendpeerinfo(C.motherShip, myinfo)

	//reciving bhai ki info
	recvpeerinfo(C.motherShip, &C.bhai)
	fmt.Println("Bhai ki info", C.bhai)

}

func (C *Comms) listenloop(F Fileserver) {
	// Listener := C.Listener
	// fmt.Print("listner is listening continuously", C.Listener.Addr().(*net.TCPAddr).Port)

	for {
		//testing
		// var buff []byte
		// con, err := C.Listener.Accept()
		// if err != nil {
		// 	log.Fatal("Not accepting Bhai")
		// }
		// con.Read(buff)
		// fmt.Print("msg from Bhai ", buff)
		// buff = append(buff, 1)
		// fmt.Println("Sending to Bhai", buff)
		// if _, err = con.Write(buff); err != nil {
		// 	log.Fatal("Unable to reply to bhai")
		// }
		// time.Sleep(2 * time.Second)

		var instruct_buf Instruct
		con, err := C.Listener.Accept()
		if err != nil {
			log.Fatal("Error Accepting Instruction From MotherShip")
		}
		recvinstruct(con, &instruct_buf)
		fmt.Println("Recv: ", instruct_buf)

		if instruct_buf.Instruction == 0 {
			fmt.Println("Recv get for ", instruct_buf.Asked_file, " request from", instruct_buf.Asker_addr)

			//check file
			if F.CheckFile(instruct_buf.Asked_file) {
				fmt.Println("Got file Sending")
				C.sendFile(instruct_buf, F)
			} else {
				fmt.Println("Didn't find asking Bhai")
				C.askBhai(instruct_buf)
			}

		} else if instruct_buf.Instruction == 1 {

			fmt.Println("Recived file", instruct_buf.Asked_file, " from ", instruct_buf.Asker_addr)
			F.SaveFile(instruct_buf.Asked_file, instruct_buf.File_data)
		} else {
			fmt.Println("Unknown Instruct", instruct_buf)
		}

	}

}
func (C Comms) sendFile(I Instruct, F Fileserver) {
	fmt.Println("Sending ", I.Asked_file, "to ", I.Asker_addr)
	con, err := net.Dial("tcp", I.Asker_addr)
	if err != nil {
		fmt.Println("Not able to Send file")
		return
	}
	sendinstruct(con, Instruct{Instruction: 1, Asker_addr: fmt.Sprintf("%s:%d", C.ip, C.portAddr), Asked_file: I.Asked_file, File_data: F.GetFile(I.Asked_file)})
	fmt.Println("File sent")
	con.Close()
}

// func (C Comms) testBhai() {
// 	fmt.Println("testing Bhai ", C.bhaiAddr())

// 	con, err := net.Dial("tcp", C.bhaiAddr())
// 	if err != nil {
// 		log.Fatal("NOt able to call Bhai")

//		}
//		// var buff []byte
//		con.Write([]byte{1})
//		con.Close()
//	}
func (C Comms) askBhai(I Instruct) {
	if C.bhai.Name == "testing" {
		return
	}
	con, err := net.Dial("tcp", C.bhaiAddr())
	if err != nil {
		fmt.Println("Failed to contact neighbour Reconnecting")
		return
	}
	sendinstruct(con, I)
	con.Close()
}
func (C Comms) bhaiAddr() string {
	return fmt.Sprintf("%s:%d", C.bhai.Ip, C.bhai.Port)
}

func (C Comms) myAddr() string {
	return fmt.Sprintf("%s:%d", C.ip, C.portAddr)
}

// func (C Comms) askBhaiFile(file string) {
// 	con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", C.bhai.Ip, C.bhai.Port))
// 	if err == nil {
// 		return
// 	}
// 	sendinstruct(con)
// 	con.Close()
// }

func (C Comms) Close() {
	C.motherShip.Close()
	C.Listener.Close()

}
func NewComms(name string, shipaddr string) Comms {
	// var a Comms
	// a.name = name
	return Comms{name: name, shipaddr: shipaddr}
}
func sendpeerinfo(con net.Conn, p peer) {
	if err := json.NewEncoder(con).Encode(p); err != nil {
		log.Fatal("Error occurred while sending peer info to child")
	}
}
func recvpeerinfo(con net.Conn, p *peer) {

	err := json.NewDecoder(con).Decode(&p)
	if err != nil {
		log.Fatal("Unable to decode Mother's msg")
	}

}

func sendinstruct(con net.Conn, p Instruct) {
	if err := json.NewEncoder(con).Encode(p); err != nil {
		log.Fatal("Error occurred while sending peer info to child")
	}
}
func recvinstruct(con net.Conn, p *Instruct) {
	err := json.NewDecoder(con).Decode(&p)
	if err != nil {
		log.Fatal("Unable to decode Mother's msg")
	}
}

// func main() {

// 	var bhai peer
// 	name := "raxx"
// 	var err error

// 	motherShip := Callhome()
// 	defer motherShip.Close()

// 	Listener, portAddr := Listening()
// 	defer Listener.Close()

// 	myinfo := peer{GetOutboundIP().String(), portAddr, name}

// 	//sending my info
// 	if err = json.NewEncoder(motherShip).Encode(&myinfo); err != nil {
// 		log.Fatal("Error occurred")
// 	}

// 	//reciving bhai ki info
// 	err = json.NewDecoder(motherShip).Decode(&bhai)
// 	if err != nil {
// 		log.Fatal("Unable to decode Mother's msg")
// 	}

// 	RecvLoop(Listener)

// }
