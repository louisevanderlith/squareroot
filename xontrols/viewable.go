package xontrols

import "github.com/louisevanderlith/squareroot/ctx"

//Viewable handles controls that handle and view items.
type Viewable interface {
	Nomad
	View(ctx.Requester) (int, interface{})
}
