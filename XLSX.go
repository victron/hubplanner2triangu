package main

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func createXLSX(data *[]S_record, total *Total, options options) {
	xlsx := excelize.NewFile()
	sheetName := (*data)[0].Resource_Name
	xlsx.SetSheetName("Sheet1", sheetName)

	// header
	headerStyle, err := xlsx.NewStyle(cellHeader)
	check(err)
	currentRow := 1
	for i, head := range headerSheet {
		column := excelize.ToAlphaString(i)
		row := strconv.Itoa(currentRow)
		axis := simpleJoin(column, row)
		xlsx.SetCellValue(sheetName, axis, head)

	}
	xlsx.SetCellStyle(sheetName, "A1", "J1", headerStyle)

	// column width
	xlsx.SetColWidth(sheetName, "A", "B", 15)
	xlsx.SetColWidth(sheetName, "G", "J", 22)

	// fill data
	for _, s_record := range *data {
		currentRow += 1
		row := strconv.Itoa(currentRow)
		date := (s_record.Date).Format("2006-01-02")
		xlsx.SetCellValue(sheetName, "A"+row, date)
		xlsx.SetCellValue(sheetName, "B"+row, s_record.Resource_Name)
		xlsx.SetCellValue(sheetName, "C"+row, s_record.Project_Name)
		xlsx.SetCellValue(sheetName, "D"+row, s_record.Project_Status)
		xlsx.SetCellValue(sheetName, "E"+row, s_record.Category)
		xlsx.SetCellValue(sheetName, "F"+row, s_record.Actual_Time)
		xlsx.SetCellValue(sheetName, "G"+row, s_record.Note)
		xlsx.SetCellValue(sheetName, "H"+row, s_record.OT10)
		xlsx.SetCellValue(sheetName, "I"+row, s_record.OT15)
		xlsx.SetCellValue(sheetName, "J"+row, s_record.OT20)
		// style
		cellSyle, err := xlsx.NewStyle(s_record.cellStyle)
		check(err)
		xlsx.SetCellStyle(sheetName, "A"+row, "J"+row, cellSyle)
	}

	// total
	currentRow += 1
	row := strconv.Itoa(currentRow)
	xlsx.SetCellFormula(sheetName, fmt.Sprintf("F%d", currentRow), fmt.Sprintf("SUM(F2:F%d)", currentRow-1))
	// put calculated value into cell (v tag), (warning exell can recalculate it accorging formula)
	xlsx.SetCellValue(sheetName, simpleJoin("F", row), (*total).Actual_Time) // v tag
	xlsx.SetCellValue(sheetName, simpleJoin("G", row), "TOTAL:")
	xlsx.SetCellFormula(sheetName, fmt.Sprintf("H%d", currentRow), fmt.Sprintf("SUM(H2:H%d)", currentRow-1))
	xlsx.SetCellValue(sheetName, simpleJoin("H", row), (*total).OT10) // v tag
	xlsx.SetCellFormula(sheetName, fmt.Sprintf("I%d", currentRow), fmt.Sprintf("SUM(I2:I%d)", currentRow-1))
	xlsx.SetCellValue(sheetName, simpleJoin("I", row), (*total).OT15) // v tag
	xlsx.SetCellFormula(sheetName, fmt.Sprintf("J%d", currentRow), fmt.Sprintf("SUM(J2:J%d)", currentRow-1))
	xlsx.SetCellValue(sheetName, simpleJoin("J", row), (*total).OT20) // v tag

	cellSyle, err := xlsx.NewStyle(cellTotal)
	check(err)
	xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("J%d", currentRow), cellSyle)

	// fotter as header
	currentRow += 1
	for i, head := range headerSheet {
		column := excelize.ToAlphaString(i)
		row := strconv.Itoa(currentRow)
		axis := column + row
		xlsx.SetCellValue(sheetName, axis, head)
	}
	xlsx.SetCellStyle(sheetName, simpleJoin("A", strconv.Itoa(currentRow)), simpleJoin("J", strconv.Itoa(currentRow)), headerStyle)

	fileName := sheetName + "_" + options.reportPeriod.Format("2006-01") + ".xlsx"
	err = xlsx.SaveAs(fileName)
	if err != nil {
		fmt.Println(err)
	}
}
