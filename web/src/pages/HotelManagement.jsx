import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Input,
  Select,
  Space,
  Tag,
  Popconfirm,
  Modal,
  Form,
  Card,
  App,
} from 'antd'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  ShopOutlined,
  ReloadOutlined,
} from '@ant-design/icons'
import axios from 'axios'

const { Option } = Select

const HotelManagement = () => {
  const { message } = App.useApp()
  const [dataSource, setDataSource] = useState([])
  const [loading, setLoading] = useState(false)
  const [filters, setFilters] = useState({ status: '', keyword: '' })
  const [modalVisible, setModalVisible] = useState(false)
  const [editingBranch, setEditingBranch] = useState(null)
  const [form] = Form.useForm()

  // 获取分店列表
  const fetchBranches = async () => {
    setLoading(true)
    try {
      const params = {}
      if (filters.status) params.status = filters.status

      const response = await axios.get('/api/v1/branches', { params })
      if (response.data.code === 200) {
        setDataSource(response.data.data.list || [])
      } else {
        message.error(response.data.msg || '获取分店列表失败')
      }
    } catch (error) {
      message.error('获取分店列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchBranches()
  }, [filters])

  // 处理添加
  const handleAdd = () => {
    setEditingBranch(null)
    form.resetFields()
    form.setFieldsValue({ status: 'ACTIVE' })
    setModalVisible(true)
  }

  // 处理编辑
  const handleEdit = (record) => {
    setEditingBranch(record)
    form.setFieldsValue({
      hotel_name: record.hotel_name,
      branch_code: record.branch_code,
      address: record.address,
      contact: record.contact,
      contact_phone: record.contact_phone,
      status: record.status,
    })
    setModalVisible(true)
  }

  // 处理删除
  const handleDelete = async (id) => {
    try {
      const response = await axios.delete(`/api/v1/branches/${id}`)
      if (response.data.code === 200) {
        message.success('删除成功')
        fetchBranches()
      } else {
        message.error(response.data.msg || '删除失败')
      }
    } catch (error) {
      message.error('删除失败：' + (error.response?.data?.msg || error.message))
    }
  }

  // 处理提交
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      
      if (editingBranch) {
        // 更新
        const response = await axios.put(`/api/v1/branches/${editingBranch.id}`, values)
        if (response.data.code === 200) {
          message.success('更新成功')
          setModalVisible(false)
          fetchBranches()
        } else {
          message.error(response.data.msg || '更新失败')
        }
      } else {
        // 创建
        const response = await axios.post('/api/v1/branches', values)
        if (response.data.code === 200) {
          message.success('创建成功')
          setModalVisible(false)
          fetchBranches()
        } else {
          message.error(response.data.msg || '创建失败')
        }
      }
    } catch (error) {
      message.error((editingBranch ? '更新' : '创建') + '失败：' + (error.response?.data?.msg || error.message))
    }
  }

  const columns = [
    {
      title: '分店名称',
      dataIndex: 'hotel_name',
      key: 'hotel_name',
      width: 200,
    },
    {
      title: '分店编码',
      dataIndex: 'branch_code',
      key: 'branch_code',
      width: 120,
    },
    {
      title: '地址',
      dataIndex: 'address',
      key: 'address',
      ellipsis: true,
    },
    {
      title: '联系人',
      dataIndex: 'contact',
      key: 'contact',
      width: 100,
    },
    {
      title: '联系电话',
      dataIndex: 'contact_phone',
      key: 'contact_phone',
      width: 130,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (status) => (
        <Tag color={status === 'ACTIVE' ? 'success' : 'default'}>
          {status === 'ACTIVE' ? '启用' : '停用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个分店吗？"
            description="删除后不可恢复，且分店下不能有房源"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card
        variant="borderless"
        title={
          <Space direction="vertical" size={0}>
            <span>
              <ShopOutlined style={{ marginRight: 8 }} />
              酒店分店管理
            </span>
            <span style={{ fontSize: 12, color: '#999', fontWeight: 'normal' }}>
              提示：相同酒店名称的分店会自动分组，可通过顶部导航栏切换酒店和分店
            </span>
          </Space>
        }
        extra={
          <Space>
            <Select
              style={{ width: 120 }}
              placeholder="状态筛选"
              allowClear
              value={filters.status}
              onChange={(value) => setFilters({ ...filters, status: value || '' })}
            >
              <Option value="ACTIVE">启用</Option>
              <Option value="INACTIVE">停用</Option>
            </Select>
            <Button icon={<ReloadOutlined />} onClick={fetchBranches}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              添加分店
            </Button>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={dataSource}
          rowKey="id"
          loading={loading}
          pagination={{
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
            showQuickJumper: true,
          }}
        />
      </Card>

      <Modal
        title={editingBranch ? '编辑分店' : '添加分店'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={handleSubmit}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="hotel_name"
            label="酒店名称"
            rules={[{ required: true, message: '请输入酒店名称' }]}
            tooltip="相同酒店名称的分店会自动分组，如：天极酒店、如家快捷"
          >
            <Input placeholder="如：天极酒店、如家快捷、锦江之星" />
          </Form.Item>

          <Form.Item 
            name="branch_code" 
            label="分店编码"
            tooltip="用于区分同一酒店的不同分店，如：总店、朝阳店、BJ001"
          >
            <Input placeholder="如：总店、朝阳店、BJ001（不填自动生成）" />
          </Form.Item>

          <Form.Item
            name="address"
            label="分店地址"
            rules={[{ required: true, message: '请输入分店地址' }]}
          >
            <Input placeholder="请输入详细地址" />
          </Form.Item>

          <Form.Item
            name="contact"
            label="联系人"
            rules={[{ required: true, message: '请输入联系人姓名' }]}
          >
            <Input placeholder="请输入联系人姓名" />
          </Form.Item>

          <Form.Item
            name="contact_phone"
            label="联系电话"
            rules={[
              { required: true, message: '请输入联系电话' },
              { pattern: /^1[3-9]\d{9}$|^0\d{2,3}-?\d{7,8}$/, message: '请输入正确的电话号码' },
            ]}
          >
            <Input placeholder="请输入联系电话" />
          </Form.Item>

          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
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

export default HotelManagement
