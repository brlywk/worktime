package main

import (
	"testing"
	"time"
)

func TestSplitColon(t *testing.T) {
	inputList := []string{
		"",
		":",
		"4:",
		":2",
		"4:2",
	}

	// NOTE: Needs to be adjusted to defaultStartTime for fallback
	expectedList := [][]int{
		{defaultStartHour, defaultStartMinute},
		{defaultStartHour, defaultStartMinute},
		{4, 0},
		{defaultStartHour, 20},
		{4, 20},
	}

	for i, input := range inputList {
		l, r := splitColon(input)
		al := expectedList[i][0]
		ar := expectedList[i][1]

		if l != al || r != ar {
			t.Errorf("Mismatch for %v\nL is %v, expected %v\tR is %v expected %v", input, l, al, r, ar)
		}
	}
}

func TestParseTimeString(t *testing.T) {
	inputList := []string{
		"",
		":",
		"9:",
		":3",
		"9:30",
		"930",
		"14:30",
		"1430",
	}

	now := time.Now()

	expectedList := []time.Time{
		time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), defaultStartHour, defaultStartMinute, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), defaultStartHour, 30, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), 14, 30, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), 14, 30, 0, 0, now.Location()),
	}

	for i, input := range inputList {
		e := expectedList[i]
		a := parseStringTime(input)

		if a != e {
			t.Errorf("Mismatch. Got %v, expected %v", a, e)
		}
	}
}

func TestCalculateEnd(t *testing.T) {
	inputList := []struct {
		start string
		pause string
	}{
		{start: "8:30", pause: "30"},
		{start: "830", pause: "30"},
		{start: "8:30", pause: "60"},
		{start: "830", pause: "60"},
	}

	expectedList := []string{
		"17:00",
		"17:00",
		"17:30",
		"17:30",
	}

	for i, input := range inputList {
		e := expectedList[i]
		a := calculateEnd(input.start, input.pause)

		if a != e {
			t.Errorf("Mismatch. Expected %v, got %v", e, a)
		}

	}
}

func TestCalculateWorkingHours(t *testing.T) {
	inputList := []struct {
		start string
		pause string
		end   string
	}{
		{start: "8:30", pause: "30", end: "17:30"},
		{start: "830", pause: "30", end: "1730"},
		{start: "8:30", pause: "60", end: "16:30"},
		{start: "830", pause: "60", end: "1630"},
	}

	expectedList := []string{
		"8:30",
		"8:30",
		"7:00",
		"7:00",
	}

	for i, input := range inputList {
		e := expectedList[i]
		a := calculateWorkingHours(input.start, input.pause, input.end)

		if a != e {
			t.Errorf("Mismatch. Expected %v, got %v", e, a)
		}

	}
}
