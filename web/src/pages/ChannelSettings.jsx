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
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons'

const { Option } = Select
import {
  getChannelConfigs,
  createChannelConfig,
  updateChannelConfig,
  deleteChannelConfig,
} from '../api/user'

const ChannelSettings = () => {
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

  // 获取渠道配置列表
  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await getChannelConfigs({ page, page_size: pageSize })
      if (res.data.code === 200) {
        setTableData(res.data.data?.list || [])
        setPagination({
          current: res.data.data?.page || 1,
          pageSize: res.data.data?.page_size || 10,
          total: res.data.data?.total || 0,
        })
      }
    } catch (error) {
      message.error('获取渠道配置列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // 打开新增/编辑弹窗
  const handleOpenModal = (record = null) => {
    setEditingRecord(record)
    if (record) {
      form.setFieldsValue(record)
    } else {
      form.resetFields()
    }
    setModalVisible(true)
  }

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingRecord) {
        await updateChannelConfig(editingRecord.id, values)
        message.success('更新成功')
      } else {
        await createChannelConfig(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchData(pagination.current, pagination.pageSize)
    } catch (error) {
      if (error.errorFields) {
        return
      }
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除渠道配置
  const handleDelete = async (id) => {
    try {
      await deleteChannelConfig(id)
      message.success('删除成功')
      fetchData(pagination.current, pagination.pageSize)
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
      title: '渠道名称',
      dataIndex: 'channel_name',
      key: 'channel_name',
    },
    {
      title: '渠道编码',
      dataIndex: 'channel_code',
      key: 'channel_code',
    },
    {
      title: '接口URL',
      dataIndex: 'api_url',
      key: 'api_url',
      ellipsis: true,
    },
    {
      title: '同步规则',
      dataIndex: 'sync_rule',
      key: 'sync_rule',
      render: (text) => (
        <Tag color={text === 'REALTIME' ? 'green' : 'orange'}>
          {text === 'REALTIME' ? '实时' : '定时'}
        </Tag>
      ),
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
            title="确定要删除这个渠道配置吗？"
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
        <h2>渠道设置</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          新建渠道
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={tableData}
        loading={loading}
        rowKey="id"
        pagination={{
          ...pagination,
          onChange: (page, pageSize) => {
            fetchData(page, pageSize)
          },
        }}
      />

      <Modal
        title={editingRecord ? '编辑渠道配置' : '新建渠道配置'}
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
            name="channel_name"
            label="渠道名称"
            rules={[{ required: true, message: '请输入渠道名称' }]}
          >
            <Input placeholder="如：携程、途游、艺龙等" />
          </Form.Item>
          <Form.Item
            name="channel_code"
            label="渠道编码"
            rules={[{ required: true, message: '请输入渠道编码' }]}
          >
            <Input placeholder="请输入渠道编码" />
          </Form.Item>
          <Form.Item
            name="api_url"
            label="对接接口URL"
            rules={[{ required: true, message: '请输入接口URL' }]}
          >
            <Input placeholder="请输入对接接口URL" />
          </Form.Item>
          <Form.Item
            name="sync_rule"
            label="同步规则"
            initialValue="REALTIME"
          >
            <Select>
              <Option value="REALTIME">实时</Option>
              <Option value="SCHEDULED">定时</Option>
            </Select>
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

export default ChannelSettings
