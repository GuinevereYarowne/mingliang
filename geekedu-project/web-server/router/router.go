package router

// InitRouter 初始化路由，定义公开接口和需要鉴权的接口，并使用中间件保护敏感接口
//路由定义，使用 Gin 框架，设置跨域配置，定义公开接口和需要鉴权的接口，并使用中间件保护敏感接口

import (
	"geekedu-project/web-server/controller"
	"geekedu-project/web-server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 跨域配置（新手必加：前端访问后端需要）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 公开接口
	public := r.Group("/api/v1")
	{
		// 认证
		public.POST("/auth/register", controller.Register)
		public.POST("/auth/login", controller.Login)
		// 课程列表
		public.GET("/courses", controller.GetCourseList)
		// 健康检查
		public.GET("/healthz", controller.HealthCheck)
	}

	// 需要鉴权的接口
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth())
	{
		// 用户信息
		auth.GET("/auth/user-info", controller.GetUserInfo)

		// 管理员接口
		admin := auth.Group("/courses")
		admin.Use(middleware.AdminAuth())
		{
			admin.POST("", controller.CreateCourse)           // 发布课程
			admin.POST("/:id/videos", controller.UploadVideo) // 上传视频
		}

		// 普通用户接口
		auth.POST("/orders", controller.CreateOrder)         // 购买课程
		auth.GET("/player/:video_id", controller.GetPlayUrl) // 获取播放链接
	}

	return r
}
