import service from '@/utils/request'

/**
 * 获取仪表盘统计数据
 * @returns {Promise} 仪表盘统计数据
 */
export const getDashboardStats = () => {
  return service({
    url: '/dashboard/stats',
    method: 'get'
  })
}