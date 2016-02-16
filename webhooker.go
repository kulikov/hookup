package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	app := cli.NewApp()
	app.Name = "webhooker"
	app.Usage = "Start Webhook Server"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 9090,
		},
		cli.StringFlag{
			Name:  "handlers",
			Usage: "Path to dir with webhook handlers scripts",
			Value: "/etc/webhooker.d/",
		},
		cli.StringSliceFlag{
			Name:  "github-secret",
			Usage: "Secret keys for github hooks",
		},
	}

	app.Action = func(c *cli.Context) {
		StartWebHookServer(c.Int("port"), c.String("handlers"), c)
	}

	app.Run(os.Args)
}

func StartWebHookServer(port int, handlersDir string, ctx *cli.Context) {
	ec := echo.New()
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())

	ec.Post("/github/events", func(c *echo.Context) error {
		defer c.Request().Body.Close()

		eventType := c.Request().Header.Get("X-GitHub-Event")
		payload, err := ioutil.ReadAll(c.Request().Body)

		if err := checkSecret(c.Request(), ctx.StringSlice("github-secret")); err != nil {
			return err
		}

		log.Printf("Receive github event '%s': \n%s\n", eventType, string(payload))

		for _, handler := range collectHandlers(handlersDir) {
			go runHandler(handler, "github", eventType, string(payload))
		}

		return err
	})

	log.Printf("Starting web hook server on :%v\n", port)

	ec.Run(":" + strconv.Itoa(port))
}

func checkSecret(req *http.Request, secrets []string) error {
	if len(secrets) > 0 {
		sig := req.Header.Get("X-Hub-Signature")

		if sig == "" {
			return http.Error("403 Forbidden - Missing X-Hub-Signature required for HMAC verification", http.StatusForbidden)
		}

		mac := hmac.New(sha1.New, []byte(s.Secret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)
		if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
			http.Error(w, "403 Forbidden - HMAC verification failed", http.StatusForbidden)
			return
		}
	}
}

func collectHandlers(handlersDir string) []string {
	handlers := make([]string, 0)

	err := filepath.Walk(handlersDir, func(path string, f os.FileInfo, err error) error {
		if err == nil && !f.IsDir() {
			handlers = append(handlers, path)
		}
		return err
	})

	if err != nil || len(handlers) == 0 {
		log.Printf("Handlers not found: %v\n", err)
	}

	return handlers
}

func runHandler(handler string, source string, eventType string, payload string) {
	log.Printf("Run %s", handler)

	out, err := exec.Command("/bin/bash", handler, "--source", source, "--event", eventType, "--payload", payload).CombinedOutput()

	if err != nil {
		log.Printf("Error: %s\n\n", err)
	} else {
		log.Printf("Out: \n%s\n\n", string(out))
	}
}
