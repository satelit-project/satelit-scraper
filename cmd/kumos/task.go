package main

import (
	"context"
	"fmt"
	"net"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"

	empty "github.com/golang/protobuf/ptypes/empty"
	uuid "shitty.moe/satelit-project/satelit-scraper/proto/common"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

type TaskService struct {
	jobs  []int32
	inner *grpc.Server
	port  int
	log   *logging.Logger
}

func NewTaskService(port int, log *logging.Logger) *TaskService {
	return &TaskService{
		inner: grpc.NewServer(),
		port:  port,
		log:   log,
	}
}

func (s *TaskService) SetJobs(j []int32) {
	s.jobs = j
}

func (s *TaskService) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	scraping.RegisterScraperTasksServiceServer(s.inner, s)
	if err = s.inner.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) CreateTask(ctx context.Context, in *scraping.TaskCreate) (*scraping.Task, error) {
	var jobs []*scraping.Job
	for _, j := range s.jobs {
		job := &scraping.Job{
			Id:      &uuid.UUID{},
			AnimeId: j,
		}
		jobs = append(jobs, job)
	}

	task := &scraping.Task{
		Id:     &uuid.UUID{},
		Source: data.Source_ANIDB,
		Jobs:   jobs,
	}

	s.log.Infof("created new task")
	return task, nil
}

func (s *TaskService) YieldResult(ctx context.Context, in *scraping.TaskYield) (*empty.Empty, error) {
	m := jsonpb.Marshaler{Indent: "  "}
	js, err := m.MarshalToString(in)
	if err != nil {
		s.log.Errorf("failed to marchal json: %v", err)
		return nil, err
	}

	s.log.Infof("received task yield:\n%s", js)
	return &empty.Empty{}, nil
}

func (s *TaskService) CompleteTask(ctx context.Context, in *scraping.TaskFinish) (*empty.Empty, error) {
	s.log.Infof("task finished")
	return &empty.Empty{}, nil
}
