package auth

import (
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Test the checkPasswordHash function
func TestCheckPasswordHash(t *testing.T) {
	password := "qwe123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}

	log.Printf("Hash: %s", hash)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     string(hash),
			want:     true,
		},
		{
			name:     "Incorrect password",
			password: "wrong-password",
			hash:     string(hash),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPasswordHash(tt.password, tt.hash); got != tt.want {
				t.Errorf("checkPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
