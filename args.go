package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type cliArgs struct {
}

func (input *cliArgs) readParams(options *options) error {
	now := time.Now()
	flag.Parse()
	if *versionL || *versionS {
		fmt.Println("ver.:", version, "build:", build)
		os.Exit(0)
	}
	// values for rate calculation
	(*options).monthRate = *monthRate
	(*options).usduah = *usduah

	// monthStr - global pointer
	yearMonth := strings.Split(*monthStr, "-")
	if len(yearMonth) == 1 {
		month, err := strconv.Atoi(yearMonth[0])
		check(err)
		if 1 > month || month > 12 {
			return errors.New("month number out of range")
		}
		(*options).reportPeriod = time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		return nil
	}
	if len(yearMonth) == 2 {
		year, err := strconv.Atoi(yearMonth[0])
		check(err)
		month, err := strconv.Atoi(yearMonth[1])
		check(err)
		if year < 2018 || year > 2020 || 1 > month || month > 12 {
			return errors.New("year (2018-2020) or month out of range")
		}
		reportPeriod, err := time.Parse("2006-01", *monthStr)
		(*options).reportPeriod = reportPeriod
		return nil
	}
	return errors.New("unknown error")
}
