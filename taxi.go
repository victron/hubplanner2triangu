package main

import (
	"fmt"
)

type pdfReader interface {
	ReadPDFs() (int, error)
}

func ReadPdfs(pdfs pdfReader) (int, error) {
	docNum, err := pdfs.ReadPDFs()
	fmt.Printf("parced= %d pdf documents\n", docNum)
	return docNum, err
}

// func taxi() {
// 	ureports := updf.NewData(reportPeriod, (*exp).cwd, (*exp).expDir)
// 	ureports.PrepareReport()
// }

// ureports := updf.NewData(reportPeriod, (*exp).cwd, (*exp).expDir)
// pdfs, err := ReadPdfs(ureports)
// check(err)
// fmt.Printf("parced= %d pdf documents\n", pdfs)
// ureports.PrepareReport()
