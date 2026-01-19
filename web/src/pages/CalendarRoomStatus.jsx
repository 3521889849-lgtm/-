import React, { useState, useEffect } from 'react'
import { Card, DatePicker, Select, Table, Tag, Button, Space, Input, Modal, Form, message } from 'antd'
import { CalendarOutlined, ReloadOutlined, EditOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import axios from 'axios'

const { RangePicker } = DatePicker
const { Option } = Select

const API_BASE = '/api/v1'

// 房态颜色映射
const statusColorMap = {
  空净房: 'green',
  入住房: 'blue',
  维修房: 'orange',
  锁定房: 'red',
  空账房: 'purple',
  预定房: 'cyan',
}

const CalendarRoomStatus = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
  const [dateRange, setDateRange] = useState([dayjs(), dayjs().add(6, 'day')])
  const [roomNo, setRoomNo] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [editModalVisible, setEditModalVisible] = useState(false)
  const [editingItem, setEditingItem] = useState(null)
  const [form] = Form.useForm()
  const [currentBranchId, setCurrentBranchId] = useState(null)

  // 从 localStorage 获取当前分店ID
  useEffect(() => {
    const branchId = localStorage.getItem('currentBranchId')
    if (branchId) {
      setCurrentBranchId(parseInt(branchId))
    }
  }, [])

  // 获取日历化房态数据
  const fetchData = async () => {
    setLoading(true)
    try {
      const params = {
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      }
      if (currentBranchId) params.branch_id = currentBranchId
      if (roomNo) params.room_no = roomNo
      if (statusFilter) params.status = statusFilter

      const response = await axios.get(`${API_BASE}/calendar-room-status`, { params })
      if (response.data.code === 200) {
        setData(response.data.data || [])
      } else {
        message.error(response.data.msg || '获取数据失败')
      }
    } catch (error) {
      console.error('获取日历化房态失败:', error)
      message.error('获取数据失败: ' + (error.response?.data?.msg || error.message))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [dateRange, roomNo, statusFilter, currentBranchId])

  // 监听分店切换
  useEffect(() => {
    const handleStorageChange = () => {
      const branchId = localStorage.getItem('currentBranchId')
      if (branchId) {
        setCurrentBranchId(parseInt(branchId))
      }
    }
    window.addEventListener('storage', handleStorageChange)
    // 也监听自定义事件（同页面切换分店时）
    window.addEventListener('branchChanged', handleStorageChange)
    return () => {
      window.removeEventListener('storage', handleStorageChange)
      window.removeEventListener('branchChanged', handleStorageChange)
    }
  }, [])

  // 处理日期范围变化
  const handleDateRangeChange = (dates) => {
    if (dates && dates.length === 2) {
      setDateRange(dates)
    }
  }

  // 处理编辑
  const handleEdit = (record) => {
    setEditingItem(record)
    form.setFieldsValue({
      room_id: record.room_id,
      date: dayjs(record.date),
      status: record.room_status,
    })
    setEditModalVisible(true)
  }

  // 保存编辑
  const handleSave = async () => {
    try {
      const values = await form.validateFields()
      const updateData = {
        room_id: values.room_id,
        date: values.date.format('YYYY-MM-DD'),
        status: values.status,
      }

      const response = await axios.put(`${API_BASE}/calendar-room-status`, updateData)
      if (response.data.code === 200) {
        message.success('更新成功')
        setEditModalVisible(false)
        fetchData()
      } else {
        message.error(response.data.msg || '更新失败')
      }
    } catch (error) {
      if (error.errorFields) {
        return
      }
      console.error('更新房态失败:', error)
      message.error('更新失败: ' + (error.response?.data?.msg || error.message))
    }
  }

  // 将数据转换为表格格式（按房间分组，按日期展示）
  const transformData = () => {
    const roomMap = new Map()

    // 按房间分组
    data.forEach((item) => {
      const key = `${item.room_id}_${item.room_no}`
      if (!roomMap.has(key)) {
        roomMap.set(key, {
          room_id: item.room_id,
          room_no: item.room_no,
          room_name: item.room_name,
          dates: {},
        })
      }
      const room = roomMap.get(key)
      room.dates[item.date] = item
    })

    // 生成日期列
    const dates = []
    let current = dateRange[0]
    while (current.isBefore(dateRange[1]) || current.isSame(dateRange[1], 'day')) {
      dates.push(current.format('YYYY-MM-DD'))
      current = current.add(1, 'day')
    }

    // 转换为表格数据
    const tableData = Array.from(roomMap.values()).map((room) => {
      const row = {
        key: room.room_id,
        room_id: room.room_id,
        room_no: room.room_no,
        room_name: room.room_name,
      }

      dates.forEach((date) => {
        if (room.dates[date]) {
          row[date] = room.dates[date]
        } else {
          row[date] = null
        }
      })

      return row
    })

    return { tableData, dates }
  }

  const { tableData, dates } = transformData()

  // 表格列定义
  const columns = [
    {
      title: '房间号',
      dataIndex: 'room_no',
      key: 'room_no',
      fixed: 'left',
      width: 100,
      render: (text, record) => (
        <div>
          <div style={{ fontWeight: 'bold' }}>{text}</div>
          <div style={{ fontSize: '12px', color: '#999' }}>{record.room_name}</div>
        </div>
      ),
    },
    ...dates.map((date) => ({
      title: (
        <div>
          <div>{dayjs(date).format('MM-DD')}</div>
          <div style={{ fontSize: '11px', color: '#999' }}>{dayjs(date).format('ddd')}</div>
        </div>
      ),
      dataIndex: date,
      key: date,
      width: 120,
      render: (item) => {
        if (!item) {
          return <Tag color="default">-</Tag>
        }
        return (
          <div>
            <Tag color={statusColorMap[item.room_status] || 'default'} style={{ marginBottom: 4 }}>
              {item.room_status}
            </Tag>
            <div style={{ fontSize: '11px', color: '#666', marginTop: 4 }}>
              {item.checked_in_count > 0 && <div>入住: {item.checked_in_count}</div>}
              {item.reserved_pending_count > 0 && <div>预定: {item.reserved_pending_count}</div>}
              {item.check_out_pending_count > 0 && <div>预退: {item.check_out_pending_count}</div>}
            </div>
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(item)}
              style={{ padding: 0, height: 'auto' }}
            >
              编辑
            </Button>
          </div>
        )
      },
    })),
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={
          <Space>
            <CalendarOutlined />
            <span>日历化房态展示</span>
          </Space>
        }
        extra={
          <Button icon={<ReloadOutlined />} onClick={fetchData} loading={loading}>
            刷新
          </Button>
        }
      >
        {/* 筛选条件 */}
        <Space style={{ marginBottom: 16 }} wrap>
          <span>日期范围：</span>
          <RangePicker
            value={dateRange}
            onChange={handleDateRangeChange}
            format="YYYY-MM-DD"
            allowClear={false}
          />

          <span>房间号：</span>
          <Input
            placeholder="输入房间号筛选"
            value={roomNo}
            onChange={(e) => setRoomNo(e.target.value)}
            style={{ width: 150 }}
            allowClear
          />

          <span>房态筛选：</span>
          <Select
            value={statusFilter}
            onChange={setStatusFilter}
            style={{ width: 120 }}
            allowClear
            placeholder="全部"
          >
            {Object.keys(statusColorMap).map((status) => (
              <Option key={status} value={status}>
                <Tag color={statusColorMap[status]}>{status}</Tag>
              </Option>
            ))}
          </Select>
        </Space>

        {/* 房态图例 */}
        <div style={{ marginBottom: 16, padding: '8px 12px', background: '#f5f5f5', borderRadius: 4 }}>
          <span style={{ marginRight: 16, fontWeight: 'bold' }}>房态说明：</span>
          {Object.entries(statusColorMap).map(([status, color]) => (
            <Tag key={status} color={color} style={{ marginRight: 8 }}>
              {status}
            </Tag>
          ))}
        </div>

        {/* 表格 */}
        <Table
          columns={columns}
          dataSource={tableData}
          loading={loading}
          scroll={{ x: 'max-content' }}
          pagination={{
            pageSize: 20,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 个房间`,
          }}
        />
      </Card>

      {/* 编辑模态框 */}
      <Modal
        title="编辑房态"
        open={editModalVisible}
        onOk={handleSave}
        onCancel={() => {
          setEditModalVisible(false)
          form.resetFields()
        }}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item label="房间号" name="room_id" hidden>
            <Input disabled />
          </Form.Item>
          <Form.Item label="日期" name="date">
            <DatePicker format="YYYY-MM-DD" disabled style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            label="房态"
            name="status"
            rules={[{ required: true, message: '请选择房态' }]}
          >
            <Select placeholder="请选择房态">
              {Object.keys(statusColorMap).map((status) => (
                <Option key={status} value={status}>
                  <Tag color={statusColorMap[status]}>{status}</Tag>
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CalendarRoomStatus
