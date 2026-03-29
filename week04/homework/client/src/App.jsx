// 搭建一个后台管理系统风格的布局（左侧菜单 + 右侧内容），通过 React Router 实现 “点击左侧菜单，右侧切换不同页面（学习心得 / 题库管理）” 的核心交互，是前端项目的 “主框架”。

import { Layout, Menu, Button } from 'antd'
import { useState, useEffect } from 'react'
import { Link, Outlet, Routes, Route } from 'react-router-dom' 
import axios from 'axios'
//这都是搭建页面骨架、布局，页面加载等等，然后下面是自定义的从另外文件导入的
// 导入pages页面
import StudyNote from './pages/StudyNote'
import QuestionManage from './pages/QuestionManage'

//解构Layout组件：提取Sider（左侧边栏）、Content（内容区）
//Layout 组件本质：Ant Design 提供的后台布局模板组件，内置 Sider（左侧边栏）、Content（内容区）等子组件；
const { Sider, Content } = Layout

function App() {
  // 左侧菜单折叠状态
  //定义状态：控制左侧菜单的折叠/展开（初始值false=展开）
  const [collapsed, setCollapsed] = useState(false)

  return (  //渲染页面结构（JSX）
    // 根布局：Layout是Ant Design的布局组件，minHeight:100vh让布局占满整个视口
    <Layout style={{ minHeight: '100vh' }}>
      {/* 左侧面板 */}
       {/* 开启折叠功能 */}
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
        {/* 左侧边栏的空白占位（美观用） */}
        <div style={{ height: 32, margin: 16 }} />
        <Menu theme="dark" defaultSelectedKeys={['1']} mode="inline">
          <Menu.Item key="1">
            <Link to="/">学习心得</Link>{/* Link是React Router的跳转组件（无刷新跳转） */}
          </Menu.Item>
          <Menu.Item key="2">
            {/* 菜单项2：题库管理，点击跳转到/question-manage */}
            <Link to="/question-manage">题库管理</Link>
            {/* Link 组件和 <a> 标签有什么区别？为什么不用 <a>?Link 组件会阻止浏览器的默认跳转行为,<a> 标签会导致整页刷新 */}
          </Menu.Item>
        </Menu>
      </Sider>

      {/* 右侧内容：配置子路由 */}
      <Layout>
        {/* 内容区：设置margin=16px（和边框保持间距） */}
        <Content style={{ margin: '16px' }}>
          <Routes>
            {/* 学习心得 */}
            {/* 规则1：根路径/ → 渲染StudyNote组件（学习心得页面） main.jsx 中的 <Routes> 匹配了所有路径，并渲染 <App />*/}
            <Route path="/" element={<StudyNote />} />
            {/* 题库管理 */}
            {/* 规则2：/question-manage → 渲染QuestionManage组件（题库管理页面） */}
            <Route path="/question-manage" element={<QuestionManage />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  )
}
// 10. 导出App组件（供main.jsx导入渲染）
export default App