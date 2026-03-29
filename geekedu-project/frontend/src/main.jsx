// main.jsx 是前端项目的入口文件，负责渲染整个 React 应用。它引入了 Ant Design 的 ConfigProvider 组件来设置全局的中文语言环境，并将 App 组件作为根组件渲染到页面上。这个文件是前端项目的启动点，确保整个应用能够正确加载和显示。
import React from 'react';
import ReactDOM from 'react-dom/client';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')).render(
  <ConfigProvider locale={zhCN}>
    <App />
  </ConfigProvider>
);