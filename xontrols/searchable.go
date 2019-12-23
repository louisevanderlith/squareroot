package xontrols

import "github.com/louisevanderlith/squareroot/ctx"

//Searchable handles controls that handle search submissions[GET]
type Searchable interface {
	Nomad
	Search(ctx.Requester) (int, interface{})
}
