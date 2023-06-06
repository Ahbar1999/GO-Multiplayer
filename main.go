package main

import (
	"fmt"
	"net"
	"os"
	"github.com/google/uuid"
	"encoding/json"
	// "strings"
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

func sendResponseToUdp(conn *net.UDPConn, remoteAddr *net.UDPAddr, payload []byte) {
	_, _, err := conn.WriteMsgUDP(payload, nil, remoteAddr)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not send data to %v: error: %v", remoteAddr.String(), err)
	}
}


func monitor(conn *net.UDPConn, players *map[string]Player, done chan bool) {
	
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
		
		if string(payload) == "Join" {
			// new player, waiting for connection
			sendResponseToUdp(conn, remoteAddr, []byte(uuid.NewString()))	
		} else {
			var p Player 
			// addr := []byte(fmt.Sprintf("udp_addr: %v", remoteAddr))
			// payload = append(payload, addr...)
			err = json.Unmarshal(payload, &p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not Unmarshal json into Player: error %v", err)
			}	
			fmt.Println(p)
			
			// update player's stats
			(*players)[p.Id] = p
			msg, err := json.Marshal(*players)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not Marshal Players: error %v", err)
			}

			// broadcast message
			for _, player := range *players {
				if player.Id == p.Id {
					continue	
				}
				playerRemoteAddr, err := net.ResolveUDPAddr("udp", player.UDPAddr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not resolve remote addr of %v error %v", player, err)
				}
				sendResponseToUdp(conn, playerRemoteAddr, msg)
			} 
		}	
	}	

	done <- true
}

func main() {
	
	// hold all the games(pair of udp connections) that are ongoing 	
	players := make(map[string]Player, 0)

	conn, err := startUDPServer(PORT)		

	if err != nil {	 
		fmt.Fprintf(os.Stderr, "Could not open port: %v on localhost, error: %v", PORT, err)
		return
	} else {
		fmt.Fprintf(os.Stdout, "Opened a UDP port: %v\n", conn.LocalAddr().String())
	}
	
	defer conn.Close()
	
	done_chan := make(chan bool)
	go monitor(conn, &players, done_chan)

	// fire up client to send data to host	
	go startClient()
	
	<- done_chan
}
