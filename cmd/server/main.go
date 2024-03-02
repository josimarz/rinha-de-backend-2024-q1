package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/josimarz/rinha-de-backend-2024-q1/configs"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/db/pgsql"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/gateway"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/handler"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/usecase"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	cfg                *configs.Config
	redisClient        *redis.Client
	rs                 *redsync.Redsync
	db                 *sql.DB
	dbGateway          gateway.DatabaseGateway
	getBankStmtUC      *usecase.GetBankStatementUseCase
	doTransactionUC    *usecase.DoTransactionUseCase
	bankStmtHandler    *handler.BankStatementHandler
	transactionHandler *handler.TransactionHandler
)

func main() {
	loadConfig()
	connectToRedis()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("unable to close Redis connection: %v", err)
		}
	}()
	connectToDatabase()
	defer db.Close()
	initDbGateway()
	initUseCases()
	initHandlers()
	startServer()
}

func loadConfig() {
	var err error
	cfg, err = configs.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
}

func connectToRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("unable to connect to Redis: %v", err)
	}
	pool := goredis.NewPool(redisClient)
	rs = redsync.New(pool)
}

func connectToDatabase() {
	var err error
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func initDbGateway() {
	dbGateway = pgsql.NewPostgreDatabaseGateway(db)
}

func initUseCases() {
	getBankStmtUC = usecase.NewGetBankStatementUseCase(dbGateway)
	doTransactionUC = usecase.NewDoTransactionUseCase(dbGateway, rs)
}

func initHandlers() {
	bankStmtHandler = handler.NewBankStatementHandler(getBankStmtUC)
	transactionHandler = handler.NewTransactionHandler(doTransactionUC)
}

func startServer() {
	mux := http.NewServeMux()
	mux.Handle("GET /clientes/{id}/extrato", bankStmtHandler)
	mux.Handle("POST /clientes/{id}/transacoes", transactionHandler)
	http.ListenAndServe(":8080", mux)
}
