package plugin

import (
	"fmt"
	"net/http"
	"strconv"

	plugin "github.com/jsmzr/boot-plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

type PrometheusPlugin struct{}

const configPrefix = "boot.prometheus"

var defaultConfig map[string]interface{} = map[string]interface{}{"enabled": true, "order": 30, "port": 9080, "path": "/prometheus"}

func (p *PrometheusPlugin) Order() int {
	return viper.GetInt(configPrefix + ".order")
}

func (p *PrometheusPlugin) Enabled() bool {
	return viper.GetBool(configPrefix + ".enabled")
}

func (p PrometheusPlugin) Load() error {

	path := viper.GetString(configPrefix + ".path")
	port := viper.GetInt(configPrefix + ".port")
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
	for key := range defaultConfig {
		viper.SetDefault(configPrefix+"."+key, defaultConfig[key])
	}
	plugin.Register("prometheus", &PrometheusPlugin{})
}
