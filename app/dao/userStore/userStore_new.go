package userStore

import (
	"context"
	"fmt"

	"github.com/mt1976/frantic-agentis/app/business/translation"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, userName, uid, realName, email, gid, host string, isActive, isSystemUser bool) (User_Store, error) {
	return new(ctx, userName, uid, realName, email, gid, host, isActive, isSystemUser, true)
}
func new(ctx context.Context, userName, uid, realName, email, gid, host string, isActive, isSystemUser, canLogin bool) (User_Store, error) {

	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.InfoLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", userName))

	// Create a new struct
	record := User_Store{}

	record.RealName = realName
	record.UID = uid
	record.UserName = userName
	record.Email = email
	record.GID = gid
	record.IsActive.Set(isActive)
	record.Host = host
	record.IsSystemUser.Set(isSystemUser)
	record.CanLogin.Set(canLogin)

	record.Key = idHelpers.Encode(record.buildUserCode())
	record.Raw = record.buildUserCode()
	record.UserCode = record.buildUserCode()

	record.Notes = "These are some notes for " + realName + " " + record.UserCode

	// Check for duplicates
	xUser, err := record.isDuplicateOf(record.UserCode)
	if err == commonErrors.ErrorDuplicate {
		// This is OK, do nothing as this is a duplicate record
		// we ignore duplicates.
		logHandler.WarningLogger.Printf(translation.Get("[%v] DUPLICATE %v already in use '%v'"), appName, domain, realName)
		clock.Stop(1)
		return xUser, nil
	}

	// Record the create action in the audit data
	auditErr := record.Audit.Action(context.TODO(), audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", domain, userName)))
	if auditErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOUpdateAuditError(domain, record.ID, auditErr))
	}

	// Log the dest instance before the creation
	xUser, err = record.prepare()
	if err == commonErrors.ErrorDuplicate {
		// This is OK, do nothing as this is a duplicate record
		// we ignore duplicate destinations.
		logHandler.WarningLogger.Printf(translation.Get("[%v] DUPLICATE %v already in use '%v'"), appName, domain, realName)
		return xUser, nil
	}

	// Save the status instance to the database
	writeErr := activeDB.Create(&record)
	if writeErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOCreateError(domain, record.ID, writeErr))
		//	panic(writeErr)
	}

	//logHandler.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, strings.ToUpper(domain), record.ID, fmt.Sprintf("New %v: %v", domain, field1))
	clock.Stop(1)
	return record, nil
}
