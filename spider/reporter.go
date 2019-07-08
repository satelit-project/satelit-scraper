package spider

import (
	"sync"

	"github.com/sirupsen/logrus"

	"satelit-project/satelit-scraper/proto/scraper"
)

type Transport interface {
	Yield(ty *scraper.TaskYield) error
	Finish(tf *scraper.TaskFinish) error
}

type TaskReporter struct {
	task  *scraper.Task
	tr    Transport
	group *sync.WaitGroup
	log   *logrus.Entry
}

func NewTaskReporter(task *scraper.Task, tr Transport) *TaskReporter {
	log := logrus.WithField("task_id", task.Id)

	return &TaskReporter{
		task:  task,
		tr:    tr,
		group: &sync.WaitGroup{},
		log:   log,
	}
}

// don't call after finish
func (r *TaskReporter) Report(anime *scraper.Anime, scheduleID int32) {
	r.group.Add(1)

	go func(r *TaskReporter) {
		defer r.group.Done()

		msg := &scraper.TaskYield{
			TaskId:     r.task.Id,
			ScheduleId: scheduleID,
			Anime:      anime,
		}

		if err := r.tr.Yield(msg); err != nil {
			r.log.Errorf("failed to yield task: %v", err)
		}
	}(r)
}

func (r *TaskReporter) Finish() {
	r.group.Wait()

	msg := &scraper.TaskFinish{
		TaskId: r.task.Id,
	}

	if err := r.tr.Finish(msg); err != nil {
		r.log.Errorf("failed to finalize task: %v", err)
	}
}
