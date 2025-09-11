/*
 * Copyright (C) 2020-2025. IntBoat <intboat@gmail.com> - All Rights Reserved.
 * Unauthorized copying or redistribution of this file in source and binary forms via any medium is strictly prohibited.
 * Proprietary and confidential
 *
 * @author    IntBoat <intboat@gmail.com>
 * @copyright 2020-2025. IntBoat <intboat@gmail.com>
 * @modified  2025-09-11 16:31:49
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 顏色定義
const (
	ColorRed    = "\033[0;31m"
	ColorGreen  = "\033[0;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[0;34m"
	ColorReset  = "\033[0m"
)

// 配置結構
type Config struct {
	FontFile   string
	DictFile   string
	OutputFile string
	FontIndex  int
	CacheSize  int
	ShowHelp   bool
}

// 進度追蹤器
type ProgressTracker struct {
	Total    int
	Current  int
	LastTime time.Time
}

// 初始化進度追蹤器
func NewProgressTracker(total int) *ProgressTracker {
	return &ProgressTracker{
		Total:    total,
		Current:  0,
		LastTime: time.Now(),
	}
}

// 更新進度
func (pt *ProgressTracker) Update(current int) {
	pt.Current = current
	now := time.Now()

	// 每 100ms 更新一次顯示，避免過於頻繁
	if now.Sub(pt.LastTime) >= 100*time.Millisecond {
		percentage := float64(current) * 100.0 / float64(pt.Total)
		fmt.Printf("\r%s[進度]%s 已處理 %d/%d (%.1f%%)",
			ColorBlue, ColorReset, current, pt.Total, percentage)
		pt.LastTime = now
	}
}

// 完成進度顯示
func (pt *ProgressTracker) Complete() {
	fmt.Printf("\r%s[完成]%s 已處理 %d/%d (100.0%%)\n",
		ColorGreen, ColorReset, pt.Total, pt.Total)
}

// 記錄訊息
func logInfo(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", ColorBlue, ColorReset, msg)
}

func logSuccess(msg string) {
	fmt.Printf("%s[SUCCESS]%s %s\n", ColorGreen, ColorReset, msg)
}

func logWarning(msg string) {
	fmt.Printf("%s[WARNING]%s %s\n", ColorYellow, ColorReset, msg)
}

func logError(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, msg)
}

// 顯示使用說明
func showUsage() {
	fmt.Printf("%sRime 字典過濾器%s\n", ColorBlue, ColorReset)
	fmt.Println()
	fmt.Println("用法: rime-filter [選項]")
	fmt.Println()
	fmt.Println("選項:")
	fmt.Println("  -f, --font <path>     指定 TTC 字型文件路徑")
	fmt.Println("  -d, --dict <path>     指定 Rime 字典文件路徑 (預設: quick5.dict.yaml)")
	fmt.Println("  -i, --index <num>     指定字型編號 (預設: 自動選擇)")
	fmt.Println("  -o, --output <path>   指定輸出文件路徑 (預設: filtered_dict.yaml)")
	fmt.Println("  -c, --cache <size>    指定緩存大小 (預設: 1000)")
	fmt.Println("  -h, --help           顯示此說明")
	fmt.Println()
	fmt.Println("範例:")
	fmt.Println("  rime-filter -f NotoSansCJK-Regular.ttc")
	fmt.Println("  rime-filter -f font.ttc -d my_dict.yaml -o result.yaml")
	fmt.Println("  rime-filter -f font.ttc -i 0 -c 2000")
}

// 解析命令行參數
func parseArgs() *Config {
	config := &Config{
		DictFile:   "quick5.dict.yaml",
		OutputFile: "filtered_dict.yaml",
		FontIndex:  -1,
		CacheSize:  1000,
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--font":
			if i+1 < len(args) {
				config.FontFile = args[i+1]
				i++
			}
		case "-d", "--dict":
			if i+1 < len(args) {
				config.DictFile = args[i+1]
				i++
			}
		case "-i", "--index":
			if i+1 < len(args) {
				if idx, err := strconv.Atoi(args[i+1]); err == nil {
					config.FontIndex = idx
				}
				i++
			}
		case "-o", "--output":
			if i+1 < len(args) {
				config.OutputFile = args[i+1]
				i++
			}
		case "-c", "--cache":
			if i+1 < len(args) {
				if size, err := strconv.Atoi(args[i+1]); err == nil {
					config.CacheSize = size
				}
				i++
			}
		case "-h", "--help":
			config.ShowHelp = true
		}
	}

	return config
}

// 檢查依賴
func checkDependencies() error {
	// 檢查 ttx 命令是否存在
	if _, err := exec.LookPath("ttx"); err != nil {
		return fmt.Errorf("ttx 命令未找到，請安裝 fonttools: pip install fonttools")
	}
	return nil
}

// 列出可用的字型編號
func listAvailableFonts(ttcFile string) ([]int, error) {
	var availableFonts []int

	for i := 0; i < 10; i++ {
		cmd := exec.Command("ttx", "-l", "-y", strconv.Itoa(i), ttcFile)
		if err := cmd.Run(); err == nil {
			availableFonts = append(availableFonts, i)
		}
	}

	return availableFonts, nil
}

// 提取字型中的中文字形
func extractFontChars(ttcFile string, fontIndex int, cacheSize int) ([]string, error) {
	logInfo("正在從字型文件中提取中文字形...")

	// 如果沒有指定字型編號，列出可用的字型
	if fontIndex == -1 {
		availableFonts, err := listAvailableFonts(ttcFile)
		if err != nil {
			return nil, err
		}

		logInfo("可用的字型編號：")
		for _, idx := range availableFonts {
			fmt.Printf("字型編號 %d: 可用\n", idx)
		}
		fmt.Println()

		fmt.Print("請輸入要提取的字型編號 (0-9): ")
		var input string
		fmt.Scanln(&input)

		if idx, err := strconv.Atoi(input); err == nil {
			fontIndex = idx
		} else {
			return nil, fmt.Errorf("無效的字型編號")
		}
	}

	// 使用 ttx 提取指定的字型
	logInfo(fmt.Sprintf("正在提取字型編號 %d...", fontIndex))

	cmd := exec.Command("ttx", "-o", "output.ttx", "-f", "-y", strconv.Itoa(fontIndex), ttcFile)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("提取字型失敗: %v", err)
	}
	defer os.Remove("output.ttx")

	// 讀取並解析 XML 文件
	logInfo("正在處理字形數據...")

	file, err := os.Open("output.ttx")
	if err != nil {
		return nil, fmt.Errorf("無法打開輸出文件: %v", err)
	}
	defer file.Close()

	// 使用正則表達式提取字形映射
	re := regexp.MustCompile(`<map code="0x([0-9a-fA-F]+)" name="[^"]+"`)
	scanner := bufio.NewScanner(file)

	var chars []string
	var cache []string

	// 先計算總行數
	totalLines := 0
	file.Seek(0, 0)
	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			totalLines++
		}
	}

	if totalLines == 0 {
		return nil, fmt.Errorf("未找到任何字形映射")
	}

	logInfo(fmt.Sprintf("找到 %d 個字形映射，開始處理...", totalLines))

	// 重新讀取文件進行處理
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	progress := NewProgressTracker(totalLines)
	processed := 0

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			// 將十六進制轉換為 Unicode 字符
			code, err := strconv.ParseInt(matches[1], 16, 32)
			if err != nil {
				continue
			}

			char := string(rune(code))
			cache = append(cache, char)

			// 當緩存滿了時，批量添加到結果中
			if len(cache) >= cacheSize {
				chars = append(chars, cache...)
				cache = cache[:0] // 清空緩存
			}

			processed++
			progress.Update(processed)
		}
	}

	// 添加剩餘的緩存數據
	if len(cache) > 0 {
		chars = append(chars, cache...)
	}

	progress.Complete()

	if len(chars) == 0 {
		return nil, fmt.Errorf("未能提取任何中文字形，請檢查字型文件或其內容")
	}

	// 將字符寫入文件
	outputFile := "chinese_characters.txt"
	charFile, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("無法創建字符文件: %v", err)
	}
	defer charFile.Close()

	writer := bufio.NewWriter(charFile)
	for _, char := range chars {
		if _, err := writer.WriteString(char + "\n"); err != nil {
			return nil, fmt.Errorf("寫入字符文件失敗: %v", err)
		}
	}
	writer.Flush()

	logSuccess(fmt.Sprintf("已提取 %d 個中文字形到 %s", len(chars), outputFile))
	return chars, nil
}

// 過濾字典文件
func filterDict(dictFile, outputFile string, validChars []string, cacheSize int) error {
	logInfo("正在過濾字典文件...")

	// 檢查輸入文件是否存在
	if _, err := os.Stat(dictFile); os.IsNotExist(err) {
		return fmt.Errorf("字典文件不存在: %s", dictFile)
	}

	// 創建有效字符的映射表
	validMap := make(map[string]bool)
	for _, char := range validChars {
		validMap[char] = true
	}

	// 打開輸入和輸出文件
	inputFile, err := os.Open(dictFile)
	if err != nil {
		return fmt.Errorf("無法打開字典文件: %v", err)
	}
	defer inputFile.Close()

	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("無法創建輸出文件: %v", err)
	}
	defer outputFileHandle.Close()

	missingFile, err := os.Create("missing_chars.txt")
	if err != nil {
		return fmt.Errorf("無法創建缺失字符文件: %v", err)
	}
	defer missingFile.Close()

	// 計算總行數
	scanner := bufio.NewScanner(inputFile)
	totalLines := 0
	for scanner.Scan() {
		totalLines++
	}

	logInfo(fmt.Sprintf("字典文件共有 %d 行，開始過濾...", totalLines))

	// 重新讀取文件進行處理
	inputFile.Seek(0, 0)
	scanner = bufio.NewScanner(inputFile)

	var outputCache []string
	var missingCache []string
	outputCacheCount := 0
	missingCacheCount := 0

	processedLines := 0
	totalDictLines := 0
	validLines := 0
	missingLines := 0

	progress := NewProgressTracker(totalLines)

	// 表頭處理狀態
	inHeader := false
	headerFound := false
	headerHandled := false

	for scanner.Scan() {
		line := scanner.Text()
		processedLines++

		trimmed := strings.TrimSpace(line)

		// 若尚未處理表頭，且第一行為 '---'，進入表頭狀態並直接輸出
		if !headerHandled {
			if !inHeader {
				if processedLines == 1 && trimmed == "---" {
					inHeader = true
					headerFound = true
					if _, err := fmt.Fprintln(outputFileHandle, line); err != nil {
						return err
					}
					continue
				}
			} else {
				// 表頭區塊內：逐行輸出，直到遇到 '...'
				if _, err := fmt.Fprintln(outputFileHandle, line); err != nil {
					return err
				}
				if trimmed == "..." {
					inHeader = false
					headerHandled = true
				}
				continue
			}
		}

		// 若已有表頭開始但尚未寫出結束分界，在遇到內容前先補一行 '...'
		if headerFound && !inHeader && !headerHandled {
			if _, err := fmt.Fprintln(outputFileHandle, "..."); err != nil {
				return err
			}
			headerHandled = true
		}

		// 保留空行與註釋
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			outputCache = append(outputCache, line)
			outputCacheCount++
			if outputCacheCount >= cacheSize {
				writeLines(outputFileHandle, outputCache)
				outputCache = outputCache[:0]
				outputCacheCount = 0
			}
			progress.Update(processedLines)
			continue
		}

		// 分界符（--- 或 ...）在內容區也原樣保留
		if trimmed == "---" || trimmed == "..." {
			outputCache = append(outputCache, line)
			outputCacheCount++
			if outputCacheCount >= cacheSize {
				writeLines(outputFileHandle, outputCache)
				outputCache = outputCache[:0]
				outputCacheCount = 0
			}
			progress.Update(processedLines)
			continue
		}

		// 內容條目：以首欄為字元進行過濾
		fields := strings.Fields(line)
		if len(fields) > 0 {
			char := fields[0]
			totalDictLines++

			if validMap[char] {
				outputCache = append(outputCache, line)
				outputCacheCount++
				validLines++
			} else {
				missingCache = append(missingCache, char)
				missingCacheCount++
				missingLines++
			}

			if outputCacheCount >= cacheSize {
				writeLines(outputFileHandle, outputCache)
				outputCache = outputCache[:0]
				outputCacheCount = 0
			}
			if missingCacheCount >= cacheSize {
				writeLines(missingFile, missingCache)
				missingCache = missingCache[:0]
				missingCacheCount = 0
			}
		}

		progress.Update(processedLines)
	}

	// 寫入剩餘的緩存數據
	if len(outputCache) > 0 {
		writeLines(outputFileHandle, outputCache)
	}
	if len(missingCache) > 0 {
		writeLines(missingFile, missingCache)
	}

	progress.Complete()

	logSuccess("過濾完成！")
	logInfo(fmt.Sprintf("總行數: %d", totalDictLines))
	logInfo(fmt.Sprintf("有效行數: %d", validLines))
	logInfo(fmt.Sprintf("缺失行數: %d", missingLines))
	logInfo(fmt.Sprintf("過濾後的文件: %s", outputFile))
	logInfo(fmt.Sprintf("缺失的字: missing_chars.txt"))

	return nil
}

// 批量寫入行到文件
func writeLines(writer io.Writer, lines []string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}
	return nil
}

// 主函數
func main() {
	config := parseArgs()

	if config.ShowHelp {
		showUsage()
		return
	}

	// 檢查是否提供了 TTC 文件
	if config.FontFile == "" {
		logError("請指定 TTC 字型文件")
		showUsage()
		os.Exit(1)
	}

	// 檢查 TTC 文件是否存在
	if _, err := os.Stat(config.FontFile); os.IsNotExist(err) {
		logError(fmt.Sprintf("TTC 字型文件不存在: %s", config.FontFile))
		os.Exit(1)
	}

	// 檢查依賴
	if err := checkDependencies(); err != nil {
		logError(err.Error())
		os.Exit(1)
	}

	// 提取字型中的中文字形
	validChars, err := extractFontChars(config.FontFile, config.FontIndex, config.CacheSize)
	if err != nil {
		logError(err.Error())
		os.Exit(1)
	}

	// 過濾字典文件
	if err := filterDict(config.DictFile, config.OutputFile, validChars, config.CacheSize); err != nil {
		logError(err.Error())
		os.Exit(1)
	}

	logSuccess("所有操作完成！")
}
