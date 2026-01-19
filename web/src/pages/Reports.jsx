import React from 'react'
import { Card } from 'antd'
import { FileTextOutlined } from '@ant-design/icons'

const Reports = () => {
  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={
          <span>
            <FileTextOutlined style={{ marginRight: 8 }} />
            报表管理
          </span>
        }
      >
        <div style={{ textAlign: 'center', padding: '100px 0', color: '#999' }}>
          <FileTextOutlined style={{ fontSize: 64, marginBottom: 16 }} />
          <div>报表管理功能开发中...</div>
        </div>
      </Card>
    </div>
  )
}

export default Reports
