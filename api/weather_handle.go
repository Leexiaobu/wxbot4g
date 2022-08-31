package api

import (
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"io"
	"net/http"
	"strings"
	"wxbot4g/config"
	"wxbot4g/db"
	"wxbot4g/logger"
)

type WeatherHandle struct {
	// 合成复用Next
	Next
}

func init() {
	//// 创建一个新的定时任务管理器
	//c := cron.New()
	//// 添加一个每三十分钟执行一次的执行器
	//
	//_ = c.AddFunc("0 */30 * * * ? ", updateWeather)
	//_ = c.AddFunc("0 */30 * * * ? ", updateAir)
	//// 新启一个协程，运行定时任务
	//go c.Start()
	//// 等待停止信号结束任务
	//defer c.Stop()
}

var weatherConfig = config.Config.ApiConfig.WeatherConfig

// Do 校验参数的逻辑
func (h *WeatherHandle) Do(ctx *openwechat.MessageContext) (bool, error) {
	if strings.EqualFold(ctx.Content, "#tq") {
		locations := weatherConfig.Location
		location := strings.Split(locations, "#")
		keys, _ := db.RedisClient.GetKeys("weather:*")
		airs, _ := db.RedisClient.GetKeys("air:*")
		if !(len(keys) == len(locations) && len(airs) == len(location)) {
			logger.Log.Info("更新天气信息")
			updateWeather()
			updateAir()
		}
		var result string
		var count int
		for i := 0; i < len(location); i++ {
			split := strings.Split(location[i], ",")
			data, _ := db.RedisClient.GetData("weather:" + split[0])
			air, _ := db.RedisClient.GetData("air:" + split[0])
			result = result + data + air + "\n"
			count++
			if count > 2 {
				ctx.ReplyText(result)
				result = ""
				count = 0
			}
		}
		if !strings.EqualFold("", result) {
			ctx.ReplyText(result)
		}
		return true, nil
	}
	return false, nil
}

func updateWeather() {
	locations := weatherConfig.Location
	location := strings.Split(locations, "#")
	for i := 0; i < len(location); i++ {
		split := strings.Split(location[i], ",")
		url := fmt.Sprintf(weatherConfig.Url, "weather", split[1], weatherConfig.Key)
		response, _ := http.Get(url)
		body, error := io.ReadAll(response.Body)
		var weatherInfo weather
		json.Unmarshal(body, &weatherInfo)
		if error != nil {
			logger.Log.Error("获取%s天气失败,%s", split[1], error.Error())
		}
		db.RedisClient.SetWithTimeout("weather:"+split[0], "【"+split[0]+"】"+weatherInfo.string(weatherInfo), "3600")
		response.Body.Close()
	}

}
func (*weather) string(weatherInfo weather) string {
	now := weatherInfo.Now
	result := fmt.Sprintf(":"+"当前%s度,体感%s度,%s,风力%s级%skm/h,湿度%s,气压%s百帕,可见度%s公里,",
		now.Temp, now.FeelsLike, now.Text,
		now.WindScale, now.WindSpeed,
		now.Humidity, now.Pressure, now.Vis)
	return result
}
func updateAir() {
	locations := weatherConfig.Location
	location := strings.Split(locations, "#")
	for i := 0; i < len(location); i++ {
		split := strings.Split(location[i], ",")
		airUrl := fmt.Sprintf(weatherConfig.Url, "air", split[1], weatherConfig.Key)
		resp, error := http.Get(airUrl)
		body, error := io.ReadAll(resp.Body)
		if error != nil {
			logger.Log.Error("获取%s天气质量失败,%s", split[1], error.Error())
		}
		var airInfo air
		json.Unmarshal(body, &airInfo)
		db.RedisClient.SetWithTimeout("air:"+split[0], airInfo.string(airInfo), "3600")
		resp.Body.Close()
	}
}
func (*air) string(airInfo air) string {
	now := airInfo.Now
	if strings.EqualFold("1", now.Level) {
		return fmt.Sprintf("空气质量%s", now.Aqi)
	}
	sprintf := fmt.Sprintf("空气质量%s,%s级,%s,主要污染物为%s,pm2.5 %sμg/m3。", now.Aqi, now.Level, now.Category,
		now.Primary, now.Pm2P5)
	result := strings.Replace(sprintf, "为NA", "无", -1)
	return result
}

type weather struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Wind360   string `json:"wind360"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		WindSpeed string `json:"windSpeed"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Pressure  string `json:"pressure"`
		Vis       string `json:"vis"`
		Cloud     string `json:"cloud"`
		Dew       string `json:"dew"`
	} `json:"now"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}
type air struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		PubTime  string `json:"pubTime"`
		Aqi      string `json:"aqi"`
		Level    string `json:"level"`
		Category string `json:"category"`
		Primary  string `json:"primary"`
		Pm10     string `json:"pm10"`
		Pm2P5    string `json:"pm2p5"`
		No2      string `json:"no2"`
		So2      string `json:"so2"`
		Co       string `json:"co"`
		O3       string `json:"o3"`
	} `json:"now"`
	Station []struct {
		PubTime  string `json:"pubTime"`
		Name     string `json:"name"`
		Id       string `json:"id"`
		Aqi      string `json:"aqi"`
		Level    string `json:"level"`
		Category string `json:"category"`
		Primary  string `json:"primary"`
		Pm10     string `json:"pm10"`
		Pm2P5    string `json:"pm2p5"`
		No2      string `json:"no2"`
		So2      string `json:"so2"`
		Co       string `json:"co"`
		O3       string `json:"o3"`
	} `json:"station"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}
