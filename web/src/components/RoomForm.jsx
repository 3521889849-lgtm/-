import React, { useEffect, useState } from 'react'
import {
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  Button,
  message,
  Upload,
  Checkbox,
} from 'antd'
import { UploadOutlined } from '@ant-design/icons'
import {
  createRoomInfo,
  updateRoomInfo,
  getRoomTypes,
  getFacilities,
  getCancellationPolicies,
  getRoomFacilities,
  setRoomFacilities,
  uploadRoomImages,
  getBranches,
} from '../api/room'

const { Option } = Select
const { TextArea } = Input

const RoomForm = ({ visible, room, onCancel, onSuccess }) => {
  const [form] = Form.useForm()
  const [branches, setBranches] = useState([])
  const [roomTypes, setRoomTypes] = useState([])
  const [facilities, setFacilities] = useState([])
  const [policies, setPolicies] = useState([])
  const [selectedFacilities, setSelectedFacilities] = useState([])
  const [roomFacilities, setRoomFacilities] = useState([])
  const [loading, setLoading] = useState(false)
  const [fileList, setFileList] = useState([])

  useEffect(() => {
    if (visible) {
      loadOptions()
      if (room) {
        form.setFieldsValue({
          ...room,
          has_breakfast: room.has_breakfast,
          has_toiletries: room.has_toiletries,
        })
        loadRoomFacilities(room.id)
      } else {
        form.resetFields()
        setSelectedFacilities([])
        setFileList([])
      }
    }
  }, [visible, room])

  const loadOptions = async () => {
    try {
      const [branchesRes, typesRes, facilitiesRes, policiesRes] = await Promise.all([
        getBranches({ page: 1, page_size: 100 }),
        getRoomTypes({ page: 1, page_size: 100 }),
        getFacilities({ page: 1, page_size: 100 }),
        getCancellationPolicies({ page: 1, page_size: 100 }),
      ])

      if (branchesRes.data.code === 200) {
        setBranches(branchesRes.data.data.list || [])
      }
      if (typesRes.data.code === 200) {
        setRoomTypes(typesRes.data.data.list || [])
      }
      if (facilitiesRes.data.code === 200) {
        setFacilities(facilitiesRes.data.data.list || [])
      }
      if (policiesRes.data.code === 200) {
        setPolicies(policiesRes.data.data.list || [])
      }
    } catch (error) {
      console.error('加载选项失败:', error)
    }
  }

  const loadRoomFacilities = async (roomId) => {
    try {
      const response = await getRoomFacilities(roomId)
      if (response.data.code === 200) {
        const facilityIds = (response.data.data || []).map((f) => f.id)
        setSelectedFacilities(facilityIds)
        setRoomFacilities(facilityIds)
      }
    } catch (error) {
      console.error('加载房源设施失败:', error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setLoading(true)

      if (room) {
        // 更新
        await updateRoomInfo(room.id, values)
        message.success('更新成功')
      } else {
        // 创建
        // 添加必填字段：room_count 和 created_by
        const createData = {
          ...values,
          room_count: values.room_count || 1, // 默认房间数量为1
          created_by: 1, // 暂时使用固定值，实际应该从用户信息中获取
        }
        const response = await createRoomInfo(createData)
        if (response.data.code === 200) {
          const newRoomId = response.data.data.id

          // 设置设施
          if (selectedFacilities.length > 0) {
            await setRoomFacilities(newRoomId, { facility_ids: selectedFacilities })
          }

          // 上传图片
          if (fileList.length > 0) {
            const formData = new FormData()
            fileList.forEach((file) => {
              formData.append('images', file.originFileObj)
            })
            await uploadRoomImages(newRoomId, formData)
          }

          message.success('创建成功')
        }
      }

      onSuccess()
    } catch (error) {
      message.error(room ? '更新失败' : '创建失败')
    } finally {
      setLoading(false)
    }
  }

  const handleFacilityChange = (checkedValues) => {
    setSelectedFacilities(checkedValues)
  }

  const uploadProps = {
    beforeUpload: (file) => {
      const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
      if (!isJpgOrPng) {
        message.error('只能上传 JPG/PNG 格式的图片!')
        return false
      }
      const isLt5M = file.size / 1024 / 1024 < 5
      if (!isLt5M) {
        message.error('图片大小不能超过 5MB!')
        return false
      }
      return false // 阻止自动上传
    },
    fileList,
    onChange: ({ fileList: newFileList }) => {
      if (newFileList.length > 16) {
        message.warning('最多只能上传 16 张图片')
        return
      }
      setFileList(newFileList)
    },
  }

  return (
    <Modal
      title={room ? '修改房源' : '添加房源'}
      open={visible}
      onCancel={onCancel}
      width={800}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          取消
        </Button>,
        <Button key="submit" type="primary" loading={loading} onClick={handleSubmit}>
          确定
        </Button>,
      ]}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="branch_id"
          label="分店"
          rules={[{ required: true, message: '请选择分店' }]}
        >
          <Select placeholder="请选择分店">
            {branches.map((branch) => (
              <Option key={branch.id} value={branch.id}>
                {branch.hotel_name} {branch.branch_code ? `(${branch.branch_code})` : ''}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="room_type_id"
          label="房型"
          rules={[{ required: true, message: '请选择房型' }]}
        >
          <Select placeholder="请选择房型">
            {roomTypes.map((type) => (
              <Option key={type.id} value={type.id}>
                {type.room_type_name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="room_no"
          label="房间号"
          rules={[{ required: true, message: '请输入房间号' }]}
        >
          <Input placeholder="请输入房间号" />
        </Form.Item>

        <Form.Item
          name="room_name"
          label="房源名称"
          rules={[{ required: true, message: '请输入房源名称' }]}
        >
          <Input placeholder="请输入房源名称" />
        </Form.Item>

        <Form.Item
          name="market_price"
          label="门市价"
          rules={[{ required: true, message: '请输入门市价' }]}
        >
          <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="请输入门市价" />
        </Form.Item>

        <Form.Item
          name="calendar_price"
          label="日历价"
          rules={[{ required: true, message: '请输入日历价' }]}
        >
          <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="请输入日历价" />
        </Form.Item>

        <Form.Item name="area" label="面积（平方米）">
          <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="请输入面积" />
        </Form.Item>

        <Form.Item
          name="room_count"
          label="房间数量"
          rules={[{ required: true, message: '请输入房间数量' }]}
          initialValue={1}
        >
          <InputNumber min={1} max={10} style={{ width: '100%' }} placeholder="请输入房间数量" />
        </Form.Item>

        <Form.Item
          name="bed_spec"
          label="床型规格"
          rules={[{ required: true, message: '请输入床型规格' }]}
        >
          <Input placeholder="如：1.8*2.0m" />
        </Form.Item>

        <Form.Item name="has_breakfast" label="是否含早" valuePropName="checked">
          <Switch />
        </Form.Item>

        <Form.Item name="has_toiletries" label="是否提供洗漱用品" valuePropName="checked">
          <Switch />
        </Form.Item>

        <Form.Item name="cancellation_policy_id" label="退订政策">
          <Select placeholder="请选择退订政策" allowClear>
            {policies.map((policy) => (
              <Option key={policy.id} value={policy.id}>
                {policy.policy_name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item label="房间设施">
          <Checkbox.Group
            value={selectedFacilities}
            onChange={handleFacilityChange}
            style={{ width: '100%' }}
          >
            {facilities.map((facility) => (
              <Checkbox key={facility.id} value={facility.id}>
                {facility.facility_name}
              </Checkbox>
            ))}
          </Checkbox.Group>
        </Form.Item>

        {!room && (
          <Form.Item label="房源图片（最多16张，400x300，jpg/png）">
            <Upload {...uploadProps} listType="picture-card" multiple>
              {fileList.length < 16 && <UploadOutlined />}
            </Upload>
          </Form.Item>
        )}
      </Form>
    </Modal>
  )
}

export default RoomForm
