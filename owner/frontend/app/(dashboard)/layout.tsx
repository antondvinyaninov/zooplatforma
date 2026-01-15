'use client';

import { ReactNode, useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import AdminLayout, { AdminTab } from '../components/admin/AdminLayout';
import {
  HomeIcon,
  HeartIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'overview', label: 'Обзор', icon: <ChartBarIcon className="w-5 h-5" /> },
  { id: 'pets', label: 'Мои питомцы', icon: <HomeIcon className="w-5 h-5" /> },
  { id: 'health', label: 'Здоровье', icon: <HeartIcon className="w-5 h-5" /> },
];

export default function DashboardLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [activeTab, setActiveTab] = useState('overview');
  const [user, setUser] = useState<{ email: string; name?: string; role: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  useEffect(() => {
    // Определяем активную вкладку по URL
    if (pathname === '/overview' || pathname === '/') {
      setActiveTab('overview');
    } else if (pathname.startsWith('/pets')) {
      setActiveTab('pets');
    } else if (pathname === '/health') {
      setActiveTab('health');
    }
  }, [pathname]);

  const checkAuth = async () => {
    try {
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        credentials: 'include',
      });

      if (!meResponse.ok) {
        router.push('/auth');
        return;
      }

      const meResult = await meResponse.json();
      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      setUser({
        email: meResult.data.email,
        name: meResult.data.name,
        role: 'owner',
      });
    } catch (error) {
      console.error('Auth check failed:', error);
      router.push('/auth');
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (tabId: string) => {
    setActiveTab(tabId);
    
    if (tabId === 'overview') {
      router.push('/overview');
    } else if (tabId === 'pets') {
      router.push('/pets');
    } else if (tabId === 'health') {
      router.push('/health');
    }
  };

  const handleLogout = async () => {
    await fetch('http://localhost:8000/api/auth/logout', {
      method: 'POST',
      credentials: 'include',
    });
    router.push('/auth');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <div className="text-lg text-gray-600">Проверка доступа...</div>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <AdminLayout
      logoSrc="/favicon.svg"
      logoText="Кабинет владельца"
      logoAlt="Кабинет владельца"
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={handleTabChange}
      adminUser={user}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      {children}
    </AdminLayout>
  );
}
