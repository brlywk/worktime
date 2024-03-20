package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

// Some default values to make this robust
const (
	defaultHoursPerDay   = 8
	defaultStartHour     = 8
	defaultStartMinute   = 30
	defaultPauseDuration = 60
)

var defaultStartTime = fmt.Sprintf("%v:%v", defaultStartHour, defaultStartMinute)

func main() {
	t := flag.Bool("t", false, "Runs the tool in \"Calculate (End) Time\" mode.")
	flag.Parse()

	timeFlag := *t

	fmt.Println("Super Awesome Worktime Calculator")
	fmt.Println("Note:\tTimes, e.g. 8:30 can be input with our without the colon.")
	fmt.Println(strings.Repeat("-", 50))

	var result, resultMessage string

	startTime, pauseDuration := getUserStartAndPause()

	if timeFlag {
		result = calculateEnd(startTime, pauseDuration)
		resultMessage = fmt.Sprintf("You will have worked %d hours at:\t%s", defaultHoursPerDay, result)
	} else {
		endTime, err := getUserEnd()
		if err != nil {
			// NOTE: Special case: If the end time cannot be read correctly, use end time calculation
			endTime = calculateEnd(startTime, pauseDuration)
		}
		result = calculateWorkingHours(startTime, pauseDuration, endTime)
		resultMessage = fmt.Sprint("Total working hours:\t\033[31m", result, " \033[0m")
	}

	fmt.Println()
	fmt.Println(resultMessage)
	fmt.Println()

	// initialise clipboard if available
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Clipboard functionality is not available :(")
	}

	clipboard.Write(clipboard.FmtText, []byte(result))
	fmt.Println("The result has been copied to your clipboard.")
}

// Split a string at the colon, making sure that right is always at least 2 characters
//
// Important: Does not check if there is a colon, but rather falls back to 'defaultStartTime'
// if no colon is found
func splitColon(str string) (left int, right int) {
	if len(str) < 1 || str == ":" {
		str = defaultStartTime
	}

	splitStr := strings.Split(str, ":")
	// log.Println("split", str, "->", splitStr)

	l := splitStr[0]
	left, err := strconv.Atoi(l)
	if err != nil {
		left = defaultStartHour
	}

	r := splitStr[1]
	right, err = strconv.Atoi(r)
	if err != nil {
		right = 0
	}

	if right < 10 {
		right *= 10
	}

	// log.Println("left", left, "right", right)
	return left, right
}

// parseStringTime tries to parse the given string into a time.Time struct
//
// Returns the current time with zeroed seconds and nanoseconds if the string is empty, or
// the defaultStartTime if the string only contains a colon (":")
func parseStringTime(str string) time.Time {
	// log.Println("[parseStringTime] parsing string", str)
	now := time.Now()

	if len(str) < 1 {
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	}

	hasColon := strings.Contains(str, ":")

	// Preprate string to be parsed
	if !hasColon {
		// (h)our and (m)inute needed for splitting and inserting colons
		var h, m string

		switch true {
		case len(str) == 1:
			str = fmt.Sprint(str, ":")

		case len(str) == 2:
			num, err := strconv.Atoi(str)
			if err != nil {
				str = defaultStartTime
			}

			// check if it's 'just an hour'
			if num <= 23 {
				str = fmt.Sprint(str, ":")
			} else {
				str = fmt.Sprint(str[0], ":", str[1])
			}

		case len(str) == 3:
			// NOTE: Also need to check case 3 characters, as it could be:
			// 830 -> 8:30  or 143 -> 14:30
			if str[0] == '1' || str[0] == '2' {
				h = str[:2]
				m = str[2:]
			} else {
				h = str[:1]
				m = str[1:]
			}

			str = fmt.Sprint(h, ":", m)

		case len(str) > 3:
			h = str[:len(str)-2]
			m = str[len(str)-2:]

			str = fmt.Sprint(h, ":", m)
		}
	}

	l, r := splitColon(str)
	t := time.Date(now.Year(), now.Month(), now.Day(), l, r, 0, 0, now.Location())

	// log.Println("[parseStringTime] string", str, "parsed to", t)

	return t
}

// calculateEnd calculates at what time defaultHoursPerDay is reached from startTime,
// taking into consideration the pauseDuration
//
// Returns a string in 24h format: %H:%M (e.g. 23:59)
func calculateEnd(startTime string, pauseDuration string) string {
	start := parseStringTime(startTime)
	pause, err := strconv.Atoi(pauseDuration)
	if err != nil {
		pause = defaultPauseDuration
	}

	end := start.Add(defaultHoursPerDay*time.Hour + time.Duration(pause)*time.Minute)

	return end.Format("15:04")
}

// calculateWorkingHours calculates the total number of hours from startTime to endTime,
// exluding pauseDuration
//
// Returns a string representing the total number of hours, rounding down to the
// nearest "15 minutes" (e.g. 8h25min -> 8:25 -> 8:15)
func calculateWorkingHours(startTime string, pauseDuration string, endTime string) string {
	start := parseStringTime(startTime)
	pause, err := strconv.Atoi(pauseDuration)
	if err != nil {
		pause = defaultPauseDuration
	}
	end := parseStringTime(endTime)

	fullDuration := end.Sub(start)
	duration := fullDuration - (time.Minute * time.Duration(pause))

	roundedDuration := duration.Round(time.Minute * 15)
	h := int(math.Floor(roundedDuration.Hours()))
	m := int(roundedDuration.Minutes() - float64(h)*60.0)

	return fmt.Sprintf("%d:%02d", h, m)
}

// getUserStartAndPause displays the input prompt for startTime and pauseDuration to
// the user
//
// Returns startTime and pauseDuration as a string
func getUserStartAndPause() (startTime, pauseDuration string) {
	var err error

	fmt.Print("Start Time:\t")
	_, err = fmt.Scanln(&startTime)
	if err != nil {
		startTime = defaultStartTime
	}

	fmt.Print("Pause (min):\t")
	_, err = fmt.Scanln(&pauseDuration)
	if err != nil {
		pauseDuration = fmt.Sprint(defaultPauseDuration)
	}

	return startTime, pauseDuration
}

// getUserEnd displays the input prompt for the endTime for duration calculation
//
// Returns the endTime as a string
func getUserEnd() (endTime string, err error) {
	fmt.Print("End Time:\t")

	_, err = fmt.Scanln(&endTime)
	if err != nil {
		return "", err
	}

	return endTime, nil
}
