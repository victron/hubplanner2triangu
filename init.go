package main

import (
	"flag"
	"time"
)

var monthStr *string
var monthRate, usduah *float64
var versionS, versionL *bool

func init() {
	now := time.Now()
	monthStr = flag.String("month", now.Format("2006-01"), "for what month need to build report, format like:\n 2018-02 or 02")
	versionS = flag.Bool("v", false, "show version build and exit")
	versionL = flag.Bool("version", false, "show version build and exit")
	monthRate = flag.Float64("mrate", 0, "monthly negotiated rate in USD")
	usduah = flag.Float64("usd", 0, "exchange rate")
	// TODO: option for dst file creation path/name
}
