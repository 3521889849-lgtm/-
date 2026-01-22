import React, { useState } from 'react';
import { Upload, ChevronDown } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { addApplication } from '../lib/mockStore';

const MerchantApply: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    merchant_name: '',
    contact_person: '',
    contact_phone: '',
  });
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    try {
      if (!formData.merchant_name) {
          alert('请填写商家名称');
          return;
      }
      setLoading(true);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      addApplication(formData);
      
      alert('提交申请成功');
      navigate('/merchant/applications');
    } catch (error) {
      console.error(error);
      alert('提交失败');
    } finally {
      setLoading(false);
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-gray-900">票务系统商家入驻原型图</h2>
        <p className="text-gray-500 mt-1">入驻申请</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Basic Info Card */}
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <h3 className="text-lg font-bold text-gray-800 mb-6">基本信息组</h3>
          
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-500 mb-2">商家名称:</label>
              <input
                type="text"
                className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={formData.merchant_name}
                onChange={(e) => setFormData({...formData, merchant_name: e.target.value})}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-2">联系人:</label>
              <div className="relative">
                <input
                  type="text"
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 pr-10"
                  value={formData.contact_person}
                  onChange={(e) => setFormData({...formData, contact_person: e.target.value})}
                />
                <ChevronDown className="absolute right-3 top-3 text-gray-400" size={16} />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-2">联系电话:</label>
              <div className="relative">
                <input
                  type="text"
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 pr-10"
                  value={formData.contact_phone}
                  onChange={(e) => setFormData({...formData, contact_phone: e.target.value})}
                />
                <ChevronDown className="absolute right-3 top-3 text-gray-400" size={16} />
              </div>
            </div>
          </div>
        </div>

        {/* Qualification Info Card */}
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100 flex flex-col">
          <h3 className="text-lg font-bold text-gray-800 mb-6">资质信息组</h3>
          
          <div className="flex-1 flex flex-col items-center justify-center border-2 border-dashed border-gray-200 rounded-lg bg-gray-50 p-8">
            <div className="w-16 h-16 bg-white rounded-lg shadow-sm flex items-center justify-center mb-4">
              <Upload className="text-blue-500" size={32} />
            </div>
            <p className="text-gray-900 font-medium mb-2">营业授权</p>
            
            <div className="flex items-center mt-4 mb-6">
              <input type="checkbox" className="mr-2" />
              <span className="text-sm text-gray-500">经营许可证上传区</span>
            </div>

            <label className="cursor-pointer w-full">
              <input type="file" className="hidden" onChange={handleFileChange} />
              <div className="w-full bg-gradient-to-r from-orange-400 to-red-500 text-white font-bold py-3 rounded-full text-center hover:opacity-90 transition-opacity">
                {file ? '已选择: ' + file.name : '上传'}
              </div>
            </label>
          </div>
        </div>
      </div>

      <div className="flex justify-center mt-12 space-x-6">
        <button
          onClick={handleSubmit}
          disabled={loading}
          className="px-12 py-3 bg-slate-900 text-white rounded-full font-medium hover:bg-slate-800 transition-colors disabled:opacity-50"
        >
          {loading ? '提交中...' : '提交申请'}
        </button>
        <button className="px-12 py-3 bg-white text-gray-800 border border-gray-300 rounded-full font-medium hover:bg-gray-50 transition-colors">
          保存草稿
        </button>
      </div>
    </div>
  );
};

export default MerchantApply;
