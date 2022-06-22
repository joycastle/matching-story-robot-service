package confmanager

import (
    "encoding/json"
    "fmt"
    "testing"
)

func TestConfManager(t *testing.T) {
    confMgr := GetConfManagerVer()
    err := confMgr.LoadCsv("template")
    if err != nil {
        t.Error(err.Error())
        return
    }

    verMgr, err := confMgr.GetConfManager()
    if err != nil {
        t.Error(err.Error())
        return
    }

    num, err := verMgr.GetConfRobotNum()
    if err != nil {
        t.Error(err.Error())
        return
    }
    fmt.Println("num:", num)


    conf, err := verMgr.GetConfRobotByIndex(0)
    if err != nil {
        t.Error(err.Error())
        return
    }

    confStr, _ := json.Marshal(conf)
    fmt.Println("index:", string(confStr))

    conf, err = verMgr.GetConfRobotByKey(1001)
    if err != nil {
        t.Error(err.Error())
        return
    }

    confStr, _ = json.Marshal(conf)
    fmt.Println("key:",string(confStr))
}