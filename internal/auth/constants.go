package auth

import (
	"errors"
	"time"
)

var (
	authCookieName        = "auth_token"
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	jwtSecretKey          = []byte("your_secret_key") // Replace with your secret key, preferably from env variables
	tokenExpiryDuration   = time.Hour * 72            // Token validity duration
)
