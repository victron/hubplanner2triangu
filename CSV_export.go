package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type exports struct {
	cwd      string   // current working dir
	expDir   string   // dir with files
	expFiles []string // list of files

}

func (exp *exports) initExp() error {
	cwd, e := os.Getwd()
	// (*exp).cwd = cwd
	exp.cwd = cwd
	check(e)
	exp.expDir = *monthStr
	expDirFullPath := filepath.Join((*exp).cwd, (*exp).expDir)
	files, e := ioutil.ReadDir(expDirFullPath)
	if e != nil {
		fmt.Println("problem with opening DIR:", expDirFullPath)
		os.Exit(1)
	}
	// TODO: it's better to move serch file logic into function of package
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".csv") && strings.HasPrefix(fileName, "Hub_Planner_Export_") {
			exp.expFiles = append(exp.expFiles, fileName)
		}
	}
	return nil
}

func readCSV(fileName string, exp *exports, data *[]S_record, wg *sync.WaitGroup, mu *sync.Mutex, options options) {
	/* worker to read CSV file */

	file, e := os.Open(fileName)
	check(e)
	defer file.Close()
	defer (*wg).Done()

	reader := csv.NewReader(file)

	// read header
	record, e := reader.Read()
	check(e)
	if len(record) != len(header) {
		fmt.Printf("headers in file= %d, expected= %d \n", len(record), len(header))
		return
	}
	for i, field := range header {
		if record[i] != field {
			fmt.Println("unexpected field=", record[i])
			return
		}
	}

	for {
		record, e := reader.Read()
		if e == io.EOF {
			break
		}
		check(e)
		s_record := new(S_record)
		e = s_record.parse(record, options)
		if e != nil {
			// don't append to data, parse another record
			continue
		}
		s_record.parseNotes()

		(*mu).Lock()
		*data = append(*data, *s_record)
		// chRecord <- *s_record
		(*mu).Unlock()
	}
}

func readCSVs(exp *exports, data *[]S_record, options options) (int, error) {
	ch_record := make(chan S_record)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for _, file := range (*exp).expFiles {
		fileName := filepath.Join((*exp).cwd, (*exp).expDir, file)
		wg.Add(1)
		go readCSV(fileName, exp, data, wg, mu, options)
	}

	wg.Wait()
	close(ch_record)

	numRecords := len(*data)

	if numRecords == 0 {
		return numRecords, errors.New("parsed num. of record == 0")
	}
	return numRecords, nil
}
