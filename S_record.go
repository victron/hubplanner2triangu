package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type S_record struct {
	Date           time.Time
	Resource_Name  string
	Project_Name   string
	Project_Status string
	Category       string
	Booked_Time    float64
	Actual_Time    float64
	Note           string
	OT10           float64
	OT15           float64
	OT20           float64
	comment        string
	cellStyle      string
}

func (S_record *S_record) parse(record []string, options options) error {
	/* parse record
	set default "cellStyle" for normal and weekend days
	*/
	dateTime, e := time.Parse("2006-01-02", record[0])
	check(e)
	if dateTime.Month() != options.reportPeriod.Month() || dateTime.Year() != options.reportPeriod.Year() {
		return errors.New("INFO: Date is out of report period, don't need to parse")
	}
	(*S_record).Date = dateTime
	(*S_record).Resource_Name = record[1]
	(*S_record).Project_Name = record[2]
	(*S_record).Project_Status = record[3]
	(*S_record).Category = record[4]
	value, err := strconv.ParseFloat(record[5], 64)
	check(err)
	(*S_record).Booked_Time = value
	value, err = strconv.ParseFloat(record[6], 64)
	(*S_record).Actual_Time = value
	(*S_record).Note = record[7]

	// default values (part of constructor)
	switch {
	case (*S_record).Category == "Vacation":
		(*S_record).cellStyle = cellVacation
	case (*S_record).Date.Weekday() == time.Sunday || (*S_record).Date.Weekday() == time.Saturday:
		(*S_record).cellStyle = cellWeekend
	case (*S_record).Category == "Public Holiday":
		(*S_record).cellStyle = cellHoliday
	default:
		(*S_record).cellStyle = cellNormal
	}
	// if (*S_record).Date.Weekday() == time.Sunday || (*S_record).Date.Weekday() == time.Saturday {
	// 	(*S_record).cellStyle = cellWeekend
	// } else {
	// 	(*S_record).cellStyle = cellNormal
	// }
	return nil
}

func (S_record *S_record) parseNotes() {
	if (*S_record).Note == "" {
		return
	}
	fields := strings.Split((*S_record).Note, fSeparator)
	if len(fields) == 0 {
		(*S_record).OT10 = (*S_record).Actual_Time
		return
	}
	for _, field := range fields {
		subfields := strings.Split(strings.TrimSpace(field), sfSeparator)
		if len(subfields) == 1 {
			continue
		}
		if len(subfields) != 2 {
			fmt.Printf("struct= %+v \n", *S_record)
			fmt.Printf("subfields= %+v \n", subfields)
			check(errors.New("expected 2 sub-field"))
		}

		key, val := strings.TrimSpace(subfields[0]), strings.TrimSpace(subfields[1])

		switch true {
		case testKey(key, keysOT10):
			value, e := strconv.ParseFloat(val, 64)
			check(e)
			(*S_record).OT10 = value
		case testKey(key, keysOT15):
			value, e := strconv.ParseFloat(val, 64)
			check(e)
			(*S_record).OT15 = value
		case testKey(key, keysOT20):
			value, e := strconv.ParseFloat(val, 64)
			check(e)
			(*S_record).OT20 = value
		case testKey(key, keys_comment):
			(*S_record).comment = strings.TrimSpace(val)
		}
	}
}

func (S_record *S_record) correctOT10() error {
	if (*S_record).Actual_Time != 0 {
		if (*S_record).OT10 == 0 {
			OT10 := (*S_record).Actual_Time - (*S_record).OT15 - (*S_record).OT20
			if OT10 >= 0 {
				(*S_record).OT10 = OT10
			} else {
				fmt.Printf("ERROR: Actual_Time= %f, OT15= %f, OT20= %f \n", (*S_record).Actual_Time, (*S_record).OT15, (*S_record).OT20)
				return errors.New("Actual_Time less then sum OT15 and OT20")
			}
		}
	}
	return nil
}

func (S_record *S_record) checker() error {
	switch false {
	case (*S_record).Actual_Time == (*S_record).OT10+(*S_record).OT15+(*S_record).OT20:
		fmt.Printf("ERROR: checker: date=%v Actual_Time= %f, OT10= %f, OT15= %f, OT20= %f \n",
			(*S_record).Date, (*S_record).Actual_Time, (*S_record).OT10, (*S_record).OT15, (*S_record).OT20)
		fmt.Printf("record= %+v \n", S_record)
		return errors.New("wrong sum OT10, OT15, OT20")
	}
	return nil
}

////////////////////////////////////////

func testKey(key string, list []string) bool {
	for _, i := range list {
		if strings.ToUpper(key) == i {
			return true
		}
	}
	return false
}
