import React, { useState } from 'react'
import { Card, Button, Space, Tag, Spin, App } from 'antd'
import { CheckCircleOutlined, CloseCircleOutlined, ReloadOutlined } from '@ant-design/icons'
import axios from 'axios'

const ConnectionTest = () => {
  const { message } = App.useApp()
  const [loading, setLoading] = useState(false)
  const [results, setResults] = useState({})

  const testAPI = async (name, url, method = 'GET', data = null) => {
    try {
      const config = {
        method,
        url,
        timeout: 5000,
      }
      if (data) {
        config.data = data
      }
      const response = await axios(config)
      return { success: true, status: response.status, data: response.data }
    } catch (error) {
      return {
        success: false,
        status: error.response?.status,
        message: error.message,
        error: error.response?.data,
      }
    }
  }

  const runTests = async () => {
    setLoading(true)
    const testResults = {}

    // 测试健康检查
    const healthResult = await testAPI('健康检查', '/health')
    testResults.health = healthResult

    // 测试获取房型列表
    const roomTypesResult = await testAPI('获取房型列表', '/api/v1/room-types?page=1&page_size=10')
    testResults.roomTypes = roomTypesResult

    // 测试获取房源列表
    const roomInfosResult = await testAPI('获取房源列表', '/api/v1/room-infos?page=1&page_size=10')
    testResults.roomInfos = roomInfosResult

    // 测试获取设施列表
    const facilitiesResult = await testAPI('获取设施列表', '/api/v1/facilities?page=1&page_size=10')
    testResults.facilities = facilitiesResult

    // 测试获取退订政策列表
    const policiesResult = await testAPI('获取退订政策列表', '/api/v1/cancellation-policies?page=1&page_size=10')
    testResults.policies = policiesResult

    setResults(testResults)
    setLoading(false)

    // 统计结果
    const successCount = Object.values(testResults).filter((r) => r.success).length
    const totalCount = Object.keys(testResults).length

    if (successCount === totalCount) {
      message.success(`所有测试通过！(${successCount}/${totalCount})`)
    } else {
      message.warning(`部分测试失败：${successCount}/${totalCount} 通过`)
    }
  }

  const testItems = [
    { key: 'health', name: '健康检查', endpoint: '/api/v1/health' },
    { key: 'roomTypes', name: '获取房型列表', endpoint: '/api/v1/room-types' },
    { key: 'roomInfos', name: '获取房源列表', endpoint: '/api/v1/room-infos' },
    { key: 'facilities', name: '获取设施列表', endpoint: '/api/v1/facilities' },
    { key: 'policies', name: '获取退订政策列表', endpoint: '/api/v1/cancellation-policies' },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title="后端连接测试"
        extra={
          <Button type="primary" icon={<ReloadOutlined />} onClick={runTests} loading={loading}>
            开始测试
          </Button>
        }
      >
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          {testItems.map((item) => {
            const result = results[item.key]
            return (
              <Card key={item.key} size="small">
                <Space>
                  {result ? (
                    result.success ? (
                      <Tag color="success" icon={<CheckCircleOutlined />}>
                        成功
                      </Tag>
                    ) : (
                      <Tag color="error" icon={<CloseCircleOutlined />}>
                        失败
                      </Tag>
                    )
                  ) : (
                    <Tag>未测试</Tag>
                  )}
                  <span style={{ fontWeight: 'bold' }}>{item.name}</span>
                  <span style={{ color: '#999' }}>{item.endpoint}</span>
                  {result && (
                    <>
                      <Tag>状态码: {result.status}</Tag>
                      {result.message && <Tag color="error">{result.message}</Tag>}
                    </>
                  )}
                </Space>
                {result && result.error && (
                  <div style={{ marginTop: 8, color: '#ff4d4f', fontSize: 12 }}>
                    错误: {JSON.stringify(result.error)}
                  </div>
                )}
              </Card>
            )
          })}
        </Space>

        <div style={{ marginTop: 24, padding: 16, background: '#f5f5f5', borderRadius: 4 }}>
          <h4>配置信息：</h4>
          <p>前端地址: http://localhost:3000</p>
          <p>后端地址: http://localhost:8080</p>
          <p>API 前缀: /api/v1</p>
          <p>代理配置: Vite 已配置代理，所有 /api 请求会转发到 http://localhost:8080</p>
        </div>
      </Card>
    </div>
  )
}

export default ConnectionTest
