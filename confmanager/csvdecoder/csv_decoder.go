package csvdecoder

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//Decode 解码
func Decode(in interface{}, filePath string) ([]interface{}, error) {
	fileIO, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := fileIO.Close()
		if err != nil {
			panic(fmt.Sprintf("defer csvdecoder Decodeerr:%s", err.Error()))
		}
	}()

	csvReader := csv.NewReader(fileIO)
	csvRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(csvRows) <= 3 {
		return nil, fmt.Errorf("decode err file:%s, rowNum:%d", filePath, len(csvRows))
	}

	//处理bom
	if len(csvRows[0]) > 0 {
		name := csvRows[0][0]
		if len(name) > 3 {
			if name[0] == 0xef || name[1] == 0xbb || name[2] == 0xbf {
				name = name[3:]
				csvRows[0][0] = name
			}
		}
	}

	headers := normalizeHeaders(csvRows[0])
	bodys := csvRows[3:]

	outValue := make([]interface{}, len(bodys))
	//err = ensureOutType(outValue.Type())
	//fmt.Println("out:", out, "in:", reflect.ValueOf(in).Elem(), "canaddr:", reflect.ValueOf(in).Elem().CanAddr(), "outtype:", err)

	oneType := reflect.ValueOf(in).Elem().Type()
	for rIdx, rowTemp := range bodys {
		objOne := reflect.New(oneType)
		objType := objOne.Elem().Type()

		fieldNum := objOne.Elem().NumField()
		for fIdx := 0; fIdx < fieldNum; fIdx++ {
			fieldName := objType.Field(fIdx).Name
			//fmt.Println("field value:", fieldName)
			isFind := false
			for hIdx, headerInfo := range headers {
				if fieldName == headerInfo {
					fieldVal := objOne.Elem().FieldByName(fieldName)
					err := setField(fieldVal, rowTemp[hIdx], true)
					if err != nil {
						return nil, err
					}
					isFind = true
				}
			}

			if !isFind {
				return nil, fmt.Errorf("not find csv field:%s filePath:%s", objOne.Elem().Type().Name(), filePath)
			}
		}

		//fmt.Println("object：", objOne, "Index:", rIdx)
		reflectedObject := reflect.ValueOf(objOne)
		reflect.ValueOf(outValue).Index(rIdx).Set(reflectedObject)
	}

	return outValue, nil
}

func normalizeHeaders(headers []string) []string {
	out := make([]string, len(headers))
	for i, h := range headers {
		out[i] = normalizeName(h)
	}
	return out
}

func normalizeName(name string) string {
	//将_x转换为X,q且首字母大写
	params := strings.Split(name, "_")
	if len(params) > 0 {
		suffix := params[len(params)-1]
		if suffix == "id" || suffix == "Id" {
			suffix = "ID"
		}
		params[len(params)-1] = suffix
	}

	var result string
	for _, paramInfo := range params {
		result += strings.ToUpper(paramInfo[0:1]) + paramInfo[1:]
	}

	if len(result) >= 2 {
		sufIdx := len(result) - 2
		suffix := result[sufIdx:]
		if suffix == "id" || suffix == "Id" {
			suffix = "ID"
		}

		result = result[0:sufIdx] + suffix
	}
	return result
}

func setField(field reflect.Value, value string, omitEmpty bool) error {
	if field.Kind() == reflect.Ptr {
		if omitEmpty && value == "" {
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Interface().(type) {
	case string:
		field.SetString(value)

	case int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))

	case []int:
		strVals := strings.Split(value, "|")
		intVals := make([]int, 0, len(strVals))
		for _, strTemp := range strVals {
			intVal, err := strconv.Atoi(strTemp)
			if err != nil {
				return err
			}
			intVals = append(intVals, intVal)
		}
		field.Set(reflect.ValueOf(intVals))

	case []string:
		strVal := strings.Split(value, "|")
		field.Set(reflect.ValueOf(strVal))

	case [][]int:
		intArr := strings.Split(value, ";")
		int2Val := make([][]int, 0, len(intArr))
		for _, intValTemp := range intArr {
			strVals := strings.Split(intValTemp, "|")
			intVals := make([]int, 0, len(strVals))
			for _, strTemp := range strVals {
				intVal, err := strconv.Atoi(strTemp)
				if err != nil {
					return err
				}
				intVals = append(intVals, intVal)
			}
			int2Val = append(int2Val, intVals)
		}
		field.Set(reflect.ValueOf(int2Val))
	}
	return nil
}
