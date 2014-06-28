package utils

import (
	"bitbucket.org/cicadaDev/storer"
	"fmt"
	"github.com/zenazn/goji/web"
	"net/http"
	"os"
)

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
