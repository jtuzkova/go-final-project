package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var (
    secret = []byte("gBElG5NThZSye")
	ExpectedPass string
)

type Request struct {
	Password string `json:"password"`
}

type Claims struct {
    PassHash string `json:"pass_hash"`
    jwt.RegisteredClaims
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"),  http.StatusMethodNotAllowed)
		return
	}

	var req Request
	var buf bytes.Buffer
	
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, fmt.Errorf("failed to read request body"), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if ExpectedPass == "" {
		writeError(w, fmt.Errorf("password not set"),  http.StatusInternalServerError)
		return
	}

	if req.Password != ExpectedPass {
		writeError(w, fmt.Errorf("incorrect password"), http.StatusUnauthorized)
		return
	}

	token, err := CreateToken(req.Password)
	if err != nil {
		writeError(w, fmt.Errorf("failed to generate token"), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"token": token})
}

func CreateToken(pass string) (string, error) {
	hashedPassword := sha256.Sum256([]byte(pass))
	hashStringPassword := hex.EncodeToString(hashedPassword[:])

	claims := &Claims{
        PassHash: hashStringPassword,
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, claims *Claims) bool {
    currentHash := sha256.Sum256([]byte(ExpectedPass))
    currentPassHash := hex.EncodeToString(currentHash[:])
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("incorrect sign method")
        }
        return []byte(secret), nil
    })
    if err != nil || !token.Valid {
        return false
    }
    return claims.PassHash == currentPassHash
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("TODO_PASSWORD: '%s'\n", ExpectedPass)

		if ExpectedPass == "" {
			log.Println("Пароль не установлен - пропускаем без проверки")
			next(w, r)
			return
		}
		
        if len(ExpectedPass) > 0 {
            var jwt string  
            cookie, err := r.Cookie("token")
            if err == nil {
                jwt = cookie.Value
            }

            if jwt == "" {
				writeError(w, fmt.Errorf("Authentication required"), http.StatusUnauthorized)
				return
			}

			claims := &Claims{}
            if !ValidateToken(jwt, claims) {
                writeError(w, fmt.Errorf("failed jdhvfjer"), http.StatusUnauthorized)
                return
            }
        }
        next(w, r)
    })
}