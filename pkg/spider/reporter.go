package spider

import (
	"satelit-project/satelit-scraper/pkg/proto/scraping"
	"sync"

	"satelit-project/satelit-scraper/pkg/proto/data"

	"github.com/sirupsen/logrus"
)

type Transport interface {
	Yield(ty *scraping.TaskYield) error
	Finish(tf *scraping.TaskFinish) error
}

type TaskReporter struct {
	task  *scraping.Task
	tr    Transport
	group *sync.WaitGroup
	log   *logrus.Entry
}

func NewTaskReporter(task *scraping.Task, tr Transport) *TaskReporter {
	log := logrus.WithField("task_id", task.Id)

	return &TaskReporter{
		task:  task,
		tr:    tr,
		group: &sync.WaitGroup{},
		log:   log,
	}
}

// don't call after finish
func (r *TaskReporter) Report(anime *data.Anime, scheduleID int32) {
	r.group.Add(1)

	go func(r *TaskReporter) {
		defer r.group.Done()

		msg := &scraping.TaskYield{
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

	msg := &scraping.TaskFinish{
		TaskId: r.task.Id,
	}

	if err := r.tr.Finish(msg); err != nil {
		r.log.Errorf("failed to finalize task: %v", err)
	}
}
