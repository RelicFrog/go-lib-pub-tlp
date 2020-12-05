package tlp

import (
	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"time"
)

//
// -- Tracing, Logging & Profiling Features (internal) --
//

func initJaegerTracing(serviceName string, jaegerServiceAddress string, log *logrus.Logger) {
	// common.GetDotEnvVariable("JAEGER_SERVICE_ADDR")
	metaJaegerServiceAddress := jaegerServiceAddress
	if metaJaegerServiceAddress == "" {
		log.Info("jaeger initialization disabled.")
		return
	}

	// Register the Jaeger exporter to be able to retrieve the collected spans.
	log.Info(metaJaegerServiceAddress)
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: fmt.Sprintf("http://%s", metaJaegerServiceAddress),
		Process: jaeger.Process{
			ServiceName: serviceName,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	trace.RegisterExporter(exporter)
	log.Info("jaeger initialization completed.")
}

func initStackdriverTracing(log *logrus.Logger) {
	// since they are not sharing packages.
	for i := 1; i <= 3; i++ {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{})
		if err != nil {
			log.Warnf("failed to initialize Stackdriver exporter: %+v", err)
		} else {
			trace.RegisterExporter(exporter)
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
			log.Info("registered Stackdriver tracing")

			// Register the views to collect server stats.
			initStackdriverStats(exporter, log)
			return
		}

		d := time.Second * 10 * time.Duration(i)
		log.Infof("sleeping %v to retry initializing Stackdriver exporter", d)
		time.Sleep(d)
	}

	log.Warn("could not initialize Stackdriver exporter after retrying, giving up :-/")
}

func initStackdriverStats(exporter *stackdriver.Exporter, log *logrus.Logger) {
	view.SetReportingPeriod(60 * time.Second)
	view.RegisterExporter(exporter)

	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Info("error registering default server views")
	} else {
		log.Info("registered default server views")
	}
}

//
// -- Public Functional Scope --
//

func InitTracing(serviceName string, jaegerServiceAddress string, log *logrus.Logger) {
	initJaegerTracing(serviceName, jaegerServiceAddress, log)
	initStackdriverTracing(log)
}

func InitProfiling(service string, version string, log *logrus.Logger) {
	for i := 1; i <= 3; i++ {
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: version,
			// ProjectID:   service, // *** ProjectID must be set if not running on GCP ***
		}); err != nil {
			log.Warnf("failed to start profiler: %+v", err)
		} else {
			log.Info("started Stackdriver profiler")
			return
		}
		d := time.Second * 10 * time.Duration(i)
		log.Infof("sleeping %v to retry initializing Stackdriver profiler", d)
		time.Sleep(d)
	}

	log.Warn("could not initialize Stackdriver profiler after retrying, giving up")
}

