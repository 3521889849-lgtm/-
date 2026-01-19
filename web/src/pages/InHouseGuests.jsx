import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Input,
  Select,
  Space,
  Card,
  App,
} from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { getInHouseGuests } from '../api/room'

const { Option } = Select

// 省份选项（示例，实际应该从后端获取）
const provinces = [
  { value: '四川省', label: '四川省' },
  { value: '北京市', label: '北京市' },
  { value: '上海市', label: '上海市' },
  { value: '广东省', label: '广东省' },
  { value: '浙江省', label: '浙江省' },
]

const InHouseGuests = () => {
  const { message } = App.useApp()
  const [loading, setLoading] = useState(false)
  const [dataSource, setDataSource] = useState([])
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 200,
    total: 0,
  })
  const [filters, setFilters] = useState({
    province: '',
    city: '',
    district: '',
    name: '',
    phone: '',
    id_number: '',
    room_no: '',
  })

  // 获取在住客人列表
  const fetchGuests = async (page = 1, pageSize = 200) => {
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

      if (filters.province) {
        params.province = filters.province
      }
      if (filters.city) {
        params.city = filters.city
      }
      if (filters.district) {
        params.district = filters.district
      }
      if (filters.name) {
        params.name = filters.name
      }
      if (filters.phone) {
        params.phone = filters.phone
      }
      if (filters.id_number) {
        params.id_number = filters.id_number
      }
      if (filters.room_no) {
        params.room_no = filters.room_no
      }

      const response = await getInHouseGuests(params)
      if (response.data.code === 200) {
        const data = response.data.data
        setDataSource(data.list || [])
        setPagination({
          current: data.page || 1,
          pageSize: data.page_size || 200,
          total: data.total || 0,
        })
      } else {
        message.error(response.data.msg || '获取在住客人列表失败')
      }
    } catch (error) {
      console.error('获取在住客人列表失败:', error)
      message.error('获取在住客人列表失败: ' + (error.response?.data?.msg || error.message))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchGuests(pagination.current, pagination.pageSize)
  }, [])

  // 处理搜索
  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    fetchGuests(1, pagination.pageSize)
  }

  // 处理分页
  const handleTableChange = (newPagination) => {
    fetchGuests(newPagination.current, newPagination.pageSize)
  }

  // 表格列定义
  const columns = [
    {
      title: '渠道来源',
      dataIndex: 'guest_source',
      key: 'guest_source',
      width: 100,
    },
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
      width: 150,
    },
    {
      title: '房间号',
      dataIndex: 'room_no',
      key: 'room_no',
      width: 100,
    },
    {
      title: '房型',
      dataIndex: 'room_type_name',
      key: 'room_type_name',
      width: 120,
    },
    {
      title: '证件类型',
      dataIndex: 'id_type',
      key: 'id_type',
      width: 100,
    },
    {
      title: '证件号',
      dataIndex: 'id_number',
      key: 'id_number',
      width: 180,
    },
    {
      title: '地址',
      dataIndex: 'address',
      key: 'address',
      width: 200,
      render: (text) => text || '-',
    },
    {
      title: '民族',
      dataIndex: 'ethnicity',
      key: 'ethnicity',
      width: 100,
      render: (text) => text || '-',
    },
    {
      title: '入住时间',
      dataIndex: 'check_in_time',
      key: 'check_in_time',
      width: 120,
      render: (text) => (text ? dayjs(text).format('YYYY-MM-DD') : '-'),
    },
    {
      title: '离店时间',
      dataIndex: 'check_out_time',
      key: 'check_out_time',
      width: 120,
      render: (text) => (text ? dayjs(text).format('YYYY-MM-DD') : '-'),
    },
    {
      title: '订单总额',
      dataIndex: 'order_amount',
      key: 'order_amount',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '已收房费',
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
      render: (amount) => (
        <span style={{ color: amount > 0 ? '#f5222d' : '#52c41a' }}>
          ¥{amount?.toFixed(2) || '0.00'}
        </span>
      ),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card title="在住客人">
        {/* 搜索筛选区域 */}
        <div style={{ marginBottom: 16 }}>
          <Space size="middle" wrap>
            <span>省份:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.province || undefined}
              onChange={(value) => setFilters({ ...filters, province: value || '' })}
            >
              {provinces.map((province) => (
                <Option key={province.value} value={province.value}>
                  {province.label}
                </Option>
              ))}
            </Select>

            <span>城市:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.city || undefined}
              onChange={(value) => setFilters({ ...filters, city: value || '' })}
            >
              <Option value="">全部</Option>
            </Select>

            <span>区县:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.district || undefined}
              onChange={(value) => setFilters({ ...filters, district: value || '' })}
            >
              <Option value="">全部</Option>
            </Select>
          </Space>
        </div>

        <div style={{ marginBottom: 16 }}>
          <Space size="middle" wrap>
            <span>姓名:</span>
            <Input
              style={{ width: 150 }}
              placeholder="姓名"
              value={filters.name}
              onChange={(e) => setFilters({ ...filters, name: e.target.value })}
              allowClear
            />

            <span>手机号</span>
            <Input
              style={{ width: 150 }}
              placeholder="手机号"
              value={filters.phone}
              onChange={(e) => setFilters({ ...filters, phone: e.target.value })}
              allowClear
            />

            <span>身份证号:</span>
            <Input
              style={{ width: 200 }}
              placeholder="身份证号"
              value={filters.id_number}
              onChange={(e) => setFilters({ ...filters, id_number: e.target.value })}
              allowClear
            />

            <span>房间号</span>
            <Input
              style={{ width: 150 }}
              placeholder="房间号"
              value={filters.room_no}
              onChange={(e) => setFilters({ ...filters, room_no: e.target.value })}
              allowClear
            />

            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              查询
            </Button>
          </Space>
        </div>

        {/* 在住客人列表表格 */}
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
            showTotal: (total) => `共${total} 条`,
            showQuickJumper: true,
            pageSizeOptions: ['50', '100', '200', '500'],
          }}
          onChange={handleTableChange}
        />
      </Card>
    </div>
  )
}

export default InHouseGuests
