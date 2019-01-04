package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var data []S_record

///////////////// sort by time ////////////////
type dataSortDate []S_record

func (sl dataSortDate) Len() int {
	return len(sl)
}

func (sl dataSortDate) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

func (sl dataSortDate) Less(i, j int) bool {
	return sl[i].Date.Before(sl[j].Date)
}

////////////////////////////////////////////////////
///////////////// sort by category ////////////////
type dataIndex struct {
	data  S_record
	index int
}

type dataSortCat []dataIndex

func (sl dataSortCat) Len() int {
	return len(sl)
}

func (sl dataSortCat) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

func (sl dataSortCat) Less(i, j int) bool {
	return !strings.HasPrefix(sl[i].data.Category, "Operation")
}

////////////////////////////////////////////////////

func sortIndex(indexes []int, data *[]S_record) ([]int, error) {
	// based on list of indexes put index with prefix "Operation" to the end of list
	// for correct working removeDuplicates
	var dataIndexs dataSortCat
	for _, i := range indexes {
		record := dataIndex{}
		record.data = (*data)[i]
		record.index = i
		dataIndexs = append(dataIndexs, record)
	}

	sort.Sort(dataSortCat(dataIndexs))
	var result []int
	for _, i := range dataIndexs {
		result = append(result, i.index)
	}
	return result, nil
}

func sameDateFinder(data *[]S_record) [][]int {
	// return list of lists with ids of record with same date
	var sameIndexs [][]int
	pDate := (*data)[0].Date
	sameIndex := []int{0}

	for i := 1; i < len(*data); i++ {
		if pDate == (*data)[i].Date {
			sameIndex = append(sameIndex, i)
			if i != len(*data)-1 {
				continue
			}
		}
		if len(sameIndex) > 1 {
			sameIndexs = append(sameIndexs, sameIndex)
		}
		sameIndex = nil
		sameIndex = append(sameIndex, i)
		pDate = (*data)[i].Date

	}
	return sameIndexs
}

func removeDuplicates(data []S_record, listDups [][]int) ([]S_record, error) {
	if listDups == nil {
		return data, nil
	}
	// additional check before remove
	// check: Actual_Time, Note
	var listToRemove []int
	for _, listDup := range listDups {
		listDup, err := sortIndex(listDup, &data)
		check(err)
	SmallList:
		for j, jj := range listDup {
			switch {
			case j == len(listDup)-1:
				// always leave one (any) record for any date
				break SmallList
			case data[jj].Actual_Time != 0:
				// ignore this index and add to remove others
				listToRemove = append(listToRemove, listDup[j+1:]...)
				break SmallList
			case data[jj].Note != "":
				// ignore this index and add to remove others
				listToRemove = append(listToRemove, listDup[j+1:]...)
				break SmallList
			case data[jj].Category == "Public Holiday" && data[jj].Actual_Time == 0:
				// remove Holiday with time == 0
				listToRemove = append(listToRemove, jj)
			case data[jj].Category == "Vacation" && data[jj].Actual_Time == 0:
				// remove Vacation with time == 0
				listToRemove = append(listToRemove, jj)
			case data[jj].Actual_Time == 0 && data[jj].Note == "":
				listToRemove = append(listToRemove, jj)
			}
		}
	}
	sort.Ints(listToRemove)
	for i := len(listToRemove) - 1; i >= 0; i-- {
		// fmt.Println("date=", data)
		// safe delete element from slise
		// https://github.com/golang/go/wiki/SliceTricks
		// fmt.Println(">>>1", data)
		copy(data[listToRemove[i]:], data[listToRemove[i]+1:])
		// fmt.Println(">>>2", data)
		data[len(data)-1] = S_record{}
		// fmt.Println(">>>3", data)
		data = data[:len(data)-1]
		// fmt.Println(">>>>", data)

	}
	return data, nil
}

func corrector(data *[]S_record) {
	// corector
	for i, _ := range *data {
		e := (*data)[i].correctOT10()
		check(e)
	}
}

func checker(data *[]S_record, reportPeriod time.Time) error {
	Resource_Name_0 := (*data)[0].Resource_Name
	firstOfMonth := reportPeriod
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	if len(*data) != lastOfMonth.Day() {
		return errors.New("number of records not eq to number of days in month")
	}
	for i, _ := range *data {
		e := (*data)[i].checker()
		check(e)

		//  check that all days of month is present
		if (*data)[i].Date != firstOfMonth.AddDate(0, 0, i) {
			fmt.Printf("expected: %s, got: %s", firstOfMonth.AddDate(0, 0, i).Format("2006-01-02"), (*data)[i].Date.Format("2006-01-02"))
			return errors.New("expected another date")
		}

		// check Resource_Name is same in all records
		if (*data)[i].Resource_Name != Resource_Name_0 {
			errString := fmt.Sprintf("Resouce name is not same; expected: %s, got: %s in date= %s \n",
				Resource_Name_0, (*data)[i].Resource_Name, (*data)[i].Date.Format("2006-01-02"))
			return errors.New(errString)
		}
	}
	return nil
}

func (data dataSortDate) removeOutOfReportPeriod(period time.Time) ([]S_record, error) {
	var newdata dataSortDate
	for i, _ := range data {
		if data[i].Date.Year() == period.Year() && data[i].Date.Month() == period.Month() {
			newdata = append(newdata, data[i])
		}
	}
	if len(newdata) == 0 {
		return newdata, errors.New("no DATA for report period. (check '-month ...' argument)")
	}
	return newdata, nil
}
