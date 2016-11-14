# lazy-korean-date-parser-go

이름 그대로:

우리말로 된 string에서 날짜/시간을 '대충' 추출해내는 Go 라이브러리.

## install

```bash
$ go get -u github.com/meinside/lazy-korean-date-parser-go
```

import는:

```go
import (
	lkdp "github.com/meinside/lazy-korean-date-parser-go"
)
```

## usage (example)

### 사용 예:

```go
package main

import (
	"fmt"

	lkdp "github.com/meinside/lazy-korean-date-parser-go"
)

func main() {
	// param2 = false일 경우, year 등 빈 값은 0으로 설정
	if date, err := lkdp.ExtractDate("5월 18일 광주민주화항쟁", false); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted date: %v\n", date)
	}

	// param2 = true인 경우는 현재 시간 기준으로 값을 채워 넣음
	if date, err := lkdp.ExtractDate("5월 18일 광주민주화항쟁 행사 진행 예정", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted date: %v\n", date)
	}

	if date, err := lkdp.ExtractDate("1950年 06월 25일 6.25사변 발발", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted date: %v\n", date)
	}

	// 내일, 모레 등의 keyword의 경우, 기준 시간에 해당 일자만큼 +/- 처리
	if date, err := lkdp.ExtractDate("모레 할 일을 글피로 미루자", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted date: %v\n", date)
	}

	// 1시간 전, 5분 뒤, 30초 후 등의 keyword의 경우, 기준 시간에 해당 시간만큼 +/- 처리
	if h, m, s, err := lkdp.ExtractTime("1시간 뒤에 알려주련?", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}
	if h, m, s, err := lkdp.ExtractTime("기상 5분 전", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}
	if h, m, s, err := lkdp.ExtractTime("30초 후 폭발하도록 되어 있다", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}

	// param2 = false인 경우 빈 값은 0으로 설정
	if h, m, s, err := lkdp.ExtractTime("수업은 오후 1시 30분에 시작합니다", false); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}

	if h, m, s, err := lkdp.ExtractTime("12시에 볼까?", false); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}

	// param2 = true인 경우 빈 값은 현재 시간값으로 채워넣음
	// AM/PM 또는 오전/오후 구분 (오후일 경우 12시간 +)
	if h, m, s, err := lkdp.ExtractTime("지진 발생 시각 PM 07:12 경", true); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}

	if h, m, s, err := lkdp.ExtractTime("현재 시각: 18:09:35.211 KST", false); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Extracted time: %02d:%02d:%02d\n", h, m, s)
	}
}
```

### 출력 예:

```
Extracted date: 0000-05-18 00:00:00 +0827 LMT
Extracted date: 2016-05-18 00:00:00 +0900 KST
Extracted date: 1950-06-25 00:00:00 +0900 KST
Extracted date: 2016-11-16 00:00:00 +0900 KST
Extracted time: 23:53:10
Extracted time: 22:48:10
Extracted time: 22:53:40
Extracted time: 13:30:00
Extracted time: 12:00:00
Extracted time: 19:12:10
Extracted time: 18:09:35
```

## TODO

- [ ] 최적화
- [ ] 패턴 추가

## license

MIT

