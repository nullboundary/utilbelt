package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
func WriteJson(res http.ResponseWriter, dataOut interface{}, pretty bool) error {

	res.Header().Add("Content-Type", "application/json")
	res.Header().Add("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	res.Header().Add("Access-Control-Allow-Credentials", "true")

	//pretty printing.
	if pretty {

		b, err := json.MarshalIndent(dataOut, "", "  ")
		if err != nil {
			return fmt.Errorf("jsonWrite Error: %v", err)
		}
		res.Write(b)

	} else {
		if err := json.NewEncoder(res).Encode(dataOut); err != nil { //encode the result struct to json and output on response writer
			return fmt.Errorf("jsonWrite Error: %v", err)
		}
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

	errorReport := map[string]string{"code": fmt.Sprintf("%d", status), "error": err.Error()}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(errorReport); err != nil {
		log.Printf("jsonWrite Error: %v", err)
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
