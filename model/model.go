package model

import (
	"log"
	"strconv"
	"strings"

	"github.com/goinggo/mapstructure"
)

// GraduationData 存储了用以打印的字段信息
type GraduationData struct {
	Name                string `jpath:"XM"`
	Gender              string `jpath:"XB"`
	Birthday            string `jpath:"CSRQ"`
	AfterBirthday       Date
	EnrollmentDate      string `jpath:"RXRQ"`
	AfterEnrollmentDate Date
	GraduationDay       string `jpath:"BYRQ"`
	AfterGraduationDay  Date
	Subject             string `jpath:"ZYH"`
	Form                string `jpath:"XXXS"`
	YearLength          string `jpath:"XZ"`
	Type                string `jpath:"PYCC"`
	ID                  string `jpath:"ZSBH"`
}

// Date 存储大写表示的内容
type Date struct {
	Year  string
	Month string
	Day   string
}

// Result 包含转换后的结果
type Result struct {
	Name            string
	Gender          string
	BirthYear       string
	BirthMonth      string
	BirthDay        string
	EnrollmentYear  string
	EnrollmentMonth string
	GraduationYear  string
	GraduationMonth string
	Subject         string
	Form            string
	Type            string
	YearLength      string
	ID              string
}

// ReadAll 从数据库中读取所有记录
func ReadAll() ([]GraduationData, error) {
	var records []GraduationData
	var i uint32
	for i = 0; i < DB.NumRecords(); i++ {
		record, err := DB.RecordToMap(i)
		if err != nil {
			log.Println("读取位于", i, "行的记录发生错误：", err)
			return nil, err
		}
		var item GraduationData
		if err := mapstructure.DecodePath(record, &item); err != nil {
			log.Println("解析位于", i, "行的记录发生错误：", err)
			return nil, err
		}
		Transform(&item)
		records = append(records, item)
	}
	return records, nil
}

// Transform 将存储的日期变成大写
func Transform(g *GraduationData) {
	g.Birthday = strings.ReplaceAll(g.Birthday, " ", "")
	g.AfterBirthday.Year = TransformYear(g.Birthday[0:4])
	g.AfterBirthday.Month = TransformMonthOrDate(g.Birthday[4:6])
	g.AfterBirthday.Day = TransformMonthOrDate(g.Birthday[6:8])

	g.EnrollmentDate = strings.ReplaceAll(g.EnrollmentDate, " ", "")
	g.AfterEnrollmentDate.Year = TransformYear(g.EnrollmentDate[0:4])
	g.AfterEnrollmentDate.Month = TransformMonthOrDate(g.EnrollmentDate[4:6])
	g.AfterEnrollmentDate.Day = TransformMonthOrDate(g.EnrollmentDate[6:8])

	g.GraduationDay = strings.ReplaceAll(g.GraduationDay, " ", "")
	g.AfterGraduationDay.Year = TransformYear(g.GraduationDay[0:4])
	g.AfterGraduationDay.Month = TransformMonthOrDate(g.GraduationDay[4:6])
	g.AfterGraduationDay.Day = TransformMonthOrDate(g.GraduationDay[6:8])
}

// TransformYear 将年份转成大写
func TransformYear(raw string) string {
	var ret string
	for _, v := range raw {
		tmp := NumberMapping[v]
		if tmp == "" {
			log.Println("非法的数字：", tmp)
			return ""
		}
		ret += tmp
	}
	return ret
}

// TransformMonthOrDate 将月份和日期转成大写
func TransformMonthOrDate(raw string) string {
	data, err := strconv.Atoi(raw)
	if err != nil {
		log.Println("非法的月份或日期：", raw, err)
	}
	if data <= 0 || data >= 31 {
		log.Println("非法的月份或日期：", raw)
	}
	var ret string
	if data%10 != 0 {
		ret += NumberMapping[rune(strconv.Itoa(data % 10)[0])]
	}
	if data/10 == 0 {
		return ret
	}
	ret = "十" + ret
	if data/10 != 1 {
		ret = NumberMapping[rune(strconv.Itoa(data / 10)[0])] + ret
	}
	return ret
}

// WriteResult 将结果写到磁盘上
func WriteResult(gl []GraduationData) error {
	sheet := XLSXFile.Sheet["毕业生信息"]
	var err error
	if sheet == nil {
		sheet, err = XLSXFile.AddSheet("毕业生信息")
		if err != nil {
			log.Println("无法添加sheet：", err)
			return err
		}
		sheet.AddRow().WriteSlice(&Header, 14)
	}
	for _, v := range gl {
		tmp := Result{
			Name:            v.Name,
			Gender:          v.Gender,
			BirthYear:       v.AfterBirthday.Year,
			BirthMonth:      v.AfterBirthday.Month,
			BirthDay:        v.AfterBirthday.Day,
			EnrollmentYear:  v.AfterEnrollmentDate.Year,
			EnrollmentMonth: v.AfterEnrollmentDate.Month,
			GraduationYear:  v.AfterGraduationDay.Year,
			GraduationMonth: v.AfterGraduationDay.Month,
			Subject:         v.Subject,
			Form:            v.Form,
			Type:            v.Type,
			YearLength:      v.YearLength,
			ID:              v.ID,
		}
		sheet.AddRow().WriteStruct(&tmp, 14)
	}
	if err = XLSXFile.Save(targetFilePath); err != nil {
		log.Println("无法存储xlsx文件", err)
		return err
	}
	return nil
}
