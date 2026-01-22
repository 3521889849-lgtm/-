import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import MerchantApply from './pages/MerchantApply';
import ApplicationList from './pages/ApplicationList';
import AuditDetail from './pages/AuditDetail';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Navigate to="/merchant/apply" replace />} />
          <Route path="/merchant/apply" element={<MerchantApply />} />
          <Route path="/merchant/applications" element={<ApplicationList />} />
          <Route path="/merchant/audit/:id" element={<AuditDetail />} />
          <Route path="*" element={<Navigate to="/merchant/apply" replace />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
