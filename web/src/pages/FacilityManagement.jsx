import React, { useState, useEffect } from 'react'
import { Table, Button, Input, Select, Space, Tag, Popconfirm, Modal, Form, App } from 'antd'
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getFacilities, createFacility, updateFacility, deleteFacility } from '../api/room'

const FacilityManagement = () => {
  const { message } = App.useApp()
  const [dataSource, setDataSource] = useState([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })
  const [filters, setFilters] = useState({ status: '', keyword: '' })
  const [modalVisible, setModalVisible] = useState(false)
  const [editingFacility, setEditingFacility] = useState(null)
  const [form] = Form.useForm()

  const fetchFacilities = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const params = { page, page_size: pageSize, ...filters }
      const response = await getFacilities(params)
      if (response.data.code === 200) {
        setDataSource(response.data.data.list || [])
        setPagination({
          current: response.data.data.page || 1,
          pageSize: response.data.data.page_size || 10,
          total: response.data.data.total || 0,
        })
      }
    } catch (error) {
      message.error('获取设施列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchFacilities(pagination.current, pagination.pageSize)
  }, [filters])

  const handleAdd = () => {
    setEditingFacility(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record) => {
    setEditingFacility(record)
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleDelete = async (id) => {
    try {
      const response = await deleteFacility(id)
      if (response.data.code === 200) {
        message.success('删除成功')
        fetchFacilities(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      message.error('删除失败')
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingFacility) {
        await updateFacility(editingFacility.id, values)
        message.success('更新成功')
      } else {
        await createFacility(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchFacilities(pagination.current, pagination.pageSize)
    } catch (error) {
      message.error(editingFacility ? '更新失败' : '创建失败')
    }
  }

  const columns = [
    { title: '设施名称', dataIndex: 'facility_name', key: 'facility_name' },
    { title: '设施描述', dataIndex: 'description', key: 'description' },
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
      render: (_, record) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            修改
          </Button>
          <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          添加设施
        </Button>
        <Space>
          <Select
            style={{ width: 150 }}
            placeholder="状态筛选"
            allowClear
            value={filters.status || undefined}
            onChange={(value) => setFilters({ ...filters, status: value || '' })}
          >
            <Select.Option value="ACTIVE">启用</Select.Option>
            <Select.Option value="INACTIVE">停用</Select.Option>
          </Select>
          <Input
            style={{ width: 200 }}
            placeholder="搜索设施名称"
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            onPressEnter={() => fetchFacilities(1, pagination.pageSize)}
          />
          <Button type="primary" icon={<SearchOutlined />} onClick={() => fetchFacilities(1, pagination.pageSize)}>
            查询
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={dataSource}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
        onChange={(newPagination) => fetchFacilities(newPagination.current, newPagination.pageSize)}
      />

      <Modal
        title={editingFacility ? '修改设施' : '添加设施'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={handleSubmit}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="facility_name" label="设施名称" rules={[{ required: true }]}>
            <Input placeholder="如：无线wifi、空调、冰箱" />
          </Form.Item>
          <Form.Item name="description" label="设施描述">
            <Input.TextArea rows={3} placeholder="请输入设施描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default FacilityManagement
