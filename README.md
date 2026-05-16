# 課題10：CSV変換・集計ツール

売上CSVデータに対してフィルタ・集計を行うCLIツールを実装する。

---

## 完成イメージ

```bash
# 金額1000以上の行を抽出して表示
$ go run . --filter "amount>=1000" sales.csv

# カテゴリ別に売上を集計
$ go run . --group category --sum amount sales.csv

# フィルタ＋集計を組み合わせる
$ go run . --filter "amount>=1000" --group category --sum amount sales.csv

# 結果をCSVファイルに出力
$ go run . --group category --sum amount sales.csv --out result.csv
```

---

## サンプルCSV (`sales.csv`)

```csv
date,category,product,amount
2026-01-01,食品,りんご,500
2026-01-01,電化製品,テレビ,80000
2026-01-02,食品,バナナ,200
2026-01-02,電化製品,冷蔵庫,120000
2026-01-03,食品,みかん,1500
2026-01-03,電化製品,電子レンジ,30000
```

---

## 学習ステップ

### Step 1 - CSVの読み込み

`encoding/csv` でCSVを読み込み、ヘッダーと行データを構造体で管理する。

```go
type Table struct {
    Headers []string
    Rows    [][]string
}

func ReadCSV(filename string) (*Table, error)
```

### Step 2 - フィルタ機能

`--filter` フラグで条件を指定し、条件に合う行だけ抽出する。

```go
// 対応する演算子: >=, <=, >, <, =
// 例: "amount>=1000"

func (t *Table) Filter(expr string) (*Table, error)
```

### Step 3 - 集計機能

`--group` でグループ化するカラム、`--sum` で合計するカラムを指定して集計する。

```go
func (t *Table) GroupSum(groupCol, sumCol string) (*Table, error)
```

### Step 4 - 結果の出力

ターミナルへの表示と `--out` フラグでCSVファイルへの出力に対応する。

```go
func (t *Table) Print()
func (t *Table) WriteCSV(filename string) error
```

### Step 5 - main.go でフラグを組み合わせる

`flag` パッケージで各フラグを受け取り、処理をパイプラインのように繋げる。

```go
filter := flag.String("filter", "", "フィルタ条件 例: amount>=1000")
group  := flag.String("group", "", "グループ化するカラム名")
sum    := flag.String("sum", "", "合計するカラム名")
out    := flag.String("out", "", "出力ファイル名")
```

### Step 6 - テストを書く

各機能をテストで検証する。

```go
func TestFilter(t *testing.T)
func TestGroupSum(t *testing.T)
```

---

## 学びのポイント

- `encoding/csv` でのCSV読み書き
- 文字列パース（フィルタ条件の解析）
- `strconv` での型変換（文字列 → 数値）
- マップを使ったグループ集計
- `os.Args` / `flag` パッケージでのCLI引数処理
- テーブル構造のデータ操作
