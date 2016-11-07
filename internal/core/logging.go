package core

import (
	"github.com/emicklei/go-restful"
	"github.com/inconshreveable/log15"
)

func AccessLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		chain.ProcessFilter(req, resp)

		dataList := []interface{}{
			"logtype", "accesslog",
			"client", req.Request.RemoteAddr,
			"verb", req.Request.Method,
			"httpversion", req.Request.Proto,
			"agent", req.Request.UserAgent(),
			"mime", resp.Header().Get("Content-Type"),
			"site", req.Request.Host,
		}
		log15.Debug(req.Request.RequestURI, dataList...)
	}
}
