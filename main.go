package main

import (
	"fmt"
	"net"
	"os"
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
const IP = "127.0.0.1"

func startUDPServer(port int) (*net.UDPConn, error) { 
	// listen to udp packets on localhost:3000	
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP: net.ParseIP(IP),
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
	
	// 4MB buffer
	var payload []byte
	var id int = 0 
	
	fmt.Println("Entered monitor()")
	for {			
		payload = make([]byte, 2048)
		n, remoteAddr, err := conn.ReadFromUDP(payload)		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read from connection: %v, error: %v\n", remoteAddr.String(), err)
			fmt.Fprintf(os.Stderr, "Partial Bytes read: %v", n)
			continue	
		}

		if n == 0 {	
			continue
		}

		payload = bytes.Trim(payload, "\x00")

		fmt.Fprintf(os.Stdout, "Message recieved from connection %v : %v \n", remoteAddr.String(), string(payload))	
		
		if string(payload) == "Join" {
			// new player, waiting for connection
			newId, _ := json.Marshal(id)

			sendResponseToUdp(conn, remoteAddr, newId)
			id += 1
		} else {
			var p Player 
			// addr := []byte(fmt.Sprintf("udp_addr: %v", remoteAddr))
			// payload = append(payload, addr...)
			err = json.Unmarshal(payload, &p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not Unmarshal json into Player: error %v", err)
			}	
			// fmt.Println(p)
			
			// update player's stats
			(*players)[p.Id] = p
			
			fmt.Println("Player pool: ", *players)

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
	// go startClient()
	
	<- done_chan
}
