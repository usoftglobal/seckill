package libs

import (
	"fmt"
	"time"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"math/rand"
	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v3"
)

// 读取配置文件
func Conf() *Config {
	conf := new(Config)
	confFile, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(confFile, conf)
	return conf
}

// 类型判断
func TypeOf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// 日期时间格式化
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 金额转换（元转分）
func UnitToCents(unit string) int {
	flt64,_ := strconv.ParseFloat(unit, 64)
	return int(flt64 * 100)
}

// 金额格式化（分转元）
func PriceFormat(price int) string {
	unit := float64(price) / 100
	return fmt.Sprintf("￥%s", strconv.FormatFloat(unit, 'f', -1, 64))
}

// 创建订单号
func CreateOrderNo() string {
	t := time.Now().Format("20060102150405")
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1000)
	return fmt.Sprintf("%s%d", t, r)
}

// 创建随机数
func CreateRandNo() int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100000)
	return r
}

// Struct 转 map
func StructToMap(s interface{}) (map[string]interface{}) {
	j, _ := json.Marshal(s)
	j = []byte(string(j))
	m := map[string]interface{} {}
	json.Unmarshal(j, &m)
	return m
}

// Map 转 JSON
func MapToJSON(m interface{}) string {
	j, _ := json.Marshal(m)
	return string(j)
}

// JSON 转 Map
func JSONToMap(j string) map[string]interface{} {
	temp := []byte(j)
    m := make(map[string]interface{})
	json.Unmarshal(temp, &m)
	return m
}

// JSON 返回成功
func Success(results interface{}) gin.H {
	return gin.H{"code": 0, "msg": "操作成功", "data": results}
}

// JSON 返回错误
func Fail(err error) gin.H {
	return gin.H{"code": 10000, "msg": fmt.Sprintf("%s", err), "data": ""}
}