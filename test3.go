package main

import (
	"fmt"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func Test3() {
	// 初始化 COM 库
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// 创建 Outlook 应用程序对象
	unknown, err := oleutil.CreateObject("Outlook.Application")
	if err != nil {
		log.Fatal("无法创建 Outlook.Application 对象:", err)
	}
	outlook, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal("无法查询 Outlook.Application 接口:", err)
	}
	defer outlook.Release()

	// 获取 MAPI 命名空间
	namespace, err := oleutil.CallMethod(outlook, "GetNamespace", "MAPI")
	if err != nil {
		log.Fatal("无法获取 MAPI 命名空间:", err)
	}
	defer namespace.ToIDispatch().Release()

	// 定义搜索关键字
	searchKeyword := "订单" // 替换为实际的搜索关键词

	// 搜索收件箱中的邮件
	inbox, err := oleutil.CallMethod(namespace.ToIDispatch(), "GetDefaultFolder", 6) // 6 表示收件箱
	if err != nil {
		log.Fatal("无法获取收件箱:", err)
	}
	defer inbox.ToIDispatch().Release()

	items := oleutil.MustGetProperty(inbox.ToIDispatch(), "Items").ToIDispatch()
	defer items.Release()

	// 执行筛选，使用 Restrict 方法
	filter := fmt.Sprintf("[Subject] = '%s'", searchKeyword)
	filteredItems := oleutil.MustCallMethod(items, "Restrict", filter).ToIDispatch()
	defer filteredItems.Release()

	// 遍历搜索结果
	count := oleutil.MustGetProperty(filteredItems, "Count").Val
	if count == 0 {
		fmt.Println("没有找到匹配的邮件")
		return
	}

	for i := 1; i <= int(count); i++ {
		item := oleutil.MustCallMethod(filteredItems, "Item", i).ToIDispatch()
		subject := oleutil.MustGetProperty(item, "Subject").ToString()
		fmt.Printf("邮件主题: %s\n", subject)
		item.Release()
	}
}
