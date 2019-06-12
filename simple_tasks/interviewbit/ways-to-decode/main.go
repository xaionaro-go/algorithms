package main

import (
	"fmt"
)


var m = map[string]int{}
func numDecodings(str string) int {
	if len(str) == 0 {
		return 1
	}
	if str[0] > '9' || str[0] < '0' {
		return 0
	}
	if len(str) == 1 {
		// abcdefghijklmnopqrstuvwxyz
		// 12345678901234567890123456
		return 1
	}
	if str[1] > '9' || str[1] < '0' {
		return 0
	}

	res, ok := m[str]
	if ok {
		return res
	}

	if str[0] == '0' || str[0] >= '3' || (str[0] == '2' && str[1] > '6' || str[1] == '0') {
		res = numDecodings(str[1:])
	} else {
		res = numDecodings(str[1:]) + numDecodings(str[2:])
	}

	m[str] = res
	return res
}

func main() {
	fmt.Println(numDecodings("1"))
	fmt.Println(numDecodings("10"))
	fmt.Println(numDecodings("12"))
	fmt.Println(numDecodings("121266"))
	fmt.Println(numDecodings("5163490394499093221199401898020270545859326357520618953580237168826696965Q537789565062429676962877038781708385575876312877941367557410101383684194057405018861234394660905712238428675120866930196204792703765204322329401298924190"))
	fmt.Println(numDecodings("5163490394499093221199401898020270545859326357520618953580237168826696965537789565062429676962877038781708385575876312877941367557410101383684194057405018861234394660905712238428675120866930196204792703765204322329401298924190"))
}
