package updf

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	com "hubplanner2triangu/pkg/common"

	"github.com/victron/pdf"
)

const tax_rate = 5.0
const UAH = '₴'

type UberReport struct {
	date        time.Time
	provider    string
	total       float64
	description string
	notes       string
}

type UReports struct {
	data   []UberReport
	period time.Time
	cwd    string // curent working dir
	dir    string
}

func (ur *UReports) Len() int           { return len((*ur).data) }
func (ur *UReports) Swap(i, j int)      { (*ur).data[i], (*ur).data[j] = (*ur).data[j], (*ur).data[i] }
func (ur *UReports) Less(i, j int) bool { return (*ur).data[i].date.Before((*ur).data[j].date) }

func NewData(period time.Time, cwd, dir string) *UReports {
	return &UReports{period: period, cwd: cwd, dir: dir, data: nil}
}

func readPdf(fileName string, ureports *UReports, wg *sync.WaitGroup, mu *sync.Mutex) {
	/* worker to read PDF file */

	f, r, err := pdf.Open(fileName)
	com.Check(err)
	defer f.Close()
	defer wg.Done()

	b, err := r.GetPlainText()
	com.Check(err)
	scan := bufio.NewScanner(b)
	scan.Split(bufio.ScanLines)
	ureport := UberReport{provider: "uber"}
	for lineNum, total := 0, false; scan.Scan(); lineNum++ {
		line := scan.Text()
		if lineNum == 1 {
			// Fri, Jan 25, 2019

			docTime, err := time.Parse("Mon, Jan 2, 2006", line)
			com.Check(err)
			firstOfMonth := (*ureports).period
			lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

			if (docTime.After(firstOfMonth) && docTime.Before(lastOfMonth)) || docTime.Equal(firstOfMonth) || docTime.Equal(lastOfMonth) {
				ureport.date = docTime
				continue
			} else {
				fmt.Printf("document %s out of report period\n", fileName)
				return
			}
		}
		if strings.HasPrefix(line, "Total") {
			total = true
			continue
		}
		if total {
			if strings.HasPrefix(line, string(UAH)) {
				line = strings.Replace(line, string(UAH), "", 1)
				ureport.total, err = strconv.ParseFloat(line, 64)
				com.Check(err)
				break
			}
		}
	}
	(*mu).Lock()
	(*ureports).data = append((*ureports).data, ureport)
	(*mu).Unlock()
}

func (ur *UReports) ReadPDFs() (int, error) {
	ch_record := make(chan UberReport)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	expDirFullPath := filepath.Join((*ur).cwd, (*ur).dir)
	files, e := ioutil.ReadDir(expDirFullPath)
	com.Check(e)

	numPdfs := 0
	for _, file := range files {
		fileName := file.Name()
		fileNameFull := filepath.Join(expDirFullPath, fileName)
		if strings.HasSuffix(fileName, ".pdf") && strings.HasPrefix(fileName, "receipt") {
			numPdfs++
			wg.Add(1)
			go readPdf(fileNameFull, ur, wg, mu)
		}
	}

	wg.Wait()
	close(ch_record)

	numRecords := len((*ur).data)

	if numRecords != numPdfs {
		return numRecords, nil
		// TODO: add checking that files not parsed in reason out of report period
		// return numPdfs - numRecords, errors.New("some files not parsed")
	}
	return numRecords, nil
}

func (ur *UReports) PrepareReport(shortComment string) (int, error) {
	if len((*ur).data) == 0 {
		return 0, nil
	}
	sort.Sort(ur)
	total := 0.0
	docsNum := 0
	fmt.Println("TAXI in " + (*ur).period.Format("2006-01") + shortComment + ":")
	for i, report := range (*ur).data {
		// TODO: change currency (if need later)
		fmt.Printf("%d %s \t %.2f UAH\n", i+1, report.date.Format("2006-01-02"), report.total)
		total += report.total
		docsNum++
	}
	totalAndTax := total * 100 / (100 - tax_rate)
	fmt.Println("--------")
	fmt.Printf("TOTAL: \t\t %.2f + TAX = %.2f UAH\n", total, totalAndTax)

	return docsNum, nil
}
