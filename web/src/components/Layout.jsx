import React, { useState, useEffect } from 'react'
import { Layout as AntLayout, Menu, Dropdown } from 'antd'
import { getBranches } from '../api/room'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  HomeOutlined,
  SettingOutlined,
  UserOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
  CloudOutlined,
  TeamOutlined,
  DollarOutlined,
  SafetyOutlined,
  FileTextOutlined,
  ShopOutlined,
  WifiOutlined,
  FileProtectOutlined,
  CalendarOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation, Outlet } from 'react-router-dom'
import Logo from './Logo'
import './Layout.css'

const { Header, Sider, Content } = AntLayout

const Layout = () => {
  const [collapsed, setCollapsed] = useState(false)
  const [branches, setBranches] = useState([])
  const [currentBranch, setCurrentBranch] = useState(null)
  const navigate = useNavigate()
  const location = useLocation()

  // 获取分店列表
  useEffect(() => {
    getBranches({ status: 'ACTIVE' }).then((res) => {
      if (res.data.code === 200) {
        const data = res.data.data
        let branchList = []
        
        if (Array.isArray(data)) {
          branchList = data
        } else if (data && Array.isArray(data.list)) {
          branchList = data.list
        } else {
          console.error('分店列表数据格式错误:', data)
        }

        if (branchList.length > 0) {
          const savedBranchId = localStorage.getItem('currentBranchId')
          const savedBranch = savedBranchId
            ? branchList.find((b) => String(b.id) === String(savedBranchId))
            : null

          setBranches(branchList)
          const nextBranch = savedBranch || branchList[0]
          setCurrentBranch(nextBranch)
          localStorage.setItem('currentBranchId', nextBranch.id)
        } else {
            setBranches([])
        }
      }
    }).catch((err) => {
      console.error('获取分店列表失败:', err)
    })
  }, [])

  // 切换分店
  const handleBranchChange = (branch) => {
    setCurrentBranch(branch)
    localStorage.setItem('currentBranchId', branch.id)
    // 触发自定义事件通知其他组件
    window.dispatchEvent(new Event('branchChanged'))
    // 触发页面刷新以更新数据
    window.location.reload()
  }

  // 顶部导航菜单
  const topMenuItems = [
    { key: '/', label: '房源管理', icon: <HomeOutlined /> },
    { key: '/calendar-room-status', label: '日历房态', icon: <CalendarOutlined /> },
    { key: '/real-time-statistics', label: '实时统计', icon: <HomeOutlined /> },
    { key: '/orders', label: '订单', icon: <FileTextOutlined /> },
    { key: '/reports', label: '报表', icon: <FileTextOutlined /> },
    { key: '/connection-test', label: '连接测试', icon: <QuestionCircleOutlined /> },
  ]

  // 左侧菜单配置
  const menuItems = [
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置管理',
      children: [
        {
          key: 'basic',
          label: '基础设置',
          children: [
            { key: '/room-types', label: '房型管理', icon: <ShopOutlined /> },
            { key: '/facilities', label: '设施管理', icon: <WifiOutlined /> },
            { key: '/cancellation-policies', label: '退订政策', icon: <FileProtectOutlined /> },
            { key: '/channels', label: '渠道设置', icon: <CloudOutlined /> },
            { key: '/payments', label: '支付管理', icon: <DollarOutlined /> },
            { key: '/consumption', label: '消费项目', icon: <FileTextOutlined /> },
            { key: '/sms-templates', label: '短信模板', icon: <FileTextOutlined /> },
            { key: '/print-settings', label: '打印设置', icon: <FileTextOutlined /> },
          ],
        },
        {
          key: 'member',
          label: '会员管理',
          children: [
            { key: '/members', label: '会员管理', icon: <TeamOutlined /> },
            { key: '/member-rights', label: '会员权益', icon: <SafetyOutlined /> },
            { key: '/points', label: '积分管理', icon: <DollarOutlined /> },
          ],
        },
        {
          key: 'finance',
          label: '财务管理',
          children: [
            { key: '/finance-overview', label: '财务概览', icon: <DollarOutlined /> },
            { key: '/financial-flows', label: '收支流水', icon: <FileTextOutlined /> },
            { key: '/room-sales', label: '房间销售', icon: <ShopOutlined /> },
            { key: '/channel-sales', label: '渠道销售', icon: <CloudOutlined /> },
            { key: '/daily-settlement', label: '房费日结', icon: <FileTextOutlined /> },
            { key: '/shift-handover', label: '交接班记录', icon: <FileTextOutlined /> },
            { key: '/inhouse-guests', label: '在住客', icon: <UserOutlined /> },
            { key: '/non-room-settlement', label: '非房费日结', icon: <FileTextOutlined /> },
          ],
        },
        {
          key: 'system',
          label: '系统设置',
          children: [
            { key: '/accounts', label: '账号管理', icon: <UserOutlined /> },
            { key: '/roles', label: '角色管理', icon: <SafetyOutlined /> },
            { key: '/system-configs', label: '基础配置', icon: <SettingOutlined /> },
            { key: '/blacklists', label: '黑名单管理', icon: <SafetyOutlined /> },
            { key: '/operation-logs', label: '操作日志', icon: <FileTextOutlined /> },
            { key: '/hotels', label: '酒店管理', icon: <ShopOutlined /> },
          ],
        },
      ],
    },
  ]

  const handleMenuClick = ({ key }) => {
    if (key.startsWith('/')) {
      navigate(key)
    }
  }

  const userMenuItems = [
    { key: 'profile', label: '个人中心' },
    { key: 'settings', label: '设置' },
    { type: 'divider' },
    { key: 'logout', label: '退出登录' },
  ]

  return (
    <AntLayout className="layout-container">
      <Sider trigger={null} collapsible collapsed={collapsed} className="sidebar">
        <div className="logo">
          <Logo size="default" showText={!collapsed} />
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          className="sidebar-menu"
        />
      </Sider>
      <AntLayout>
        <Header className="header">
          <div className="header-left">
            {React.createElement(collapsed ? MenuUnfoldOutlined : MenuFoldOutlined, {
              className: 'trigger',
              onClick: () => setCollapsed(!collapsed),
            })}
            <Logo size="small" showText={true} />
            <span className="hotel-name">
              {currentBranch ? (currentBranch.branch_code ? `${currentBranch.hotel_name} (${currentBranch.branch_code})` : currentBranch.hotel_name) : '客栈酒店名称'}
            </span>
            <Dropdown
              menu={{
                items: (Array.isArray(branches) ? branches : []).map((branch) => ({
                  key: String(branch.id),
                  label: branch.branch_code ? `${branch.hotel_name} (${branch.branch_code})` : branch.hotel_name,
                })),
                onClick: ({ key }) => {
                  const branch = (Array.isArray(branches) ? branches : []).find((b) => String(b.id) === String(key))
                  if (branch) {
                    handleBranchChange(branch)
                  }
                },
              }}
            >
              <span className="switch-branch">
                切换分店 <span className="arrow">▼</span>
              </span>
            </Dropdown>
          </div>
          <div className="header-center">
            <Menu
              mode="horizontal"
              selectedKeys={[location.pathname === '/' ? '/' : location.pathname]}
              items={topMenuItems}
              onClick={handleMenuClick}
              className="top-menu"
            />
          </div>
          <div className="header-right">
            <SearchOutlined className="header-icon" />
            <Dropdown menu={{ items: userMenuItems }}>
              <span className="user-info">
                用户登录名 <span className="arrow">▼</span>
              </span>
            </Dropdown>
            <SettingOutlined className="header-icon" />
            <QuestionCircleOutlined className="header-icon" />
          </div>
        </Header>
        <Content className="content">
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  )
}

export default Layout
