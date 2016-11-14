package lkdp

import (
	"testing"
)

func TestExtractDate(t *testing.T) {
	for _, str := range []string{
		`광복절은 1945-08-15`,
		`광주민주화항쟁: 1980년 05월 18일`,
		`기억하라, 기억하라, 11월 5일을.`,
		`방영 시간: 오늘 오후 9시 30분부터`,
		`모레까지 이 과제를 끝마치도록.`,
	} {
		if _, err := ExtractDate(str, false); err != nil {
			t.Error("ExtractDate failed with: " + str)
		}
	}
}

func TestExtractTime(t *testing.T) {
	for _, str := range []string{
		`5시 01분 <= 앞에 0이 padding 되어 있거나 말거나 잘 나와야 합니다`,
		`어제 오후 01시 01분 <= 오후일 때에는 12시간을 더해줘야 합니다`,
		`PM 03:30 <= '오전/오후' 말고 'am/pm'도 잘 구분해야 합니다`,
	} {
		if _, _, _, err := ExtractTime(str, false); err != nil {
			t.Error("ExtractTime failed with: " + str)
		}
	}
}

func BenchmarkExtractDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExtractDate(`아무도 알고 싶어하진 않지만, 내 생일은 1981년 06월 02일이다.`, true)
	}
}

func BenchmarkExtractTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExtractDate(`게임은 오후 7시 30분부터 시작했다. 아직 아내에게 혼날 정도로 오래 하진 않았다고 생각한다.`, true)
	}
}