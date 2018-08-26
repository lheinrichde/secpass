package shorts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	// implement postgresql driver
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
)

// LogInfo log only error or everything
var LogInfo = true

// InitLogging Enable logging to file with automatic error printing
func InitLogging(logInfo bool) {
	// set LogInfo
	LogInfo = logInfo

	// create directories
	os.MkdirAll("logs", os.ModePerm)

	// create new file and check for error
	file, err := os.Create("logs/" + time.Now().Format("2006-01-02.txt"))
	defer file.Close()

	if err == nil {
		log.SetOutput(file)
	} else {
		fmt.Println(err)
	}
}

// InitLoggingFile Enable logging to specified file
func InitLoggingFile(fileName string, logInfo bool) error {
	// set LogInfo
	LogInfo = logInfo

	// split directories from filename and create them
	directory, _ := path.Split(fileName)
	os.MkdirAll(directory, os.ModePerm)

	// create new file and check for error
	file, err := os.Create(fileName)
	defer file.Close()
	if err == nil {
		log.SetOutput(file)
	}

	return err
}

// UUID generate random uuid and return it as a string
func UUID() string {
	// generate uuid
	id, _ := uuid.NewV4()
	return id.String()
}

// Check Print error if exists
func Check(err error, isError bool) {
	// check whether error exists
	if err != nil {
		// check whether to print or not
		if isError || LogInfo == true {
			// print error
			log.Println(err)
		}
	}
}

// ConnectPostgreSQL Connection to postgresql database
func ConnectPostgreSQL(host, port, database, username, password string, ssl string) *sql.DB {
	// open database connection and check for errors
	db, err := sql.Open("postgres", "postgres://"+username+":"+password+"@"+host+"/"+database+"?sslmode="+ssl)
	Check(err, true)
	Check(db.Ping(), true)

	// return database connection
	return db
}

// Hash return SHA-3 512 hash string
func Hash(input string) string {
	hasher := sha3.New512()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateKey generate AES key
func GenerateKey(key string) []byte {
	length := len(key)
	switch {
	case length >= 32:
		return []byte(key[:32])
	case length >= 24:
		return []byte(key[:24])
	case length >= 16:
		return []byte(key[:16])
	}

	return GenerateKey(key + "_secpass_key1gen")
}

// Encrypt text with key
func Encrypt(text string, key []byte) string {
	plain := []byte(text)

	block, err := aes.NewCipher(key)
	Check(err, true)

	cipherText := make([]byte, aes.BlockSize+len(plain))
	iv := cipherText[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	Check(err, true)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plain)

	return base64.URLEncoding.EncodeToString(cipherText)
}

// Decrypt text with key
func Decrypt(text string, key []byte) string {
	cipherText, err := base64.URLEncoding.DecodeString(text)
	Check(err, true)

	block, err := aes.NewCipher(key)
	Check(err, true)

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}
