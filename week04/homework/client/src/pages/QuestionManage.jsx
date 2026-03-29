// 这份代码是你前端题库系统的核心业务页面，也是功能最完整、前后端交互最密集的页面，核心实现了题目的全流程 CRUD 操作（查询 / 添加 / 编辑 / 删除），还包含手工出题 / AI 出题双模式、AI 生成题目预览勾选、单个 / 批量删除、分页 / 筛选 / 搜索等高频后台功能，所有操作都直接对接你之前写的 Go 后端接口（/api/questions//api/ai-generate）。
// 页面基于 React 状态管理 + Antd 海量 UI 组件 + Axios 异步请求实现，是典型的后台管理系统业务页面写法，下面按「核心定位→状态定义→核心业务函数→表格列配置→UI 渲染模块」的逻辑逐部分拆解

// 功能层面：实现题库的查询（分页 / 筛选 / 搜索）、添加（手工 / AI）、编辑、删除（单个 / 批量） 全流程操作，是整个题库系统的核心功能页；
// 交互层面：对接 Go 后端所有题库相关接口，完成前后端交互闭环，同时提供友好的操作反馈（成功 / 失败提示、确认弹窗、表单验证）；
// 技术层面：综合使用 React（useState/useEffect状态管理）、Antd（Table/Modal/Form 等核心 UI 组件）、Axios（异步请求 + 错误处理），是前端后台开发的典型实战案例。

import { useState, useEffect } from 'react'
// Antd UI组件：表格/按钮/输入框/下拉选择/确认弹窗/模态框/表单/消息提示/复选框
import { Table, Button, Input, Select, Popconfirm, Modal, Form, message, Checkbox } from 'antd'
import axios from 'axios'
// 异步请求：调用后端接口

const { Option } = Select// 解构Select的子组件

function QuestionManage() {
  // 定义所有状态
  // 题目列表相关状态
  const [questions, setQuestions] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [size, setSize] = useState(10)
  const [type, setType] = useState('')
  const [keyword, setKeyword] = useState('')
  
  // 出题弹窗相关状态
  const [modalVisible, setModalVisible] = useState(false)
  const [isAiMode, setIsAiMode] = useState(false) // 手工/AI模式
  const [form] = Form.useForm()
  //Form.useForm() 的作用是什么？和直接用 useState 管理表单数据相比，有什么优势？
  //Form.useForm() 是 Ant Design 提供的表单实例，可以方便地控制表单字段（设置值、重置、校验）。相比手动管理每个字段的状态，它减少了样板代码，并自动处理校验逻辑
  
  // 批量删除相关状态
  const [selectedRowKeys, setSelectedRowKeys] = useState([]) // 勾选的题目ID
  
  // AI出题预览相关状态
  const [previewVisible, setPreviewVisible] = useState(false) // 预览弹窗显示
  const [previewList, setPreviewList] = useState([]) // AI生成的待预览题目
  const [selectedPreviewKeys, setSelectedPreviewKeys] = useState([]) 

  // 获取题目列表（分页、筛选、搜索）
  const getQuestions = () => {
    axios.get('/api/questions', {
      params: { page, size, type, keyword }
      //params是 Axios 的 GET 请求参数配置，会自动拼接到 URL 后
    }).then(res => {
      setQuestions(res.data.list)
      setTotal(res.data.total)
    }).catch(err => {
      message.error('加载题目失败！')
      console.log('加载失败：', err)
    })
  }

  // 页面加载/参数变化时刷新列表
  useEffect(() => {
    getQuestions()
  }, [page, size, type, keyword])

  // 编辑题目：打开弹窗并填充原有数据
  const editQuestion = (q) => {
    setModalVisible(true)
    setIsAiMode(false)
    form.setFieldsValue({
      id: q.id,
      type: q.type,
      difficulty: q.difficulty,
      language: q.language,
      content: q.content,
      options: q.options,
      answer: q.answer
    })
  }

  // 删除题目
  const deleteQuestion = (ids) => {
    axios.delete('/api/questions', { data: { ids } }).then(() => {
      message.success('删除成功！')
      getQuestions() // 刷新列表
      setSelectedRowKeys([]) // 清空勾选状态
    }).catch(err => {
      message.error('删除失败！')
      console.log('删除失败：', err)
    })
  }

  // 批量删除
  const handleBatchDelete = () => {
    Popconfirm.confirm({
      title: `确定删除选中的${selectedRowKeys.length}道题吗？`,
      okText: '确定',
      cancelText: '取消',
      onConfirm: () => deleteQuestion(selectedRowKeys),
      icon: <span style={{ color: 'red' }}>⚠️</span>
    })
  }

  // AI预览题目勾选
  const handlePreviewSelect = (index) => {
    const isSelected = selectedPreviewKeys.includes(index)
    if (isSelected) {
      setSelectedPreviewKeys(selectedPreviewKeys.filter(key => key !== index))
    } else {
      setSelectedPreviewKeys([...selectedPreviewKeys, index])
    }
  }

  // AI预览全选/取消全选
  const handlePreviewSelectAll = () => {
    if (selectedPreviewKeys.length === previewList.length) {
      setSelectedPreviewKeys([])
    } else {
      setSelectedPreviewKeys(previewList.map((_, index) => index))
    }
  }

  // 提交出题
  const handleSubmit = () => {
    form.validateFields().then(values => {
      if (isAiMode) {
        // AI出题：先调用接口生成题目，再打开预览
        axios.post('/api/ai-generate', {
          type: values.type,
          count: values.count,
          difficulty: values.difficulty,
          language: values.language
        }).then(res => {
          if (!res.data.list || res.data.list.length === 0) {
            message.error('AI生成的题目为空！')
            return
          }
          // 存储生成的题目，打开预览弹窗
          setPreviewList(res.data.list)
          setSelectedPreviewKeys([]) // 初始化预览勾选状态
          setModalVisible(false) // 关闭出题弹窗
          setPreviewVisible(true) // 打开预览弹窗
        }).catch(err => {
          message.error('AI生成题目失败！')
          console.log('AI生成失败原因：', err.response?.data || err.message)
        })
      } else {
        // 手工出题：直接提交到数据库
        const submitData = {
          type: values.type,
          content: values.content,
          options: values.options,
          answer: values.answer,
          difficulty: values.difficulty,
          language: values.language
        }
        // 添加id
        if (values.id) {
          submitData.id = values.id
          // 调用编辑接口
          axios.put('/api/questions', submitData).then(() => {
            message.success('编辑题目成功！')
            setModalVisible(false)
            getQuestions()
            form.resetFields()
          }).catch(err => {
            message.error('编辑题目失败！')
            console.log('编辑失败：', err)
          })
        } else {
          axios.post('/api/questions', submitData).then(() => {
            message.success('手工添加题目成功！')
            setModalVisible(false)
            getQuestions()
            form.resetFields()
          }).catch(err => {
            message.error('手工添加题目失败！')
            console.log('手工添加失败：', err)
          })
        }
      }
    }).catch(err => {
      message.error('表单填写有误，请检查！')
      console.log('表单验证失败：', err)
    })
  }

  // 确认添加预览中的题目
  const confirmAddPreviewQuestions = () => {
    if (selectedPreviewKeys.length === 0) {
      message.warning('请选择要添加的题目！')
      return
    }
    // 筛选出勾选的题目
    const selectedQuestions = selectedPreviewKeys.map(index => previewList[index])
    // 批量添加到数据库
    const addPromises = selectedQuestions.map(q => axios.post('/api/questions', q))
    Promise.all(addPromises)
      .then(() => {
        message.success(`成功添加${selectedQuestions.length}道题目！`)
        setPreviewVisible(false) // 关闭预览弹窗
        getQuestions() // 刷新题库列表
        setSelectedPreviewKeys([]) // 清空勾选状态
      })
      .catch(() => {
        message.error('批量添加题目失败！')
      })
  }

  // 表格列配置
  const columns = [
    {
      title: '选择',
      type: 'selection',
      key: 'selection',
      width: 60,
    },
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '题型', dataIndex: 'type', width: 100 },
    { title: '题目内容', dataIndex: 'content', ellipsis: true, width: 300 },
    { title: '难度', dataIndex: 'difficulty', width: 100 },
    { title: '编程语言', dataIndex: 'language', width: 120 },
    {
      title: '操作',
      width: 180,
      render: (_, record) => (
        <div style={{ display: 'flex', gap: 8 }}>
          <Button size="small" onClick={() => editQuestion(record)}>编辑</Button>
          <Popconfirm
            title="确定删除这道题吗？"
            onConfirm={() => deleteQuestion([record.id])}
            okText="确定"
            cancelText="取消"
          >
            <Button size="small" danger>删除</Button>
          </Popconfirm>
        </div>
      )
    }
  ]

  // AI预览表格列配置
  const previewColumns = [
    {
      title: '选择',
      key: 'select',
      width: 60,
      render: (_, record, index) => (
        <Checkbox
          checked={selectedPreviewKeys.includes(index)}
          onChange={() => handlePreviewSelect(index)}
        />
      )
    },
    { title: '题型', dataIndex: 'type', width: 100 },
    { title: '题目内容', dataIndex: 'content', ellipsis: true, width: 350 },
    { title: '选项', dataIndex: 'options', ellipsis: true, width: 250, render: (options) => options || '无' },
    { title: '答案', dataIndex: 'answer', width: 100, render: (answer) => answer || '无' },
    { title: '难度', dataIndex: 'difficulty', width: 100 },
    { title: '编程语言', dataIndex: 'language', width: 120, render: (lang) => lang || '无' },
  ]

  return (
    <div style={{ background: '#fff', padding: 20, borderRadius: 4 }}>
      {/* 顶部操作栏 */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16, alignItems: 'center' }}>
        {/* 出题按钮 */}
        <Button type="primary" onClick={() => setModalVisible(true)}>
          出题
        </Button>
        {/* 批量删除按钮 */}
        <Button
          danger
          disabled={selectedRowKeys.length === 0}
          onClick={handleBatchDelete}
        >
          批量删除（{selectedRowKeys.length}）
        </Button>
        {/* 题型筛选 */}
        <Select
          placeholder="筛选题型"
          value={type}
          onChange={setType}
          style={{ width: 120 }}
        >
          <Option value="">全部题型</Option>
          <Option value="单选题">单选题</Option>
          <Option value="多选题">多选题</Option>
          <Option value="编程题">编程题</Option>
        </Select>
        {/* 关键词搜索 */}
        <Input
          placeholder="搜索题目内容"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          style={{ flex: 1, maxWidth: 300 }}
          onPressEnter={getQuestions} // 按回车搜索
        />
      </div>

      {/* 题库表格 */}
      <Table
        columns={columns}
        dataSource={questions}
        rowKey="id"
        pagination={{
          current: page,
          pageSize: size,
          total: total,
          onChange: (p, s) => { setPage(p); setSize(s) },
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 道题目`
        }}
        bordered
        rowSelection={{
          selectedRowKeys,
          onChange: (keys) => setSelectedRowKeys(keys),
          columnWidth: 60
        }}
        scroll={{ x: 'max-content' }}
      />

      {/* 出题弹窗 */}
      <Modal
        title={isAiMode ? "AI出题" : form.getFieldValue('id') ? "编辑题目" : "手工出题"}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
          setIsAiMode(false)
        }}
        onOk={handleSubmit}
        width={600}
        destroyOnClose={true}
      >
        <Form form={form} layout="vertical" labelCol={{ span: 4 }} wrapperCol={{ span: 20 }}>
          {/* 隐藏ID字段 */}
          <Form.Item name="id" hidden />
          
          {/* 通用字段 */}
          <Form.Item
            name="type"
            label="题型"
            rules={[{ required: true, message: '请选择题型！' }]}
          >
            <Select disabled={form.getFieldValue('id') ? true : false}>
              <Option value="单选题">单选题</Option>
              <Option value="多选题">多选题</Option>
              <Option value="编程题">编程题</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="difficulty"
            label="难度"
            rules={[{ required: true, message: '请选择难度！' }]}
          >
            <Select>
              <Option value="简单">简单</Option>
              <Option value="中等">中等</Option>
              <Option value="困难">困难</Option>
            </Select>
          </Form.Item>
          <Form.Item name="language" label="编程语言">
            <Select>
              <Option value="">无</Option>
              <Option value="Go">Go</Option>
              <Option value="JavaScript">JavaScript</Option>
            </Select>
          </Form.Item>

          {/* AI出题 */}
          {isAiMode && (
            <Form.Item
              name="count"
              label="题目数量"
              rules={[
                { required: true, message: '请输入题目数量！' },
                { type: 'number', min: 1, max: 10, message: '数量需在1-10之间！' }
              ]}
              valuePropName="valueAsNumber"
            >
              <Input type="number" min={1} max={10} placeholder="请输入1-10的数字" />
            </Form.Item>
          )}

          {/* 手工出题 */}
          {!isAiMode && (
            <>
              <Form.Item
                name="content"
                label="题目内容"
                rules={[{ required: true, message: '请输入题目内容！' }]}
              >
                <Input.TextArea rows={4} placeholder="请输入题目内容" />
              </Form.Item>
              <Form.Item name="options" label="选项（逗号分隔）">
                <Input placeholder="例：A.选项1,B.选项2,C.选项3,D.选项4" />
              </Form.Item>
              <Form.Item name="answer" label="正确答案">
                <Input placeholder="例：A 或 AB（多选题）" />
              </Form.Item>
            </>
          )}
        </Form>

        {/* 手工/AI切换 */}
        {!form.getFieldValue('id') && (
          <div style={{ textAlign: 'center', marginTop: 16 }}>
            <Button
              onClick={() => setIsAiMode(false)}
              type={!isAiMode ? "primary" : ""}
            >
              手工出题
            </Button>
            <Button
              onClick={() => setIsAiMode(true)}
              type={isAiMode ? "primary" : ""}
              style={{ marginLeft: 8 }}
            >
              AI出题
            </Button>
          </div>
        )}
      </Modal>

      {/* AI出题预览 */}
      <Modal
        title="AI生成题目预览"
        open={previewVisible}
        onCancel={() => {
          setPreviewVisible(false)
          setSelectedPreviewKeys([])
        }}
        onOk={confirmAddPreviewQuestions}
        width={1200}
        destroyOnClose={true}
        maskClosable={false}
      >
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
          <Button onClick={handlePreviewSelectAll}>
            {selectedPreviewKeys.length === previewList.length ? '取消全选' : '全选'}
          </Button>
          <span>已选择：{selectedPreviewKeys.length}/{previewList.length} 道题</span>
        </div>
        <Table
          columns={previewColumns}
          dataSource={previewList}
          rowKey={(record, index) => `preview-${index}`}
          pagination={false}
          bordered
          scroll={{ x: 'max-content' }}
        />
      </Modal>
    </div>
  )
}

export default QuestionManage