package utilbelt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//ReadJson read json from request and marshal the data
func ReadJson(req *http.Request, data interface{}) error {

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return fmt.Errorf("jsonRead Error: %v", err)
	}

	return nil
}

//WriteJsonStatus writes json to response and set specific http status code
func WriteJsonStatus(res http.ResponseWriter, status int, dataOut interface{}, pretty bool) error {

	res.Header().Add("Content-Type", "application/json")
	res.Header().Add("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	res.Header().Add("Access-Control-Allow-Credentials", "true")
	res.WriteHeader(status)

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

//WriteJson writes json to response with http.StatusOK
func WriteJson(res http.ResponseWriter, dataOut interface{}, pretty bool) error {

	//default to status ok, unless specified
	return WriteJsonStatus(res, http.StatusOK, dataOut, pretty)
}

//JsonErrorResponse writes your error to response with a specific http status
func JsonErrorResponse(res http.ResponseWriter, err error, status int) {

	errorReport := map[string]string{"code": fmt.Sprintf("%d", status), "error": err.Error()}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(errorReport); err != nil {
		log.Printf("jsonWrite Error: %v", err)
	}
}

//DebugPrintJson prints json data for debugging
func DebugPrintJson(Data interface{}) {

	printJSon := json.NewEncoder(os.Stdout)
	printJSon.Encode(Data)
}
