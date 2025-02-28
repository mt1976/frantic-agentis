package userStore

// Data Access Object User
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"

	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var activeDB *database.DB
var initialised bool = false // default to false
var cfg *commonConfig.Settings
var appName string
var useIsolatedDB bool = true

func Initialise(ctx context.Context, isolateDB bool) {
	//logHandler.EventLogger.Printf("Initialising %v", domain)
	timing := timing.Start(domain, actions.INITIALISE.GetCode(), "Initialise")
	cfg = commonConfig.Get()
	appName = "frantic-agentis"
	useIsolatedDB = isolateDB
	// For a specific database connection, use NamedConnect, otherwise use Connect
	if useIsolatedDB {
		activeDB = database.ConnectToNamedDB("agentis")
	} else {
		activeDB = database.Connect()
	}
	initialised = true
	timing.Stop(1)
	logHandler.EventLogger.Printf("[%v] Initialised %v", appName, domain)
}
