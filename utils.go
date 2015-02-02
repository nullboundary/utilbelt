package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"hash/fnv"
	"strings"
)

var clientEtcd = etcd.NewClient([]string{"http://10.1.42.1:4001"}) //TODO: find a better way to set this!

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func SetEtcdKey(key string, value string) error {
	// SET the value "bar" to the key "foo" with zero TTL
	// returns a: *store.Response
	_, err := clientEtcd.Set(key, value, 0)
	if err != nil {
		return fmt.Errorf("set etcd key error: %v", err)
	}

	return nil

}

//////////////////////////////////////////////////////////////////////////
//
//
//
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
