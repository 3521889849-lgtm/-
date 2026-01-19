import axios from 'axios'

const API_BASE = '/api/v1'

// 会员管理 API
export const getMembers = async (params = {}) => {
  const response = await axios.get(`${API_BASE}/members`, { params })
  return response.data
}

export const getMember = async (id) => {
  const response = await axios.get(`${API_BASE}/members/${id}`)
  return response.data
}

export const getMemberByGuestID = async (guestId) => {
  const response = await axios.get(`${API_BASE}/members/guest/${guestId}`)
  return response.data
}

export const createMember = async (data) => {
  const response = await axios.post(`${API_BASE}/members`, data)
  return response.data
}

export const updateMember = async (id, data) => {
  const response = await axios.put(`${API_BASE}/members/${id}`, data)
  return response.data
}

export const deleteMember = async (id) => {
  const response = await axios.delete(`${API_BASE}/members/${id}`)
  return response.data
}

// 会员权益管理 API
export const getMemberRights = async (params = {}) => {
  const response = await axios.get(`${API_BASE}/member-rights`, { params })
  return response.data
}

export const getMemberRight = async (id) => {
  const response = await axios.get(`${API_BASE}/member-rights/${id}`)
  return response.data
}

export const getRightsByMemberLevel = async (memberLevel) => {
  const response = await axios.get(`${API_BASE}/member-rights/level/${memberLevel}`)
  return response.data
}

export const createMemberRights = async (data) => {
  const response = await axios.post(`${API_BASE}/member-rights`, data)
  return response.data
}

export const updateMemberRights = async (id, data) => {
  const response = await axios.put(`${API_BASE}/member-rights/${id}`, data)
  return response.data
}

export const deleteMemberRights = async (id) => {
  const response = await axios.delete(`${API_BASE}/member-rights/${id}`)
  return response.data
}

// 会员积分管理 API
export const getPointsRecords = async (params = {}) => {
  const response = await axios.get(`${API_BASE}/points-records`, { params })
  return response.data
}

export const createPointsRecord = async (data) => {
  const response = await axios.post(`${API_BASE}/points-records`, data)
  return response.data
}

export const getMemberPointsBalance = async (memberId) => {
  const response = await axios.get(`${API_BASE}/members/${memberId}/points-balance`)
  return response.data
}
