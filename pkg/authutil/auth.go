package authutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/atmshang/nuclear-nest/pkg/logutil"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

const (
	headerInternalServiceAuth = "X-LincService-Auth"
	headerVerifiedByTraefik   = "X-Verified-By-Traefik"
	validateTime              = time.Second * time.Duration(10) // 十秒内有效
)

var (
	emptyAuthHeader   = errors.New("auth header is empty")
	invalidAuthHeader = errors.New("invalid auth header")
	expiredAuthHeader = errors.New("auth header is expired")
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	debugMode  = true
)

/*****************************************************************
*							调试模式
*****************************************************************/

// SetDebugMode 设置调试模式
func SetDebugMode(debug bool) {
	debugMode = debug
	logutil.Println("调试模式：", debugMode)
}

/*****************************************************************
*							密钥设置
*****************************************************************/

// SetPublicKey 设置公钥
func SetPublicKey(pemStr string) error {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return errors.New("failed to parse public key")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return errors.New("failed to parse public key")
	}

	var ok bool
	publicKey, ok = parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("not a valid RSA public key")
	}

	return nil
}

// SetPrivateKey 设置私钥
func SetPrivateKey(pemStr string) error {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return errors.New("failed to parse private key")
	}

	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("failed to parse private key")
	}

	privateKey = parsedPrivateKey
	return nil
}

/*****************************************************************
*							可信访问
*****************************************************************/

type AuthHeader struct {
	Expiration int64 `json:"expiration"` // 过期时间，UTC时间戳
}

// GenerateAuthHeaderValue 生成可信访问的请求头参数
func GenerateAuthHeaderValue() (string, string) {
	header := AuthHeader{
		Expiration: time.Now().Add(validateTime).UnixMilli(),
	}
	jsonBytes, err := json.Marshal(header)
	if err != nil {
		panic(err)
	}

	bytes, err := encryptRSA(jsonBytes)
	if err != nil {
		panic(err)
	}

	return headerInternalServiceAuth, base64.StdEncoding.EncodeToString(bytes)
}

func parseAuthHeaderValue(headerStr string) (AuthHeader, error) {
	var auth AuthHeader

	encryptedBytes, err := base64.StdEncoding.DecodeString(headerStr)
	if err != nil {
		return auth, err
	}
	jsonStrBytes, err := decryptRSA(encryptedBytes)
	if err != nil {
		return auth, err
	}

	err = json.Unmarshal(jsonStrBytes, &auth)
	if err != nil {
		return auth, err
	}
	return auth, nil
}

type verify struct {
	UserId    string    `json:"userId"`
	IsAdmin   bool      `json:"isAdmin"`
	Timestamp time.Time `json:"timestamp"`
}

func verifiedByTraefik(ctx *gin.Context) bool {
	verifiedStr := ctx.GetHeader(headerInternalServiceAuth)
	if len(verifiedStr) == 0 {
		logutil.Println("[verifiedByTraefik] 来自网关的请求头不存在")
		return false
	}
	var encryptedModel EncryptedData
	err := json.Unmarshal([]byte(verifiedStr), &encryptedModel)
	if err != nil {
		logutil.Println("[verifiedByTraefik] 加密的请求头反序列化失败")
		return false
	}
	bytes, err := DecryptAESString(encryptedModel)
	if err != nil {
		logutil.Println("[verifiedByTraefik] 加密的请求头解密失败")
		return false
	}
	var verifyModel verify
	err = json.Unmarshal([]byte(bytes), &verifyModel)
	if err != nil {
		logutil.Println("[verifiedByTraefik] 解密的鉴权内容反序列化失败")
		return false
	}
	return true
}

// InternalServiceAuth 内部服务间调用的认证中间件,若是经过traefik验证,则直接放行
func InternalServiceAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if debugMode {
			logutil.Println("[InternalServiceAuth] 调试模式放行")
		}

		if verifiedByTraefik(ctx) {
			ctx.Next()
			return
		}

		if verifyByAuthHeader(ctx) {
			ctx.Next()
			return
		}

		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func verifyByAuthHeader(ctx *gin.Context) bool {
	var header string
	header = ctx.GetHeader(headerVerifiedByTraefik)
	if len(header) == 0 {
		header = ctx.Query(headerVerifiedByTraefik)
	}
	if len(header) == 0 {
		_ = ctx.AbortWithError(http.StatusUnauthorized, emptyAuthHeader)
		logutil.Println("可信请求的字段不存在")
		return false
	}

	authHeader, err := parseAuthHeaderValue(header)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, invalidAuthHeader)
		logutil.Println("可信请求的解析失败")
		return false
	}

	if time.Now().After(time.UnixMilli(authHeader.Expiration)) {
		_ = ctx.AbortWithError(http.StatusUnauthorized, expiredAuthHeader)
		logutil.Println("可信请求已过期")
		return false
	}
	return true
}

/*****************************************************************
*							加密部分
*****************************************************************/

func encryptRSA(input []byte) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("public key is not set")
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, input)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

func decryptRSA(encryptedData []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is not set")
	}
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

const (
	aesKeyLength = 32
)

func EncryptAESString(input string) (EncryptedData, error) {
	// Generate a random AES key
	aesKey := make([]byte, aesKeyLength)
	_, err := rand.Read(aesKey)
	if err != nil {
		return EncryptedData{}, err
	}

	// Encrypt the input using AES with the generated IV
	encryptedData, encryptedNonce, err := encryptAES([]byte(input), aesKey)
	if err != nil {
		return EncryptedData{}, err
	}

	// Encrypt the AES key using RSA public key
	encryptedAESKey, err := encryptRSA(aesKey)
	if err != nil {
		return EncryptedData{}, err
	}

	encrypted := EncryptedData{
		Data:  base64.StdEncoding.EncodeToString(encryptedData),
		Key:   base64.StdEncoding.EncodeToString(encryptedAESKey),
		Nonce: encryptedNonce,
	}
	return encrypted, nil
}

func DecryptAESString(encryptedData EncryptedData) (string, error) {
	// Decode the base64-encoded data
	decodedEncryptedData, err := base64.StdEncoding.DecodeString(encryptedData.Data)
	if err != nil {
		return "", err
	}

	decodedEncryptedAESKey, err := base64.StdEncoding.DecodeString(encryptedData.Key)
	if err != nil {
		return "", err
	}

	decodedIv, err := base64.StdEncoding.DecodeString(encryptedData.Nonce)
	if err != nil {
		return "", err
	}

	// Decrypt the AES key using RSA private key
	aesKey, err := decryptRSA(decodedEncryptedAESKey)
	if err != nil {
		return "", err
	}

	// Decrypt the data using AES key and IV
	decryptedData, err := decryptAES(decodedEncryptedData, aesKey, decodedIv)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

func encryptAES(input []byte, key []byte) ([]byte, string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, "", err
	}

	return gcm.Seal(nil, nonce, input, nil), base64.StdEncoding.EncodeToString(nonce), nil
}

func decryptAES(encryptedData []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, iv, encryptedData, nil)
}

type EncryptedData struct {
	Data    string `json:"data"`
	Key     string `json:"key"`
	Nonce   string `json:"nonce"`
	Message string `json:"message"`
}
