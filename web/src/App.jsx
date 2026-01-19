import React from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import RoomStatus from './pages/RoomStatus'
import RoomTypeManagement from './pages/RoomTypeManagement'
import FacilityManagement from './pages/FacilityManagement'
import CancellationPolicyManagement from './pages/CancellationPolicyManagement'
import ConnectionTest from './pages/ConnectionTest'
import CalendarRoomStatus from './pages/CalendarRoomStatus'
import RealTimeStatistics from './pages/RealTimeStatistics'
import Orders from './pages/Orders'
import Reports from './pages/Reports'
import ChannelSettings from './pages/ChannelSettings'
import HotelManagement from './pages/HotelManagement'
import InHouseGuests from './pages/InHouseGuests'
import FinancialFlows from './pages/FinancialFlows'
import UserAccountManagement from './pages/UserAccountManagement'
import RoleManagement from './pages/RoleManagement'
import SystemConfigManagement from './pages/SystemConfigManagement'
import BlacklistManagement from './pages/BlacklistManagement'
import MemberManagement from './pages/MemberManagement'
import MemberRightsManagement from './pages/MemberRightsManagement'
import MemberPointsManagement from './pages/MemberPointsManagement'
import OperationLogManagement from './pages/OperationLogManagement'
import Placeholder from './pages/Placeholder'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<RoomStatus />} />
          <Route path="calendar-room-status" element={<CalendarRoomStatus />} />
          <Route path="real-time-statistics" element={<RealTimeStatistics />} />
          <Route path="orders" element={<Orders />} />
          <Route path="reports" element={<Reports />} />
          <Route path="financial-flows" element={<FinancialFlows />} />
          <Route path="finance-overview" element={<Placeholder title="财务概览" subTitle="财务概览功能正在开发中" />} />
          <Route path="room-sales" element={<Placeholder title="房间销售" subTitle="房间销售统计功能正在开发中" />} />
          <Route path="channel-sales" element={<Placeholder title="渠道销售" subTitle="渠道销售统计功能正在开发中" />} />
          <Route path="daily-settlement" element={<Placeholder title="房费日结" subTitle="房费日结功能正在开发中" />} />
          <Route path="shift-handover" element={<Placeholder title="交接班记录" subTitle="交接班记录功能正在开发中" />} />
          <Route path="non-room-settlement" element={<Placeholder title="非房费日结" subTitle="非房费日结功能正在开发中" />} />
          <Route path="room-types" element={<RoomTypeManagement />} />
          <Route path="facilities" element={<FacilityManagement />} />
          <Route path="cancellation-policies" element={<CancellationPolicyManagement />} />
          <Route path="channels" element={<ChannelSettings />} />
          <Route path="payments" element={<Placeholder title="支付管理" subTitle="支付管理功能正在开发中" />} />
          <Route path="consumption" element={<Placeholder title="消费项目" subTitle="消费项目管理功能正在开发中" />} />
          <Route path="sms-templates" element={<Placeholder title="短信模板" subTitle="短信模板管理功能正在开发中" />} />
          <Route path="print-settings" element={<Placeholder title="打印设置" subTitle="打印设置功能正在开发中" />} />
          <Route path="hotels" element={<HotelManagement />} />
          <Route path="inhouse-guests" element={<InHouseGuests />} />
          <Route path="accounts" element={<UserAccountManagement />} />
          <Route path="roles" element={<RoleManagement />} />
          <Route path="system-configs" element={<SystemConfigManagement />} />
          <Route path="blacklists" element={<BlacklistManagement />} />
          <Route path="members" element={<MemberManagement />} />
          <Route path="member-rights" element={<MemberRightsManagement />} />
          <Route path="points" element={<MemberPointsManagement />} />
          <Route path="operation-logs" element={<OperationLogManagement />} />
          <Route path="connection-test" element={<ConnectionTest />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
