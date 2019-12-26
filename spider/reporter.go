package spider

import (
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
)

// Represents an object that can send jobs reports to external service.
type Transport interface {
	// Signals that task has made progress.
	Yield(ty *scraping.TaskYield) error

	// Signals that task has been finished.
	Finish(tf *scraping.TaskFinish) error
}

// Reports scraping progress for a given task.
type TaskReporter struct {
	// Task to report progress for.
	Task      *scraping.Task

	// An object which can communicate with remote service.
	Transport Transport
}

// Reports that there's new scraped anime entity for the task.
func (r *TaskReporter) Report(job *scraping.Job, anime *data.Anime) error {
	msg := &scraping.TaskYield{
		TaskId: r.Task.Id,
		JobId:  job.Id,
		Anime:  anime,
	}

	return r.Transport.Yield(msg)
}

// Reports that scraping has been finished for the task.
func (r *TaskReporter) Finish() error {
	msg := &scraping.TaskFinish{
		TaskId: r.Task.Id,
	}

	return r.Transport.Finish(msg)
}
