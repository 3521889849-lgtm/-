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
  App,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons'

const { Option } = Select
import {
  getBlacklists,
  createBlacklist,
  updateBlacklist,
  deleteBlacklist,
} from '../api/user'

const { TextArea } = Input

const BlacklistManagement = () => {
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

  // 获取黑名单列表
  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await getBlacklists({ page, page_size: pageSize })
      if (res.data.code === 200) {
        setTableData(res.data.data?.list || [])
        setPagination({
          current: res.data.data?.page || 1,
          pageSize: res.data.data?.page_size || 10,
          total: res.data.data?.total || 0,
        })
      }
    } catch (error) {
      message.error('获取黑名单列表失败')
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
      // 添加 operator_id（实际应该从登录信息获取）
      values.operator_id = 1
      if (editingRecord) {
        await updateBlacklist(editingRecord.id, values)
        message.success('更新成功')
      } else {
        await createBlacklist(values)
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

  // 删除黑名单
  const handleDelete = async (id) => {
    try {
      await deleteBlacklist(id)
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
      title: '客人姓名',
      dataIndex: 'guest_name',
      key: 'guest_name',
      render: (text) => text || '-',
    },
    {
      title: '证件号',
      dataIndex: 'id_number',
      key: 'id_number',
      ellipsis: true,
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '拉黑原因',
      dataIndex: 'reason',
      key: 'reason',
      ellipsis: true,
    },
    {
      title: '拉黑时间',
      dataIndex: 'black_time',
      key: 'black_time',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 'VALID' ? 'red' : 'default'}>
          {status === 'VALID' ? '有效' : '无效'}
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
            title="确定要删除这条黑名单记录吗？"
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
        <h2>黑名单管理</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          添加黑名单
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
        title={editingRecord ? '编辑黑名单' : '添加黑名单'}
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
            name="id_number"
            label="证件号"
            rules={[{ required: true, message: '请输入证件号' }]}
          >
            <Input placeholder="请输入证件号" />
          </Form.Item>
          <Form.Item
            name="phone"
            label="手机号"
            rules={[{ required: true, message: '请输入手机号' }]}
          >
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item
            name="reason"
            label="拉黑原因"
            rules={[{ required: true, message: '请输入拉黑原因' }]}
          >
            <TextArea rows={3} placeholder="请输入拉黑原因" />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="VALID">
            <Select>
              <Option value="VALID">有效</Option>
              <Option value="INVALID">无效</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default BlacklistManagement
