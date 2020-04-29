package main

import (
	jaegerExp "contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	jaegerClient "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

func main() {
	router := initializeRouter()

	withJaegerClientOpentracing() // add the exporter

	err := router.Run()
	if err != nil {
		log.Fatal("error while starting the service")
	}
}

func initializeRouter() *gin.Engine {

	router := gin.Default()
	logrus.Info("initializing router!")

	router.GET("/hello", handleHelloWithOpenTracing) // add respective handler func for the exporter

	return router
}

func handleHelloWithOpenCensus(ctx *gin.Context) {

	_, span := trace.StartSpan(ctx, "/hello-service")

	defer span.End()

	span.AddAttributes(trace.StringAttribute("aKey", "aValue"))
	span.Annotate([]trace.Attribute{trace.StringAttribute("anAnnotatedKey", "anAnnotatedValue")}, "an annotation string")

	mapOfAtt := map[string]interface{}{}

	span.AddLink(trace.Link{Attributes: mapOfAtt})

	logrus.Info("in hello handler!")

	ctx.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

func handleHelloWithOpenTracing(ctx *gin.Context){
	tracer := opentracing.GlobalTracer()
	tracer.StartSpan("hello-opertation")

	logrus.Info("logrus info in handle func")
	log.Println("log println in handle func")

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

func withOCJaegerExporter() {
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"

	je, err := jaegerExp.NewExporter(jaegerExp.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "HEALTH_SERVICE_JAEGER_EXPORTER",
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

func withJaegerClientOpentracing() {

	jaegerConfig := &config.Configuration{
		Sampler:  &config.SamplerConfig{Type: "const", Param: 1},
		Reporter: &config.ReporterConfig{LogSpans: true},
	}

	closer, err := jaegerConfig.InitGlobalTracer(
		"HEALTH_SERVICE_JAEGER_CLIENT",
		config.Logger(jaegerClient.StdLogger))
	if err != nil {
		logrus.Fatal("error in jaeger init", err)
	}

	defer closer.Close()
}
