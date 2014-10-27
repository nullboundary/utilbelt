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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	res.Header().Set("Access-Control-Allow-Credentials", "true")

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

	res.WriteHeader(status)

	type errorMap struct {
		ErrorStatus int    `json:"code"`
		Error       string `json:"error"`
	}

	errorStruct := &errorMap{}

	errorStruct.Error = err.Error()
	errorStruct.ErrorStatus = status

	if err := WriteJson(res, errorStruct, true); err != nil {
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
