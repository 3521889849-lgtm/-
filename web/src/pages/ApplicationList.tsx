import React, { useState, useEffect } from 'react';
import { Search, ChevronDown, Eye, Edit, Trash, Plus, MoreHorizontal } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Application, getApplications } from '../lib/mockStore';

const ApplicationList: React.FC = () => {
  const navigate = useNavigate();
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    setLoading(true);
    try {
      // Use Mock Store
      const data = getApplications();
      setApplications(data);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'passed':
        return <span className="px-3 py-1 bg-green-500 text-white rounded text-xs">通过</span>;
      case 'pending':
        return <span className="px-3 py-1 bg-orange-500 text-white rounded text-xs">待审核</span>;
      case 'rejected':
        return <span className="px-3 py-1 bg-red-500 text-white rounded text-xs">拒回</span>;
      default:
        return <span className="px-3 py-1 bg-gray-500 text-white rounded text-xs">未知</span>;
    }
  };

  return (
    <div className="bg-white rounded-lg shadow min-h-[600px] flex flex-col">
      <div className="p-6 border-b border-gray-100">
        <h2 className="text-xl font-bold text-gray-800 mb-6">入驻申请列表</h2>
        
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center space-x-2">
            <span className="font-bold text-gray-700">查询</span>
            <div className="relative">
              <input type="text" placeholder="搜索" className="pl-3 pr-8 py-1.5 border border-gray-200 rounded text-sm w-48" />
            </div>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="text-sm text-gray-600">状态下</span>
            <div className="relative w-32 border border-gray-200 rounded px-3 py-1.5 flex items-center justify-between bg-white">
              <span className="text-sm text-gray-400">选择状态</span>
              <ChevronDown size={14} className="text-gray-400" />
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <span className="text-sm text-gray-600">日期</span>
             <div className="w-32 h-8 border border-gray-200 rounded flex items-center px-2 text-gray-400 text-sm">
               2022-11-01
             </div>
             <span className="text-gray-400">-</span>
             <div className="w-32 h-8 border border-gray-200 rounded flex items-center px-2 text-gray-400 text-sm">
               2022-12-01
             </div>
             <button className="px-4 py-1.5 bg-blue-600 text-white rounded text-sm">日期</button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <table className="w-full">
          <thead className="bg-gray-50 text-gray-700 text-sm font-medium border-b border-gray-200">
            <tr>
              <th className="py-4 px-6 text-left w-12"><input type="checkbox" /></th>
              <th className="py-4 px-6 text-left">商家名称</th>
              <th className="py-4 px-6 text-left">申请日期</th>
              <th className="py-4 px-6 text-left">更新日期</th>
              <th className="py-4 px-6 text-left">ID</th>
              <th className="py-4 px-6 text-center">状态</th>
              <th className="py-4 px-6 text-center">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {applications.map((app) => (
              <tr key={app.id} className="hover:bg-blue-50 transition-colors">
                <td className="py-4 px-6"><input type="checkbox" /></td>
                <td className="py-4 px-6 font-medium text-gray-800">{app.merchant_name}</td>
                <td className="py-4 px-6 text-sm text-gray-500">{app.created_at}</td>
                <td className="py-4 px-6 text-sm text-gray-500">{app.updated_at}</td>
                <td className="py-4 px-6 text-sm text-gray-900">{app.id_number}</td>
                <td className="py-4 px-6 text-center">
                  {getStatusBadge(app.status)}
                </td>
                <td className="py-4 px-6">
                  <div className="flex items-center justify-center space-x-2">
                    {/* Action Buttons based on status */}
                    {app.status === 'passed' ? (
                       <button className="px-3 py-1 bg-green-500 text-white rounded text-xs">通过</button>
                    ) : app.status === 'rejected' ? (
                       <button className="px-3 py-1 bg-orange-200 text-orange-700 rounded text-xs">拒回</button>
                    ) : (
                       <button 
                         onClick={() => navigate(`/merchant/audit/${app.id}`)}
                         className="px-3 py-1 bg-red-500 text-white rounded text-xs hover:bg-red-600 transition-colors"
                       >
                         审核
                       </button>
                    )}
                    
                    <button 
                      onClick={() => navigate(`/merchant/audit/${app.id}`)}
                      className="p-1.5 bg-blue-100 text-blue-600 rounded hover:bg-blue-200 transition-colors"
                    >
                      <Eye size={14} />
                    </button>
                    <button className="p-1.5 bg-gray-100 text-gray-600 rounded">
                      <MoreHorizontal size={14} />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="p-4 border-t border-gray-200 flex justify-center">
        <div className="flex space-x-2">
          <button className="w-8 h-8 flex items-center justify-center border border-gray-300 rounded hover:bg-gray-50 text-sm">{'<'}</button>
          <button className="w-8 h-8 flex items-center justify-center bg-blue-600 text-white rounded text-sm">1</button>
          <button className="w-8 h-8 flex items-center justify-center border border-gray-300 rounded hover:bg-gray-50 text-sm">2</button>
          <button className="w-8 h-8 flex items-center justify-center border border-gray-300 rounded hover:bg-gray-50 text-sm">3</button>
          <span className="w-8 h-8 flex items-center justify-center text-gray-400">...</span>
          <button className="w-8 h-8 flex items-center justify-center border border-gray-300 rounded hover:bg-gray-50 text-sm">{'>'}</button>
        </div>
      </div>
    </div>
  );
};

export default ApplicationList;
