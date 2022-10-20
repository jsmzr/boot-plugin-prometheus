package plugin

import (
	"fmt"
	"net/http"
	"strconv"

	config "github.com/jsmzr/bootstrap-config"
	plugin "github.com/jsmzr/bootstrap-plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusPlugin struct{}

type PrometheusProperties struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

func (p *PrometheusPlugin) Order() int {
	return 200
}

func (p PrometheusPlugin) Load() error {
	var properties PrometheusProperties

	if err := config.Resolve("prometheus", &properties); err != nil {
		properties.Port = 9080
		properties.Path = "/prometheus"
	}
	go func() {
		mux := http.NewServeMux()
		mux.Handle(properties.Path, promhttp.Handler())
		err := http.ListenAndServe(":"+strconv.Itoa(properties.Port), mux)
		if err != nil {
			fmt.Printf("[BOOTSTRAP-plugin]  Init prometheus error: %s \n", err.Error())
		}
	}()
	return nil
}

func init() {
	plugin.Register("prometheus", &PrometheusPlugin{})
}
