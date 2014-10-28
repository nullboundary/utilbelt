package utils

import (
	"bitbucket.org/cicadaDev/storer"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/zenazn/goji/web"
	"net/http"
	"os"
)

var clientEtcd = etcd.NewClient([]string{os.Getenv("ETCD_IP")})

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
// function call to set filesystem path
//
//
//////////////////////////////////////////////////////////////////////////
func SafeFileSystem(path http.Dir) filterFS {
	return filterFS{path}
}

//////////////////////////////////////////////////////////////////////////
//
// interface struct to limit access to os file system
//
//
//////////////////////////////////////////////////////////////////////////
type filterFS struct {
	http.FileSystem
}

//////////////////////////////////////////////////////////////////////////
//
//  Open func is FileSystem interface implementation
//
//
//////////////////////////////////////////////////////////////////////////
func (fs filterFS) Open(name string) (http.File, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	return noReadDirFile{f}, nil
}

//////////////////////////////////////////////////////////////////////////
//
// interface struct to implement a special http.File
//
//
//////////////////////////////////////////////////////////////////////////
type noReadDirFile struct {
	http.File
}

//////////////////////////////////////////////////////////////////////////
//
// Overwrite Readdir for noReadDirFile interface implmentation
//
//
//////////////////////////////////////////////////////////////////////////
func (f noReadDirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
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
