package xontrols

import "github.com/louisevanderlith/squareroot/ctx"

//Createable handles controls that create content [POST]
type Createable interface {
	Searchable
	Create(ctx.Requester) (int, interface{})
}
