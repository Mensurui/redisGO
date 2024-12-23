package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type config struct {
	address  string
	password string
	db       int
}

type application struct {
	config *config
	client *redis.Client
}

func main() {
	cfg := loadConfig()

	client := redisConnect(cfg)
	app := &application{
		config: &cfg,
		client: client,
	}

	runApp(app)
}

func loadConfig() config {
	return config{
		address:  "redis-11913.c98.us-east-1-4.ec2.redns.redis-cloud.com:11913",
		password: "cPtQ3slOp3GLelXM4mVpnyy0jcMpvOUZ",
		db:       0,
	}
}

func redisConnect(cfg config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.address,
		Password: cfg.password,
		DB:       cfg.db,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}

	return client
}

func runApp(app *application) {
	set := flag.Bool("s", false, "Use if you want to set the OTP sent to you")
	get := flag.Bool("g", false, "Use if you want to get your data")
	flag.Parse()

	if *set {
		handleSetOTP(app)
	} else if *get {
		handleGetOTP(app)
	} else {
		fmt.Println("Invalid flag. Use -s to set OTP or -g to get OTP.")
	}
}

func handleSetOTP(app *application) {
	phoneNumber := getInput("Enter your phone number (+251-XXX..)")
	otp := getInput("Enter your OTP")

	if err := app.SetOTP(phoneNumber, otp); err != nil {
		fmt.Printf("Failed to set OTP: %v\n", err)
		return
	}
	fmt.Println("OTP successfully set.")
}

func handleGetOTP(app *application) {
	phoneNumber := getInput("Enter your phone number (+251-XXX..)")
	otp, err := app.GetOTP(phoneNumber)
	if err != nil {
		fmt.Printf("Failed to retrieve OTP: %v\n", err)
		return
	}
	fmt.Printf("Your OTP is: %s\n", otp)
}

func getInput(prompt string) string {
	fmt.Println(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning input: %v\n", err)
		os.Exit(1)
	}
	return scanner.Text()
}
