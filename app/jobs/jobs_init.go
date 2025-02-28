package jobs

import (
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/robfig/cron/v3"
)

var scheduledTasks *cron.Cron
var cfg *commonConfig.Settings
var domain = "jobs"
var appName string

func init() {

	cfg = commonConfig.Get()
	appName = "frantic-agentis"
	err := jobs.Initialise(cfg)
	if err != nil {
		logHandler.ServiceLogger.Fatalf("[%v] Error: %v", domain, err)
		panic(err)
	}
}
