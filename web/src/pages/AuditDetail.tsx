import React, { useState } from 'react';
import { ChevronRight, ChevronLeft, Maximize2, FileText, Bell, Settings, User } from 'lucide-react';
import clsx from 'clsx';
import client from '../api/client';
import { useParams, useNavigate } from 'react-router-dom';

const AuditDetail: React.FC = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [comment, setComment] = useState('');
  const [auditStatus, setAuditStatus] = useState<'approved' | 'rejected' | null>(null);
  const [loading, setLoading] = useState(false);

  // Mock data for images
  const images = [
    { id: 1, title: '营业执照正本', url: 'https://via.placeholder.com/600x400?text=Business+License' },
    { id: 2, title: '经营许可证', url: 'https://via.placeholder.com/600x400?text=Permit' },
    { id: 3, title: '法人身份证', url: 'https://via.placeholder.com/600x400?text=ID+Card' },
  ];
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  const handleSubmit = async (status: 'approved' | 'rejected') => {
    if (!comment) {
      alert('请填写审核意见');
      return;
    }
    setLoading(true);
    try {
      // Mapping status to int32 as per thrift definition: 2=Approved, 3=Rejected
      const resultValue = status === 'approved' ? 2 : 3;

      await client.post('/audit/manual/submit', {
        audit_no: id, // audit_no from URL params
        audit_result: resultValue,
        audit_opinion: comment,
        auditor_id: '1001', // Mock auditor ID
        auditor_name: '管理员' // Mock auditor name
      });
      alert('审核提交成功');
      navigate('/merchant/applications');
    } catch (error: any) {
      console.error(error);
      alert('提交失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col h-[calc(100vh-80px)]">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">人工审核原型</h2>
          <div className="flex items-center space-x-2 mt-1">
            <span className="text-blue-600 font-medium">人工风格设计</span>
          </div>
        </div>
        
        <div className="flex items-center space-x-4 bg-white p-2 rounded-lg shadow-sm border border-gray-100">
           <div className="px-4 py-1.5 bg-blue-50 text-blue-600 rounded-md text-sm font-medium">工单编号: {id || '20230912'}</div>
           <div className="text-gray-500 text-sm">申请企业名称: {appData?.merchant_name || '示例票务公司'}</div>
           <div className="text-gray-400 text-sm">时间: {appData?.created_at || '9月12日'}</div>
        </div>
        
        <div className="flex space-x-3">
          <button className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm">时间时间</button>
          <button className="px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg text-sm">审核操作</button>
        </div>
      </div>

      <div className="flex flex-1 gap-6 overflow-hidden">
        {/* Left: Workflow */}
        <div className="w-48 bg-white rounded-xl shadow-sm p-6 border border-gray-100 flex flex-col items-center">
           <div className="space-y-8 relative">
              <div className="absolute left-1/2 top-8 bottom-8 w-0.5 border-l-2 border-dashed border-red-200 -z-10 transform -translate-x-1/2"></div>
              
              <div className="bg-white border border-gray-300 text-gray-600 rounded-lg px-4 py-3 text-center w-32 text-sm">
                --自动审核
              </div>
              
              <div className="bg-red-500 text-white rounded-lg px-4 py-3 text-center w-32 text-sm shadow-lg shadow-red-200 font-bold relative">
                人工审核
                <div className="absolute -right-2 top-1/2 w-2 h-2 bg-red-500 transform rotate-45 -translate-y-1/2"></div>
              </div>
              
              <div className="bg-white border border-gray-300 text-gray-600 rounded-lg px-4 py-3 text-center w-32 text-sm">
                审核完成
              </div>
           </div>
        </div>

        {/* Center: Content */}
        <div className="flex-1 bg-white rounded-xl shadow-sm border border-gray-100 flex flex-col overflow-hidden">
           <div className="flex border-b border-gray-800 bg-gray-900 text-gray-400">
             <div className="px-6 py-3 text-white border-b-2 border-blue-500 font-medium">工单信息</div>
             <div className="px-6 py-3 hover:text-white cursor-pointer">审核预览</div>
           </div>
           
           <div className="flex-1 bg-gray-800 p-8 relative flex flex-col items-center justify-center">
              {/* Modal-like view for document */}
              <div className="bg-white rounded-lg w-full max-w-3xl h-full flex flex-col overflow-hidden shadow-2xl">
                 <div className="h-12 bg-gray-800 flex items-center justify-between px-4 text-white">
                   <span className="font-medium">资质文件查看</span>
                   <button className="text-gray-400 hover:text-white">×</button>
                 </div>
                 
                 <div className="flex-1 bg-gray-100 p-4 flex items-center justify-center relative">
                    <img 
                      src={images[currentImageIndex].url} 
                      alt={images[currentImageIndex].title} 
                      className="max-h-full max-w-full object-contain shadow-lg"
                    />
                    
                    <button 
                      onClick={() => setCurrentImageIndex(i => (i > 0 ? i - 1 : images.length - 1))}
                      className="absolute left-4 p-2 bg-white/80 rounded-full hover:bg-white shadow"
                    >
                      <ChevronLeft />
                    </button>
                    <button 
                      onClick={() => setCurrentImageIndex(i => (i < images.length - 1 ? i + 1 : 0))}
                      className="absolute right-4 p-2 bg-white/80 rounded-full hover:bg-white shadow"
                    >
                      <ChevronRight />
                    </button>
                 </div>
                 
                 <div className="h-32 bg-white border-t border-gray-200 p-4 overflow-x-auto whitespace-nowrap">
                    {images.map((img, idx) => (
                      <div 
                        key={img.id}
                        onClick={() => setCurrentImageIndex(idx)}
                        className={clsx(
                          "inline-block w-32 h-24 mr-4 border-2 rounded cursor-pointer overflow-hidden relative",
                          currentImageIndex === idx ? "border-blue-500" : "border-transparent"
                        )}
                      >
                        <img src={img.url} className="w-full h-full object-cover" />
                        <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white text-xs p-1 truncate">
                          {img.title}
                        </div>
                      </div>
                    ))}
                 </div>
              </div>
           </div>
        </div>

        {/* Right: Audit Form */}
        <div className="w-80 bg-gray-900 text-white rounded-xl shadow-sm p-6 flex flex-col">
           <div className="mb-6">
             <span className="px-3 py-1 bg-blue-600 rounded-full text-xs font-bold">审核流程</span>
           </div>
           
           <div className="space-y-6 flex-1">
             <div>
               <label className="block text-sm font-medium text-gray-300 mb-2">审核意见</label>
               <div className="bg-white rounded p-1">
                 <input 
                   type="text" 
                   className="w-full px-3 py-2 text-gray-900 text-sm border-none focus:outline-none" 
                   placeholder="审核文本铺"
                 />
               </div>
               <p className="text-xs text-gray-500 mt-2">
                 审核方部提交参考选递逻辑中，参全文本，专达委提交支持点击大大看组
               </p>
             </div>
             
             <div>
                <textarea
                  className="w-full h-32 bg-white rounded p-3 text-gray-900 text-sm focus:outline-none"
                  placeholder="审核点击最大查看"
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                ></textarea>
             </div>
             
             <div>
               <label className="block text-sm font-medium text-gray-300 mb-4">审核结果</label>
               <div className="space-y-3">
                 <label className="flex items-center space-x-3 cursor-pointer">
                   <div className={clsx("w-5 h-5 rounded-full border-2 flex items-center justify-center", auditStatus === 'approved' ? "border-white" : "border-gray-500")}>
                     {auditStatus === 'approved' && <div className="w-3 h-3 bg-white rounded-full"></div>}
                   </div>
                   <input type="radio" className="hidden" name="status" onChange={() => setAuditStatus('approved')} />
                   <span className="text-sm">通/拆回</span>
                 </label>
                 
                 <label className="flex items-center space-x-3 cursor-pointer">
                   <div className={clsx("w-5 h-5 rounded-full border-2 flex items-center justify-center", auditStatus === 'rejected' ? "border-white" : "border-gray-500")}>
                     {auditStatus === 'rejected' && <div className="w-3 h-3 bg-white rounded-full"></div>}
                   </div>
                   <input type="radio" className="hidden" name="status" onChange={() => setAuditStatus('rejected')} />
                   <span className="text-sm">退回审核上一步</span>
                 </label>
               </div>
             </div>
           </div>
           
           <div className="mt-8 flex space-x-4">
             <button 
                onClick={() => handleSubmit('approved')}
                className="flex-1 py-3 bg-green-500 hover:bg-green-600 rounded text-white font-medium transition-colors"
             >
               通过
             </button>
             <button 
                onClick={() => handleSubmit('rejected')}
                className="flex-1 py-3 bg-red-500 hover:bg-red-600 rounded text-white font-medium transition-colors"
             >
               退回
             </button>
           </div>
        </div>
      </div>
    </div>
  );
};

export default AuditDetail;
