package main

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

func main() {
	router := initializeRouter()

	withOCJaegerExporter()

	err := router.Run()
	if err != nil {
		log.Fatal("error while starting the service")
	}
}

func initializeRouter() *gin.Engine {

	router := gin.Default()
	logrus.Info("initializing router!")

	router.GET("/hello", handleHello)

	return router
}

func handleHello(ctx *gin.Context) {

	_, span := trace.StartSpan(ctx, "/hello-service")

	defer span.End()

	span.AddAttributes(trace.StringAttribute("aKey", "aValue"))
	span.Annotate([]trace.Attribute{trace.StringAttribute("anAnnotatedKey", "anAnnotatedValue")}, "an annotation string")

	mapOfAtt := map[string]interface{}{}

	span.AddLink(trace.Link{Attributes: mapOfAtt})

	logrus.Info("in hello handler!")

	ctx.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

func withOCAgentExporter() {
	oce, _ := ocagent.NewExporter(ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(1*time.Second),
		ocagent.WithAddress("localhost:55678"),
		ocagent.WithServiceName("HEALTH_SERVICE_OCAGENT"))

	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

func withOCJaegerExporter(){
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:          agentEndpointURI,
		CollectorEndpoint:      collectorEndpointURI,
		ServiceName:            "HEALTH_SERVICE_JAEGER_EXPORTER",
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
