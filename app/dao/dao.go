package dao

import (
	"context"

	"github.com/asdine/storm/v3"

	"github.com/mt1976/frantic-agentis/app/dao/passwordStore"
	"github.com/mt1976/frantic-agentis/app/dao/userStore"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var name = "DAO"
var Version = 1
var DB *storm.DB
var tableName = "database"
var appName = "frantic-agentis"

func Initialise(cfg *commonConfig.Settings) error {
	clock := timing.Start(name, actions.INITIALISE.GetCode(), "")
	logHandler.EventLogger.Printf("[%v] Initialising %v - Started", appName, name)

	userStore.Initialise(context.TODO(), false)

	passwordStore.Initialise(context.TODO(), false)

	logHandler.EventLogger.Printf("[%v] Initialising %v - Complete", appName, name)
	clock.Stop(1)
	return nil
}

func ExportAllToCSV() {
	userStore.ExportCSV()
	passwordStore.ExportCSV()
}
