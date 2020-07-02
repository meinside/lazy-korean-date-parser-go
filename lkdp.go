package lkdp

// Lazy Korean Date(+Time) Parser

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// constants
const (
	DefaultLocation = "Asia/Seoul"
)

// constant strings
const (
	ExpressionTheDayBeforeYesterday1 = `그저께`
	ExpressionTheDayBeforeYesterday2 = `그제`
	ExpressionYesterday1             = `어제`
	ExpressionYesterday2             = `작일`
	ExpressionToday1                 = `오늘`
	ExpressionToday2                 = `금일`
	ExpressionTomorrow1              = `내일`
	ExpressionTomorrow2              = `명일`
	ExpressionTheDayAfterTomorrow1   = `모레`
	ExpressionTwoDaysAfterTomorrow1  = `글피`

	ExpressionYear1  = `년`
	ExpressionYear2  = `年`
	ExpressionMonth1 = `월`
	ExpressionMonth2 = `月`
	ExpressionMonth3 = `개월`
	ExpressionDay1   = `일`
	ExpressionDay2   = `日`

	ExpressionPeriodAM1 = `오전`
	ExpressionPeriodAM2 = `AM`
	ExpressionPeriodPM1 = `오후`
	ExpressionPeriodPM2 = `PM`

	ExpressionHour1   = `시`
	ExpressionHour2   = `時`
	ExpressionHour3   = `:`
	ExpressionMinute1 = `분`
	ExpressionMinute2 = `分`
	ExpressionMinute3 = `:`
	ExpressionSecond1 = `초`
	ExpressionSecond2 = `秒`

	ExpressionTimeHour1   = `시간`
	ExpressionTimeMinute1 = ExpressionMinute1
	ExpressionTimeSecond1 = ExpressionSecond1

	ExpressionBefore1 = `전`
	ExpressionAfter1  = `후`
	ExpressionAfter2  = `뒤`

	ExpressionMinuteThirty = `반` // xx시 '반' = xx시 '30분'

	ExpressionDateSeparator1 = `\-`
	ExpressionDateSeparator2 = `\.`
	ExpressionDateSeparator3 = `/`
)

// Verbose flag for debugging
var Verbose bool

// Hms struct for hh:mm:ss
type Hms struct {
	Hours          int
	Minutes        int
	Seconds        int
	NumDaysChanged int

	Ambiguous bool // whether this time is ambiguous or not (eg: AM/PM)
}

var _location *time.Location

var dateExactRe1, dateExactRe2 *regexp.Regexp // 특정 일자
var dateRelRe1, dateRelRe2 *regexp.Regexp     // 상대 일자
var timeRelRe1 *regexp.Regexp                 // 상대 시간
var timeExactRe1, timeExactRe2 *regexp.Regexp // 특정 시간

func init() {
	_location, _ = time.LoadLocation(DefaultLocation)

	dateExactRe1 = regexp.MustCompile(fmt.Sprintf(`((\d{2,})\s*[%s])?\s*((\d{1,2})\s*[%s])?\s*(\d{1,2})\s*[%s]`,
		strings.Join([]string{
			ExpressionYear1,
			ExpressionYear2,
		}, ""),
		strings.Join([]string{
			ExpressionMonth1,
			ExpressionMonth2,
		}, ""),
		strings.Join([]string{
			ExpressionDay1,
			ExpressionDay2,
		}, ""),
	))
	dateExactRe2 = regexp.MustCompile(fmt.Sprintf(`((\d{2,})\s*[%s])?\s*((\d{1,2})\s*[%s]\s*(\d{1,2})\s*[%s]?)`,
		strings.Join([]string{
			ExpressionDateSeparator1,
			ExpressionDateSeparator2,
			ExpressionDateSeparator3,
		}, ""),
		strings.Join([]string{
			ExpressionDateSeparator1,
			ExpressionDateSeparator2,
			ExpressionDateSeparator3,
		}, ""),
		strings.Join([]string{
			ExpressionDateSeparator2,
		}, ""),
	))
	dateRelRe1 = regexp.MustCompile(fmt.Sprintf(`(\d+)\s*(%s)\s*(%s)`, strings.Join([]string{
		ExpressionYear1,
		ExpressionYear2,
		ExpressionMonth1,
		ExpressionMonth3,
		ExpressionDay1,
		ExpressionDay2,
	}, "|"), strings.Join([]string{
		ExpressionBefore1,
		ExpressionAfter1,
		ExpressionAfter2,
	}, "|")))
	dateRelRe2 = regexp.MustCompile(fmt.Sprintf(`(%s)`, strings.Join([]string{
		ExpressionTheDayBeforeYesterday1,
		ExpressionTheDayBeforeYesterday2,
		ExpressionYesterday1,
		ExpressionYesterday2,
		ExpressionToday1,
		ExpressionTomorrow1,
		ExpressionTomorrow2,
		ExpressionTheDayAfterTomorrow1,
		ExpressionTwoDaysAfterTomorrow1,
	}, "|")))
	timeRelRe1 = regexp.MustCompile(fmt.Sprintf(`(\d+)\s*(%s)\s*(%s)`,
		strings.Join([]string{
			ExpressionTimeHour1,
			ExpressionTimeMinute1,
			ExpressionTimeSecond1,
		}, "|"),
		strings.Join([]string{
			ExpressionBefore1,
			ExpressionAfter1,
			ExpressionAfter2,
		}, "|"),
	))
	timeExactRe1 = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)?\s*((\d{1,2})\s*[%s])\s*%s`,
		strings.Join([]string{
			ExpressionPeriodAM1,
			ExpressionPeriodAM2,
			ExpressionPeriodPM1,
			ExpressionPeriodPM2,
		}, "|"),
		strings.Join([]string{
			ExpressionHour1,
			ExpressionHour2,
			ExpressionHour3,
		}, "|"),
		ExpressionMinuteThirty,
	))
	timeExactRe2 = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)?\s*((\d{1,2})\s*[%s])\s*((\d{1,2})(\s*[%s]?(\d{1,2})\s*[%s]?)?)?`,
		strings.Join([]string{
			ExpressionPeriodAM1,
			ExpressionPeriodAM2,
			ExpressionPeriodPM1,
			ExpressionPeriodPM2,
		}, "|"),
		strings.Join([]string{
			ExpressionHour1,
			ExpressionHour2,
			ExpressionHour3,
		}, "|"),
		strings.Join([]string{
			ExpressionMinute1,
			ExpressionMinute2,
			ExpressionMinute3,
		}, "|"),
		strings.Join([]string{
			ExpressionMinute1,
			ExpressionMinute2,
		}, "|"),
	))
}

// SetLocation sets location
// 지역 설정 (timezone)
//
// https://golang.org/pkg/time/#Location
func SetLocation(str string) error {
	var err error

	_location, err = time.LoadLocation(str)

	return err
}

// ExtractDates extracts all dates from given string
//
// returns `nil` dates on error
//
// priority of regexs is:
//   dateRelRe1 > dateRelRe2 > dateExactRe1 > dateExactRe2
func ExtractDates(str string, ifEmptyFillAsToday bool) (dates map[string]time.Time, err error) {
	// initialize values
	dates = map[string]time.Time{}
	var year, month, day int = 0, 0, 0

	// indices of processed matches: not to extract duplicated matches
	alreadyProcessed := map[int]struct{}{}

	var matches []string
	matches = dateRelRe1.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := dateRelRe1.FindStringSubmatch(match)

			debugPrint("dateRelRe1: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			date := time.Now() // today

			number, _ := strconv.ParseInt(slices[1], 10, 16)

			multiply := 1
			switch slices[3] {
			case ExpressionBefore1: // before
				multiply = -1
			case ExpressionAfter1, ExpressionAfter2: // after
				// do nothing (+1)
			}
			switch slices[2] {
			case ExpressionYear1, ExpressionYear2: // year
				date = date.AddDate(multiply*int(number), 0, 0)
			case ExpressionMonth1, ExpressionMonth3: // month
				date = date.AddDate(0, multiply*int(number), 0)
			case ExpressionDay1, ExpressionDay2: // day
				date = date.AddDate(0, 0, multiply*int(number))
			default:
				// do nothing
			}

			year, month, day = date.Year(), int(date.Month()), date.Day()
			if ifEmptyFillAsToday {
				year, month, _ = fillEmptyYearMonthDay(year, month, day)
			}

			debugPrint("dateRelRe1: extracted ymd = %04d-%02d-%02d", year, month, day)

			// append extracted date
			dates[match] = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, _location)
		}
	}
	matches = dateRelRe2.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := dateRelRe2.FindStringSubmatch(match)

			debugPrint("dateRelRe2: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			match := slices[0] // take the first slice

			date := time.Now() // today

			switch match {
			case ExpressionTheDayBeforeYesterday1, ExpressionTheDayBeforeYesterday2: // 2 days before
				date = date.AddDate(0, 0, -2)
			case ExpressionYesterday1, ExpressionYesterday2: // 1 day before
				date = date.AddDate(0, 0, -1)
			case ExpressionToday1, ExpressionToday2: // today
				// do nothing (= today)
			case ExpressionTomorrow1, ExpressionTomorrow2: // 1 day after
				date = date.AddDate(0, 0, 1)
			case ExpressionTheDayAfterTomorrow1: // 2 days after
				date = date.AddDate(0, 0, 2)
			case ExpressionTwoDaysAfterTomorrow1: // 3 days after
				date = date.AddDate(0, 0, 3)
			default:
				// do nothing
			}

			year, month, day = date.Year(), int(date.Month()), date.Day()
			if ifEmptyFillAsToday {
				year, month, _ = fillEmptyYearMonthDay(year, month, day)
			}

			debugPrint("dateRelRe2: extracted ymd = %04d-%02d-%02d", year, month, day)

			// append extracted date
			dates[match] = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, _location)
		}
	}
	matches = dateExactRe1.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := dateExactRe1.FindStringSubmatch(match)

			debugPrint("dateExactRe1: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			year64, _ := strconv.ParseInt(slices[2], 10, 16)
			month64, _ := strconv.ParseInt(slices[4], 10, 16)
			day64, _ := strconv.ParseInt(slices[5], 10, 16)
			year, month, day = int(year64), int(month64), int(day64)
			if ifEmptyFillAsToday {
				year, month, _ = fillEmptyYearMonthDay(year, month, day)
			}

			debugPrint("dateExactRe1: extracted ymd = %04d-%02d-%02d", year, month, day)

			// append extracted date
			dates[match] = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, _location)
		}
	}
	matches = dateExactRe2.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := dateExactRe2.FindStringSubmatch(match)

			debugPrint("dateExactRe2: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			year64, _ := strconv.ParseInt(slices[2], 10, 16)
			month64, _ := strconv.ParseInt(slices[4], 10, 16)
			day64, _ := strconv.ParseInt(slices[5], 10, 16)
			year, month, day = int(year64), int(month64), int(day64)
			if ifEmptyFillAsToday {
				year, month, _ = fillEmptyYearMonthDay(year, month, day)
			}

			debugPrint("dateExactRe2: extracted ymd = %04d-%02d-%02d", year, month, day)

			// append extracted date
			dates[match] = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, _location)
		}
	}

	if len(dates) <= 0 {
		return nil, fmt.Errorf("해당하는 날짜 표현이 없습니다: '%s'", str)
	}

	return dates, nil
}

// ExtractDate extracts date from given string
//
// 주어진 한글 string으로부터 패턴에 가장 먼저 맞는 날짜값 추출
func ExtractDate(str string, ifEmptyFillAsToday bool) (date time.Time, err error) {
	var dates map[string]time.Time
	dates, err = ExtractDates(str, ifEmptyFillAsToday)

	if err != nil {
		return time.Time{}, err
	}

	// get the left-most(with the least index) matched date
	var index int
	minIndex := len(str)
	for s, d := range dates {
		index = strings.Index(str, s)
		if index < minIndex {
			minIndex = index
			date = d
		}
	}

	return date, nil
}

// ExtractTimes extracts all times from given string
//
// returns `nil` times on error
//
// priority of regexs is:
//   timeRelRe1 > timeExactRe1 > timeExactRe2
//
// 주어진 한글 string으로부터 시간 추출
func ExtractTimes(str string, ifEmptyFillAsNow bool) (hmss map[string]Hms, err error) {
	// initialize values
	hmss = map[string]Hms{}
	var parseError error

	// indices of processed matches: not to extract duplicated matches
	alreadyProcessed := map[int]struct{}{}

	var matches []string

	// relative time
	matches = timeRelRe1.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := timeRelRe1.FindStringSubmatch(match)

			debugPrint("timeRelRe1: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			now := time.Now() // now

			var number int64
			if number, parseError = strconv.ParseInt(slices[1], 10, 16); parseError != nil {
				continue
			}
			multiply := 1
			switch slices[3] {
			case ExpressionBefore1: // before
				multiply = -1
			case ExpressionAfter1, ExpressionAfter2: // after
				// do nothing (+1)
			}

			var when time.Time

			switch slices[2] {
			case ExpressionTimeHour1: // hour
				when = now.Add(time.Duration(multiply) * time.Duration(number) * time.Hour)
			case ExpressionTimeMinute1: // minute
				when = now.Add(time.Duration(multiply) * time.Duration(number) * time.Minute)
			case ExpressionTimeSecond1: // second
				when = now.Add(time.Duration(multiply) * time.Duration(number) * time.Second)
			}

			debugPrint("timeRelRe1: extracted hms = %02d:%02d:%02d", when.Hour(), when.Minute(), when.Second())

			// append extracted time
			hmss[match] = Hms{Hours: when.Hour(), Minutes: when.Minute(), Seconds: when.Second(), NumDaysChanged: when.Day() - now.Day(), Ambiguous: false}
		}
	}

	// exact time (pattern 1)
	matches = timeExactRe1.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := timeExactRe1.FindStringSubmatch(match)

			debugPrint("timeExactRe1: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			var hour64 int64
			now := time.Now()
			if hour64, parseError = strconv.ParseInt(slices[3], 10, 16); parseError != nil && ifEmptyFillAsNow {
				hour64 = int64(now.Hour())
			}

			ambiguous := false
			ampm := slices[1]
			if strings.EqualFold(ampm, ExpressionPeriodPM1) || strings.EqualFold(ampm, ExpressionPeriodPM2) {
				if hour64 <= 12 {
					hour64 += 12
				}
			} else if !strings.EqualFold(ampm, ExpressionPeriodAM1) && !strings.EqualFold(ampm, ExpressionPeriodAM2) {
				if hour64 < 12 {
					ambiguous = true
				}
			}

			debugPrint("timeExactRe1: extracted hms = %02d:%02d:%02d", hour64, 30, 0)

			// append extracted time
			hmss[match] = Hms{Hours: int(hour64), Minutes: 30, Seconds: 0, NumDaysChanged: 0, Ambiguous: ambiguous}
		}
	}

	// exact time (pattern 2)
	matches = timeExactRe2.FindAllString(str, -1)
	if matches != nil {
		for _, match := range matches {
			// skip already processed string
			index := strings.Index(str, match)
			if _, exists := alreadyProcessed[index]; exists {
				continue
			}
			alreadyProcessed[index] = struct{}{} // mark it as 'already processed'

			slices := timeExactRe2.FindStringSubmatch(match)

			debugPrint("timeExactRe2: matched string = '%s', slices = [%s]", match, strings.Join(slices, ", "))

			var hour64, minute64, second64 int64 = 0, 0, 0
			now := time.Now()
			if hour64, parseError = strconv.ParseInt(slices[3], 10, 16); parseError != nil && ifEmptyFillAsNow {
				hour64 = int64(now.Hour())
			}
			if minute64, parseError = strconv.ParseInt(slices[5], 10, 16); parseError != nil && ifEmptyFillAsNow {
				minute64 = int64(now.Minute())
			}
			if second64, parseError = strconv.ParseInt(slices[7], 10, 16); parseError != nil && ifEmptyFillAsNow {
				second64 = int64(now.Second())
			}

			ambiguous := false
			ampm := slices[1]
			if strings.EqualFold(ampm, ExpressionPeriodPM1) || strings.EqualFold(ampm, ExpressionPeriodPM2) {
				if hour64 <= 12 {
					hour64 += 12
				}
			} else if !strings.EqualFold(ampm, ExpressionPeriodAM1) && !strings.EqualFold(ampm, ExpressionPeriodAM2) {
				if hour64 < 12 {
					ambiguous = true
				}
			}

			debugPrint("timeExactRe2: extracted hms = %02d:%02d:%02d", hour64, minute64, second64)

			// append extracted time
			hmss[match] = Hms{Hours: int(hour64), Minutes: int(minute64), Seconds: int(second64), NumDaysChanged: 0, Ambiguous: ambiguous}
		}
	}

	if len(hmss) <= 0 {
		return nil, fmt.Errorf("해당하는 시간 패턴이 없습니다: %s", str)
	}

	return hmss, nil
}

// ExtractTime extracts time from given string
//
// 주어진 한글 string으로부터 패턴에 가장 '먼저' 맞는 시간값 추출
func ExtractTime(str string, ifEmptyFillAsNow bool) (hms Hms, err error) {
	var times map[string]Hms
	times, err = ExtractTimes(str, ifEmptyFillAsNow)

	if err != nil {
		return Hms{}, err
	}

	// get the left-most(with the least index) matched time
	var index int
	minIndex := len(str)
	for s, d := range times {
		index = strings.Index(str, s)
		if index < minIndex {
			minIndex = index
			hms = d
		}
	}

	return hms, nil
}

// 주어진 연/월/일이 0  이하일 경우 '오늘' 날짜 기준으로 값을 채워줌
func fillEmptyYearMonthDay(year, month, day int) (int, int, int) {
	today := time.Now()

	if year <= 0 {
		year = int(today.Year())
	}
	if month <= 0 {
		month = int(today.Month())
	}
	if day <= 0 {
		day = int(today.Day())
	}

	return year, month, day
}

// print debug messages
func debugPrint(format string, v ...interface{}) {
	if Verbose {
		log.Printf(format, v...)
	}
}
