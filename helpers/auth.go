package helpers

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("acvsd"), nil)

	// Test key
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"email": "test@gmail.com", "user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return tokenAuth
}

func GenerateTokenString(email string, id int) string {
	d, _ := time.ParseDuration("3h")
	a, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"email": email, "user_id": id,
		"exp": time.Now().UTC().Add(d)})
	fmt.Println(a.Expiration())
	return tokenString
}
