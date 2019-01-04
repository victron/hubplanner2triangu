package main

import (
	"errors"
	"testing"
	"time"
)

type testPairSameDateFinder struct {
	data        []S_record
	dupLists    [][]int
	datacleaned []S_record
	dataError   error
}

var tests = []testPairSameDateFinder{
	{[]S_record{{Date: time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "", OT10: 0, OT15: 0, OT20: 0}},
		[][]int{{2, 3, 4}, {6, 7}},
		[]S_record{{Date: time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 0, OT15: 0, OT20: 0}},
		nil},

	{[]S_record{{Date: time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0}},
		[][]int{},
		[]S_record{{Date: time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 0, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0}},
		nil},

	{[]S_record{{Date: time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 0, OT15: 0, OT20: 0},
		{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "should not remove this", OT10: 0, OT15: 0, OT20: 0}},
		[][]int{{2, 3, 4}, {6, 7}},
		// this data should return error
		[]S_record{{Date: time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 8, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "OT10:0 OT15:0 OT20:0", OT10: 0, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 8, Note: "", OT10: 0, OT15: 0, OT20: 0},
			{Date: time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC), Actual_Time: 0, Note: "should not remove this", OT10: 0, OT15: 0, OT20: 0}},
		errors.New("record is duplicated, marked to remove but Notes exists")},
}

func TestSameDateFinder(t *testing.T) {

	for _, pair := range tests {
		res := sameDateFinder(&(pair.data))
		if len(res) != len(pair.dupLists) {
			t.Fatal("For:", pair.data,
				"\n expected:", len(pair.dupLists),
				"\n got:", len(res),
				"\n res:", res)

		}
		for i, sameIndex := range pair.dupLists {
			for j, v := range sameIndex {
				if res[i][j] != v {
					t.Fatal("Result:", res,
						"\n expected:", v,
						"\n got:", res[i][j])
				}
			}

		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	for _, pair := range tests {
		_, err := removeDuplicates(pair.data, pair.dupLists)
		if err != nil {
			if err.Error() != pair.dataError.Error() {
				t.Fatal("For:", pair.data,
					"\n expected:", pair.dataError, "got:", err)
			}

		}
	}

}
