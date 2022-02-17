package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	//Ordem: 1 - containerdPath 2 - loggingPath 3 - stateManager 4 - registry
	argsWithoutProg := os.Args[1:]
	SetContainerdPath(argsWithoutProg[0])
	SetLoggingPath(argsWithoutProg[1])
	SetStateManagerUrl(argsWithoutProg[2])
	SetRegistry(argsWithoutProg[3])
	app := fiber.New()

	// GET /john
	app.Get("/:service/:snapshot-id", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	app.Post("/snapshot", func(c *fiber.Ctx) error {
		var body Snapshot
		err := c.BodyParser(&body)

		// if error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Cannot parse JSON",
			})
		}

		GenerateSnapshot(body)

		return c.SendString("ok")
	})

	app.Post("/restore", func(c *fiber.Ctx) error {
		var body Snapshot
		err := c.BodyParser(&body)

		// if error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Cannot parse JSON",
			})
		}

		Restore(body)

		return c.SendString("Ok")
	})

	log.Fatal(app.Listen(":6666"))
}
