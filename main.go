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

func saveDataToGob() (err error) {
	gobFile, err := os.Create("data.gob")
	if err != nil {
		fmt.Println("Failed to create data.gob: " + err.Error())
		return
	}
	dataEncoder := gob.NewEncoder(gobFile)
	dataEncoder.Encode(data)

	gobFile.Close()
	return
}

func loadGob() (err error) {
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
	return
}

func writeToClient(conn net.Conn, text []byte) {
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
				writeToClient(conn, []byte("OK"))
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
					writeToClient(conn, []byte("OK"))
				} else {
					writeToClient(conn, []byte("SYNTAX INVALID"))
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
					writeToClient(conn, []byte("OK"))
				} else {
					writeToClient(conn, []byte("SYNTAX INVALID"))
				}

			case "get":
				if len(chunks) == 2 {
					values := get(chunks[1])
					output, err := json.Marshal(values)
					if err != nil {
						writeToClient(conn, []byte("FAIL"))
					} else {
						writeToClient(conn, output)
					}
				} else {
					writeToClient(conn, []byte("SYNTAX INVALID"))
				}
			case "fpop":
				if len(chunks) == 2 {
					value := popFront(chunks[1])
					output, err := json.Marshal(value)
					if err != nil {
						writeToClient(conn, []byte("FAIL"))
					} else {
						writeToClient(conn, output)
					}
				} else {
					writeToClient(conn, []byte("SYNTAX INVALID"))
				}

			case "bpop":
				if len(chunks) == 2 {
					value := popBack(chunks[1])
					output, err := json.Marshal(value)
					if err != nil {
						writeToClient(conn, []byte("FAIL"))
					} else {
						writeToClient(conn, output)
					}
				} else {
					writeToClient(conn, []byte("SYNTAX INVALID"))
				}
			case "keys":
				search := ""
				if len(chunks) >= 2 {
					search = chunks[1]
				}
				keys := getKeys(search)
				output, err := json.Marshal(keys)
				if err != nil {
					writeToClient(conn, []byte("FAIL"))
				} else {
					writeToClient(conn, output)
				}

			case "empty":
				if len(chunks) == 2 {
					empty(chunks[1])
					writeToClient(conn, []byte("OK"))
				} else {
					writeToClient(conn, []byte("SYNTAX INVALID"))
				}

			case "deletekey":
				if len(chunks) == 2 {
					deleteKey(chunks[1])
					writeToClient(conn, []byte("OK"))
				}
			case "loaddata":
				err := loadGob()
				if err != nil {
					writeToClient(conn, []byte("FAIL"))
				} else {
					writeToClient(conn, []byte("OK"))
				}

			case "savedata":
				err := saveDataToGob()
				if err != nil {
					writeToClient(conn, []byte("FAIL"))
				} else {
					writeToClient(conn, []byte("OK"))
				}

			case "autosave":
				if len(chunks) == 2 {
					if strings.ToLower(chunks[1]) == "false" {
						settings.UseDiskWriter = false
						writeToClient(conn, []byte("OK"))
					} else if strings.ToLower(chunks[1]) == "true" {
						settings.UseDiskWriter = true
						writeToClient(conn, []byte("OK"))
					} else {
						writeToClient(conn, []byte("SYNTAX INVALID"))
					}
				} else {
					output, err := json.Marshal(strconv.FormatBool(settings.UseDiskWriter))
					if err != nil {
						writeToClient(conn, []byte("FAIL"))
					} else {
						writeToClient(conn, output)
					}
				}

			default:
				writeToClient(conn, []byte("SYNTAX INVALID"))
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
		if len(data[key]) > 0 {
			result = data[key][0]
			data[key] = append(data[key][:0], data[key][1:]...)
		}
	}
	return
}

func popBack(key string) (result string) {
	_, ok := data[key]
	if ok {
		if len(data[key]) > 0 {
			result, data[key] = data[key][len(data[key])-1], data[key][:len(data[key])-1]
		}
	}
	return
}

func empty(key string) {
	_, ok := data[key]
	if ok {
		data[key] = nil
	}
}

func getKeys(search string) (keys []string) {
	for k := range data {
		if strings.Contains(k, search) {
			keys = append(keys, k)
		}
	}
	return
}

func deleteKey(key string) {
	_, ok := data[key]
	if ok {
		delete(data, key)
	}
}
