package spider

import (
	"sync"

	"go.uber.org/zap"

	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
)

type Transport interface {
	Yield(ty *scraping.TaskYield) error
	Finish(tf *scraping.TaskFinish) error
}

type TaskReporter struct {
	task  *scraping.Task
	tr    Transport
	group *sync.WaitGroup
	log   *zap.SugaredLogger
}

func NewTaskReporter(task *scraping.Task, tr Transport) *TaskReporter {
	log := logging.DefaultLogger().With("task_id", task.Id)

	return &TaskReporter{
		task:  task,
		tr:    tr,
		group: &sync.WaitGroup{},
		log:   log,
	}
}

// don't call after finish
func (r *TaskReporter) Report(job *scraping.Job, anime *data.Anime) {
	r.group.Add(1)

	go func(r *TaskReporter) {
		defer r.group.Done()

		msg := &scraping.TaskYield{
			TaskId: r.task.Id,
			JobId:  job.Id,
			Anime:  anime,
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
