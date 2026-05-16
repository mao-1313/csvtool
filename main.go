package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	// オプションの処理
	outFile := flag.String("out", "", "Output CSV file path")
	filterOpt := flag.String("filter", "", " Filter condition (e.g. 'amount>=1000,date=20240101')")
	groupOpt := flag.String("group", "", " group condition (e.g. 'date,store,item')")
	sumOpt := flag.String("sum", "", " sum condition (e.g. 'amount,qy')")

	flag.Parse()

	// 引数処理
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: csvtool [options] input.csv")
		os.Exit(1)
	}

	inputFile := args[0]
	// CSVファイルの読み込み
	table, err := ReadCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	// フィルタリング
	if *filterOpt != "" {
		table, err = table.Filter(*filterOpt)
		if err != nil {
			fmt.Printf("Error filtering data: %v\n", err)
			os.Exit(1)
		}
	}

	// 集計
	if *groupOpt != "" && *sumOpt != "" {
		table, err = table.GroupAndSum(*groupOpt, *sumOpt)
		if err != nil {
			fmt.Printf("Error grouping and summing data: %v\n", err)
			os.Exit(1)
		}
	} else if *groupOpt != "" || *sumOpt != "" {
		fmt.Println("Both -group and -sum options must be specified together")
		os.Exit(1)
	}

	// 出力
	if *outFile != "" {
		err = table.WriteCSV(*outFile)
		if err != nil {
			fmt.Printf("Error writing CSV: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 標準出力にタブ区切りで出力
		err = table.WriteTSV(os.Stdout)
		if err != nil {
			fmt.Printf("Error writing TSV: %v\n", err)
			os.Exit(1)
		}
	}
}
