import axios from 'axios';
import message from 'antd/es/message';
// 这个文件是前端项目的 “请求工具”，封装了 axios 实例，配置了基础路径、请求/响应拦截器，统一处理 Token 注入和错误提示，供整个前端项目调用后端接口时使用。

// 创建 axios 实例（基础路径需与后端 Nginx 配置一致）
const request = axios.create({
  baseURL: '/api/v1', // 后端接口统一前缀（如后端配置为 /api/v1，需对应）
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 1. 请求拦截器：自动注入 Token（关键！）
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      // Token 格式需与后端 JWT 解析逻辑一致（通常为 Bearer + 空格 + Token）
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    message.error('请求发送失败，请检查网络');
    return Promise.reject(error);
  }
);

// 2. 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => {
    // 后端返回格式统一为 { code, msg, data }，直接返回响应体
    return response.data;
  },
  (error) => {
    // Token 过期（401）：清空状态并提示
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.reload(); // 刷新页面触发重新登录
      message.error('登录已过期，请重新登录');
    } else {
      // 其他错误：显示后端返回的错误信息
      const errMsg = error.response?.data?.msg || '接口请求失败';
      message.error(errMsg);
    }
    return Promise.reject(error);
  }
);

export default request;