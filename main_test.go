package main

import (
	"testing"
)

func TestExtractDate(t *testing.T) {
	// テストケースの定義（テーブル駆動テスト）
	tests := []struct {
		name        string // テストケース名
		filename    string // 入力するファイル名
		wantDateStr string // 期待する日付文字列（検証しやすいようにYYYY-MM-DD形式にする）
		wantErr     bool   // エラーが発生することを期待するかどうか
	}{
		{
			name:        "標準的なInsta360ファイル",
			filename:    "VID_20260414_160540_003.mp4",
			wantDateStr: "2026-04-14",
			wantErr:     false,
		},
		{
			name:        "Proモードのファイル",
			filename:    "PRO_VID_20260415_102030_001.mp4",
			wantDateStr: "2026-04-15",
			wantErr:     false,
		},
		{
			name:        "プロキシ・独自形式",
			filename:    "LRV_20260416_112233_002.insv",
			wantDateStr: "2026-04-16",
			wantErr:     false,
		},
		{
			name:        "日付を含まないファイル",
			filename:    "just_a_normal_video.mp4",
			wantDateStr: "",
			wantErr:     true, // エラーになるのが正解
		},
		{
			name:        "8桁だが存在しない日付（13月45日）",
			filename:    "VID_20261345_160540_003.mp4",
			wantDateStr: "",
			wantErr:     true, // time.Parseで弾かれるのでエラーになるのが正解
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := extractDate(tt.filename)

			// 期待通りのエラー状態かチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("extractDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーを期待していない（正常系）のテストの場合は、日付の値をチェック
			if !tt.wantErr {
				// 比較しやすいようにフォーマットして比較
				gotDateStr := gotTime.Format("2006-01-02")
				if gotDateStr != tt.wantDateStr {
					t.Errorf("extractDate() got = %v, want %v", gotDateStr, tt.wantDateStr)
				}
			}
		})
	}
}
