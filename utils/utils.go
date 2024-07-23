package utils

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

func GetUuid(c context.Context) string {
	val := c.Value("uuid")
	return fmt.Sprintf("%s", val)
}

func GetUUID() string {
	return uuid.NewString()
}

func GetKDAStr(k int64, d int64, a int64) string {
	var result float64
	if d <= 0 {
		result = float64(k) + float64(a)
	} else {
		result = (float64(k) + float64(a)) / float64(d)
	}
	res := fmt.Sprintf("%.2f", result)
	return res
}

func GetKDA(k int64, d int64, a int64, f int) float64 {
	var result float64
	if d <= 0 {
		result = float64(k) + float64(a)
	} else {
		result = (float64(k) + float64(a)) / float64(d)
	}
	p := math.Pow10(f)
	res := math.Round(result*p) / p
	return res
}

func GetKASTStr(k int64, d int64) string {
	var result float64
	if d <= 0 {
		result = float64(0)
	} else if k > d {
		result = float64(100)
	} else {
		result = float64(k) / float64(d)
	}
	res := fmt.Sprintf("%.1f", result)
	return res
}

func GetPoint(val float64, f int) float64 {
	p := math.Pow10(f)
	res := math.Round(val*p) / p
	return res
}

func GetFloat64(k int64, d int64) float64 {
	if d <= 0 {
		return 0
	} else {
		return float64(k) / float64(d)
	}
}

func GetInt64(k int64, d int64) int64 {
	if d <= 0 {
		return 0
	} else {
		return k / d
	}
}

func GetNowTimestamp() string {
	nowTime := time.Unix(time.Now().Unix(), 0)
	layOut := "2006-01-02 15:04:05"
	date := nowTime.Format(layOut)
	return date
}

func GetTimestampByTime(val int64) string {
	nowTime := time.Unix(val, 0)
	layOut := "2006-01-02 15:04:05"
	date := nowTime.Format(layOut)
	return date
}

func GetTimestamp(val time.Time) string {
	layOut := "2006-01-02 15:04:05"
	date := val.Format(layOut)
	return date
}

func GetDateStr(time time.Time) string {
	layOut := "2006-01-02"
	date := time.Format(layOut)
	return date
}

func GetDayStart(val time.Time) int64 {
	year, month, day := val.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, val.Location())
	return midnight.Unix()
}

func GetWeekStartEnd() (startDate time.Time, endDate time.Time) {
	now := time.Now()
	weekday := int(now.Weekday())
	location := now.Location()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	num := weekday
	if weekday == 0 {
		num = 7
	}
	startDate = nowDate.AddDate(0, 0, -(num - 1))
	endDate = startDate.AddDate(0, 0, 7)
	return startDate, endDate
}

func GetTimeDate(val string) time.Time {
	layOut := "2006-01-02 15:04:05"
	date, _ := time.Parse(layOut, val)
	return date.Add(-8 * time.Hour)
}

func GetNewTimeDate(val string) time.Time {
	layOut := "2006-01-02 15:04:05"
	date, _ := time.Parse(layOut, val)
	return date
}

// GetTimeByTimeZone 返回带东八区时区的time
func GetTimeByTimeZone(val string) time.Time {
	layOut := "2006-01-02 15:04:05"
	date, _ := time.Parse(layOut, val)
	shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
	return date.In(shanghaiLocation).Add(-8 * time.Hour)
}

func GetYearWeek() int {
	now := time.Now()
	year, week := now.ISOWeek()
	weeks, _ := strconv.Atoi(fmt.Sprintf("%d%d", year, week))
	return weeks
}

func GetYearMonth() int {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	months, _ := strconv.Atoi(fmt.Sprintf("%d%d", year, month))
	return months
}

func GetDateLine(val time.Time) int64 {
	date := val.Unix()
	return date
}

func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func IsDev() bool {
	appEnv := os.Getenv("APP_ENV")
	return len(appEnv) == 0 || appEnv == "develop"
}

func IsTest() bool {
	appEnv := os.Getenv("APP_ENV")
	return appEnv == "test"
}

func IsProd() bool {
	appEnv := os.Getenv("APP_ENV")
	return appEnv == "prod"
}

func IsQQ(qq string) bool {
	// 匹配QQ号的正则表达式
	pattern := `^[0-9A-Za-z]{5,12}$`
	match, _ := regexp.MatchString(pattern, qq)
	return match
}

func IsIDCardValid(id string) bool {
	pattern := `^\d{17}(\d|X)$`
	match, _ := regexp.MatchString(pattern, id)
	return match
}

func RandomString(length int) string {
	rand.NewSource(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return strings.ToUpper(string(result))
}

func IsPhoneNumber(phone string) bool {
	pattern := `^1[3456789]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

func IsEmail(email string) bool {
	pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func GetStrLen(str string) (all, a, b, c, d int) {
	a = 0
	b = 0
	c = 0
	d = 0
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			a++
		} else if unicode.IsLetter(v) {
			b++
		} else if unicode.IsNumber(v) {
			c++
		} else {
			d++
		}
	}
	all = a + b + c + d
	return all, a, b, c, d
}

func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Base64Decode(s string) []byte {
	if strings.Index(s, ",") != -1 {
		s = strings.Split(s, ",")[1]
	}
	str, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logx.Errorf("utils Base64Decode fail:%s", err.Error())
		return nil
	}
	return str
}

func SaltMD5(s, salt string) string {
	return MD5(MD5(s) + salt)
}

func GenerateUniqueIdentifier() string {
	// 使用当前时间戳和随机数生成唯一标识符
	rand.NewSource(time.Now().UnixNano())
	timestamp := time.Now().Unix()
	randomNumber := rand.Intn(100000)
	uniqueID := fmt.Sprintf("%d%d", timestamp, randomNumber)
	return uniqueID
}

// ConvertOffsetLimit page转换offset
func ConvertOffsetLimit(page, pageSize int) (offset int, limit int) {
	if page, pageSize = Max(page, 1), Max(pageSize, 10); page > 1 {
		offset = (page - 1) * pageSize
	}
	limit = pageSize

	return
}

func GetFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ZipFile(url []string, fileName []string) (*bytes.Buffer, error) {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	// 压缩临时文件到ZIP
	zipBuf := new(bytes.Buffer)
	zw := zip.NewWriter(zipBuf)

	for k, value := range url {
		suffix := "pdf"
		arr := strings.Split(value, ".")
		if len(arr) > 1 {
			suffix = arr[len(arr)-1]
		}
		wg.Add(1)
		// 从远程URL下载文件
		go func(k int, value string) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			resp, err := http.Get(value)
			if err != nil {
				logx.Errorf("utils ZipFile http get err:%s", err.Error())
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					logx.Errorf("utils ZipFile err:%s", err.Error())
				}
			}(resp.Body)

			// 创建一个临时文件来保存下载的文件
			tmpFile, err := os.CreateTemp("", fileName[k]+"-*"+"."+suffix)
			if err != nil {
				logx.Errorf("utils ZipFile CreateTemp err:%s", err.Error())
				return
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					logx.Errorf("utils ZipFile Remove err:%s", err.Error())
				}
			}(tmpFile.Name()) // 清理临时文件
			defer func(tmpFile *os.File) {
				err := tmpFile.Close()
				if err != nil {
					logx.Errorf("utils ZipFile Close err:%s", err.Error())
				}
			}(tmpFile)

			// 将远程文件内容写入临时文件
			_, err = io.Copy(tmpFile, resp.Body)
			if err != nil {
				logx.Errorf("utils ZipFile Copy err:%s", err.Error())
				return
			}

			// 添加文件到ZIP
			fileInfo, err := tmpFile.Stat()
			if err != nil {
				logx.Errorf("utils ZipFile Stat err:%s", err.Error())
				return
			}
			fileHeader, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				logx.Errorf("utils ZipFile FileInfoHeader err:%s", err.Error())
				return
			}
			w, err := zw.CreateHeader(fileHeader)
			if err != nil {
				logx.Errorf("utils ZipFile CreateHeader err:%s", err.Error())
				return
			}

			// 将临时文件内容复制到ZIP条目
			_, err = tmpFile.Seek(0, 0)
			if err != nil {
				logx.Errorf("utils ZipFile Seek err:%s", err.Error())
				return
			}
			// 重置文件读取位置到开始
			_, err = io.Copy(w, tmpFile)
			if err != nil {
				logx.Errorf("utils ZipFile Copy err:%s", err.Error())
				return
			}
		}(k, value)
	}
	wg.Wait()
	// 关闭ZIP写入器，完成ZIP文件的创建
	err := zw.Close()
	if err != nil {
		return nil, err
	}
	return zipBuf, nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func HasIntersection(slice1, slice2 []string) bool {
	// 创建一个map来存储slice1的元素，用于快速查找
	map1 := make(map[string]bool)
	for _, value := range slice1 {
		map1[value] = true
	}

	// 遍历slice2，检查每个元素是否在map1中
	for _, value := range slice2 {
		if _, exists := map1[value]; exists {
			return true // 如果找到交集，返回true
		}
	}

	return false // 如果没有找到交集，返回false
}

// MonthDays 获取当前月份的天数
func MonthDays() int {
	now := time.Now()
	year, month, _ := now.Date()
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}
	lastDayOfMonth := time.Date(year, nextMonth, 0, 0, 0, 0, 0, now.Location())
	daysInMonth := lastDayOfMonth.Day()
	return daysInMonth
}

func ParseDateString(dateStr string) (time.Time, error) {
	layOut := "2006-01-02"
	date, err := time.Parse(layOut, dateStr)
	return date, err
}

func GetScheduleServer(server string) string {
	if server == "港服" {
		return server
	}
	return "国服"
}

func GetLeagueIDs(leagueIds string) []int64 {
	leagueArr := make([]int64, 0)
	if len(leagueIds) == 0 {
		return leagueArr
	}
	arr := strings.Split(leagueIds, ",")
	for _, v := range arr {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		leagueArr = append(leagueArr, int64(atoi))
	}
	return leagueArr
}

func FormatDateString(timestamp int64) (string, time.Time) {
	layOut := "2006-01-02"
	date := time.Unix(timestamp, 0)
	dateStr := date.Format(layOut)
	return dateStr, date
}

func SliceToMap[K comparable, T any](slice []T, keyMapper func(T) K) map[K]T {
	res := make(map[K]T)
	for _, item := range slice {
		key := keyMapper(item)
		value := item
		res[key] = value
	}
	return res
}

func SliceToMapSlice[K comparable, T any](slice []T, keyMapper func(T) K) map[K][]T {
	res := make(map[K][]T)
	for _, item := range slice {
		key := keyMapper(item)
		value := item

		sliceValue, isOk := res[key]
		if isOk {
			sliceValue = append(sliceValue, value)
		} else {
			sliceValue = []T{value}
		}
		res[key] = sliceValue
	}
	return res
}

func FindPageList[T any](query *gorm.DB, callBack []*T, page int, pageSize int, searchKey string, searchValue string, order string) (resp []*T, count int64, err error) {
	offset, limit := ConvertOffsetLimit(page, pageSize)
	if len(searchKey) > 0 && len(searchValue) > 0 {
		searchKeyArr := strings.Split(searchKey, ",")
		searchValueArr := strings.Split(searchValue, ",")
		for k, v := range searchKeyArr {
			query = query.Where(v+"= ? ", searchValueArr[k])
		}
	}
	resp = callBack
	err = query.Offset(offset).Limit(limit).Order(order).Find(&resp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return callBack, 0, nil
		}
		return callBack, 0, err
	}
	if size := len(resp); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return resp, count, nil
	}
	query.Offset(-1).Limit(-1).Count(&count)
	return resp, count, nil
}

func SubStringWithFormat(str string, length int, ellipsis ...bool) string {
	var newLength int
	var newStr string
	var chineseRegex = regexp.MustCompile(`[^\x00-\xff]`)
	strLength := len(chineseRegex.ReplaceAllString(str, "**"))
	for _, char := range str {
		if chineseRegex.MatchString(string(char)) {
			newLength += 2
		} else {
			newLength++
		}
		if newLength > length {
			break
		}
		newStr += string(char)
	}

	if len(ellipsis) > 0 && ellipsis[0] && strLength > length {
		newStr = fmt.Sprintf("%s%s", newStr, "...")
	}

	return newStr
}

// 实现一下 js 的 substring 方法
func SubString(str string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > end {
		start, end = end, start
	}
	if start > len(str) {
		start = len(str)
	}
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

// Ternary 虚拟的三元运算符
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
