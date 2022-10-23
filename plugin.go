package plugin

import (
	"fmt"
	"net/http"
	"strconv"

	config "github.com/jsmzr/boot-config"
	plugin "github.com/jsmzr/boot-plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusPlugin struct{}

type PrometheusProperties struct {
	Port    *int
	Path    *string
	Enabled *bool
}

const prefix = "boot.prometheus"

func (p *PrometheusPlugin) Order() int {
	return 200
}

func (p *PrometheusPlugin) Enabled() bool {
	enabled, ok := config.Get(prefix + ".enabled")
	if ok {
		return enabled.Bool()
	}
	return true
}

func (p PrometheusPlugin) Load() error {
	var properties PrometheusProperties
	var path string
	var port int
	_ = config.Resolve(prefix, &properties)
	if properties.Path == nil {
		path = "/prometheus"
	} else {
		path = *properties.Path
	}
	if properties.Port == nil {
		port = 9080
	} else {
		port = *properties.Port
	}
	fmt.Printf("[BOOT-plugin]  start prometheus by [:%d%s]\n", port, path)
	go func() {
		mux := http.NewServeMux()
		mux.Handle(path, promhttp.Handler())
		err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
		if err != nil {
			fmt.Printf("[BOOT-plugin]  Init prometheus error: %s \n", err.Error())
		}
	}()
	return nil
}

func init() {
	plugin.Register("prometheus", &PrometheusPlugin{})
}
