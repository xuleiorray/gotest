package task

import (
	"net/http"
	"perftest/http/model"
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
        Timeout: time.Minute,
        OnDone: func(config *axios.Config, resp *axios.Response, err error) {
            if err != nil {
                log.Info(err)
            }
        },
    })
    axiosClient = New(axiosIns)
}

func (client *AxiosClient) Dispatch(request *model.HttpRequest) *model.HttpResponse {

    resp, err := client.axiosIns.Request(buildRequestConfig(request))
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

func buildRequestConfig(request *model.HttpRequest) *axios.Config {
    return &axios.Config{
        URL: request.Url,
        Method: request.Method,
        Body: request.Body,
        Query: toMapArray(request.Params),
        Headers: toMapArray(request.Headers),
    }
}

func toMapArray(maps map[string]string) map[string][]string {
    mapParams := make(map[string][]string)
    for k, v := range maps {
        mapParams[k] = []string{v}
    }
    return mapParams
}