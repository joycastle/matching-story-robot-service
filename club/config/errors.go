package config

import "fmt"

//confmanager error
func errConfManangerInit(confType string, err error) error {
	return fmt.Errorf("config: confmanager init error: %s, conf_type:%s", err.Error(), confType)
}

func errConfManangerRead(confType string, err error) error {
	return fmt.Errorf("config: confmanager read error: %s, conf_type:%s", err.Error(), confType)
}

//line error
func errLineNumEmpty(confType string) error {
	return fmt.Errorf("config: line number is empty, conf_type:%s", confType)
}

func errLineNumLimit(confType string, limit int) error {
	return fmt.Errorf("config: line number only need %d lines, conf_type:%s", limit, confType)
}

//data error
func errDataArrayNumNotMatch(confType string, id int, args ...string) error {
	return fmt.Errorf("config: data num not match, id:%d, fileds:[%v], conf_type:%s", id, confType, args)
}

func errDataArrayNumLimit(confType string, id int, filed string, limit int) error {
	return fmt.Errorf("config: data only need %d numbers, id:%d, fileds:[%v], conf_type:%s", limit, id, filed, confType)
}

func errDataOnlyNeedLimit(confType string, id int, filed string, limit1, limit2 int) error {
	return fmt.Errorf("config: data only need %d or %d numbers, id:%d, fileds:[%v], conf_type:%s", limit1, limit2, id, filed, confType)
}

func errDataResultNotMatch(confType string, args ...string) error {
	return fmt.Errorf("config: data result not match in fileds[%v], conf_type:%s", args, confType)
}

func errDataResultEmpty(confType string, typ string) error {
	return fmt.Errorf("config: data contents is empty in %s, conf_type:%s", typ, confType)
}

//data parse
func errParseIndexNotFound(confType string, inmap string, index string) error {
	return fmt.Errorf("config: index %d not found in %s, conf_type:%s", index, inmap, confType)
}

func errParseIndexLimit(confType string, index string, inmap string, limit int) error {
	return fmt.Errorf("config: index %d only need %s number in %s, conf_type:%s", index, limit, inmap, confType)
}
