package wlPlugin 

import (
	"regexp"
        "strings"
        "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/go-plugins-helpers/authorization"
	"golang.org/x/net/context"
)

var whiteList []string

func UpdateList(input []string) {
	whiteList = input
	logrus.Debug("Updated whitelist %s", whiteList)
}

func chkList(id string) bool {
        for i := range whiteList {
                if strings.Compare(id, whiteList[i]) == 0 {
                        return true
                }
        }
        return false
}

type wlPlugin struct {
        wlclient *client.Client
}

func NewPlugin(dockerHost string) (*wlPlugin, error) {

	defaultHeaders := map[string]string{"wlPlugin-Agent": "engine-api-cli-1.0"}

        logrus.Infof("host '%s' ", dockerHost)
	c, err := client.NewClient(dockerHost, "v1.22", nil, defaultHeaders)
	if err != nil {
		return nil, err
	}
	list := []string{"sha256:70c557e50ed630deed07cbb0dc4d28aa0f2a485cf7af124cc48f06bce83f784b"}
	UpdateList(list)
	return &wlPlugin{wlclient: c}, nil
}


var re = regexp.MustCompile(`/containers/(.*)/start$`)

func (p *wlPlugin) AuthZReq(req authorization.Request) authorization.Response {
        logrus.Debug("Request URI: %s", req.RequestURI)
        res := re.FindStringSubmatch(req.RequestURI)
        if req.RequestMethod == "POST" && len(res) > 0 {
                logrus.Debug("Cntnr Id %s", res[1])
                container, err := p.wlclient.ContainerInspect(context.Background(), res[1])
                if err != nil {
                        logrus.Error("Error ContainerInspect %s", err.Error())
                        return authorization.Response{Allow: true, Err: err.Error()}
                }
                logrus.Debug("Image id: %s", container.Image)
                if !chkList(container.Image) {
                        options := types.ContainerRemoveOptions{
                                ContainerID:   res[1],
                                RemoveVolumes: false,
                                RemoveLinks:   false,
                                Force:         false,
                        }

                        if err := p.wlclient.ContainerRemove(context.Background(), options); err != nil {
                                logrus.Error("Error remove Container %s", err.Error())
                        }

                        return authorization.Response{Allow: false, Msg: "Unauthorized Image"}
                }
        }
        return authorization.Response{Allow: true}
}

func (p *wlPlugin) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}
