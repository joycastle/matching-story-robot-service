package config

func Startup() error {
	if err := ReadRobotTeamConfigFromConfManager(); err != nil {
		return err
	}

	if err := ReadRobotTeamChatFromConfManager(); err != nil {
		return err
	}

	if err := ReadRobotTeamInitialFromConfManager(); err != nil {
		return err
	}

	if err := ReadRobotTeamFromConfManager(); err != nil {
		return err
	}

	return nil
}
