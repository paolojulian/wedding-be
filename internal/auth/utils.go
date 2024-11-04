package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func setCookieHandler(w http.ResponseWriter, cookieValue string) {
	cookie := http.Cookie{
		Name:     authCookieName,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   int(tokenExpiryDuration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
}

func clearCookieHandler(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     authCookieName,
		Value:    "", // Clear the cookie value
		Path:     "/",
		MaxAge:   -1, // Set MaxAge to -1 to delete the cookie
		HttpOnly: true,
		Secure:   true,                  // Set to true to enforce HTTPS
		SameSite: http.SameSiteNoneMode, // Ensure SameSite is consistent
	}
	http.SetCookie(w, &cookie)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateJWT generates a JWT token for a user
func generateJWT(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(tokenExpiryDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}
