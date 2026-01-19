import React from 'react'
import { Result, Button } from 'antd'
import { useNavigate } from 'react-router-dom'

const Placeholder = ({ title = '功能开发中', subTitle = '该功能正在开发中，敬请期待' }) => {
  const navigate = useNavigate()

  return (
    <Result
      status="info"
      title={title}
      subTitle={subTitle}
      extra={
        <Button type="primary" onClick={() => navigate('/')}>
          返回首页
        </Button>
      }
    />
  )
}

export default Placeholder
