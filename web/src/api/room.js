import axios from 'axios'

const API_BASE_URL = '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 可以在这里添加 token 等
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response) {
      console.error('API Error:', error.response.data)
    }
    return Promise.reject(error)
  }
)

// 房源管理 API
export const getRoomInfos = (params) => api.get('/room-infos', { params })
export const getRoomInfo = (id) => api.get(`/room-infos/${id}`)
export const createRoomInfo = (data) => api.post('/room-infos', data)
export const updateRoomInfo = (id, data) => api.put(`/room-infos/${id}`, data)
export const deleteRoomInfo = (id) => api.delete(`/room-infos/${id}`)
export const updateRoomStatus = (id, data) => api.put(`/room-infos/${id}/status`, data)
export const batchUpdateRoomStatus = (data) => api.put('/room-infos/batch-status', data)

// 房源设施 API
export const getRoomFacilities = (id) => api.get(`/room-infos/${id}/facilities`)
export const setRoomFacilities = (id, data) => api.put(`/room-infos/${id}/facilities`, data)
export const addRoomFacility = (id, data) => api.post(`/room-infos/${id}/facilities`, data)
export const removeRoomFacility = (id, facilityId) => api.delete(`/room-infos/${id}/facilities/${facilityId}`)

// 房源图片 API
export const getRoomImages = (id) => api.get(`/room-infos/${id}/images`)
export const uploadRoomImages = (id, formData) => api.post(`/room-infos/${id}/images`, formData, {
  headers: { 'Content-Type': 'multipart/form-data' },
})
export const deleteRoomImage = (id) => api.delete(`/room-infos/images/${id}`)

// 关联房 API
export const getRoomBindings = (id) => api.get(`/room-infos/${id}/bindings`)
export const createRoomBinding = (data) => api.post('/room-infos/bindings', data)
export const deleteRoomBinding = (id) => api.delete(`/room-infos/bindings/${id}`)

// 房型管理 API
export const getRoomTypes = (params) => api.get('/room-types', { params })
export const getRoomType = (id) => api.get(`/room-types/${id}`)
export const createRoomType = (data) => api.post('/room-types', data)
export const updateRoomType = (id, data) => api.put(`/room-types/${id}`, data)
export const deleteRoomType = (id) => api.delete(`/room-types/${id}`)

// 设施管理 API
export const getFacilities = (params) => api.get('/facilities', { params })
export const getFacility = (id) => api.get(`/facilities/${id}`)
export const createFacility = (data) => api.post('/facilities', data)
export const updateFacility = (id, data) => api.put(`/facilities/${id}`, data)
export const deleteFacility = (id) => api.delete(`/facilities/${id}`)

// 退订政策 API
export const getCancellationPolicies = (params) => api.get('/cancellation-policies', { params })
export const getCancellationPolicy = (id) => api.get(`/cancellation-policies/${id}`)
export const createCancellationPolicy = (data) => api.post('/cancellation-policies', data)
export const updateCancellationPolicy = (id, data) => api.put(`/cancellation-policies/${id}`, data)
export const deleteCancellationPolicy = (id) => api.delete(`/cancellation-policies/${id}`)

// 分店管理 API
export const getBranches = (params) => api.get('/branches', { params })
export const getBranch = (id) => api.get(`/branches/${id}`)

// 渠道同步 API
export const syncRoomStatusToChannel = (data) => api.post('/sync-room-status', data)

// 订单管理 API
export const getOrders = (params) => api.get('/orders', { params })
export const getOrder = (id) => api.get(`/orders/${id}`)

// 在住客人管理 API
export const getInHouseGuests = (params) => api.get('/in-house-guests', { params })

// 财务管理 API
export const getFinancialFlows = (params) => api.get('/financial-flows', { params })