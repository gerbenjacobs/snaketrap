package core

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-selfdiagnose"
)

func handleSelfdiagnose(r *restful.Request, w *restful.Response) {
	ctx := selfdiagnose.NewContext()
	// prepare for ReportHttpRequest
	ctx.Variables["http.request"] = r
	var reporter selfdiagnose.Reporter
	if strings.HasSuffix(r.Request.URL.Path, ".json") || r.Request.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		reporter = selfdiagnose.JSONReporter{w}
	} else if strings.HasSuffix(r.Request.URL.Path, ".xml") || r.Request.URL.Query().Get("format") == "xml" {
		w.Header().Set("Content-Type", "application/xml")
		reporter = selfdiagnose.XMLReporter{w}
	} else {
		w.Header().Set("Content-Type", "text/html")
		reporter = selfdiagnose.HtmlReporter{w}
	}
	selfdiagnose.DefaultRegistry.RunWithContext(reporter, ctx)
}

func SetupSelfdiagnose(app *restful.WebService) {
	// Register routes
	app.Route(app.GET("/internal/selfdiagnose.html").To(handleSelfdiagnose))
	app.Route(app.GET("/internal/selfdiagnose.json").To(handleSelfdiagnose))
	app.Route(app.GET("/internal/selfdiagnose.xml").To(handleSelfdiagnose))

	// Register tasks
	selfdiagnose.Register(ReportHTTPRequest{})
	selfdiagnose.Register(ReportHostname{})
	cpu := selfdiagnose.ReportMessage{
		Message: fmt.Sprintf("%d CPU available. %d goroutines active", runtime.NumCPU(), runtime.NumGoroutine()),
	}
	cpu.SetComment("Num CPU")
	selfdiagnose.Register(cpu)
}

type ReportBuildAndDate struct {
	Version string
	Date    string
}

func (r ReportBuildAndDate) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	result.Passed = true
	result.Reason = fmt.Sprintf("Version: %s Date: %s", r.Version, r.Date)
}

func (r ReportBuildAndDate) Comment() string {
	return "Build information"
}

type ReportHTTPRequest struct{}

func (r ReportHTTPRequest) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	request, ok := ctx.Variables["http.request"].(*restful.Request)
	req := request.Request
	if !ok {
		result.Passed = false
		result.Reason = "missing variable 'http.request'"
		return
	}
	var buf bytes.Buffer
	for k, v := range req.Header {
		buf.WriteString(fmt.Sprintf("%s = %s<br/>", k, v))
	}
	result.Passed = true
	result.Reason = buf.String()
	result.Severity = selfdiagnose.SeverityNone
}

func (r ReportHTTPRequest) Comment() string { return "HTTP Request" }

type ReportHostname struct{}

func (r ReportHostname) Comment() string { return "Hostname" }

func (r ReportHostname) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	h, err := os.Hostname()
	result.Severity = selfdiagnose.SeverityNone
	result.Passed = err == nil
	result.Reason = h
}
