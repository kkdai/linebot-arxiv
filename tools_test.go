package main

import (
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestTruncateString(t *testing.T) {
	str := "Hello, World"
	res := truncateString(str, 5)
	assert.Equal(t, "Hello", res, "Should return first 5 characters")

	res = truncateString(str, 20)
	assert.Equal(t, str, res, "Should return original string if maxLength > len(str)")
}

func TestInArray(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	val := 3
	exists, index := InArray(val, arr)
	assert.True(t, exists, "Should return true if val exists in array")
	assert.Equal(t, 2, index, "Should return correct index if val exists in array")

	val = 10
	exists, index = InArray(val, arr)
	assert.False(t, exists, "Should return false if val does not exist in array")
	assert.Equal(t, -1, index, "Should return -1 if val does not exist in array")
}

func TestRemoveStringItem(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e"}
	res := RemoveStringItem(arr, 2)
	assert.Equal(t, []string{"a", "b", "d", "e"}, res, "Should remove element at specified index")
}

func TestIsGroupEvent(t *testing.T) {
	event := &linebot.Event{
		Source: &linebot.EventSource{
			GroupID: "",
			RoomID:  "",
		},
	}
	res := isGroupEvent(event)
	assert.False(t, res, "Should return false if both GroupID and RoomID are empty")

	event.Source.GroupID = "group1"
	res = isGroupEvent(event)
	assert.True(t, res, "Should return true if GroupID is not empty")

	event.Source.GroupID = ""
	event.Source.RoomID = "room1"
	res = isGroupEvent(event)
	assert.True(t, res, "Should return true if RoomID is not empty")
}

func TestGetGroupID(t *testing.T) {
	event := &linebot.Event{
		Source: &linebot.EventSource{
			GroupID: "",
			RoomID:  "",
		},
	}
	res := getGroupID(event)
	assert.Equal(t, "", res, "Should return empty string if both GroupID and RoomID are empty")

	event.Source.GroupID = "group1"
	res = getGroupID(event)
	assert.Equal(t, "group1", res, "Should return GroupID if GroupID is not empty")

	event.Source.GroupID = ""
	event.Source.RoomID = "room1"
	res = getGroupID(event)
	assert.Equal(t, "room1", res, "Should return RoomID if RoomID is not empty")
}

func TestGetRandomIntSet(t *testing.T) {
	res := GetRandomIntSet(100, 10)
	assert.Equal(t, 10, len(res), "Should return slice of specified length")

	res = GetRandomIntSet(10, 5)
	assert.Equal(t, 5, len(res), "Should return slice of specified length")
}

func TestAddLineBreaksAroundURLs(t *testing.T) {
	input := "Check out this website https://example.com and this one http://another-example.com"
	expected := "Check out this website \nhttps://example.com\n and this one \nhttp://another-example.com\n"
	res := AddLineBreaksAroundURLs(input)
	assert.Equal(t, expected, res, "Should correctly insert line breaks around URLs")
}
