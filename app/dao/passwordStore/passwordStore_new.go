package passwordStore

import (
	"context"
	"fmt"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, userKey, password, source string) (Password_Store, error) {

	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.SecurityLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", userKey))

	//uid := strconv.Itoa(userID)

	// Create a new struct
	record := Password_Store{}
	record.Key = userKey
	record.Raw = userKey
	record.UserKey = userKey
	record.Source = source // Source of the password
	if source == "" {
		record.Source = cfg.GetApplication_Name()
	}
	record.Expired.Set(false)

	pwd, err := encode(password)
	if err != nil {
		logHandler.ErrorLogger.Panicf("[%v] %v", appName, commonErrors.WrapDAOCreateError(domain, record.ID, err))
	}
	record.Password = pwd
	// Record the create action in the audit data
	auditErr := record.Audit.Action(ctx, audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", domain, userKey)))
	if auditErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panicf("[%v] %v", appName, commonErrors.WrapDAOCreateError(domain, record.ID, auditErr))
	}

	// Save the status instance to the database
	writeErr := activeDB.Create(&record)
	if writeErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panicf("[%v] %v", appName, commonErrors.WrapDAOCreateError(domain, record.ID, writeErr))
	}

	clock.Stop(1)
	return record, nil
}
