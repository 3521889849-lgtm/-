import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Space,
  Popconfirm,
  Tag,
  InputNumber,
  App,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import {
  getMembers,
  getMember,
  createMember,
  updateMember,
  deleteMember,
} from '../api/member'

const { Option } = Select
const { TextArea } = Input

const MEMBER_LEVELS = [
  { label: '普通会员', value: 'NORMAL' },
  { label: '黄金会员', value: 'GOLD' },
  { label: '钻石会员', value: 'DIAMOND' },
]

const MemberManagement = () => {
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
      const response = await getMembers(params)
      if (response && response.list) {
        setTableData(response.list || [])
        setPagination({
          ...pagination,
          total: response.total || 0,
        })
      } else {
        message.error('获取会员列表失败')
      }
    } catch (error) {
      message.error('获取会员列表失败: ' + error.message)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenModal = async (record = null) => {
    setEditingRecord(record)
    if (record) {
      try {
        const response = await getMember(record.id)
        if (response) {
          form.setFieldsValue({
            guest_id: response.guest_id,
            member_level: response.member_level,
            points_balance: response.points_balance,
            status: response.status,
          })
        }
      } catch (error) {
        message.error('获取会员详情失败')
      }
    } else {
      form.resetFields()
    }
    setModalVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingRecord) {
        const response = await updateMember(editingRecord.id, values)
        if (response && response.code === 200) {
          message.success('更新成功')
          setModalVisible(false)
          fetchData()
        } else {
          message.error(response?.msg || '更新失败')
        }
      } else {
        const response = await createMember(values)
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
      const response = await deleteMember(id)
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
      title: '会员ID',
      dataIndex: 'id',
      key: 'id',
      width: 100,
    },
    {
      title: '客人姓名',
      dataIndex: 'guest_name',
      key: 'guest_name',
    },
    {
      title: '客人手机',
      dataIndex: 'guest_phone',
      key: 'guest_phone',
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
      title: '积分余额',
      dataIndex: 'points_balance',
      key: 'points_balance',
      render: (balance) => balance || 0,
    },
    {
      title: '注册时间',
      dataIndex: 'register_time',
      key: 'register_time',
    },
    {
      title: '最后入住时间',
      dataIndex: 'last_check_in_time',
      key: 'last_check_in_time',
      render: (time) => time || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 'ACTIVE' ? 'green' : 'red'}>
          {status === 'ACTIVE' ? '启用' : '冻结'}
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
            title="确定要删除这个会员吗？"
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
        <h2>会员管理</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          新建会员
        </Button>
      </div>

      <div style={{ marginBottom: 16, padding: '16px', background: '#f5f5f5', borderRadius: '4px' }}>
        <Space>
          <Input
            placeholder="搜索姓名或手机号"
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
            <Option value="FROZEN">冻结</Option>
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
        title={editingRecord ? '编辑会员' : '新建会员'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="guest_id"
            label="客人ID"
            rules={[{ required: true, message: '请输入客人ID' }]}
          >
            <InputNumber
              placeholder="请输入客人ID"
              style={{ width: '100%' }}
              min={1}
            />
          </Form.Item>
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
          <Form.Item name="points_balance" label="积分余额">
            <InputNumber
              placeholder="请输入积分余额"
              style={{ width: '100%' }}
              min={0}
            />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="ACTIVE">
            <Select>
              <Option value="ACTIVE">启用</Option>
              <Option value="FROZEN">冻结</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MemberManagement
