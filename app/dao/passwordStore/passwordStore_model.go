package passwordStore

// Data Access Object Password
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"github.com/mt1976/frantic-core/dao"
	audit "github.com/mt1976/frantic-core/dao/audit"
)

// Password_Store represents a Password_Store entity.
type Password_Store struct {
	ID       int           `csv:"-" storm:"id,increment=100000"` // primary key with auto increment
	Key      string        `csv:"-" storm:"unique"`              // key
	Raw      string        `csv:"-" storm:"unique"`              // raw ID before encoding
	UserKey  string        `csv:"-" storm:"index"`               // user key
	Password string        `csv:"-"`
	Expired  dao.StormBool `csv:"-" storm:"index"` // is expired
	Source   string        `csv:"-"`               // Source Applicaiton of the user (if applicable, for future use)
	Audit    audit.Audit   `csv:"-"`               // audit data
}

// Define the field set as names
var (
	FIELD_ID       = "ID"
	FIELD_Key      = "Key"
	FIELD_Raw      = "Raw"
	FIELD_UserKey  = "UserKey"
	FIELD_Password = "Password"
	FIELD_Expired  = "Expired"
	FIELD_Source   = "Source" // Effectively the app that set the password
	FIELD_Audit    = "Audit"
)

var domain = "Password"
