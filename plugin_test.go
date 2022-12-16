package plugin

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestOrder(t *testing.T) {
	plugin := PrometheusPlugin{}
	if plugin.Order() != 30 {
		t.Fatal("default order should be 30")
	}
	newOrder := 40
	viper.Set(configPrefix+".order", newOrder)
	if plugin.Order() != newOrder {
		t.Fatalf("order should be %d", newOrder)
	}
}

func TestEnabled(t *testing.T) {
	plugin := PrometheusPlugin{}
	if !plugin.Enabled() {
		t.Fatal("default enable should be true")
	}
	viper.Set(configPrefix+".enabled", false)
	if plugin.Enabled() {
		t.Fatal("enable should be false")
	}

}

func TestLoad(t *testing.T) {
	plugin := PrometheusPlugin{}
	if err := plugin.Load(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	checkMertrics(viper.GetInt(configPrefix+".port"), viper.GetString(configPrefix+".path"))
}

func checkMertrics(port int, path string) bool {
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port) + path)
	return err == nil && resp.StatusCode == 200
}
