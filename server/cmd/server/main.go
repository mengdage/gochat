package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"githum.com/mengdage/gochat/pkg/api"
	"githum.com/mengdage/gochat/pkg/app"
	"githum.com/mengdage/gochat/pkg/storage"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("Failed to run: %v", err)
		os.Exit(1)
	}
}

func run() error {
	initConfig()
	initLog()

	cache := createCache()

	dbPool, err := createPGPool()
	if err != nil {
		return err
	}
	dbs := storage.NewDBStorage(dbPool)

	storage := storage.NewStorage(dbs, cache)

	clientManager := api.NewClientManager()
	go clientManager.Start()

	userService := api.NewUserService(storage)
	wsService := api.NewWSService(storage, clientManager)
	convService := api.NewConversationService(storage)

	messageOperator := api.NewMessageOperator(clientManager, storage)

	router := gin.Default()

	server := app.NewServer(router, userService, wsService, convService, clientManager, messageOperator)

	err = server.Run(*addr)
	if err != nil {
		return err
	}

	return nil
}

func initConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	fmt.Println("app configurations:")
	for k, v := range viper.GetStringMapString("app") {
		fmt.Printf("\t%s: %s\n", k, v)
	}
}

func initLog() {
	gin.DisableConsoleColor()

	logFilePath := viper.GetString("app.log_file")
	f, err := os.Create(logFilePath)

	if err != nil {
		panic(fmt.Errorf("Failed to create logFile %s: %s", logFilePath, err.Error()))
	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetOutput(gin.DefaultWriter)
}

func createPGPool() (*pgxpool.Pool, error) {
	pgURL := viper.GetString("pg.url")

	pool, err := pgxpool.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatalf("Failed to connect to postgreSQL %s: %v", pgURL, err)
		return nil, err
	}

	return pool, nil

}

func createCache() *cache.Cache {
	client := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolsize"),
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),
	})

	c := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return c
}
