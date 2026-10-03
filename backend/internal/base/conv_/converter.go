package conv_

import (
	"handfree-work/octo-backup/internal/base/log_"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

func ConvertList(sourceList interface{}, targetList interface{}, mapping ...map[string]string) {
	if err := gconv.Structs(sourceList, targetList, mapping...); err != nil {
		panic(err)
	}
}

func Convert(source interface{}, target interface{}) {
	if err := gconv.Struct(source, target); err != nil {
		if strings.HasPrefix(err.Error(), "convert params from") && strings.HasSuffix(err.Error(), "to \"map[string]interface{}\" failed") {
			return
		}
		log_.Error("convert error", err)
		panic(err)
	}
}
func ArrayStringToInt(strArr []string) []int64 {
	res := make([]int64, len(strArr))

	for index, val := range strArr {
		res[index], _ = strconv.ParseInt(val, 10, 64)
	}

	return res
}
func ArrayInt64ToString(intArr []int64) []string {
	res := make([]string, len(intArr))

	for index, val := range intArr {
		res[index] = strconv.FormatInt(val, 10)
	}

	return res
}

func Int64ToStr(val int64) string {
	return strconv.FormatInt(val, 10)
}

func StrToInt64(val string) int64 {
	res, _ := strconv.ParseInt(val, 10, 64)
	return res
}
