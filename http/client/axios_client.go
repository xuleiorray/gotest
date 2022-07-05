package client

import (
	"net/http"
	"perftest/http/config"
	"perftest/http/logger"
	"perftest/http/model"
	"time"

	"github.com/vicanso/go-axios"
)

type AxiosClient struct {
    axiosIns *axios.Instance
}

var log = logger.LOGGER
var HttpClient *AxiosClient

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
    HttpClient = New(axiosIns)
}

func (client *AxiosClient) Dispatch(request *model.HttpRequest) *model.HttpResponse {

    slowTime := config.INSTANCE.GetInt(config.HTTP_REQUEST_SLOW_THRESHOLD)
    timeStart := time.Now()
    resp, err := client.axiosIns.Request(buildRequestConfig(request))
    duration := time.Since(timeStart).Milliseconds()
    
    if duration >= int64(slowTime) {
        log.Infof("Slow request execute takes %dms", duration)
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