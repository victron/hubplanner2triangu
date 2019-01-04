package main

import (
	"flag"
	"time"
)

var monthStr *string
var versionS, versionL *bool

func init() {
	now := time.Now()
	monthStr = flag.String("month", now.Format("2006-01"), "for what moth need to build report, format like:\n 2018-02 or 02")
	versionS = flag.Bool("v", false, "show version build and exit")
	versionL = flag.Bool("version", false, "show version build and exit")
	// TODO: option for dst file creation path/name
}
