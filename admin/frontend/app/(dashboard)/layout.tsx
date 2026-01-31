'use client';

import { useRouter, usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import AdminLayout, { AdminTab } from '../components/admin/AdminLayout';
import {
  UsersIcon,
  DocumentTextIcon,
  ChartBarIcon,
  DocumentDuplicateIcon,
  BuildingOfficeIcon,
  ServerIcon,
  FlagIcon,
} from '@heroicons/react/24/outline';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const [adminUser, setAdminUser] = useState<{ email: string; name?: string; avatar?: string; role: string } | null>(null);
  const [activeTab, setActiveTab] = useState('dashboard');

  useEffect(() => {
    // Проверка авторизации через Auth Service
    const checkAuth = async () => {
      try {
        const response = await fetch('http://localhost:7100/api/auth/me', {
          credentials: 'include',
        });

        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data && data.data.user) {
            // Проверяем что пользователь - superadmin
            if (data.data.user.role !== 'superadmin') {
              alert('Доступ запрещен. Требуются права суперадмина.');
              router.push('http://localhost:3000');
              return;
            }
            setAdminUser({
              email: data.data.user.email,
              name: data.data.user.name,
              avatar: data.data.user.avatar,
              role: data.data.user.role,
            });
          } else {
            router.push('/auth');
          }
        } else {
          router.push('/auth');
        }
      } catch (error) {
        console.error('Auth check failed:', error);
        router.push('/auth');
      }
    };

    checkAuth();
  }, [router]);

  useEffect(() => {
    // Определяем активную вкладку по URL
    if (pathname.includes('/dashboard')) {
      setActiveTab('dashboard');
    } else if (pathname.includes('/posts')) {
      setActiveTab('posts');
    } else if (pathname.includes('/logs')) {
      setActiveTab('logs');
    } else if (pathname.includes('/monitoring')) {
      setActiveTab('health');
    } else if (pathname.includes('/organizations')) {
      setActiveTab('organizations');
    } else if (pathname.includes('/moderation')) {
      setActiveTab('moderation');
    } else if (pathname.includes('/users')) {
      setActiveTab('users');
    } else {
      setActiveTab('dashboard'); // По умолчанию
    }
  }, [pathname]);

  const tabs: AdminTab[] = [
    {
      id: 'dashboard',
      label: 'Дашборд',
      icon: <ChartBarIcon className="w-5 h-5" />,
    },
    {
      id: 'users',
      label: 'Пользователи',
      icon: <UsersIcon className="w-5 h-5" />,
    },
    {
      id: 'posts',
      label: 'Посты',
      icon: <DocumentTextIcon className="w-5 h-5" />,
    },
    {
      id: 'moderation',
      label: 'Модерация',
      icon: <FlagIcon className="w-5 h-5" />,
    },
    {
      id: 'logs',
      label: 'Логирование',
      icon: <DocumentDuplicateIcon className="w-5 h-5" />,
    },
    {
      id: 'organizations',
      label: 'Организации',
      icon: <BuildingOfficeIcon className="w-5 h-5" />,
    },
    {
      id: 'health',
      label: 'Мониторинг',
      icon: <ServerIcon className="w-5 h-5" />,
    },
  ];

  const handleTabChange = (tabId: string) => {
    setActiveTab(tabId);
    
    // Навигация по табам
    const routes: Record<string, string> = {
      dashboard: '/dashboard',
      users: '/users',
      posts: '/posts',
      moderation: '/moderation',
      logs: '/logs',
      organizations: '/organizations',
      health: '/monitoring',
    };

    if (routes[tabId]) {
      router.push(routes[tabId]);
    }
  };

  const handleLogout = async () => {
    try {
      await fetch('http://localhost:7100/api/auth/logout', {
        method: 'POST',
        credentials: 'include',
      });
      router.push('/auth');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  if (!adminUser) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Загрузка...</div>
      </div>
    );
  }

  return (
    <AdminLayout
      logoSrc="/logo.svg"
      logoText="ЗооАдминка"
      logoAlt="ЗооПлатформа Админка"
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={handleTabChange}
      adminUser={adminUser}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      {children}
    </AdminLayout>
  );
}
