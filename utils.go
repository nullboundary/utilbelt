package utils

import (
	"bitbucket.org/cicadaDev/storer"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/zenazn/goji/web"
	"net/http"
)

var clientEtcd = etcd.NewClient([]string{"http://10.1.42.1:4001"}) //TODO: find a better way to set this!

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func setEtcdKey(key string, value string) error {
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
func getEtcdKey(key string) (string, error) {

	// GET the value that is stored for the key
	resp, err := clientEtcd.Get(key, false, false)
	if err != nil {
		return "", fmt.Errorf("get etcd key error: %v", err)
	}

	return resp.Node.Value, nil

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
			var err error
			rt.Url, err = getEtcdKey("db/url") //os.Getenv("PASS_APP_DB_URL")
			Check(err)
			rt.Port, err = getEtcdKey("db/port") //os.Getenv("PASS_APP_DB_PORT")
			Check(err)
			rt.DbName, err = getEtcdKey("db/name") //os.Getenv("PASS_APP_DB_NAME")
			Check(err)

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
