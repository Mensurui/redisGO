package main

import (
	"bufio"
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
	var cfg config

	cfg.address = "redis-11913.c98.us-east-1-4.ec2.redns.redis-cloud.com:11913"
	cfg.password = "cPtQ3slOp3GLelXM4mVpnyy0jcMpvOUZ"
	cfg.db = 0

	set := flag.Bool("s", false, "Use if you want to set the otp sent to you")
	get := flag.Bool("g", false, "Use if you want to get your data")
	all := flag.Bool("a", false, "Use if you want all the data")
	flag.Parse()

	client := redisConnect(cfg)

	app := application{
		config: &cfg,
		client: client,
	}

	switch {
	case *set:
		fmt.Println("Enter you phone number(+251-XXX..)")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()

		if err != nil {
			fmt.Println("Error scanning")
			os.Exit(1)
		}

		phoneNumber := scanner.Text()
		fmt.Printf("Your phone number is: %s\n", phoneNumber)

		fmt.Println("Enter your password")
		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			fmt.Println("Error scanning")
			os.Exit(1)
		}
		otp := scanner.Text()
		fmt.Println("OTP Set")

		err = app.SetOTP(phoneNumber, otp)
		if err != nil {
			panic(err)
		}

	case *get:
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Println("Error scanning")
			os.Exit(1)
		}
		phoneNumber := scanner.Text()
		val, err := app.GetOTP(phoneNumber)
		if err != nil {
			panic(err)
		}
		fmt.Printf("your OTP is:%v\t\n", val)
	case *all:
		data, err := app.GetAll()
		if err != nil {
			fmt.Printf("Error fetching data due to: %v\n", err)
			break // Exit the case if there's an error
		}

		// Print user data in a readable format
		for _, user := range data {
			fmt.Printf("Phone Number: %s, OTP: %s\n", user.PhoneNumber, user.OTP)
		}

	}
}

func redisConnect(cfg config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.address,
		Password: cfg.password,
		DB:       cfg.db,
	})

	return client
}
