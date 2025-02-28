package userStore

import (
	"context"
	"fmt"

	"github.com/mt1976/frantic-core/logHandler"
)

// UserImportData is a struct to hold the data from the CSV file
// it is used to import the data into the database
// The struct tags are used to map the fields to the CSV columns
// this struct should be customised to suit the specific requirements of the entryination table/DAO.

var COMMA = '|'

// load is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func importProcessor(inOriginal **User_Store) (string, error) {

	importedData := **inOriginal

	//	logHandler.ImportLogger.Printf("Importing %v [%v] [%v]", domain, original.Raw, original.Field1)
	tempUserCode := importedData.buildUserCode()

	// Check if the user already exists
	record, err := GetByUserCode(tempUserCode)
	if err == nil {
		// Update the existing record
		record.UID = importedData.UID
		record.GID = importedData.GID
		record.RealName = importedData.RealName
		record.UserName = importedData.UserName
		record.Email = importedData.Email
		record.IsActive = importedData.IsActive
		record.Host = importedData.Host

		err := record.Update(context.TODO(), "Imported/Refreshed")
		if err != nil {
			logHandler.ImportLogger.Panicf("[%v] Error updating %v: %v [%v]", appName, domain, tempUserCode, err.Error())
			return record.UserCode, err
		}

		return fmt.Sprintf("[%v] Updated '%v'", appName, record.UserCode), nil
	} else {

		// New(ctx context.Context, userName, uid, realName, email, gid, host string,active bool)
		newUser, err := New(context.TODO(), importedData.UserName, importedData.UID, importedData.RealName, importedData.Email, importedData.GID, importedData.Host, importedData.IsActive.Bool(), importedData.IsSystemUser.Bool())
		if err != nil {
			logHandler.ImportLogger.Panicf("[%v] Error importing %v: %v [%v]", appName, domain, tempUserCode, err.Error())
			return newUser.UserCode, err
		}
		return fmt.Sprintf("[%v] Created '%v'", appName, newUser.UserCode), nil
	}
}
