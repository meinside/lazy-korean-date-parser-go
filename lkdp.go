package lkdp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultLocation = "Asia/Seoul"
)

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

	ExpressionPeriodAm1 = `오전`
	ExpressionPeriodAm2 = `AM`
	ExpressionPeriodPm1 = `오후`
	ExpressionPeriodPm2 = `PM`

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
)

var _location *time.Location

var dateExactRe1, dateExactRe2 *regexp.Regexp // 특정 일자
var dateRelRe1, dateRelRe2 *regexp.Regexp     // 상대 일자
var timeRelRe1 *regexp.Regexp                 // 상대 시간
var timeExactRe1 *regexp.Regexp               // 특정 시간

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
	dateExactRe2 = regexp.MustCompile(`((\d{2,})\s*[\-\./])?\s*((\d{1,2})\s*[\-\./]\s*(\d{1,2}))`)
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
	timeExactRe1 = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)?\s*((\d{1,2})\s*[%s])\s*((\d{1,2})(\s*[%s]?(\d{1,2})\s*[%s]?)?)?`,
		strings.Join([]string{
			ExpressionPeriodAm1,
			ExpressionPeriodAm2,
			ExpressionPeriodPm1,
			ExpressionPeriodPm2,
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

// 지역 설정 (timezone)
//
// https://golang.org/pkg/time/#Location
func SetLocation(str string) error {
	var err error = nil

	_location, err = time.LoadLocation(str)

	return err
}

// 주어진 한글 string으로부터 가장 먼저 패턴에 맞는 날짜값 추출
func ExtractDate(str string, ifEmptyFillAsToday bool) (date time.Time, err error) {
	var year, month, day int = 0, 0, 0

	bytes := []byte(str)

	if dateRelRe1.Match(bytes) {
		slices := dateRelRe1.FindStringSubmatch(str)

		date := time.Now() // today

		number, _ := strconv.ParseInt(slices[1], 10, 16)

		var multiply int = 1
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
			return date, fmt.Errorf("해당하는 날짜 표현이 없습니다: %s", str)
		}

		year, month, day = date.Year(), int(date.Month()), date.Day()
	} else if dateRelRe2.Match(bytes) {
		match := dateRelRe2.FindStringSubmatch(str)[0]

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
			return date, fmt.Errorf("해당하는 날짜 표현이 없습니다: %s", str)
		}

		year, month, day = date.Year(), int(date.Month()), date.Day()
	} else if dateExactRe1.Match(bytes) {
		slices := dateExactRe1.FindStringSubmatch(str)

		year64, _ := strconv.ParseInt(slices[2], 10, 16)
		month64, _ := strconv.ParseInt(slices[4], 10, 16)
		day64, _ := strconv.ParseInt(slices[5], 10, 16)
		year, month, day = int(year64), int(month64), int(day64)
	} else if dateExactRe2.Match(bytes) {
		slices := dateExactRe2.FindStringSubmatch(str)

		year64, _ := strconv.ParseInt(slices[2], 10, 16)
		month64, _ := strconv.ParseInt(slices[4], 10, 16)
		day64, _ := strconv.ParseInt(slices[5], 10, 16)
		year, month, day = int(year64), int(month64), int(day64)
	} else {
		return date, fmt.Errorf("해당하는 날짜 패턴이 없습니다: %s", str)
	}

	if ifEmptyFillAsToday {
		year, month, _ = fillEmptyYearMonthDay(year, month, day)
	}

	// set date
	date = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, _location)

	return date, nil
}

// 주어진 한글 string으로부터 시간 추출
func ExtractTime(str string, ifEmptyFillAsNow bool) (hour, min, sec int, err error) {
	bytes := []byte(str)

	var parseError error
	if timeRelRe1.Match(bytes) {
		slices := timeRelRe1.FindStringSubmatch(str)

		when := time.Now() // now

		var number int64 = 0
		if number, parseError = strconv.ParseInt(slices[1], 10, 16); parseError != nil {
			return 0, 0, 0, fmt.Errorf("해당하는 시간 패턴이 없습니다: %s", str)
		}
		var multiply int = 1
		switch slices[3] {
		case ExpressionBefore1: // before
			multiply = -1
		case ExpressionAfter1, ExpressionAfter2: // after
			// do nothing (+1)
		}
		switch slices[2] {
		case ExpressionTimeHour1: // hour
			when = when.Add(time.Duration(multiply) * time.Duration(number) * time.Hour)
		case ExpressionTimeMinute1: // minute
			when = when.Add(time.Duration(multiply) * time.Duration(number) * time.Minute)
		case ExpressionTimeSecond1: // second
			when = when.Add(time.Duration(multiply) * time.Duration(number) * time.Second)
		}

		hour, min, sec = when.Hour(), when.Minute(), when.Second()
	} else if timeExactRe1.Match(bytes) {
		slices := timeExactRe1.FindStringSubmatch(str)

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

		ampm := slices[1]
		if strings.EqualFold(ampm, ExpressionPeriodPm1) || strings.EqualFold(ampm, ExpressionPeriodPm2) {
			if hour64 <= 12 {
				hour64 += 12
			}
		}

		hour, min, sec = int(hour64), int(minute64), int(second64)
	} else {
		return 0, 0, 0, fmt.Errorf("해당하는 시간 패턴이 없습니다: %s", str)
	}

	return hour, min, sec, nil
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
