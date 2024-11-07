package main

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	// "C"
)

/*
*
0 for successful
*/
func AddReminder(title string, date string) int {
	// 初始化 OLE
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// 打开 Outlook 应用
	outlookApp, err := oleutil.CreateObject("Outlook.Application")
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
	oleutil.PutProperty(appointment, "Duration", 60) // 持续时间（分钟）
	// oleutil.PutProperty(appointment, "Body", content)
	oleutil.PutProperty(appointment, "ReminderSet", true)
	oleutil.PutProperty(appointment, "ReminderMinutesBeforeStart", 15)

	// 保存日历事件
	oleutil.MustCallMethod(appointment, "Save")
	fmt.Println("Outlook 日历事件已成功添加！")
	return 0
}

func main() {
	AddReminder("test subject", "")
}
