package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"path"
)

var CellCols = []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}

func ChunkExcel() {
	f, err := excelize.OpenFile(fmt.Sprintf("%s%s", ExecPath, ConfigData.InputFile))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("---------- 开始生成 ----------")

	var titles []string
	excelFiles := map[string]*excelize.File{}
	fileName := fmt.Sprintf("output%s%d-%d.xlsx", string(os.PathSeparator), 1, ConfigData.ChunkSize)
	excelFiles[fileName] = excelize.NewFile()
	firstSheet := excelFiles[fileName].NewSheet("Sheet1")
	excelFiles[fileName].SetActiveSheet(firstSheet)
	sheets := f.GetSheetMap()
	for _, sheetName := range sheets {
		rows := f.GetRows(sheetName)
		rowIndex := 2
		numIndex := 0
		for i, row := range rows {
			if i == 0 {
				//标题行
				for j, colCell := range row {
					titles = append(titles, colCell)
					//写入表头
					excelFiles[fileName].SetCellValue("Sheet1", fmt.Sprintf("%s1", CellCols[j]), colCell)
				}
			} else {
				//开始写入文件
				for j, colCell := range row {
					excelFiles[fileName].SetCellValue("Sheet1", fmt.Sprintf("%s%d", CellCols[j], rowIndex), colCell)
				}
				numIndex++
				if numIndex % ConfigData.ChunkSize == 0 {
					//保存
					fmt.Println(fmt.Sprintf("开始生成文件%s", path.Base(fileName)))
					if err := excelFiles[fileName].SaveAs(fileName); err != nil {
						fmt.Println(err)
						return
					}
					//生成新文件
					fileName = fmt.Sprintf("output%s%d-%d.xlsx", string(os.PathSeparator), numIndex+1, numIndex+ConfigData.ChunkSize)
					excelFiles[fileName] = excelize.NewFile()
					firstSheet = excelFiles[fileName].NewSheet("Sheet1")
					excelFiles[fileName].SetActiveSheet(firstSheet)
					rowIndex = 1
					for j, colCell := range titles {
						//写入表头
						excelFiles[fileName].SetCellValue("Sheet1", fmt.Sprintf("%s1", CellCols[j]), colCell)
					}
				}

				rowIndex++
			}
		}
		//最后的保存
		fmt.Println(fmt.Sprintf("开始生成文件%s", path.Base(fileName)))
		if err := excelFiles[fileName].SaveAs(fileName); err != nil {
			fmt.Println(err)
			return
		}
		break
	}

	fmt.Println("---------- 生成结束 ----------")
}