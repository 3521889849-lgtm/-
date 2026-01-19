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
  Popconfirm,
  Tag,
  DatePicker,
  InputNumber,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import {
  getMemberRights,
  getMemberRight,
  createMemberRights,
  updateMemberRights,
  deleteMemberRights,
} from '../api/member'
import dayjs from 'dayjs'

const { Option } = Select
const { TextArea } = Input

const MEMBER_LEVELS = [
  { label: '普通会员', value: 'NORMAL' },
  { label: '黄金会员', value: 'GOLD' },
  { label: '钻石会员', value: 'DIAMOND' },
]

const MemberRightsManagement = () => {
  const { message } = App.useApp()
  const [form] = Form.useForm()
  const [tableData, setTableData] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState(null)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })
  const [filters, setFilters] = useState({
    member_level: undefined,
    status: undefined,
    keyword: undefined,
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
      const response = await getMemberRights(params)
      if (response && response.list) {
        setTableData(response.list || [])
        setPagination({
          ...pagination,
          total: response.total || 0,
        })
      } else {
        message.error('获取权益列表失败')
      }
    } catch (error) {
      message.error('获取权益列表失败: ' + error.message)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenModal = async (record = null) => {
    setEditingRecord(record)
    if (record) {
      try {
        const response = await getMemberRight(record.id)
        if (response) {
          form.setFieldsValue({
            member_level: response.member_level,
            rights_name: response.rights_name,
            description: response.description,
            discount_ratio: response.discount_ratio,
            effective_time: response.effective_time ? dayjs(response.effective_time) : null,
            expire_time: response.expire_time ? dayjs(response.expire_time) : null,
            status: response.status,
          })
        }
      } catch (error) {
        message.error('获取权益详情失败')
      }
    } else {
      form.resetFields()
      form.setFieldsValue({
        effective_time: dayjs(),
        status: 'ACTIVE',
      })
    }
    setModalVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const submitData = {
        ...values,
        effective_time: values.effective_time
          ? values.effective_time.format('YYYY-MM-DD HH:mm:ss')
          : dayjs().format('YYYY-MM-DD HH:mm:ss'),
        expire_time: values.expire_time
          ? values.expire_time.format('YYYY-MM-DD HH:mm:ss')
          : undefined,
      }
      if (editingRecord) {
        const response = await updateMemberRights(editingRecord.id, submitData)
        if (response && response.code === 200) {
          message.success('更新成功')
          setModalVisible(false)
          fetchData()
        } else {
          message.error(response?.msg || '更新失败')
        }
      } else {
        const response = await createMemberRights(submitData)
        if (response && response.code === 200) {
          message.success('创建成功')
          setModalVisible(false)
          fetchData()
        } else {
          message.error(response?.msg || '创建失败')
        }
      }
    } catch (error) {
      console.error('提交失败:', error)
    }
  }

  const handleDelete = async (id) => {
    try {
      const response = await deleteMemberRights(id)
      if (response && response.code === 200) {
        message.success('删除成功')
        fetchData()
      } else {
        message.error(response?.msg || '删除失败')
      }
    } catch (error) {
      message.error('删除失败: ' + error.message)
    }
  }

  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    fetchData()
  }

  const handleReset = () => {
    setFilters({
      member_level: undefined,
      status: undefined,
      keyword: undefined,
    })
    setPagination({ ...pagination, current: 1 })
  }

  const columns = [
    {
      title: '权益ID',
      dataIndex: 'id',
      key: 'id',
      width: 100,
    },
    {
      title: '会员等级',
      dataIndex: 'member_level',
      key: 'member_level',
      render: (level) => {
        const levelMap = {
          NORMAL: { text: '普通会员', color: 'default' },
          GOLD: { text: '黄金会员', color: 'gold' },
          DIAMOND: { text: '钻石会员', color: 'purple' },
        }
        const info = levelMap[level] || { text: level, color: 'default' }
        return <Tag color={info.color}>{info.text}</Tag>
      },
    },
    {
      title: '权益名称',
      dataIndex: 'rights_name',
      key: 'rights_name',
    },
    {
      title: '权益描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '折扣比例',
      dataIndex: 'discount_ratio',
      key: 'discount_ratio',
      render: (ratio) => (ratio ? `${ratio}%` : '-'),
    },
    {
      title: '生效时间',
      dataIndex: 'effective_time',
      key: 'effective_time',
    },
    {
      title: '失效时间',
      dataIndex: 'expire_time',
      key: 'expire_time',
      render: (time) => time || '永久有效',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 'ACTIVE' ? 'green' : 'red'}>
          {status === 'ACTIVE' ? '启用' : '停用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个权益吗？"
            onConfirm={() => handleDelete(record.id)}
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>会员权益管理</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          新建权益
        </Button>
      </div>

      <div style={{ marginBottom: 16, padding: '16px', background: '#f5f5f5', borderRadius: '4px' }}>
        <Space>
          <Input
            placeholder="搜索权益名称或描述"
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            style={{ width: 200 }}
            allowClear
          />
          <Select
            placeholder="会员等级"
            value={filters.member_level}
            onChange={(value) => setFilters({ ...filters, member_level: value })}
            style={{ width: 150 }}
            allowClear
          >
            {MEMBER_LEVELS.map((level) => (
              <Option key={level.value} value={level.value}>
                {level.label}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="状态"
            value={filters.status}
            onChange={(value) => setFilters({ ...filters, status: value })}
            style={{ width: 120 }}
            allowClear
          >
            <Option value="ACTIVE">启用</Option>
            <Option value="INACTIVE">停用</Option>
          </Select>
          <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
            查询
          </Button>
          <Button onClick={handleReset}>重置</Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={tableData}
        loading={loading}
        rowKey="id"
        pagination={{
          ...pagination,
          onChange: (page, pageSize) => {
            setPagination({ ...pagination, current: page, pageSize })
          },
        }}
      />

      <Modal
        title={editingRecord ? '编辑权益' : '新建权益'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        width={700}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="member_level"
            label="会员等级"
            rules={[{ required: true, message: '请选择会员等级' }]}
          >
            <Select placeholder="请选择会员等级">
              {MEMBER_LEVELS.map((level) => (
                <Option key={level.value} value={level.value}>
                  {level.label}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="rights_name"
            label="权益名称"
            rules={[{ required: true, message: '请输入权益名称' }]}
          >
            <Input placeholder="请输入权益名称（如：专属预订、积分兑换、房价折扣等）" />
          </Form.Item>
          <Form.Item name="description" label="权益描述">
            <TextArea rows={3} placeholder="请输入权益描述" />
          </Form.Item>
          <Form.Item name="discount_ratio" label="折扣比例（%）">
            <InputNumber
              placeholder="请输入折扣比例"
              style={{ width: '100%' }}
              min={0}
              max={100}
              precision={2}
            />
          </Form.Item>
          <Form.Item
            name="effective_time"
            label="生效时间"
            rules={[{ required: true, message: '请选择生效时间' }]}
          >
            <DatePicker
              showTime
              format="YYYY-MM-DD HH:mm:ss"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item name="expire_time" label="失效时间">
            <DatePicker
              showTime
              format="YYYY-MM-DD HH:mm:ss"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="ACTIVE">
            <Select>
              <Option value="ACTIVE">启用</Option>
              <Option value="INACTIVE">停用</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MemberRightsManagement
