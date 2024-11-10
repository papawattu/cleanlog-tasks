package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	"github.com/papawattu/cleanlog-tasks/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventInterceptor struct {
	ctx   context.Context
	Next  TaskRepository
	conn  *amqp.Connection
	queue amqp.Queue
}

func (ei *EventInterceptor) GetTask(id int) (*Task, error) {
	log.Printf("Getting task: %v", id)
	return ei.Next.GetTask(id)
}

func (ei *EventInterceptor) SaveTask(t *Task) error {
	log.Printf("Saving task: %v", t)

	var data bytes.Buffer

	ev := types.TaskEvent{
		EntityType:    "NewTask",
		EntityVersion: 1,
		TaskId:        *t.TaskID,
		Description:   t.TaskDescription,
	}
	enc := gob.NewEncoder(&data)
	err := enc.Encode(ev)
	if err != nil {
		log.Fatalf("Error encoding task: %v", err)
	}
	ch, err := ei.conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}
	defer ch.Close()

	err = ch.PublishWithContext(
		ei.ctx,
		"",
		ei.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        data.Bytes(),
		})

	if err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	return ei.Next.SaveTask(t)
}

func NewEventInterceptor(ctx context.Context, queueName string, amqUri string, next TaskRepository) TaskRepository {

	conn, err := amqp.Dial(amqUri)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("Error declaring queue: %v", err)
	}

	return &EventInterceptor{
		Next:  next,
		conn:  conn,
		queue: q,
	}
}
