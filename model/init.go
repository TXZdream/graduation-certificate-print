package model

import (
	"errors"
	"log"
	"os"

	dbf "github.com/SebastiaanKlippert/go-foxpro-dbf"
	"github.com/tealeg/xlsx"
	"gopkg.in/ini.v1"
)

var (
	// DB 存储到数据库的连接
	DB *dbf.DBF
	// XLSXFile 存储到目标文件的连接
	XLSXFile *xlsx.File
)

var (
	dbFilePath     string
	targetFilePath string
)

func init() {
	// 读取ini配置
	conf, err := ini.Load("data/conf.ini")
	if err != nil {
		log.Println("打开配置文件失败： ", err)
	}
	dbFilePath = conf.Section("DB").Key("db").String()
	targetFilePath = conf.Section("DB").Key("target").String()
	// 设置数据库版本
	dbf.SetValidFileVersionFunc(func(version byte) error {
		if version == 0x03 {
			return nil
		}
		return errors.New("not 0x03")
	})
	// 设置编码格式
	var decoder dbf.Decoder
	switch conf.Section("Content").Key("encoding").String() {
	case "GB18030":
		decoder = &GB18030{}
	default:
		decoder = nil
	}
	if decoder == nil {
		log.Println("不支持的编码格式：", conf.Section("Content").Key("encoding").String())
		os.Exit(1)
	}
	// 打开数据库文件
	db, err := dbf.OpenFile(dbFilePath, decoder)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	DB = db
	// 打开转换后存储的文件
	if f, err := os.Open(targetFilePath); os.IsNotExist(err) {
		XLSXFile = xlsx.NewFile()
		return
	} else if err != nil {
		log.Println("无法打开xlsx文件：", err)
		os.Exit(1)
	} else {
		f.Close()
	}
	if XLSXFile, err = xlsx.OpenFile(targetFilePath); err != nil {
		log.Println("无法打开xlsx文件：", err)
		os.Exit(1)
	}
}
