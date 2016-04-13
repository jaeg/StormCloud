package main

import (
	"fmt"
)

var data = make(map[string][]string)

func main() {
	fmt.Println("StormCloud running")

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

func deleteKey(key string) {
	_, ok := data[key]
	if ok {
		delete(data, key)
	}
}
