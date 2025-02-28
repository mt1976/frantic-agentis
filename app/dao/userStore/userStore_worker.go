package userStore

import (
	"context"

	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func jobProcessor(job jobs.Job) {
	// Do something every day at midnight
	name := jobs.CodedName(job)

	clock := timing.Start(name, job.Name(), "Scheduled Service")
	// Get all the users
	userList, err := GetAll()
	if err != nil {
		logHandler.ServiceLogger.Printf("[%v] [%v] - Error=[%e]", appName, job.Name(), err)
		return
	}
	// Loop through the users
	for _, user := range userList {
		// Get the status of the user
		if user.IsActive.Bool() {
			go userWorker(&user)
		} else {
			logHandler.ServiceLogger.Printf("[%v] [%v] - user=[%v] host=[%v] - Skipping", appName, job.Name(), user.UserCode, user.Host)
		}
	}

	clock.Stop(0)
}

func userWorker(user *User_Store) error {
	Initialise(context.TODO())

	//logHandler.InfologHandler.Printf("user Worker - [%v] - Running", h.Name, cfg.GetApplicationName())

	job := timing.Start(domain, actions.RUN.GetCode(), user.UserCode)

	logHandler.ServiceLogger.Printf("[%v] User=[%v] Host=[%v] - Running", appName, user.UserCode, user.Host)

	job.Stop(1)

	return nil
}
