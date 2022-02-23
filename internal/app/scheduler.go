package app

import (
	"context"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/dentych/taskeroo/internal/util"
	"log"
	"time"
)

type Scheduler struct {
	notificationLogic *NotificationLogic
	taskLogic         *TaskLogic
	groupRepo         *database.GroupRepo

	context context.Context
	cancel  context.CancelFunc
}

func NewScheduler(notificationLogic *NotificationLogic, taskLogic *TaskLogic, groupRepo *database.GroupRepo) *Scheduler {
	return &Scheduler{notificationLogic: notificationLogic, taskLogic: taskLogic, groupRepo: groupRepo}
}

func (s *Scheduler) Start() {
	s.context, s.cancel = context.WithCancel(context.Background())
	go s.noonTask()
}

func (s *Scheduler) noonTask() {
	for {
		durationUntilTomorrowNoon := util.DurationToNextNoon(time.Now())
		log.Printf("SCHEDULER: Waiting till noon. Duration in hours: %f\n", durationUntilTomorrowNoon.Hours())
		time.Sleep(durationUntilTomorrowNoon)

		log.Printf("SCHEDULER: Running notify tasks due today")
		err := s.taskLogic.NotifyTasksDueToday(s.context)
		if err != nil {
			log.Printf("ERROR: NoonTask: Error during notification of tasks due today: %s", err)
		}
		log.Printf("SCHEDULER: Done running notify tasks due today")
	}
}
