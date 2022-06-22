package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joycastle/matching-story-robot-service/confmanager/tools/csv2go"
)

func main() {
	packageStr := flag.String("package", "csvauto", "csv package name")
	templateStr := flag.String("template", "../template", "csv get path")
	casAutoPath := flag.String("csvsave", "../csvauto", "csv save path")
	managerFileStr := flag.String("managerfile", "../conf_manager_auto.go", "manager save file")
	flag.Parse()
	fmt.Println(os.Getwd())
	csvReader := csv2go.NewCsv2Struct()
	csvReader.SetPackageName(*packageStr)
	csvReader.SetCsvPath(*templateStr)
	csvReader.SetSavePath(*casAutoPath)
	csvReader.Run()
	fmt.Println("succeed : csv to go struct")

	confMgr := csv2go.NewConfManagerGenerator()
	confMgr.SetCsvPath(*templateStr)
	confMgr.SetOutFile(*managerFileStr)
	confMgr.Run()
	fmt.Println("succeed : csv to conf manager")
}
