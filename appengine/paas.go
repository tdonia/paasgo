package paas

import (
  "appengine"
  "net/http"
)

type Context appengine.Context

func NewContext(r *http.Request) Context{
  return appengine.NewContext(r)
}

func IsDevAppServer() bool{
  return appengine.IsDevAppServer()
}
