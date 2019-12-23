package mix

import (
	"io"
)

//Mixer is used by the Contexer to ApplyHeaders and Write the Response from the Reader
type Mixer interface {
	Reader() (io.Reader, error)
	Headers() map[string]string
}

//InitFunc is a function that returns a Mixer that is able to serve requests.
//@name is the name of the current route
//@obj is the data that has to be returned
//@d is the identity of the services' state
type InitFunc func(name string, obj interface{}) Mixer

//type InitFunc func(name string, obj interface{}, d *element.Identity, avoc *models.ClaimIdentity) Mixer

//DefaultHeaders returns a set of Headers that apply to all mixers
func DefaultHeaders() map[string]string {
	result := make(map[string]string)

	result["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	result["Access-Control-Allow-Credentials"] = "true"
	result["Server"] = "square"
	result["X-Content-Type-Options"] = "nosniff"

	return result
}
