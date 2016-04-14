package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type settingsStruct struct {
	Port                string `json:"port"`
	UseDiskWriter       bool   `json:"usediskwriter"`
	ReadFromDiskAtStart bool   `json:"readfromdiskatstart"`
}

var settings = &settingsStruct{
	Port: "6464", UseDiskWriter: false, ReadFromDiskAtStart: false}

var data = make(map[string][]string)

func main() {
	loadConfig()
	if settings.ReadFromDiskAtStart == true {
		loadGob()
	}

	fmt.Println("StormCloud running")
	fmt.Println("Operating on port " + settings.Port)
	ln, err := net.Listen("tcp", ":"+settings.Port)
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

func loadConfig() {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Unable to load config.json.  Using default settings.")
		return
	}

	err = json.Unmarshal(file, &settings)
	if err != nil {
		fmt.Println("Error parsing config file: " + err.Error())
	}
}

func saveDataToGob() {
	gobFile, err := os.Create("data.gob")
	if err != nil {
		fmt.Println("Failed to create data.gob: " + err.Error())
		return
	}
	dataEncoder := gob.NewEncoder(gobFile)
	dataEncoder.Encode(data)

	gobFile.Close()
}

func loadGob() {
	gobFile, err := os.Open("data.gob")
	if err != nil {
		fmt.Println("Failed to load data.gob: " + err.Error())
		return
	}

	dataDecoder := gob.NewDecoder(gobFile)
	err = dataDecoder.Decode(&data)
	if err != nil {
		fmt.Println("Failed to decode data.gob: " + err.Error())
	}

	gobFile.Close()
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
			log.Printf("Client disconnected. Read error: %s ", err)
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
				} else {
					writeToClient(conn, "Syntax Invalid")
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
				} else {
					writeToClient(conn, "Syntax Invalid")
				}

			case "get":
				if len(chunks) == 2 {
					values := get(chunks[1])
					output := "Number of values: " + strconv.Itoa(len(values)) + "\r\n"
					for index, value := range values {
						output += "Value " + strconv.Itoa(index) + ": " + value + "\r\n"
					}
					writeToClient(conn, output)
				} else {
					writeToClient(conn, "Syntax Invalid")
				}
			case "fpop":
				if len(chunks) == 2 {
					value := popFront(chunks[1])
					writeToClient(conn, "Value: "+value)
				} else {
					writeToClient(conn, "Syntax Invalid")
				}

			case "bpop":
				if len(chunks) == 2 {
					value := popBack(chunks[1])
					writeToClient(conn, "Value: "+value)
				} else {
					writeToClient(conn, "Syntax Invalid")
				}
			case "keys":
				keys := getKeys()
				output := "Number of keys: " + strconv.Itoa(len(keys)) + "\r\n"
				for _, key := range keys {
					output += " " + key + "\r\n"
				}
				writeToClient(conn, output)

			case "empty":
				if len(chunks) == 2 {
					empty(chunks[1])
					writeToClient(conn, "OK")
				} else {
					writeToClient(conn, "Syntax Invalid")
				}

			case "deletekey":
				if len(chunks) == 2 {
					deleteKey(chunks[1])
					writeToClient(conn, "OK")
				}
			case "loaddata":
				loadGob()
				writeToClient(conn, "OK")

			case "savedata":
				saveDataToGob()
				writeToClient(conn, "OK")

			case "autosave":
				if len(chunks) == 2 {
					if strings.ToLower(chunks[1]) == "false" {
						settings.UseDiskWriter = false
						writeToClient(conn, "OK")
					} else if strings.ToLower(chunks[1]) == "true" {
						settings.UseDiskWriter = true
						writeToClient(conn, "OK")
					} else {
						writeToClient(conn, "Syntax Invalid")
					}
				} else {
					writeToClient(conn, strconv.FormatBool(settings.UseDiskWriter))
				}

			default:
				writeToClient(conn, "Syntax Invalid")
			}
			if settings.UseDiskWriter == true {
				saveDataToGob()
			}
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
