package core

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func EnableSwagger(container *restful.Container) {
	config := swagger.Config{
		WebServices:    container.RegisteredWebServices(),
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/internal/apidocs.json",

		SwaggerPath:     "/internal/apidocs/",
		SwaggerFilePath: "./swagger",
	}
	swagger.RegisterSwaggerService(config, container)
}
