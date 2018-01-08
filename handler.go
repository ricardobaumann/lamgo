package main

import (
	"net/http"

	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/go-redis/redis"
	"io/ioutil"
	"os"
	"strings"
	"io"
	"log"
	"fmt"
	"encoding/json"
)
const redisEndpoint = "localhost:6379"

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

	client := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	switch r.Method {
	case "GET":
		id := strings.TrimPrefix(r.URL.Path, "/test/")
		w.Write([]byte("Hello, World! from api for id "+id))
	case "PUT":
		body, _ := ioutil.ReadAll(r.Body)
		results := make(map[string][]string)
		json.Unmarshal(body, &results)
		for k,v := range results {
			fmt.Printf("key[%s] value[%s]\n", k, v)
		}

		w.Write([]byte(string(body)))
	}

}

func HandlePlain(evt interface{}, ctx *runtime.Context) (string, error) {
	getEndpoint :=  os.Getenv("apiBaseUrl")+"/test/122"
	println(getEndpoint)
	resp, err := http.Get(getEndpoint)
	if err != nil {
		return "",err
	}
	bodyContent,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	putEndpoint := os.Getenv("apiBaseUrl")+"/test"
	putRequest(putEndpoint,strings.NewReader("{}"))
	return string(bodyContent), nil
}

func putRequest(url string, data io.Reader)  {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		// handle error
		log.Fatal(err)
	} else {
		log.Print(req.Body)
	}
	_, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}


}