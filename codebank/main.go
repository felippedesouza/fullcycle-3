package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/felippedesouza/fullcycle3-codebank/infrastructure/grpc/server"
	"github.com/felippedesouza/fullcycle3-codebank/infrastructure/kafka"
	"github.com/felippedesouza/fullcycle3-codebank/infrastructure/repository"
	"github.com/felippedesouza/fullcycle3-codebank/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := setupDb()
	defer db.Close()
	producer := setupKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	fmt.Println("rodando grpc server")
	serveGrpc(processTransactionUseCase)

	// cc := domain.NewCreditCard()
	// cc.Number = "1234"
	// cc.Name = "Felippe"
	// cc.ExpirationYear = 2021
	// cc.ExpirationMonth = 9
	// cc.CVV = 123
	// cc.Limit = 1000
	// cc.Balance = 0

	// repo := repository.NewTransactionRepositoryDb(db)
	// err := repo.CreateCreditCard(*cc)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer("host.docker.internal:9094")
	return producer
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	usecase := usecase.NewUseCaseTransaction(transactionRepository)
	usecase.KafkaProducer = producer
	return usecase
}

func setupDb() *sql.DB {
	psqlinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}

func serveGrpc(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUseCase
	grpcServer.Serve()
}
