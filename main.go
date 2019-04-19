package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

var (
	port        = 1314
	destination = ""
	client      = &http.Client{}
)

type handler struct{}

func (s *handler) CORS(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "" {
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	w.Header().Add("Access-Control-Max-Age", "3600")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With, Origin, Content-Type, X-Auth-Token, Authorization")
}

func (s *handler) WriteError(w http.ResponseWriter, err error) {
	log.Println("locap error", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

func (s *handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var url = destination + r.URL.RequestURI()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	req, err := http.NewRequest(r.Method, url, bytes.NewBuffer(body))
	if err != nil {
		s.WriteError(w, err)
		return
	}
	for k, v := range r.Header {
		req.Header.Set(k, strings.Join(v, ","))
	}
	resp, err := client.Do(req)
	if err != nil {
		s.WriteError(w, err)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.WriteError(w, err)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.RequestURI())
	s.CORS(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		s.Transfer(w, r)
	}
}

func serverHTTP() {
	var p = fmt.Sprintf(":%d", port)
	log.Println("Accepting connections from:", p)
	http.ListenAndServe(p, &handler{})
}

func main() {
	app := cli.NewApp()
	app.Name = "locap"
	app.Usage = "locap"
	app.UsageText = "locap [options]"
	app.HideVersion = true
	app.Author = "Yoann Cerda"
	app.Email = "tuxlinuxien@gmail.com"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 1314,
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "",
		},
	}
	app.Action = func(c *cli.Context) error {
		port = c.Int("port")
		destination = c.String("destination")
		if destination == "" {
			log.Println("--destination should be set")
			return nil
		}
		serverHTTP()
		return nil
	}

	app.Run(os.Args)
}
