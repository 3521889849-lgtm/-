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
  Checkbox,
  Tree,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons'
import {
  getRoles,
  getRole,
  createRole,
  updateRole,
  deleteRole,
  getPermissions,
} from '../api/user'

const { TextArea } = Input
const { Option } = Select

const RoleManagement = () => {
  const { message } = App.useApp()
  const [form] = Form.useForm()
  const [tableData, setTableData] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState(null)
  const [permissions, setPermissions] = useState([])
  const [selectedPermissions, setSelectedPermissions] = useState([])
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  // 获取权限列表
  useEffect(() => {
    getPermissions().then((res) => {
      if (res.data.code === 200) {
        const permList = res.data.data?.list || []
        setPermissions(permList)
      }
    })
  }, [])

  // 获取角色列表
  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await getRoles({ page, page_size: pageSize })
      if (res.data.code === 200) {
        setTableData(res.data.data?.list || [])
        setPagination({
          current: res.data.data?.page || 1,
          pageSize: res.data.data?.page_size || 10,
          total: res.data.data?.total || 0,
        })
      }
    } catch (error) {
      message.error('获取角色列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // 打开新增/编辑弹窗
  const handleOpenModal = async (record = null) => {
    setEditingRecord(record)
    if (record) {
      // 获取角色详情（包含权限）
      try {
        const { getRole } = await import('../api/user')
        const res = await getRole(record.id)
        if (res.data.code === 200) {
          const roleData = res.data.data
          form.setFieldsValue({
            role_name: roleData.role_name,
            description: roleData.description,
            status: roleData.status,
          })
          setSelectedPermissions(roleData.permission_ids || [])
        }
      } catch (error) {
        form.setFieldsValue(record)
        setSelectedPermissions([])
      }
    } else {
      form.resetFields()
      setSelectedPermissions([])
    }
    setModalVisible(true)
  }

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const data = {
        ...values,
        permission_ids: selectedPermissions,
      }
      if (editingRecord) {
        await updateRole(editingRecord.id, data)
        message.success('更新成功')
      } else {
        await createRole(data)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      setSelectedPermissions([])
      fetchData(pagination.current, pagination.pageSize)
    } catch (error) {
      if (error.errorFields) {
        return
      }
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除角色
  const handleDelete = async (id) => {
    try {
      await deleteRole(id)
      message.success('删除成功')
      fetchData(pagination.current, pagination.pageSize)
    } catch (error) {
      message.error(error.response?.data?.error || '删除失败')
    }
  }

  // 构建权限树
  const buildPermissionTree = (perms) => {
    return perms.map((perm) => ({
      title: (
        <span>
          {perm.permission_name}
          <span style={{ color: '#999', marginLeft: 8, fontSize: 12 }}>
            ({perm.permission_url})
          </span>
        </span>
      ),
      key: perm.id,
      children: perm.children ? buildPermissionTree(perm.children) : undefined,
    }))
  }

  // 处理权限选择
  const handlePermissionCheck = (checkedKeys) => {
    setSelectedPermissions(checkedKeys)
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '角色名称',
      dataIndex: 'role_name',
      key: 'role_name',
    },
    {
      title: '角色描述',
      dataIndex: 'description',
      key: 'description',
      render: (text) => text || '-',
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
            title="确定要删除这个角色吗？"
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
        <h2>角色管理</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          新建角色
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
        title={editingRecord ? '编辑角色' : '新建角色'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
          setSelectedPermissions([])
        }}
        width={800}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="role_name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>
          <Form.Item name="description" label="角色描述">
            <TextArea rows={3} placeholder="请输入角色描述" />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="ACTIVE">
            <Select>
              <Option value="ACTIVE">启用</Option>
              <Option value="INACTIVE">停用</Option>
            </Select>
          </Form.Item>
          <Form.Item label="业务权限" required>
            <div style={{ border: '1px solid #d9d9d9', borderRadius: 4, padding: 8, maxHeight: 300, overflow: 'auto' }}>
              <Tree
                checkable
                checkedKeys={selectedPermissions}
                onCheck={handlePermissionCheck}
                treeData={buildPermissionTree(permissions)}
              />
            </div>
            <div style={{ marginTop: 8, color: '#999', fontSize: 12 }}>
              双击删除URL,仅支持
            </div>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default RoleManagement
