package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 1. 数据库模型（题目）
type Question struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Type       string `json:"type"`       // 单选题/多选题/编程题
	Content    string `json:"content"`    // 题目内容
	Options    string `json:"options"`    // 选项
	Answer     string `json:"answer"`     // 答案
	Difficulty string `json:"difficulty"` // 简单/中等/困难
	Language   string `json:"language"`   // 编程语言
}

var db *gorm.DB // 全局数据库连接对象，所有数据库操作都用它
//全局变量简化了代码，避免在函数间传递 *gorm.DB 对象。

// 初始化数据库
func initDB() {
	var err error
	// 连接SQLite数据库，文件名为question.db
	db, err = gorm.Open(sqlite.Open("question.db"), &gorm.Config{})
	//&gorm.Config{} 是 GORM 的配置对象，可以设置日志级别、命名策略、连接池等。如果不传，则使用默认配置。对于简单项目，默认配置足够，但生产环境可能需要自定义（如慢日志、跳过默认事务等）。
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	// 自动创建Question对应的数据库表（表名默认是questions）
	// AutoMigrate 会根据结构体自动创建表、添加缺失的字段、修改字段类型（但不会删除字段）。不安全用于生产环境，因为可能意外删除数据或引起数据丢失。
	db.AutoMigrate(&Question{})

	// 	gorm.Open：建立数据库连接
	// AutoMigrate：自动迁移表结构（如果表不存在则创建，字段变更则更新）
}

// 加载.env配置
func initEnv() {
	err := godotenv.Load() //读取根目录的.env文件
	if err != nil {
		log.Println("未找到.env文件，可能影响大模型调用") // 容错，不退出程序
	}
}

// 后端接口！所有的接口都挂载在/api前缀下，使用 Gin 的gin.Context处理请求和响应。

// 获取学习心得
func getStudyNote(c *gin.Context) {
	content, err := os.ReadFile("../学习心得.md") //读取文件内容（字节数组），转成字符串返回给前端
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "读取学习心得失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

// 查询题目列表（分页、筛选、搜索）
func listQuestions(c *gin.Context) {
	// 1. 获取前端传的查询参数（分页+筛选+搜索）
	page, _ := strconv.Atoi(c.Query("page")) // 页码（字符串转数字）
	size, _ := strconv.Atoi(c.Query("size")) // 每页条数
	questionType := c.Query("type")          // 题目类型筛选
	keyword := c.Query("keyword")            // 关键词搜索

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset := (page - 1) * size

	var questions []Question       // 存储查询结果的切片
	query := db.Model(&Question{}) // 基于Question模型构建查询
	//db.Model(&Question{}) 和 db.Table("questions") 有什么区别？这里为什么用 Model？
	// Model 指定了要操作的结构体，GORM 会根据结构体名推断表名。Table 则直接指定表名字符串。使用 Model 更符合 ORM 风格，并且可以利用结构体中的 gorm 标签（如主键、索引等）

	// 筛选条件
	//query.Where() 是 GORM 框架提供的条件查询方法，作用是给数据库查询语句添加 WHERE 子句（条件筛选），只返回符合条件的数据。
	if questionType != "" {
		query = query.Where("type = ?", questionType)
	}
	// 搜索关键词
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 分页查询
	var total int64
	query.Count(&total)
	// Count 必须在 Find 之前调用，因为 Count 会执行一个 SELECT COUNT(*) 查询，而 Find 会执行 SELECT *。如果先执行 Find，Count 会被添加到同一个查询链中，导致 SQL 语句错误。GORM 的设计是链式调用，但每个终止方法（Count、Find）都会触发执行。
	query.Offset(offset).Limit(size).Find(&questions) //Find：执行查询并把结果存入切片

	c.JSON(http.StatusOK, gin.H{
		"list":  questions,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// 添加题目（手工/AI出题）
func addQuestion(c *gin.Context) {
	var q Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		//ShouldBindJSON：自动把请求体的 JSON 数据绑定到结构体，参数错误返回 400。
		return
	}
	db.Create(&q) //db.Create：插入一条新记录到数据库。
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功"})
}

// 编辑题目
func editQuestion(c *gin.Context) {
	var q Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	// 根据ID更新数据（只更新有值的字段）
	db.Model(&Question{}).Where("id = ?", q.ID).Updates(q)
	//Updates：GORM 的更新方法，会自动忽略空值字段（比如前端只传了 content，就只更 content）
	c.JSON(http.StatusOK, gin.H{"msg": "编辑成功"})
}

// 删除题目（单个/批量）
func deleteQuestion(c *gin.Context) {
	// 定义接收参数的结构体（IDs是数组）
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	db.Delete(&Question{}, req.IDs)
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// AI 生成题目
func aiGenerateQuestion(c *gin.Context) {
	// 接收前端参数：类型、数量、难度、编程语言
	var params struct {
		Type       string `json:"type"`       // 题型：单选题/多选题/编程题
		Count      int    `json:"count"`      // 题目数量
		Difficulty string `json:"difficulty"` // 难度：简单/中等/困难
		Language   string `json:"language"`   // 编程语言（仅编程题）
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		//尝试把前端传过来的 JSON 数据，绑定到我们定义的 params 结构体里
		//ShouldBindJSON：Gin 提供的「参数绑定方法」，专门处理 POST/PUT 等请求体中的 JSON 数据；
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	// 限制数量
	if params.Count < 1 || params.Count > 10 {
		params.Count = 1
	}

	// 构造清晰的提示词，强制按格式返回
	prompt := fmt.Sprintf(`你是专业的题库生成助手，请严格按照要求生成且仅生成%d道%s，难度为%s。`, params.Count, params.Type, params.Difficulty)
	if params.Type == "编程题" {
		if params.Language != "" {
			prompt += fmt.Sprintf(` 编程语言为%s。`, params.Language)
		}
	} else {
		prompt += ` 要求包含选项和正确答案。`
	}
	prompt += `
### 输出格式要求（必须严格遵守，否则无效）：
1. 每个题目之间用---分隔；
2. 每个题目包含【题目内容】【选项】【答案】三个字段；
3. 选项格式：A.内容,B.内容,C.内容,D.内容；
4. 答案格式：单选题填单个字母（如A），多选题填多个字母（如AB）；
5. 编程题不需要【选项】和【答案】，仅保留【题目内容】；
6. 只返回题目列表，不要任何额外解释、说明文字。

### 示例：
【题目内容】Go语言中，声明变量的关键字是？
【选项】A.const,B.var,C.let,D.func
【答案】B
---
【题目内容】JavaScript中，以下哪些是基本数据类型？
【选项】A.Number,B.Object,C.String,D.Array
【答案】AC
---
【题目内容】用Go语言实现一个计算两数之和的函数
【选项】
【答案】
`

	// 调用通义千问API
	client := resty.New()
	//resty 提供了链式 API、自动 JSON 编解码、超时设置、重试机制等，简化了 HTTP 客户端的使用
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+os.Getenv("BAILIAN_API_KEY")).
		SetBody(map[string]interface{}{
			"model": "qwen-turbo", // 通义千问免费模型
			"messages": []map[string]string{
				{"role": "user", "content": prompt}, // 用户指令
			},
			"temperature": 0.7,
			"max_tokens":  2000,
		}).
		Post(os.Getenv("BAILIAN_API_URL"))

	// 检查请求是否成功
	if err != nil {
		log.Println("API请求失败：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "大模型接口请求失败"})
		return
	}
	if resp.StatusCode() != http.StatusOK {
		log.Println("API返回状态码错误：", resp.StatusCode(), "响应内容：", string(resp.Body()))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "大模型返回异常"})
		return
	}

	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"` // 大模型返回的题目内容
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(resp.Body(), &aiResp); err != nil {
		log.Println("解析响应失败：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "解析题目失败"})
		return
	}
	// 检查API是否返回错误
	if aiResp.Error.Message != "" {
		log.Println("大模型返回错误：", aiResp.Error.Message)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": aiResp.Error.Message})
		return
	}
	aiContent := aiResp.Choices[0].Message.Content
	log.Println("大模型返回内容：", aiContent) // 打印日志

	questionStrs := strings.Split(aiContent, "---")
	var previewQuestions []Question
	for _, qStr := range questionStrs {
		qStr = strings.TrimSpace(qStr)
		if qStr == "" {
			continue
		}

		// 提取字段
		content := ""
		if strings.Contains(qStr, "【题目内容】") {
			contentPart := strings.Split(qStr, "【题目内容】")[1]
			if strings.Contains(contentPart, "【选项】") {
				content = strings.Split(contentPart, "【选项】")[0]
			} else {
				content = contentPart
			}
		}

		options := ""
		if strings.Contains(qStr, "【选项】") {
			optionsPart := strings.Split(qStr, "【选项】")[1]
			if strings.Contains(optionsPart, "【答案】") {
				options = strings.Split(optionsPart, "【答案】")[0]
			}
		}

		answer := ""
		if strings.Contains(qStr, "【答案】") {
			answer = strings.Split(qStr, "【答案】")[1]
		}

		// 清理空格/换行
		content = strings.TrimSpace(content)
		options = strings.TrimSpace(options)
		answer = strings.TrimSpace(answer)

		// 列表
		previewQuestions = append(previewQuestions, Question{
			Type:       params.Type,
			Content:    content,
			Options:    options,
			Answer:     answer,
			Difficulty: params.Difficulty,
			Language:   params.Language,
		})
	}
	if len(previewQuestions) > params.Count {
		previewQuestions = previewQuestions[:params.Count]
		// 只取前params.Count道，确定题目数量
	}

	// 返回真实生成的题目
	c.JSON(http.StatusOK, gin.H{
		"list": previewQuestions,
		"msg":  fmt.Sprintf("成功生成%d道题目", len(previewQuestions)),
	})
}

func serveFrontend(r *gin.Engine) {
	// 1. 把前端静态文件托管到 /static 路径,防止和/api冲突
	r.Static("/static", "../client/dist")

	r.GET("/", func(c *gin.Context) {
		c.File("../client/dist/index.html")
	})
	// 处理单页应用路由（刷新页面不404）
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.File("../client/dist/index.html")
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	//作用：让 Go 后端同时托管前端静态文件（Vue/React 打包后的 dist 目录），实现前后端一体化部署。
	//NoRoute：处理前端路由（比如 Vue 的/questions），刷新页面时返回 index.html，避免 404。

	//在生产环境中，通常用 Nginx 来托管静态文件，为什么这里选择用 Go 来托管？有什么好处和坏处？
	//用 Go 托管静态文件的好处是部署简单，只需一个二进制文件即可运行，无需额外安装 Nginx。坏处是静态文件处理性能不如 Nginx，且会占用应用服务器的资源。
}

func main() {
	// 初始化数据库和环境变量
	initDB()
	initEnv()
	r := gin.Default()

	// 后端接口
	// 定义API路由组（所有接口都以/api开头）
	//Group 用于创建路由分组，可以将共同前缀的接口组织在一起，便于添加中间件（如鉴权、日志）。
	api := r.Group("/api")
	{
		api.GET("/study-note", getStudyNote)
		api.GET("/questions", listQuestions)
		api.POST("/questions", addQuestion)
		api.PUT("/questions", editQuestion)
		//为什么使用 PUT 而不是 POST 来编辑题目？
		// PUT 通常用于更新整个资源，POST 用于创建新资源。这里 PUT 是合理的，因为编辑接口期望接收完整的题目对象。
		api.DELETE("/questions", deleteQuestion)
		api.POST("/ai-generate", aiGenerateQuestion)
	}
	// 托管前端
	serveFrontend(r)

	r.Run(":8080")
}
