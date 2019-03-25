package lkdp

// $ go test .
// $ go test -bench=.

import (
	"testing"
)

func TestExtractDate(t *testing.T) {
	for str, b := range map[string]bool{
		`광복절은 1945-08-15`:        false,
		`광주민주화항쟁: 1980년 05월 18일`: false,
		`기억하라, 기억하라, 11월 5일을.`:   true,
		`방영 시간: 오늘 오후 9시 30분부터`:  true,
		`모레까지 이 과제를 끝마치도록.`:      false,
		`2일 전만 해도 이런게 없었는데...`:   false,
		`대략 10개월 뒤면 또 크리스마스라네~`:  false,
		`1년 후 2020 원더키디 등장`:      false,
		`12월 12일부터 다음 해 6월 2일까지`: true,
		`6월 2일부터 다시 12월 12일까지`:   true,
	} {
		if d, err := ExtractDate(str, b); err == nil {
			t.Logf("ExtractDate extracted date: %s from string: '%s'", d.Format("2006-01-02"), str)
		} else {
			t.Errorf("ExtractDate failed with string: '%s' (error: %s)", str, err)
		}
	}
}

func TestExtractDates(t *testing.T) {
	for str, b := range map[string]bool{
		`2019년 3월 1일에 3.1 만세운동(1919.03.01) 100주년이라고 알려다오`: false,
		`6일 후면 3월 31일, 이달의 마지막 날이다`:                       false,
	} {
		if ds, err := ExtractDates(str, b); err == nil {
			for m, d := range ds {
				t.Logf("ExtractDates extracted date: %s from match: '%s' in string: '%s'", d.Format("2006-01-02"), m, str)
			}
		} else {
			t.Errorf("ExtractDates failed with string: '%s' (error: %s)", str, err)
		}
	}
}

func TestExtractTime(t *testing.T) {
	for str, b := range map[string]bool{
		`5시 01분`:          false,
		`어제 오후 01시 01분`:   false,
		`PM 03:30`:        false,
		`9시 반에 일어났다`:      false,
		`1시간 전이면 몇 시일까요?`: false,
		`5분 뒤 30분 후에 약먹으라고 알려줄래?`: false,
	} {
		if hms, err := ExtractTime(str, b); err == nil {
			t.Logf("ExtractTime extracted time: %02d:%02d:%02d from string: '%s'", hms.Hours, hms.Minutes, hms.Seconds, str)
		} else {
			t.Errorf("ExtractTime failed with string: '%s' (error: %s)", str, err)
		}
	}
}

func TestExtractTimes(t *testing.T) {
	for str, b := range map[string]bool{
		`5시 01분 ~ 15시 6분`: false,
		`어제 오후 01시 01분과 오늘 오전 05:00 사이에 대체 무슨 일이 있었던걸까?`: false,
		`PM 03:30, 3시간 10분 뒤`:      false,
		`9시 반에 일어났다가 10시에 다시 잠들었다`: false,
	} {
		if hmss, err := ExtractTimes(str, b); err == nil {
			for m, hms := range hmss {
				t.Logf("ExtractTimes extracted time: %02d:%02d:%02d from match: '%s' in string: '%s'", hms.Hours, hms.Minutes, hms.Seconds, m, str)
			}
		} else {
			t.Errorf("ExtractTime failed with string: '%s' (error: %s)", str, err)
		}
	}
}

func BenchmarkExtractDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExtractDate(`아무도 알고 싶어하진 않지만, 내 생일은 1981년 06월 02일이다.`, true)
		_, _ = ExtractDate(`1977년 12월 12일은 누구 생일일까요?`, true)
	}
}

func BenchmarkExtractDates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExtractDates(`내 생일이 1981년 06월 02일이니 2020년 6월 2일은 내 삼십대의 마지막 생일이겠구나.`, true)
		_, _ = ExtractDates(`내일은 2020년 원더키디가 나오는 날보다 앞일까 뒤일까`, true)
	}
}

func BenchmarkExtractTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExtractTime(`게임은 오후 7시 30분부터 시작했다. 아직 아내에게 혼날 정도로 오래 하진 않았다고 생각한다.`, true)
	}
}

func BenchmarkExtractTimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExtractTimes(`3시간 뒤, 18:00에는 23:00:00에 있을 임시점검을 대비한 시뮬레이션이 있을 예정입니다.`, true)
		_, _ = ExtractTimes(`3시 타임에 늦지 않도록 1시간 전까지는 도착해야 한다`, true)
	}
}
