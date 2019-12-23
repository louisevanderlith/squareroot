package ctx

import "github.com/louisevanderlith/squareroot/mix"

type Responder interface {
	Requester
	SetHeader(key string, val string) //SetHeader sets a value on the Response Header
	SetStatus(code int)               //SetStatus set the final Response Status

	//Serve(int, mix.Mixer) error
	Serve(mix.InitFunc, ServeFunc) error
}
