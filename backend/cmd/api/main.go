package main

import (
	"fmt"
	"log"

	"github.com/InnoServe/blockSBOM/internal/api"
	"github.com/InnoServe/blockSBOM/internal/config"
	"github.com/InnoServe/blockSBOM/internal/dal"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := dal.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建 Hertz 服务器
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", cfg.Server.Port)),
	)

	// 注册路由
	api.RegisterRoutes(h)

	// 启动服务器
	h.Spin()
}
