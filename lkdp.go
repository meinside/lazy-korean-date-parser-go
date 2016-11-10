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
)

var _location *time.Location
var re0, re1, re2, re3 *regexp.Regexp

func init() {
	_location, _ = time.LoadLocation(DefaultLocation)

	re0 = regexp.MustCompile(`((\d{2,})\s*[년年])?\s*((\d{1,2})\s*[월月])?\s*(\d{1,2})\s*[일日]`)
	re1 = regexp.MustCompile(`((\d{2,})\s*[\-\./])?\s*((\d{1,2})\s*[\-\./])?\s*(\d{1,2})`)
	re2 = regexp.MustCompile(fmt.Sprintf(`(%s)`, strings.Join([]string{
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
	re3 = regexp.MustCompile(`(?i)(오전|오후|AM|PM)?\s*((\d{1,2})\s*[시時:])\s*((\d{1,2})(\s*[분分:]?(\d{1,2})\s*[초]?)?)?`)
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

	if re0.Match(bytes) {
		slices := re0.FindStringSubmatch(str)

		year64, _ := strconv.ParseInt(slices[2], 10, 16)
		month64, _ := strconv.ParseInt(slices[4], 10, 16)
		day64, _ := strconv.ParseInt(slices[5], 10, 16)
		year, month, day = int(year64), int(month64), int(day64)
	} else if re1.Match(bytes) {
		slices := re1.FindStringSubmatch(str)

		year64, _ := strconv.ParseInt(slices[2], 10, 16)
		month64, _ := strconv.ParseInt(slices[4], 10, 16)
		day64, _ := strconv.ParseInt(slices[5], 10, 16)
		year, month, day = int(year64), int(month64), int(day64)
	} else if re2.Match(bytes) {
		match := re2.FindStringSubmatch(str)[0]

		date := time.Now() // today

		switch match {
		case ExpressionTheDayBeforeYesterday1, ExpressionTheDayBeforeYesterday2:
			date = date.AddDate(0, 0, -2)
		case ExpressionYesterday1, ExpressionYesterday2:
			date = date.AddDate(0, 0, -1)
		case ExpressionToday1, ExpressionToday2:
			// do nothing (= today)
		case ExpressionTomorrow1, ExpressionTomorrow2:
			date = date.AddDate(0, 0, 1)
		case ExpressionTheDayAfterTomorrow1:
			date = date.AddDate(0, 0, 2)
		case ExpressionTwoDaysAfterTomorrow1:
			date = date.AddDate(0, 0, 3)
		default:
			return date, fmt.Errorf("해당하는 날짜 표현이 없습니다: %s", str)
		}

		year, month, day = date.Year(), int(date.Month()), date.Day()
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

	if re3.Match(bytes) {
		slices := re3.FindStringSubmatch(str)

		var hour64, minute64, second64 int64 = 0, 0, 0
		var parseError error
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
		if strings.EqualFold(ampm, "오후") || strings.EqualFold(ampm, "PM") {
			if hour64 <= 12 {
				hour64 += 12
			}
		}

		hour, min, sec = int(hour64), int(minute64), int(second64)
	} else {
		return hour, min, sec, fmt.Errorf("해당하는 시간 패턴이 없습니다: %s", str)
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
