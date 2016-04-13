package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var data = make(map[string][]string)

func main() {
	fmt.Println("StormCloud running")
	scanner := bufio.NewScanner(os.Stdin)
	for fmt.Print("StormCloud > "); scanner.Scan(); fmt.Print("StormCloud > ") {
		line := scanner.Text()
		if len(line) > 0 {
			if line == "quit" {
				break
			}
			chunks := strings.Split(line, " ")
			chunks[0] = strings.ToLower(chunks[0])
			switch chunks[0] {
			case "fpush":
				if len(chunks) == 3 {
					pushFront(chunks[1], chunks[2])
					continue
				}
			case "bpush":
				if len(chunks) == 3 {
					pushBack(chunks[1], chunks[2])
					continue
				}
			case "popfront":
				if len(chunks) == 2 {
					value := popFront(chunks[1])
					fmt.Println("Value: " + value)
					continue
				}
			case "popback":
				if len(chunks) == 2 {
					value := popBack(chunks[1])
					fmt.Println("Value: " + value)
					continue
				}
			case "keys":
				keys := getKeys()
				fmt.Println("Number of Keys: " + strconv.Itoa(len(keys)))
				for _, key := range keys {
					fmt.Println(" " + key)
				}
				continue

			case "empty":
				if len(chunks) == 2 {
					empty(chunks[1])
					continue
				}

			case "deletekey":
				if len(chunks) == 2 {
					deleteKey(chunks[1])
					continue
				}
			}
			fmt.Println("Syntax Invalid")
		}
	}

}

func pushFront(key string, value string) {
	data[key] = append([]string{value}, data[key]...)
}

func pushBack(key string, value string) {
	data[key] = append(data[key], value)
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
