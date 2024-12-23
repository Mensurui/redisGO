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
		return fmt.Errorf("failed to set OTP: %w", err)
	}
	return nil
}

func (app *application) GetOTP(phoneNumber string) (string, error) {
	ctx := context.Background()
	otp, err := app.client.HGet(ctx, phoneNumber, "otp").Result()
	if err == redis.Nil {
		return "", fmt.Errorf("no OTP found for phone number: %s", phoneNumber)
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve OTP: %w", err)
	}
	return otp, nil
}
