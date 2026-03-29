// 前端的 “学习心得展示页面”，对应 Go 后端的/api/study-note接口 —— 从后端获取 Markdown 文本，把纯文本转换成带格式的网页内容（比如标题、列表、代码块），展示给用户

//导入React的状态/副作用钩子（管理数据和页面加载逻辑）
import { useState, useEffect } from 'react'
//导入ReactMarkdown组件（把Markdown文本渲染成HTML）
import ReactMarkdown from 'react-markdown'
// 导入axios（发送HTTP请求，调用后端接口）
import axios from 'axios'

function StudyNote() {
  const [content, setContent] = useState('')
//   content：存储学习心得的 Markdown 文本（初始为空）；
// setContent：更新content的函数 —— 当后端返回数据后，调用这个函数更新content，页面会自动重新渲染（显示新内容）；

  //学习心得,页面加载时执行：调用后端接口获取学习心得
  useEffect(() => {  //useEffect的作用：在 “页面挂载（第一次显示）” 时执行括号里的代码；
    //发送GET请求到后端/api/study-note接口（走Vite代理到8080端口）
    axios.get('/api/study-note').then(res => {
      setContent(res.data.content) //res.data：后端返回的 JSON 数据
    })
  }, [])// 空依赖数组：只在页面第一次加载时执行一次，确保只执行一次

  //渲染页面结构
  return (
    <div style={{ padding: 20, background: '#fff' }}>
      <h2>学习心得</h2>
      <ReactMarkdown>{content}</ReactMarkdown>{/* 把 Markdown 文本转换成美观的 HTML */}
    </div>
  )
  //如果调用失败，可以用.catch捕捉错误，console.error打印错误，或者显示一个错误提示
}

export default StudyNote