package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/xordataexchange/crypt/config"
	"hash/fnv"
	"os"
	"strings"
	"time"
)

var clientEtcdURL = []string{"http://172.17.42.1:4001"} //Default
var clientEtcd = etcd.NewClient(clientEtcdURL)

//////////////////////////////////////////////////////////////////////////
//
//	set the url address and port of the etcd service from environment variables
//
//
//////////////////////////////////////////////////////////////////////////
func SetEtcdURL() string {
	addr := os.Getenv("ETCD") //"http://10.1.42.1:4001"
	if addr != "" {
		clientEtcdURL = []string{addr}
		clientEtcd = etcd.NewClient(clientEtcdURL)
		return addr
	}
	return clientEtcdURL[0]

}

//////////////////////////////////////////////////////////////////////////
//
//	SetEtcdKey sets key/value pairs to etcd disrtibuted store
//
//
//////////////////////////////////////////////////////////////////////////
func SetEtcdKey(key string, value string, ttl int) error {
	// SET the value "bar" to the key "foo" with zero TTL
	// returns a: *store.Response
	_, err := clientEtcd.Set(key, value, uint64(ttl))
	if err != nil {
		return fmt.Errorf("set etcd key error: %v", err)
	}

	return nil

}

//////////////////////////////////////////////////////////////////////////
//
//	GetEtcdKey retrives key/value pairs from etcd disrtibuted store
//	//TODO: Return []byte?
//
//////////////////////////////////////////////////////////////////////////
func GetEtcdKey(key string) (string, error) {

	// GET the value that is stored for the key
	resp, err := clientEtcd.Get(key, false, false)
	if err != nil {
		return "", fmt.Errorf("get etcd key error: %v", err)
	}

	return resp.Node.Value, nil

}

//////////////////////////////////////////////////////////////////////////
//
//	HeartbeatEtcd sets key/value pairs to etcd disrtibuted store at an interval
//  Used to renew a ttl set
//
//////////////////////////////////////////////////////////////////////////
func HeartBeatEtcd(key string, value string, ttl int) {

	interval := (ttl * 1000) - 500
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	go func() {
		for _ = range ticker.C {
			SetEtcdKey(key, value, ttl)
		}
	}()

}

//////////////////////////////////////////////////////////////////////////
//
//	Used for loading encrypted etcd key/value pairs.
//	crypt set -keyring .pubring.gpg -endpoint http://172.17.42.1:4001 /catagory/variable filename
//
//////////////////////////////////////////////////////////////////////////
func GetCryptKey(keyringPath string, key string) ([]byte, error) {

	//get key ring
	kr, err := os.Open(keyringPath)
	if err != nil {
		return nil, fmt.Errorf("open keyring error: %v", err)
	}
	defer kr.Close()

	//setup etcd manager
	cm, err := config.NewEtcdConfigManager(clientEtcdURL, kr)
	if err != nil {
		return nil, fmt.Errorf("setup etcd manager error: %v", err)
	}

	value, err := cm.Get(key)
	if err != nil {
		return nil, fmt.Errorf("get crypt %v error: %v", clientEtcdURL, err)
	}

	return value, nil
}

//////////////////////////////////////////////////////////////////////////
//
//	generate a hash fnv1a hash. Fast, unique, but insecure! use only for ids and such.
//  https://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
//
//////////////////////////////////////////////////////////////////////////
func GenerateFnvHashID(hashSeeds ...string) uint32 {

	inputString := strings.Join(hashSeeds, "")

	var randomness int32
	binary.Read(rand.Reader, binary.LittleEndian, &randomness) //add a little randomness
	inputString = fmt.Sprintf("%s%x", inputString, randomness)

	h := fnv.New32a()
	h.Write([]byte(inputString))
	return h.Sum32()

}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
