# go-lib-pub-tlp

[![Software License](https://img.shields.io/badge/license-APACHE%202.0-CB212D.svg)](LICENSE)
[![Software Version](https://img.shields.io/badge/version-1.0.0-green.svg)]()

_This lightweight library can be used for various microservices (grpc support available also) or standalone applications in the golang context that want to use active profiling, tracing or debugging using stackdriver or jaeger._

## functional scope

| function      | params  | description  |
|:--------------|:--------|:-------------|
| InitTracing   | `serviceName` `<string>` <br/> `jaegerServiceAddress` `<string>` <br/> `log` `<*logrus.Logger` | activates tracing for stackdriver and Jaeger |
| InitProfiling | `service` `<string>` <br/> `version` `<string>` <br/> `log` `<*logrus.Logger` | activates profiling for stackdriver |

## example of use

```
import {
    tlp "github.com/RelicFrog/go-lib-pub-tlp"
    "github.com/sirupsen/logrus"
    "os"
}

var (
	log *logrus.Logger
    err error
)

func init() {

	//
	// -- define logging mechanics --
	//

	log = logrus.New()
	log.Level = logrus.DebugLevel
	if metaDebugMode := os.Getenv("ENABLE_DEBUG"); metaDebugModeEnabled == "1" {
		log.Level = logrus.ErrorLevel
	}

	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg: "message",
		},
		TimestampFormat: time.RFC3339Nano,
	};  log.Out = os.Stdout

	//
	// -- define config bound mechanics for logging/profiling --
	//

	if metaTracerDisabled := os.Getenv("ENABLE_TRACING"); metaTracerEnabled == "1" {
		log.Infof("%s: tracing enabled.","myservice")
		go tlp.InitTracing("myservice",os.Getenv("JAEGER_SERVICE_ADDR"),log)
	}

	if metaProfilerDisabled := os.Getenv("ENABLE_PROFILER"); metaProfilerEnabled == "1" {
		log.Infof("%s: profiling enabled.","myservice")
		go tlp.InitProfiling("myserver","1.0.0",log)
	}
}
```

## links

- https://www.jaegertracing.io/
- https://cloud.google.com/go/docs/stackdriver

## copyright

(c) 2020-now Patrick Paechnatz <post@dunkelfrosch.com> All Rights Reserved. This code, all used media assets, vendor libs are proprietary protected and aught to the copyright holder.

## license

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
