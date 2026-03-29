// App.jsx 是前端项目的主组件，负责整体布局、状态管理、用户交互和与后端 API 的通信。它使用 React Hooks 管理状态，Ant Design 组件构建 UI，React Player 播放视频，并通过封装的 request 工具与后端接口交互，实现用户登录、注册、课程展示、购买和视频播放等功能。
import { useState, useEffect } from 'react';
import { Button, Input, List, Card, Modal, Form, message, Typography, Space, Skeleton } from 'antd';
import ReactPlayer from 'react-player';
import request from './utils/request'; // 需配合下方 request.js 配置（关键！）

const { Title, Text } = Typography;
const { Item } = Form;

function App() {
  // 1. 状态管理：补充用户信息字段（解决购买课程缺 user_id 问题）
  const [isLoginModalVisible, setIsLoginModalVisible] = useState(false);
  const [isRegisterModalVisible, setIsRegisterModalVisible] = useState(false);
  const [courses, setCourses] = useState([]); // 初始空数组，避免 null 报错
  const [isCourseLoading, setIsCourseLoading] = useState(true); // 课程列表加载态
  const [currentUser, setCurrentUser] = useState(null); // 登录后存 { token, id, username, role }
  const [playUrl, setPlayUrl] = useState('');
  const [isPlayerVisible, setIsPlayerVisible] = useState(false);
  const [isPlayerLoading, setIsPlayerLoading] = useState(false); // 播放器加载态

  // 2. 表单实例：注册表单加角色验证
  const [loginForm] = Form.useForm();
  const [registerForm] = Form.useForm();

  // 3. 初始化：检查登录状态 + 加载课程（解决刷新后用户信息丢失）
  useEffect(() => {
    const initApp = async () => {
      try {
        // 先加载课程
        await fetchCourses();
      } catch (err) {
        console.error('课程加载失败:', err);
      }
      
      // 检查本地 Token，有则获取用户信息（补全 user_id 等字段）
      const token = localStorage.getItem('token');
      if (token) {
        try {
          await getUserInfo(token);
        } catch (err) {
          console.error('用户信息加载失败:', err);
          localStorage.removeItem('token');
        }
      }
    };
    initApp();
  }, []);

  // 4. 获取课程列表：res.data 兜底空数组（彻底解决 length 报错）
  const fetchCourses = async () => {
    setIsCourseLoading(true); // 显示加载骨架屏
    try {
      const res = await request.get('/courses');
      // 关键修复：确保 res.data 不为 null，且总是返回数组
      const courseData = res?.code === 200 && res.data ? res.data : [];
      setCourses(Array.isArray(courseData) ? courseData : []);
    } catch (err) {
      console.error('课程加载错误:', err);
      message.error('课程加载失败，请稍后重试');
      setCourses([]); // 错误时强制设空数组
    } finally {
      setIsCourseLoading(false); // 关闭加载态
    }
  };

  // 5. 获取用户信息（解决登录后缺 user_id 问题）
  const getUserInfo = async (token) => {
    try {
      const res = await request.get('/auth/user-info'); // 需后端配合实现该接口
      if (res.code === 200) {
        // 存储完整用户信息（含 id，用于购买课程）
        setCurrentUser({
          token,
          id: res.data.id, // 后端返回的用户 ID
          username: res.data.username, // 后端返回的用户名
          role: res.data.role // 后端返回的用户角色
        });
      }
    } catch (err) {
      // Token 过期/无效，清空本地存储
      localStorage.removeItem('token');
      setCurrentUser(null);
      message.warning('登录状态已失效，请重新登录');
    }
  };

  // 6. 登录：补充用户信息获取（解决仅存 token 问题）
  const handleLogin = async () => {
    try {
      const values = await loginForm.validateFields();
      const res = await request.post('/auth/login', values);
      
      if (res.code === 200) {
        const { token } = res.data;
        localStorage.setItem('token', token);
        
        // 登录成功后获取用户信息（补全 id）
        try {
          const userRes = await request.get('/auth/user-info');
          if (userRes.code === 200) {
            const userData = {
              token,
              id: userRes.data.id,
              username: userRes.data.username,
              role: userRes.data.role
            };
            setCurrentUser(userData);
            setIsLoginModalVisible(false);
            message.success(`欢迎回来，${userRes.data.username}!`);
          }
        } catch (err) {
          message.error('获取用户信息失败');
          localStorage.removeItem('token');
        }
      } else {
        message.error(res.msg || '登录失败');
      }
    } catch (err) {
      message.error(err?.message || '登录表单验证失败');
    }
  };

  // 7. 注册：角色字段加验证（避免无效角色值）
  const handleRegister = async () => {
    try {
      const values = await registerForm.validateFields();
      const res = await request.post('/auth/register', values);
      
      if (res.code === 200) {
        setIsRegisterModalVisible(false);
        message.success('注册成功！请登录');
        // 自动打开登录弹窗
        setIsLoginModalVisible(true);
      } else {
        message.error(res.msg || '注册失败');
      }
    } catch (err) {
      message.error(err?.message || '注册表单验证失败');
    }
  };

  // 8. 退出登录：清空状态 + 本地存储
  const handleLogout = () => {
    localStorage.removeItem('token');
    setCurrentUser(null);
    message.success('退出登录成功');
  };

  // 9. 购买课程：传递 user_id（解决后端缺用户标识问题）
  const handleBuyCourse = async (courseId) => {
    if (!currentUser) {
      message.warning('请先登录再购买课程');
      setIsLoginModalVisible(true);
      return;
    }

    try {
      const res = await request.post('/orders', {
        user_id: currentUser.id, // 关键：传递用户 ID
        course_id: courseId
      });
      
      if (res.code === 200) {
        message.success('课程购买成功！可前往我的课程查看');
      } else {
        message.error(res.msg || '购买失败');
      }
    } catch (err) {
      message.error('购买请求失败，请检查网络');
    }
  };

  // 10. 播放视频：修复 videoId 硬编码（用课程实际 video_id）
  const handlePlayVideo = async (videoId) => {
    if (!currentUser) {
      message.warning('请先登录再观看视频');
      setIsLoginModalVisible(true);
      return;
    }

    setIsPlayerLoading(true); // 播放器加载态
    try {
      const res = await request.get(`/player/${videoId}`);
      if (res.code === 200) {
        setPlayUrl(res.data.play_url);
        setIsPlayerVisible(true);
      } else {
        message.error(res.msg || '获取播放链接失败');
      }
    } catch (err) {
      message.error('播放链接请求失败，请检查网络');
    } finally {
      setIsPlayerLoading(false); // 关闭加载态
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '1400px', margin: '0 auto' }}>
      {/* 头部：用户状态切换 */}
      <Space style={{ marginBottom: '24px', alignItems: 'center' }}>
        <Title level={2} style={{ margin: 0 }}>在线视频学习平台</Title>
        <Space style={{ marginLeft: 'auto' }}>
          {currentUser ? (
            <Space>
              <Text>欢迎，{currentUser.username}</Text>
              <Button onClick={handleLogout} type="text">退出登录</Button>
            </Space>
          ) : (
            <>
              <Button onClick={() => setIsLoginModalVisible(true)}>登录</Button>
              <Button onClick={() => setIsRegisterModalVisible(true)} type="primary">注册</Button>
            </>
          )}
        </Space>
      </Space>

      {/* 课程列表：加空状态 + 加载骨架屏（提升体验） */}
      <List
        grid={{ gutter: 24, column: 3, xs: 1, sm: 2, md: 3, lg: 3 }} // 响应式布局
        dataSource={courses}
        loading={isCourseLoading}
        renderEmpty={() => (
          <div style={{ textAlign: 'center', padding: '40px 0' }}>
            <Text type="secondary">暂无课程数据，敬请期待～</Text>
          </div>
        )}
        renderItem={(course) => (
          <List.Item key={course.id} style={{ width: '100%' }}>
            <Card
              style={{ width: '100%' }}
              title={course.title}
              cover={
                <Skeleton loading={isCourseLoading} active>
                  <img 
                    alt={course.title} 
                    src={course.cover_url || 'https://via.placeholder.com/400x200?text=课程封面'} 
                    style={{ height: '200px', objectFit: 'cover', width: '100%' }} 
                  />
                </Skeleton>
              }
            >
              <Skeleton loading={isCourseLoading} active paragraph={{ rows: 2 }}>
                <div style={{ color: '#666', fontSize: '14px', lineHeight: '1.5', marginBottom: '12px', display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                  {course.intro || '暂无课程介绍'}
                </div>
              </Skeleton>
              
              {/* 底部操作栏：价格 + 按钮 */}
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingTop: '12px', borderTop: '1px solid #f0f0f0' }}>
                <span style={{ fontSize: '16px', fontWeight: 'bold', color: '#ff4d4f' }}>
                  ¥{(course.price || 0).toFixed(2)}
                </span>
                <div style={{ display: 'flex', gap: '8px' }}>
                  <Button 
                    type="primary" 
                    size="small"
                    onClick={() => handleBuyCourse(course.id)}
                  >
                    立即购买
                  </Button>
                  <Button 
                    type="default"
                    size="small"
                    onClick={() => handlePlayVideo(course.video_id)}
                    disabled={!course.video_id}
                  >
                    {course.video_id ? '观看视频' : '无视频'}
                  </Button>
                </div>
              </div>
            </Card>
          </List.Item>
        )}
      />

      {/* 登录弹窗：简化表单 + 错误提示 */}
      <Modal
        title="用户登录"
        open={isLoginModalVisible}
        onCancel={() => setIsLoginModalVisible(false)}
        onOk={handleLogin}
        destroyOnClose // 关闭时清空表单
      >
        <Form form={loginForm} layout="vertical" name="login_form">
          <Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }, { min: 3, message: '用户名至少 3 个字符' }]}
          >
            <Input placeholder="请输入用户名" />
          </Item>
          <Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '密码至少 6 个字符' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Item>
        </Form>
      </Modal>

      {/* 注册弹窗：角色字段加验证 */}
      <Modal
        title="用户注册"
        open={isRegisterModalVisible}
        onCancel={() => setIsRegisterModalVisible(false)}
        onOk={handleRegister}
        destroyOnClose // 关闭时清空表单
      >
        <Form form={registerForm} layout="vertical" name="register_form">
          <Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }, { min: 3, message: '用户名至少 3 个字符' }]}
          >
            <Input placeholder="请输入用户名" />
          </Item>
          <Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '密码至少 6 个字符' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Item>
          <Item
            name="role"
            label="角色"
            rules={[
              { required: true, message: '请选择角色' },
              { enum: ['student', 'admin'], message: '角色只能是 student 或 admin' } // 限制角色值
            ]}
          >
            <Input placeholder="请输入角色（student/admin）" defaultValue="student" />
          </Item>
        </Form>
      </Modal>

      {/* 播放器弹窗：加加载态 + 空状态 */}
      <Modal
        title="视频播放"
        open={isPlayerVisible}
        onCancel={() => {
          setIsPlayerVisible(false);
          setPlayUrl(''); // 关闭时清空播放链接
        }}
        width={800}
        footer={null}
        destroyOnClose
      >
        {isPlayerLoading ? (
          <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Skeleton active paragraph={{ rows: 1 }} />
          </div>
        ) : playUrl ? (
          <ReactPlayer 
            url={playUrl} 
            width="100%" 
            height="400px" 
            controls 
            fallback={<div>视频加载中...</div>}
          />
        ) : (
          <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Text type="secondary">暂无有效播放链接</Text>
          </div>
        )}
      </Modal>
    </div>
  );
}

export default App;