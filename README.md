# performance_timer_go

Simple package used to time the use of functions over time to see where and when the time loss occurs

## Install

> go get github.com/Pineapples27/performance_timer_go

> import "github.com/Pineapples27/performance_timer_go"

## Usage

To start a timer use StartTimer() before a function:

```
performance_timer.StartTimer("someLargeFunction")
someLargeFunction()
```

To print a specific timer use PrintTimer(${key}) after the function:

```
performance_timer.StartTimer("someLargeFunction")
someLargeFunction()
performance_timer.PrintTimer("someLargeFunction")
```

To add a specific timer to a cumulative total to be printed later; then use GetTotalTime(${key}) after the function:

```
performance_timer.StartTimer("someLargeFunction")
someLargeFunction()
performance_timer.GetTotalTime("someLargeFunction")
```

To print all the total times; then use PrintTotalTime() after the function:

```
for{
    performance_timer.StartTimer("someLargeFunction")
    someLargeFunction()
    performance_timer.GetTotalTime("someLargeFunction")
    
    
    performance_timer.StartTimer("someOtherLargeFunction")
    someOtherLargeFunction()
    performance_timer.GetTotalTime("someOtherLargeFunction")
    
    
    performance_timer.PrintTotalTime()
}
```

To return a specific time; then use GetTime(${key}) after the function:

```
    performance_timer.StartTimer("someLargeFunction")
    someLargeFunction()
    elapsedTime := performance_timer.GetTime("someLargeFunction")
```

To return a specific time without removing it from the timer map; then use GetTimeWithoutDelete(${key}) after the
function:

```
    performance_timer.StartTimer("someLargeFunction")
    someLargeFunction()
    elapsedTime := performance_timer.GetTimeWithoutDelete("someLargeFunction")
```

If a lot of timers have been added to the service, you may want to turn them all off without removing all the code; then
use TurnTimerOff() at the init() of a service:

This is used in conjunction with the ...Override() methods which will do the same as their overridden function but
without worrying if the timer is turned on or not. This allows you to turn off all the timers except the few you are
actually interested in.

```
init() {
    performance_timer.TurnTimerOff()
}

for{
    performance_timer.StartTimer("someLargeFunction")
    someLargeFunction()
    performance_timer.GetTotalTime("someLargeFunction")
    
    
    performance_timer.StartTimerOverride("someOtherLargeFunction")
    someOtherLargeFunction()
    performance_timer.GetTotalTimeOverride("someOtherLargeFunction")
    
    
    performance_timer.PrintTotalTimeOverride()
}
```