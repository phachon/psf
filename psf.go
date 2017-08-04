package main

import (
	"fmt"
	"path/filepath"
	"os"
	"flag"
	"io/ioutil"
	"strings"
	"encoding/base64"
)

var fileList = []string{}

var fileInfoMapping = map[string]string{}

func main()  {

	filePath := flag.String("path", "", "static file path, more than one , separated. ")
	psfResFile := flag.String("name", "psf_resouce", "save go file name")
	flag.Parse()

	filePathValues := strings.Split(*filePath, ",");
	if(len(filePathValues) == 0) {
		checkErr("static path is not bee empty!")
	}
	for _, filePathValue := range filePathValues {
		err := getFileList(filePathValue)
		if(err != nil) {
			checkErr("get file list error: " + err.Error())
		}
	}

	for _, file := range fileList {
		fileInfo, err := readFileInfo(file)
		if(err != nil) {
			checkErr("read file "+ file +"data error: " + err.Error())
		}
		fileInfoMapping[file] = fileInfo
	}

	createGoFile(*psfResFile + ".go")
}

// get file list
func getFileList(path string) (err error) {
	err = filepath.Walk(path, fileHandleFunc)
	return
}

// handle file list
func fileHandleFunc(path string, f os.FileInfo, err error) error {
	if(f == nil) {
		return err
	}
	if(f.IsDir()) {
		return nil
	}

	path = strings.Replace(path, "\\", "/", -1)

	fileList = append(fileList, path)
	return nil
}

//read file data
func readFileInfo(file string) (data string, err error) {
	var	BASE_64_TABLE = "<,./?~!@#$CDVWX%^&*ABYZabcghijkpqrstuvwxyz01EFKLMNOPQRSTU2345678"

	fileBytes, err := ioutil.ReadFile(file)
	var coder = base64.NewEncoding(BASE_64_TABLE)
	var src []byte = []byte(string(fileBytes))
	data = string([]byte(coder.EncodeToString(src)))
	if(err != nil) {
		return
	}
	return
}

//create file value
func createGoFile(file string)  {

	templateFile := `package main

var psfResValues = map[string]string{
#data#
}

func PsfRes(name string) string {
	const BASE_64_TABLE = "<,./?~!@#$CDVWX%^&*ABYZabcghijkpqrstuvwxyz01EFKLMNOPQRSTU2345678"
	str := psfResValues[name]
	coder := base64.NewEncoding(BASE_64_TABLE)
	by, _:= coder.DecodeString(str)
	return string(by)
}`

	var contents = ""
	for fileName, content := range fileInfoMapping {
		contents += `"`+fileName+`":` + `"` + content + `",`
		contents += "\n\n"
	}

	templateData := strings.Replace(templateFile, "#data#", contents, 1);

	err := ioutil.WriteFile(file, []byte(templateData), 0777)
	if(err != nil) {
		checkErr("write file error : " + err.Error())
	}
}

//check error
func checkErr(err string)  {
	fmt.Println(err + "\n")
	os.Exit(0);
}