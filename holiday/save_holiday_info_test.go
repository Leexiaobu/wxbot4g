package holiday

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestUpdateYear(t *testing.T) {
	// 初始化WechatBotMap
	UpdateYear("2022")
}

func TestSomething(t *testing.T) {
	format := time.Now().Format("20060102")
	fmt.Println(format)
	atoi, _ := strconv.Atoi(format)
	fmt.Println(atoi)
}

func Test_updateEnum(t *testing.T) {
	updateEnum()
}
