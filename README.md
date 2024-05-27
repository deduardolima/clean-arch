# Clean Architecture

## Descrição

Este projeto implementa clean architecture na linguagem Go, utilizando três servidores (Web, gRPC e GraphQL).


Os três servidores estão rodando em goroutines nas seguintes portas:
- Web: 8000
- gRPC: 50051
- GraphQL: 8080

## Estrutura do Projeto

- Clean-Arch/
  - cmd/
    - server/
      - main.go
  - configs/
    - config.go   
  - internal/
    - entity/
      - order.go
    - event/
      - handler/
        - event_handler.go
        - order_created_handler.go
        - orders_listed_handler.go
      - order_created.go
      - orders_listed.go
    - infra/
      - database/
        - interface.go
        - order_repository.go
      - graph/
        - schema.graphqls
        - resolver.go
      - grpc/
        - pb/
          - order.pb.go
        - service/
          - order_service.go
      - web/
        - webserver/
          - webserver.go
        - order_handler.go
    - usecase/
      - create_order.go
      - list_orders.go
  - pkg/
     - events/
      - event_dispatcher_test.go
      - event_dispatcher.go
      - interface.go       
  - sql/
    - migrations/
      - 000001_init.up.sql
      - 000001_init.down.sql
  - .env  
  - Dockerfile
  - docker-compose.yml
  - go.mod
  - go.sum
  - gqlgen.yml
  - Makefile
  - README.md
  - tools.go

## Pré-requisitos

- Docker
- Docker Compose

## Configuração
```
git clone https://github.com/deduardolima/clean-arch.git
cd clean-arch

```

## Instalação e Execução com Docker
Construa e inicie os containers:
```
docker-compose up --build -d
```

isso irá construir a imagem do aplicativo e iniciar os serviços definidos no docker-compose.yml, incluindo o banco de dados, rabbitMQ e o aplicativo.

## Criação das Tabelas

Após iniciar os containers, é necessário criar as tabelas no banco de dados. Siga os comandos abaixo:

Entre no container da aplicação:


```
docker exec -it go-app /bin/sh

```
Crie as tabelas com o comando:
```
make migrate
```
Para sair do container:
```
exit
```


## Execucão 

Agora é só utilizar!
- Para GraphQL: 

```
http://localhost:8080
```

- Para Webserver: 

```
http://localhost:8000
```

- Para gRPC, é necessário ter instalado o cliente Evans. Instale usando:

```
go install github.com/ktr0731/evans@latest

```

#### Usando Evans para gRPC
-  Para acessar o evans:

```
 evans --host localhost --port 50051 -r repl --proto internal/infra/grpc/protofiles/order.proto
```
- Dentro do REPL do Evans, selecione o pacote:
```
package pb
```

- Acesse os serviços disponíveis:

```
show service
```



### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis exemplo :

- DB_DRIVER=mysql
- DB_USER=root
- DB_PASSWORD=root
- DB_HOST=mysql
- DB_PORT=3306
- DB_NAME=orders
- WEB_SERVER_PORT=8000
- GRPC_SERVER_PORT=50051
- GRAPHQL_SERVER_PORT=8080
- RABBITMQ_USER=guest
- RABBITMQ_PASSWORD=guest
- RABBITMQ_HOST=rabbitmq
- RABBITMQ_PORT=5672

