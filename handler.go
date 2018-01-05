package main

import (
	"net/http"

	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"io/ioutil"
)

// Handle is the exported handler called by AWS Lambda.
var Handle apigatewayproxy.Handler

func init() {
	ln := net.Listen()

	// Amazon API Gateway binary media types are supported out of the box.
	// If you don't send or receive binary data, you can safely set it to nil.
	Handle = apigatewayproxy.New(ln, []string{"image/png"}).Handle

	// Any Go framework complying with the Go http.Handler interface can be used.
	// This includes, but is not limited to, Vanilla Go, Gin, Echo, Gorrila, Goa, etc.
	go http.Serve(ln, http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World! from api"))
}

func HandlePlain(evt interface{}, ctx *runtime.Context) (string, error) {
	putEndpoint :=  "http://wm68vs7yg8.execute-api.eu-west-1.amazonaws.com/dev/test/122" //os.Getenv("apiBaseUrl")
	println(putEndpoint)
	resp, err := http.Get(putEndpoint)
	if err != nil {
		return "",err
	}
	bodyContent,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}

	return string(bodyContent), nil
}