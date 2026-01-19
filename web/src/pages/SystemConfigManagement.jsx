import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  message,
  Space,
  Popconfirm,
  Tag,
  Tabs,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons'

const { Option } = Select
import {
  getSystemConfigs,
  createSystemConfig,
  updateSystemConfig,
  deleteSystemConfig,
  getSystemConfigsByCategory,
} from '../api/user'

const { TextArea } = Input

const CATEGORIES = [
  { key: 'CONSUMPTION', label: '消费项设置' },
  { key: 'SMS_TEMPLATE', label: '短信模板' },
  { key: 'PRINT_SETTINGS', label: '打印设置' },
  { key: 'MEMBER_RULES', label: '会员管理规则' },
]

const SystemConfigManagement = () => {
  const [form] = Form.useForm()
  const [tableData, setTableData] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState(null)
  const [activeCategory, setActiveCategory] = useState('CONSUMPTION')
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  // 获取系统配置列表
  const fetchData = async (category, page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await getSystemConfigs({
        page,
        page_size: pageSize,
        config_category: category,
      })
      if (res.data.code === 200) {
        setTableData(res.data.data?.list || [])
        setPagination({
          current: res.data.data?.page || 1,
          pageSize: res.data.data?.page_size || 10,
          total: res.data.data?.total || 0,
        })
      }
    } catch (error) {
      message.error('获取系统配置列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData(activeCategory)
  }, [activeCategory])

  // 打开新增/编辑弹窗
  const handleOpenModal = (record = null) => {
    setEditingRecord(record)
    if (record) {
      form.setFieldsValue(record)
    } else {
      form.resetFields()
      form.setFieldsValue({ config_category: activeCategory })
    }
    setModalVisible(true)
  }

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      // 添加 updated_by（实际应该从登录信息获取）
      values.updated_by = 1
      if (editingRecord) {
        await updateSystemConfig(editingRecord.id, values)
        message.success('更新成功')
      } else {
        await createSystemConfig(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchData(activeCategory, pagination.current, pagination.pageSize)
    } catch (error) {
      if (error.errorFields) {
        return
      }
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除系统配置
  const handleDelete = async (id) => {
    try {
      await deleteSystemConfig(id)
      message.success('删除成功')
      fetchData(activeCategory, pagination.current, pagination.pageSize)
    } catch (error) {
      message.error('删除失败')
    }
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '配置键',
      dataIndex: 'config_key',
      key: 'config_key',
    },
    {
      title: '配置值',
      dataIndex: 'config_value',
      key: 'config_value',
      ellipsis: true,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
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
            title="确定要删除这个配置吗？"
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
        <h2>基础配置</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          新建配置
        </Button>
      </div>

      <Tabs
        activeKey={activeCategory}
        onChange={setActiveCategory}
        items={CATEGORIES.map((cat) => ({
          label: cat.label,
          key: cat.key,
        }))}
      />

      <Table
        columns={columns}
        dataSource={tableData}
        loading={loading}
        rowKey="id"
        pagination={{
          ...pagination,
          onChange: (page, pageSize) => {
            fetchData(activeCategory, page, pageSize)
          },
        }}
      />

      <Modal
        title={editingRecord ? '编辑配置' : '新建配置'}
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
            name="config_category"
            label="配置分类"
            rules={[{ required: true, message: '请选择配置分类' }]}
          >
            <Select disabled={!!editingRecord}>
              {CATEGORIES.map((cat) => (
                <Option key={cat.key} value={cat.key}>
                  {cat.label}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="config_key"
            label="配置键"
            rules={[{ required: true, message: '请输入配置键' }]}
          >
            <Input placeholder="请输入配置键" />
          </Form.Item>
          <Form.Item
            name="config_value"
            label="配置值"
            rules={[{ required: true, message: '请输入配置值' }]}
          >
            <TextArea rows={4} placeholder="请输入配置值" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={2} placeholder="请输入配置描述" />
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

export default SystemConfigManagement
