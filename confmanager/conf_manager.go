package confmanager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/joycastle/matching-story-robot-service/confmanager/csvdecoder"
)

// ConfAllManager confmanager.ConfAllManager define version 配置管理器
type ConfAllManager struct {
	sync.RWMutex
	confVerMap map[int]*ConfManager
	ver        int
}

var gConfAllMgr = &ConfAllManager{confVerMap: make(map[int]*ConfManager)}

// GetConfManagerVer 初始化配置管理器
func GetConfManagerVer() *ConfAllManager {
	return gConfAllMgr
}

func (c *ConfAllManager) getNextVersion() int {
	c.ver++
	return c.ver
}

// GetConfManager 获取管理器实例
func (c *ConfAllManager) GetConfManager() (*ConfManager, error) {
	c.RLock()
	defer c.RUnlock()

	confMgr, ok := c.confVerMap[c.ver]
	if !ok {
		return nil, fmt.Errorf("conf manager version not find the version type ver:%d", c.ver)
	}
	return confMgr, nil
}

// LoadCsv 加载CSV
func (c *ConfAllManager) LoadCsv(csvPath string) error {
	c.Lock()
	defer c.Unlock()

	confMgr := &ConfManager{confMap: make(map[string][]interface{})}
	err := confMgr.loadCsv(csvPath)
	if err != nil {
		return err
	}

	ver := c.getNextVersion()
	c.confVerMap[ver] = confMgr

	// 加载新的版本是删除过期不再使用的版本，允许内存存在两个版本的配置
	for key, val := range c.confVerMap {
		if val.isExpired() {
			delete(c.confVerMap, key)
		}
	}
	return nil
}

type file2Struct struct {
	File2Str map[string]interface{}
}

func (f *file2Struct) getFileStruct(fileName string) (interface{}, error) {
	inter, ok := f.File2Str[fileName]
	if !ok {
		return nil, fmt.Errorf("not find csv struct maybe need to gen code file:%s", fileName)
	}
	return inter, nil
}

// 版本更新60s后视为过期，可以被删除
const expireTime = 60

// ConfManager define 单一版本配置管理器
type ConfManager struct {
	confMap    map[string][]interface{}
	expireTime int64
}

func (c *ConfManager) isExpired() bool {
	timeNow := time.Now().Unix()
	if c.expireTime > 0 && c.expireTime < timeNow {
		return true
	}
	return false
}

func (c *ConfManager) loadCsv(csvPath string) error {
	file2Str := &file2Struct{File2Str: make(map[string]interface{})}
	file2Str.init()

	// 读取csv文件夹，获取所有的文件名与文件路径映射，并且将文件名去除后缀
	struct2FilePath := make(map[string]string)
	err := filepath.Walk(csvPath, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			return errWalk
		}

		if !strings.Contains(path, ".csv") {
			return nil
		}

		_, fileName := filepath.Split(path)
		structName := strings.TrimSuffix(fileName, ".csv")
		struct2FilePath[structName] = path
		return nil
	})

	if err != nil {
		return fmt.Errorf("ConfManager LoadCsv filepath.Walk err:%s", err.Error())
	}

	for structName, filePath := range struct2FilePath {
		inter, ok := file2Str.File2Str[structName]
		if !ok {
			log.Printf("not find file struct file:%s", structName)
			continue
		}

		err := c.loadFile(inter, structName, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfManager) loadFile(obj interface{}, structName, filePath string) error {
	confVals, err := csvdecoder.Decode(obj, filePath)
	if err != nil {
		return err
	}

	confs := make([]interface{}, 0, len(confVals))
	for _, conf := range confVals {
		confVal := conf.(reflect.Value)
		confInter := confVal.Interface()
		confs = append(confs, confInter)
	}

	// str, _ := json.Marshal(confs)
	// fmt.Println("confs:", string(str))
	c.confMap[structName] = confs
	return nil
}
