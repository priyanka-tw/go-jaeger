package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/stdout"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

func main() {
	fmt.Println("in main")
	router := newRouter()

	initTracer()

	err := router.Run()
	if err != nil {
		log.Fatal("error while starting the service")
	}
}

func initTracer() {
	//exporter, err := jaeger.NewRawExporter(
	//	jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
	//)

	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal("unable to init tracer")
	}

	config := sdkTrace.WithConfig(sdkTrace.Config{DefaultSampler: sdkTrace.AlwaysSample()})

	provider, err := sdkTrace.NewProvider(config, sdkTrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal("unable to init provider")
	}
	global.SetTraceProvider(provider)
}

func newRouter() *gin.Engine {
	router := gin.Default()
	logrus.Info("initializing router!")

	router.GET("/hey", routeHandler)

	return router
}

func routeHandler(ctx *gin.Context) {

	tracer := global.TraceProvider().Tracer("/hey")
	_, span := tracer.Start(ctx, "hey-span")

	span.SetAttributes(core.KeyValue{Key: "aNewKey", Value: core.String("nNewValue")})

	log.Println("in router handler")

	ctx.JSON(http.StatusOK, "hey there!")
}
