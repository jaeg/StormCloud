package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

var data = make(map[string][]string)

func main() {
	fmt.Println("StormCloud running")
	ln, err := net.Listen("tcp", ":6464")
	if err != nil {
		fmt.Println("Error opening listener: " + err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming connection: " + err.Error())
		}
		go handleConnection(conn)
	}
}

func writeToClient(conn net.Conn, text string) {
	_, err := conn.Write([]byte(text))
	if err != nil {
		fmt.Println("Error writing to client: " + err.Error())
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("Client connected")
	buffer := make([]byte, 81920)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s", err)
			}
			break
		}
		bufferAsString := string(buffer)
		line := bufferAsString[0:n]
		if len(line) > 0 {
			if line == "quit" {
				writeToClient(conn, "OK")
				err := conn.Close()
				if err != nil {
					fmt.Println("Error disconnecting from client: " + err.Error())
				}
				break
			}
			chunks := strings.Split(line, " ")

			chunks[0] = strings.ToLower(chunks[0])
			switch chunks[0] {
			case "fpush":
				if len(chunks) >= 3 {
					newValue := chunks[2]
					if len(chunks) > 3 {
						for index, value := range chunks {
							if index > 2 {
								newValue += " " + value
							}
						}
					}
					pushFront(chunks[1], newValue)
					writeToClient(conn, "OK")
					continue
				}
			case "bpush":
				if len(chunks) >= 3 {
					newValue := chunks[2]
					if len(chunks) > 3 {
						for index, value := range chunks {
							if index > 2 {
								newValue += " " + value
							}
						}
					}
					pushBack(chunks[1], newValue)
					writeToClient(conn, "OK")
					continue
				}

			case "get":
				if len(chunks) == 2 {
					values := get(chunks[1])
					output := "Number of values: " + strconv.Itoa(len(values)) + "\r\n"
					for index, value := range values {
						output += "Value " + strconv.Itoa(index) + ": " + value + "\r\n"
					}
					writeToClient(conn, output)
					continue
				}
			case "popfront":
				if len(chunks) == 2 {
					value := popFront(chunks[1])
					writeToClient(conn, "Value: "+value)
					continue
				}

			case "popback":
				if len(chunks) == 2 {
					value := popBack(chunks[1])
					writeToClient(conn, "Value: "+value)
					continue
				}
			case "keys":
				keys := getKeys()
				output := "Number of keys: " + strconv.Itoa(len(keys)) + "\r\n"
				for _, key := range keys {
					output += " " + key + "\r\n"
				}
				writeToClient(conn, output)
				continue

			case "empty":
				if len(chunks) == 2 {
					empty(chunks[1])
					writeToClient(conn, "OK")
					continue
				}

			case "deletekey":
				if len(chunks) == 2 {
					deleteKey(chunks[1])
					writeToClient(conn, "OK")
					continue
				}
			}
			writeToClient(conn, "Syntax Invalid")
		}
	}
}

func pushFront(key string, value string) {
	data[key] = append([]string{value}, data[key]...)
}

func pushBack(key string, value string) {
	data[key] = append(data[key], value)
}

func get(key string) (result []string) {
	_, ok := data[key]
	if ok {
		result = data[key]
	}
	return
}

func popFront(key string) (result string) {
	_, ok := data[key]
	if ok {
		result, data[key] = data[key][0], data[key][:0]
	}
	return
}

func popBack(key string) (result string) {
	_, ok := data[key]
	if ok {
		result, data[key] = data[key][len(data[key])-1], data[key][:len(data[key])-1]
	}
	return
}

func empty(key string) {
	_, ok := data[key]
	if ok {
		data[key] = nil
	}
}

func getKeys() (keys []string) {
	for k := range data {
		keys = append(keys, k)
	}
	return
}

func deleteKey(key string) {
	_, ok := data[key]
	if ok {
		delete(data, key)
	}
}
