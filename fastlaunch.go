package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

import "C"

type ProgramData struct {
	Name       string
	SearchName string
	APPID      string
	Weight     int
}

var apps []ProgramData

//export loadApps
func loadApps() int {
	a, err := getStartApps()
	apps = a

	if err != nil {
		print("err code : 1200")
		return 1
	}

	return 0

	// for _, a := range apps {
	// 	fmt.Println(a)
	// }
	// fmt.Println(apps)
}

//export searchAndRun
func searchAndRun(queryChar *C.char) int {
	query := C.GoString(queryChar)

	app := searchApp(query)
	// fmt.Println(app)

	if len(app) > 0 {
		openApp(app[0].APPID)

		return 1
	}

	return 0
}

// 模糊搜索函数
func searchApp(query string) []ProgramData {
	// fmt.Println(query)
	var results []ProgramData
	for _, program := range apps {

		fmt.Println(program)

		ret, weight := isFuzzyMatch(program.SearchName, query)
		if ret {
			program.Weight = weight
			results = insertInOrder(results, program)
		}
	}
	return results
}

func openApp(app string) {

	print("Start-Process \"Shell:AppsFolder\\" + app + "\"")

	out, err := exec.Command("powershell", "Start-Process", "\"Shell:AppsFolder\\"+app+"\"").CombinedOutput()

	if err != nil {
		fmt.Printf("err %v\n", err)
	}

	fmt.Printf("output %s\n", out)
}

// 通过 PowerShell 获取所有开始菜单的应用
func getStartApps() ([]ProgramData, error) {
	// PowerShell 命令获取开始菜单应用
	cmd := exec.Command("powershell", "chcp 65001; Get-StartApps | Select Name")

	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令并捕获输出
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	pyconf := pinyin.NewArgs()
	pyconf.Heteronym = true

	// 解析 PowerShell 输出
	lines := strings.Split(out.String(), "\n")

	var programs []ProgramData

	// 跳过第一行标题并遍历每一行数据
	for _, line := range lines[4:] { // 跳过标题行

		// print(line + "\n")
		// continue
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		program := ProgramData{
			Name:       strings.TrimSpace(line),
			SearchName: strings.TrimSpace(line) + uniqueAndJoin(pinyin.Pinyin(strings.TrimSpace(line), pyconf)), // 应用程序的名称
			APPID:      "",
			Weight:     0, // 应用程序的 APPID
		}
		programs = append(programs, program)
	}

	cmd = exec.Command("powershell", "chcp 65001; Get-StartApps | Select AppID")

	var out2 bytes.Buffer
	cmd.Stdout = &out2

	// 执行命令并捕获输出
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// 解析 PowerShell 输出
	lines = strings.Split(out2.String(), "\n")
	for index, line := range lines[4:] { // 跳过标题行

		// 假设 Name 和 APPID 之间有空格
		// parts := strings.Fields(line)

		// print(index)
		line = strings.TrimSpace(line)

		if index < len(programs) {
			programs[index].APPID = strings.TrimSpace(line)
		}

	}

	return programs, nil
}

func uniqueAndJoin(arr [][]string) string {
	uniqueMap := make(map[string]bool)
	var result []string

	// 遍历二维数组，将每个元素放入 map 以去重
	for _, subArr := range arr {
		for _, str := range subArr {
			if _, exists := uniqueMap[str]; !exists {
				uniqueMap[str] = true
				result = append(result, str)
			}
		}
	}

	// 使用空格连接去重后的字符串
	return strings.Join(result, " ")
}

// 判断 programName 是否包含 query 的所有字符并保持顺序
func isFuzzyMatch(programName, query string) (bool, int) {
	programName = strings.ToLower(programName)
	query = strings.ToLower(query)

	// fmt.Println(programName, query)

	weight := 0

	delta := -1
	j := 0 // query 的索引
	for i := 0; i < len(programName) && j < len(query); i++ {
		if programName[i] == query[j] {

			if i-j == delta {
				weight += 10
			} else {

				delta = i - j
			}

			j++
		}

		if j == len(query) {
			break
		}
	}

	return j == len(query), weight // 如果 query 的所有字符都按顺序匹配到，则返回 true
}

// 按照 Weight 值顺序插入数据
func insertInOrder(items []ProgramData, newItem ProgramData) []ProgramData {
	// 找到第一个 Weight 大于 newItem.Weight 的位置
	index := 0
	for i, item := range items {
		if item.Weight < newItem.Weight {
			index = i
			break
		}
		index = i + 1
	}
	// 插入数据到指定位置，保持顺序
	items = append(items[:index], append([]ProgramData{newItem}, items[index:]...)...)
	return items
}

func main() {
	// loadApps()

	// app := searchApp("word")
	// fmt.Println(app)

	// if len(app) > 0 {
	// 	openApp(app[0].APPID)
	// }

	// return
}
