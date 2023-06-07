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

func startClient() {
	var err error
	var  conn *net.UDPConn
	// info about other player
	// uuid => Player
	var otherPlayers map[string]Player

	currAddr := net.UDPAddr{
		Port: 3001,
		IP: net.ParseIP("0.0.0.0"),
	}

	hostAddr := net.UDPAddr{
		Port: 3000,
		IP: net.ParseIP("0.0.0.0"),
	}	

	// this player
	thisPlayer := Player{currAddr.String(), [2]float32{0, 0}, uuid.NewString()}
	
	conn, err = net.DialUDP("udp", &currAddr, &hostAddr)
	if err != nil {	 
		fmt.Fprintf(os.Stderr, "Could not connect to HOST error: %v", err)
		return	
	}

	msg := make([]byte, 2048)
	
	// on startup 
	send, _ := json.Marshal(Player{ currAddr.String(), [2]float32{0, 0}, "None"})	
	_, err = conn.Write(send)
	
	if err != nil {	 
		fmt.Fprintf(os.Stderr, "Could not write to HOST: %v, error: %\n", hostAddr.String(), err)
		return	
	}
	
	// start comm
	for {
		// from now onwards ther
		n, remoteAddr, err := conn.ReadFromUDP(msg)
			
		msg = bytes.Trim(msg, "\x00")

		if n == 0 {
			continue	
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if remoteAddr.String() != hostAddr.String() {
			// cannot accept packets from any other address	
			continue
		}
		
		// update the player object 
		json.Unmarshal(msg, &otherPlayers)
		send, _ = json.Marshal(thisPlayer)

		conn.WriteToUDP(send, &hostAddr)
	}		
}
