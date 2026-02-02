'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import AdminLayout, { AdminTab } from '../../components/admin/AdminLayout';
import {
  HomeIcon,
  ClipboardDocumentListIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'overview', label: 'Обзор', icon: <ChartBarIcon className="w-5 h-5" /> },
  { id: 'pets', label: 'Мои подопечные', icon: <HomeIcon className="w-5 h-5" /> },
  { id: 'tasks', label: 'Задачи', icon: <ClipboardDocumentListIcon className="w-5 h-5" /> },
];

export default function PetsLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [user, setUser] = useState<{ email: string; name?: string; avatar?: string; role: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        method: 'GET',
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
        email: meResult.data.user.email,
        name: meResult.data.user.name,
        avatar: meResult.data.user.avatar,
        role: 'volunteer',
      });
    } catch (error) {
      console.error('Auth check failed:', error);
      router.push('/auth');
      return;
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; domain=localhost';
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    router.push('/auth');
  };

  const handleTabChange = (tabId: string) => {
    router.push('/dashboard');
    if (typeof window !== 'undefined') {
      localStorage.setItem('volunteerActiveTab', tabId);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
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
      logoText="Кабинет зоопомощника"
      logoAlt="Кабинет зоопомощника"
      tabs={tabs}
      activeTab="pets"
      onTabChange={handleTabChange}
      adminUser={user}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      {children}
    </AdminLayout>
  );
}
