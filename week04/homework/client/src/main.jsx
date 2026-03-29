//整个前端项目的 “总开关”，负责初始化 React 应用、配置核心依赖（路由 / UI 组件库），并把页面挂载到浏览器中。
//像 Go 后端的main.go,所有前端代码的执行都从这里开始，同时配置了项目的核心基础能力（路由、UI 组件库）。
import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import App from './App'
//导入Ant Design的全局配置组件（统一配置UI组件的风格/语言等）
import { ConfigProvider } from 'antd'
//导入Ant Design的基础样式（必须，否则UI组件会没有样式）
import 'antd/dist/reset.css' // Ant Design样式

//核心渲染逻辑：把React组件挂载到浏览器页面中
ReactDOM.createRoot(document.getElementById('root')).render(
  //Ant Design全局配置（包裹整个应用，可统一设置主题/语言等）
  //剩下的都是路由器相关，用来渲染App组件的
  <ConfigProvider>
    <BrowserRouter>
      <Routes>
        <Route path="*" element={<App />} />
      </Routes>
    </BrowserRouter>
  </ConfigProvider>
)

//<BrowserRouter> 和 <Routes> 的作用是什么？为什么使用 path="*"？
//BrowserRouter 使用 HTML5 history API 管理路由。Routes 是 React Router v6 的新组件，用于声明路由规则，path="*" 表示匹配所有路径，将所有路由交给 App 组件内部处理。这样做是为了在 App 组件中再定义子路由，实现嵌套路由。