## 一、基础信息
- 接口前缀：所有接口均以 `/api` 开头
- 后端地址：http://localhost:8080
- 响应格式：统一返回 JSON，格式为 `{ "msg": "提示信息", "data": {} }`（成功）或 `{ "msg": "错误信息" }`（失败）

## 二、具体接口列表
### 1. 获取学习心得
- 请求方式：GET
- 接口地址：/api/study-note
- 请求参数：无
- 返回示例（成功，200 OK）：
{
  "msg": "success",
  "content": "# 我的学习心得..."
}

### 2. 查询题目列表（分页/筛选/搜索）
- 请求方式：GET
- 接口地址：/api/questions
- 返回示例（成功，200 OK）：
{
  "msg": "success",
  "list": [
    {
      "id": 1,
      "type": "单选题",
      "content": "Go语言中声明变量的关键字是？",
      "options": "A.const,B.var,C.let,D.func",
      "answer": "B",
      "difficulty": "简单",
      "language": "Go"
    }
  ],
  "total": 10,
  "page": 1,
  "size": 10
}

### 3. 添加题目
- 请求方式：POST
- 接口地址：/api/questions
- 请求体（JSON）：
{
  "type": "单选题", 
  "content": "题目内容",
  "options": "A.选项1,B.选项2",
  "answer": "A",
  "difficulty": "中等",
  "language": "Go"
}
- 返回示例（成功，200 OK）：{ "msg": "添加成功" }

### 4. 编辑题目
- 请求方式：PUT
- 接口地址：/api/questions
- 请求体（JSON）：同上（必须包含 id 字段）
- 返回示例（成功，200 OK）：{ "msg": "编辑成功" }

### 5. 删除题目
- 请求方式：DELETE
- 接口地址：/api/questions
- 请求体（JSON）：{ "ids": [1,2,3] } // 题目 ID 数组
- 返回示例（成功，200 OK）：{ "msg": "删除成功" }

### 6. AI 生成题目
- 请求方式：POST
- 接口地址：/api/ai-generate
- 请求体（JSON）：
{
  "type": "单选题",
  "count": 3,
  "difficulty": "简单",
  "language": "JavaScript"
}
- 返回示例（成功，200 OK）：
{
  "msg": "成功生成3道题目",
  "list": [
    {
      "type": "单选题",
      "content": "题目内容",
      "options": "A.选项1,B.选项2",
      "answer": "A",
      "difficulty": "简单",
      "language": "JavaScript"
    }
  ]
}