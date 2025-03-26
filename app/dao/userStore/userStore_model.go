package userStore

// Data Access Object User
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"time"

	"github.com/mt1976/frantic-core/dao"
	audit "github.com/mt1976/frantic-core/dao/audit"
)

// User_Store represents a User_Store entity.
type User_Store struct {
	ID           int           `csv:"-" storm:"id,increment=100000"` // primary key with auto increment
	Key          string        `csv:"-" storm:"unique"`              // key
	Raw          string        `csv:"-" storm:"unique"`              // raw ID before encoding
	Audit        audit.Audit   `csv:"-"`                             // audit data
	UID          string        `validate:"required"`
	GID          string        `storm:"index" validate:"required"`
	RealName     string        `validate:"required,min=5"` // this field will not be indexed
	UserName     string        `validate:"required,min=5"`
	UserCode     string        `csv:"-" storm:"index" validate:"required,min=5"`
	Email        string        `storm:"index"`
	Notes        string        `csv:"-" validate:"max=75"`
	IsActive     dao.StormBool `csv:"-" storm:"index"`
	CanLogin     dao.StormBool `csv:"-" storm:"index"`
	LastLogin    time.Time     `csv:"-"`
	LastHost     string        `csv:"-"`
	Host         string        `storm:"index"`
	IsSystemUser dao.StormBool `csv:"-" storm:"index"` // is a system user
	Locale       string        ``
	Display      string        `csv:"-"`       // Display Name
	Source       string        ``              // Source Applicaiton of the user (if applicable, for future use)
	Theme        string        ``              // Theme
	Timezone     string        ``              // Timezone
	Role         string        `storm:"index"` // Role
}

// Define the field set as names
var (
	FIELD_ID           = "ID"
	FIELD_Key          = "Key"
	FIELD_Raw          = "Raw"
	FIELD_UID          = "UID"
	FIELD_GID          = "GID"
	FIELD_RealName     = "RealName"
	FIELD_UserName     = "UserName"
	FIELD_UserCode     = "UserCode"
	FIELD_Email        = "Email"
	FIELD_Notes        = "Notes"
	FIELD_IsActive     = "IsActive"
	FIELD_CanLogin     = "CanLogin"
	FIELD_LastLogin    = "LastLogin"
	FIELD_LastHost     = "LastHost"
	FIELD_Host         = "Host"
	FIELD_IsSystemUser = "IsSystemUser"
	FIELD_Locale       = "Locale"
	FIELD_Display      = "Display"
	FIELD_Source       = "Source"
	FIELD_Audit        = "Audit"
	FIELD_Theme        = "Theme"
	FIELD_Timezone     = "Timezone"
)

var domain = "User"
