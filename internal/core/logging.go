package core

import (
	"github.com/emicklei/go-restful"
	"github.com/inconshreveable/log15"
)

func AccessLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var username = "-"
		if req.Request.URL.User != nil {
			if name := req.Request.URL.User.Username(); name != "" {
				username = name
			}
		}
		chain.ProcessFilter(req, resp)

		dataList := []interface{}{
			"logtype", "accesslog",
			"client", req.Request.RemoteAddr,
			"verb", req.Request.Method,
			"httpversion", req.Request.Proto,
			"referrer", req.Request.Referer(),
			"agent", req.Request.UserAgent(),
			"mime", resp.Header().Get("Content-Type"),
			"site", req.Request.Host,
			"user", username,
		}
		log15.Debug(req.Request.RequestURI, dataList...)
	}
}
