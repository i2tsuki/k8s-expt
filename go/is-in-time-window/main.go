package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
}

func main() {
	now := time.Now().UTC()
	window := "Tue:10:00-Mon:23:30"
	in, err := isInWindow(now, window)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	println(in)
}

// Ref: https://cs.opensource.google/go/go/+/release-branch.go1.21:src/time/format.go;l=371
func match(s1, s2 string) bool {
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 {
			// Switch to lower-case; 'a'-'A' is known to be a single bit.
			c1 |= 'a' - 'A'
			c2 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' {
				return false
			}
		}
	}
	return true
}

// Ref: https://cs.opensource.google/go/go/+/release-branch.go1.21:src/time/format.go;l=387
func lookupWeekday(shortDay string) (time.Weekday, error) {
	for i, v := range shortDayNames {
		if len(shortDay) >= len(v) && match(shortDay[0:len(v)], v) {
			return time.Weekday(i), nil
		}
	}
	// Ref: cs.opensource.google/go/go/+/release-branch.go1.21:src/time/format.go;l=816
	return time.Weekday(0), errors.New("bad value for field") // placeholder not passed to user
}

func isInWindow(t time.Time, window string) (bool, error) {
	period := strings.Split(window, "-")

	start, err := time.Parse("Mon:15:04", period[0])
	if err != nil {
		return false, err
	}
	startWeekday, err := lookupWeekday(strings.Split(period[0], ":")[0])
	if err != nil {
		return false, err
	}

	end, err := time.Parse("Mon:15:04", period[1])
	if err != nil {
		return false, err
	}
	endWeekday, err := lookupWeekday(strings.Split(period[1], ":")[0])
	if err != nil {
		return false, err
	}

	if int(startWeekday) > int(endWeekday) {
		return false, errors.New("specify the weekday in the window Sun-Sat")
	}

	start = start.AddDate(t.Year(), int(t.Month()-1), t.Day()-1+int(startWeekday-t.Weekday()))
	end = end.AddDate(t.Year(), int(t.Month()-1), t.Day()-1+int(endWeekday-t.Weekday()))
	fmt.Printf("start: %s, end: %s, t: %s\n", start, end, t)
	if t.After(start) && t.Before(end) {
		return true, nil
	}
	return false, nil
}
