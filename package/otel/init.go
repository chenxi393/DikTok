package otel

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

// 还会开启 内存检测 不过内存检测放在一个进程不合适
func Init(url, serviceName string) func() {
	tp := initTracerProvider(url, serviceName)
	// mp := initMeterProvider(url, serviceName)
	// 内存监测 还没啥效果
	runtimeMemory()
	// 这里 defer 直接关了 应该返回 然后defer
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
		// if err := mp.Shutdown(context.Background()); err != nil {
		// 	log.Printf("Error shutting down meter provider: %v", err)
		// }
	}
}

// 挂了怎么办
func runtimeMemory() {
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}
}
