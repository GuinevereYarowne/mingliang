// 之前写的 Go 后端跑在 http://localhost:8080，而 React 前端用 Vite 启动后跑在 http://localhost:5173—— 这两个地址的端口不同，浏览器会触发「同源策略限制」，直接拦截前端对后端的请求（比如前端调 /api/questions 会报错）
import { defineConfig } from 'vite'
//导入Vite的React插件（让Vite支持React语法）
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
  return {
    //配置Vite插件：启用React插件
    plugins: [react()],
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true
          //代理配置的作用是什么？为什么要设置 changeOrigin: true？
          //代理用于将前端的 /api 请求转发到后端服务器，解决开发时的跨域问题。changeOrigin: true 会修改请求头中的 Origin 字段为目标地址的域名，避免后端做 Host 校验时出现问题
        }
      },
      port: 5173
    },
    base: mode === 'production' ? '/static/' : '/'
    //base 配置在开发和生产环境下的不同值有什么意义？如果生产环境下不设置 base 会怎样？
    //base 决定了静态资源的引用路径。开发时使用根路径 / 即可，因为 Vite 的开发服务器会正确处理资源请求。生产环境下设置 /static/ 是为了匹配后端托管静态文件的路径。如果不设置 base，生产环境下资源引用路径会错误，导致页面无法正确加载 CSS/JS 文件。
  }
})