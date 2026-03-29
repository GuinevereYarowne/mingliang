package controller

//定义控制器函数，处理用户注册、登录、课程列表、发布课程、购买课程、获取播放链接等请求，并调用 gRPC 客户端与 Logic Server 交互

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"geekedu-project/proto"
	"geekedu-project/web-server/client"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	// 1. 绑定参数
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 2. 调用gRPC服务
	res, err := client.GeekEduClient.Register(c, &proto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	// 3. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	// 1. 绑定参数
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 2. 调用gRPC服务
	res, err := client.GeekEduClient.Login(c, &proto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	// 3. 解析返回的Token
	var loginRes proto.LoginResponse
	if len(res.Data) > 0 {
		json.Unmarshal(res.Data, &loginRes)
	}

	// 4. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
		"data": gin.H{"token": loginRes.Token},
	})
}

// GetCourseList 获取课程列表（公开接口）
func GetCourseList(c *gin.Context) {
	// 调用gRPC服务
	res, err := client.GeekEduClient.GetCourseList(c, &proto.CourseListRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	// 解析课程列表
	var courses []interface{}
	if len(res.Data) > 0 {
		json.Unmarshal(res.Data, &courses)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
		"data": courses,
	})
}

// CreateCourse 发布课程（管理员）
func CreateCourse(c *gin.Context) {
	// 1. 获取用户ID（从JWT中间件）
	userID, _ := c.Get("user_id")

	// 2. 获取表单数据
	title := c.PostForm("title")
	price := c.PostForm("price")
	intro := c.PostForm("intro")

	// 3. 获取封面文件
	file, err := c.FormFile("cover_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请上传封面文件",
		})
		return
	}

	// 4. 读取文件内容
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "读取文件失败",
		})
		return
	}
	defer fileBytes.Close()
	content := make([]byte, file.Size)
	fileBytes.Read(content)

	// 5. 调用gRPC服务
	res, err := client.GeekEduClient.CreateCourse(c, &proto.CreateCourseRequest{
		Title:     title,
		Price:     price,
		Intro:     intro,
		CoverFile: content,
		CreatedBy: userID.(int64),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	// 解析课程ID
	var courseRes struct {
		CourseID int64 `json:"course_id"`
	}
	if len(res.Data) > 0 {
		json.Unmarshal(res.Data, &courseRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
		"data": gin.H{"course_id": courseRes.CourseID},
	})
}

// CreateOrder 购买课程
func CreateOrder(c *gin.Context) {
	// 1. 获取用户ID（从JWT中间件）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录或Token无效",
		})
		return
	}

	// 2. 绑定参数
	var req struct {
		CourseID int64 `json:"course_id"`
		UserID   int64 `json:"user_id"` // 接收前端的user_id，但优先使用JWT中的
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 3. 调用gRPC服务
	res, err := client.GeekEduClient.CreateOrder(c, &proto.CreateOrderRequest{
		UserId:   userID.(int64),
		CourseId: req.CourseID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
	})
}

// GetPlayUrl 获取播放链接
func GetPlayUrl(c *gin.Context) {
	// 1. 获取用户ID
	userID, _ := c.Get("user_id")

	// 2. 获取视频ID
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "视频ID错误",
		})
		return
	}

	// 3. 调用gRPC服务
	res, err := client.GeekEduClient.GetPlayUrl(c, &proto.GetPlayUrlRequest{
		UserId:  userID.(int64),
		VideoId: videoID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	// 解析播放链接
	var playRes struct {
		PlayUrl string `json:"play_url"`
	}
	if len(res.Data) > 0 {
		json.Unmarshal(res.Data, &playRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
		"data": gin.H{"play_url": playRes.PlayUrl},
	})
}

// GetUserInfo 获取当前用户信息（解决前端登录后获取user_id问题）
func GetUserInfo(c *gin.Context) {
	// 从JWT中间件获取用户信息
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"id":       userID,
			"username": username,
			"role":     role,
		},
	})
}

// UploadVideo 上传视频（新增：实现核心上传视频功能）
func UploadVideo(c *gin.Context) {
	// 1. 获取课程ID
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "课程ID错误",
		})
		return
	}

	// 2. 获取视频文件
	file, err := c.FormFile("video_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请上传视频文件",
		})
		return
	}

	// 3. 打开文件准备分片读取
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "读取文件失败",
		})
		return
	}
	defer fileBytes.Close()

	// 4. 调用gRPC流式服务
	stream, err := client.GeekEduClient.UploadVideo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "连接gRPC失败",
		})
		return
	}

	// 5. 发送视频数据：按 1MB 分片发送，避免 gRPC 单消息大小限制
	const chunkSize = 1 << 20 // 1MB
	buf := make([]byte, chunkSize)
	first := true
	total := 0
	fmt.Printf("开始通过 gRPC 流发送视频数据，文件名=%s\n", file.Filename)
	for {
		n, readErr := fileBytes.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			req := &proto.CreateCourseRequest{
				Title:     file.Filename,
				CreatedBy: courseID,
				CoverFile: chunk,
			}
			if !first {
				// 之后的分片可不重复传文件名/课程ID（服务端使用第一片的信息）
				req.Title = ""
				req.CreatedBy = 0
			}
			if err := stream.Send(req); err != nil {
				fmt.Printf("stream.Send 错误: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "发送视频数据失败",
				})
				return
			}
			total += n
			first = false
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			fmt.Printf("读取文件分片错误: %v\n", readErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "读取文件失败",
			})
			return
		}
	}
	fmt.Printf("总共发送字节: %d\n", total)

	// 6. 关闭流并接收响应
	res, err := stream.CloseAndRecv()
	if err != nil {
		// 打印详细错误到日志，便于调试
		fmt.Printf("stream.CloseAndRecv 错误: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务调用失败",
		})
		return
	}

	// 7. 解析视频ID
	var videoRes struct {
		VideoID int64 `json:"video_id"`
	}
	if len(res.Data) > 0 {
		json.Unmarshal(res.Data, &videoRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
		"data": gin.H{"video_id": videoRes.VideoID},
	})
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
