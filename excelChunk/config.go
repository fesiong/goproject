package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

type configData struct {
	InputFile string `json:"input_file"`
	ChunkSize int `json:"chunk_size"`
}

var ExecPath string
var ConfigData configData

func init(){
	sep := string(os.PathSeparator)
	root := filepath.Dir(os.Args[0])
	ExecPath, _ = filepath.Abs(root)
	length := utf8.RuneCountInString(ExecPath)
	lastChar := ExecPath[length-1:]
	if lastChar != sep {
		ExecPath = ExecPath + sep
	}
	_, err := os.Stat(ExecPath + "output")
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(ExecPath+"output", os.ModePerm)
		if err != nil {
			fmt.Println("Can't create output directory: ", err.Error())
			os.Exit(-1)
		}
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%sconfig.json", ExecPath))
	if err != nil {
		fmt.Println("Read config error: ", err.Error())
		os.Exit(-1)
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

	if err := json.Unmarshal(bytes, &ConfigData); err != nil {
		fmt.Println("Invalid config: ", err.Error())
		os.Exit(-1)
	}
}