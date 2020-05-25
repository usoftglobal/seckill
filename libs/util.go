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
func TypeOf(v interface{}) {
	fmt.Println("格式：", fmt.Sprintf("%T", v), "值：", v)
}

// 日期时间格式化
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 金额转换（元转分）
func UnitToCents(m uint) uint {
	return m * 100
}

// 金额格式化（分转元）
func PriceFormat(price int) string {
	unit := float64(price) / 100
	return fmt.Sprintf("￥%s", strconv.FormatFloat(unit, 'f', -1, 64))
}

// 字符串转整型
func StringToUint(str string) uint {
	num, _ := strconv.ParseUint(str, 10, 0);
    return uint(num)
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

// Struct 转 Map
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

// JSON 返回成功
func Success(results interface{}) gin.H {
	data := gin.H{"code": 0, "msg": "操作成功", "data": results}
	return data
}

// JSON 返回错误
func Fail(err error) gin.H {
	data := gin.H{"code": 10000, "msg": fmt.Sprintf("%s", err), "data": ""}
	return data
}
