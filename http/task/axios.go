package task

import (
	"net/http"
	"perftest/http/model"
	"strings"
	"time"

	"github.com/vicanso/go-axios"
)

type AxiosClient struct {
    axiosIns *axios.Instance
}

var axiosClient *AxiosClient

func New(axiosIns *axios.Instance) *AxiosClient {
    return &AxiosClient{
        axiosIns: axiosIns,
    }
}
func init() {
	axiosIns := axios.NewInstance(&axios.InstanceConfig{
        //BaseURL:     "https://www.xxx.com",
        EnableTrace: true,
        Client: &http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyFromEnvironment,
            },
        },
        Timeout: 10 * time.Second,
        OnDone: func(config *axios.Config, resp *axios.Response, err error) {
            if err != nil {
                log.Info(err)
            }
        },
    })
    axiosClient = New(axiosIns)
}

func (client *AxiosClient) Dispatch(request *model.HttpRequest) *model.HttpResponse {
    var resp *axios.Response
    var err error
    if strings.ToUpper(request.Method) == "GET" {
        resp, err = client.axiosIns.Get(request.Url, queryString(request))
    } else if strings.ToUpper(request.Method) == "POST" {
        resp, err = client.axiosIns.Post(request.Url, request.Body, queryString(request))
    } else {
        log.Errorln("Unsupported http request method, just support GET and POST method")
        return nil
    }
    if err != nil {
        log.Infof("Get request execute failed with error: %s", err.Error())
        return nil
    }

    httpResponse := &model.HttpResponse {
        HttpRequest: *request,
        StatusCode: resp.Status,
        Body: string(resp.Data),
        Headers: resp.Headers,
    }
    return httpResponse
}


func queryString(request *model.HttpRequest) map[string][]string {
    mapParams := make(map[string][]string)
    for k, v := range request.Params {
        mapParams[k] = []string{v}
    }
    return mapParams
}