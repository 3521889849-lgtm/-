import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, FileText, Search, UserPlus, FileBox, Bell, Settings, User } from 'lucide-react';
import clsx from 'clsx';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const location = useLocation();

  const navItems = [
    { name: '查询筛选', icon: Search, path: '/merchant/applications' },
    { name: '商家入驻', icon: UserPlus, path: '/merchant/apply' },
  ];

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <div className="w-64 bg-slate-900 text-white flex flex-col">
        <div className="p-6 flex items-center space-x-2 border-b border-slate-800">
          <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
            <span className="font-bold">票</span>
          </div>
          <span className="text-xl font-bold">票务系统</span>
        </div>
        
        <nav className="flex-1 py-6">
          <ul className="space-y-1 px-3">
            {navItems.map((item) => (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={clsx(
                    'flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors',
                    location.pathname === item.path
                      ? 'bg-blue-600 text-white'
                      : 'text-slate-400 hover:bg-slate-800 hover:text-white'
                  )}
                >
                  <item.icon size={20} />
                  <span>{item.name}</span>
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="bg-white border-b border-gray-200 h-16 flex items-center justify-between px-6">
          <h1 className="text-xl font-bold text-gray-800">
            {navItems.find(i => i.path === location.pathname)?.name || '票务系统'}
          </h1>
          
          <div className="flex items-center space-x-4">
            <button className="p-2 text-gray-500 hover:bg-gray-100 rounded-full">
              <Bell size={20} />
            </button>
            <button className="p-2 text-gray-500 hover:bg-gray-100 rounded-full">
              <Settings size={20} />
            </button>
            <div className="flex items-center space-x-2 pl-4 border-l border-gray-200">
              <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center overflow-hidden">
                <User size={20} className="text-gray-500" />
              </div>
              <span className="text-sm font-medium text-gray-700">管理员</span>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;
