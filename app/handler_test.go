package main

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	ctx := context.Background()

	// Start the Redis container
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	defer redisC.Terminate(ctx)

	// Get the Redis container's host and port
	host, err := redisC.Host(ctx)
	require.NoError(t, err)
	port, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)

	// Connect the Redis client to the container
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port.Port(),
		Password: "", // No password for the test container
		DB:       0,
	})
	defer client.Close()

	// Test the application's logic
	app := &application{
		config: &config{address: host + ":" + port.Port(), password: "", db: 0},
		client: client,
	}

	// Run tests
	t.Run("Set and Get OTP", func(t *testing.T) {
		phoneNumber := "+251-123456789"
		otp := "123456"

		// Test SetOTP
		err := app.SetOTP(phoneNumber, otp)
		require.NoError(t, err)

		// Test GetOTP
		retrievedOTP, err := app.GetOTP(phoneNumber)
		require.NoError(t, err)
		require.Equal(t, otp, retrievedOTP)
	})
}
