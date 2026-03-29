package service

// GeekEduService 实现 gRPC 服务接口，处理注册、登录、课程管理、订单和视频播放等核心业务逻辑
//gRPC 服务实现，处理注册、登录、课程管理、订单和视频播放等核心业务逻辑

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"geekedu-project/common/jwt"
	"geekedu-project/common/oss"
	"geekedu-project/logic-server/db"
	"geekedu-project/logic-server/model"
	"geekedu-project/proto"
)

// GeekEduService 实现gRPC服务接口
type GeekEduService struct {
	proto.UnimplementedGeekEduServiceServer
}

// Register 用户注册
func (s *GeekEduService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.Response, error) {
	// 1. 检查参数
	if req.Username == "" || req.Password == "" {
		return &proto.Response{Code: 400, Msg: "用户名/密码不能为空"}, nil
	}

	// 2. 密码加密
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &proto.Response{Code: 500, Msg: "密码加密失败"}, nil
	}

	// 3. 写入数据库
	user := &model.User{
		Username: req.Username,
		Password: string(hashPassword),
		Role:     req.Role,
	}
	if err := db.DB.Create(user).Error; err != nil {

		fmt.Printf("[注册失败] 数据库插入错误：%v\n", err) // 这行必须加，会显示真实原因

		if strings.Contains(err.Error(), "Duplicate entry") { // MySQL唯一索引冲突的错误信息包含这个关键词
			return &proto.Response{Code: 400, Msg: "用户名已存在"}, nil
		}
		// 其他错误（表不存在、字段不匹配等）返回“数据库错误”，并打印具体错误（方便排查）
		fmt.Printf("数据库插入失败：%v\n", err) // 打印错误到日志，方便你看具体原因
		return &proto.Response{Code: 500, Msg: "数据库错误，注册失败"}, nil
	}
	return &proto.Response{Code: 200, Msg: "注册成功"}, nil
}

// Login 用户登录
func (s *GeekEduService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.Response, error) {
	// 1. 查询用户
	var user model.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &proto.Response{Code: 400, Msg: "用户名不存在"}, nil
		}
		return &proto.Response{Code: 500, Msg: "数据库错误"}, nil
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &proto.Response{Code: 400, Msg: "密码错误"}, nil
	}

	// 3. 生成JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return &proto.Response{Code: 500, Msg: "生成Token失败"}, nil
	}

	// 4. 返回Token
	loginRes := &proto.LoginResponse{Token: token}
	data, _ := json.Marshal(loginRes)
	return &proto.Response{Code: 200, Msg: "登录成功", Data: data}, nil
}

// GetCourseList 获取课程列表
func (s *GeekEduService) GetCourseList(ctx context.Context, req *proto.CourseListRequest) (*proto.Response, error) {
	// 1. 查询所有课程
	var courses []model.Course
	if err := db.DB.Find(&courses).Error; err != nil {
		emptyCourseList, _ := json.Marshal([]interface{}{}) // 空数组序列化
		return &proto.Response{
			Code: 500,
			Msg:  "查询课程失败",
			Data: emptyCourseList, // 失败时返回空数组，前端可安全访问 length
		}, nil
	}

	// 2. 为每个课程生成封面的预签名URL
	type CourseVO struct {
		ID        int64   `json:"id"`
		Title     string  `json:"title"`
		Price     float64 `json:"price"`
		Intro     string  `json:"intro"`
		CoverUrl  string  `json:"cover_url"` // 签名URL
		CreatedBy int64   `json:"created_by"`
	}
	var courseVOList []CourseVO
	for _, course := range courses {
		// 生成封面的签名URL
		coverUrl, err := oss.GeneratePresignedURL(course.CoverOssKey)
		if err != nil {
			coverUrl = ""
		}
		courseVOList = append(courseVOList, CourseVO{
			ID:        course.ID,
			Title:     course.Title,
			Price:     course.Price,
			Intro:     course.Intro,
			CoverUrl:  coverUrl,
			CreatedBy: course.CreatedBy,
		})
	}

	// 3. 返回数据
	data, _ := json.Marshal(courseVOList)
	return &proto.Response{Code: 200, Msg: "成功", Data: data}, nil
}

// CreateCourse 发布课程
func (s *GeekEduService) CreateCourse(ctx context.Context, req *proto.CreateCourseRequest) (*proto.Response, error) {
	// 1. 检查参数
	if req.Title == "" || req.Price == "" || req.CreatedBy == 0 {
		return &proto.Response{Code: 400, Msg: "课程标题/价格/创建者不能为空"}, nil
	}

	// 2. 解析价格
	price, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		return &proto.Response{Code: 400, Msg: "价格格式错误"}, nil
	}

	// 3. 上传封面到OSS
	// 生成唯一的OSS存储路径
	objectKey := fmt.Sprintf("covers/%d_%d.jpg", time.Now().Unix(), rand.Intn(10000))
	err = oss.UploadFile(objectKey, req.CoverFile)
	if err != nil {
		return &proto.Response{Code: 500, Msg: "封面上传失败"}, nil
	}

	// 4. 保存课程到数据库
	course := &model.Course{
		Title:       req.Title,
		Price:       price,
		Intro:       req.Intro,
		CoverOssKey: objectKey,
		CreatedBy:   req.CreatedBy,
	}
	if err := db.DB.Create(course).Error; err != nil {
		return &proto.Response{Code: 500, Msg: "创建课程失败"}, nil
	}

	// 5. 返回课程ID
	type ResData struct {
		CourseID int64 `json:"course_id"`
	}
	data, _ := json.Marshal(ResData{CourseID: course.ID})
	return &proto.Response{Code: 200, Msg: "发布成功", Data: data}, nil
}

// CreateOrder 购买课程
func (s *GeekEduService) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.Response, error) {
	// 1. 检查参数
	if req.UserId == 0 || req.CourseId == 0 {
		return &proto.Response{Code: 400, Msg: "用户ID/课程ID不能为空"}, nil
	}

	// 2. 检查是否已购买
	var order model.Order
	if err := db.DB.Where("user_id = ? AND course_id = ?", req.UserId, req.CourseId).First(&order).Error; err == nil {
		return &proto.Response{Code: 400, Msg: "已购买该课程"}, nil
	}

	// 3. 创建订单（模拟支付，直接设为paid）
	order = model.Order{
		UserID:   req.UserId,
		CourseID: req.CourseId,
		Status:   "paid",
	}
	if err := db.DB.Create(&order).Error; err != nil {
		return &proto.Response{Code: 500, Msg: "创建订单失败"}, nil
	}

	return &proto.Response{Code: 200, Msg: "购买成功"}, nil
}

// GetPlayUrl 获取视频播放链接（核心考核点）
func (s *GeekEduService) GetPlayUrl(ctx context.Context, req *proto.GetPlayUrlRequest) (*proto.Response, error) {
	// 1. 检查参数
	if req.UserId == 0 || req.VideoId == 0 {
		return &proto.Response{Code: 400, Msg: "用户ID/视频ID不能为空"}, nil
	}

	// 2. 查询视频关联的课程ID
	var video model.Video
	if err := db.DB.Where("id = ?", req.VideoId).First(&video).Error; err != nil {
		return &proto.Response{Code: 400, Msg: "视频不存在"}, nil
	}

	// 3. 检查用户是否购买该课程
	var order model.Order
	if err := db.DB.Where("user_id = ? AND course_id = ? AND status = ?", req.UserId, video.CourseID, "paid").First(&order).Error; err != nil {
		return &proto.Response{Code: 403, Msg: "未购买该课程，无法播放"}, nil
	}

	// 4. 生成OSS预签名URL（有效期1小时）
	playUrl, err := oss.GeneratePresignedURL(video.OssKey)
	if err != nil {
		return &proto.Response{Code: 500, Msg: "生成播放链接失败"}, nil
	}

	// 5. 返回播放链接
	type ResData struct {
		PlayUrl string `json:"play_url"`
	}
	data, _ := json.Marshal(ResData{PlayUrl: playUrl})
	return &proto.Response{Code: 200, Msg: "成功", Data: data}, nil
}

// UploadVideo 上传视频到OSS（流式上传）
func (s *GeekEduService) UploadVideo(stream grpc.ClientStreamingServer[proto.CreateCourseRequest, proto.Response]) error {
	// 接收流式上传的分片，合并后上传到 OSS，并在 videos 表中创建记录，返回 video_id。

	var courseID int64
	var filename string
	buf := bytes.NewBuffer(nil)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 记录 courseID/filename（以第一条消息为准）
		if courseID == 0 {
			courseID = req.CreatedBy
			filename = req.Title
		}

		if len(req.CoverFile) > 0 {
			if _, err := buf.Write(req.CoverFile); err != nil {
				return err
			}
		}
	}

	// 检查必须字段
	if courseID == 0 || filename == "" {
		return fmt.Errorf("缺少课程ID或文件名")
	}

	// 生成 OSS 存储路径并上传
	objectKey := fmt.Sprintf("videos/%d_%d_%s", courseID, time.Now().Unix(), filename)
	if err := oss.UploadFile(objectKey, buf.Bytes()); err != nil {
		return err
	}

	// 在数据库创建视频记录
	video := &model.Video{
		CourseID:  courseID,
		OssKey:    objectKey,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(video).Error; err != nil {
		return err
	}

	// 返回结果给客户端
	type ResData struct {
		VideoID int64 `json:"video_id"`
	}
	data, _ := json.Marshal(ResData{VideoID: video.ID})
	return stream.SendAndClose(&proto.Response{Code: 200, Msg: "上传成功", Data: data})
}
