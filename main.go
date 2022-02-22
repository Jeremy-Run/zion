package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"
)

func main() {
	// External parameters
	factor := flag.Int64("factor", 1, "factor")
	dataSize := flag.Int("dataSize", 1, "dataSize")
	readTime := flag.Int("readTime", 1, "readTime")
	flag.Parse()

	// Initialize basic data
	c := InitConfig(*factor, *dataSize, *readTime)
	d := InitDict()
	m := InitMetrics()

	// Write data to dict
	for i := 0; i < c.DataSize; i++ {
		d.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i), m, c)
	}

	// Update dict metrics
	d.AllDB(m)

	// Read data by dict and calculate spend time
	startTime := time.Now()
	for j := 0; j < c.ReadTime; j++ {
		d.Get(m.LastEntry.Key)
	}
	elapsedTime := time.Since(startTime) / time.Millisecond
	fmt.Printf("Time: segment finished in %dms \n", elapsedTime)

	// Print config
	cj, _ := json.Marshal(c)
	fmt.Printf("Config: %s \n", cj)

	// Print metrics
	mj, _ := json.Marshal(m)
	fmt.Printf("Metrics: %s \n", mj)
	fmt.Printf("")
}
