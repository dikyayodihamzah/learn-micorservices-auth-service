package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/learn-micorservices/auth-service/config"
	"gitlab.com/learn-micorservices/auth-service/controller"
	"gitlab.com/learn-micorservices/auth-service/helper"
	"gitlab.com/learn-micorservices/auth-service/repository"
	"gitlab.com/learn-micorservices/auth-service/service"
)

func controllers() {
	time.Local = time.UTC

	serverConfig := config.NewServerConfig()
	db := config.NewDB
	validate := validator.New()

	authRepository := repository.NewAuthRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	authService := service.NewAuthService(authRepository, roleRepository, validate)
	authController := controller.NewAuthController(authService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	authController.NewAuthRouter(app)

	err := app.Listen(serverConfig.URI)
	log.Println(err)
}

func main() {
	time.Local = time.UTC
	topics := []string{helper.KafkaTopicProfile, helper.KafkaTopicAuth, helper.KafkaTopicRole}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumer := helper.NewKafkaConsumer()

	log.Printf("Created Consumer %v\n", consumer)

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		log.Fatal("Error on subscribe topics", err.Error())
	}

	go controllers()

	db := config.NewDB
	authRepository := repository.NewAuthRepository(db)
	kafkaUserConsumerService := service.NewKafkaUserConsumerService(authRepository)

	roleRepository := repository.NewRoleRepository(db)
	kafkaRoleConsumerService := service.NewKafkaRoleConsumerService(roleRepository)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("Message on %s:\n", e.TopicPartition)
				if e.Headers != nil {
					log.Printf("Headers: %v\n", e.Headers)
				}

				method := fmt.Sprintf("%v", e.Headers)

				switch method {

				// USER
				case `[method="POST.USER"]`:
					err := kafkaUserConsumerService.Insert(e.Value)
					helper.FatalIfError(err)
				case `[method="PUT.USER"]`:
					err := kafkaUserConsumerService.Update(e.Value)
					helper.FatalIfError(err)
				case `[method="DELETE.USER"]`:
					err := kafkaUserConsumerService.Delete(e.Value)
					helper.FatalIfError(err)

				// ROLE
				case `[method="POST.ROLE"]`:
					err := kafkaRoleConsumerService.Insert(e.Value)
					helper.FatalIfError(err)
				case `[method="PUT.ROLE"]`:
					err := kafkaRoleConsumerService.Update(e.Value)
					helper.FatalIfError(err)
				case `[method="DELETE.ROLE"]`:
					err := kafkaRoleConsumerService.Delete(e.Value)
					helper.FatalIfError(err)
				}

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				log.Printf("Ignored %v\n", e)
			}
		}
	}

	consumer.Close()
	log.Println("Kafka consumer closed")
}
