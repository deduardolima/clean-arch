package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/deduardolima/clean-arch/configs"
	"github.com/deduardolima/clean-arch/internal/event/handler"
	"github.com/deduardolima/clean-arch/internal/infra/graph"
	"github.com/deduardolima/clean-arch/internal/infra/grpc/pb"
	"github.com/deduardolima/clean-arch/internal/infra/grpc/service"
	"github.com/deduardolima/clean-arch/internal/infra/web/webserver"
	"github.com/deduardolima/clean-arch/pkg/events"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		configs.DBUser,
		configs.DBPassword,
		configs.DBHost,
		configs.DBPort,
		configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("OrdersListed", &handler.OrdersListedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db, eventDispatcher)

	// Iniciar servidor web
	go func() {
		webserver := webserver.NewWebServer(configs.WebServerPort)
		webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
		webserver.AddHandler("/order", webOrderHandler.Create)
		webserver.AddHandler("/orders", webOrderHandler.List)
		fmt.Println("Starting web server on port", configs.WebServerPort)
		webserver.Start()
	}()

	// Iniciar servidor gRPC
	go func() {
		grpcServer := grpc.NewServer()
		orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
		pb.RegisterOrderServiceServer(grpcServer, orderService)
		reflection.Register(grpcServer)

		fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
		if err != nil {
			panic(err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		}
	}()

	// Iniciar servidor GraphQL
	go func() {
		srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			CreateOrderUseCase: *createOrderUseCase,
			ListOrdersUseCase:  *listOrdersUseCase,
		}}))
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)

		fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
		if err := http.ListenAndServe(":"+configs.GraphQLServerPort, nil); err != nil {
			fmt.Printf("GraphQL server error: %v\n", err)
		}
	}()

	select {}
}

func getRabbitMQChannel() *amqp.Channel {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	var conn *amqp.Connection
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
			configs.RabbitMQUser,
			configs.RabbitMQPassword,
			configs.RabbitMQHost,
			configs.RabbitMQPort,
		))
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to RabbitMQ, retrying in 2 seconds... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("Could not connect to RabbitMQ: %v", err))
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
