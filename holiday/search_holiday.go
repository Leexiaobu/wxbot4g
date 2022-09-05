package holiday

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"wxbot4g/db"
)

var DB = db.MysqlCon.Client

func GetToDay() string {
	today := time.Now().Format("20060102")
	todayInt, error := strconv.Atoi(today)
	var holidayInfo *HolidayInfo
	redisTodayKey := "wechat:holiday:" + today
	db.RedisClient.Del(redisTodayKey)
	data, _ := db.RedisClient.GetData(redisTodayKey)
	if error != nil || data == "" {
		holidayInfo = genTodayInfo(todayInt)
		marshal, _ := json.Marshal(&holidayInfo)
		db.RedisClient.SetWithTimeout(redisTodayKey, string(marshal), "86400")
	} else {
		json.Unmarshal([]byte(data), &holidayInfo)
	}
	msg := getReplyMsg(holidayInfo)
	return msg
}
func GetDay(today string) string {
	todayInt, error := strconv.Atoi(today)
	var holidayInfo *HolidayInfo
	redisTodayKey := "wechat:holiday:" + today
	db.RedisClient.Del(redisTodayKey)
	data, _ := db.RedisClient.GetData(redisTodayKey)
	if error != nil || data == "" {
		holidayInfo = genTodayInfo(todayInt)
		marshal, _ := json.Marshal(&holidayInfo)
		db.RedisClient.SetWithTimeout(redisTodayKey, string(marshal), "86400")
	} else {
		json.Unmarshal([]byte(data), &holidayInfo)
	}
	msg := getReplyMsg(holidayInfo)
	return msg
}
func getReplyMsg(info *HolidayInfo) string {
	today := info.Today
	var msgs []HolidayMessage
	rand.Seed(time.Now().Unix())
	var reply string
	if today.Workday == 1 {
		//日常上班
		DB.Table("holiday_message").Where("type=0").Find(&msgs)
		msg := msgs[rand.Intn(len(msgs))]
		reply = fmt.Sprintf(msg.Message, info.NextHoliday.WorkLength)
		if today.HolidayOvertime != 10 {
			//调休
			var overMsgs []HolidayMessage
			DB.Table("holiday_message").Where("type=1").Find(&overMsgs)
			overMsg := overMsgs[rand.Intn(len(overMsgs))]
			reply = reply + overMsg.Message
		}
	} else {
		//日常周末
		DB.Table("holiday_message").Where("type=2").Find(&msgs)
		msg := msgs[rand.Intn(len(msgs))]
		reply = fmt.Sprintf(msg.Message, info.NextHoliday.HolidaySpend, info.NextHoliday.HolidayLength)
		if today.Holiday != 10 {
			holidayEnum := EnumFiled{}
			DB.Table("holiday_enum").Where("group_name='holiday'").Where("enum_key=?", today.Holiday).Find(&holidayEnum)
			var holidayMsgs []HolidayMessage
			//法定节假日
			DB.Table("holiday_message").Where("type=3").Where("enum_id=?", holidayEnum.ID).Find(&holidayMsgs)
			holidayMsg := holidayMsgs[rand.Intn(len(holidayMsgs))]
			reply = reply + holidayMsg.Message
		}
	}
	return reply
}

func genTodayInfo(todayInt int) (data *HolidayInfo) {
	var today Today
	var nextHoliday NextHoliday
	var holidayLength int
	var holidaySpend int
	var holiday HolidayMin

	DB.Table("holiday").Where("date = (?)", todayInt).Find(&today)
	filed := fmt.Sprintf("( date - %d ) AS work_length", todayInt)
	DB.Table("holiday").Select("date", filed).Where("date >= (?)", todayInt).Where("workday = (?)", 2).Order("date").Limit(1).Find(&nextHoliday)
	nextDate := nextHoliday.Date
	filed = fmt.Sprintf("( date - %d ) AS holiday_length", nextDate)
	DB.Table("holiday").Select(filed).Where("date >= (?)", nextDate).Where("workday = (?)", 1).Order("date").Limit(1).Find(&holidayLength)
	filed = fmt.Sprintf("(  %d -date -1 ) AS holiday_spend", nextDate)
	DB.Table("holiday").Select(filed).Where("date <= (?)", nextDate).Where("workday = (?)", 1).Order("date desc").Limit(1).Find(&holidaySpend)
	DB.Table("holiday").Where("date = (?)", nextDate).Find(&holiday)
	nextHoliday.HolidayLength = holidayLength
	nextHoliday.HolidaySpend = holidaySpend
	info := HolidayInfo{
		Today:       today,
		NextHoliday: nextHoliday,
		Holiday:     holiday,
	}
	return &info
}

type HolidayInfo struct {
	Today       Today
	NextHoliday NextHoliday
	Holiday     HolidayMin
}
type Today struct {
	Date              int `gorm:"date"`
	Week              int
	Workday           int
	Weekend           int
	Holiday           int
	HolidayToday      int    `gorm:"holiday_today"`
	HolidayOvertime   int    `gorm:"holiday_overtime"`
	HolidayOvertimeCn string `gorm:"holiday_overtime_cn"`
}
type NextHoliday struct {
	Date          int
	WorkLength    int
	HolidayLength int
	HolidaySpend  int
}

type HolidayMin struct {
	Date         int `gorm:"date"`
	Weekend      int
	HolidayToday int
	WeekendCn    string
	HolidayCn    string
	HolidayOr    int
}
type HolidayMessage struct {
	ID      int    `json:"id" gorm:"column:id"`
	EnumID  int    `json:"enum_id" gorm:"column:enum_id"`
	Type    int    `json:"type" gorm:"column:type"`       // 0:上班 1调休 2周末 3放假
	Message string `json:"message" gorm:"column:message"` //  上班、调休需要配置下次放假，调休则需要三倍工资
}
