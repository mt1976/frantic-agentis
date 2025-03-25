package userStore

import (
	"context"

	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

func InitialiseServiceUser(cfg *commonConfig.Settings) (messageHelpers.UserMessage, error) {

	serviceUserUID := cfg.GetServiceUser_UID()
	serviceUserName := cfg.GetServiceUser_Name()
	serviceUserCode := cfg.GetServiceUser_UserCode()

	if serviceUserUID == "" || serviceUserName == "" || serviceUserCode == "" {
		logHandler.ErrorLogger.Fatalf("[%v] Service User - Configuration Error", appName)
		panic("Service User - Configuration Error")
	}

	logHandler.SecurityLogger.Printf("[%v] System Service User - Generation Starting... '%v'", appName, serviceUserCode)
	//check for a system user
	u, err := GetBy(FIELD_UserCode, serviceUserCode)
	if err == nil {
		logHandler.SecurityLogger.Printf("[%v] System Service User - already exists 😀 '%v'", appName, serviceUserCode)
		return messageHelpers.UserMessage{Key: u.Key, Code: u.UserCode, Source: u.Source}, nil
	}
	logHandler.SecurityLogger.Printf("[%v] System Service User - does not exist, creating '%v'", appName, serviceUserCode)
	su, err := new(context.TODO(), serviceUserName, serviceUserUID, "System Service User "+serviceUserName, serviceUserName+"@"+application.SystemIdentity(), serviceUserUID, application.SystemIdentity(), true, true, false, "")
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] ERROR Unable to Create System Service User [%v]", appName, err.Error())
		return messageHelpers.UserMessage{}, err
	}

	logHandler.SecurityLogger.Printf("[%v] System Service User - Created '%v'", appName, serviceUserCode)
	if su.Source == "" {
		// Set the source to the application name if its not been set
		su.Source = cfg.GetApplication_Name()
	}
	return messageHelpers.UserMessage{Key: su.Key, Code: su.UserCode, Source: su.Source}, nil
}
