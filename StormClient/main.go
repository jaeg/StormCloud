package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var ip = flag.String("ip", "localhost", "Location of StormCloud server.  DEFAULT: localhost")
var port = flag.String("port", "6464", "Port for StormCloud server.  DEFAULT: 6464")

func main() {
	flag.Parse()
	fmt.Println("Storm Client")
	fmt.Println("Connecting to " + *ip + ":" + *port)
	conn, err := net.Dial("tcp", *ip+":"+*port)
	if err != nil {
		fmt.Println("Error connecting to the Storm Cloud. " + err.Error())
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for fmt.Print("StormCloud> "); scanner.Scan(); fmt.Print("StormCloud> ") {
		//Write
		line := scanner.Text()
		_, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("Error write to stormcloud: " + err.Error())
		}

		//ReadResponse
		buffer := make([]byte, 81920)
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s", err)
			}
			break
		}
		bufferAsString := string(buffer)
		result := bufferAsString[0:n]

		fmt.Println(result)

	}

}
