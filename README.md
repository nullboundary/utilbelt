# utilbelt
#### A grab bag of various useful functions for golang.
--
    import "github.com/nullboundary/utilbelt"


## Usage

#### func  CalcHMAC

```go
func CalcHMAC(message string, key []byte) []byte
```
CalcHMAC makes an HMAC of a message and key

--
#### func  DebugPrintJson

```go
func DebugPrintJson(Data interface{})
```
DebugPrintJson prints json data for debugging

--
#### func  DecodeUriToBytes

```go
func DecodeUriToBytes(str string, fileType string) (string, []byte)
```
DecodeUriToBytes decodes a data uri into bytes

--
#### func  DecryptAESCFB

```go
func DecryptAESCFB(key []byte, cryptoText string) string
```
DecryptAESCFB decrypts from base64 to decrypted string

--
#### func  EncodetoDataUri

```go
func EncodetoDataUri(fileName string, mimeType string, allowTypes ...string) (string, error)
```
EncodetoDataUri reads a file and generates a data uri

--
#### func  EncryptAESCFB

```go
func EncryptAESCFB(key []byte, text string) string
```
EncryptAESCFB encrypts string to base64 crypto using AES

--
#### func  GenerateFnvHashID

```go
func GenerateFnvHashID(hashSeeds ...string) uint32
```
GenerateFnvHashID generates a hash fnv1a hash. Fast, unique, but insecure! use
only for ids and such.
https://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed

--
#### func  GenerateToken

```go
func GenerateToken(key []byte, seeds ...string) string
```
GenerateToken creates a urlEnocded hmac token

--
#### func  GetEtcdKey

```go
func GetEtcdKey(key string) (string, error)
```
GetEtcdKey retrives key/value pairs from etcd disrtibuted store

--
#### func  HashSha1Bytes

```go
func HashSha1Bytes(hashBytes []byte) []byte
```
HashSha1Bytes computes a Sha1 hash of byte data

--
#### func  HashSha1Json

```go
func HashSha1Json(jsonData interface{}) []byte
```
HashSha1Json computes a Sha1 hash of json data

--
#### func  HeartBeatEtcd

```go
func HeartBeatEtcd(key string, value string, ttl int)
```
HeartbeatEtcd sets key/value pairs to etcd disrtibuted store at an interval Used
to renew a ttl set

--
#### func  JsonErrorResponse

```go
func JsonErrorResponse(res http.ResponseWriter, err error, status int)
```
JsonErrorResponse writes your error to response with a specific http status

--
#### func  RandomStr

```go
func RandomStr(n int) string
```
RandomStr generates a string of random letters and numbers with crypto/rand

--
#### func  ReadJson

```go
func ReadJson(req *http.Request, data interface{}) error
```
ReadJson read json from request and marshal the data

--
#### func  SetEtcdKey

```go
func SetEtcdKey(key string, value string, ttl int) error
```
SetEtcdKey sets key/value pairs to etcd disrtibuted store

--
#### func  SetEtcdURL

```go
func SetEtcdURL(etcdURL ...string) string
```
SetEtcdURL sets the url address and port of the etcd service from environment
variables

--
#### func  VerifyHMAC

```go
func VerifyHMAC(message string, macOfMessage []byte, key []byte) bool
```
VerifyHMAC verifies a HMAC message

--
#### func  VerifyToken

```go
func VerifyToken(key []byte, authToken string, seeds ...string) (bool, error)
```
VerifyToken returns true if token has a valid HMAC.

--
#### func  WriteJson

```go
func WriteJson(res http.ResponseWriter, dataOut interface{}, pretty bool) error
```
WriteJson writes json to response with http.StatusOK

--
#### func  WriteJsonStatus

```go
func WriteJsonStatus(res http.ResponseWriter, status int, dataOut interface{}, pretty bool) error
```
WriteJsonStatus writes json to response and set specfic http status code

