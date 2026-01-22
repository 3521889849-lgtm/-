export interface Application {
  id: string;
  merchant_name: string;
  contact_person: string;
  contact_phone: string;
  created_at: string;
  updated_at: string;
  status: 'pending' | 'passed' | 'rejected';
  audit_comment?: string;
  id_number: number;
}

const STORAGE_KEY = 'mock_applications';

export const getApplications = (): Application[] => {
  const data = localStorage.getItem(STORAGE_KEY);
  if (data) {
    return JSON.parse(data);
  }
  
  const initialData: Application[] = [
    { 
      id: '1', 
      merchant_name: '示例票务公司A', 
      contact_person: '张三', 
      contact_phone: '13800138000', 
      created_at: '2022-12-14 17:03', 
      updated_at: '2022-11-00', 
      status: 'passed', 
      id_number: 32 
    },
    { 
      id: '2', 
      merchant_name: '示例票务公司B', 
      contact_person: '李四', 
      contact_phone: '13900139000', 
      created_at: '2022-12-13 17:03', 
      updated_at: '2022-11-02', 
      status: 'pending', 
      id_number: 30 
    },
    { 
      id: '3', 
      merchant_name: '示例票务公司C', 
      contact_person: '王五', 
      contact_phone: '13700137000', 
      created_at: '2020-12-14 17:03', 
      updated_at: '2022-11-00', 
      status: 'pending', 
      id_number: 21 
    },
    { 
      id: '4', 
      merchant_name: '示例票务公司D', 
      contact_person: '赵六', 
      contact_phone: '13600136000', 
      created_at: '2023-12-14 17:03', 
      updated_at: '2022-11-00', 
      status: 'pending', 
      id_number: 32 
    },
    { 
      id: '5', 
      merchant_name: '示例票务公司E', 
      contact_person: '孙七', 
      contact_phone: '13500135000', 
      created_at: '1918-12-12 17:03', 
      updated_at: '2022-11-00', 
      status: 'rejected', 
      id_number: 99 
    },
  ];
  localStorage.setItem(STORAGE_KEY, JSON.stringify(initialData));
  return initialData;
};

export const addApplication = (app: { merchant_name: string; contact_person: string; contact_phone: string }) => {
  const apps = getApplications();
  const now = new Date();
  const formattedDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`;
  
  const newApp: Application = {
    ...app,
    id: Date.now().toString(),
    created_at: formattedDate,
    updated_at: formattedDate,
    status: 'pending',
    id_number: Math.floor(Math.random() * 100),
  };
  
  apps.unshift(newApp);
  localStorage.setItem(STORAGE_KEY, JSON.stringify(apps));
  return newApp;
};

export const updateApplicationStatus = (id: string, status: 'passed' | 'rejected', comment: string) => {
  const apps = getApplications();
  const index = apps.findIndex(a => a.id === id);
  if (index !== -1) {
    apps[index].status = status;
    apps[index].audit_comment = comment;
    const now = new Date();
    apps[index].updated_at = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(apps));
    return apps[index];
  }
  return null;
};

export const getApplicationById = (id: string) => {
    const apps = getApplications();
    return apps.find(a => a.id === id);
};
