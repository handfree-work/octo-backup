package main

import (
	"flag"
	"log"
	"os"

	"handfree-work/web-restic/internal/base/db_"
	"handfree-work/web-restic/internal/config"
	"handfree-work/web-restic/internal/handler"
	"handfree-work/web-restic/internal/models"
	"handfree-work/web-restic/internal/svc"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

var (
	port      = flag.String("port", "", "覆盖配置中的监听地址")
	mode      = flag.String("mode", getenv("APP_MODE", "dev"), "运行模式")
	configDir = flag.String("config", "./etc", "配置文件目录")
	prod      = flag.Bool("prod", false, "启用 Fiber prefork")
)

func main() {
	flag.Parse()
	cfg, err := config.Load(*configDir, *mode)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	listenAddr := cfg.Server.Port
	if *port != "" {
		listenAddr = *port
	}
	if listenAddr == "" {
		listenAddr = ":3000"
	}
	databasePath := cfg.Database.Path
	if databasePath == "" {
		databasePath = "./data/db.sqlite"
	}
	database, err := db_.OpenSQLite(databasePath)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	if err := db_.Migrate(database, &models.User{}); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	handler.UserHandlersRegister(app, &svc.ServiceContext{Db: database})
	app.Get("/*", static.New("./static/public"))

	log.Printf("web-restic mode=%s listening on %s", cfg.Mode, listenAddr)
	log.Fatal(app.Listen(listenAddr, fiber.ListenConfig{EnablePrefork: *prod}))
}

func getenv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
