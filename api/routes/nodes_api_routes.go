package routes

import (
	"loadbalancer/data/models"
	"loadbalancer/data/repositories"

	"github.com/gofiber/fiber/v2"
)

func GetAllNodes(c *fiber.Ctx) error {
	allNodes, err := repositories.GetAllNodes()

	if err != nil {
		return c.Status(500).SendString("Could not query the nodes table in the sqlite database")
	}

	return c.JSON(allNodes)
}

func CreateNode(c *fiber.Ctx) error {
	node := models.Node{}

	if err := c.BodyParser(&node); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	err := repositories.CreateNode(node.Address)

	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.Status(201).SendString("Created node")
}

func DeleteNode(c *fiber.Ctx) error {
	node := models.Node{}

	if err := c.BodyParser(&node); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	err := repositories.DeleteNodeByAddress(node.Address)

	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.Status(200).SendString("Deleted node")
}

func SetupNodesApiRoutes(app *fiber.App) {
	app.Get("/nodes", GetAllNodes)
	app.Post("/nodes", CreateNode)
	app.Delete("/nodes", DeleteNode)
}
