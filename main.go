package main

import (
	"./hookup"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	app := cli.NewApp()
	app.Name = "hookup"
	app.Usage = "Start Webhook Server"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 9090,
		},
		cli.StringFlag{
			Name:  "handlers",
			Usage: "Path to dir with webhook handlers scripts",
			Value: "/etc/hookup.d/",
		},
	}

	app.Action = func(c *cli.Context) {
		hookup.StartWebhookServer(c.Int("port"), func(source string, eventType string, payload string) {
			for _, cmd := range findHandlerCmds(c.String("handlers")) {
				go execHandler(cmd, source, eventType, payload)
			}
		})
	}

	app.Run(os.Args)
}

func findHandlerCmds(handlersDir string) []string {
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

func execHandler(handler string, source string, eventType string, payload string) {
	log.Printf("Run %s", handler)

	out, err := exec.Command("/bin/bash", handler, "--source", source, "--event", eventType, "--payload", payload).CombinedOutput()

	if err != nil {
		log.Printf("Error: %s\n\n", err)
	} else {
		log.Printf("Out: \n%s\n\n", string(out))
	}
}
