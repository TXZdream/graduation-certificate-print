package main

import (
	"log"
	"os"

	"github.com/txzdream/graduation-certificate-print/model"
)

func main() {
	log.Println("读取源文件")
	records, err := model.ReadAll()
	if err != nil {
		os.Exit(1)
	}
	log.Println("正在写入结果")
	if err = model.WriteResult(records); err != nil {
		os.Exit(1)
	}
}
