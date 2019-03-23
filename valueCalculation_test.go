package main

import (
	"reflect"
	"testing"
	"time"
)

type testmonthInfoConstruct struct {
	input_reportPeriod time.Time
	res_expected       monthInfo
}

type testTotalCalc struct {
	hourCost    float64
	usduah      float64
	input_total Total
	exp_total   Total
}

var testsTableMonthInfo = []testmonthInfoConstruct{
	{
		time.Date(2018, 12, 1, 0, 0, 0, 0, time.UTC),
		// TODO: Holidays not in test currently
		monthInfo{woDays: 21, woHours: 168, weDays: 10, weHours: 80, hoDays: 0, hoHours: 0},
	},
	{
		time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		monthInfo{woDays: 23, woHours: 184, weDays: 8, weHours: 64, hoDays: 0, hoHours: 0},
	},
}

var testTableTotalCalc = []testTotalCalc{
	{
		5.9, 27.8,
		Total{Actual_Time: 15, OT10: 17, OT15: 8, OT20: 5, OT10USD: 0, OT15USD: 0, OT20USD: 0,
			OT10UAH: 0, OT15UAH: 0, OT20UAH: 0, TotalUSD: 0, TotalUAH: 0},
		Total{Actual_Time: 15, OT10: 17, OT15: 8, OT20: 5, OT10USD: 100.3, OT15USD: 70.8, OT20USD: 59,
			OT10UAH: 2788.34, OT15UAH: 1968.24, OT20UAH: 1640.2, TotalUSD: 230.1, TotalUAH: 6396.78},
	},
	{
		5.32, 27.34,
		Total{Actual_Time: 180, OT10: 160, OT15: 26, OT20: 14, OT10USD: 0, OT15USD: 0, OT20USD: 0,
			OT10UAH: 0, OT15UAH: 0, OT20UAH: 0, TotalUSD: 0, TotalUAH: 0},
		Total{Actual_Time: 180, OT10: 160, OT15: 26, OT20: 14, OT10USD: 851.2, OT15USD: 207.48, OT20USD: 148.96,
			OT10UAH: 23271.81, OT15UAH: 5672.50, OT20UAH: 4072.57, TotalUSD: 1207.64, TotalUAH: 33016.88},
	},
}

func TestMonthInfoConstruct(t *testing.T) {
	for _, pair := range testsTableMonthInfo {
		res := monthInfoConstruct(pair.input_reportPeriod)
		if !reflect.DeepEqual(res, pair.res_expected) {
			t.Fatal("Result for:", pair.input_reportPeriod,
				"\nGot:", res,
				"\nExpected:", pair.res_expected)
		}
	}
}

func TestCalcValues(t *testing.T) {
	for _, pair := range testTableTotalCalc {
		total := &(pair.input_total)
		total.calcValues(pair.hourCost, pair.usduah)
		if !reflect.DeepEqual(pair.exp_total, *total) {
			t.Fatal("Result for:", pair.input_total,
				"\nGot     :", *total,
				"\nExpected:", pair.exp_total)
		}
	}
}
