package holiday

import (
	"fmt"
	"testing"
)

func Test_genTodayInfo(t *testing.T) {
	gotData := genTodayInfo(20220911)
	fmt.Println(gotData)
}

func Test_getToDay(t *testing.T) {
	day, _ := GetToDay()
	fmt.Println(day)
}
func Test_getDay(t *testing.T) {
	fmt.Println(GetDay("20220910"))
	fmt.Println(GetDay("20220911"))
	fmt.Println(GetDay("20220912"))
}
func Test_nil(t *testing.T) {
	var extra string
	if extra == "" {
		fmt.Println("123")
	}
}
