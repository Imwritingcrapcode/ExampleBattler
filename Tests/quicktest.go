package main

import (
	"fmt"
	"time"
)

func main() {
	var lastActivityTime int64
	lastActivityTime = 1571949034666704600
	timePassed := time.Now().UnixNano() - lastActivityTime
	fmt.Println(timePassed, timePassed > (15*time.Minute).Nanoseconds(), timePassed > (20*time.Minute).Nanoseconds(), timePassed > (65*time.Minute).Nanoseconds())
}
