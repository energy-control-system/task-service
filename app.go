package main

import (
	"context"
	"fmt"
	"io/fs"
	"task-service/api"
	"task-service/config"
	dbtask "task-service/database/task"
	"task-service/service/task"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const (
	serviceName = "task-service"
	dbTimeout   = 15 * time.Second
)

type App struct {
	/* main */
	mainCtx  context.Context
	log      golog.Logger
	settings config.Settings

	/* http */
	server goserver.Server

	/* db */
	postgres           *sqlx.DB
	kafka              gokafka.Kafka
	taskProducer       gokafka.Producer
	inspectionConsumer gokafka.Consumer

	/* services */
	taskService *task.Service
}

func NewApp(mainCtx context.Context, log golog.Logger, settings config.Settings) *App {
	return &App{
		mainCtx:  mainCtx,
		log:      log,
		settings: settings,
	}
}

func (a *App) InitDatabases(fs fs.FS, path string) (err error) {
	postgresCtx, cancelPostgresCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelPostgresCtx()

	a.postgres, err = db.NewPgx(postgresCtx, a.settings.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}

	err = db.Migrate(fs, a.log, a.postgres, path)
	if err != nil {
		return fmt.Errorf("migrate postgres: %w", err)
	}

	a.kafka = gokafka.NewKafka(a.settings.Databases.Kafka.Brokers)

	a.taskProducer = a.kafka.Producer(a.settings.Databases.Kafka.Topics.Tasks)

	a.inspectionConsumer, err = a.kafka.Consumer(a.log.WithTags("inspectionConsumer"), func() (context.Context, context.CancelFunc) {
		return context.WithCancel(a.mainCtx)
	}, gokafka.WithTopic(a.settings.Databases.Kafka.Topics.Inspections), gokafka.WithConsumerGroup(serviceName))
	if err != nil {
		return fmt.Errorf("init inspection consumer: %w", err)
	}

	return nil
}

func (a *App) InitServices() error {
	taskRepository := dbtask.NewRepository(a.postgres)

	taskPublisher := task.NewPublisher(a.mainCtx, a.taskProducer)

	a.taskService = task.NewService(taskRepository, taskPublisher)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.mainCtx, a.log, a.settings)
	sb.AddDebug()
	sb.AddTasks(a.taskService)

	a.server = sb.Build()
}

func (a *App) Start() {
	a.server.Start()
	a.inspectionConsumer.Subscribe(a.taskService.SubscriberOnInspectionEvent(a.mainCtx, a.log.WithTags("inspectionSubscriber")))
}

func (a *App) Stop(ctx context.Context) {
	consumerCtx, cancelConsumerCtx := context.WithTimeout(ctx, dbTimeout)
	defer cancelConsumerCtx()

	err := a.inspectionConsumer.Close(consumerCtx)
	if err != nil {
		a.log.Errorf("failed to close inspection consumer: %v", err)
	}

	a.server.Stop()

	producerCtx, cancelProducerCtx := context.WithTimeout(ctx, dbTimeout)
	defer cancelProducerCtx()

	err = a.taskProducer.Close(producerCtx)
	if err != nil {
		a.log.Errorf("failed to close task producer: %v", err)
	}

	err = a.postgres.Close()
	if err != nil {
		a.log.Errorf("failed to close postgres connection: %v", err)
	}
}
