import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  App,
  Space,
  Tag,
  InputNumber,
  DatePicker,
  Card,
  Row,
  Col,
} from 'antd'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
} from '@ant-design/icons'
import {
  getPointsRecords,
  createPointsRecord,
} from '../api/member'
import dayjs from 'dayjs'
import './MemberPointsManagement.css'

const { Option } = Select
const { TextArea } = Input

const CHANGE_TYPES = [
  { label: '获取', value: 'EARN' },
  { label: '消费', value: 'CONSUME' },
]

const MemberPointsManagement = () => {
  const { message } = App.useApp()
  const [form] = Form.useForm()
  const [tableData, setTableData] = useState([])
  const [loading, setLoading] = useState(false)
  const [submitLoading, setSubmitLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
    showSizeChanger: true,
    showQuickJumper: true,
  })
  const [filters, setFilters] = useState({
    member_id: undefined,
    order_id: undefined,
    change_type: undefined,
    start_time: undefined,
    end_time: undefined,
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
      const response = await getPointsRecords(params)
      if (response && response.list) {
        setTableData(response.list || [])
        setPagination({
          ...pagination,
          total: response.total || 0,
        })
      } else {
        message.error('获取积分记录失败')
      }
    } catch (error) {
      message.error('获取积分记录失败: ' + error.message)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenModal = () => {
    form.resetFields()
    form.setFieldsValue({
      change_type: 'EARN',
    })
    setModalVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setSubmitLoading(true)
      const submitData = {
        ...values,
        points_value: values.change_type === 'CONSUME' ? -Math.abs(values.points_value) : Math.abs(values.points_value),
      }
      const response = await createPointsRecord(submitData)
      if (response && response.code === 200) {
        message.success('创建成功')
        setModalVisible(false)
        fetchData()
      } else {
        message.error(response?.msg || '创建失败')
      }
    } catch (error) {
      console.error('提交失败:', error)
    } finally {
      setSubmitLoading(false)
    }
  }

  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    // useEffect will trigger fetchData
  }

  const handleReset = () => {
    setFilters({
      member_id: undefined,
      order_id: undefined,
      change_type: undefined,
      start_time: undefined,
      end_time: undefined,
    })
    setPagination({ ...pagination, current: 1 })
  }

  const columns = [
    {
      title: '记录ID',
      dataIndex: 'id',
      key: 'id',
      width: 100,
    },
    {
      title: '会员ID',
      dataIndex: 'member_id',
      key: 'member_id',
      width: 120,
    },
    {
      title: '会员姓名',
      dataIndex: 'member_name',
      key: 'member_name',
      render: (name) => name || '-',
    },
    {
      title: '订单ID',
      dataIndex: 'order_id',
      key: 'order_id',
      render: (id) => id || '-',
    },
    {
      title: '变动类型',
      dataIndex: 'change_type',
      key: 'change_type',
      render: (type) => {
        const typeMap = {
          EARN: { text: '获取', color: 'success' },
          CONSUME: { text: '消费', color: 'error' },
        }
        const info = typeMap[type] || { text: type, color: 'default' }
        return <Tag color={info.color} className="status-tag">{info.text}</Tag>
      },
    },
    {
      title: '积分值',
      dataIndex: 'points_value',
      key: 'points_value',
      render: (value) => {
        const isPositive = value > 0
        return (
          <span style={{ 
            color: isPositive ? '#52c41a' : '#ff4d4f',
            fontWeight: 500 
          }}>
            {isPositive ? '+' : ''}{value}
          </span>
        )
      },
    },
    {
      title: '变动原因',
      dataIndex: 'change_reason',
      key: 'change_reason',
      ellipsis: true,
    },
    {
      title: '变动时间',
      dataIndex: 'change_time',
      key: 'change_time',
      width: 180,
    },
    {
      title: '操作人ID',
      dataIndex: 'operator_id',
      key: 'operator_id',
      width: 120,
    },
  ]

  return (
    <div className="page-container">
      <div className="page-header">
        <h2 className="page-title">会员积分管理</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={handleOpenModal}
          size="large"
        >
          新建积分记录
        </Button>
      </div>

      <Card className="content-card" bordered={false}>
        <Row gutter={[16, 16]} align="middle" className="search-form">
          <Col xs={24} sm={12} md={6} lg={4}>
            <InputNumber
              placeholder="会员ID"
              value={filters.member_id}
              onChange={(value) => setFilters({ ...filters, member_id: value })}
              style={{ width: '100%' }}
              min={1}
            />
          </Col>
          <Col xs={24} sm={12} md={6} lg={4}>
            <InputNumber
              placeholder="订单ID"
              value={filters.order_id}
              onChange={(value) => setFilters({ ...filters, order_id: value })}
              style={{ width: '100%' }}
              min={1}
            />
          </Col>
          <Col xs={24} sm={12} md={6} lg={4}>
            <Select
              placeholder="变动类型"
              value={filters.change_type}
              onChange={(value) => setFilters({ ...filters, change_type: value })}
              style={{ width: '100%' }}
              allowClear
            >
              {CHANGE_TYPES.map((type) => (
                <Option key={type.value} value={type.value}>
                  {type.label}
                </Option>
              ))}
            </Select>
          </Col>
          <Col xs={24} sm={12} md={6} lg={4}>
            <DatePicker
              placeholder="开始时间"
              value={filters.start_time ? dayjs(filters.start_time) : null}
              onChange={(date) => setFilters({ ...filters, start_time: date ? date.format('YYYY-MM-DD HH:mm:ss') : undefined })}
              showTime
              format="YYYY-MM-DD HH:mm:ss"
              style={{ width: '100%' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6} lg={4}>
            <DatePicker
              placeholder="结束时间"
              value={filters.end_time ? dayjs(filters.end_time) : null}
              onChange={(date) => setFilters({ ...filters, end_time: date ? date.format('YYYY-MM-DD HH:mm:ss') : undefined })}
              showTime
              format="YYYY-MM-DD HH:mm:ss"
              style={{ width: '100%' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6} lg={4} style={{ display: 'flex', gap: '8px' }}>
            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              查询
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleReset}>
              重置
            </Button>
          </Col>
        </Row>
      </Card>

      <Card className="content-card" bordered={false} bodyStyle={{ padding: 0 }}>
        <Table
          columns={columns}
          dataSource={tableData}
          loading={loading}
          rowKey="id"
          className="custom-table"
          pagination={{
            ...pagination,
            onChange: (page, pageSize) => {
              setPagination({ ...pagination, current: page, pageSize })
            },
          }}
          scroll={{ x: 1000 }}
        />
      </Card>

      <Modal
        title="新建积分记录"
        open={modalVisible}
        onOk={handleSubmit}
        confirmLoading={submitLoading}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        width={520}
        centered
        destroyOnClose
        maskClosable={false}
      >
        <Form form={form} layout="vertical" requiredMark="optional">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="member_id"
                label="会员ID"
                rules={[{ required: true, message: '请输入会员ID' }]}
              >
                <InputNumber
                  placeholder="请输入"
                  style={{ width: '100%' }}
                  min={1}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="order_id" label="订单ID（可选）">
                <InputNumber
                  placeholder="请输入"
                  style={{ width: '100%' }}
                  min={1}
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="change_type"
                label="变动类型"
                rules={[{ required: true, message: '请选择变动类型' }]}
              >
                <Select placeholder="请选择">
                  {CHANGE_TYPES.map((type) => (
                    <Option key={type.value} value={type.value}>
                      {type.label}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="points_value"
                label="积分值"
                rules={[{ required: true, message: '请输入积分值' }]}
              >
                <InputNumber
                  placeholder="请输入正数"
                  style={{ width: '100%' }}
                  min={1}
                />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="change_reason"
            label="变动原因"
            rules={[{ required: true, message: '请输入变动原因' }]}
          >
            <TextArea rows={4} placeholder="请输入详细的变动原因" showCount maxLength={200} />
          </Form.Item>

          <Form.Item
            name="operator_id"
            label="操作人ID"
            rules={[{ required: true, message: '请输入操作人ID' }]}
          >
            <InputNumber
              placeholder="请输入"
              style={{ width: '100%' }}
              min={1}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MemberPointsManagement
