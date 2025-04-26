package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Page struct {
	Title  string
	Style  string
	Year   string
	Reload int
	Date   string
}

type Config struct {
	Target          string       `yaml:"target"`
	ServerName      string       `yaml:"server_name"`
	Username        string       `yaml:"username"`
	Password        string       `yaml:"password"`
	BearerToken     string       `yaml: bearer_token`
	TLS_Skip_Verify bool         `yaml:"tls_skip_verify"`
	ListMetric      []ListMetric `yaml:"list_metrics"`
}

type ListMetric struct {
	MetricName string `yaml:"metric_name"`
	Name       string `yaml:"name"`
	Size       string `yaml:"size"`
	Max        int    `yaml:"max"`
}

type Server struct {
	Port             int      `yaml:"port"`
	Timeout_Duration int      `yaml:"timeout"`
	BufferSize       int      `yaml:"buffer_size"`
	Certfile         string   `yaml:"certfile"`
	Keyfile          string   `yaml:"keyfile"`
	Config           []Config `yaml:"config"`
	ScrapeInterval   int      `yaml:"scrape_interval"`
	UTC              int      `yaml:"utc"`
	UpdateMetricData string
	Now              string
}

type Content struct {
	Page    Page
	Content string
}

type ListServer struct {
	ServerName string
	ListMetric []Metric
}

type Metric struct {
	MetricName string
	Value      string
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "config.file", "", "Location of config file")
	flag.Parse()
	if len(fileName) == 0 {
		log.Println("Usage: \n./status-page --config.file=server.yaml")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	var server Server
	err = yaml.Unmarshal(data, &server)
	if err != nil {
		log.Fatal(err)
		return
	}
	Website(server)
}
