package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

// version
const version = "0.1.1.1"
const build = "2019-01-04"

/////////// export settings //////////////
const expDir = "export"
const fSeparator = ";"
const sfSeparator = ":"

var header = [...]string{"Date", "Resource Name", "Project Name", "Project Status", "Category", "Booked Time", "Actual Time", "Note"}
var keysOT10 = []string{"OT10", "OT1_0", "OT0"}
var keysOT15 = []string{"OT15", "OT1_5", "OT1.5", "OT1"}
var keysOT20 = []string{"OT20", "OT2_0", "OT2.0", "OT2"}
var keys_comment = []string{"C", "COMMENT", "COM"}

// xlsx styles
const borderPrefix = `"border":[{"type":"left","color":"#000000","style":1}, {"type":"right","color":"#000000","style":1}, {"type":"top","color":"#000000","style":1}, {"type":"bottom","color":"#000000","style":1}]}`
const cellHeader = `{"fill":{"type":"pattern","color":["#E0EBF5"],"pattern":1}}`
const cellTotal = `{"fill":{"type":"pattern","color":["#999999"],"pattern":1}}` // gray
// const cellWeekend = `{"fill":{"type":"pattern","color":["#FFFF00"],"pattern":1},` + borderPrefix  //yelow
const cellWeekend = `{"fill":{"type":"pattern","color":["#FFFF00"],"pattern":1}}` //yelow
const cellHoliday = `{"fill":{"type":"pattern","color":["#FF9900"],"pattern":1}}` // orange
// const cellVacation = `{"fill":{"type":"pattern","color":["#CCFF33"],"pattern":1},` + borderPrefix // green
const cellVacation = `{"fill":{"type":"pattern","color":["#CCFF33"],"pattern":1}}` // green
const cellNormal = `{"fill":{"type":"pattern","pattern":1}}`

var headerSheet = [...]string{"Date", "Resource Name", "Project Name", "Project Status", "Category",
	"Actual Time", "Note", "Normal hours: rate 1.0", "Overtime hours: rate 1.5", "Night work hours: rate 2.0"}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Total struct {
	Actual_Time float32
	OT10        float32
	OT15        float32
	OT20        float32
}

func (total *Total) total() {
	for _, s_record := range data {
		(*total).OT10 = (*total).OT10 + s_record.OT10
		(*total).OT15 = (*total).OT15 + s_record.OT15
		(*total).OT20 = (*total).OT20 + s_record.OT20
		(*total).Actual_Time = (*total).Actual_Time + s_record.Actual_Time
	}
}

func simpleJoin(strs ...string) string {
	var result string
	for _, str := range strs {
		result += str
	}
	return result
}

type options struct {
	reportPeriod time.Time
}

type paramSource interface {
	readParams(options *options) error
}

func main() {
	var options options
	var params paramSource
	params = &cliArgs{}
	err := params.readParams(&options)
	check(err)
	reportPeriod := options.reportPeriod

	exp := new(exports)
	exp.initExp()
	numRecords, err := readCSVs(exp, &data, options)
	if err != nil {
		if numRecords == 0 {
			fmt.Println("ERROR:", err)
			fmt.Println("no DATA for report period. (check '-month ...' argument)")
			os.Exit(1)
		}
	}

	// sorting and removing not related to report period data
	sort.Sort(dataSortDate(data))

	corrector(&data)
	data, err := removeDuplicates(data, sameDateFinder(&data))
	check(err)
	err = checker(&data, reportPeriod)
	check(err)
	total := new(Total)
	total.total()
	fmt.Printf("Total= %+v \n", *total)
	createXLSX(&data, total, options)
	os.Exit(0)
}
