import React, { useState, useEffect } from 'react'
import { Card, Row, Col, Statistic, DatePicker, Select, Input, Button, Table, Tag, Space, message } from 'antd'
import { ReloadOutlined, HomeOutlined, UserOutlined, LogoutOutlined, CalendarOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import axios from 'axios'

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

const RealTimeStatistics = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState(null)
  const [date, setDate] = useState(dayjs())
  const [roomNo, setRoomNo] = useState('')
  const [roomTypeId, setRoomTypeId] = useState('')
  const [roomTypes, setRoomTypes] = useState([])
  const [currentBranchId, setCurrentBranchId] = useState(null)

  // 从 localStorage 获取当前分店ID
  useEffect(() => {
    const branchId = localStorage.getItem('currentBranchId')
    if (branchId) {
      setCurrentBranchId(parseInt(branchId))
    }
  }, [])

  // 获取房型列表
  useEffect(() => {
    axios.get(`${API_BASE}/room-types?page_size=100`).then((res) => {
      if (res.data.code === 200) {
        setRoomTypes(res.data.data?.list || [])
      }
    })
  }, [])

  // 获取实时统计数据
  const fetchData = async () => {
    setLoading(true)
    try {
      const params = {
        date: date.format('YYYY-MM-DD'),
      }
      if (currentBranchId) params.branch_id = currentBranchId
      if (roomNo) params.room_no = roomNo
      if (roomTypeId) params.room_type_id = roomTypeId

      const response = await axios.get(`${API_BASE}/real-time-statistics`, { params })
      if (response.data.code === 200) {
        setData(response.data.data)
      } else {
        message.error(response.data.msg || '获取数据失败')
      }
    } catch (error) {
      console.error('获取实时统计数据失败:', error)
      message.error('获取数据失败: ' + (error.response?.data?.msg || error.message))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [date, roomNo, roomTypeId, currentBranchId])

  // 监听分店切换
  useEffect(() => {
    const handleStorageChange = () => {
      const branchId = localStorage.getItem('currentBranchId')
      if (branchId) {
        setCurrentBranchId(parseInt(branchId))
      }
    }
    window.addEventListener('storage', handleStorageChange)
    window.addEventListener('branchChanged', handleStorageChange)
    return () => {
      window.removeEventListener('storage', handleStorageChange)
      window.removeEventListener('branchChanged', handleStorageChange)
    }
  }, [])

  // 表格列定义
  const columns = [
    {
      title: '房间号',
      dataIndex: 'room_no',
      key: 'room_no',
      width: 100,
    },
    {
      title: '房间名称',
      dataIndex: 'room_name',
      key: 'room_name',
      width: 150,
    },
    {
      title: '房态',
      dataIndex: 'room_status',
      key: 'room_status',
      width: 100,
      render: (status) => (
        <Tag color={statusColorMap[status] || 'default'}>{status}</Tag>
      ),
    },
    {
      title: '剩余数量',
      dataIndex: 'remaining_count',
      key: 'remaining_count',
      width: 100,
      align: 'right',
    },
    {
      title: '已入住人数',
      dataIndex: 'checked_in_count',
      key: 'checked_in_count',
      width: 120,
      align: 'right',
    },
    {
      title: '预退房人数',
      dataIndex: 'check_out_pending_count',
      key: 'check_out_pending_count',
      width: 120,
      align: 'right',
    },
    {
      title: '预定待入住',
      dataIndex: 'reserved_pending_count',
      key: 'reserved_pending_count',
      width: 120,
      align: 'right',
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={
          <Space>
            <HomeOutlined />
            <span>实时数据统计</span>
          </Space>
        }
        extra={
          <Button icon={<ReloadOutlined />} onClick={fetchData} loading={loading}>
            刷新
          </Button>
        }
      >
        {/* 筛选条件 */}
        <Space style={{ marginBottom: 24 }} wrap>
          <span>统计日期：</span>
          <DatePicker
            value={date}
            onChange={setDate}
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

          <span>房型：</span>
          <Select
            value={roomTypeId}
            onChange={setRoomTypeId}
            style={{ width: 150 }}
            allowClear
            placeholder="全部房型"
          >
            {roomTypes.map((type) => (
              <Option key={type.id} value={type.id}>
                {type.room_type_name}
              </Option>
            ))}
          </Select>
        </Space>

        {/* 核心数据卡片 */}
        {data && (
          <>
            <Row gutter={16} style={{ marginBottom: 24 }}>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="总房间数"
                    value={data.total_rooms}
                    prefix={<HomeOutlined />}
                    valueStyle={{ color: '#1890ff' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="剩余房间数"
                    value={data.remaining_rooms}
                    prefix={<HomeOutlined />}
                    valueStyle={{ color: '#52c41a' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="已入住人数"
                    value={data.checked_in_count}
                    prefix={<UserOutlined />}
                    valueStyle={{ color: '#1890ff' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="预退房人数"
                    value={data.check_out_pending_count}
                    prefix={<LogoutOutlined />}
                    valueStyle={{ color: '#faad14' }}
                  />
                </Card>
              </Col>
            </Row>

            <Row gutter={16} style={{ marginBottom: 24 }}>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="预定待入住"
                    value={data.reserved_pending_count}
                    prefix={<CalendarOutlined />}
                    valueStyle={{ color: '#722ed1' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="已入住房间"
                    value={data.occupied_rooms}
                    valueStyle={{ color: '#1890ff' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="维修房间"
                    value={data.maintenance_rooms}
                    valueStyle={{ color: '#fa8c16' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="锁定房间"
                    value={data.locked_rooms}
                    valueStyle={{ color: '#f5222d' }}
                  />
                </Card>
              </Col>
            </Row>

            <Row gutter={16} style={{ marginBottom: 24 }}>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="空净房间"
                    value={data.empty_rooms}
                    valueStyle={{ color: '#52c41a' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="预定房间"
                    value={data.reserved_rooms}
                    valueStyle={{ color: '#13c2c2' }}
                  />
                </Card>
              </Col>
            </Row>

            {/* 房态分组统计 */}
            {data.status_breakdown && data.status_breakdown.length > 0 && (
              <Card title="房态分布" style={{ marginBottom: 24 }}>
                <Row gutter={16}>
                  {data.status_breakdown.map((item) => (
                    <Col key={item.status} span={6} style={{ marginBottom: 16 }}>
                      <Card size="small">
                        <Statistic
                          title={
                            <Tag color={statusColorMap[item.status] || 'default'}>
                              {item.status}
                            </Tag>
                          }
                          value={item.count}
                        />
                      </Card>
                    </Col>
                  ))}
                </Row>
              </Card>
            )}

            {/* 房间明细表格 */}
            {data.room_details && data.room_details.length > 0 && (
              <Card title="房间明细">
                <Table
                  columns={columns}
                  dataSource={data.room_details.map((item, index) => ({
                    ...item,
                    key: item.room_id || index,
                  }))}
                  loading={loading}
                  pagination={{
                    pageSize: 20,
                    showSizeChanger: true,
                    showTotal: (total) => `共 ${total} 个房间`,
                  }}
                />
              </Card>
            )}
          </>
        )}
      </Card>
    </div>
  )
}

export default RealTimeStatistics
