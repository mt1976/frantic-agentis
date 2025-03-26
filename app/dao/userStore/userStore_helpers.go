package userStore

// Data Access Object User - User_Store
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func (record *User_Store) prepare() (User_Store, error) {
	//os.Exit(0)
	//logger.ErrorLogger.Printf("ACT: VAL Validate")
	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Insert any defaults/validation below  here

	user, err := record.dup()
	if err == commonErrors.ErrorDuplicate {
		return *record, nil
	}
	if err != nil {
		return user, err
	}

	// Insert any defaults/validation above here

	return *record, nil
}

func (record *User_Store) calculate() error {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Calculate the duration in days between the start and end dates
	return nil
}

func (record *User_Store) isDuplicateOf(id string) (User_Store, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.
	if id == "" {
	}
	return record.dup()
	// // Get all status
	// recordList, err := GetAll()
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Getting all %v failed %v", domain, err.Error())
	// 	return User_Store{}, err
	// }

	// // range through status list, if status code is found and deletedby is empty then return error
	// for _, checkRecord := range recordList {
	// 	//s.Dump(!,strings.ToUpper(code) + "-uchk-" + s.Code)
	// 	testValue := strings.ToUpper(id)
	// 	checkValue := strings.ToUpper(checkRecord.Key)
	// 	//logger.InfoLogger.Printf("CHK: TestValue:[%v] CheckValue:[%v]", testValue, checkValue)
	// 	//logger.InfoLogger.Printf("CHK: Code:[%v] s.Code:[%v] s.Audit.DeletedBy:[%v]", testCode, s.Code, s.Audit.DeletedBy)
	// 	if checkValue == testValue && checkRecord.Audit.DeletedBy == "" {
	// 		logHandler.WarningLogger.Printf("Duplicate %v, %v already in use", strings.ToUpper(domain), record.ID)
	// 		return checkRecord, commonErrors.ErrorDuplicate
	// 	}
	// }

	// return User_Store{}, nil
}

func (record *User_Store) IsDuplicateOf() (User_Store, error) {
	return record.dup()
}

func (u *User_Store) dup() (User_Store, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Get all status
	foundUser, err := GetBy(FIELD_UserCode, u.UserCode)
	if err == nil {
		//	logHandler.WarningLogger.Printf("[%v] DUPLICATE %v already exist '%v' (%v)", appName, domain, u.UserCode, foundUser.UserName)
		return foundUser, commonErrors.ErrorDuplicate
	}

	//logger.InfoLogger.Printf("CHK: %v is unique", strings.ToUpper(name))

	// Return nil if the code is unique

	return *u, nil
}

func ClearDown(ctx context.Context) error {
	logHandler.EventLogger.Printf("[%v] Clearing %v", appName, domain)

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.CLEAR.GetCode(), "Cleardown")

	// Delete all active session recordList
	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, err).Error())
		clock.Stop(0)
		return commonErrors.WrapDAOInitialisationError(domain, err)
	}

	noRecords := len(recordList)
	count := 0

	for thisRecord, record := range recordList {
		logHandler.SecurityLogger.Printf("[%v] Deleting %v (%v/%v) %v", appName, domain, thisRecord, noRecords, record.Key)
		delErr := Delete(ctx, record.ID, fmt.Sprintf("Clearing %v %v @ initialisation ", domain, record.ID))
		if delErr != nil {
			logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, delErr).Error())
			continue
		}
		count++
	}

	clock.Stop(count)

	return nil
}

func GetByUserName(userName string) (User_Store, error) {
	xx, err := GetBy(FIELD_UserName, userName)
	if err != nil {
		logHandler.SecurityLogger.Printf("[%v] User Not Found!: UserName=[%v]", appName, userName)
		return User_Store{}, err
	}
	logHandler.SecurityLogger.Printf("[%v] User Found!: UserName=[%v] UserCode=[%v]", appName, userName, xx.UserCode)
	return xx, nil
}

func (u *User_Store) SetName(name string) error {
	if name == "" {
		return commonErrors.ErrorEmptyName
	}
	if len(name) > 50 {
		return commonErrors.ErrorNameTooLong
	}
	u.RealName = name
	return nil
}

func (u *User_Store) buildUserCode() string {
	return fmt.Sprintf("%v_%v", u.UID, u.UserName)
}

func GetByUserCode(code string) (User_Store, error) {
	return GetBy(FIELD_UserCode, code)
}

func GetByUID(code int) (User_Store, error) {
	if code == 0 {
		return User_Store{}, fmt.Errorf("[%v] Reading UID=[%v] %v ", strings.ToUpper(domain), code, "UID is blank/zero")
	}
	return GetBy(FIELD_UID, code)
}
