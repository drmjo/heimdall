package proxy

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/drmjo/heimdall/spec"
	api "github.com/drmjo/heimdall/spec"
	"github.com/gorilla/mux"

	uuid "github.com/satori/go.uuid"
)

var routes map[string]string

func setRequestIdOnResponse(res *http.Response) {
	req := res.Request
	res.Header.Set(HEADER_REQUEST_ID, req.Header.Get(HEADER_REQUEST_ID))
	res.Header.Set(HEADER_REQUEST_HOST, req.URL.Hostname())
	res.Header.Set(HEADER_REQUEST_PATH, req.URL.EscapedPath())
}

func logRequest(r *http.Request) {
	log.Print("========== Vars ==========")
	vars := mux.Vars(r)
	log.Printf("vars: %s", vars)
	log.Print("========== End Vars ==========")
	log.Print("========== Request ==========")
	for key, value := range r.Header {
		log.Printf("%s: %s", key, value)
	}
	log.Printf("Uri: %s", r.URL.RequestURI())
	log.Print("========== End Request ==========")
}

func logResponse(headers http.Header, body []byte) {

	// defer res.Body.Close()

	log.Print("========== Response ==========")
	log.Print("+++ HEADERS +++")
	for key, value := range headers {
		log.Printf("%s: %s", key, value)
	}
	log.Print("+++ END HEADERS +++")
	log.Print("+++ BODY +++")
	log.Printf("%s", body)
	log.Print("+++ END BODY +++")
	log.Print("========== End Response ==========")

}

func decodeBody(res *http.Response) []byte {
	var bodyBuffer bytes.Buffer
	var body []byte

	defer bodyBuffer.Reset()

	io.Copy(&bodyBuffer, res.Body)
	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewReader(bodyBuffer.Bytes()))

	// check content encoding
	switch res.Header.Get(HEADER_CONTENT_ENCODING) {
	case "":
		// nothing to decode just return the raw bytes
		body = bodyBuffer.Bytes()
	case "gzip":
		log.Printf("Gzipped response...")
		reader, _ := gzip.NewReader(bytes.NewReader(bodyBuffer.Bytes()))
		body, _ = ioutil.ReadAll(reader)
	default:
		log.Printf("Unsupported encoding refusing to pass through")
		log.Fatalf(
			"%s: %s",
			HEADER_CONTENT_ENCODING,
			res.Header.Get(HEADER_CONTENT_ENCODING))
	}

	return body
}

func validateGetResponse(get *api.Operation, res *http.Response) error {
	return nil
}

func newResponseModifier(item *api.PathItem) func(res *http.Response) error {
	return func(res *http.Response) error {

		// modify headers
		setRequestIdOnResponse(res)

		// modify body
		body := decodeBody(res)
		logResponse(res.Header, body)

		var o *api.Operation
		if o = item.GetOperation(res.Request.Method); o == nil {
			return errors.New("Undefined request Method")
		}

		responses := o.Responses
		// get the proper response definition
		if response, ok := responses[strconv.Itoa(res.StatusCode)]; ok {
			// validate to see if the returned response has the proper content type
			if _, ok := response.Content[res.Header.Get(HEADER_CONTENT_TYPE)]; ok {
				body = nil
				return nil
			}
		}

		body = nil
		return errors.New("Invalid response from upstream")
	}
}

// this does nothing but to make sure the url defined in the schema
// is valid and proxy passes the request to it
func newDirector(servers []*api.Server) func(req *http.Request) {
	return func(req *http.Request) {

		server := servers[0]
		url, _ := url.Parse(server.URL)

		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host

		log.Printf("forwarding to %v", req.URL)

		// also set the host header
		req.Host = url.Host
	}
}

func newRoute(
	h *mux.Router,
	path string,
	item *api.PathItem,
	disableEncoding bool) *mux.Route {
	route := h.NewRoute().Path(path).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		id, _ := uuid.NewV4()
		req.Header.Set(HEADER_REQUEST_ID, id.String())

		logRequest(req)

		if disableEncoding {
			req.Header.Del("Accept-Encoding")
		}

		director := newDirector(item.Servers)
		responseModifier := newResponseModifier(item)

		proxy := &httputil.ReverseProxy{
			Transport:      Transport,
			Director:       director,
			ModifyResponse: responseModifier}

		proxy.ServeHTTP(w, req)
	})

	return route
}

func setMethods(route *mux.Route, item *api.PathItem) {
	var methods []string
	if item.Get != nil {
		methods = append(methods, "GET")
	}
	if item.Post != nil {
		methods = append(methods, "POST")
	}
	if item.Put != nil {
		methods = append(methods, "PUT")
	}

	route.Methods(methods...)
}

// generate a new mux.Router
func NewRouter(api *spec.Api, disableEncoding bool) *mux.Router {
	router := mux.NewRouter()
	// loop through the paths in the shema to gnerate the mux
	for path, item := range api.Paths {
		route := newRoute(router, path, item, disableEncoding)
		setMethods(route, item)
	}

	return router
}
