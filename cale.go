package main

// "C"

import "C"

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"
)

type WorkingHours struct {
	StartTime  string   `json:"startTime"`
	EndTime    string   `json:"endTime"`
	DaysOfWeek []string `json:"daysOfWeek"`
}

type CalendarEvent struct {
	FreeBusyStatus string `json:"freeBusyStatus"`
	StartTime      int64  `json:"startTime"`
	EndTime        int64  `json:"endTime"`
}

type User struct {
	UserId        string          `json:"userId"`
	WorkingHours  WorkingHours    `json:"workingHours"`
	CalendarEvent []CalendarEvent `json:"calendarEvent"`
}

type Response struct {
	Code int    `json:"code"`
	Data []User `json:"data"`
}

func parseEvents(users []User) []timeInterval {
	var busyIntervals []timeInterval
	for _, user := range users {
		for _, event := range user.CalendarEvent {
			if event.FreeBusyStatus != "Free" {
				busyIntervals = append(busyIntervals, timeInterval{event.StartTime / 1000, event.EndTime / 1000})
			}
		}
	}
	return mergeIntervals(busyIntervals)
}

type timeInterval struct {
	start int64
	end   int64
}

// 合并时间区间
func mergeIntervals(intervals []timeInterval) []timeInterval {
	// func mergeIntervals(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return []timeInterval{}
	}

	// 按起始时间排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	// fmt.Println(intervals)

	var merged = []timeInterval{}
	for _, interval := range intervals {
		// 如果 merged 为空，或者当前区间不与 merged 最后一个区间重叠，则直接加入
		if len(merged) == 0 || merged[len(merged)-1].end < interval.start {
			merged = append(merged, interval)
		} else {
			// 发生重叠，合并区间
			merged[len(merged)-1].end = max(merged[len(merged)-1].end, interval.end)
		}
	}

	// fmt.Println(merged)
	for _, i := range merged {
		fmt.Printf("merge time: %s - %s\n", localTime(i.start), localTime(i.end))
	}

	return merged
}

// 求最大值
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func findAvailableSlots(users []User) timeInterval {
	workingStart, _ := time.Parse("15:04:05", users[0].WorkingHours.StartTime)
	workingEnd, _ := time.Parse("15:04:05", users[0].WorkingHours.EndTime)

	dailyStart := workingStart.Hour()*3600 + workingStart.Minute()*60
	dailyStart = 9*3600 + 30*60 //9:30
	var dailyEnd = int64(workingEnd.Hour()*3600 + workingEnd.Minute()*60)
	dailyEnd = 18*3600 + 30*60 //18：30

	busyIntervals := parseEvents(users)
	loc, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		fmt.Println("load location err,", err.Error())
		loc = time.Local
	}

	t := time.Unix(busyIntervals[0].start, 0).In(loc)

	// 获取当天 00:00:00 的时间
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)

	var availableSlots timeInterval

	fmt.Println(zeroTime)

	prevEnd := int64(zeroTime.Unix() + int64(dailyStart))
	// fmt.Println(prevEnd)

	for _, interval := range busyIntervals {
		if interval.start-prevEnd >= 3600 {
			availableSlots = timeInterval{prevEnd, prevEnd + 3600}
			return availableSlots
		}
		prevEnd = interval.end
	}
	if dailyEnd-prevEnd >= 3600 {
		availableSlots = timeInterval{prevEnd, int64(dailyEnd)}
	}

	return availableSlots
}

func localTime(ts int64) string {

	loc := time.Local

	t := time.Unix(ts, 0).In(loc)

	return t.Format("2006-01-02 15:04:05")
}

func getString(input *C.wchar_t) string {

	length := 0

	ptr := uintptr(unsafe.Pointer(input))
	for *(*uint16)(unsafe.Pointer(ptr + uintptr(length*2))) != 0 {
		length++
	}

	// Step 2: 将 UTF-16 数据读取到 uint16 数组中
	utf16Data := make([]uint16, length)
	for i := 0; i < length; i++ {
		utf16Data[i] = *(*uint16)(unsafe.Pointer(ptr + uintptr(i*2)))
	}

	// Step 3: 将 UTF-16 转换为 Go 字符串（UTF-8）
	goStr := string(utf16.Decode(utf16Data))

	return goStr
}

func main() {
	// emails := []string{"songkun2@lenovo.com", "wangsy91@lenovo.com", "liulh15@lenovo.com"} //"shijun7@lenovo.com",
	emails := []string{"yuhai@lenovo.com"}

	jsonData, ret := getSchedule(emails, "2025-03-10")

	if ret != 0 {
		fmt.Print("req error ", ret)
		return
	}

	var response Response
	json.Unmarshal([]byte(jsonData), &response)

	ignorelist := []string{
		"songkun2@lenovo.com",
		"wangsy91@lenovo.com",
		"shijun7@lenovo.com",
	}

	var availableSlots timeInterval

	userData := response.Data

	fmt.Println("--------------")
	fmt.Println(userData)
	fmt.Println("==============")

	for {

		if len(userData) == 0 {
			fmt.Println("len(userData)==0")
			break
		}

		fmt.Println("-- start find --")

		availableSlots = findAvailableSlots(userData)

		fmt.Println("-- availableSlots: ", availableSlots)

		if availableSlots.start != 0 {
			break
		} else {
			found := false

			for _, ig := range ignorelist {

				for i, user := range userData {
					if user.UserId == ig {
						found = true
						userData = append(userData[:i], userData[i+1:]...)
						break
					}
				}

				if found {
					break
				}
			}

			if !found {
				break
			}
		}
	}

	fmt.Printf("Available meeting time: %s - %s\n", localTime(availableSlots.start), localTime(availableSlots.end))

	// createEvent(emails, availableSlots.start*1000, availableSlots.end*1000, "[AINow] test")
}

func getSchedule(emails []string, dateStr string) (string, int) {

	url := "https://dev.cochat.lenovo.com/calendar/api/calendar/getschedule"
	method := "POST"

	// 指定时区（本地时区）
	loc, _ := time.LoadLocation("Local")

	// 解析日期，设定时间为当天的 00:00:00
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		fmt.Println("解析时间出错:", err)
		return "", -1
	}

	if len(emails) < 1 {
		return "", -2
	}

	emailStr := strings.Join(emails, ";")

	// 转换为时间戳
	startTime := t.Unix()

	params := fmt.Sprintf(`{
		"Emails":"%s",
		"Start":%d,
		"End":%d,
		"Token":"EpbF0q/M3haVMwUV0ca/IwmbZ7PMjhMVcF4TTObheOVSFrR0XObIyCJAhRjlwZdpJ1C6xC8izIXml2HuUgUdJ0/HkGjO3sz1XV8Ca4p/xi0ThnL/M9yNALN2MCMnGGJO1nlGMrraoarONSV8GQ++AUKANsgVYR8pbDndyTnF/H/5ziHtBP9jwzcOxuSNlYV56pmK0YhjYShYsZpP4uPRPsm+kcmRWT3NlWnqBAnuSFvDnO6WWG/2deTr3VD3GOToaRrhKWSKHpcA+k0JcnGmCTFMPRFk+Qhs0z+cPR3Qnyxn1xYkndmgKfVntkiYgH0fIUbWpCxCRUz0940xKYy4MQ==myhubeyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJvcyI6bnVsbCwidXNlcl9uYW1lIjoic29uZ2t1bjIiLCJzY29wZSI6WyJhbGwiXSwiZXh0cmEiOiIlN0IlMjJleF9lbk5hbWUlMjIlM0ElMjJLdW4rS3VuMitTb25nJTIyJTJDJTIyZXhfbWFpbCUyMiUzQSUyMnNvbmdrdW4yJTQwbGVub3ZvLmNvbSUyMiUyQyUyMmV4X2NuTmFtZSUyMiUzQSUyMiVFNSVBRSU4QiVFNSVBMCU4MyUyMiUyQyUyMmV4X2NvdW50cnklMjIlM0ElMjJDTiUyMiU3RCIsImV4cCI6MTc0OTg3MTI0OSwianRpIjoiYmQzOWRkOTEtNzU1Ni00Y2E3LTkzMGYtNzM1YzNmM2Y5NzhmIiwiY2xpZW50X2lkIjoiY2xpZW50LWFwcCJ9.TiZ8vr8lyyxECobT-co0GzczpXd7jcTl3RPEMgQvM3IjpmXUiRaK4warNwNUU-7pOLbMa8TLhrOmJo-aeh1ZkzSU8_sjwaOc0fbcr-f41J1sONc6Yd5pMQj7Zm2iAcUWmBbEQWXDC76Y7dqCEGPUjQlJYaXLrMuE6kyXViQj7u3DSqRGS_TU6Qn21fhducsKxCOq2g4Y5mJj-Oqr-bUDNpv5jnkE9Rdycgy4b3rxKrVL8hdEeyhxmdOYF82x5acMGHFSo1bOE1F4WoX3GHvNf5i8ln5fUEEXyE6FNEuiW7pwsmg2qB9nnFgLIFlQSTaQGxVTiC43v_Pvq881WhWwUw"}`,
		emailStr, startTime*1000, startTime*1000+24*3600*1000)

	payload := strings.NewReader(params)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", -3
	}
	req.Header.Add("loc", "O365")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", -4
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", -5
	}
	// fmt.Println(string(body))
	return string(body), 0
}

func createEvent(emails []string, startTs int64, endTs int64, topic string) int {
	url := "https://dev.cochat.lenovo.com/calendar/api/calendar/sendevent"
	method := "POST"

	if len(emails) < 1 {
		return -20
	}

	emailStr := strings.Join(emails, ";")

	params := fmt.Sprintf(`{
		"Subject": "%s",
		"Body": "<p><br></p>",
		"Location": "",
		"Start": %d,
		"End": %d,
		"ReminderMinutesBeforeStart": 15,
		"RequiredAttendees": "%s",
		"Token":"EpbF0q/M3haVMwUV0ca/IwmbZ7PMjhMVcF4TTObheOVSFrR0XObIyCJAhRjlwZdpJ1C6xC8izIXml2HuUgUdJ0/HkGjO3sz1XV8Ca4p/xi0ThnL/M9yNALN2MCMnGGJO1nlGMrraoarONSV8GQ++AUKANsgVYR8pbDndyTnF/H/5ziHtBP9jwzcOxuSNlYV56pmK0YhjYShYsZpP4uPRPsm+kcmRWT3NlWnqBAnuSFvDnO6WWG/2deTr3VD3GOToaRrhKWSKHpcA+k0JcnGmCTFMPRFk+Qhs0z+cPR3Qnyxn1xYkndmgKfVntkiYgH0fIUbWpCxCRUz0940xKYy4MQ==myhubeyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJvcyI6bnVsbCwidXNlcl9uYW1lIjoic29uZ2t1bjIiLCJzY29wZSI6WyJhbGwiXSwiZXh0cmEiOiIlN0IlMjJleF9lbk5hbWUlMjIlM0ElMjJLdW4rS3VuMitTb25nJTIyJTJDJTIyZXhfbWFpbCUyMiUzQSUyMnNvbmdrdW4yJTQwbGVub3ZvLmNvbSUyMiUyQyUyMmV4X2NuTmFtZSUyMiUzQSUyMiVFNSVBRSU4QiVFNSVBMCU4MyUyMiUyQyUyMmV4X2NvdW50cnklMjIlM0ElMjJDTiUyMiU3RCIsImV4cCI6MTc0OTg3MTI0OSwianRpIjoiYmQzOWRkOTEtNzU1Ni00Y2E3LTkzMGYtNzM1YzNmM2Y5NzhmIiwiY2xpZW50X2lkIjoiY2xpZW50LWFwcCJ9.TiZ8vr8lyyxECobT-co0GzczpXd7jcTl3RPEMgQvM3IjpmXUiRaK4warNwNUU-7pOLbMa8TLhrOmJo-aeh1ZkzSU8_sjwaOc0fbcr-f41J1sONc6Yd5pMQj7Zm2iAcUWmBbEQWXDC76Y7dqCEGPUjQlJYaXLrMuE6kyXViQj7u3DSqRGS_TU6Qn21fhducsKxCOq2g4Y5mJj-Oqr-bUDNpv5jnkE9Rdycgy4b3rxKrVL8hdEeyhxmdOYF82x5acMGHFSo1bOE1F4WoX3GHvNf5i8ln5fUEEXyE6FNEuiW7pwsmg2qB9nnFgLIFlQSTaQGxVTiC43v_Pvq881WhWwUw"}`,
		topic, startTs, endTs, emailStr)

	payload := strings.NewReader(params)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return -21
	}
	req.Header.Add("loc", "O365")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return -22
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return -23
	}
	fmt.Println(string(body))
	return 0
}

func parseData(emails []string, dateStr string) (timeInterval, int) {
	// emails := []string{"songkun2@lenovo.com", "wangsy91@lenovo.com", "shijun7@lenovo.com"}

	jsonData, ret := getSchedule(emails, dateStr) // "2025-03-07")

	if ret != 0 {
		fmt.Print("req error ", ret)
		return timeInterval{}, -30
	}

	var response Response
	json.Unmarshal([]byte(jsonData), &response)

	ignorelist := []string{
		"songkun2@lenovo.com",
		"wangsy91@lenovo.com",
		"shijun7@lenovo.com",
		"liulh15@lenovo.com",
	}

	var availableSlots timeInterval

	userData := response.Data

	for {
		fmt.Println("for--")
		fmt.Println(userData)

		availableSlots = findAvailableSlots(userData)

		fmt.Println("availableSlots: ", availableSlots)

		if availableSlots.start != 0 {
			break
		} else {
			found := false

			for _, ig := range ignorelist {

				for i, user := range userData {
					if user.UserId == ig {
						found = true
						userData = append(userData[:i], userData[i+1:]...)
						break
					}
				}

				if found {
					break
				}
			}

			if !found {
				break
			}
		}
	}

	fmt.Printf("Available meeting time: %s - %s\n", localTime(availableSlots.start), localTime(availableSlots.end))

	if availableSlots.start == 0 {
		return timeInterval{}, -31
	}

	return availableSlots, 0
}

//export AddCale
func AddCale(titleChar, dateChar, attendeeChar *C.wchar_t) int64 {

	title := getString(titleChar)
	date := getString(dateChar)

	attendees := ""

	if attendeeChar != nil {
		attendees = getString(attendeeChar)
	}

	slice := strings.Split(attendees, ";")
	emails := []string{}

	emailtest := map[string]string{
		// zhangwei, zhang wei
		"zhang wei": "zhangweig@lenovo.com",
		"zhangwei":  "zhangweig@lenovo.com",

		"song kun": "songkun2@lenovo.com",
		"songkun":  "songkun2@lenovo.com",

		//shijun, shi jun
		"shi jun": "shijun7@lenovo.com",
		"shijun":  "shijun7@lenovo.com",

		// zheng aiguo, ai guo, AG
		"ai guo":      "zhengag@lenovo.com",
		"zheng aiguo": "zhengag@lenovo.com",
		"ag":          "zhengag@lenovo.com",

		// "cui qi":       "cuiqi@lenovo.com",
		"shuangyang": "wangsy91@lenovo.com",
		// "eric":         "yuhai@lenovo.com",
		"charlie":      "xcharlie@lenovo.com",
		"peng jinzhen": "pengjz1@lenovo.com",
		"pengjinzhen":  "pengjz1@lenovo.com",
		"jinzhen":      "pengjz1@lenovo.com",

		"lihua":     "liulh15@lenovo.com",
		"liu lihua": "liulh15@lenovo.com",
	}

	for _, v := range slice {

		value, exists := emailtest[strings.ToLower(strings.TrimSpace(v))]

		fmt.Println(strings.ToLower(strings.TrimSpace(v)))
		if exists {

			emails = append(emails, value)
		}

	}

	// func parseData(emails []string, dateStr string, topic string) int
	fmt.Println(emails, date, title)

	availableSlots, ret := parseData(emails, date)

	if ret == 0 {
		ret = createEvent(emails, availableSlots.start*1000, availableSlots.end*1000, "[AINow] "+title)
		return availableSlots.start
	} else {
		return int64(ret)
	}
}

// go build -o ./cale.dll -buildmode=c-shared cale.go
