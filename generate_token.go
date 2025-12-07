package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {

	// ! 2 ortam için de 2 farklı token üretmemiz gerekiyor.

	dockerJwt := "jwt-secret"       // burada docker ortamı için bir json jwt token üretiyoruz
	localJwt := "my-256-bit-secret" // burada da local ortam için ayrıca bir json jwt token üretmemiz gerekiyor.

	generateToken("Docker", dockerJwt)
	generateToken("Local", localJwt)
}

// 24 saat geçerli bir jwt token üretiyoruz.
func generateToken(envName, jwtToken string) {

	claims := jwt.MapClaims{
		"userId":  "testUser1",
		"role":    "admin",
		"expTime": time.Now().Add(time.Hour * 24).Unix(), // TOKEN SÜRESİ = 24 SAAT
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // token burada oluşturulur. HS256 (HMAC-SHA256)
	signedToken, err := token.SignedString([]byte(jwtToken))   // token sign(imzalama). token başarısız olursa err döner.
	if err != nil {
		fmt.Printf("Hata %s: %v\n", envName, err)
	}

	fmt.Printf("Ortam: %s\n", envName)
	fmt.Printf("Bearer: %s\n", signedToken)
}
