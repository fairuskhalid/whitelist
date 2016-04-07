package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/fairuskhalid/whitelist/util"
	"github.com/robfig/cron"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	pluginSocket  = "/run/docker/plugins/whitelist-plugin.sock"
	wlHost =  "http://localhost:8080/getlist"
	interval = "@every 1m"
	debug = false
)

var (
	flDockerHost = flag.String("dhost", dockerHost, "Specifies the host where to contact the docker daemon")
	flWhiteListHost = flag.String("wlhost", wlHost, "Specifies the host where to get whitelist")
	flInterval = flag.String("i", interval, "Specifies the interval to query the whitelist")
	flDebug = flag.Bool("d", debug, "Show debug or not ")
)

type wlJob struct{}

// cron job to query the whitelist from specified server
func (d wlJob) Run() {
	resp, err := http.Get(*flWhiteListHost)
	if err != nil {
		logrus.Error("Get whitelist: %s", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                logrus.Error("Get whitelist: %s", err.Error())
                return
        }

	dataDec, err := b64.URLEncoding.DecodeString(string(body))
	if err != nil {
		logrus.Error("Encoding error %s", err.Error())
		return
	}

	r := csv.NewReader(bytes.NewReader(dataDec))
	r.TrimLeadingSpace = true
	list, err := r.Read()
	if err != nil {
		logrus.Info("Server return empty list")
		return
	}
	wlPlugin.UpdateList(list)

}

func main() {
	flag.Parse()

	if *flDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
                logrus.SetLevel(logrus.InfoLevel)
        }

	plugin, err := wlPlugin.NewPlugin(*flDockerHost)
	if err != nil {
		logrus.Fatal(err)
	}

        logrus.Infof("wlhost %s", *flWhiteListHost)

	// add a cron job
        var job wlJob
        c:= cron.New()
        c.AddJob(*flInterval, job)
        c.Start()

	h := authorization.NewHandler(plugin)

	if err := h.ServeUnix("root", pluginSocket); err != nil {
		logrus.Fatal(err)
	}

}
