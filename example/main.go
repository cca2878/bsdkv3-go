package main

import (
	"errors"
	"os"

	"bsdkv3-go/sdk"
	"bsdkv3-go/sdk/config"
	"bsdkv3-go/sdk/log"
)

// testValidator demonstrates the captcha validation functionality
func testValidator() {
	log.SetLevel(log.LevelDebug)
	log.Info("Testing captcha validator...")

	validator := bsdkv3.NewRemoteValidator()
	result, err := validator.Validate()
	if err != nil {
		log.Error("Captcha validation failed: %v", err)
		return
	}

	// Display validation results
	log.Info("Captcha validation successful!")
	log.Info("Challenge: %s", result.Challenge)
	log.Info("Gt: %s", result.Gt)
	log.Info("GtUserId: %s", result.GtUserId)
	log.Info("Validate: %s", result.Validate)
}

// testLogin demonstrates the login functionality with proper error handling
func testLogin() {
	log.SetLevel(log.LevelDebug)
	log.Info("Testing login functionality...")

	// Get credentials from environment variables or use defaults
	username := getEnvOrDefault("BSDK_USERNAME", "your_username")
	password := getEnvOrDefault("BSDK_PASSWORD", "your_password")

	if username == "your_username" || password == "your_password" {
		log.Warn("Using default credentials. Set BSDK_USERNAME and BSDK_PASSWORD environment variables for real testing.")
	}

	// Create user info
	user := bsdkv3.UserInfo{
		Username: username,
		Password: password,
		Platform: "android",
		Channel:  "bilibili",
	}

	// Create client with default configuration
	client, err := bsdkv3.NewClient(config.AppkeyPcr)
	if err != nil {
		handleLoginError(err)
		return
	}
	defer client.Close()

	// Attempt login
	account, err := client.Login(user)
	if err != nil {
		handleLoginError(err)
		return
	}

	// Login successful
	log.Info("Login successful!")
	log.Info("User ID: %s", account.Uid)
	log.Info("Access Key: %s", account.AccessKey[:10]+"...") // Only show first 10 chars for security
	log.Info("Platform: %s", account.Platform)
	log.Info("Channel: %s", account.Channel)
}

// testCustomConfig demonstrates using custom configuration
func testCustomConfig() {
	log.Info("Testing custom configuration...")

	// Create custom configuration
	customConfig := config.NewConfig(
		config.WithRequestTimeout(10), // 10 second timeout
		config.WithAppId("custom_app"),
		config.WithGameId("custom_game"),
		config.WithPlatform("ios"),
		config.WithChannelId("2"),
	)

	// Create client with custom config
	client, err := bsdkv3.NewClient(config.AppkeyPcr, bsdkv3.WithConfig(customConfig))
	if err != nil {
		log.Error("Failed to create client with custom config: %v", err)
		return
	}
	defer client.Close()

	log.Info("Client created successfully with custom configuration")
}

// handleLoginError provides specific error handling for different types of login errors
func handleLoginError(err error) {
	var clientErr *bsdkv3.ClientError
	var apiErr *bsdkv3.APIError
	var validationErr *bsdkv3.ValidationError

	switch {
	case errors.As(err, &validationErr):
		log.Error("Validation error: %v", validationErr)
		log.Error("Please check the field '%s' with value '%s'", validationErr.Field, validationErr.Value)

	case errors.As(err, &apiErr):
		log.Error("API error: %v", apiErr)
		log.Error("Error code: %s", apiErr.Code)

	case errors.As(err, &clientErr):
		log.Error("Client error in operation '%s': %v", clientErr.Op, clientErr)

		// Check for specific error types
		if errors.Is(err, bsdkv3.ErrInvalidCredentials) {
			log.Error("Invalid username or password provided")
		} else if errors.Is(err, bsdkv3.ErrCaptchaFailed) {
			log.Error("Captcha verification failed")
		} else if errors.Is(err, bsdkv3.ErrNetworkError) {
			log.Error("Network connection failed - check your internet connection")
		}

	default:
		log.Error("Unexpected error: %v", err)
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	log.Info("BSdk V3 Go Client Example")
	log.Info("========================")

	// Test custom configuration
	testCustomConfig()

	// Test login (will fail with default credentials due to network restrictions)
	testLogin()

	log.Info("Example completed. Note: Some operations may fail due to network restrictions or invalid credentials.")
}
