import React, { useState, useEffect } from 'react'
import { Table, Button, Form, Input, Select, DatePicker, Space, Tag, App } from 'antd'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons'
import axios from 'axios'
import dayjs from 'dayjs'

const { Option } = Select
const { RangePicker } = DatePicker

const API_BASE = '/api/v1'

const MODULES = [
  { label: '房源管理', value: '房源管理' },
  { label: '订单处理', value: '订单处理' },
  { label: '客人管理', value: '客人管理' },
  { label: '会员管理', value: '会员管理' },
  { label: '财务管理', value: '财务管理' },
  { label: '系统配置', value: '系统配置' },
]

const OPERATION_TYPES = [
  { label: '查询', value: '查询' },
  { label: '添加', value: '添加' },
  { label: '修改', value: '修改' },
  { label: '删除', value: '删除' },
  { label: '导出', value: '导出' },
]

const OperationLogManagement = () => {
  const { message } = App.useApp()
  const [form] = Form.useForm()
  const [tableData, setTableData] = useState([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })
  const [filters, setFilters] = useState({
    operator_id: undefined,
    module: undefined,
    operation_type: undefined,
    start_time: undefined,
    end_time: undefined,
    is_success: undefined,
  })

  useEffect(() => {
    fetchData()
  }, [pagination.current, pagination.pageSize, filters])

  const fetchData = async () => {
    setLoading(true)
    try {
      const params = {
        page: pagination.current,
        page_size: pagination.pageSize,
        ...filters,
      }
      if (params.start_time) {
        params.start_time = dayjs(params.start_time).format('YYYY-MM-DD HH:mm:ss')
      }
      if (params.end_time) {
        params.end_time = dayjs(params.end_time).format('YYYY-MM-DD HH:mm:ss')
      }
      const response = await axios.get(`${API_BASE}/operation-logs`, { params })
      if (response.data && response.data.list) {
        setTableData(response.data.list || [])
        setPagination({
          ...pagination,
          total: response.data.total || 0,
        })
      } else {
        message.error('获取操作日志失败')
      }
    } catch (error) {
      message.error('获取操作日志失败: ' + error.message)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (values) => {
    const newFilters = {
      operator_id: values.operator_id,
      module: values.module,
      operation_type: values.operation_type,
      is_success: values.is_success,
    }
    if (values.dateRange && values.dateRange.length === 2) {
      newFilters.start_time = values.dateRange[0]
      newFilters.end_time = values.dateRange[1]
    }
    setFilters(newFilters)
    setPagination({ ...pagination, current: 1 })
  }

  const handleReset = () => {
    form.resetFields()
    setFilters({})
    setPagination({ ...pagination, current: 1 })
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '操作人',
      dataIndex: 'operator_name',
      key: 'operator_name',
      width: 120,
    },
    {
      title: '操作模块',
      dataIndex: 'module',
      key: 'module',
      width: 120,
    },
    {
      title: '操作类型',
      dataIndex: 'operation_type',
      key: 'operation_type',
      width: 100,
      render: (text) => {
        const colorMap = {
          查询: 'blue',
          添加: 'green',
          修改: 'orange',
          删除: 'red',
          导出: 'purple',
        }
        return <Tag color={colorMap[text] || 'default'}>{text}</Tag>
      },
    },
    {
      title: '操作内容',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
    },
    {
      title: '操作时间',
      dataIndex: 'operation_time',
      key: 'operation_time',
      width: 180,
    },
    {
      title: '操作IP',
      dataIndex: 'operation_ip',
      key: 'operation_ip',
      width: 140,
    },
    {
      title: '关联ID',
      dataIndex: 'related_id',
      key: 'related_id',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'is_success',
      key: 'is_success',
      width: 80,
      render: (isSuccess) => (
        <Tag color={isSuccess ? 'success' : 'error'}>
          {isSuccess ? '成功' : '失败'}
        </Tag>
      ),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Form
        form={form}
        layout="inline"
        onFinish={handleSearch}
        style={{ marginBottom: '16px' }}
      >
        <Form.Item name="operator_id" label="操作人ID">
          <Input placeholder="请输入操作人ID" style={{ width: 150 }} />
        </Form.Item>
        <Form.Item name="module" label="操作模块">
          <Select placeholder="请选择模块" style={{ width: 150 }} allowClear>
            {MODULES.map((item) => (
              <Option key={item.value} value={item.value}>
                {item.label}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item name="operation_type" label="操作类型">
          <Select placeholder="请选择类型" style={{ width: 120 }} allowClear>
            {OPERATION_TYPES.map((item) => (
              <Option key={item.value} value={item.value}>
                {item.label}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item name="is_success" label="状态">
          <Select placeholder="请选择状态" style={{ width: 100 }} allowClear>
            <Option value="true">成功</Option>
            <Option value="false">失败</Option>
          </Select>
        </Form.Item>
        <Form.Item name="dateRange" label="操作时间">
          <RangePicker showTime format="YYYY-MM-DD HH:mm:ss" />
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
              查询
            </Button>
            <Button onClick={handleReset} icon={<ReloadOutlined />}>
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>

      <Table
        columns={columns}
        dataSource={tableData}
        rowKey="id"
        loading={loading}
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: (page, pageSize) => {
            setPagination({ ...pagination, current: page, pageSize })
          },
        }}
      />
    </div>
  )
}

export default OperationLogManagement
