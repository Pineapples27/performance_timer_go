package performance_timer

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

var mapOfTimers map[string]*time.Time
var totalTimers map[string]float64
var timerLock sync.RWMutex
var totalTimerLock sync.RWMutex
var timerOff = false

//Set up timers
func init() {
	mapOfTimers = make(map[string]*time.Time)
}

//StartTimer is used to start a timer and add it to the map of timers but will be ignored if the timer is turned off
func StartTimer(key string) string {
	if timerOff {
		return key
	}
	return StartTimerOverride(key)
}

//StartTimerOverride is used to start a timer and add it to the map of timers
func StartTimerOverride(key string) string {
	timerLock.Lock()
	if mapOfTimers[key] != nil {
		newKey := findNewKey(key)
		key = newKey
	}
	timeNow := time.Now()
	mapOfTimers[key] = &timeNow
	timerLock.Unlock()
	return key
}

//PrintTimer is used to print a specific timer by key; but will be ignored if the timer is turned off
func PrintTimer(key, comment string) {
	if timerOff {
		return
	}
	PrintTimerOverride(key, comment)
}

//PrintTimerOverride is used to print a specific timer by key
func PrintTimerOverride(key, comment string) {
	timerLock.Lock()
	if mapOfTimers[key] != nil {
		elapsed := time.Since(*mapOfTimers[key])
		log.Printf("%s took %s -> %s\n", key, elapsed, comment)
	}
	timerLock.Unlock()
}

//PrintTotalTime is used to print all times in the timer map; but will be ignored if the timer is turned off
func PrintTotalTime() {
	if timerOff {
		return
	}
	PrintTotalTimeOverride()
}

//PrintTotalTimeOverride is used to print all times in the timer map
func PrintTotalTimeOverride() {
	totalTimerLock.Lock()
	printFormatted()
	totalTimerLock.Unlock()
}

//PrintTotalTimeRaw is used to print all times in the timer map without formatting
func PrintTotalTimeRaw() {
	if timerOff {
		return
	}
	totalTimerLock.Lock()
	printUnformatted()
	totalTimerLock.Unlock()
}

//GetTime is used to return a specific timer by key; but will be ignored if the timer is turned off
func GetTime(key string) (minutesSince float64) {
	if timerOff {
		return
	}
	return GetTimeOverride(key)
}

//GetTimeOverride is used to return a specific timer by key
func GetTimeOverride(key string) (minutesSince float64) {
	timerLock.Lock()
	if mapOfTimers[key] != nil {
		minutesSince = time.Since(*mapOfTimers[key]).Minutes()
		delete(mapOfTimers, key)
	}
	timerLock.Unlock()
	return minutesSince
}

//GetTimeWithoutDelete is used to return a specific timer by key without deleting it; but will be ignored if the timer is turned off
func GetTimeWithoutDelete(key string) (minutesSince float64) {
	if timerOff {
		return
	}
	return GetTimeWithoutDeleteOverride(key)
}

//GetTimeWithoutDeleteOverride is used to return a specific timer by key without deleting it
func GetTimeWithoutDeleteOverride(key string) (minutesSince float64) {
	timerLock.Lock()
	if mapOfTimers[key] != nil {
		minutesSince = time.Since(*mapOfTimers[key]).Minutes()
	}
	timerLock.Unlock()
	return minutesSince
}

//GetUnsafeTimeWithoutDelete is used to return a specific timer by key without deleting it and without concurrency locks which could make it unsafe;
//but will be ignored if the timer is turned off
func GetUnsafeTimeWithoutDelete(key string) float64 {
	if timerOff {
		return 0
	}
	return GetUnsafeTimeWithoutDeleteOverride(key)
}

//GetUnsafeTimeWithoutDeleteOverride is used to return a specific timer by key without deleting it and without concurrency locks which could make it unsafe
func GetUnsafeTimeWithoutDeleteOverride(key string) float64 {
	if mapOfTimers[key] == nil {
		return 0
	}
	return time.Since(*mapOfTimers[key]).Minutes()
}

//GetTotalTime is used to set all the current timers to a total timer map for printing later; but will be ignored if the timer is turned off
func GetTotalTime(key string) {
	if timerOff {
		return
	}
	GetTotalTimeOverride(key)
}

//GetTotalTimeOverride is used to set all the current timers to a total timer map for printing later; but will be ignored if the timer is turned off
func GetTotalTimeOverride(key string) {
	var minutesSince float64
	timerLock.Lock()
	if mapOfTimers[key] != nil {
		minutesSince = time.Since(*mapOfTimers[key]).Minutes()
		delete(mapOfTimers, key)
	}
	timerLock.Unlock()
	totalTimerLock.Lock()
	if totalTimers == nil {
		totalTimers = make(map[string]float64)
	}

	if strings.Contains(key, "___") {
		splitKey := strings.Split(key, "___")
		key = splitKey[0]
	}

	totalTimers[key] += minutesSince
	totalTimerLock.Unlock()
}

//TurnTimerOff is used to as a flag to quickly turn all the timers and logs off if there are already a lot of them and the users needs to see the program run without the logs or without the thread safe locks
func TurnTimerOff() {
	timerOff = true
}

func printFormatted() {
	orderedKeys := make([]string, 0)
	for s := range totalTimers {
		orderedKeys = append(orderedKeys, s)
	}
	sort.Strings(orderedKeys)
	for _, s := range orderedKeys {
		f := totalTimers[s]
		log.Println(fmt.Sprintf("%s:%.2f", s, f))
	}
}
func printUnformatted() {
	log.Println(totalTimers)
}
func findNewKey(key string) string {
	i := 0
	for {
		potentialNewKey := fmt.Sprintf("%s___%v", key, i)
		if mapOfTimers[potentialNewKey] == nil {
			return potentialNewKey
		}
		i++
	}
}
