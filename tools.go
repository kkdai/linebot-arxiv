package main

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// InArray: Check if string item is in array
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// RemoveStringItem: Remove string item from slice
func RemoveStringItem(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func isGroupEvent(event *linebot.Event) bool {
	return event.Source.GroupID != "" || event.Source.RoomID != ""
}

func getGroupID(event *linebot.Event) string {
	if event.Source.GroupID != "" {
		return event.Source.GroupID
	} else if event.Source.RoomID != "" {
		return event.Source.RoomID
	}

	return ""
}

// GetRandomIntSet: Get random int set
func GetRandomIntSet(max int, count int) (randInts []int) {
	rand.Seed(time.Now().UnixNano())
	list := rand.Perm(max)
	randInts = list[:count]
	return randInts
}
