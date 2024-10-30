package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mozillazg/go-pinyin"
)

func main() {

	Test()
	return

	installedPrograms, err := GetInstalledPrograms()
	if err != nil {
		log.Fatalf("Failed to get installed programs: %v", err)
	}

	// hans := "中国人民银行"
	pyconf := pinyin.NewArgs()
	pyconf.Heteronym = true
	// ret := pinyin.Pinyin(hans, pyconf)

	// fmt.Println(ret)
	// fmt.Println(uniqueAndJoin(ret))
	// return

	for i := range installedPrograms {
		// 将 pinyin 拼接到原程序名称中
		installedPrograms[i].Name += uniqueAndJoin(pinyin.Pinyin(installedPrograms[i].Name, pyconf))

		// 打印更新后的名称和路径
		// fmt.Printf("Name: %s\nPath: %s\n\n", installedPrograms[i].Name, installedPrograms[i].Path)
	}

	reader := bufio.NewReader(os.Stdin)
	for {

		// var input string
		fmt.Print("---------")
		fmt.Print("请输入 app name: ")
		// fmt.Scanln(&input)
		input, _ := reader.ReadString('\n') // 读取输入直到换行
		input = strings.TrimSpace(input)

		if input == "quit" {
			return
		} else if input == "" {
			continue
		}

		start := time.Now()
		matchedPrograms := searchPrograms(installedPrograms, input)

		elapsed := time.Since(start)

		// 输出结果
		fmt.Printf("函数执行时间: %s\n", elapsed)

		// 输出结果
		if len(matchedPrograms) == 0 {
			fmt.Println("No matching programs found.")
		} else {
			fmt.Println("Matching programs:")
			for _, program := range matchedPrograms {
				fmt.Printf("Name: %s\nPath: %s\n\n", program.Name, program.Path)
			}
		}
	}
}

type Program struct {
	Name string
	Path string
}

// 去重并合并二维数组中的字符串
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
func isFuzzyMatch(programName, query string) bool {
	programName = strings.ToLower(programName)
	query = strings.ToLower(query)

	// fmt.Println(programName, query)

	j := 0 // query 的索引
	for i := 0; i < len(programName) && j < len(query); i++ {
		if programName[i] == query[j] {
			j++
		}
	}

	return j == len(query) // 如果 query 的所有字符都按顺序匹配到，则返回 true
}

// 模糊搜索函数
func searchPrograms(programs []Program, query string) []Program {
	fmt.Println(query)
	var results []Program
	for _, program := range programs {
		if isFuzzyMatch(program.Name, query) {
			results = append(results, program)
		}
	}
	return results
}
