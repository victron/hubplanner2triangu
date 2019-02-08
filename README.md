hubplanner2triangu
----------

Description:
=============

Get .csv files from `export` folder and create `.xlsx` file.

Usage:
=================

### Options:

`-month` - for what period prepare report.
Possible formats `2018-12`, `12`.

- if below options provided calculating total money values:  
`-mrate` - monthly negotiated rate in USD  
`-usd` - USDUAH exchange rate at day of payment   

### Testing Notes

#### Integrational test

cli to run for current test data  
`go test -args -month 2018-12`  
`go test -args -month 2018-12 -mrate 1000 -usd 27.9`  
`go test -args -month 2019-01`  

### Builds:

in releases tab.


### Notes
[Why Go Interfaces are Awesome](https://blog.teamtreehouse.com/go-interfaces-awesome)  
