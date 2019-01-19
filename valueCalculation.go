package main

import (
	"math"
	"time"
	// "github.com/rhymond/go-money"   //https://github.com/Rhymond/go-money
)

type monthInfo struct {
	woDays  int // working
	woHours float64
	weDays  int // weekend
	weHours float64
	hoDays  int // public holiday
	hoHours float64
}

type Total struct {
	Actual_Time float64
	OT10        float64 // including vacation
	OT15        float64
	OT20        float64
	OT10val     float64 // including vacation
	OT15val     float64
	OT20val     float64
	OT10valUAH  float64 // including vacation
	OT15valUAH  float64
	OT20valUAH  float64
	TotalUAH    float64
}

func (total *Total) total() {
	for _, s_record := range data {
		(*total).OT10 = (*total).OT10 + s_record.OT10
		(*total).OT15 = (*total).OT15 + s_record.OT15
		(*total).OT20 = (*total).OT20 + s_record.OT20
		(*total).Actual_Time = (*total).Actual_Time + s_record.Actual_Time
	}
}

func monthInfoConstruct(reportPeriod time.Time) monthInfo {
	var monthInfo monthInfo
	firstOfMonth := reportPeriod
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	for day := firstOfMonth; day.Before(firstOfNextMonth); {
		// fmt.Println("day=", day)
		switch {
		case day.Weekday() == time.Saturday || day.Weekday() == time.Sunday:
			monthInfo.weDays += 1
			monthInfo.weHours += 8
		default:
			monthInfo.woDays += 1
			monthInfo.woHours += 8
		}
		day = day.AddDate(0, 0, 1)
	}
	return monthInfo
}

func (monthInfo monthInfo) oneHourCost(monthRate float64) float64 {
	// return: one OT10 working hour cost (can be any currency)
	return monthRate / monthInfo.woHours
}

func (total *Total) calcValues(hourCost, usduah float64) {
	(*total).OT10val = (*total).OT10 * hourCost
	(*total).OT10val = math.Round((*total).OT10val*100) / 100
	(*total).OT15val = (*total).OT15 * hourCost * 1.5
	(*total).OT15val = math.Round((*total).OT15val*100) / 100
	(*total).OT20val = (*total).OT20 * hourCost * 2
	(*total).OT20val = math.Round((*total).OT20val*100) / 100

	(*total).OT10valUAH = (*total).OT10val * usduah
	(*total).OT10valUAH = math.Round((*total).OT10valUAH*100) / 100
	(*total).OT15valUAH = (*total).OT15val * usduah
	(*total).OT15valUAH = math.Round((*total).OT15valUAH*100) / 100
	(*total).OT20valUAH = (*total).OT20val * usduah
	(*total).OT20valUAH = math.Round((*total).OT20valUAH*100) / 100

	(*total).TotalUAH = (*total).OT10valUAH + (*total).OT15valUAH + (*total).OT20valUAH
	(*total).TotalUAH = math.Round((*total).TotalUAH*100) / 100

}
