// Package confmanager GENERATED BY CSV MANAGER AUTO; DO NOT EDIT
package confmanager

import (
    confcsv "github.com/joycastle/matching-story-robot-service/confmanager/csvauto"
    "fmt"
)

func (f *file2Struct) init() {
    f.File2Str["robot_name"] = &confcsv.RobotName{}
    f.File2Str["robot_team"] = &confcsv.RobotTeam{}
    f.File2Str["robot_team_chat"] = &confcsv.RobotTeamChat{}
    f.File2Str["robot_team_config"] = &confcsv.RobotTeamConfig{}
    f.File2Str["robot_team_initial"] = &confcsv.RobotTeamInitial{}

}

//IConfManager 配置管理中心
type IConfManager interface {

    GetConfRobotNameNum() (int, error)
    GetConfRobotNameByKey(key int) (confcsv.IRobotName, error)
    GetConfRobotNameByIndex(index int) (confcsv.IRobotName, error)

    GetConfRobotTeamNum() (int, error)
    GetConfRobotTeamByKey(key int) (confcsv.IRobotTeam, error)
    GetConfRobotTeamByIndex(index int) (confcsv.IRobotTeam, error)

    GetConfRobotTeamChatNum() (int, error)
    GetConfRobotTeamChatByKey(key int) (confcsv.IRobotTeamChat, error)
    GetConfRobotTeamChatByIndex(index int) (confcsv.IRobotTeamChat, error)

    GetConfRobotTeamConfigNum() (int, error)
    GetConfRobotTeamConfigByKey(key int) (confcsv.IRobotTeamConfig, error)
    GetConfRobotTeamConfigByIndex(index int) (confcsv.IRobotTeamConfig, error)

    GetConfRobotTeamInitialNum() (int, error)
    GetConfRobotTeamInitialByKey(key int) (confcsv.IRobotTeamInitial, error)
    GetConfRobotTeamInitialByIndex(index int) (confcsv.IRobotTeamInitial, error)

}



//GetConfRobotNameNum auto
func (c * ConfManager) GetConfRobotNameNum() (int, error) {
    inters, ok := c.confMap["robot_name"]
    if !ok {
        return 0, fmt.Errorf("not find conf file name:%s", "robot_name")
    }
    return len(inters), nil
}

//GetConfRobotNameByKey auto
func (c * ConfManager) GetConfRobotNameByKey(key int) (confcsv.IRobotName, error) {
    inters, ok := c.confMap["robot_name"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_name")
    }

    for _, inter := range inters {
        obj := inter.(confcsv.IRobotName)
        if obj.GetID() == key {
            return obj, nil
        }
    }
    return nil, fmt.Errorf("conf not find robot_name file key:%v", key)
}

//GetConfRobotNameByIndex auto
func (c * ConfManager) GetConfRobotNameByIndex(index int) (confcsv.IRobotName, error) {
    inters, ok := c.confMap["robot_name"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_name")
    }

    if len(inters) <= index {
        return nil, fmt.Errorf("conf robot_name index crash index:%d, len:%d", len(inters), index)
    }

    obj := inters[index].(confcsv.IRobotName)
    return obj, nil
}

//GetConfRobotTeamNum auto
func (c * ConfManager) GetConfRobotTeamNum() (int, error) {
    inters, ok := c.confMap["robot_team"]
    if !ok {
        return 0, fmt.Errorf("not find conf file name:%s", "robot_team")
    }
    return len(inters), nil
}

//GetConfRobotTeamByKey auto
func (c * ConfManager) GetConfRobotTeamByKey(key int) (confcsv.IRobotTeam, error) {
    inters, ok := c.confMap["robot_team"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team")
    }

    for _, inter := range inters {
        obj := inter.(confcsv.IRobotTeam)
        if obj.GetID() == key {
            return obj, nil
        }
    }
    return nil, fmt.Errorf("conf not find robot_team file key:%v", key)
}

//GetConfRobotTeamByIndex auto
func (c * ConfManager) GetConfRobotTeamByIndex(index int) (confcsv.IRobotTeam, error) {
    inters, ok := c.confMap["robot_team"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team")
    }

    if len(inters) <= index {
        return nil, fmt.Errorf("conf robot_team index crash index:%d, len:%d", len(inters), index)
    }

    obj := inters[index].(confcsv.IRobotTeam)
    return obj, nil
}

//GetConfRobotTeamChatNum auto
func (c * ConfManager) GetConfRobotTeamChatNum() (int, error) {
    inters, ok := c.confMap["robot_team_chat"]
    if !ok {
        return 0, fmt.Errorf("not find conf file name:%s", "robot_team_chat")
    }
    return len(inters), nil
}

//GetConfRobotTeamChatByKey auto
func (c * ConfManager) GetConfRobotTeamChatByKey(key int) (confcsv.IRobotTeamChat, error) {
    inters, ok := c.confMap["robot_team_chat"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team_chat")
    }

    for _, inter := range inters {
        obj := inter.(confcsv.IRobotTeamChat)
        if obj.GetID() == key {
            return obj, nil
        }
    }
    return nil, fmt.Errorf("conf not find robot_team_chat file key:%v", key)
}

//GetConfRobotTeamChatByIndex auto
func (c * ConfManager) GetConfRobotTeamChatByIndex(index int) (confcsv.IRobotTeamChat, error) {
    inters, ok := c.confMap["robot_team_chat"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team_chat")
    }

    if len(inters) <= index {
        return nil, fmt.Errorf("conf robot_team_chat index crash index:%d, len:%d", len(inters), index)
    }

    obj := inters[index].(confcsv.IRobotTeamChat)
    return obj, nil
}

//GetConfRobotTeamConfigNum auto
func (c * ConfManager) GetConfRobotTeamConfigNum() (int, error) {
    inters, ok := c.confMap["robot_team_config"]
    if !ok {
        return 0, fmt.Errorf("not find conf file name:%s", "robot_team_config")
    }
    return len(inters), nil
}

//GetConfRobotTeamConfigByKey auto
func (c * ConfManager) GetConfRobotTeamConfigByKey(key int) (confcsv.IRobotTeamConfig, error) {
    inters, ok := c.confMap["robot_team_config"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team_config")
    }

    for _, inter := range inters {
        obj := inter.(confcsv.IRobotTeamConfig)
        if obj.GetID() == key {
            return obj, nil
        }
    }
    return nil, fmt.Errorf("conf not find robot_team_config file key:%v", key)
}

//GetConfRobotTeamConfigByIndex auto
func (c * ConfManager) GetConfRobotTeamConfigByIndex(index int) (confcsv.IRobotTeamConfig, error) {
    inters, ok := c.confMap["robot_team_config"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team_config")
    }

    if len(inters) <= index {
        return nil, fmt.Errorf("conf robot_team_config index crash index:%d, len:%d", len(inters), index)
    }

    obj := inters[index].(confcsv.IRobotTeamConfig)
    return obj, nil
}

//GetConfRobotTeamInitialNum auto
func (c * ConfManager) GetConfRobotTeamInitialNum() (int, error) {
    inters, ok := c.confMap["robot_team_initial"]
    if !ok {
        return 0, fmt.Errorf("not find conf file name:%s", "robot_team_initial")
    }
    return len(inters), nil
}

//GetConfRobotTeamInitialByKey auto
func (c * ConfManager) GetConfRobotTeamInitialByKey(key int) (confcsv.IRobotTeamInitial, error) {
    inters, ok := c.confMap["robot_team_initial"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team_initial")
    }

    for _, inter := range inters {
        obj := inter.(confcsv.IRobotTeamInitial)
        if obj.GetID() == key {
            return obj, nil
        }
    }
    return nil, fmt.Errorf("conf not find robot_team_initial file key:%v", key)
}

//GetConfRobotTeamInitialByIndex auto
func (c * ConfManager) GetConfRobotTeamInitialByIndex(index int) (confcsv.IRobotTeamInitial, error) {
    inters, ok := c.confMap["robot_team_initial"]
    if !ok {
        return nil, fmt.Errorf("not find conf file name:%s", "robot_team_initial")
    }

    if len(inters) <= index {
        return nil, fmt.Errorf("conf robot_team_initial index crash index:%d, len:%d", len(inters), index)
    }

    obj := inters[index].(confcsv.IRobotTeamInitial)
    return obj, nil
}

