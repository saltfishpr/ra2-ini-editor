package ra2

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTranslation(t *testing.T) {
	translation, err := LoadTranslation("testdata/ra2md.ini", "zh-TW")
	if err != nil {
		t.Fatalf("failed to load translation: %v", err)
	}
	_ = translation

	// 打开文件
	file, err := os.Open("testdata/ra2md.ini")
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 创建一个按行读取的 Scanner
	scanner := bufio.NewScanner(file)

	// 按行读取
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(pair[0])
		fmt.Println(key)
	}
}
