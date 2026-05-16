package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

type FilterCond struct {
	Field    string
	Op       string
	NumValue float64
	StrValue string
	idx      int
	isNum    bool
}

// 条件式を解析して、フィールド名、演算子、値を抽出する
func parseCondition(cond string) (FilterCond, error) {
	filterCond := FilterCond{}

	if strings.Contains(cond, ">=") {
		parts := strings.SplitN(cond, ">=", 2)
		filterCond.Field = strings.TrimSpace(parts[0])
		filterCond.Op = ">="
		filterCond.StrValue = strings.TrimSpace(parts[1])
	} else if strings.Contains(cond, "<=") {
		parts := strings.SplitN(cond, "<=", 2)
		filterCond.Field = strings.TrimSpace(parts[0])
		filterCond.Op = "<="
		filterCond.StrValue = strings.TrimSpace(parts[1])
	} else if strings.Contains(cond, "!=") {
		parts := strings.SplitN(cond, "!=", 2)
		filterCond.Field = strings.TrimSpace(parts[0])
		filterCond.Op = "!="
		filterCond.StrValue = strings.TrimSpace(parts[1])
	} else if strings.Contains(cond, "=") {
		parts := strings.SplitN(cond, "=", 2)
		filterCond.Field = strings.TrimSpace(parts[0])
		filterCond.Op = "="
		filterCond.StrValue = strings.TrimSpace(parts[1])
	} else if strings.Contains(cond, ">") {
		parts := strings.SplitN(cond, ">", 2)
		filterCond.Field = strings.TrimSpace(parts[0])
		filterCond.Op = ">"
		filterCond.StrValue = strings.TrimSpace(parts[1])
	} else if strings.Contains(cond, "<") {
		parts := strings.SplitN(cond, "<", 2)
		filterCond.Field = strings.TrimSpace(parts[0])
		filterCond.Op = "<"
		filterCond.StrValue = strings.TrimSpace(parts[1])
	} else {
		return filterCond, fmt.Errorf("invalid filter condition: %s", cond)
	}

	if n, err := strconv.ParseFloat(filterCond.StrValue, 64); err == nil {
		filterCond.NumValue = n
		filterCond.isNum = true
	}

	return filterCond, nil
}

func (t *Table) Filter(expr string) (*Table, error) {
	// 条件がカンマ区切りで複数指定されるので、分割して処理する
	conditions := strings.Split(expr, ",")
	for i := range conditions {
		conditions[i] = strings.TrimSpace(conditions[i])
	}

	// 条件式をfilterCondsに変換
	filterConds := make([]FilterCond, 0, len(conditions))
	for _, cond := range conditions {
		// 条件を解析して、フィルタリングする
		// 例: "amount>=1000" -> フィールド名: "amount", 演算子: ">=", 値: "1000"

		filterCond, err := parseCondition(cond)
		if err != nil {
			return nil, err
		}

		fieldIdx, err := t.colIndex(filterCond.Field)
		if err != nil {
			return nil, err
		}

		filterCond.idx = fieldIdx
		filterConds = append(filterConds, filterCond)
	}

	// 行をフィルタリング
	var filteredRows [][]string
	for _, row := range t.Rows {
		if matchesConditions(row, filterConds) {
			filteredRows = append(filteredRows, row)
		}
	}

	return &Table{Headers: t.Headers, Rows: filteredRows}, nil
}

// 行文字列をパースして、条件式にマッチするかどうかを判定する
func matchesConditions(row []string, conds []FilterCond) bool {
	for _, cond := range conds {
		cellValue := row[cond.idx]

		// 数値条件の場合は、セルの値を数値に変換して比較
		if cond.isNum {
			cellNum, err := strconv.ParseFloat(cellValue, 64)
			if err != nil {
				return false
			}

			switch cond.Op {
			case ">=":
				if !(cellNum >= cond.NumValue) {
					return false
				}

			case "<=":
				if !(cellNum <= cond.NumValue) {
					return false
				}

			case ">":
				if !(cellNum > cond.NumValue) {
					return false
				}
			case "<":
				if !(cellNum < cond.NumValue) {
					return false
				}
			case "!=":
				if !(cellNum != cond.NumValue) {
					return false
				}
			case "=":
				if !(cellNum == cond.NumValue) {
					return false
				}
			}
		} else {
			// 文字列条件の場合は単純に比較
			switch cond.Op {

			case "!=":
				if cellValue == cond.StrValue {
					return false
				}
			case "=":
				if cellValue != cond.StrValue {
					return false
				}
			default:
				return false

			}
		}

	}

	return true
}

// ReadCSV CSVファイルを読み込んで、Table構造体を返す
func ReadCSV(filename string) (*Table, error) {

	// ファイルオープン
	in, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	reader := csv.NewReader(in)

	table := &Table{}

	// ヘッダを読む
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	table.Headers = header

	// データ行を読む

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		table.Rows = append(table.Rows, record)

	}

	return table, nil

}

type SumField struct {
	Field string
	idx   int
}

type GroupField struct {
	Field string
	idx   int
}

func (t *Table) colIndex(colName string) (int, error) {
	for i, h := range t.Headers {
		if h == colName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("field %q not found in headers", colName)
}

func (t *Table) GroupAndSum(groupStr, sumStr string) (*Table, error) {
	// 条件がカンマ区切りで複数指定されるので、分割して処理する
	groupNames := strings.Split(groupStr, ",")
	for i := range groupNames {
		groupNames[i] = strings.TrimSpace(groupNames[i])
	}

	sumNames := strings.Split(sumStr, ",")
	for i := range sumNames {
		sumNames[i] = strings.TrimSpace(sumNames[i])
	}

	// 指定カラム名をsumFiledsに変換
	sumFields := make([]SumField, 0, len(sumNames))
	for _, sumName := range sumNames {
		fieldIdx, err := t.colIndex(sumName)
		if err != nil {
			return nil, err
		}
		sumFields = append(sumFields, SumField{Field: sumName, idx: fieldIdx})
	}

	// 指定カラム名をgroupFiledsに変換
	groupFields := make([]GroupField, 0, len(groupNames))
	for _, groupName := range groupNames {
		fieldIdx, err := t.colIndex(groupName)
		if err != nil {
			return nil, err
		}
		groupFields = append(groupFields, GroupField{Field: groupName, idx: fieldIdx})
	}

	// sumFields,groupFieldsでヘッダを作り直す
	var newHeaders []string
	for _, gf := range groupFields {
		newHeaders = append(newHeaders, gf.Field)
	}
	for _, sf := range sumFields {
		newHeaders = append(newHeaders, sf.Field)
	}

	// groupFieldsをキーにして、sumfieldsを集計する
	groupSumMap := make(map[string][]float64)

	for _, row := range t.Rows {
		// groupFieldsをキーにする
		groupKeyParts := make([]string, len(groupFields))
		for i, gf := range groupFields {
			groupKeyParts[i] = row[gf.idx]
		}
		groupKey := strings.Join(groupKeyParts, "\x00")

		// sumFieldsのカラム値を数値に変換して、スライスを作成
		sumValues := make([]float64, len(sumFields))
		for i, sf := range sumFields {
			n, err := strconv.ParseFloat(row[sf.idx], 64)
			if err != nil {
				n = 0
			}
			sumValues[i] = n
		}

		// mapにスライスに設定した値を集計
		if existing, ok := groupSumMap[groupKey]; ok {
			for i := range sumValues {
				existing[i] += sumValues[i]
			}

		} else {
			groupSumMap[groupKey] = sumValues
		}
	}

	// groupSumMapに集計結果がまとまっているので、行のスライスに変換する
	// groupKey, sumValues -> groupKeyParts, sumValues, ->row
	var groupSumRows [][]string
	for groupKey, sumValues := range groupSumMap {
		groupKeyParts := strings.Split(groupKey, "\x00")
		var row []string
		row = append(row, groupKeyParts...)
		for _, sumValue := range sumValues {
			row = append(row, strconv.FormatFloat(sumValue, 'f', -1, 64))
		}
		groupSumRows = append(groupSumRows, row)
	}

	return &Table{Headers: newHeaders, Rows: groupSumRows}, nil
}

func (t *Table) WriteCSV(filename string) error {
	// ファイルオープン
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(out)
	defer out.Close()

	// ヘッダを書き込む
	err = writer.Write(t.Headers)
	if err != nil {
		return err
	}

	// データ行を書き込む
	for _, row := range t.Rows {
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

func (t *Table) WriteTSV(w io.Writer) error {
	// ヘッダを書き込む
	fmt.Fprintln(w, strings.Join(t.Headers, "\t"))

	// データ行を書き込む
	for _, row := range t.Rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return nil
}
