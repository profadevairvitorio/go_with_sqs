package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	localStackEndpoint = "http://localhost:4566"
	queueName          = "my-local-queue"
)

func createSQSSession() *sqs.SQS {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:   aws.String(localStackEndpoint),
		Region:     aws.String("us-east-1"),
		DisableSSL: aws.Bool(true),
		MaxRetries: aws.Int(3),
	}))

	sqsSvc := sqs.New(sess)

	return sqsSvc
}

func createQueue(svc *sqs.SQS) (*string, error) {
	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Fila criada com sucesso: %s\n", *result.QueueUrl)
	return result.QueueUrl, nil
}

func sendMessage(svc *sqs.SQS, queueURL *string, messageBody string) error {
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    queueURL,
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		return err
	}

	log.Printf("Mensagem enviada com sucesso: %s\n", messageBody)
	return nil
}

func receiveMessage(svc *sqs.SQS, queueURL *string) error {
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(20),
	})
	if err != nil {
		return err
	}

	for _, message := range result.Messages {
		log.Printf("Mensagem recebida: %s\n", *message.Body)

		// Deletar mensagem da fila após processamento
		_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      queueURL,
			ReceiptHandle: message.ReceiptHandle,
		})
		if err != nil {
			return err
		}

		log.Println("Mensagem deletada da fila.")
	}

	return nil
}

func main() {
	sqsSvc := createSQSSession()

	queueURL, err := createQueue(sqsSvc)
	if err != nil {
		log.Fatal("Erro ao criar fila:", err)
		os.Exit(1)
	}

	messageBody := "Olá, LocalStack!"
	err = sendMessage(sqsSvc, queueURL, messageBody)
	if err != nil {
		log.Fatal("Erro ao enviar mensagem:", err)
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	err = receiveMessage(sqsSvc, queueURL)
	if err != nil {
		log.Fatal("Erro ao receber mensagem:", err)
		os.Exit(1)
	}
}
