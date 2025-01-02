/*

MAIN SOURCE FILE

*/

package main // package declaration

import (
	"fmt" // println package
	"log" // error handling
	"os"

	"github.com/gofiber/fiber/v2" // fiber
	"github.com/joho/godotenv"
)

// Todo represents a task with specific attributes.
type Todo struct {
	ID        int    `json:"id"`        // Unique identifier for the task (integer).
	Completed bool   `json:"completed"` // Status of the task: true if completed, false otherwise (boolean).
	Body      string `json:"body"`      // Description or content of the task (string).
}

func main() {

	// create a new application / server
	app := fiber.New()


	// loads in environment variables
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading .env file")
	}
	
	PORT := os.Getenv("PORT") // specifies what port the api is runing on


	// 1. API WITHOUT DB / IN MEMORY

	// array of todos for the list
	todos := []Todo{}

	// first main route
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos) // basically just sends all get requests this json with hello world
	})

	// CREATE A TODO ENDPOINT
	// post endpoint that allows you to add todos via the api
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // {id: 0, completed: false, body: "task #1"}
		// this should be a pointer ebcause we are going to be altering it

		err := c.BodyParser(todo) // takes the json that the user sends form the request and parses into struct
		// check for errors (ie invalid json etc)
		if err != nil {
			// ERROR THROWN
			return err
		}

		// shouldn't have an empty task
		if todo.Body == "" {
			// throw an error
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		// SUCCESSFUL --> PASSED ALL CHECKS
		todo.ID = len(todos) + 1
		todos = append(todos, *todo) // add todo to todos list --> derefence b/c initially memory address

		return c.Status(201).JSON(todo) // returns success status and json of todo

	})

	// UPDATE A TODO
	// could either be put or patch endpoint
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") // extracts the id from the endpoint url

		for i, todo := range todos { // iterate through all todos
			// fmt sprint = converts integer todo.ID to string and compares with
			// id parameter
			if fmt.Sprint(todo.ID) == id {
				// updates if it is completed or not
				todos[i].Completed = !todos[i].Completed
				// returns success status
				return c.Status(200).JSON(todos[i])
			}
		}

		// todo not found
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// DELETE A TODO
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				// basically will just make the new array of every element excluding the current one
				// the ... gets all values till the end --> variatic operator
				todos = append(todos[:i], todos[i+1:]...)
				// return success status
				return c.Status(200).JSON(fiber.Map{"success": "true"})
			}
		}

		// todo isn't found
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// list on port 4000 with thebackend
	log.Fatal(app.Listen(":" + PORT)) // log.Fatal will show us the errors

}
