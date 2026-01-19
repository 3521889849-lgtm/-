import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Input,
  Select,
  DatePicker,
  Space,
  Tag,
  Popconfirm,
  message,
  Modal,
  Form,
  InputNumber,
  Switch,
  Upload,
} from 'antd'
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  LinkOutlined,
  CloudOutlined,
} from '@ant-design/icons'
import { getRoomInfos, updateRoomStatus, deleteRoomInfo, syncRoomStatusToChannel } from '../api/room'
import RoomForm from '../components/RoomForm'
import './RoomStatus.css'

const { RangePicker } = DatePicker

const RoomStatus = () => {
  const [dataSource, setDataSource] = useState([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })
  const [filters, setFilters] = useState({
    status: '',
    keyword: '',
    dateRange: null,
  })
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRoom, setEditingRoom] = useState(null)
  const [syncLoading, setSyncLoading] = useState(false)

  // 获取房源列表
  const fetchRoomList = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const currentBranchId = localStorage.getItem('currentBranchId')
      const params = {
        page,
        page_size: pageSize,
        ...filters,
      }
      // 根据当前分店筛选
      if (currentBranchId) {
        params.branch_id = parseInt(currentBranchId)
      }
      const response = await getRoomInfos(params)
      if (response.data.code === 200) {
        setDataSource(response.data.data.list || [])
        setPagination({
          current: response.data.data.page || 1,
          pageSize: response.data.data.page_size || 10,
          total: response.data.data.total || 0,
        })
      }
    } catch (error) {
      message.error('获取房源列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRoomList(pagination.current, pagination.pageSize)
  }, [filters])

  // 处理状态切换
  const handleStatusChange = async (roomId, newStatus) => {
    try {
      const response = await updateRoomStatus(roomId, { status: newStatus })
      if (response.data.code === 200) {
        message.success('状态更新成功')
        fetchRoomList(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      message.error('状态更新失败')
    }
  }

  // 处理删除
  const handleDelete = async (roomId) => {
    try {
      const response = await deleteRoomInfo(roomId)
      if (response.data.code === 200) {
        message.success('删除成功')
        fetchRoomList(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      message.error('删除失败')
    }
  }

  // 处理搜索
  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    fetchRoomList(1, pagination.pageSize)
  }

  // 处理分页
  const handleTableChange = (newPagination) => {
    fetchRoomList(newPagination.current, newPagination.pageSize)
  }

  const columns = [
    {
      title: '房源名称',
      dataIndex: 'room_name',
      key: 'room_name',
      width: 200,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => {
        const statusMap = {
          ACTIVE: { color: 'green', text: '启用' },
          INACTIVE: { color: 'red', text: '停用' },
          MAINTENANCE: { color: 'orange', text: '维修' },
        }
        const statusInfo = statusMap[status] || { color: 'default', text: status }
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>
      },
    },
    {
      title: '售卖周期',
      dataIndex: 'sales_period',
      key: 'sales_period',
      width: 100,
      render: () => '天',
    },
    {
      title: '售卖方式',
      dataIndex: 'sales_method',
      key: 'sales_method',
      width: 100,
      render: () => '间',
    },
    {
      title: '操作',
      key: 'action',
      width: 300,
      render: (_, record) => (
        <Space size="small">
          <Button type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            查看
          </Button>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            修改
          </Button>
          <Popconfirm
            title="确定要删除这个房源吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
          {record.status === 'ACTIVE' ? (
            <Button
              type="link"
              onClick={() => handleStatusChange(record.id, 'INACTIVE')}
            >
              停用
            </Button>
          ) : (
            <Button
              type="link"
              onClick={() => handleStatusChange(record.id, 'ACTIVE')}
            >
              启用
            </Button>
          )}
          <Button type="link" icon={<LinkOutlined />} onClick={() => handleBind(record)}>
            关联房源
          </Button>
        </Space>
      ),
    },
  ]

  const handleView = (record) => {
    setEditingRoom(record)
    setModalVisible(true)
  }

  const handleEdit = (record) => {
    setEditingRoom(record)
    setModalVisible(true)
  }

  const handleBind = (record) => {
    // TODO: 实现关联房源功能
    message.info('关联房源功能开发中')
  }

  const handleAdd = () => {
    setEditingRoom(null)
    setModalVisible(true)
  }

  const handleModalClose = () => {
    setModalVisible(false)
    setEditingRoom(null)
    fetchRoomList(pagination.current, pagination.pageSize)
  }

  // 同步房态到渠道
  const handleSyncToChannel = async () => {
    const currentBranchId = localStorage.getItem('currentBranchId')
    if (!currentBranchId) {
      message.warning('请先选择分店')
      return
    }

    setSyncLoading(true)
    try {
      // 默认同步到途游渠道（channel_id=1），可以根据实际情况调整
      const syncData = {
        branch_id: parseInt(currentBranchId),
        channel_id: 1, // 途游渠道ID，需要根据实际渠道配置调整
        start_date: new Date().toISOString().split('T')[0],
        end_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0], // 未来7天
      }

      const response = await syncRoomStatusToChannel(syncData)
      if (response.data.code === 200) {
        const result = response.data.data
        message.success(
          `同步完成！成功: ${result.success_count}，失败: ${result.fail_count}`
        )
      } else {
        message.error(response.data.msg || '同步失败')
      }
    } catch (error) {
      console.error('同步失败:', error)
      message.error('同步失败: ' + (error.response?.data?.msg || error.message))
    } finally {
      setSyncLoading(false)
    }
  }

  return (
    <div className="room-status-page">
      <div className="page-header">
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          添加房源
        </Button>
      </div>

      <div className="search-section">
        <Space size="middle" wrap>
          <span>状态:</span>
          <Select
            style={{ width: 150 }}
            placeholder="全部"
            allowClear
            value={filters.status || undefined}
            onChange={(value) => setFilters({ ...filters, status: value || '' })}
          >
            <Select.Option value="ACTIVE">启用</Select.Option>
            <Select.Option value="INACTIVE">停用</Select.Option>
            <Select.Option value="MAINTENANCE">维修</Select.Option>
          </Select>

          <span>房源名称:</span>
          <Input
            style={{ width: 200 }}
            placeholder="请输入房源名称"
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            onPressEnter={handleSearch}
          />

          <span>添加时间:</span>
          <RangePicker
            onChange={(dates) =>
              setFilters({
                ...filters,
                dateRange: dates
                  ? [dates[0].format('YYYY-MM-DD'), dates[1].format('YYYY-MM-DD')]
                  : null,
              })
            }
          />

          <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
            查询
          </Button>
        </Space>
      </div>

      <div className="sync-section">
        <Button
          type="primary"
          icon={<CloudOutlined />}
          onClick={handleSyncToChannel}
          loading={syncLoading}
        >
          同步到途游
        </Button>
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
          showQuickJumper: true,
        }}
        onChange={handleTableChange}
      />

      <RoomForm
        visible={modalVisible}
        room={editingRoom}
        onCancel={handleModalClose}
        onSuccess={handleModalClose}
      />
    </div>
  )
}

export default RoomStatus
