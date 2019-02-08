package main

import (
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/xuri/excelize"
)

// TODO: clean created files after test

const testData = "testData"
const etalonFile = "Etalon.xlsx"
const etalonSheet = "Sheet1"
const reportFile = "Viktor Tsymbalyuk_2018-12.xlsx"
const sheetName = "Viktor Tsymbalyuk"

func TestAll(t *testing.T) {
	cwd, err := os.Getwd()
	check(err)
	testDir := path.Join(cwd, testData)
	etalonFileName := path.Join(testDir, etalonFile)
	reportFileName := path.Join(testDir, reportFile)
	err = os.Chdir(path.Join(cwd, testData))
	check(err)
	cwd, err = os.Getwd()
	check(err)
	if !strings.HasSuffix(cwd, testData) {
		t.Fatal("not test environment!!!!")
	}
	// exec main
	main()

	// open files to compare
	etalon, err := excelize.OpenFile(etalonFileName)
	if err != nil {
		t.Fatal("can't open etalon file", etalonFile, "\n")
	}

	report, err := excelize.OpenFile(reportFileName)
	if err != nil {
		t.Fatal("can't open report file", reportFile, "\n")
	}

	// compare
	for row, last_row := 1, false; last_row == false && row < 50; row++ {
		for i, _ := range headerSheet {
			aix := excelize.ToAlphaString(i) + strconv.Itoa(row)

			// stop condition
			if etalon.GetCellValue(etalonSheet, aix) == "Night work hours: rate 2.0" && row > 1 {
				last_row = true
			}

			if etalon.GetCellFormula(etalonSheet, aix) != report.GetCellFormula(sheetName, aix) {
				t.Fatal("aix:", aix, "expected Formula:", etalon.GetCellFormula(etalonSheet, aix), "got Formula:", report.GetCellFormula(sheetName, aix))
			}
			// value compare
			if etalon.GetCellValue(etalonSheet, aix) != report.GetCellValue(sheetName, aix) {
				t.Fatal("aix:", aix, "expected:", etalon.GetCellValue(etalonSheet, aix), "got:", report.GetCellValue(sheetName, aix))
			}

		}
	}
}
