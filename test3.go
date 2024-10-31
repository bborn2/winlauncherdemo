package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// 定义 Windows API 函数
var (
	modole32         = syscall.NewLazyDLL("ole32.dll")
	coInitializeEx   = modole32.NewProc("CoInitializeEx")
	coCreateInstance = modole32.NewProc("CoCreateInstance")
	coUninitialize   = modole32.NewProc("CoUninitialize")
)

// CLSID 和 IID 的定义
var (
	CLSID_OutlookApplication = syscall.GUID{0x0006F03A, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IDispatch            = syscall.GUID{0x00020400, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
)

func main() {
	// 初始化 COM 库
	hr, _, _ := coInitializeEx.Call(0, 0x20)
	if hr != 0 {
		fmt.Println("CoInitializeEx failed")
		return
	}
	defer coUninitialize.Call()

	// 创建 Outlook 应用程序实例
	var outApp unsafe.Pointer
	hr, _, _ = coCreateInstance.Call(
		uintptr(unsafe.Pointer(&CLSID_OutlookApplication)),
		0,
		0x17,
		uintptr(unsafe.Pointer(&IID_IDispatch)),
		uintptr(unsafe.Pointer(&outApp)),
	)
	if hr != 0 {
		fmt.Println("CoCreateInstance failed")
		return
	}

	// 执行搜索操作
	searchString := "subject:test"
	hr, _, _ = syscall.Syscall(
		outApp,
		0x60020000,
		uintptr(unsafe.Pointer(&searchString)),
	)
	if hr != 0 {
		fmt.Println("Search failed")
		return
	}

	fmt.Println("Search completed successfully")
}
