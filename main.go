package main

import (
	"fmt"
	"net"
	"os"
	"github.com/google/uuid"
	"encoding/json"
	"strings"
	"bytes"
)

type Game struct {
	P1 Player 
	P2 Player 
	Id string 
} 

type Player struct {
	UDPAddr string 	`json:"udp_addr"` 
	Pos [2]int		`json:"pos"`
	Id string		`json:"id"`
}

const PORT = 3000

func startUDPServer(port int) (*net.UDPConn, error) { 
	// listen to udp packets on localhost:3000	
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP: net.ParseIP("0.0.0.0"),
	})

	if err != nil {
		return nil, err
	}
		
	// return connection  
	return conn, nil
}

func sendResponseToUdp(conn *net.UDPConn, remoteAddr *net.UDPAddr, payload []byte) error {
	_, _, err := conn.WriteMsgUDP(payload, nil, remoteAddr)
	if err != nil {
		return err	
	}

	return nil
}


func monitor(conn *net.UDPConn, done chan bool) {
	defer conn.Close()
	
	// 2MB buffer
	payload := make([]byte, 2048)
	
	fmt.Println("Entered monitor()")
	for {
		n, remoteAddr, err := conn.ReadFromUDP(payload)	
			
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read from connection: %v, error: %v", conn.LocalAddr().String(), err) 
		}

		if n == 0 {	
			continue
		}
		payload = bytes.Trim(payload, "\x00")

		fmt.Fprintf(os.Stdout, "Message recieved from connection %v : %v \n", remoteAddr.String(), string(payload))	
		
		if string(payload) == "0" {
			// new player, waiting for connection
			new_id := uuid.NewString()
			err := sendResponseToUdp(conn, remoteAddr, []byte(new_id))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not send data to %v: error: %v", remoteAddr.String(), err)
			}
		} else if strings.HasPrefix(string(payload), "1") {
			// 1: Connect <udp add>	
			// TODO 
		} else {
			// establish connection or provide communication between already established conn 
			// 2: id <Data>	
			// TODO
			// unmarshal into player struct
			var p Player 
			// addr := []byte(fmt.Sprintf("udp_addr: %v", remoteAddr))
			// payload = append(payload, addr...)
			err := json.Unmarshal(payload, &p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not Unmarshal json into Player: error %v", err)
			}	
			fmt.Println(p)
		}	
	}	

	done <- true
}


func main() {
	
	// hold all the games(pair of udp connections) that are ongoing 
	// games := make(map[string]Game, 0)
	conn, err := startUDPServer(PORT)	
	if err != nil {	 
		fmt.Fprintf(os.Stderr, "Could not open port: %v on localhost, error: %v", PORT, err)
		return
	} else {
		fmt.Fprintf(os.Stdout, "Opened a UDP port: %v\n", conn.LocalAddr().String())
	}
	
	done_chan := make(chan bool)
	go monitor(conn, done_chan)

	// fire up client to send data to host	
	go startClient()
	
	<- done_chan
}
