package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"goexpert-clean-architecture/configs"
	event_handler "goexpert-clean-architecture/internal/event/handler"
	"goexpert-clean-architecture/internal/infra/graph"
	"goexpert-clean-architecture/internal/infra/grpc/pb"
	"goexpert-clean-architecture/internal/infra/grpc/service"
	"goexpert-clean-architecture/internal/infra/rest"
	usecase "goexpert-clean-architecture/internal/usecase/order"
	"goexpert-clean-architecture/pkg/events"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db := connectToDatabase()
	defer db.Close()

	applyDatabaseMigrations(db)

	rabbitMQChannel := connectToRabbitMQ()

	eventDispatcher := setupEventDispatcher(rabbitMQChannel)

	createOrderUseCase := setupCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := setupListOrderUseCase(db)

	startRestServer(cfg.WebServerPort, db, eventDispatcher)
	startGRPCServer(cfg.GRPCServerPort, *createOrderUseCase, *listOrderUseCase)
	startGraphQLServer(cfg.GraphQLServerPort, *createOrderUseCase, *listOrderUseCase)
}

func connectToDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	return db
}

func applyDatabaseMigrations(db *sql.DB) {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("./migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}

func connectToRabbitMQ() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func setupEventDispatcher(ch *amqp.Channel) *events.EventDispatcher {
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &event_handler.OrderCreatedHandler{
		RabbitMQChannel: ch,
	})
	return eventDispatcher
}

func setupCreateOrderUseCase(db *sql.DB, eventDispatcher *events.EventDispatcher) *usecase.CreateOrderUseCase {
	return NewCreateOrderUseCase(db, eventDispatcher)
}

func setupListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	return NewListOrderUseCase(db)
}

func startRestServer(port string, db *sql.DB, eventDispatcher *events.EventDispatcher) {
	ws := rest.NewServer(port)
	orderPath := "/order"
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	ws.AddHandler(rest.NewRoute(orderPath, "POST", webOrderHandler.Create))
	ws.AddHandler(rest.NewRoute(orderPath, "GET", webOrderHandler.GetOrders))
	fmt.Println("Starting REST server on port", port)
	go ws.Start()
}

func startGRPCServer(port string, createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) {
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(createOrderUseCase, listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
}

func startGraphQLServer(port string, createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CreateOrderUseCase: createOrderUseCase,
			ListOrderUseCase:   listOrderUseCase,
		},
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start GraphQL server: %v", err)
	}
}
