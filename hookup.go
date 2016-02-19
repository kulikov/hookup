package hookup

import (
	"io/ioutil"
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)


func StartWebhookServer(port int, handler func(source string, eventType string, payload string)) {
	ec := echo.New()
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())

	ec.Post("/github/events", func(c *echo.Context) error {
		defer c.Request().Body.Close()

		eventType := c.Request().Header.Get("X-GitHub-Event")
		payload, err := ioutil.ReadAll(c.Request().Body)

		log.Printf("Receive github event '%s': \n%s\n", eventType, string(payload))

		handler("github", eventType, string(payload))

		return err
	})

	log.Printf("Starting web hook server on :%v\n", port)

	ec.Run(":" + strconv.Itoa(port))
}
