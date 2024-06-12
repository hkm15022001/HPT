package main

import (
	"log"
	"os"
	"sync"

	grpc "github.com/hkm15022001/Supply-Chain-Event-Management/api/grpc"
	"github.com/hkm15022001/Supply-Chain-Event-Management/api/middleware"
	httpServer "github.com/hkm15022001/Supply-Chain-Event-Management/api/server"
	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/handler"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if os.Getenv("RUNENV") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Initial web auth middleware
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		runWebAuth()
	}

	// Select app auth middleware
	if os.Getenv("RUN_APP_AUTH") == "redis" {
		runAppAuthRedis()
	} else if os.Getenv("RUN_APP_AUTH") == "buntdb" {
		log.Println("Selected BuntDB to run app auth")
	}

	dbRW, dbRO := connectPostgress()
	Store := handler.NewDataBaseConnection(dbRW, dbRO)

	// if err := handler.RefreshDatabase(Store); err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// log.Print("Database refreshed!")

	// WaitGroup to both server close
	var wg sync.WaitGroup
	wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	brokerList := []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")}
	// 	order_topic := "order-topic"
	// 	kafka.StartProducer(brokerList, order_topic)
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	brokerList := []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")}
	// 	longship_topic := "longship-topic"
	// 	kafka.StartConsumer(brokerList, longship_topic)
	// }()

	//  HTTP server
	go func() {
		defer wg.Done()
		httpServer.RunServer(Store)
	}()

	//  gRPC server
	go func() {
		defer wg.Done()
		grpc.RunServer(os.Getenv("GRPC_URL"))
	}()
	//wait all goroutine finish
	wg.Wait()
}

// Source code: https://www.devdungeon.com/content/working-files-go#read_all
func runWebAuth() {
	sessionKey := []byte(os.Getenv("SESSION_KEY"))
	if err := middleware.RunWebAuth(sessionKey); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Web authenticate activated!")
}

func runAppAuthRedis() {
	if err := middleware.RunAppAuthRedis(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Selected Redis to run app auth!")
}

func connectPostgress() (dbRW *gorm.DB, dbRO *gorm.DB) {
	var err error
	dbRO, err = ConnectReadOnlyPostgres()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	dbRW, err = ConnectReadWritePostgres()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Connected with posgres database!")
	return dbRW, dbRO
}

// ConnectPostgres to open connect with leader database
func ConnectReadOnlyPostgres() (db *gorm.DB, err error) {
	dbRO, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("READ_ONLY_POSTGRES_DSN"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return dbRO, nil

}

// ConnectPostgres to open connect with replica database
func ConnectReadWritePostgres() (db *gorm.DB, err error) {
	dbRW, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("READ_WRITE_POSTGRES_DSN"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return dbRW, nil

}
