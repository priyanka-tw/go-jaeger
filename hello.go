package main

import (
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

	oce, _ := ocagent.NewExporter(ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(1*time.Second),
		ocagent.WithAddress("localhost:55678"),
		ocagent.WithServiceName("HEALTH_SERVICE"))


	trace.RegisterExporter(oce)
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
