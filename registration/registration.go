package registration

import (
	"Databases/postgresql"
	cryptorand "crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	mathrand "math/rand"
	"net/http"

	"github.com/jackc/pgtype"
)

type regAnswer struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type newUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	pwSaltBytes = 32
	rcBytes     = 5
	charset     = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func makeSalt() []byte {
	salt := make([]byte, pwSaltBytes)
	_, err := io.ReadFull(cryptorand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	return salt
}

func makeRefCode() string {
	brefCode := make([]byte, rcBytes)
	for i := range brefCode {
		brefCode[i] = charset[mathrand.Intn(len(charset))]
	}
	for checkRefCode(string(brefCode)) {
		for i := range brefCode {
			brefCode[i] = charset[mathrand.Intn(len(charset))]
		}
	}
	return string(brefCode)
}

func hashPassword(sourcePassword string, salt []byte) string {
	hasher := sha512.New()
	hasher.Write([]byte(sourcePassword))
	hash1 := hex.EncodeToString(hasher.Sum(nil))
	hasher.Write(append([]byte(hash1), salt...))
	hash2 := hex.EncodeToString(hasher.Sum(nil))
	return hash2
}

//RegistrationHandler for Registration user
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	var err error

	var newUser newUser
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Println(newUser.Password)
		log.Println(err)
	}
	username := newUser.Username
	password := newUser.Password

	if username == "" || password == "" {
		jsonResponse, err = json.Marshal(regAnswer{"Bad", "Empty username or password"})
		w.Write(jsonResponse)
		return
	}
	if len(password) < 8 {
		jsonResponse, err = json.Marshal(regAnswer{"Bad", "Please make safe password"})
		w.Write(jsonResponse)
		return
	}
	if checkUser(username) {
		jsonResponse, err = json.Marshal(regAnswer{"Bad", "User already exist, please take another nickname"})
	} else {
		salt := makeSalt()
		refcode := makeRefCode()
		password = hashPassword(password, salt)
		jsonResponse, err = json.Marshal(addUserToDatabase(username, password, salt, refcode, false, false, 0))
	}
	if err != nil {
		panic(err)
	}
	w.Write(jsonResponse)
}

func addUserToDatabase(username string, password string, salt []byte, refcode string, showAd bool, isPremium bool, inv_pers int8) regAnswer {
	var saltArray pgtype.Bytea
	saltArray.Set(salt)
	params := postgresql.NewParams("INSERT INTO accounts VALUES (default, $1, $2, $3, $4, $5, $6, $7)", []interface{}{username, password, saltArray, refcode, showAd, isPremium, inv_pers})
	result, err := postgresql.RunQuery(params)
	result.Close()
	if err != nil {
		return regAnswer{"Bad", "Error while adding user"}
	}
	return regAnswer{"Ok", "User added"}
}

func checkUser(username string) bool {
	params := postgresql.NewParams("SELECT id FROM accounts WHERE username=$1", []interface{}{username})
	result, _ := postgresql.RunQuery(params)
	return result.Next()
}

func checkRefCode(refcode string) bool {
	params := postgresql.NewParams("SELECT id FROM accounts WHERE refcode=$1", []interface{}{refcode})
	result, _ := postgresql.RunQuery(params)
	return result.Next()
}
