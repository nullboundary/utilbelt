package utils

import (
	"bitbucket.org/cicadaDev/storer"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"os"
	"strings"
)

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func ReadJson(req *http.Request, data interface{}) error {

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return fmt.Errorf("jsonRead Error: %v", err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func WriteJson(res http.ResponseWriter, dataOut interface{}) error {

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	res.Header().Set("Access-Control-Allow-Credentials", "true")

	//TODO: Implement pretty printing. //err = json.MarshalIndent(result, "", "  ")

	if err := json.NewEncoder(res).Encode(dataOut); err != nil { //encode the result struct to json and output on response writer
		return fmt.Errorf("jsonWrite Error: %v", err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func JsonErrorResponse(res http.ResponseWriter, err error, status int) {

	res.WriteHeader(status)

	type errorMap struct {
		ErrorStatus int    `json:"code"`
		Error       string `json:"error"`
	}

	errorStruct := &errorMap{}

	errorStruct.Error = err.Error()
	errorStruct.ErrorStatus = status

	if err := WriteJson(res, errorStruct); err != nil {
		log.Printf("json write Error: %s", err.Error())
	}
}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func DebugPrintJson(Data interface{}) {

	printJSon := json.NewEncoder(os.Stdout)
	printJSon.Encode(Data)
}

//////////////////////////////////////////////////////////////////////////
//
//	getDbType
//
//
//////////////////////////////////////////////////////////////////////////
func GetDbType(c web.C) (storer.Storer, error) {

	if v, ok := c.Env["db"]; ok {

		if db, ok := v.(storer.Storer); ok {

			return db, nil //all good

		} else {
			err := fmt.Errorf("value could not convert to type Storer")
			return nil, err
		}

	} else {
		err := fmt.Errorf("value for key db, not found")
		return nil, err
	}

}

//////////////////////////////////////////////////////////////////////////
//
//	addDb Middleware
//
//
//////////////////////////////////////////////////////////////////////////
func AddDb(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if c.Env == nil {
			c.Env = make(map[string]interface{})
		}

		if _, ok := c.Env["db"]; !ok { //test is the db is already added

			rt := storer.NewReThink()
			rt.Url = os.Getenv("PASS_APP_DB_URL")
			rt.Port = os.Getenv("PASS_APP_DB_PORT")
			rt.DbName = os.Getenv("PASS_APP_DB_NAME")

			s := storer.Storer(rt) //abstract cb to a Storer
			s.Conn()

			c.Env["db"] = s //add db
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
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

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func GenerateToken(key string, seeds ...string) string {

	tokenSeed := strings.Join(seeds, "|")
	hmac := CalcHMAC(tokenSeed, key)
	return base64.URLEncoding.EncodeToString(hmac)

}

//////////////////////////////////////////////////////////////////////////
//
// verifyToken returns true if messageMAC is a valid HMAC tag for message.
//
//
//////////////////////////////////////////////////////////////////////////
func VerifyToken(key string, authToken string, seeds ...string) (bool, error) {

	decodedMac, err := base64.URLEncoding.DecodeString(authToken)
	if err != nil {
		return false, fmt.Errorf("base64 Decode Error: %s", err)
	}
	tokenSeed := strings.Join(seeds, "|")
	return VerifyHMAC(tokenSeed, decodedMac, key), nil

}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func CalcHMAC(message string, key string) []byte {

	mac := hmac.New(sha256.New, []byte(key))
	n, err := mac.Write([]byte(message))
	if n != len(message) || err != nil {
		panic(err)
	}
	return mac.Sum(nil)
}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func VerifyHMAC(message string, macOfMessage []byte, key string) bool {

	mac := hmac.New(sha256.New, []byte(key))
	n, err := mac.Write([]byte(message))
	if n != len(message) || err != nil {
		panic(err)
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(macOfMessage, expectedMAC)
}
