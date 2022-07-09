package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}

	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}

	defer channelRabbitMQ.Close()

	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1",
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(
		logger.New(),
		) // add simple logger


//Add router
app.Get("/send", func(c *fiber.Ctx) error {
	message :=  amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(c.Query("msg")),
}

if err := channelRabbitMQ.Publish(
	"", // exchange
	"QueueService1", // routing key
	false, // mandatory
	false, // immediate
	message,
); err != nil {
		return err
}
	return nil
})

// start fiber server
log.Fatal(app.Listen(":3000"))
}