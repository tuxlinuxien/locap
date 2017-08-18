package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
)

var (
	port        = 1314
	destination = ""
	client      = &http.Client{}
)

func transfer(c echo.Context) error {
	var url = destination + c.Request().URL.RequestURI()
	defer c.Request().Body.Close()
	body, _ := ioutil.ReadAll(c.Request().Body)
	req, err := http.NewRequest(c.Request().Method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("locap error", err)
		return err
	}
	for k, v := range c.Request().Header {
		req.Header.Set(k, strings.Join(v, ","))
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("locap error", err)
		return err
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		c.Response().Header().Set(k, strings.Join(v, ","))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("locap error", err)
		return err
	}
	return c.String(resp.StatusCode, string(respBody))
}

func server() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Any("*", transfer)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
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
		server()
		return nil
	}

	app.Run(os.Args)
}
