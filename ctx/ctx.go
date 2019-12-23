package ctx

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/squareroot/mix"
)

//Ctx provides context around Requests and Responses
type c struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	//client         models.ClientCred
	//introspect     client.Inspector
}

func New(response http.ResponseWriter, request *http.Request /*, client models.ClientCred , introspect client.Inspector*/) Responder {
	return &c{
		responseWriter: response,
		request:        request,
		//client:         client,
		//introspect:     introspect,
	}
}

func (c *c) Request() *http.Request {
	return c.request
}

func (c *c) Responder() http.ResponseWriter {
	return c.responseWriter
}

//Method returns the Requests' Method
func (c *c) Method() string {
	return c.request.Method
}

//GetHeader returns a Request Header
func (c *c) GetHeader(key string) (string, error) {
	headers := c.request.Header[key]

	if len(headers) == 0 {
		return "", fmt.Errorf("no header '%s' found", key)
	}

	return headers[0], nil
}

//SetHeader sets a value on the Response Header
func (c *c) SetHeader(key string, val string) {
	c.responseWriter.Header().Set(key, val)
}

//SetStatus set the final Response Status
func (c *c) SetStatus(code int) {
	c.responseWriter.WriteHeader(code)
}

//File returns the Uploaded file.
func (c *c) File(name string) (multipart.File, *multipart.FileHeader, error) {
	err := c.request.ParseMultipartForm(32 << 20)

	if err != nil {
		return nil, nil, err
	}

	return c.request.FormFile(name)
}

//FindFormValue is used to read additional information from File Uploads
func (c *c) FindFormValue(name string) string {
	return c.request.FormValue(name)
}

//FindQueryParam returns the requested querystring parameter
func (c *c) FindQueryParam(name string) string {
	results, ok := c.request.URL.Query()[name]

	if !ok {
		return ""
	}

	return results[0]
}

//FindParam returns the requested path variable
func (c *c) FindParam(name string) string {
	vars := mux.Vars(c.request)

	result, ok := vars[name]

	if !ok {
		return ""
	}

	return result
}

func (c *c) Redirect(status int, url string) {
	http.Redirect(c.responseWriter, c.request, url, status)
}

func (c *c) WriteResponse(data []byte) (int, error) {
	return c.responseWriter.Write(data)
}

func (c *c) WriteStreamResponse(data io.Reader) (int64, error) {
	return io.Copy(c.responseWriter, data)
}

func (c *c) RequestURI() string {
	return c.request.URL.RequestURI()
}

func (c *c) GetCookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

func (c *c) Scheme() string {
	return c.request.URL.Scheme
}

func (c *c) Host() string {
	return c.request.Host
}

//Body returns an error when unable to Decode the JSON request
func (c *c) Body(container interface{}) error {
	decoder := json.NewDecoder(c.request.Body)

	return decoder.Decode(container)
}

func (c *c) GetInstanceID() string {
	return "nothing" //c.client.ID
}

//Serve is usually sent a Mixer. Serve(mixer.JSON(500, nil))
func (c *c) Serve(mxFunc mix.InitFunc, srvFunc ServeFunc) error {
	status, data := srvFunc(c)
	mxr := mxFunc(c.RequestURI(), data)

	for key, head := range mxr.Headers() {
		c.SetHeader(key, head)
	}

	if status != http.StatusOK {
		c.SetStatus(status)
	}

	readr, err := mxr.Reader()

	if err != nil {
		return err
	}

	_, err = io.Copy(c.responseWriter, readr)

	return err
}

//GetKeyedRequest will return the Key and update the Target when Requests are sent for updates.
func (c *c) GetKeyedRequest(target interface{}) (husk.Key, error) {
	result := struct {
		Key  husk.Key
		Body interface{}
	}{
		Body: target,
	}

	err := c.Body(&result)

	if err != nil {
		return husk.CrazyKey(), err
	}

	return result.Key, nil
}

//GetPageData turns /B1 into page 1. size 1
func (c *c) GetPageData() (page, pageSize int) {
	pageData := c.FindParam("pagesize")
	return getPageData(pageData)
}

func getPageData(pageData string) (int, int) {
	defaultPage := 1
	defaultSize := 10

	if len(pageData) < 2 {
		return defaultPage, defaultSize
	}

	pChar := []rune(pageData[:1])

	if len(pChar) != 1 {
		return defaultPage, defaultSize
	}

	page := int(pChar[0]) % 32
	pageSize, err := strconv.Atoi(pageData[1:])

	if err != nil {
		return defaultPage, defaultSize
	}

	return page, pageSize
}

func (c *c) GetMyToken() string {
	cooki, err := c.GetCookie("avosession")

	if err != nil {
		return ""
	}

	return cooki.Value
}

func (c *c) GetMyUser() interface{} /*models.ClaimIdentity*/ {
	//token := c.GetMyToken()

	//avoc, err := bodies.GetAvoCookie(token, ctx.publicKey)

	//idn, err := c.introspect.Introspect(token, c.client)

	/*if err != nil {
		log.Println(err)
		return nil
	}

	return */

	return nil
}
