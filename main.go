package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 日付抽出用の正規表現
var dateRegex = regexp.MustCompile(`_?(\d{8})_`)

// extractDate はファイル名から日付を抽出し、time.Time型で返します
func extractDate(filename string) (time.Time, error) {
	matches := dateRegex.FindStringSubmatch(filename)
	if len(matches) > 1 {
		parsedDate, err := time.Parse("20060102", matches[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("日付の変換に失敗 (%s): %v", filename, err)
		}
		return parsedDate, nil
	}
	return time.Time{}, fmt.Errorf("日付フォーマットなし")
}

// ensureDir はディレクトリが存在するかチェックし、なければ作成します
func ensureDir(dirName string) error {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("同名のファイルが存在します: %s", dirName)
	}
	return nil
}

func main() {
	// カレントディレクトリを対象
	dir := "."
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("ディレクトリの読み込みに失敗しました: %v", err)
	}

	successCount := 0
	skipCount := 0

	fmt.Println("=== 仕分け処理を開始します ===")

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		ext := strings.ToLower(filepath.Ext(filename))

		// 対象となる拡張子をフィルタリング
		if ext != ".mp4" && ext != ".insv" && ext != ".lrv" {
			continue
		}

		// 1. 日付を抽出
		t, err := extractDate(filename)
		if err != nil {
			fmt.Printf("[スキップ] %-30s (理由: %v)\n", filename, err)
			skipCount++
			continue
		}

		// 2. 移動先のディレクトリ名を決定し、存在を保証する
		folderName := t.Format("2006-01-02")
		if err := ensureDir(folderName); err != nil {
			fmt.Printf("[エラー] %-32s (理由: ディレクトリ作成失敗 %v)\n", filename, err)
			skipCount++
			continue
		}

		// 3. ファイルパスを構築して移動（リネーム）
		// filepath.Joinを使うことで、OSごとのパス区切り文字の違いを吸収します
		newPath := filepath.Join(folderName, filename)

		if err := os.Rename(filename, newPath); err != nil {
			fmt.Printf("[エラー] %-32s (理由: 移動失敗 %v)\n", filename, err)
			skipCount++
			continue
		}

		fmt.Printf("[移動完了] %-30s -> %s/\n", filename, folderName)
		successCount++
	}

	fmt.Println("=== 処理完了 ===")
	fmt.Printf("移動成功: %d 件, スキップ/エラー: %d 件\n", successCount, skipCount)
}
