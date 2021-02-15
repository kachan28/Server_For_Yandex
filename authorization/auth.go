package authorization

import (
	"Databases/postgresql"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type curUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type regAnswer struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var mySecret = []byte("secret")

//AuthHandler for Authorization user
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	var curUser curUser
	var err error
	err = json.NewDecoder(r.Body).Decode(&curUser)
	salt := getSalt(curUser)
	curUser.Password = hashPassword(curUser.Password, salt)
	if err != nil {
		log.Println(err)
	}
	isuserValid := checkUser(curUser)
	if isuserValid {
		t, rt, err := generateTokenPair(curUser.Username)
		if err != nil {
			jsonResponse, err = json.Marshal(regAnswer{"Bad", err.Error()})
			w.Write(jsonResponse)
			return
		} else {
			jsonResponse, err = json.Marshal(
				map[string]string{
					"access_token":  *t,
					"refresh_token": *rt,
				})
		}
	} else {
		jsonResponse, err = json.Marshal(regAnswer{"Bad", "Wrong username or password, please, retry"})
	}
	w.Write(jsonResponse)
}

//RefreshToken for refreshing jwtoken
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySecret, nil
		})
		if error != nil {
			jsonResponse, err := json.Marshal(regAnswer{"Bad", error.Error()})
			if err != nil {
				log.Println(err)
			} else {
				w.Write(jsonResponse)
			}
			return
		}
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			if claims["method"] == "refresh" {
				t, rt, err := generateTokenPair(claims["name"].(string))
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				jsonResponse, err := json.Marshal(map[string]string{
					"access_token":  *t,
					"refresh_token": *rt,
				})
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				w.Write(jsonResponse)
			} else {
				jsonResponse, err := json.Marshal(regAnswer{"Bad", "Token not for refresh"})
				if err != nil {
					w.Write([]byte(err.Error()))
				}
				w.Write(jsonResponse)
			}
		}
	} else {
		jsonResponse, err := json.Marshal(regAnswer{"Bad", "No token"})
		if err != nil {
			log.Println(err)
		} else {
			w.Write(jsonResponse)
		}
	}
}

//JWTAccessHandler for processing jwtoken
func JWTAccessHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySecret, nil
			})
			if error != nil {
				log.Println(error)
				jsonResponse, err := json.Marshal(regAnswer{"Bad", error.Error()})
				if err != nil {
					log.Println(err)
				} else {
					w.Write(jsonResponse)
				}
				return
			}
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				if claims["method"] == "access" {
					next.ServeHTTP(w, r)
				} else {
					jsonResponse, err := json.Marshal(regAnswer{"Bad", "Token not for access"})
					if err != nil {
						w.Write([]byte(err.Error()))
					}
					w.Write(jsonResponse)
				}
			}
		} else {
			jsonResponse, err := json.Marshal(regAnswer{"Bad", "No token"})
			if err != nil {
				log.Println(err)
			} else {
				w.Write(jsonResponse)
			}
		}
	})
}

func hashPassword(sourcePassword string, salt []byte) string {
	hasher := sha512.New()
	hasher.Write([]byte(sourcePassword))
	hash1 := hex.EncodeToString(hasher.Sum(nil))
	hasher.Write(append([]byte(hash1), salt...))
	hash2 := hex.EncodeToString(hasher.Sum(nil))
	return hash2
}

func getSalt(curUser curUser) []byte {
	params := postgresql.NewParams("SELECT salt FROM accounts WHERE username=$1", []interface{}{curUser.Username})
	result, _ := postgresql.RunQuery(params)
	var salt []byte
	for result.Next() {
		result.Scan(&salt)
	}
	return salt
}

func checkUser(curUser curUser) bool {
	params := postgresql.NewParams("SELECT id FROM accounts WHERE username=$1 and password=$2", []interface{}{curUser.Username, curUser.Password})
	result, _ := postgresql.RunQuery(params)
	return result.Next()
}

func generateTokenPair(username string) (*string, *string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["method"] = "access"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // Generate encoded token and send it as response.
	t, err := token.SignedString(mySecret)

	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["name"] = username
	rtClaims["method"] = "refresh"
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString(mySecret)

	if err != nil {
		return nil, nil, err
	}

	return &t, &rt, nil
}
