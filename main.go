package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"
)

// matchAnyRegex 接受一段文本和一组正则表达式，如果任何一个正则表达式匹配到文本，则返回true
func matchAnyRegex(ctx context.Context, text string, patterns []string) bool {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // 确保所有goroutines最终都会被取消

	var wg sync.WaitGroup
	matchFound := make(chan bool, 1) // 用于通知找到匹配

	for _, pattern := range patterns {
		wg.Add(1)
		go func(pattern string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// 如果context被取消，提前退出goroutine
				return
			default:
				// 检查正则表达式是否匹配
				if matched, _ := regexp.MatchString("(?i)"+pattern, text); matched {
					select {
					case matchFound <- true:
						// 找到匹配，发送信号
					default:
						// 如果已经有信号发送，防止阻塞
					}
					cancel() // 取消其他goroutines
				}
			}
		}(pattern)
	}

	go func() {
		wg.Wait()
		close(matchFound) // 确保所有goroutine完成后关闭channel
	}()

	// 检查是否找到匹配
	found, ok := <-matchFound
	return ok && found
}

// matchInFile 打开一个文件，逐行读取并检查是否有任何行匹配给定的正则表达式列表
func matchInFile(filePath string, patterns []string, outputFilePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()
	outputFile, err := os.Create(outputFilePath) // 创建或打开输出文件
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		return
	}
	defer outputFile.Close() // 确保输出文件最终被关闭
	writer := bufio.NewWriter(outputFile)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 逐行读取文件内容
		line := scanner.Text()
		// 调用之前定义的匹配函数
		if matchAnyRegex(context.Background(), line, patterns) {
			fmt.Fprintln(writer, line) // 将匹配的行写入输出文件
			writer.Flush()             // 确保数据被写入文件
		}
	}

	// 检查扫描过程中是否有错误（非EOF）
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
	}

}

// 读取正则表达式文件
func readFileA(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err // 返回错误
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // 将读取的行添加到切片中
	}

	if err := scanner.Err(); err != nil {
		return nil, err // 如果读取过程中发生错误，返回错误
	}

	return lines, nil // 返回读取到的所有行的切片
}

func main() {
	urlFile := flag.String("f", "", "Path to the file containing URLs to scan")
	regexFile := flag.String("r", "", "Path to the file containing regular expressions")
	outputFile := flag.String("o", "", "Path to the output file for results")
	help := flag.Bool("help", false, "Display help message")

	// 解析命令行参数
	flag.Parse()

	// 如果指定了help选项或没有提供必要的参数，打印帮助信息
	if *help || *urlFile == "" || *regexFile == "" || *outputFile == "" {
		fmt.Println("Usage of this program:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	filePath := *urlFile          // url文件
	mafile := *regexFile          //正则表达式文件
	outfile := *outputFile        //输出结果的文件
	mastr, _ := readFileA(mafile) //读取正则表达式，返回string数组
	matchInFile(filePath, mastr, outfile)
}
