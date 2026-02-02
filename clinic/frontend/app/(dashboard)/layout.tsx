'use client';

import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import AdminLayout, { AdminTab } from '../components/admin/AdminLayout';
import {
  UserGroupIcon,
  CalendarIcon,
  ChartBarIcon,
  Cog6ToothIcon,
  ClipboardDocumentCheckIcon,
} from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'overview', label: 'Обзор', icon: <ChartBarIcon className="w-5 h-5" />, href: '/overview' },
  { id: 'registration', label: 'Регистрация животных', icon: <ClipboardDocumentCheckIcon className="w-5 h-5" />, href: '/registration' },
  { id: 'patients', label: 'Пациенты', icon: <UserGroupIcon className="w-5 h-5" />, href: '/patients' },
  { id: 'appointments', label: 'Записи на приём', icon: <CalendarIcon className="w-5 h-5" />, href: '/appointments' },
  { id: 'settings', label: 'Настройки', icon: <Cog6ToothIcon className="w-5 h-5" />, href: '/settings' },
];

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const [user, setUser] = useState<{ email: string; name?: string; avatar?: string; role: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      // Проверяем авторизацию через Main API (SSO)
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
        role: 'clinic',
      });

      // Проверяем выбрана ли клиника
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) {
        router.push('/select');
        return;
      }

      setLoading(false);
    } catch (error) {
      console.error('Auth check failed:', error);
      router.push('/auth');
    }
  };

  const handleLogout = () => {
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; domain=localhost';
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    localStorage.removeItem('selectedClinicId');
    router.push('/auth');
  };

  const handleTabChange = (tabId: string) => {
    const tab = tabs.find(t => t.id === tabId);
    if (tab?.href) {
      router.push(tab.href);
    }
  };

  // Определяем активную вкладку по URL
  const getActiveTab = () => {
    if (pathname === '/overview') return 'overview';
    if (pathname === '/registration') return 'registration';
    if (pathname === '/patients') return 'patients';
    if (pathname === '/appointments') return 'appointments';
    if (pathname === '/settings') return 'settings';
    return 'overview';
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
      logoText="Кабинет ветклиники"
      logoAlt="Кабинет ветклиники"
      tabs={tabs}
      activeTab={getActiveTab()}
      onTabChange={handleTabChange}
      adminUser={user}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      {children}
    </AdminLayout>
  );
}
