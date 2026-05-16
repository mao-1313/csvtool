package main

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	// テスト用のCSVファイルを作成
	csvContent := `category,amount
食品,500
電化製品,80000
食品,200`
	inputFile := "test_input.csv"
	err := os.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test CSV file: %v", err)
	}
	defer os.Remove(inputFile)

	// ReadCSV を実行
	table, err := ReadCSV(inputFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 結果を検証
	if len(table.Headers) != 2 || table.Headers[0] != "category" || table.Headers[1] != "amount" {
		t.Errorf("unexpected headers: %v", table.Headers)
	}
	if len(table.Rows) != 3 {
		t.Errorf("got %d rows, want 3", len(table.Rows))
	}
	if table.Rows[0][0] != "食品" || table.Rows[0][1] != "500" {
		t.Errorf("unexpected first row: %v", table.Rows[0])
	}

}

func TestReadFail(t *testing.T) {
	// テスト用のCSVファイルを作成
	inputFile := "test_input.csv"

	// ReadCSV を実行
	_, err := ReadCSV(inputFile)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

}

func TestFilterNum(t *testing.T) {
	// 1. テスト用の Table を直接作る（ファイル不要）
	table := &Table{
		Headers: []string{"category", "amount"},
		Rows: [][]string{
			{"食品", "500"},
			{"電化製品", "80000"},
			{"食品", "200"},
		},
	}

	// 2. Filter を実行
	got, err := table.Filter("amount>=1000")

	// 3. 結果を検証
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Rows) != 1 {
		t.Errorf("got %d rows, want 1", len(got.Rows))
	}
	if got.Rows[0][0] != "電化製品" || got.Rows[0][1] != "80000" {
		t.Errorf("unexpected row content: %v", got.Rows[0])
	}
}

func TestFilterStr(t *testing.T) {
	// 1. テスト用の Table を直接作る（ファイル不要）
	table := &Table{
		Headers: []string{"category", "amount"},
		Rows: [][]string{
			{"食品", "500"},
			{"電化製品", "80000"},
			{"食品", "200"},
		},
	}

	// 2. Filter を実行
	got, err := table.Filter("category=食品")

	// 3. 結果を検証
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Rows) != 2 {
		t.Errorf("got %d rows, want 2", len(got.Rows))
	}
}

func TestFilter2(t *testing.T) {
	// 1. テスト用の Table を直接作る（ファイル不要）
	table := &Table{
		Headers: []string{"category", "amount"},
		Rows: [][]string{
			{"食品", "500"},
			{"電化製品", "80000"},
			{"食品", "200"},
		},
	}

	// 2. Filter を実行
	got, err := table.Filter("category=食品,amount>=300")

	// 3. 結果を検証
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Rows) != 1 {
		t.Errorf("got %d rows, want 1", len(got.Rows))
	}
	if got.Rows[0][0] != "食品" || got.Rows[0][1] != "500" {
		t.Errorf("unexpected row content: %v", got.Rows[0])
	}
}

func TestGroupAndSum(t *testing.T) {
	table := &Table{
		Headers: []string{"category", "amount"},
		Rows: [][]string{
			{"食品", "500"},
			{"電化製品", "50000"},
			{"電化製品", "40000"},
			{"食品", "200"},
			{"食品", "300"},
			{"食品", "400"},
			{"食品", "100"},
		},
	}

	got, err := table.GroupAndSum("category", "amount")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Rows) != 2 {
		t.Errorf("got %d rows, want 2", len(got.Rows))
	}

	// 集計結果を確認
	sums := map[string]string{}
	for _, row := range got.Rows {
		sums[row[0]] = row[1]
	}

	if sums["食品"] != "1500" {
		t.Errorf("食品 sum =%s, want 1500", sums["食品"])
	}
	if sums["電化製品"] != "90000" {
		t.Errorf("電化製品 sum =%s, want 90000", sums["電化製品"])
	}
}

func TestFilters(t *testing.T) {
	table := &Table{
		Headers: []string{"category", "amount"},
		Rows: [][]string{
			{"食品", "500"},
			{"電化製品", "50000"},
			{"電化製品", "40000"},
			{"食品", "200"},
			{"食品", "300"},
			{"食品", "400"},
			{"食品", "100"},
		},
	}

	tests := []struct {
		expr    string
		wantLen int
	}{
		{"amount>=1000", 2},
		{"category=食品", 5},
		{"category=食品,amount>=300", 3},
	}

	for _, test := range tests {
		// forループでテストしてエラーチェックをする場合はサブテストとする
		t.Run(test.expr, func(t *testing.T) {
			got, err := table.Filter(test.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got.Rows) != test.wantLen {
				t.Errorf("got %d rows, want %d", len(got.Rows), test.wantLen)
			}

		})
	}
}
