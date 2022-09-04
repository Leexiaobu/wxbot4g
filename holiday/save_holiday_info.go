package holiday

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wxbot4g/db"
)

// 数据来源 http://www.apihubs.cn/#/holiday
var url = "https://api.apihubs.cn/holiday/get?month=%s&cn=1&size=31"

func UpdateYear(year string) {
	for i := 2; i <= 12; i++ {
		var month string
		if i < 10 {
			month = fmt.Sprintf("%s0", year) + strconv.Itoa(i)
		} else {
			month = fmt.Sprintf("%s", year) + strconv.Itoa(i)
		}
		fmt.Println(month)
		updateMonInfo(month)
	}
}
func updateMonInfo(mon string) {
	realUrl := fmt.Sprintf(url, mon)
	response, _ := http.Get(realUrl)
	body, _ := io.ReadAll(response.Body)
	var resData HolidayResponse
	json.Unmarshal(body, &resData)
	dataList := resData.Data.List
	//for data := range dataList {
	//	db.MysqlCon.Client.Debug().Create(&data)
	//}
	db.MysqlCon.Client.Debug().Create(&dataList)
}

type Holiday struct {
	ID                int    `json:"id" gorm:"column:id"`
	Year              int    `json:"year" gorm:"column:year"`
	Month             int    `json:"month" gorm:"column:month"`
	Date              int    `json:"date" gorm:"column:date"`
	YearWeek          int    `json:"year_week" gorm:"column:year_week"`
	YearDay           int    `json:"year_day" gorm:"column:year_day"`
	LunarYear         int    `json:"lunar_year" gorm:"column:lunar_year"`
	LunarMonth        int    `json:"lunar_month" gorm:"column:lunar_month"`
	LunarDate         int    `json:"lunar_date" gorm:"column:lunar_date"`
	LunarYearDay      int    `json:"lunar_year_day" gorm:"column:lunar_year_day"`
	Week              int    `json:"week" gorm:"column:week"`
	Weekend           int    `json:"weekend" gorm:"column:weekend"`
	Workday           int    `json:"workday" gorm:"column:workday"`
	Holiday           int    `json:"holiday" gorm:"column:holiday"`
	HolidayOr         int    `json:"holiday_or" gorm:"column:holiday_or"`
	HolidayOvertime   int    `json:"holiday_overtime" gorm:"column:holiday_overtime"`
	HolidayToday      int    `json:"holiday_today" gorm:"column:holiday_today"`
	HolidayLegal      int    `json:"holiday_legal" gorm:"column:holiday_legal"`
	HolidayRecess     int    `json:"holiday_recess" gorm:"column:holiday_recess"`
	YearCn            string `json:"year_cn" gorm:"column:year_cn"`
	MonthCn           string `json:"month_cn" gorm:"column:month_cn"`
	DateCn            string `json:"date_cn" gorm:"column:date_cn"`
	YearWeekCn        string `json:"year_week_cn" gorm:"column:year_week_cn"`
	YearDayCn         string `json:"year_day_cn" gorm:"column:year_day_cn"`
	LunarYearCn       string `json:"lunar_year_cn" gorm:"column:lunar_year_cn"`
	LunarMonthCn      string `json:"lunar_month_cn" gorm:"column:lunar_month_cn"`
	LunarDateCn       string `json:"lunar_date_cn" gorm:"column:lunar_date_cn"`
	LunarYearDayCn    string `json:"lunar_year_day_cn" gorm:"column:lunar_year_day_cn"`
	WeekCn            string `json:"week_cn" gorm:"column:week_cn"`
	WeekendCn         string `json:"weekend_cn" gorm:"column:weekend_cn"`
	WorkdayCn         string `json:"workday_cn" gorm:"column:workday_cn"`
	HolidayCn         string `json:"holiday_cn" gorm:"column:holiday_cn"`
	HolidayOrCn       string `json:"holiday_or_cn" gorm:"column:holiday_or_cn"`
	HolidayOvertimeCn string `json:"holiday_overtime_cn" gorm:"column:holiday_overtime_cn"`
	HolidayTodayCn    string `json:"holiday_today_cn" gorm:"column:holiday_today_cn"`
	HolidayLegalCn    string `json:"holiday_legal_cn" gorm:"column:holiday_legal_cn"`
	HolidayRecessCn   string `json:"holiday_recess_cn" gorm:"column:holiday_recess_cn"`
}

func (m *Holiday) TableName() string {
	return "holiday"
}

type HolidayResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List  []Holiday `json:"list"`
		Page  int       `json:"page"`
		Size  int       `json:"size"`
		Total int       `json:"total"`
	} `json:"data"`
}
