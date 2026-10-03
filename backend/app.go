package main

import (
	"flag"
	"fmt"
	"os"

	adminui "handfree-work/octo-backup/internal/admin"
	"handfree-work/octo-backup/internal/base/db_"
	"handfree-work/octo-backup/internal/base/log_"
	"handfree-work/octo-backup/internal/base/web_"
	"handfree-work/octo-backup/internal/config"
	"handfree-work/octo-backup/internal/handler"
	"handfree-work/octo-backup/internal/models"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins"
	"handfree-work/octo-backup/internal/svc"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"go.uber.org/zap"
)

// @title OctoBackup API
// @version 1.0
// @description OctoBackup 的认证与用户管理 API。
// @BasePath /
// @schemes http
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
// @description 输入 JWT，格式为 "Bearer {token}"。

var (
	port      = flag.String("port", "", "覆盖配置中的监听地址")
	mode      = flag.String("mode", getenv("APP_MODE", "dev"), "运行模式")
	configDir = flag.String("config", "./etc", "配置文件目录")
	prod      = flag.Bool("prod", false, "启用 Fiber prefork")
	adminMode = flag.Bool("admin", false, "进入终端管理界面")
)

func main() {
	flag.Parse()
	cfg, err := config.Load(*configDir, *mode)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		return
	}
	if _, err := log_.InitZap(log_.ZapConfig{
		Mode:      cfg.Mode,
		Directory: cfg.Log.Directory,
		Level:     cfg.Log.Level,
	}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		return
	}
	defer log_.ToDefer()
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
	jwtSecret := getenv("JWT_SECRET", cfg.Auth.JWTSecret)
	if cfg.Mode == "prod" && os.Getenv("JWT_SECRET") == "" {
		log_.Logger.Error("生产环境必须设置 JWT_SECRET")
		return
	}
	authConfig, err := web_.NewConfig(jwtSecret, cfg.Auth.TokenTTL)
	if err != nil {
		log_.Logger.Error("初始化认证配置失败", zap.Error(err))
		return
	}
	database, err := db_.OpenSQLite(databasePath)
	if err != nil {
		log_.Logger.Error("连接数据库失败", zap.Error(err))
		return
	}
	if err := db_.Migrate(database, &models.User{}, &models.SysSetting{}, &models.Plugin{}); err != nil {
		log_.Logger.Error("初始化数据库失败", zap.Error(err))
		return
	}
	if *adminMode {
		if err := adminui.Run(database, os.Stdin, os.Stdout); err != nil {
			log_.Logger.Error("管理界面退出", zap.Error(err))
		}
		return
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(log_.HTTPMiddleware())
	pluginRegistry := plugin.NewRegistry()
	if err := plugins.RegisterAll(pluginRegistry); err != nil {
		log_.Logger.Error("注册插件失败", zap.Error(err))
		return
	}
	handler.Register(app, &svc.ServiceContext{Db: database, Auth: authConfig, Plugins: pluginRegistry})
	app.Get("/*", static.New("./static/public"))

	log_.Logger.Info("OctoBackup 启动", zap.String("mode", cfg.Mode), zap.String("address", listenAddr))
	if err := app.Listen(listenAddr, fiber.ListenConfig{EnablePrefork: *prod}); err != nil {
		log_.Logger.Error("服务停止", zap.Error(err))
	}
}

func getenv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
