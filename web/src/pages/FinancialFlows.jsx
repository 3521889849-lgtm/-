import React, { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Input,
  Select,
  DatePicker,
  Space,
  Card,
  App,
  Tag,
} from 'antd'
import { SearchOutlined, ExportOutlined, EyeOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { getFinancialFlows, getOrder } from '../api/room'

const { Option } = Select
const { RangePicker } = DatePicker

// 收支类型选项
const flowTypeOptions = [
  { value: '收入', label: '收入' },
  { value: '支出', label: '支出' },
]

// 收支项目选项（示例）
const flowItemOptions = [
  { value: '收取房费', label: '收取房费' },
  { value: '收取押金', label: '收取押金' },
  { value: '违约金', label: '违约金' },
  { value: '退款', label: '退款' },
]

// 支付方式选项
const payTypeOptions = [
  { value: '现金', label: '现金' },
  { value: '支付宝', label: '支付宝' },
  { value: '微信', label: '微信' },
  { value: '银联', label: '银联' },
  { value: '刷卡', label: '刷卡' },
  { value: '途游代收', label: '途游代收' },
  { value: '携程代收', label: '携程代收' },
  { value: '去哪儿代收', label: '去哪儿代收' },
]

const FinancialFlows = () => {
  const { message } = App.useApp()
  const [loading, setLoading] = useState(false)
  const [dataSource, setDataSource] = useState([])
  const [summary, setSummary] = useState(null)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 200,
    total: 0,
  })
  const [filters, setFilters] = useState({
    flow_type: '',
    flow_item: '',
    pay_type: '',
    operator_id: '',
    occur_range: null,
  })

  // 获取收支流水列表
  const fetchFinancialFlows = async (page = 1, pageSize = 200) => {
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

      if (filters.flow_type) {
        params.flow_type = filters.flow_type
      }
      if (filters.flow_item) {
        params.flow_item = filters.flow_item
      }
      if (filters.pay_type) {
        params.pay_type = filters.pay_type
      }
      if (filters.operator_id) {
        params.operator_id = parseInt(filters.operator_id)
      }
      if (filters.occur_range && filters.occur_range.length === 2) {
        params.occur_start = filters.occur_range[0].format('YYYY-MM-DD')
        params.occur_end = filters.occur_range[1].format('YYYY-MM-DD')
      }

      const response = await getFinancialFlows(params)
      if (response.data.code === 200) {
        const data = response.data.data
        setDataSource(data.list || [])
        setSummary(data.summary)
        setPagination({
          current: data.page || 1,
          pageSize: data.page_size || 200,
          total: data.total || 0,
        })
      } else {
        message.error(response.data.msg || '获取收支流水列表失败')
      }
    } catch (error) {
      console.error('获取收支流水列表失败:', error)
      message.error('获取收支流水列表失败: ' + (error.response?.data?.msg || error.message))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchFinancialFlows(pagination.current, pagination.pageSize)
  }, [])

  // 处理搜索
  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    fetchFinancialFlows(1, pagination.pageSize)
  }

  // 处理分页
  const handleTableChange = (newPagination) => {
    fetchFinancialFlows(newPagination.current, newPagination.pageSize)
  }

  // 查看订单详情
  const handleViewOrder = async (orderId) => {
    if (!orderId) return
    try {
      const response = await getOrder(orderId)
      if (response.data.code === 200) {
        // 这里可以打开订单详情模态框
        message.info('订单详情功能待实现')
      } else {
        message.error(response.data.msg || '获取订单详情失败')
      }
    } catch (error) {
      console.error('获取订单详情失败:', error)
      message.error('获取订单详情失败: ' + (error.response?.data?.msg || error.message))
    }
  }

  // 导出Excel（占位功能）
  const handleExport = () => {
    message.info('导出Excel功能待实现')
  }

  // 汇总表格列
  const summaryColumns = [
    {
      title: '收入/支出',
      dataIndex: 'type',
      key: 'type',
      width: 100,
    },
    {
      title: '合计',
      dataIndex: 'total',
      key: 'total',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '现金',
      dataIndex: 'cash',
      key: 'cash',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '支付宝',
      dataIndex: 'alipay',
      key: 'alipay',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '微信',
      dataIndex: 'wechat',
      key: 'wechat',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '银联',
      dataIndex: 'unionpay',
      key: 'unionpay',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '刷卡',
      dataIndex: 'card_swipe',
      key: 'card_swipe',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '途游代收',
      dataIndex: 'tuyou_collection',
      key: 'tuyou_collection',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '携程代收',
      dataIndex: 'ctrip_collection',
      key: 'ctrip_collection',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '去哪儿代收',
      dataIndex: 'qunar_collection',
      key: 'qunar_collection',
      width: 120,
      align: 'right',
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
  ]

  // 汇总数据
  const summaryData = summary ? [
    {
      key: 'income',
      type: '收入',
      total: summary.income?.total || 0,
      cash: summary.income?.cash || 0,
      alipay: summary.income?.alipay || 0,
      wechat: summary.income?.wechat || 0,
      unionpay: summary.income?.unionpay || 0,
      card_swipe: summary.income?.card_swipe || 0,
      tuyou_collection: summary.income?.tuyou_collection || 0,
      ctrip_collection: summary.income?.ctrip_collection || 0,
      qunar_collection: summary.income?.qunar_collection || 0,
    },
    {
      key: 'expense',
      type: '支出',
      total: summary.expense?.total || 0,
      cash: summary.expense?.cash || 0,
      alipay: summary.expense?.alipay || 0,
      wechat: summary.expense?.wechat || 0,
      unionpay: summary.expense?.unionpay || 0,
      card_swipe: summary.expense?.card_swipe || 0,
      tuyou_collection: summary.expense?.tuyou_collection || 0,
      ctrip_collection: summary.expense?.ctrip_collection || 0,
      qunar_collection: summary.expense?.qunar_collection || 0,
    },
    {
      key: 'balance',
      type: '结余',
      total: summary.balance?.total || 0,
      cash: summary.balance?.cash || 0,
      alipay: summary.balance?.alipay || 0,
      wechat: summary.balance?.wechat || 0,
      unionpay: summary.balance?.unionpay || 0,
      card_swipe: summary.balance?.card_swipe || 0,
      tuyou_collection: summary.balance?.tuyou_collection || 0,
      ctrip_collection: summary.balance?.ctrip_collection || 0,
      qunar_collection: summary.balance?.qunar_collection || 0,
    },
  ] : []

  // 明细表格列
  const columns = [
    {
      title: '记录时间',
      dataIndex: 'occur_time',
      key: 'occur_time',
      width: 180,
      render: (text) => (text ? dayjs(text).format('YYYY-MM-DD HH:mm:ss') : '-'),
    },
    {
      title: '操作人',
      dataIndex: 'operator_name',
      key: 'operator_name',
      width: 100,
      render: (text) => text || '-',
    },
    {
      title: '房间号',
      dataIndex: 'room_no',
      key: 'room_no',
      width: 100,
      render: (text) => text || '-',
    },
    {
      title: '客人姓名',
      dataIndex: 'guest_name',
      key: 'guest_name',
      width: 120,
      render: (text) => text || '-',
    },
    {
      title: '联系电话',
      dataIndex: 'contact_phone',
      key: 'contact_phone',
      width: 120,
      render: (text) => text || '-',
    },
    {
      title: '收支项目',
      dataIndex: 'flow_item',
      key: 'flow_item',
      width: 120,
      render: (text, record) => (
        <Tag color={record.flow_type === '收入' ? 'green' : 'red'}>
          {text}
        </Tag>
      ),
    },
    {
      title: '支付方式',
      dataIndex: 'pay_type',
      key: 'pay_type',
      width: 100,
    },
    {
      title: '收入',
      dataIndex: 'amount',
      key: 'income',
      width: 120,
      align: 'right',
      render: (amount, record) => 
        record.flow_type === '收入' ? (
          <span style={{ color: '#52c41a', fontWeight: 'bold' }}>
            ¥{amount?.toFixed(2) || '0.00'}
          </span>
        ) : '-',
    },
    {
      title: '支出',
      dataIndex: 'amount',
      key: 'expense',
      width: 120,
      align: 'right',
      render: (amount, record) => 
        record.flow_type === '支出' ? (
          <span style={{ color: '#f5222d', fontWeight: 'bold' }}>
            ¥{amount?.toFixed(2) || '0.00'}
          </span>
        ) : '-',
    },
    {
      title: '订单详情',
      key: 'action',
      width: 100,
      render: (_, record) => (
        record.order_id ? (
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleViewOrder(record.order_id)}
          >
            查看订单
          </Button>
        ) : '-'
      ),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card title="收支流水">
        {/* 搜索筛选区域 */}
        <div style={{ marginBottom: 16 }}>
          <Space size="middle" wrap>
            <span>支出项目:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.flow_item || undefined}
              onChange={(value) => setFilters({ ...filters, flow_item: value || '' })}
            >
              {flowItemOptions.map((option) => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>

            <span>发生时间:</span>
            <RangePicker
              value={filters.occur_range}
              onChange={(dates) => setFilters({ ...filters, occur_range: dates })}
              format="YYYY-MM-DD"
            />

            <span>支付方式:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.pay_type || undefined}
              onChange={(value) => setFilters({ ...filters, pay_type: value || '' })}
            >
              {payTypeOptions.map((option) => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>

            <span>收入项目:</span>
            <Select
              style={{ width: 150 }}
              placeholder="全部"
              allowClear
              value={filters.flow_item || undefined}
              onChange={(value) => setFilters({ ...filters, flow_item: value || '' })}
            >
              {flowItemOptions.map((option) => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>
          </Space>
        </div>

        <div style={{ marginBottom: 16 }}>
          <Space size="middle" wrap>
            <span>操作人</span>
            <Input
              style={{ width: 150 }}
              placeholder="操作人ID"
              value={filters.operator_id}
              onChange={(e) => setFilters({ ...filters, operator_id: e.target.value })}
              allowClear
            />

            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              查询
            </Button>
          </Space>
        </div>

        {/* 收入/支出汇总表 */}
        <div style={{ marginBottom: 16 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
            <h3>收入/支出</h3>
            <Button icon={<ExportOutlined />} onClick={handleExport}>
              导出excle
            </Button>
          </div>
          <Table
            columns={summaryColumns}
            dataSource={summaryData}
            pagination={false}
            size="small"
            bordered
          />
        </div>

        {/* 操作按钮 */}
        <div style={{ marginBottom: 16 }}>
          <Space>
            <Button type="primary" icon={<ExportOutlined />}>
              订单收支明细
            </Button>
            <Button icon={<ExportOutlined />}>
              记一笔收支明细
            </Button>
          </Space>
        </div>

        {/* 收支明细列表表格 */}
        <div style={{ marginBottom: 8 }}>
          <Button icon={<ExportOutlined />} onClick={handleExport} style={{ float: 'right' }}>
            导出excle
          </Button>
        </div>
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

export default FinancialFlows
