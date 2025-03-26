package userIdentity

import (
	"github.com/mt1976/frantic-agentis/app/dao/userStore"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

var appName = "frantic-agentis"

func Validate(userKey string) (messageHelpers.UserMessage, error) {

	logHandler.SecurityLogger.Printf("[%v] Validate: UserID=[%v]", appName, userKey)

	user, err := userStore.GetByKey(userKey)
	if err != nil {
		logHandler.SecurityLogger.Printf("[%v] ERROR: UserKey=[%v] NO USER FOUND", appName, userKey)
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		return messageHelpers.UserMessage{}, commonErrors.ErrorUserNotFound
	}

	if user.Audit.DeletedBy != "" {
		logHandler.SecurityLogger.Printf("[%v] ERROR: UserKey=[%v] UserCode=[%v] DELETED USER", appName, userKey, user.UserCode)
		logHandler.ErrorLogger.Printf("Error=[%v]", "User Deleted")
		return messageHelpers.UserMessage{}, commonErrors.ErrorUserNotActive
	}

	if !user.IsActive.Bool() {
		logHandler.SecurityLogger.Printf("[%v] ERROR: UserKey=[%v] UserCode=[%v] INACTIVE USER", appName, userKey, user.UserCode)
		logHandler.ErrorLogger.Printf("Error=[%v]", "User Not Active")
		return messageHelpers.UserMessage{}, commonErrors.ErrorUserNotActive
	}

	logHandler.SecurityLogger.Printf("[%v] User Validated!: UserKey=[%v] UserCode=[%v]", appName, userKey, user.UserCode)
	x := messageHelpers.UserMessage{}
	x.Request(user.Key, user.UserCode, user.Source, user.Locale, user.Theme, user.Timezone, user.Role)
	return x, nil
}

func ValidateUserName(userName string) (messageHelpers.UserMessage, error) {
	//	logHandler.SecurityLogger.Printf("[%v] ValidateUserName: UserName=[%v]", domain, userName)

	user, err := userStore.GetByUserName(userName)
	if err != nil {
		logHandler.SecurityLogger.Printf("[%v] Error: UserName=[%v] NO USER FOUND", appName, userName)
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		return messageHelpers.UserMessage{}, commonErrors.ErrorUserNotFound
	}

	x := messageHelpers.UserMessage{}
	x.Request(user.Key, user.UserCode, user.Source, user.Locale, user.Theme, user.Timezone, user.Role)
	return x, nil
}
