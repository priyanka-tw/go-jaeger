package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

func main() {
	router := initializeRouter()


	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:          agentEndpointURI,
		CollectorEndpoint:      collectorEndpointURI,
		ServiceName:            "demoooooo",
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(je)

	// oce, _ := ocagent.NewExporter(ocagent.WithInsecure(),
	// 	ocagent.WithReconnectionPeriod(1*time.Second),
	// 	ocagent.WithAddress("localhost:55678"),
	// 	ocagent.WithServiceName("HEALTH_SERVICE"))


	// trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	router.Run()
}

func initializeRouter() *gin.Engine {

	router := gin.Default()
	logrus.Info("initializing router!")

	router.GET("/hello", handleHello)

	return router
}

func handleHello(ctx *gin.Context) {

	_, span := trace.StartSpan(ctx, "/hellooooo-service")
	defer span.End()

	span.AddAttributes(trace.StringAttribute("akey", "avalue"))
	span.Annotate([]trace.Attribute{trace.StringAttribute("annotated", "anAnnotatedValue")}, "annotation check")

	mapOfAtt := map[string]interface{}{}

	span.AddLink(trace.Link{Attributes: mapOfAtt})

	logrus.Info("in hello handler!")
	log.Println("A fatal log!")


	ctx.JSON(http.StatusOK, gin.H{"data": "hello world"})
}
