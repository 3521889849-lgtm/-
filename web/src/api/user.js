import axios from 'axios'

const API_BASE = '/api/v1'

// 账号管理
export const getUserAccounts = (params) => {
  return axios.get(`${API_BASE}/user-accounts`, { params })
}

export const getUserAccount = (id) => {
  return axios.get(`${API_BASE}/user-accounts/${id}`)
}

export const createUserAccount = (data) => {
  return axios.post(`${API_BASE}/user-accounts`, data)
}

export const updateUserAccount = (id, data) => {
  return axios.put(`${API_BASE}/user-accounts/${id}`, data)
}

export const deleteUserAccount = (id) => {
  return axios.delete(`${API_BASE}/user-accounts/${id}`)
}

// 角色管理
export const getRoles = (params) => {
  return axios.get(`${API_BASE}/roles`, { params })
}

export const getRole = (id) => {
  return axios.get(`${API_BASE}/roles/${id}`)
}

export const createRole = (data) => {
  return axios.post(`${API_BASE}/roles`, data)
}

export const updateRole = (id, data) => {
  return axios.put(`${API_BASE}/roles/${id}`, data)
}

export const deleteRole = (id) => {
  return axios.delete(`${API_BASE}/roles/${id}`)
}

// 权限管理
export const getPermissions = (params) => {
  return axios.get(`${API_BASE}/permissions`, { params })
}

// 渠道配置
export const getChannelConfigs = (params) => {
  return axios.get(`${API_BASE}/channel-configs`, { params })
}

export const getChannelConfig = (id) => {
  return axios.get(`${API_BASE}/channel-configs/${id}`)
}

export const createChannelConfig = (data) => {
  return axios.post(`${API_BASE}/channel-configs`, data)
}

export const updateChannelConfig = (id, data) => {
  return axios.put(`${API_BASE}/channel-configs/${id}`, data)
}

export const deleteChannelConfig = (id) => {
  return axios.delete(`${API_BASE}/channel-configs/${id}`)
}

// 系统配置
export const getSystemConfigs = (params) => {
  return axios.get(`${API_BASE}/system-configs`, { params })
}

export const getSystemConfig = (id) => {
  return axios.get(`${API_BASE}/system-configs/${id}`)
}

export const createSystemConfig = (data) => {
  return axios.post(`${API_BASE}/system-configs`, data)
}

export const updateSystemConfig = (id, data) => {
  return axios.put(`${API_BASE}/system-configs/${id}`, data)
}

export const deleteSystemConfig = (id) => {
  return axios.delete(`${API_BASE}/system-configs/${id}`)
}

export const getSystemConfigsByCategory = (category) => {
  return axios.get(`${API_BASE}/system-configs/category/${category}`)
}

// 黑名单管理
export const getBlacklists = (params) => {
  return axios.get(`${API_BASE}/blacklists`, { params })
}

export const getBlacklist = (id) => {
  return axios.get(`${API_BASE}/blacklists/${id}`)
}

export const createBlacklist = (data) => {
  return axios.post(`${API_BASE}/blacklists`, data)
}

export const updateBlacklist = (id, data) => {
  return axios.put(`${API_BASE}/blacklists/${id}`, data)
}

export const deleteBlacklist = (id) => {
  return axios.delete(`${API_BASE}/blacklists/${id}`)
}
