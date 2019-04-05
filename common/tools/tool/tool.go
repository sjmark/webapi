package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"adwetec.com/spider/lib_common/tools/typeconv"
	"bytes"
	"github.com/klauspost/compress/zlib"
	"io"
	"regexp"
	"html/template"
)

const (
	sidrandom    = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	verifyrandom = "1234567890"
	regular      = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\d{8}$"
	reg_num      = "^\\+?(0|[1-9][0-9]*)$"
	reg_amount   = "^(([0-9]+\\.[0-9]*[1-9][0-9]*)|([0-9]*[1-9][0-9]*\\.[0-9]+)|([0-9]*[1-9][0-9]*))$"
)

//StrToBytes string 转为byte 高效
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//BytesToStr byte 转为string 高效
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func BoolToInt32(b bool) int32 {
	var r int32
	if b {
		r = 1
	}
	return r
}

func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func StrToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func StrToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

//GetIPAddress 获取IP地址
func GetIPAddress(r *http.Request) string {
	ip := ""
	proxyip := ""
	if r.Header.Get("X-Forwarded-For") != "" {
		proxyip = r.Header.Get("X-Forwarded-For")
	} else if r.Header.Get("X-Real-IP") != "" {
		proxyip = r.Header.Get("X-Real-IP")
	} else if r.Header.Get("Host") != "" {
		proxyip = r.Header.Get("Host")
	} else {
		proxyip = r.RemoteAddr
	}
	ips := strings.Split(proxyip, ":")
	if len(ips) > 0 {
		ip = ips[0]
	}
	return ip
}

//MD5f 获取MD5值
func MD5f(str string) string {
	h := md5.New()
	if _, err := h.Write([]byte(str)); err != nil {
		fmt.Println("MD5", err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

//JoinArray 连接数组为一个字符串
//@param sep:连接符
//@param array:
//@return 连接后的字符串
func JoinArray(sep string, array ...interface{}) string {
	strVals := make([]string, 0, len(array))
	for _, val := range array {
		strVals = append(strVals, fmt.Sprint(val))
	}
	return strings.Join(strVals, sep)
}

//MakeKey 根据指定信息数组生成一个字段名称
//@param keys:用于生成字段名称的信息数组
//@return 使用:连接的字符串
func MakeKey(keys ...interface{}) string {
	return JoinArray(":", keys...)
}

//OpenKey 拆解字段名
func OpenKey(key string) []string {
	return strings.Split(key, ":")
}

//JoinKeys 拼装keys
func JoinKeys(keys ...string) string {
	return strings.Join(keys, ":")
}
func JoinStrings(strs []string) string {
	return strings.Join(strs, ":")
}

//JoinKeyID 拼装keyID
func JoinKeyID(majorKey string, id int64) string {
	return JoinKeys(strconv.FormatInt(id, 10), majorKey)
}

//JoinKeyID 拼装keyID
func JoinUserFiled(gsid, uid int64, f string) string {
	return JoinKeys(strconv.FormatInt(gsid, 10), strconv.FormatInt(uid, 10), f)
}

//JoinKeyID 拼装keyID
func JoinUserUidFiled(uid int64, f string) string {
	return JoinKeys(strconv.FormatInt(uid, 10), f)
}

func SubTimes() int64 {
	nowtime := time.Now()
	lastTime := nowtime.Add(time.Hour * 24).Format("2006-01-02")
	var loc, _ = time.LoadLocation("Local")
	var tm, _ = time.ParseInLocation("2006-01-02", lastTime, loc)
	t := tm.Sub(nowtime)
	return int64(t.Seconds())
}

func InitSlice(id, slice int64) int64 {
	return (id % slice) + 1
}

func InitStrSlice(id string, slice int64) int64 {
	idInt64, _ := strconv.ParseInt(id, 10, 64)
	return (idInt64 % slice) + 1
}
func InitUid(defUid int64) int64 {
	var uid int64
	if err := typeconv.SetInt64FromStr(&uid, fmt.Sprintf("%d", defUid)); err != nil {
		fmt.Println("InitUid", err)
	}
	return uid
}

// 用户财富类型
func UserStateType(userState int64) int64 {
	if userState == 0 {
		return 0
	}
	return userState / 10000
}

func RedisServerNum(id int64) int {
	var serverNum int
	if err := typeconv.SetIntFromStr(&serverNum, fmt.Sprintf("%d", id)[0:3]); err != nil {
		fmt.Println("RedisServerNum", err)
	}
	return serverNum
}

func SplitArr(arr []interface{}, index, offset int) []interface{} {
	var size = len(arr)
	if offset >= size {
		return arr[index:size]
	} else {
		return arr[index:offset]
	}
}

// 获取当前时间的0点
func GetZeroPoint() (int64, int64) {

	var tm = time.Now()
	ntm := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())

	return ntm.Unix(), ntm.Add(time.Hour * 24).Unix()
}

// 获取当前时间的0点
func GetZeroPointTime() (time.Time, time.Time) {

	var tm = time.Now()
	ntm := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())

	return ntm, ntm.Add(time.Hour * 24)
}

// 获取时间
func GetStartAndEnd(format, start, end string) (int64, int64) {

	var loc, _ = time.LoadLocation("Local")
	var startTime, endTime int64

	if start != "" {
		var stm, _ = time.ParseInLocation(format, start, loc)
		startTime = stm.Unix()
	}

	if end != "" {
		var etm, _ = time.ParseInLocation(format, end, loc)
		endTime = etm.Unix()
	}

	return startTime, endTime
}

//ToDayZeroPoint 今日凌晨
func ToDayZeroPoint() time.Time {
	var tm = time.Now()
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
}

//TimeSubDay 计算时间相差天数
func TimeSubDay(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}

//IsTimeToDay 是今天
func IsTimeToDay(t1 time.Time) bool {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	now := time.Now()
	t2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	return t1.Sub(t2).Hours() == 0
}

// ----------------------------------------------------------------
func GetOtype(id, typeNum int64) int32 {
	return int32(id / typeNum)
}

// ----------------------------------------------------------------
func GetActType(id, typeNum int64) int32 {
	return int32(id % typeNum)
}

//ToDayZeroPoint 某日凌晨
func ToZeroPoint(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
}

//ToDayZeroPoint 23:59:59
func ToLastZeroPoint(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day()+1, 0, 0, 0, 0, tm.Location()).Add(-time.Second)
}

//GetEndTime 某日凌晨
func GetEndTime(tm string) time.Time {

	if tm != "" {
		var loc, _ = time.LoadLocation("Local")
		stm, _ := time.ParseInLocation("2006-01-02 15:04:05", tm, loc)
		return time.Date(stm.Year(), stm.Month(), stm.Day(), 0, 0, 0, 0, loc)
	}

	return time.Time{}
}

func GetStrTime(tm string) time.Time {

	if tm != "" {
		var loc, _ = time.LoadLocation("Local")
		stm, _ := time.ParseInLocation("2006-01-02 15:04:05", tm, loc)
		return stm
	}

	return time.Time{}
}

//GetEndTime 某日凌晨
func GetStartTime(tm int64) time.Time {
	if tm > 0 {
		var stm = time.Unix(tm, 0)
		return stm
	}

	return time.Time{}
}

// 获取某个时间对应当前月对第几周
func GetCurrentWeek(t time.Time) int {
	_, nowWeek := t.ISOWeek()
	tm := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local)
	_, pm := tm.ISOWeek()
	return nowWeek - pm + 1
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func MysqlInjectCheck(str string) string {
	var newStr string
	newStr = strings.Replace(str, "'", "\\'", -1)
	newStr = strings.Replace(newStr, "--", "\\--", -1)
	// newStr = strings.Replace(str, "select", "", -1)
	// newStr = strings.Replace(newStr, "insert", "", -1)
	// newStr = strings.Replace(newStr, "and", "", -1)
	// newStr = strings.Replace(newStr, "or", "", -1)
	// newStr = strings.Replace(newStr, "drop", "", -1)
	// newStr = strings.Replace(newStr, "delete", "", -1)
	// newStr = strings.Replace(newStr, "load_file", "", -1)
	// newStr = strings.Replace(newStr, "outfile", "", -1)
	return strings.TrimSpace(newStr)
}

//验证手机
func MobileNumCheck(MobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(MobileNum)
}

//验证纯数字
func ValidateNum(num string) (bool, string) {
	nums := strings.Trim(num, "")
	reg := regexp.MustCompile(reg_num)
	return reg.MatchString(nums), nums
}

//验证金额
func ValidateAmount(amount string) (bool, string) {
	amounts := strings.Trim(amount, "")
	reg := regexp.MustCompile(reg_amount)
	return reg.MatchString(amounts), amounts
}

//过滤字符中
func ValidateStr(str string) string {

	//替换HTML的空白字符为空格
	re := regexp.MustCompile(`\s`) //ns*r
	str = re.ReplaceAllString(str, " ")
	//格式化用户输入信息html标签，返回 <script>alert("sss")</script>--->&lt;script&gt;
	str = template.HTMLEscapeString(str)
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	str = re.ReplaceAllString(str, "\n")
	//防sql注入
	re, _ = regexp.Compile("(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)")
	str = re.ReplaceAllString(str, "")

	return strings.TrimSpace(str)
}

//过滤字符中
func ValidateImgStr(str string) bool {
	for _, v := range []string{"png", "jpg", "git", "pneg"} {
		if strings.Contains(str, v) {
			return true
		}
	}

	return false
}
