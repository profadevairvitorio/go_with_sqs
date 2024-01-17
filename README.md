# Go + AWS SQS LocalStack

## Configuração

Certifique-se de ter o LocalStack instalado e em execução antes de começar.

1. Instale o SDK AWS Go:

   ```bash
   go get -u github.com/aws/aws-sdk-go
   ```
2. Configure o ambiente:
   ```bash
    LocalStack Endpoint: http://localhost:4566
    Região: us-east-1 
    SSL: false
   ```
3. Suba o docker
   ```bash
   docker compose up
   ```
4. Execute 
    ```bash
    go run main.go
    ```

### Funcionalidades
1)createSQSSession()

Cria e retorna uma instância do serviço SQS configurada para o LocalStack.

2)createQueue(svc *sqs.SQS) (*string, error)

Cria uma fila SQS e retorna a URL da fila.

3)sendMessage(svc *sqs.SQS, queueURL *string, messageBody string) error

Envia uma mensagem para a fila especificada.

4)receiveMessage(svc *sqs.SQS, queueURL *string) error

Recebe e processa mensagens da fila.
