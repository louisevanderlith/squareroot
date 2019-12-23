package routing

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/squareroot/ctx"
	"github.com/louisevanderlith/squareroot/mix"
	"github.com/louisevanderlith/squareroot/xontrols"
)

func CRUDRouter(r *mux.Router, path string, mxFunc mix.InitFunc, ctrls ...xontrols.Nomad) http.Handler {
	//sub := e.router.(*mux.Router).PathPrefix(path).Subrouter()
	sub := r.PathPrefix(path).Subrouter()

	for _, ctrl := range ctrls {
		ctrlName := getControllerName(ctrl)
		ctrlPath := "/" + strings.ToLower(ctrlName)
		ctrlSub := sub.PathPrefix(ctrlPath).Subrouter()

		//Get
		ctrlSub.Handle("", execute(ctrlName, mxFunc, ctrl.Get)).Methods(http.MethodGet)
		//r.JoinPath(ctrlSub, "", ctrlName, http.MethodGet, required, mxFunc, ctrl.Get)

		//Search
		searchCtrl, isSearch := ctrl.(xontrols.Searchable)

		if isSearch {
			ctrlSub.Handle("/{pagesize:[A-Z][0-9]+}", execute(ctrlName, mxFunc, searchCtrl.Search)).Methods(http.MethodGet)
			//r.JoinPath(ctrlSub, , ctrlName, http.MethodGet, required, mxFunc, searchCtrl.Search)

			ctrlSub.Handle("/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", execute(ctrlName, mxFunc, searchCtrl.Search)).Methods(http.MethodGet)
			//r.JoinPath(ctrlSub, "/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", ctrlName, http.MethodGet, required, mxFunc, searchCtrl.Search)
		}

		//View
		viewCtrl, isView := ctrl.(xontrols.Viewable)

		if isView {
			ctrlSub.Handle("/{key:[0-9]+\x60[0-9]+}", execute(ctrlName, mxFunc, viewCtrl.View)).Methods(http.MethodGet)
			//r.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName, http.MethodGet, required, mxFunc, viewCtrl.View)
		}

		//Create
		createCtrl, isCreate := ctrl.(xontrols.Createable)

		if isCreate {
			ctrlSub.Handle("", execute(ctrlName, mxFunc, createCtrl.Create)).Methods(http.MethodPost)
			//r.JoinPath(ctrlSub, "", ctrlName, http.MethodPost, required, mxFunc, createCtrl.Create)
		}

		//Update
		updatCtrl, isUpdate := ctrl.(xontrols.Updateable)

		if isUpdate {
			ctrlSub.Handle("/{key:[0-9]+\x60[0-9]+}", execute(ctrlName, mxFunc, updatCtrl.Update)).Methods(http.MethodPut)
			//r.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName, http.MethodPut, required, mxFunc, updatCtrl.Update)
		}

		//Delete
		delCtrl, isDelete := ctrl.(xontrols.Deleteable)

		if isDelete {
			ctrlSub.Handle("/{key:[0-9]+\x60[0-9]+}", execute(ctrlName, mxFunc, delCtrl.Delete)).Methods(http.MethodDelete)
			//r.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName, http.MethodDelete, required, mxFunc, delCtrl.Delete)
		}

		//Queries
		qryCtrl, isQueried := ctrl.(xontrols.Queries)

		if isQueried {
			for qkey, qval := range qryCtrl.AcceptsQuery() {
				ctrlSub.Queries(qkey, qval)
			}
		}
	}

	return sub
}

func execute(name string, mxFunc mix.InitFunc, srvFunc ctx.ServeFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		cntx := ctx.New(resp, req)

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		//status, data := srvFunc(cntx)
		//mxer := mxFunc(cntx.RequestURI(), srvFunc)

		err := cntx.Serve(mxFunc, srvFunc)

		if err != nil {
			log.Panicln(err)
		}
	}
}

func getControllerName(ctrl xontrols.Nomad) string {
	tpe := reflect.TypeOf(ctrl).String()
	lstDot := strings.LastIndex(tpe, ".")

	if lstDot != -1 {
		return tpe[(lstDot + 1):]
	}

	return tpe
}
