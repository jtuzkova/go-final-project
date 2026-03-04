package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AfterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid data format: %s", dstart)
	}
	slice := strings.Split(repeat, " ")

	if len(slice) == 0 {
		return "", fmt.Errorf("invalid data format")
	}

	var years, months, days int
	switch slice[0] {
	case "d":
		if len(slice) != 2 {
			return "", fmt.Errorf("interval in days not specified")
		}

		interval, err := strconv.Atoi(slice[1])
		if err != nil {
			return "", fmt.Errorf("error parsing string to int")
		}

		if interval <= 0 || interval >= 400 {
			return "", fmt.Errorf("interval must be in the range from 1 to 400")
		}

		days = interval

	case "y":
		years = 1
	
	default:
		return "", fmt.Errorf("unsupported format: %s", slice[0])
	}

	for {
    	date = date.AddDate(years, months, days)
    	if AfterNow(date, now) || date.Equal(now){
        	break
    	}
	}
	return date.Format(DateFormat), nil
}

func NextDayHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	strNow := req.FormValue("now")
	strDate := req.FormValue("date")
	repeat := req.FormValue("repeat")

	var now time.Time
	var err error
	if strNow != "" {
		now, err = time.Parse(DateFormat, strNow)
		if err != nil {
			http.Error(w, "invalid format 'now'", http.StatusBadRequest)
			return
		} 
	} else {
		now = time.Now()
	}

	nextDate, err := NextDate(now, strDate, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}