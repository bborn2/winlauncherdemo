package main

// "C"

import "C"
import (
	"fmt"
	"strings"
	"unicode/utf16"
	"unsafe"

	"net/url"
	"os/exec"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

/*
*
0 for successful
*/

// //export Add
// func Add(a, b int) int {
// 	return a + b
// }

//export PrintMessage
func PrintMessage(msg *C.char) {
	// Convert C string to Go string
	message := C.GoString(msg)
	fmt.Println("Message from C#: ", message)
}

func test() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// 打开 Outlook 应用
	outlookApp, err := oleutil.CreateObject("Outlook.Application")
	// outlookApp, err := oleutil.GetActiveObject("Outlook.Application")
	if err != nil {
		fmt.Println("Error creating Outlook application:", err)
		return
	}
	defer outlookApp.Release()
}

//export AddReminder
func AddReminder(titleChar, dateChar, locationChar, attendeeChar *C.wchar_t, dur int) int {

	title := getString(titleChar)
	date := getString(dateChar)
	location := getString(locationChar)
	attendees := getString(attendeeChar)

	if dur == 0 {
		dur = 60
	}

	print(title, date, location, attendees)

	// C.free(unsafe.Pointer(titleChar))
	// c.free(unsafe.Pointer(dateChar))

	// 初始化 OLE
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// 打开 Outlook 应用
	outlookApp, err := oleutil.CreateObject("Outlook.Application")
	// outlookApp, err := oleutil.GetActiveObject("Outlook.Application")
	if err != nil {
		fmt.Println("Error creating Outlook application:", err)
		return 1001
	}
	defer outlookApp.Release()

	outlook, err := outlookApp.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("Error querying Outlook interface:", err)
		return 1002
	}
	defer outlook.Release()

	// 访问 MAPI 命名空间
	mapiNamespace, err := oleutil.CallMethod(outlook, "GetNamespace", "MAPI")
	if err != nil {
		fmt.Println("Error accessing MAPI namespace:", err)
		return 1003
	}
	defer mapiNamespace.ToIDispatch().Release()

	// 获取默认的日历文件夹（9 表示日历文件夹）
	calendarFolder := oleutil.MustCallMethod(mapiNamespace.ToIDispatch(), "GetDefaultFolder", 9).ToIDispatch()
	defer calendarFolder.Release()

	// 创建新的日历事件
	newAppointment := oleutil.MustCallMethod(calendarFolder, "Items").ToIDispatch()
	defer newAppointment.Release()

	// 配置事件属性
	appointment := oleutil.MustCallMethod(newAppointment, "Add").ToIDispatch()
	defer appointment.Release()

	// 设置事件信息
	oleutil.PutProperty(appointment, "Subject", title)
	// oleutil.PutProperty(appointment, "Location", "会议室A")
	// startTime := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	// print(startTime)
	// oleutil.PutProperty(appointment, "Start", startTime)
	oleutil.PutProperty(appointment, "Start", date)
	oleutil.PutProperty(appointment, "Duration", dur) // 持续时间（分钟）

	if strings.TrimSpace(location) != "" {
		oleutil.PutProperty(appointment, "Location", location)
	}
	// oleutil.PutProperty(appointment, "Body", content)
	oleutil.PutProperty(appointment, "ReminderSet", true)
	oleutil.PutProperty(appointment, "ReminderMinutesBeforeStart", 15)

	// 添加与会人
	attendeelist := strings.Split(attendees, ",")

	if len(attendeelist) > 0 {
		recipients := oleutil.MustGetProperty(appointment, "Recipients").ToIDispatch()
		defer recipients.Release()

		for _, a := range attendeelist {
			// 添加必需与会人
			addRecipient(recipients, a, 1) // 必需与会人
		}

	}

	// 保存日历事件
	oleutil.MustCallMethod(appointment, "Save")
	fmt.Println("Outlook 日历事件已成功添加！", date)
	return 0
}

func addRecipient(recipients *ole.IDispatch, email string, recipientType int) {
	recipient := oleutil.MustCallMethod(recipients, "Add", email).ToIDispatch()
	defer recipient.Release()

	// 设置与会人类型
	oleutil.PutProperty(recipient, "Type", recipientType)
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

func test2() {
	icsContent := `
BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:Team Meeting
DTSTART:20241120T100000Z
DTEND:20241120T110000Z
LOCATION:Conference Room
DESCRIPTION:Monthly team meeting.
END:VEVENT
END:VCALENDAR`

	mailTo := fmt.Sprintf("mailto:?subject=%s&body=%s",
		url.QueryEscape("Meeting Reminder"),
		url.QueryEscape(icsContent),
	)

	// Open in default mail client
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", mailTo).Start()
	if err != nil {
		fmt.Println("Failed to open mail client:", err)
	}
}

func main() {
	// AddReminder(C.CString("和川普吃烧烤"), C.CString("2024-11-28 13:00:00"))
	// test2()
}

// go build -o ./reminder.dll -buildmode=c-shared reminder.go
