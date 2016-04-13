package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Storm Client")

	fmt.Println("Connecting to localhost:6464")
	conn, err := net.Dial("tcp", "localhost:6464")
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
