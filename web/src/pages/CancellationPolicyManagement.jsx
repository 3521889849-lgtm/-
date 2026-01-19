import React, { useState, useEffect } from 'react'
import { Table, Button, Input, Select, Space, Tag, Popconfirm, Modal, Form, InputNumber, App } from 'antd'
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getCancellationPolicies, createCancellationPolicy, updateCancellationPolicy, deleteCancellationPolicy, getRoomTypes } from '../api/room'

const CancellationPolicyManagement = () => {
  const { message } = App.useApp()
  const [dataSource, setDataSource] = useState([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })
  const [filters, setFilters] = useState({ status: '', keyword: '', room_type_id: '' })
  const [modalVisible, setModalVisible] = useState(false)
  const [editingPolicy, setEditingPolicy] = useState(null)
  const [roomTypes, setRoomTypes] = useState([])
  const [form] = Form.useForm()

  const fetchPolicies = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const params = { page, page_size: pageSize, ...filters }
      const response = await getCancellationPolicies(params)
      if (response.data.code === 200) {
        setDataSource(response.data.data.list || [])
        setPagination({
          current: response.data.data.page || 1,
          pageSize: response.data.data.page_size || 10,
          total: response.data.data.total || 0,
        })
      }
    } catch (error) {
      message.error('获取退订政策列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchPolicies(pagination.current, pagination.pageSize)
    loadRoomTypes()
  }, [filters])

  const loadRoomTypes = async () => {
    try {
      const response = await getRoomTypes({ page: 1, page_size: 100 })
      if (response.data.code === 200) {
        setRoomTypes(response.data.data.list || [])
      }
    } catch (error) {
      console.error('加载房型失败:', error)
    }
  }

  const handleAdd = () => {
    setEditingPolicy(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record) => {
    setEditingPolicy(record)
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleDelete = async (id) => {
    try {
      const response = await deleteCancellationPolicy(id)
      if (response.data.code === 200) {
        message.success('删除成功')
        fetchPolicies(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      message.error('删除失败')
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingPolicy) {
        await updateCancellationPolicy(editingPolicy.id, values)
        message.success('更新成功')
      } else {
        await createCancellationPolicy(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchPolicies(pagination.current, pagination.pageSize)
    } catch (error) {
      message.error(editingPolicy ? '更新失败' : '创建失败')
    }
  }

  const columns = [
    { title: '政策名称', dataIndex: 'policy_name', key: 'policy_name' },
    { title: '规则描述', dataIndex: 'rule_description', key: 'rule_description' },
    { title: '违约金比例', dataIndex: 'penalty_ratio', key: 'penalty_ratio', render: (val) => `${val}倍房费` },
    {
      title: '适用房型',
      dataIndex: 'room_type_id',
      key: 'room_type_id',
      render: (id) => {
        const roomType = roomTypes.find((t) => t.id === id)
        return roomType ? roomType.room_type_name : '全部'
      },
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
          添加退订政策
        </Button>
        <Space>
          <Select
            style={{ width: 150 }}
            placeholder="房型筛选"
            allowClear
            value={filters.room_type_id || undefined}
            onChange={(value) => setFilters({ ...filters, room_type_id: value || '' })}
          >
            {roomTypes.map((type) => (
              <Select.Option key={type.id} value={type.id}>
                {type.room_type_name}
              </Select.Option>
            ))}
          </Select>
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
            placeholder="搜索政策名称"
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            onPressEnter={() => fetchPolicies(1, pagination.pageSize)}
          />
          <Button type="primary" icon={<SearchOutlined />} onClick={() => fetchPolicies(1, pagination.pageSize)}>
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
        onChange={(newPagination) => fetchPolicies(newPagination.current, newPagination.pageSize)}
      />

      <Modal
        title={editingPolicy ? '修改退订政策' : '添加退订政策'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={handleSubmit}
        width={700}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="policy_name" label="政策名称" rules={[{ required: true }]}>
            <Input placeholder="如：入住前24小时不可取消" />
          </Form.Item>
          <Form.Item name="rule_description" label="规则描述" rules={[{ required: true }]}>
            <Input.TextArea
              rows={3}
              placeholder="如：入住前X小时内不可取消，否则收取X倍房费"
            />
          </Form.Item>
          <Form.Item name="penalty_ratio" label="违约金比例（倍房费）" rules={[{ required: true }]}>
            <InputNumber min={0} max={10} precision={2} style={{ width: '100%' }} placeholder="如：1.5" />
          </Form.Item>
          <Form.Item name="room_type_id" label="适用房型（可选）">
            <Select placeholder="选择房型，留空则适用于所有房型" allowClear>
              {roomTypes.map((type) => (
                <Select.Option key={type.id} value={type.id}>
                  {type.room_type_name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CancellationPolicyManagement
