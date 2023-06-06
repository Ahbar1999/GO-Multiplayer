package main 

import (
	"fmt"
	"net"
	"os"
	// "github.com/google/uuid"
	"encoding/json"
	// "strings"
)

func startClient() {
	// p := make([]byte, 2048)	
	var err error
	var  conn *net.UDPConn

	currAddr := net.UDPAddr{
		Port: 3001,
		IP: net.ParseIP("0.0.0.0"),
	}

	hostAddr := net.UDPAddr{
		Port: 3000,
		IP: net.ParseIP("0.0.0.0"),
	}	

	conn, err = net.DialUDP("udp", &currAddr, &hostAddr)
	if err != nil {	 
		fmt.Fprintf(os.Stderr, "Could not connect to HOST error: %v", err)
		return	
	}

	msg, _ := json.Marshal(Player{ currAddr.String(), [2]int{0, 0}, "Fake Id"})	
	_, err = conn.Write(msg)
	
	if err != nil {	 
		fmt.Fprintf(os.Stderr, "Could not write to HOST: %v, error: %\n", hostAddr.String(), err)
		return	
	}	
	
	fmt.Fprintln(os.Stdout, "Sent data to host successfully")
}
