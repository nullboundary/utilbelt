package utilbelt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"hash/fnv"
	"os"
	"path"
	"strings"
	"time"
	"unicode"
)

var clientEtcdURL = []string{"http://172.17.0.1:4001"} //Default
var clientEtcd = etcd.NewClient(clientEtcdURL)

//SetEtcdURL sets the url address and port of the etcd service from environment variables
func SetEtcdURL(etcdURL ...string) string {

	//first try the list of etcdURLs
	if len(etcdURL) > 0 {
		clientEtcd = etcd.NewClient(etcdURL)
		return etcdURL[0]
	}

	//Then try an env variable
	addr := os.Getenv("ETCD")
	if addr != "" {
		clientEtcdURL = []string{addr}
		clientEtcd = etcd.NewClient(clientEtcdURL)
		return addr
	}

	//otherwise fall back on default
	return clientEtcdURL[0]

}

//SetEtcdKey sets key/value pairs to etcd disrtibuted store
func SetEtcdKey(key string, value string, ttl int) error {

	// SET the value "bar" to the key "foo" with zero TTL
	_, err := clientEtcd.Set(key, value, uint64(ttl))
	if err != nil {
		return fmt.Errorf("set etcd key error: %v", err)
	}

	return nil

}

//GetEtcdKey retrives key/value pairs from etcd disrtibuted store
func GetEtcdKey(key string) (string, error) {

	// GET the value that is stored for the key
	resp, err := clientEtcd.Get(key, false, false)
	if err != nil {
		return "", fmt.Errorf("get etcd key error: %v", err)
	}

	return resp.Node.Value, nil //TODO: Return []byte?

}

//HeartbeatEtcd sets key/value pairs to etcd disrtibuted store at an interval
//Used to renew a ttl set
func HeartBeatEtcd(key string, value string, ttl int) {

	interval := (ttl * 1000) - 500
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	go func() {
		for _ = range ticker.C {
			SetEtcdKey(key, value, ttl)
		}
	}()

}

//GenerateFnvHashID generates a hash fnv1a hash. Fast, unique, but insecure! use only for ids and such.
//https://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
func GenerateFnvHashID(hashSeeds ...string) uint32 {

	inputString := strings.Join(hashSeeds, "")

	var randomness int32
	binary.Read(rand.Reader, binary.LittleEndian, &randomness) //add a little randomness
	inputString = fmt.Sprintf("%s%x", inputString, randomness)

	h := fnv.New32a()
	h.Write([]byte(inputString))
	return h.Sum32()

}

//EncodetoDataUri reads a file and generates a data uri
func EncodetoDataUri(fileName string, mimeType string, allowTypes ...string) (string, error) {

	file, _ := os.Open(fileName) //TODO: change to read in through form or json
	defer file.Close()

	fileInfo, _ := file.Stat() // FileInfo interface
	size := fileInfo.Size()    // file size

	data := make([]byte, size)

	contentType := path.Ext(fileName)

	typeFound := false
	for _, fileType := range allowTypes { //match the type with the allowed types
		if contentType == fileType {
			typeFound = true
			break
		}
	}

	if !typeFound {
		err := fmt.Errorf("[Error] file type: %s not allowed", contentType)
		return "", err
	}

	file.Read(data)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}

//DecodeUriToBytes decodes a data uri into bytes
func DecodeUriToBytes(str string, fileType string) (string, []byte) {

	dataStr := strings.SplitN(str, ",", 2) //seperate data:image/png;base64, from the DataURI

	fieldTest := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fields := strings.FieldsFunc(dataStr[0], fieldTest) //Fields are: ["data" "image" "png" "base64"]
	dataExt := fields[2]                                //only need the file extension

	if dataExt != fileType {
		err := fmt.Errorf("[Error] file type: %s not allowed", dataExt)
		if err != nil {
			fmt.Printf("file type error %s", err) //TODO. Return error
		}
	}

	data, err := base64.StdEncoding.DecodeString(dataStr[1]) // [] byte
	if err != nil {
		fmt.Printf("base64 encode error %s", err)
	}

	return dataExt, data

}
