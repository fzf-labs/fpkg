package asynq

import (
	"context"
	"errors"
	"time"

	"github.com/fzf-labs/fpkg/transport"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

var (
	_ transport.Server = (*Server)(nil)
)

type ServerOption func(o *Server)

func WithConcurrency(concurrency int) ServerOption {
	return func(s *Server) {
		s.asynqConfig.Concurrency = concurrency
	}
}

func WithQueues(queues map[string]int) ServerOption {
	return func(s *Server) {
		s.asynqConfig.Queues = queues
	}
}

func WithRetryDelayFunc(fn asynq.RetryDelayFunc) ServerOption {
	return func(s *Server) {
		s.asynqConfig.RetryDelayFunc = fn
	}
}

func WithStrictPriority(val bool) ServerOption {
	return func(s *Server) {
		s.asynqConfig.StrictPriority = val
	}
}

func WithErrorHandler(fn asynq.ErrorHandler) ServerOption {
	return func(s *Server) {
		s.asynqConfig.ErrorHandler = fn
	}
}

func WithHealthCheckFunc(fn func(error)) ServerOption {
	return func(s *Server) {
		s.asynqConfig.HealthCheckFunc = fn
	}
}

func WithHealthCheckInterval(tm time.Duration) ServerOption {
	return func(s *Server) {
		s.asynqConfig.HealthCheckInterval = tm
	}
}

func WithDelayedTaskCheckInterval(tm time.Duration) ServerOption {
	return func(s *Server) {
		s.asynqConfig.DelayedTaskCheckInterval = tm
	}
}

func WithGroupGracePeriod(tm time.Duration) ServerOption {
	return func(s *Server) {
		s.asynqConfig.GroupGracePeriod = tm
	}
}

func WithGroupMaxDelay(tm time.Duration) ServerOption {
	return func(s *Server) {
		s.asynqConfig.GroupMaxDelay = tm
	}
}

func WithGroupMaxSize(sz int) ServerOption {
	return func(s *Server) {
		s.asynqConfig.GroupMaxSize = sz
	}
}

func WithMiddleware(m ...asynq.MiddlewareFunc) ServerOption {
	return func(o *Server) {
		o.mux.Use(m...)
	}
}

func WithLocation(loc *time.Location) ServerOption {
	return func(s *Server) {
		s.schedulerOpts.Location = loc
	}
}

type Server struct {
	started bool

	baseCtx context.Context
	err     error

	asynqServer    *asynq.Server
	asynqClient    *asynq.Client
	asynqScheduler *asynq.Scheduler

	mux           *asynq.ServeMux
	asynqConfig   asynq.Config
	redisOpt      *asynq.RedisClientOpt
	schedulerOpts *asynq.SchedulerOpts
}

func NewServer(redisOpt *asynq.RedisClientOpt, opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx:  context.Background(),
		started:  false,
		redisOpt: redisOpt,
		asynqConfig: asynq.Config{
			Concurrency: 10,
			Logger:      NewLogger(),
		},
		schedulerOpts: &asynq.SchedulerOpts{},
		mux:           asynq.NewServeMux(),
	}

	srv.init(opts...)

	return srv
}

func (s *Server) Name() string {
	return "asynq"
}

func (s *Server) HandleFunc(pattern string, handler func(context.Context, *asynq.Task) error) error {
	if s.started {
		slog.Error("[asynq] HandleFunc [%s] failed", pattern)
		return errors.New("cannot handle func, server already started")
	}
	s.mux.HandleFunc(pattern, handler)
	return nil
}

func (s *Server) Handle(pattern string, handler asynq.Handler) error {
	if s.started {
		slog.Error("[asynq] Handle [%s] failed", pattern)
		return errors.New("cannot handle, server already started")
	}
	s.mux.Handle(pattern, handler)
	return nil
}

// NewTask enqueue a new task
func (s *Server) NewTask(typeName string, payload []byte, opts ...asynq.Option) error {
	if s.asynqClient == nil {
		if err := s.createAsynqClient(); err != nil {
			return err
		}
	}

	task := asynq.NewTask(typeName, payload)
	info, err := s.asynqClient.Enqueue(task, opts...)
	if err != nil {
		slog.Error("[asynq] [%s] Enqueue failed: %s", typeName, err.Error())
		return err
	}
	slog.Debug("[asynq] enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return nil
}

// NewPeriodicTask enqueue a new crontab task
func (s *Server) NewPeriodicTask(cronSpec, typeName string, payload []byte, opts ...asynq.Option) error {
	if s.asynqScheduler == nil {
		if err := s.createAsynqScheduler(); err != nil {
			return err
		}
		if err := s.runAsynqScheduler(); err != nil {
			return err
		}
	}

	task := asynq.NewTask(typeName, payload)

	entryID, err := s.asynqScheduler.Register(cronSpec, task, opts...)
	if err != nil {
		slog.Error("[asynq] [%s] Enqueue failed: %s", typeName, err.Error())
		return err
	}
	slog.Debug("[asynq] registered an entry: id=%q", entryID)

	return nil
}

func (s *Server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	if s.started {
		return nil
	}

	if err := s.runAsynqServer(); err != nil {
		return err
	}

	slog.Info("[asynq] server listening on: %s", s.redisOpt.Addr)

	s.baseCtx = ctx
	s.started = true

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	slog.Info("[asynq] server stopping")
	s.started = false

	if s.asynqClient != nil {
		_ = s.asynqClient.Close()
		s.asynqClient = nil
	}

	if s.asynqServer != nil {
		s.asynqServer.Stop()
		s.asynqServer = nil
	}

	if s.asynqScheduler != nil {
		s.asynqScheduler.Shutdown()
		s.asynqScheduler = nil
	}

	return nil
}

func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
	_ = s.createAsynqServer()
}

func (s *Server) createAsynqServer() error {
	if s.asynqServer != nil {
		slog.Error("[asynq] asynq server already created")
		return errors.New("asynq server already created")
	}

	s.asynqServer = asynq.NewServer(s.redisOpt, s.asynqConfig)
	if s.asynqServer == nil {
		slog.Error("[asynq] create asynq server failed")
		return errors.New("create asynq server failed")
	}
	return nil
}

func (s *Server) runAsynqServer() error {
	if s.asynqServer == nil {
		slog.Error("[asynq] asynq server is nil")
		return errors.New("asynq server is nil")
	}

	if err := s.asynqServer.Run(s.mux); err != nil {
		slog.Error("[asynq] asynq server run failed: %s", err.Error())
		return err
	}
	return nil
}

func (s *Server) createAsynqClient() error {
	if s.asynqClient != nil {
		slog.Error("[asynq] asynq client already created")
		return errors.New("asynq client already created")
	}

	s.asynqClient = asynq.NewClient(s.redisOpt)
	if s.asynqClient == nil {
		slog.Error("[asynq] create asynq client failed")
		return errors.New("create asynq client failed")
	}

	return nil
}

func (s *Server) createAsynqScheduler() error {
	if s.asynqScheduler != nil {
		slog.Error("[asynq] asynq scheduler already created")
		return errors.New("asynq scheduler already created")
	}

	s.asynqScheduler = asynq.NewScheduler(s.redisOpt, s.schedulerOpts)
	if s.asynqScheduler == nil {
		slog.Error("[asynq] create asynq scheduler failed")
		return errors.New("create asynq scheduler failed")
	}

	return nil
}

func (s *Server) runAsynqScheduler() error {
	if s.asynqScheduler == nil {
		slog.Error("[asynq] asynq scheduler is nil")
		return errors.New("asynq scheduler is nil")
	}

	if err := s.asynqScheduler.Start(); err != nil {
		slog.Error("[asynq] asynq scheduler start failed: %s", err.Error())
		return err
	}
	return nil
}
