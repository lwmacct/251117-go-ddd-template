package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-bd-vmalert/internal/adapters/http"
	"github.com/lwmacct/251117-bd-vmalert/internal/infrastructure/config"
)

// Container 依赖注入容器
type Container struct {
	Config *config.Config
	Router *gin.Engine
}

func NewContainer(cfg *config.Config) (*Container, error) {

	// 7. 初始化路由
	router := http.SetupRouter(cfg)

	return &Container{
		Config: cfg,
		Router: router,
	}, nil

}
