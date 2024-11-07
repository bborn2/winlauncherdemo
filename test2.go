package main

import (
	"fmt"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func Test2() {
	// 初始化 COM
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

	// 获取命名空间 MAPI
	namespace, err := oleutil.CallMethod(outlook, "GetNamespace", "MAPI")
	if err != nil {
		log.Fatal("无法获取命名空间:", err)
	}
	defer namespace.ToIDispatch().Release()

	// 搜索关键字
	searchKeyword := "your search keyword" // 替换为你的搜索关键字

	// 执行高级搜索
	search, err := oleutil.CallMethod(outlook, "AdvancedSearch", "Inbox", fmt.Sprintf("urn:schemas:httpmail:subject LIKE '%s'", searchKeyword))
	if err != nil {
		log.Fatal("无法执行高级搜索:", err)
	}
	defer search.ToIDispatch().Release()

	fmt.Println("Outlook 搜索已执行，搜索关键字为:", searchKeyword)
}
