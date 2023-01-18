package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/websocket/v2"
	"github.com/liuquanhao/moyu/controller"
	"github.com/liuquanhao/moyu/middleware"
)

//go:embed dist
var frontend embed.FS

func main() {
	base := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	baseurl := os.Getenv("BASEURL")
	app := base.Group("/"+baseurl)
	app.Use(cors.New())

	app.Use("/ws/*", middleware.UpgradeOptions)

	app.Get("/sys_info", controller.GetSysInfo)
	app.Get("/sys_status", controller.GetSysStatus)
	app.Get("/ws/sys_status", websocket.New(controller.PushSysStatus))

	stripped, err := fs.Sub(frontend, "dist")
	if err != nil {
		log.Fatal(err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(stripped),
		Browse: true,
	}))

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	log.Println(host + ":" + port)
	log.Fatal(base.Listen(host + ":" + port))
}
