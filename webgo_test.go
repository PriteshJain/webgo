package webgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type response struct {
	Data   map[string]string `json:"data"`
	Status int               `json:"status"`
}

// appConfig is a sample struct to hold app configurations
type appConfig struct {
	Name string
}

const p1 = "world"
const p2 = "spiderman"

const baseapi = "http://127.0.0.1:9696"
const baseapiHTTPS = "http://127.0.0.1:9696"

var BenchAPIs = map[string]string{
	"GETNOPARAM":    strings.Join([]string{baseapi, "nparams"}, "/"),
	"GETWITHPARAM":  strings.Join([]string{baseapi, "wparams", p1, "goblin", p2}, "/"),
	"POSTWITHPARAM": strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
}
var GETAppContextAPI = []string{
	strings.Join([]string{baseapi, "appcontext"}, "/"),
	strings.Join([]string{baseapiHTTPS, "appcontext"}, "/"),
}

var GETAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var POSTAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var PUTAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var DELETEAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var PATCHAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var OPTIONSAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

func setup() (*Router, *httptest.ResponseRecorder) {
	// Initializing router with all the required routes
	router := NewRouter(&Config{
		Host:               "127.0.0.1",
		Port:               "9696",
		HTTPSPort:          "8443",
		CertFile:           "tests/ssl/server.crt",
		KeyFile:            "tests/ssl/server.key",
		ReadTimeout:        15,
		WriteTimeout:       60,
		InsecureSkipVerify: true,
	}, getRoutes())

	router.AppContext = map[string]interface{}{
		"config": &appConfig{
			Name: "WebGo",
		},
	}
	return router, httptest.NewRecorder()
}

func TestInvalidHTTPMethod(t *testing.T) {
	router, respRec := setup()

	for _, url := range GETAPI {
		req, err := http.NewRequest("ABC", url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)
			continue
		}

		router.ServeHTTP(respRec, req)
		if respRec.Code != http.StatusNotImplemented {
			t.Fatalf(`Expected response HTTP status code %d, received %d`, http.StatusNotImplemented, respRec.Code)
		}
	}
}

func TestGet(t *testing.T) {
	router, respRec := setup()

	for _, url := range GETAPI {
		req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)
			continue
		}

		router.ServeHTTP(respRec, req)

		resp := response{}

		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)
			continue
		}

		if resp.Data["method"] != http.MethodGet {
			t.Fatal("URL:", url, "response method:", resp.Data["method"], " required method:", http.MethodGet)
		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])
		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])
		}
	}
}

func TestMiddleware(t *testing.T) {
	router, respRec := setup()
	router.Use(mware)
	url := baseapi + "/"
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)

	}
	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Fatal(err, respRec.Code, url)

	}

	v := respRec.Header().Get("k1")
	if respRec.Header().Get("k1") != "v1" {
		t.Fatal("Expected response header value `v1` for key `k1`, received", v)

	}
}
func TestGetPostResponse(t *testing.T) {
	router, respRec := setup()
	url := baseapi + "/"
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)

	}
	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Fatal(err, respRec.Code, url)

	}
}
func TestGet404(t *testing.T) {
	router, respRec := setup()
	url := baseapi + "/random"
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)

	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 404 {
		t.Fatal(err)

	}
}
func TestHead(t *testing.T) {
	router, respRec := setup()

	for _, url := range GETAPI {
		req, err := http.NewRequest(http.MethodHead, url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		router.ServeHTTP(respRec, req)

		resp := response{}

		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodHead {
			t.Fatal("URL:", url, "response method:", resp.Data["method"], " required method:", http.MethodGet)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}
	}
}

func TestPost(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)

	for _, url := range POSTAPI {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}
		router.ServeHTTP(respRec, req)
		resp := response{}
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodPost {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodPost)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestPut(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range PUTAPI {

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}
		router.ServeHTTP(respRec, req)
		resp := response{}
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodPut {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodPut)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestPatch(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range PATCHAPI {
		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodPatch {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodPatch)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestDelete(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range DELETEAPI {
		req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodDelete {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodDelete)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestOptions(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)

	for _, url := range OPTIONSAPI {
		req, err := http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodOptions {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodOptions)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestAppContext(t *testing.T) {
	router, respRec := setup()
	for _, url := range GETAppContextAPI {
		req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["Name"] != "WebGo" {
			t.Fatal("Invalid App context config received")

		}
	}

}

func TestStart(t *testing.T) {
	router, _ := setup()
	go router.Start()
	time.Sleep(time.Second * 5)
	err := router.Shutdown()
	if err != nil {
		t.Fatal(err)

	}
}
func TestStartHTTPS(t *testing.T) {
	router, _ := setup()
	go router.StartHTTPS()
	time.Sleep(time.Second * 5)
	err := router.ShutdownHTTPS()
	if err != nil {
		t.Fatal(err)

	}
}

func dummy(w http.ResponseWriter, r *http.Request) {

	output := ""

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		R400(w, err)
		return
	}
	defer r.Body.Close()

	output = string(b)

	wctx := Context(r)

	R200(
		w,
		map[string]string{
			"p1":      wctx.Params["p1"],
			"p2":      wctx.Params["p2"],
			"payload": output,
			"pattern": r.URL.Path,
			"method":  r.Method,
		},
	)
}

func getAppConfig(w http.ResponseWriter, r *http.Request) {
	wctx := Context(r)
	aC, ok := wctx.AppContext["config"].(*appConfig)
	if !ok {
		R400(w, "No app config found")
		return
	}
	R200(w, aC)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	R200(w, "Hello world")
}

func postResp(w http.ResponseWriter, r *http.Request) {
	log.Println("This is a post response handler")
}

func getRoutes() []*Route {
	return []*Route{
		{
			Name:                    "root",         // A label for the API/URI
			Method:                  http.MethodGet, // request type
			Pattern:                 "/",
			FallThroughPostResponse: true, // Pattern for the route
			TrailingSlash:           true,
			Handlers:                []http.HandlerFunc{dummy, postResp}, // route handler
		},
		{
			Name:     "appcontext",                     // A label for the API/URI
			Method:   http.MethodGet,                   // request type
			Pattern:  "/appcontext",                    // Pattern for the route
			Handlers: []http.HandlerFunc{getAppConfig}, // route handler
		},
		{
			Name:     "hw-noparams",                  // A label for the API/URI
			Method:   http.MethodGet,                 // request type
			Pattern:  "/nparams",                     // Pattern for the route
			Handlers: []http.HandlerFunc{helloWorld}, // route handler
		},
		{
			Name:          "hw-withparams", // A label for the API/URI
			Method:        http.MethodGet,
			TrailingSlash: true,                           // request type
			Pattern:       "/wparams/:p1/goblin/:p2",      // Pattern for the route
			Handlers:      []http.HandlerFunc{helloWorld}, // route handler
		},
		{
			Name:     "params-get",              // A label for the API/URI
			Method:   http.MethodGet,            // request type
			Pattern:  "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},
		{
			Name:     "params-head",             // A label for the API/URI
			Method:   http.MethodHead,           // request type
			Pattern:  "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},

		{
			Name:     "params-post-sameuri",     // A label for the API/URI
			Method:   http.MethodPost,           // request type
			Pattern:  "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},
		{
			Name:     "params-put-sameuri",      // A label for the API/URI
			Method:   http.MethodPut,            // request type
			Pattern:  "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},
		{
			Name:     "params-patch-sameuri",    // A label for the API/URI
			Method:   http.MethodPatch,          // request type
			Pattern:  "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},
		{
			Name:     "params-delete-sameuri",   // A label for the API/URI
			Method:   http.MethodDelete,         // request type
			Pattern:  "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},
		{
			Name:    "params-options-sameuri", // A label for the API/URI
			Method:  http.MethodOptions,       // request type
			Pattern: "/hello/:p1/goblin/:p2",  // Pattern for the route
			// Handler: []http.HandlerFunc{dummy}, // route handler
			Handlers: []http.HandlerFunc{dummy}, // route handler
		},
	}
}

func mware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	rw.Header().Add("k1", "v1")
	next(rw, req)
}
