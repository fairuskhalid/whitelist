package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/Sirupsen/logrus"
)

const port = ":8080"
var flPort = flag.String("port", port, "specifies the server port number")

func main() {
	flag.Parse()


	http.HandleFunc("/getlist", func(w http.ResponseWriter, r *http.Request) {
        	dat, err := ioutil.ReadFile("./whitelist.dat")

	        if err != nil {
        	        logrus.Errorf(err.Error())
                	return
	        }
	
		dataEnc := b64.URLEncoding.EncodeToString([]byte(dat))
		fmt.Fprintf(w, "%s", dataEnc)
	})


	logrus.Infof("Server running @ %s", *flPort)
	err := http.ListenAndServe(*flPort, nil)
	if err != nil {
		logrus.Errorf("Unable to start server: %s", err.Error())
	}

}
