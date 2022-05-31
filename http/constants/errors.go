package constants

import "errors"

var ErrHttpUrl 				= errors.New("invalid http url")
var ErrHttpMethod 			= errors.New("invalid http method")
var ErrHttpRequestParams 	= errors.New("invalid http request parameters")