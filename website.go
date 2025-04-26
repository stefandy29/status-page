package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

func (srv *Server) Home(w http.ResponseWriter, r *http.Request) {
	srv.TemplateSkeleton(w, srv.UpdateMetricData)
	r.Body.Close()
	return
}

func (srv *Server) TemplateSkeleton(w http.ResponseWriter, content string) error {
	tmplt, err := template.New("article").Parse(skeleton)
	if err != nil {
		http.Error(w, "template not found", http.StatusNotFound)
		return err
	}
	// set timezone,
	// now := time.Now().Format("2 Jan 2006 15:04")
	value := Content{
		Page: Page{
			Title:  "Status Server",
			Style:  style_html,
			Year:   strconv.Itoa(time.Now().Year()),
			Reload: srv.ScrapeInterval,
			Date:   srv.Now,
		},
		Content: content,
	}
	tmplt.ExecuteTemplate(w, "article", value)
	return nil
}

func (srv *Server) FetchExporter(target string, username string, password string, bearer string, tls_skip_verify bool) string {
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(srv.Timeout_Duration) * time.Second,
			KeepAlive: time.Duration(srv.Timeout_Duration) * time.Second,
		}).Dial,
		ReadBufferSize:  1024 * srv.BufferSize,
		WriteBufferSize: 1024 * srv.BufferSize,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: tls_skip_verify,
		},
	}
	client_settings := &http.Client{
		Transport: tr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(srv.Timeout_Duration)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}
	if bearer != "" {
		req.Header.Add("Authorization", bearer)
	}

	resp, err := client_settings.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if len(string(body)) == 0 {
		log.Println("body empty")
		return ""
	}
	if resp.StatusCode != 200 {
		log.Println("error status code is not" + fmt.Sprint(resp.StatusCode))
		return ""
	}
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	resp.Body.Close()
	return string(body)
}

func (srv *Server) FetchMetricData(config Config) []string {
	m := make(map[string]string)
	metrics := srv.FetchExporter(config.Target, config.Username, config.Password, config.BearerToken, config.TLS_Skip_Verify)
	if metrics == "" {
		return nil
	}
	re := regexp.MustCompile(`(?m)^[^#\n]+`)
	for _, match := range re.FindAllString(metrics, -1) {
		metric := strings.Split(match, " ")
		if len(metric) >= 2 {
			m[metric[0]] = metric[1]
		}
	}

	//filter
	var metric_slice []string
	for k, v := range m {
		for _, y := range config.ListMetric {
			if y.MetricName == k {
				var name string
				if y.Name == "" {
					name = k
				} else {
					name = y.Name
				}
				metric_slice = append(metric_slice, srv.Metric(name, v, y.Max, y.Size))
			}
		}
	}
	if metric_slice != nil {
		slices.Sort(metric_slice)
	}
	return metric_slice
}

func (srv *Server) Content() string {
	var wg sync.WaitGroup
	var list_metric []string
	var content []string

	for _, config := range srv.Config {
		wg.Add(1)
		go func() {
			list_metric = srv.FetchMetricData(config)
			defer wg.Done()
		}()
	}
	wg.Wait()

	for _, config := range srv.Config {
		if list_metric != nil {
			content = append(content, srv.ListServer(config.ServerName, srv.FetchMetricData(config)))
		}
	}
	return strings.Join(content, "")
}

func (srv *Server) ListServer(server_name string, list_metric []string) string {
	server_list, err := template.New("").Parse(server_list_html)
	if err != nil {
		log.Fatal(err)
	}
	var tpl bytes.Buffer
	err = server_list.Execute(&tpl, map[string]interface{}{
		"ServerName":  server_name,
		"ListMetrics": list_metric,
	},
	)
	if err != nil {
		log.Fatal(err)
	}
	result := tpl.String()
	return result
}

func (srv *Server) Metric(key string, value string, max_size int, type_size string) string {
	metric, err := template.New("").Funcs(template.FuncMap{
		"Divide": func(a string) float64 {
			value_str, _ := strconv.ParseFloat(a, 64)
			value := value_str / float64(max_size)
			switch {
			case max_size == 0:
				value = 0
			case value > 1:
				value = 1
			case value < 0:
				value = 0
			}
			return value
		},
	}).Parse(metric_html)
	if err != nil {
		log.Fatal(err)
	}
	var tpl bytes.Buffer
	err = metric.Execute(&tpl, map[string]interface{}{
		"Key":       key,
		"Value":     value,
		"Type_Size": type_size,
	},
	)
	if err != nil {
		log.Fatal(err)
	}
	result := tpl.String()
	return result
}

func (srv *Server) UpdateMetric() {
	if srv.UpdateMetricData == "" {
		srv.UpdateMetricData = srv.Content()
	}
	ticker := time.NewTicker(time.Second * time.Duration(srv.ScrapeInterval))

	go func() {
		for range ticker.C {
			srv.UpdateMetricData = srv.Content()
			srv.Now = time.Now().UTC().Add(time.Hour * time.Duration(srv.UTC)).Format("2 Jan 2006 15:04")
		}
	}()
}

func checkFileExist(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func Website(srv Server) {
	srv.UpdateMetric()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	mux := http.NewServeMux()
	timeout := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         ":" + strconv.Itoa(srv.Port),
		Handler:      mux,
	}
	mux.HandleFunc("/", srv.Home)

	log.Printf("%s", "Starting website at Port "+strconv.Itoa(srv.Port))
	go func() {
		if checkFileExist(srv.Certfile) && checkFileExist(srv.Keyfile) {
			log.Fatal(timeout.ListenAndServeTLS(srv.Certfile, srv.Keyfile))
		} else {
			log.Fatal(timeout.ListenAndServe())
		}
	}()
	wg.Wait()
	defer wg.Done()
}
