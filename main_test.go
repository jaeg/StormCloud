package main

import "testing"

func TestPushFront(t *testing.T) {
	pushFront("push", "front")
	if data["push"][0] != "front" {
		t.Error("PushFront Test failed")
	}
}

func TestPushBack(t *testing.T) {
	pushBack("push", "back")
	if data["push"][1] != "back" {
		t.Error("PushBack Test failed")
	}
}

func TestPopFront(t *testing.T) {
	pushFront("front", "back")
	pushFront("front", "front")
	value := popFront("front")
	if value != "front" {
		t.Error("PopFront Test failed")
	}
}

func TestPopBack(t *testing.T) {
	pushFront("back", "back")
	pushFront("back", "front")
	value := popBack("back")
	if value != "back" {
		t.Error("PopBack Test failed")
	}
}

func TestEmpty(t *testing.T) {
	empty("push")

	if len(data["push"]) > 0 {
		t.Error("Empty test failed")
	}
}

func TestGetKeys(t *testing.T) {
	keys := getKeys("")
	if len(keys) == 0 {
		t.Error("GetKeys test failed")
	}
}

func TestGetKeysSearch(t *testing.T) {
	keys := getKeys("ont")
	if len(keys) != 1 {
		t.Error("GetKeys test failed")
	}
}

func TestDeleteKey(t *testing.T) {
	deleteKey("push")
	_, ok := data["push"]
	if ok {
		t.Error("DeleteKey test failed")
	}
}
