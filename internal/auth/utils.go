package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/config"
	"golang.org/x/crypto/bcrypt"
)

func setCookieHandler(c *gin.Context, cookieValue string) {
	// Get the origin domain
	origin := c.Request.Header.Get("Origin")

	// Parse the origin URL to get the domain
	cookieDomain := config.GetCookieDomain()

	// In production, set this to your specific domain
	// Example: if frontend is on app.example.com and backend is on api.example.com
	// cookieDomain should be ".example.com" to work across subdomains

	c.SetCookie(
		authCookieName,                     // name
		cookieValue,                        // value
		int(tokenExpiryDuration.Seconds()), // max age in seconds
		"/",                                // path
		cookieDomain,                       // domain (important for cross-origin)
		true,                               // secure
		true,                               // httpOnly
	)

	// Set SameSite attribute in header manually since Gin doesn't support it directly
	c.Header("Set-Cookie", c.Writer.Header().Get("Set-Cookie")+"; SameSite=None")

	// Set CORS headers
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}

func clearCookieHandler(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")

	// Parse the origin URL to get the domain
	cookieDomain := config.GetCookieDomain()

	c.SetCookie(
		authCookieName,
		"",
		-1,
		"/",
		cookieDomain,
		true,
		true,
	)

	// Set SameSite attribute in header manually
	c.Header("Set-Cookie", c.Writer.Header().Get("Set-Cookie")+"; SameSite=None")

	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}

// func clearCookieHandler(w http.ResponseWriter) {
// 	cookie := http.Cookie{
// 		Name:     authCookieName,
// 		Value:    "", // Clear the cookie value
// 		Path:     "/",
// 		MaxAge:   -1, // Set MaxAge to -1 to delete the cookie
// 		HttpOnly: true,
// 		Secure:   true,                  // Set to true to enforce HTTPS
// 		SameSite: http.SameSiteNoneMode, // Ensure SameSite is consistent
// 	}
// 	http.SetCookie(w, &cookie)
// }

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
