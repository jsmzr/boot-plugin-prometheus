package plugin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jsmzr/bootstrap-config/config"
	"github.com/jsmzr/bootstrap-plugin/plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusPlugin struct{}

type PrometheusProperties struct {
	Port int
	Path string
}

func (p *PrometheusPlugin) Order() int {
	return 200
}

func (p PrometheusPlugin) Load() error {
	var properties PrometheusProperties

	if err := config.Resolve("prometheus", &properties); err != nil {
		properties.Port = 9000
		properties.Path = "/prometheus"
	}
	go func() {
		mux := http.NewServeMux()
		mux.Handle(properties.Path, promhttp.Handler())
		fmt.Printf("prometheus init [:%d%s]\n", properties.Port, properties.Path)
		err := http.ListenAndServe(":"+strconv.Itoa(properties.Port), mux)
		if err != nil {
			fmt.Println("Prometheus plugin init err:", err)
		}
	}()
	return nil
}

func init() {
	plugin.Register("prometheus", &PrometheusPlugin{})
}
