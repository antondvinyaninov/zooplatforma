'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import AdminLayout, { AdminTab } from '../components/admin/AdminLayout';
import TableWidget from '../components/admin/widgets/TableWidget';
import {
  HomeIcon,
  ClipboardDocumentListIcon,
  ChartBarIcon,
  PlusIcon,
} from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'overview', label: 'Обзор', icon: <ChartBarIcon className="w-5 h-5" /> },
  { id: 'pets', label: 'Мои подопечные', icon: <HomeIcon className="w-5 h-5" /> },
  { id: 'tasks', label: 'Задачи', icon: <ClipboardDocumentListIcon className="w-5 h-5" /> },
];

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  photo?: string;
  status: string;
  created_at: string;
}

export default function VolunteerDashboard() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('overview');
  const [isClient, setIsClient] = useState(false);
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<{ email: string; name?: string; avatar?: string; role: string } | null>(null);

  useEffect(() => {
    setIsClient(true);
    const savedTab = localStorage.getItem('volunteerActiveTab');
    if (savedTab) {
      setActiveTab(savedTab);
    }
  }, []);

  useEffect(() => {
    if (isClient) {
      loadData();
    }
  }, [activeTab, isClient]);

  const loadData = async () => {
    setLoading(true);

    try {
      // Проверяем авторизацию через Auth Service
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      if (!meResponse.ok) {
        console.error('Auth check failed with status:', meResponse.status);
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

      // Загружаем подопечных питомцев
      const petsResponse = await fetch('http://localhost:8500/api/my-pets', {
        method: 'GET',
        credentials: 'include',
      });

      if (petsResponse.ok) {
        const petsResult = await petsResponse.json();
        if (petsResult.success && petsResult.data) {
          setPets(petsResult.data);
        }
      } else {
        console.error('Failed to load pets, status:', petsResponse.status);
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      router.push('/auth');
      return;
    }

    setLoading(false);
  };

  const handleLogout = () => {
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; domain=localhost';
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    router.push('/auth');
  };

  const handleTabChange = (tabId: string) => {
    setActiveTab(tabId);
    localStorage.setItem('volunteerActiveTab', tabId);
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
      activeTab={activeTab}
      onTabChange={handleTabChange}
      adminUser={user}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      {activeTab === 'overview' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-4xl font-bold text-gray-900 mb-2">Добро пожаловать!</h2>
            <p className="text-base text-gray-600">
              {user.name || user.email}, спасибо за помощь животным
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <HomeIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{pets.length}</div>
              </div>
              <div className="text-base font-medium opacity-90">Моих подопечных</div>
            </div>
            
            <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <ClipboardDocumentListIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">0</div>
              </div>
              <div className="text-base font-medium opacity-90">Активных задач</div>
            </div>
            
            <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <ChartBarIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">0</div>
              </div>
              <div className="text-base font-medium opacity-90">Выполнено за месяц</div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <h3 className="text-2xl font-bold text-gray-900 mb-4">О кабинете</h3>
              <div className="space-y-4 text-gray-600 text-base leading-relaxed">
                <p>
                  Кабинет зоопомощника - это ваш инструмент для координации помощи животным.
                  Здесь вы можете отслеживать своих подопечных, управлять задачами
                  и взаимодействовать с приютами.
                </p>
                <p>
                  Спасибо за вашу помощь! Каждое доброе дело делает мир лучше.
                </p>
              </div>
            </div>

            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <h3 className="text-2xl font-bold text-gray-900 mb-6">Возможности</h3>
              <div className="space-y-4">
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">🐾</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Подопечные</div>
                    <div className="text-sm text-gray-600">Животные под вашей опекой</div>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">📋</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Задачи</div>
                    <div className="text-sm text-gray-600">Планирование и отчетность</div>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">🏠</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Приюты</div>
                    <div className="text-sm text-gray-600">Связь с организациями</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'pets' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Мои подопечные</h2>
            <p className="text-gray-600">Животные, за которыми вы ухаживаете</p>
          </div>

          <TableWidget
            title={`Подопечных (${pets.length})`}
            actions={
              <button
                onClick={() => router.push('/dashboard/pets/create')}
                className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
              >
                <PlusIcon className="w-4 h-4" />
                Добавить животное
              </button>
            }
          >
            {loading ? (
              <div className="text-center py-12 text-gray-400">Загрузка...</div>
            ) : pets.length === 0 ? (
              <div className="text-center py-12 text-gray-400">
                <HomeIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Подопечных пока нет</p>
                <p className="text-sm mt-2">Возьмите животное под опеку из каталога</p>
              </div>
            ) : (
              <div className="overflow-x-auto p-6">
                <table className="w-full min-w-[800px]">
                  <thead>
                    <tr className="border-b border-gray-100">
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Фото</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Кличка</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Вид</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Порода</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Статус</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Действия</th>
                    </tr>
                  </thead>
                  <tbody>
                    {pets.map((pet) => (
                      <tr key={pet.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                        <td className="py-3 px-4">
                          {pet.photo ? (
                            <img
                              src={pet.photo}
                              alt={pet.name}
                              className="w-12 h-12 rounded-full object-cover"
                            />
                          ) : (
                            <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-xl">
                              🐾
                            </div>
                          )}
                        </td>
                        <td className="py-3 px-4 text-sm font-medium text-gray-900">{pet.name}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{pet.species}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{pet.breed || '-'}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{pet.status}</td>
                        <td className="py-3 px-4 text-sm">
                          <button
                            onClick={() => window.location.href = `http://localhost:3000/pets/${pet.id}`}
                            className="text-blue-600 hover:text-blue-700"
                          >
                            Открыть
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </TableWidget>
        </div>
      )}

      {activeTab === 'tasks' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Задачи</h2>
            <p className="text-gray-600">Планирование и отчетность</p>
          </div>

          <TableWidget
            title="Задачи (0)"
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2">
                <PlusIcon className="w-4 h-4" />
                Создать задачу
              </button>
            }
          >
            <div className="text-center py-12 text-gray-400">
              <ClipboardDocumentListIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>Задач пока нет</p>
              <p className="text-sm mt-2">Создайте первую задачу для планирования работы</p>
            </div>
          </TableWidget>
        </div>
      )}
    </AdminLayout>
  );
}
