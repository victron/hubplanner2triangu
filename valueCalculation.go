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
	vaDays  int // vacation
	vaHours float64
}

type Total struct {
	Actual_Time float64
	OT10        float64
	OT15        float64
	OT20        float64
	OT10USD     float64
	OT15USD     float64
	OT20USD     float64
	OT10UAH     float64
	OT15UAH     float64
	OT20UAH     float64
	TotalUSD    float64
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

// construct monthInfo structure,
// filling only working and weekend days
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

func (monthInfo *monthInfo) fillHolidays(data []S_record) {
	for _, s_record := range data {
		switch {
		case s_record.cellStyle == cellVacation:
			(*monthInfo).woDays -= 1
			(*monthInfo).woHours -= 8
			(*monthInfo).vaDays += 1
			(*monthInfo).vaHours += 8

		case s_record.cellStyle == cellHoliday:
			(*monthInfo).woDays -= 1
			(*monthInfo).woHours -= 8
			(*monthInfo).hoDays += 1
			(*monthInfo).hoHours += 8

		}
	}
}

// calculate one hour cost, based on va0 and ho flags
func (monthInfo monthInfo) oneHourCost(monthRate float64) float64 {
	var normHoursInMonth = monthInfo.woHours
	if !*vacation0 {
		normHoursInMonth += monthInfo.vaHours
	}
	if !*holiday0 {
		normHoursInMonth += monthInfo.hoHours
	}
	return math.Round(monthRate/normHoursInMonth*100) / 100
}

func (total *Total) calcValues(hourCost, usduah float64) {
	(*total).OT10USD = (*total).OT10 * hourCost
	(*total).OT10USD = math.Round((*total).OT10USD*100) / 100
	(*total).OT15USD = (*total).OT15 * hourCost * 1.5
	(*total).OT15USD = math.Round((*total).OT15USD*100) / 100
	(*total).OT20USD = (*total).OT20 * hourCost * 2
	(*total).OT20USD = math.Round((*total).OT20USD*100) / 100

	(*total).OT10UAH = (*total).OT10USD * usduah
	(*total).OT10UAH = math.Round((*total).OT10UAH*100) / 100
	(*total).OT15UAH = (*total).OT15USD * usduah
	(*total).OT15UAH = math.Round((*total).OT15UAH*100) / 100
	(*total).OT20UAH = (*total).OT20USD * usduah
	(*total).OT20UAH = math.Round((*total).OT20UAH*100) / 100

	(*total).TotalUSD = (*total).OT10USD + (*total).OT15USD + (*total).OT20USD
	(*total).TotalUSD = math.Round((*total).TotalUSD*100) / 100
	(*total).TotalUAH = (*total).TotalUSD * usduah
	(*total).TotalUAH = math.Round((*total).TotalUAH*100) / 100

}
