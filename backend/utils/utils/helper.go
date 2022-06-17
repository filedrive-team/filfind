package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

const (
	GenesisUnixTime = 1598306400
)

func PaginationHelper(page, pageSize, defaultPageSize int) (offset int, size int) {
	size = pageSize
	if size == 0 {
		size = defaultPageSize
	}
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * size
	}
	return
}

func TmpFileName(name string) string {
	return fmt.Sprintf("%s_%d_%s", name, CurrentTimestamp(), GenRandStr(4))
}

func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

func GetFileList(args []string) (fileList []string, err error) {
	fileList = make([]string, 0)
	for _, path := range args {
		finfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if finfo.IsDir() {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return nil, err
			}
			templist := make([]string, 0)
			for _, n := range files {
				templist = append(templist, fmt.Sprintf("%s/%s", path, n.Name()))
			}
			list, err := GetFileList(templist)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, list...)
		} else {
			fileList = append(fileList, path)
		}
	}

	return
}

func GetEpochByTime(t time.Time) int64 {
	epoch := (t.Unix() - GenesisUnixTime) / 30
	if epoch < 0 {
		epoch = 0
	}
	return epoch
}

func GetGenesisTime() time.Time {
	return time.Unix(GenesisUnixTime, 0)
}

func GetBlockTimeByEpoch(epoch int64) time.Time {
	return time.Unix(GenesisUnixTime+epoch*30, 0)
}

func GetDurationByEpoch(epoch int64) time.Duration {
	return time.Duration(epoch*30) * time.Second
}

func MonthBegin(t time.Time) time.Time {
	begin, err := time.ParseInLocation("2006-01-02 15:04:05",
		fmt.Sprintf("%d-%02d-01 00:00:00", t.Year(), t.Month()),
		time.Now().Location())
	if err != nil {
		panic(err)
	}
	return begin
}

func GetClientIP(c *gin.Context) string {
	remoteIPHeaders := []string{
		"X-Real-Ip",
		"X-Forwarded-For",
	}

	for _, headerName := range remoteIPHeaders {
		ip, valid := validateHeader(c.GetHeader(headerName))
		if valid {
			return ip
		}
	}
	return c.ClientIP()
}

func validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",")
	for i, ipStr := range items {
		ipStr = strings.TrimSpace(ipStr)
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return "", false
		}

		// We need to return the first IP in the list, but,
		// we should not early return since we need to validate that
		// the rest of the header is syntactically valid
		if i == 0 {
			clientIP = ipStr
			valid = true
		}
	}
	return
}
