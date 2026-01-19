import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Input,
  Select,
  DatePicker,
  Space,
  Tag,
  message,
  Modal,
  Descriptions,
  Card,
} from 'antd'
import { SearchOutlined, EyeOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { getOrders, getOrder } from '../api/room'

const { Option } = Select
const { RangePicker } = DatePicker

// 订单状态映射
const orderStatusMap = {
  RESERVED: { color: 'blue', text: '已预定' },
  CHECKED_IN: { color: 'green', text: '已入住' },
  CHECKED_OUT: { color: 'default', text: '已退房' },
  CANCELLED: { color: 'red', text: '已失效' },
  已预定: { color: 'blue', text: '已预定' },
  已入住: { color: 'green', text: '已入住' },
  已退房: { color: 'default', text: '已退房' },
  已失效: { color: 'red', text: '已失效' },
}

// 客人来源选项
const guestSourceOptions = [
  { value: '散客', label: '散客' },
  { value: '携程', label: '携程' },
  { value: '途游', label: '途游' },
  { value: '艺龙', label: '艺龙' },
  { value: '歌客/自来客', label: '歌客/自来客' },
]

const Orders = () => {
  const [loading, setLoading] = useState(false)
  const [dataSource, setDataSource] = useState([])
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })
  const [filters, setFilters] = useState({
    guest_source: '',
    order_status: '',
    keyword: '',
    check_in_range: null,
    check_out_range: null,
    reserve_range: null,
  })
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [orderDetail, setOrderDetail] = useState(null)

  // 获取订单列表
  const fetchOrders = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const currentBranchId = localStorage.getItem('currentBranchId')
      const params = {
        page,
        page_size: pageSize,
      }

      if (currentBranchId) {
        params.branch_id = parseInt(currentBranchId)
      }

      if (filters.guest_source) {
        params.guest_source = filters.guest_source
      }

      if (filters.order_status) {
        params.order_status = filters.order_status
      }

      if (filters.keyword) {
        params.keyword = filters.keyword
      }

      if (filters.check_in_range && filters.check_in_range.length === 2) {
        params.check_in_start = filters.check_in_range[0].format('YYYY-MM-DD')
        params.check_in_end = filters.check_in_range[1].format('YYYY-MM-DD')
      }

      if (filters.check_out_range && filters.check_out_range.length === 2) {
        params.check_out_start = filters.check_out_range[0].format('YYYY-MM-DD')
        params.check_out_end = filters.check_out_range[1].format('YYYY-MM-DD')
      }

      if (filters.reserve_range && filters.reserve_range.length === 2) {
        params.reserve_start = filters.reserve_range[0].format('YYYY-MM-DD HH:mm:ss')
        params.reserve_end = filters.reserve_range[1].format('YYYY-MM-DD HH:mm:ss')
      }

      const response = await getOrders(params)
      if (response.data.code === 200) {
        const data = response.data.data
        setDataSource(data.list || [])
        setPagination({
          current: data.page || 1,
          pageSize: data.page_size || 10,
          total: data.total || 0,
        })
      } else {
        message.error(response.data.msg || '获取订单列表失败')
      }
    } catch (error) {
      console.error('获取订单列表失败:', error)
      message.error('获取订单列表失败: ' + (error.response?.data?.msg || error.message))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchOrders(pagination.current, pagination.pageSize)
  }, [])

  // 处理搜索
  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    fetchOrders(1, pagination.pageSize)
  }

  // 处理分页
  const handleTableChange = (newPagination) => {
    fetchOrders(newPagination.current, newPagination.pageSize)
  }

  // 查看订单详情
  const handleViewDetail = async (orderId) => {
    try {
      const response = await getOrder(orderId)
      if (response.data.code === 200) {
        setOrderDetail(response.data.data)
        setDetailModalVisible(true)
      } else {
        message.error(response.data.msg || '获取订单详情失败')
      }
    } catch (error) {
      console.error('获取订单详情失败:', error)
      message.error('获取订单详情失败: ' + (error.response?.data?.msg || error.message))
    }
  }

  // 表格列定义
  const columns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
      width: 150,
      fixed: 'left',
    },
    {
      title: '客人来源',
      dataIndex: 'guest_source',
      key: 'guest_source',
      width: 120,
    },
    {
      title: '房间号',
      dataIndex: 'room_nos',
      key: 'room_nos',
      width: 150,
      render: (roomNos) => {
        if (!roomNos || roomNos.length === 0) {
          return '-'
        }
        return (
          <div>
            {roomNos.map((roomNo, index) => (
              <div key={index}>{roomNo}</div>
            ))}
          </div>
        )
      },
    },
    {
      title: '入住时间',
      dataIndex: 'check_in_time',
      key: 'check_in_time',
      width: 150,
      render: (text) => (text ? dayjs(text).format('YYYY-MM-DD') : '-'),
    },
    {
      title: '离店时间',
      dataIndex: 'check_out_time',
      key: 'check_out_time',
      width: 150,
      render: (text) => (text ? dayjs(text).format('YYYY-MM-DD') : '-'),
    },
    {
      title: '联系人',
      dataIndex: 'contact',
      key: 'contact',
      width: 100,
      render: (text) => text || '-',
    },
    {
      title: '手机号',
      dataIndex: 'contact_phone',
      key: 'contact_phone',
      width: 120,
      render: (text) => text || '-',
    },
    {
      title: '订单金额',
      dataIndex: 'order_amount',
      key: 'order_amount',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '已收押金',
      dataIndex: 'deposit_received',
      key: 'deposit_received',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '欠补费用',
      dataIndex: 'outstanding_amount',
      key: 'outstanding_amount',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '状态',
      dataIndex: 'order_status',
      key: 'order_status',
      width: 100,
      render: (status) => {
        const statusInfo = orderStatusMap[status] || { color: 'default', text: status }
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>
      },
    },
    {
      title: '预定时间',
      dataIndex: 'reserve_time',
      key: 'reserve_time',
      width: 180,
      render: (text) => (text ? dayjs(text).format('YYYY-MM-DD HH:mm:ss') : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record.id)}
        >
          查看订单
        </Button>
      ),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card title="订单管理">
        {/* 搜索筛选区域 */}
        <div style={{ marginBottom: 16 }}>
          <Space size="middle" wrap>
            <span>客人来源:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.guest_source || undefined}
              onChange={(value) => setFilters({ ...filters, guest_source: value || '' })}
            >
              {guestSourceOptions.map((option) => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>

            <span>入住时间:</span>
            <DatePicker.RangePicker
              value={filters.check_in_range}
              onChange={(dates) => setFilters({ ...filters, check_in_range: dates })}
              format="YYYY-MM-DD"
            />

            <span>离店时间:</span>
            <DatePicker.RangePicker
              value={filters.check_out_range}
              onChange={(dates) => setFilters({ ...filters, check_out_range: dates })}
              format="YYYY-MM-DD"
            />

            <span>预定时间:</span>
            <DatePicker.RangePicker
              showTime
              value={filters.reserve_range}
              onChange={(dates) => setFilters({ ...filters, reserve_range: dates })}
              format="YYYY-MM-DD HH:mm:ss"
            />
          </Space>
        </div>

        <div style={{ marginBottom: 16 }}>
          <Space size="middle" wrap>
            <span>订单状态:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.order_status || undefined}
              onChange={(value) => setFilters({ ...filters, order_status: value || '' })}
            >
              <Option value="RESERVED">已预定</Option>
              <Option value="CHECKED_IN">已入住</Option>
              <Option value="CHECKED_OUT">已退房</Option>
              <Option value="CANCELLED">已失效</Option>
            </Select>

            <span>关键词:</span>
            <Input
              style={{ width: 300 }}
              placeholder="订单号/房间号/手机号/联系人"
              value={filters.keyword}
              onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
              allowClear
            />

            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              查询
            </Button>
          </Space>
        </div>

        {/* 订单列表表格 */}
        <Table
          columns={columns}
          dataSource={dataSource}
          loading={loading}
          rowKey="id"
          scroll={{ x: 1500 }}
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
            showQuickJumper: true,
          }}
          onChange={handleTableChange}
        />
      </Card>

      {/* 订单详情模态框 */}
      <Modal
        title="订单详情"
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false)
          setOrderDetail(null)
        }}
        footer={null}
        width={800}
      >
        {orderDetail && (
          <Descriptions column={2} bordered>
            <Descriptions.Item label="订单号">{orderDetail.order_no}</Descriptions.Item>
            <Descriptions.Item label="订单状态">
              <Tag color={orderStatusMap[orderDetail.order_status]?.color || 'default'}>
                {orderStatusMap[orderDetail.order_status]?.text || orderDetail.order_status}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="客人来源">{orderDetail.guest_source}</Descriptions.Item>
            <Descriptions.Item label="分店名称">{orderDetail.branch_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="房间号" span={2}>
              {orderDetail.room_nos && orderDetail.room_nos.length > 0
                ? orderDetail.room_nos.join('、')
                : orderDetail.room_no || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="房间名称">{orderDetail.room_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="房型">{orderDetail.room_type_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="入住时间">
              {orderDetail.check_in_time
                ? dayjs(orderDetail.check_in_time).format('YYYY-MM-DD HH:mm:ss')
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="离店时间">
              {orderDetail.check_out_time
                ? dayjs(orderDetail.check_out_time).format('YYYY-MM-DD HH:mm:ss')
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="预定时间">
              {orderDetail.reserve_time
                ? dayjs(orderDetail.reserve_time).format('YYYY-MM-DD HH:mm:ss')
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="联系人">{orderDetail.contact || '-'}</Descriptions.Item>
            <Descriptions.Item label="手机号">{orderDetail.contact_phone || '-'}</Descriptions.Item>
            <Descriptions.Item label="客人姓名">{orderDetail.guest_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="入住人数">{orderDetail.guest_count || '-'}</Descriptions.Item>
            <Descriptions.Item label="房间数量">{orderDetail.room_count || '-'}</Descriptions.Item>
            <Descriptions.Item label="订单金额">
              <span style={{ color: '#f5222d', fontWeight: 'bold' }}>
                ¥{orderDetail.order_amount?.toFixed(2) || '0.00'}
              </span>
            </Descriptions.Item>
            <Descriptions.Item label="已收押金">
              ¥{orderDetail.deposit_received?.toFixed(2) || '0.00'}
            </Descriptions.Item>
            <Descriptions.Item label="欠补费用">
              <span style={{ color: orderDetail.outstanding_amount > 0 ? '#f5222d' : '#52c41a' }}>
                ¥{orderDetail.outstanding_amount?.toFixed(2) || '0.00'}
              </span>
            </Descriptions.Item>
            <Descriptions.Item label="违约金">
              ¥{orderDetail.penalty_amount?.toFixed(2) || '0.00'}
            </Descriptions.Item>
            <Descriptions.Item label="支付方式">{orderDetail.pay_type || '-'}</Descriptions.Item>
            {orderDetail.special_request && (
              <Descriptions.Item label="特殊需求" span={2}>
                {orderDetail.special_request}
              </Descriptions.Item>
            )}
            <Descriptions.Item label="创建时间">
              {orderDetail.created_at
                ? dayjs(orderDetail.created_at).format('YYYY-MM-DD HH:mm:ss')
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {orderDetail.updated_at
                ? dayjs(orderDetail.updated_at).format('YYYY-MM-DD HH:mm:ss')
                : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  )
}

export default Orders
