//Package csv2go Manager Gen
package csv2go

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const iCsvMgr = `// Package confmanager GENERATED BY CSV MANAGER AUTO; DO NOT EDIT
package confmanager

import (
    confcsv "github.com/joycastle/matching-story-robot-service/confmanager/csvauto"
    "fmt"
)

func (f *file2Struct) init() {
{{range .}}    f.File2Str["{{.FileName}}"] = &confcsv.{{.StructName }}{}{{"\n"}}{{end}}
}

//IConfManager 配置管理中心
type IConfManager interface {
{{range .}}
    GetConf{{.StructName }}Num() (int, error)
    GetConf{{.StructName }}ByKey(key {{.FirstColumnType}}) (confcsv.I{{.StructName}}, error)
    GetConf{{.StructName }}ByIndex(index int) (confcsv.I{{.StructName}}, error)
{{end}}
}


{{range .}}
//GetConf{{.StructName }}Num auto
func (c * ConfManager) GetConf{{.StructName }}Num() (int, error) {
    inters, ok := c.confMap["{{.FileName}}"]
    if !ok {
        return 0, fmt.Errorf("not find conf file name:%s", "{{.FileName}}")
    }
    return len(inters), nil
}

//GetConf{{.StructName }}ByKey auto
func (c * ConfManager) GetConf{{.StructName }}ByKey(key {{.FirstColumnType}}) (confcsv.I{{.StructName}}, error) {
    inters, ok := c.confMap["{{.FileName}}"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "{{.FileName}}")
    }

    for _, inter := range inters {
        obj := inter.(confcsv.I{{.StructName }})
        if obj.Get{{.FirstColumnName}}() == key {
            return obj, nil
        }
    }
    return nil, fmt.Errorf("conf not find {{.FileName}} file key:%v", key)
}

//GetConf{{.StructName }}ByIndex auto
func (c * ConfManager) GetConf{{.StructName }}ByIndex(index int) (confcsv.I{{.StructName}}, error) {
    inters, ok := c.confMap["{{.FileName}}"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "{{.FileName}}")
    }

    if len(inters) <= index {
        return nil, fmt.Errorf("conf {{.FileName}} index crash index:%d, len:%d", len(inters), index)
    }

    obj := inters[index].(confcsv.I{{.StructName }})
    return obj, nil
}
{{end}}
`

//PrepareConfStruct 模板参数结构
type PrepareConfStruct struct {
	FileName        string
	StructName      string
	FirstColumnName string
	FirstColumnType string
}

//PrepareConfStuctGenerator 模板文件生成器
type PrepareConfStuctGenerator struct {
	csvPath    string
	filePath   string
	dealFiles  []string //需要处理的文件名，不输入默认全选
	fiterFiles []string //需要过滤的文件名，不输入默认不过滤
}

//NewConfManagerGenerator 创建模板文件生成器
func NewConfManagerGenerator() *PrepareConfStuctGenerator {
	return &PrepareConfStuctGenerator{}
}

//SetCsvPath 设置csv路径
func (p *PrepareConfStuctGenerator) SetCsvPath(csvPath string) *PrepareConfStuctGenerator {
	p.csvPath = csvPath
	return p
}

//SetOutFile 设置输出文件路径
func (p *PrepareConfStuctGenerator) SetOutFile(filePath string) *PrepareConfStuctGenerator {
	p.filePath = filePath
	return p
}

//SetDealFiles 设置需要处理的文件
func (p *PrepareConfStuctGenerator) SetDealFiles(files []string) *PrepareConfStuctGenerator {
	p.dealFiles = files
	return p
}

//SetFiterFiles 设置需要过滤的文件
func (p *PrepareConfStuctGenerator) SetFiterFiles(files []string) *PrepareConfStuctGenerator {
	p.fiterFiles = files
	return p
}

//Run 执行
func (p *PrepareConfStuctGenerator) Run() {
	if len(p.filePath) <= 0 {
		panic("no out file")
	}
	if _, err := os.Stat(p.filePath); os.IsExist(err) {
		errR := os.Remove(p.filePath)
		if errR != nil {
			panic(fmt.Sprintf("create conf manager os.Remove filePath:%s, err:%s", p.filePath, errR))
		}
	}

	f, err := os.Create(p.filePath)
	if err != nil {
		panic(fmt.Sprintf("create conf manager os.Create filePath:%s, err:%s", p.filePath, err))
	}
	defer f.Close()

	args := p.readCsvInfo()
	t := template.Must(template.New("manager").Parse(iCsvMgr))
	if err := t.Execute(f, args); err != nil {
		panic(err.Error())
	}

	fmt.Println("succeed")
}

func (p *PrepareConfStuctGenerator) readCsvInfo() []*PrepareConfStruct {
	dealFiles := p.dealFiles

	fmt.Println("csvPath:", p.csvPath)
	//读取csv文件夹，获取所有的文件名与文件路径映射，并且将文件名去除后缀
	fileNameMap := make(map[string]string)
	if len(dealFiles) <= 0 {
		err := filepath.Walk(p.csvPath, func(path string, info os.FileInfo, errWalk error) error {
			if errWalk != nil {
				return errWalk
			}

			if !strings.Contains(path, ".csv") {
				return nil
			}

			_, fileName := filepath.Split(path)

			if isInArray(fileName, p.fiterFiles) {
				return nil
			}

			fileName = strings.TrimSuffix(fileName, ".csv")
			fileNameMap[fileName] = path
			return nil
		})

		if err != nil {
			panic(fmt.Sprintf("readCsvInfo filepath.Walk err:%s", err.Error()))
		}
	}

	//将csv文件转化为对应的CsvName2Columns 结构体
	csv2Prepare := make([]*PrepareConfStruct, 0, len(fileNameMap))
	for fileName, filePath := range fileNameMap {
		fileIO, err := os.Open(filePath)
		if err != nil {
			panic(fmt.Sprintf("readCsvInfo ioutil.ReadFile err:%s", err.Error()))
		}
		defer fileIO.Close()

		csvReader := csv.NewReader(fileIO)
		columns, err := csvReader.Read()
		if err != nil {
			panic(fmt.Sprintf("readCsvInfo columns csvReader.Read err:%s", err.Error()))
		}

		//处理bom
		if len(columns) > 0 {
			name := columns[0]
			if len(name) > 3 {
				if name[0] == 0xef || name[1] == 0xbb || name[2] == 0xbf {
					name = name[3:]
					columns[0] = name
				}
			}
		}

		if len(columns) < 2 {
			panic(fmt.Sprintf("readCsvInfo columns csvReader.Read err columnslen < 2 len:%d", len(columns)))
		}

		types, err := csvReader.Read()
		if err != nil {
			panic(fmt.Sprintf("readCsvInfo types csvReader.Read err:%s", err.Error()))
		}

		if len(columns) > len(types) {
			panic(fmt.Sprintf("readCsvInfo types less than columns typesLen:%d, columnsLen:%d", len(types), len(columns)))
		}

		typeFirst, err := dealTypeToGo(types[0])
		if err != nil {
			panic(fmt.Sprintf("readCsvInfo err type info file:%s err:%s", fileName, err.Error()))
		}

		csv2PreTemp := &PrepareConfStruct{
			FileName:        fileName,
			StructName:      dealName(fileName),
			FirstColumnName: dealName(columns[0]),
			FirstColumnType: typeFirst,
		}
		csv2Prepare = append(csv2Prepare, csv2PreTemp)
	}

	for _, csvInfo := range csv2Prepare {
		fmt.Println("prepare manager csv file to csvInfo:", csvInfo.FileName)
	}
	return csv2Prepare
}
