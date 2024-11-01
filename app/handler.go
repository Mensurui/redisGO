package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type User struct {
	PhoneNumber string `redis:"phoneNumber"`
	OTP         string `redis:"otp"`
}

func (app *application) SetOTP(phoneNumber, otp string) error {
	ctx := context.Background()
	err := app.client.HSet(ctx, phoneNumber, map[string]interface{}{
		"otp": otp,
	}).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (app *application) GetOTP(phoneNumber string) (string, error) {
	ctx := context.Background()
	otp, err := app.client.HGet(ctx, phoneNumber, "otp").Result()
	if err == redis.Nil {
		return "", fmt.Errorf("no OTP found for phone number: %s", phoneNumber)
	} else if err != nil {
		return "", err
	}
	return otp, nil
}

func (app *application) GetAll() ([]User, error) {
	ctx := context.Background()
	var users []User

	// Fetch all keys matching the user pattern (assuming phone numbers are the keys).
	keys, err := app.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving keys: %w", err)
	}

	// Iterate over the keys and get the user data for each phone number.
	for _, key := range keys {
		var user User
		err := app.client.HGetAll(ctx, key).Scan(&user)
		if err == redis.Nil {
			// No data found for this key, continue to next.
			continue
		} else if err != nil {
			return nil, fmt.Errorf("error retrieving user data for phone number %s: %w", key, err)
		}
		user.PhoneNumber = key // Store the key (phone number) in the user struct.
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no users found in data store")
	}

	return users, nil
}
