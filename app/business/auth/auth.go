package auth

import (
	"context"
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/mt1976/frantic-agentis/app/business/translation"
	"github.com/mt1976/frantic-agentis/app/dao/passwordStore"
	"github.com/mt1976/frantic-agentis/app/dao/userStore"
	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
)

var cfg commonConfig.Settings
var appName string

func init() {
	cfg = *commonConfig.Get()
	appName = "frantic-agentis"
}

func ValidateUserName(userName string) bool {
	// Validate the user name
	// If the user name is valid, return nil
	// If the user name is invalid, return an error
	_, err := userStore.GetByUserName(userName)
	return err == nil
}

func ValidateUserNameAndPassword(userName, enteredPassword string) error {
	// Validate the password
	// If the password is valid, return nil
	// If the password is invalid, return an error
	user, err := userStore.GetByUserName(userName)
	if err != nil {
		return err
	}

	return ValidateUserIDAndPassword(user.Key, enteredPassword)
}

func encode(str string) string {
	// Encode the string
	// Return the encoded string
	return idHelpers.Encode(str)
}

func SetPassword(ctx context.Context, userKey, inPassword string) error {

	logHandler.SecurityLogger.Printf("Setting password for user [key=%v]", userKey)
	// Set the password
	// If the password is set, return nil
	// If the password is not set, return an error
	user, err := userStore.GetByKey(userKey)
	if err != nil {
		return err
	}

	logHandler.SecurityLogger.Printf("Setting password for user [key=%v] (%v)", user.Key, user.UserCode)

	//id := user.ID
	//idstring := encode(strconv.Itoa(id))
	p, err := passwordStore.GetByKey(userKey)
	if err != nil {
		// Create a new password object
		_, err := passwordStore.New(ctx, userKey, inPassword)
		if err != nil {
			return err
		}
		return nil
	}

	p.Password = encode(inPassword)
	err = p.Update(ctx, "Updated Password")
	if err != nil {
		return err
	}
	return nil
}

func ValidateUserIDAndPassword(userKey string, password string) error {
	// Validate the password

	record, err := passwordStore.GetByKey(userKey)
	if err != nil {
		logHandler.SecurityLogger.Printf("Error getting password: %v", err)
		return err
	}

	encodedPassword := encode(password)

	//fmt.Printf("record.Password: [%v] (%v)\n", record.Password, len(record.Password))
	logHandler.SecurityLogger.Printf("stored  password hash: [%v] (%v)", record.Password, len(record.Password))
	//fmt.Printf("encodedPassword: [%v] (%v)\n", encodedPassword, len(encodedPassword))
	logHandler.SecurityLogger.Printf("entered password hash: [%v] (%v)", encodedPassword, len(encodedPassword))

	if record.Password != encodedPassword {
		logHandler.SecurityLogger.Printf("Password does not match")
		return commonErrors.ErrorPasswordMismatch
	}

	logHandler.SecurityLogger.Printf("Password matches")

	return nil
}

func LoginCurrentOSUser(ctx context.Context) userStore.User_Store {

	logHandler.InfoLogger.Printf("[%v] Attempting to Login Current OS User", appName)

	temp := buildUserStub()

	usr, err := userStore.GetByUserCode(temp.UserCode)
	if err != nil {
		logHandler.WarningLogger.Printf("[%v] Warning=[%v] User=[%v]", appName, err.Error(), temp.UserCode)
		logHandler.InfoLogger.Printf("[%v] User=[%v] does not exist, creating", appName, temp.UserCode)
		usr, err = AddCurrentOSUser(ctx)
		if err != nil {
			logHandler.ErrorLogger.Printf("[%v] Warning=[%v] User=[%v]", appName, err.Error(), temp.UserCode)
			panic(err)
		}
	}

	usr.LastLogin = time.Now()
	usr.LastHost = application.HostName()

	//support.SetActiveUserInformation(usr.UID, usr.RealName, usr.UserName, usr.UserCode) // Adds the current user to the application context

	//u.Dump(!,"login")
	err = usr.Update(ctx, "login")
	if err != nil {
		logHandler.WarningLogger.Printf("[%v] Warning=[%v] User=[%v]", appName, err.Error(), usr.UserName)
		panic(err)
	}
	logHandler.EventLogger.Printf("[%v] User=[%v] Logged In [%v]", appName, usr.UserCode, usr.UserName)

	//setupBehaviours()
	//usr.Spew()
	return usr

}

func AddCurrentOSUser(ctx context.Context) (userStore.User_Store, error) {
	// Create a new User object
	// Check if the user already exists, if not create
	testu := buildUserStub()

	oldu, dupErr := testu.IsDuplicateOf() // Check if the user already exists
	if dupErr == commonErrors.ErrorDuplicate {
		logHandler.WarningLogger.Printf(translation.Get("[%v ]A user already exists for [%v]"), appName, testu.UserName)
		return oldu, nil
	}

	u, err := userStore.New(ctx, testu.UserName, testu.UID, testu.RealName, testu.Email, testu.GID, testu.Host, true, false)
	if err != nil {
		logHandler.ErrorLogger.Printf(translation.Get("[%v] Error=[%v]"), appName, err.Error())
		return userStore.User_Store{}, err
	}

	// Return the new User object
	return u, nil
}

func buildUserStub() userStore.User_Store {
	currentUser := GetOSUserDetails()

	stub := userStore.User_Store{}
	stub.ID = 0
	stub.UID = currentUser.Uid
	userName := currentUser.Username
	if application.IsRunningOnWindows() {
		// Take everthing after the \ in the username
		tkID := strings.Split(currentUser.Username, "\\")
		userName = tkID[1]
	}
	stub.UserName = userName
	stub.RealName = currentUser.Name
	stub.GID = currentUser.Gid
	stub.Email = strings.ToLower(fmt.Sprintf("%v@%v.com", userName, application.HostName()))
	stub.UserCode = strings.ToLower(fmt.Sprintf("%v_%s", currentUser.Uid, userName))
	stub.Notes = "Auto Generated User"
	stub.IsActive.Set(true)
	stub.LastLogin = time.Now()
	stub.LastHost = application.HostName()
	stub.Host = application.HostName()
	stub.IsSystemUser.Set(false)
	//spew.Dump(stub)

	return stub
}

func GetOSUserDetails() *user.User {
	currentUser, err := user.Current()
	if err != nil {
		logHandler.ErrorLogger.Fatalln(err.Error())
	}
	if application.IsRunningOnWindows() {
		//Example UID on windows "S-1-5-21-3849575818-2607088806-4266749144"
		//Tokenize the UID
		tkID := strings.Split(currentUser.Uid, "-")
		//Concatenate the token in 0, 1, 2, 3
		currentUser.Uid = fmt.Sprintf("%s%s%s%s", tkID[0], tkID[1], tkID[2], tkID[3])
	}
	//currentUser.Uid = fmt.Sprintf("%s_%s", currentUser.Uid, Application.HostName())
	currentUser.Uid = strings.Replace(currentUser.Uid, "-", "", -1)
	return currentUser
}
