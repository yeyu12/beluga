package helpers

import (
	"github.com/astaxie/beego"
	"sort"
	"crypto/md5"
	"fmt"
	"crypto/sha1"
	"time"
	"math/rand"
	"strconv"
)

// 表前缀
func GetTablePrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "beluga_")
}

// 字符串存在数组中
func In_array(str string, arr []string) bool {
	if str == "" {
		return false
	}

	if len(arr) == 0 {
		return false
	}

	sort.Strings(arr)

	for _, value := range arr {
		if str == value {
			return true
		}
	}

	return false
}

// 密码加密
func EncryptionPasswd(str string) string {
	sha1Has := sha1.Sum([]byte(str))
	sha1_str := fmt.Sprintf("%x", sha1Has)
	md5Has := md5.Sum([]byte(sha1_str)[10:])
	md5_str := fmt.Sprintf("%x", md5Has)

	return md5_str
}

// 参数判断

// 生成随机数
func GenerateRangeNum(min int, max int) int {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

// 字符串转md5
func StrToMd5(str string) string {
	str_md5_has := md5.Sum([]byte(str))
	str_md5 := fmt.Sprintf("%x", str_md5_has)

	return str_md5
}

// 生成全局唯一32位字符串
func RandStr() string {
	nano := time.Now().UnixNano()
	nano_str := strconv.FormatInt(nano, 10)
	nano_md5_str := StrToMd5(nano_str)

	rand_num := GenerateRangeNum(1, 1000)
	rand_num_str := strconv.Itoa(rand_num)
	rand_md5 := StrToMd5(rand_num_str)

	rand_str := nano_md5_str + rand_md5

	return StrToMd5(rand_str)
}

// 生成版本号
func GenVersion() string {
	return time.Now().Format("20060102150405") + strconv.Itoa(GenerateRangeNum(10000, 99999))
}
