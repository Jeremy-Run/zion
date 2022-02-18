package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// Initialize the dict and metrics
	d := InitDict()
	t := InitMetrics(1, 262144)

	// write data to dict
	for i := 0; i < t.DataSize; i++ {
		d.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i), t)
	}

	// dict metrics
	d.AllDB(t)

	// read data by dict and calculate spend time
	startTime := time.Now()
	for j := 0; j < 1000000; j++ {
		d.Get(t.LastEntry.Key)
	}
	elapsedTime := time.Since(startTime) / time.Millisecond
	fmt.Printf("Segment finished in %dms \n", elapsedTime)

	// print metrics
	tj, _ := json.Marshal(t)
	fmt.Printf("%s \n", tj)
}
