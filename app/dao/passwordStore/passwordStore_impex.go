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

// importProcessor is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func importProcessor(inOriginal **Password_Store) (string, error) {

	importedData := **inOriginal

	pw, err := load(context.TODO(), importedData.UserKey, importedData.Password)
	if err != nil {
		logHandler.ImportLogger.Panicf("[%v] Error importing %v: %v [%v]", appName, domain, importedData.Key, err.Error())
		return pw.Key, err
	}

	return pw.Key, nil
}

func load(ctx context.Context, userKey, password string) (Password_Store, error) {

	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.SecurityLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", userKey))

	//uid := strconv.Itoa(userID)

	// Create a new struct
	record := Password_Store{}
	record.Key = userKey
	record.Raw = userKey
	record.UserKey = userKey
	record.Password = password
	// Record the create action in the audit data
	auditErr := record.Audit.Action(ctx, audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", domain, userKey)))
	if auditErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panicf("[%v] %v", appName, commonErrors.WrapDAOUpdateAuditError(domain, record.ID, auditErr))
	}
	// Check if the password already exists
	checkRecord, err := GetByKey(userKey)
	checkRecord.Password = password
	var writeErr error
	if err == nil {
		logHandler.WarningLogger.Printf("[%v] Updating %v: %v", appName, domain, userKey)
		writeErr = activeDB.Create(&checkRecord)
	} else {
		// Save the status instance to the database
		logHandler.WarningLogger.Printf("[%v] Creating %v: %v", appName, domain, userKey)
		writeErr = activeDB.Create(&record)
	}
	if writeErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panicf("[%v] %v", appName, commonErrors.WrapDAOInitialisationError(domain, writeErr))
	}

	clock.Stop(1)
	return record, nil
}
