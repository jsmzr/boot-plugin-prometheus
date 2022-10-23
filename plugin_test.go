package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	config "github.com/jsmzr/boot-config"
)

type TestConfig struct{}

type TestContainer struct {
	content string
}

var testPort = 9000
var testPath = "/mertrics"
var testProperties = PrometheusProperties{Port: &testPort, Path: &testPath}

func (t *TestConfig) Load(filename string) (config.Configer, error) {
	data := make(map[string]PrometheusProperties)
	data["prometheus"] = testProperties
	resouce := make(map[string]interface{})
	resouce["boot"] = data
	configJson, toJsonErr := json.Marshal(resouce)
	if toJsonErr != nil {
		return nil, toJsonErr
	}
	return &TestContainer{content: string(configJson)}, nil
}

func (t *TestContainer) GetJson() string {
	return t.content
}

func checkMertrics(port int, path string) bool {
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port) + path)
	fmt.Println(resp, err)
	return err == nil && resp.StatusCode == 200
}

func TestLoad(t *testing.T) {
	prometheusPlugin := PrometheusPlugin{}
	if err := prometheusPlugin.Load(); err != nil {
		t.Fatal("load plugin should not be error")
	}
	// TODO maybe > 2s
	time.Sleep(2 * time.Second)
	if !checkMertrics(9080, "/prometheus") {
		t.Fatal("check should be true")
	}

}

func TestLoadByConfig(t *testing.T) {
	config.Register("test", &TestConfig{})
	if err := config.InitInstance("test", "ignore"); err != nil {
		t.Fatal(err)
	}
	prometheusPlugin := PrometheusPlugin{}
	if err := prometheusPlugin.Load(); err != nil {
		t.Fatal("load plugin should not be error")
	}
	// TODO maybe > 2s
	time.Sleep(2 * time.Second)
	if !checkMertrics(*testProperties.Port, *testProperties.Path) {
		t.Fatal("check should be true")
	}

}
