import React, { useState, useEffect } from 'react'
import { Table, Button, Input, Select, Space, Tag, Popconfirm, Modal, Form, InputNumber, Switch, message } from 'antd'
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getRoomTypes, createRoomType, updateRoomType, deleteRoomType } from '../api/room'

const RoomTypeManagement = () => {
  const [dataSource, setDataSource] = useState([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })
  const [filters, setFilters] = useState({ status: '', keyword: '' })
  const [modalVisible, setModalVisible] = useState(false)
  const [editingType, setEditingType] = useState(null)
  const [form] = Form.useForm()

  const fetchRoomTypes = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const params = { page, page_size: pageSize, ...filters }
      const response = await getRoomTypes(params)
      if (response.data.code === 200) {
        setDataSource(response.data.data.list || [])
        setPagination({
          current: response.data.data.page || 1,
          pageSize: response.data.data.page_size || 10,
          total: response.data.data.total || 0,
        })
      }
    } catch (error) {
      message.error('获取房型列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRoomTypes(pagination.current, pagination.pageSize)
  }, [filters])

  const handleAdd = () => {
    setEditingType(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record) => {
    setEditingType(record)
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleDelete = async (id) => {
    try {
      const response = await deleteRoomType(id)
      if (response.data.code === 200) {
        message.success('删除成功')
        fetchRoomTypes(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      message.error('删除失败')
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingType) {
        await updateRoomType(editingType.id, values)
        message.success('更新成功')
      } else {
        await createRoomType(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchRoomTypes(pagination.current, pagination.pageSize)
    } catch (error) {
      message.error(editingType ? '更新失败' : '创建失败')
    }
  }

  const columns = [
    { title: '房型名称', dataIndex: 'room_type_name', key: 'room_type_name' },
    { title: '床型规格', dataIndex: 'bed_spec', key: 'bed_spec' },
    { title: '面积', dataIndex: 'area', key: 'area', render: (val) => val ? `${val}㎡` : '-' },
    { title: '含早', dataIndex: 'has_breakfast', key: 'has_breakfast', render: (val) => val ? '是' : '否' },
    { title: '洗漱用品', dataIndex: 'has_toiletries', key: 'has_toiletries', render: (val) => val ? '是' : '否' },
    { title: '默认价格', dataIndex: 'default_price', key: 'default_price', render: (val) => `¥${val}` },
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
          添加房型
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
            placeholder="搜索房型名称"
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            onPressEnter={() => fetchRoomTypes(1, pagination.pageSize)}
          />
          <Button type="primary" icon={<SearchOutlined />} onClick={() => fetchRoomTypes(1, pagination.pageSize)}>
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
        onChange={(newPagination) => fetchRoomTypes(newPagination.current, newPagination.pageSize)}
      />

      <Modal
        title={editingType ? '修改房型' : '添加房型'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={handleSubmit}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="room_type_name" label="房型名称" rules={[{ required: true }]}>
            <Input placeholder="如：大床房、标准间" />
          </Form.Item>
          <Form.Item name="bed_spec" label="床型规格" rules={[{ required: true }]}>
            <Input placeholder="如：1.8*2.0m" />
          </Form.Item>
          <Form.Item name="area" label="面积（平方米）">
            <InputNumber min={0} precision={2} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="has_breakfast" label="是否含早" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="has_toiletries" label="是否提供洗漱用品" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="default_price" label="默认门市价" rules={[{ required: true }]}>
            <InputNumber min={0} precision={2} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default RoomTypeManagement
