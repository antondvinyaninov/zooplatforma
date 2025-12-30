'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import AdminLayout, { AdminTab } from '../components/admin/AdminLayout';
import TableWidget from '../components/admin/widgets/TableWidget';
import { shelterApi } from '@/lib/api';
import {
  HomeIcon,
  HeartIcon,
  UserGroupIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'stats', label: 'Статистика', icon: <ChartBarIcon className="w-5 h-5" /> },
  { id: 'animals', label: 'Животные приюта', icon: <HomeIcon className="w-5 h-5" /> },
  { id: 'adoption', label: 'Пристройство', icon: <HeartIcon className="w-5 h-5" /> },
  { id: 'volunteers', label: 'Волонтеры', icon: <UserGroupIcon className="w-5 h-5" /> },
];

interface ShelterAnimal {
  id: number;
  name: string;
  species: string;
  breed?: string;
  age?: string;
  gender?: string;
  status: string;
  arrival_date: string;
  photo?: string;
}

export default function ShelterDashboard() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('stats');
  const [isClient, setIsClient] = useState(false);
  const [animals, setAnimals] = useState<ShelterAnimal[]>([]);
  const [loading, setLoading] = useState(true);
  const [adminUser, setAdminUser] = useState<{ email: string; role: string } | null>(null);
  const [stats, setStats] = useState({ total_animals: 0, adopted_this_month: 0, active_volunteers: 0, pending_requests: 0 });
  const [organization, setOrganization] = useState<{ name: string } | null>(null);

  useEffect(() => {
    setIsClient(true);
    const savedTab = localStorage.getItem('shelterActiveTab');
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

    // Проверяем авторизацию через Admin API (SSO)
    try {
      const meResponse = await fetch('http://localhost:9000/api/admin/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      // Проверяем роль (shelter_admin или выше)
      const allowedRoles = ['shelter_admin', 'moderator', 'admin', 'superadmin'];
      if (!allowedRoles.includes(meResult.data?.role)) {
        router.push('/auth');
        return;
      }

      setAdminUser(meResult.data);

      // Загружаем информацию о приюте
      const orgResult = await shelterApi.getOrganization();
      if (orgResult.success && orgResult.data) {
        setOrganization(orgResult.data);
      }

      // Загружаем статистику
      const statsResult = await shelterApi.getStats();
      if (statsResult.success && statsResult.data) {
        setStats(statsResult.data);
      }

      // Загружаем данные в зависимости от вкладки
      if (activeTab === 'animals') {
        const animalsResult = await shelterApi.getAnimals();
        if (animalsResult.success && animalsResult.data) {
          setAnimals(animalsResult.data);
        }
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
    localStorage.setItem('shelterActiveTab', tabId);
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

  if (!adminUser) {
    return null;
  }

  return (
    <AdminLayout
      logoSrc="/favicon.svg"
      logoText={organization?.name || "Приют"}
      logoAlt="Кабинет приюта"
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={handleTabChange}
      adminUser={adminUser}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      {activeTab === 'stats' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-4xl font-bold text-gray-900 mb-2">Статистика приюта</h2>
            <p className="text-base text-gray-600">Общая информация о животных и деятельности</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <HomeIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{stats.total_animals}</div>
              </div>
              <div className="text-base font-medium opacity-90">Животных в приюте</div>
            </div>
            
            <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <HeartIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{stats.adopted_this_month}</div>
              </div>
              <div className="text-base font-medium opacity-90">Пристроено за месяц</div>
            </div>
            
            <div className="bg-gradient-to-br from-orange-500 to-orange-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <UserGroupIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{stats.active_volunteers}</div>
              </div>
              <div className="text-base font-medium opacity-90">Активных волонтеров</div>
            </div>
            
            <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <ChartBarIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{stats.pending_requests}</div>
              </div>
              <div className="text-base font-medium opacity-90">Заявок на пристройство</div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <h3 className="text-2xl font-bold text-gray-900 mb-4">О приюте</h3>
              <div className="space-y-4 text-gray-600 text-base leading-relaxed">
                <p>
                  Кабинет приюта - это инструмент для управления животными, находящимися в приюте,
                  координации волонтеров и обработки заявок на пристройство.
                </p>
                <p>
                  Здесь вы можете вести учет животных, отслеживать их статус, публиковать объявления
                  о поиске дома и взаимодействовать с потенциальными хозяевами.
                </p>
              </div>
            </div>

            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <h3 className="text-2xl font-bold text-gray-900 mb-6">Возможности</h3>
              <div className="space-y-4">
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">🏠</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Учет животных</div>
                    <div className="text-sm text-gray-600">Полная информация о каждом питомце</div>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">💚</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Пристройство</div>
                    <div className="text-sm text-gray-600">Поиск новых хозяев для питомцев</div>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">👥</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Волонтеры</div>
                    <div className="text-sm text-gray-600">Координация помощников приюта</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'animals' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Животные приюта</h2>
            <p className="text-gray-600">Управление питомцами, находящимися в приюте</p>
          </div>

          <TableWidget
            title={`Животных (${animals.length})`}
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Добавить животное
              </button>
            }
          >
            {loading ? (
              <div className="text-center py-12 text-gray-400">Загрузка...</div>
            ) : animals.length === 0 ? (
              <div className="text-center py-12 text-gray-400">
                <HomeIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Животные не найдены</p>
                <p className="text-sm mt-2">Добавьте первое животное, чтобы начать работу</p>
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
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Возраст</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Статус</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Дата поступления</th>
                    </tr>
                  </thead>
                  <tbody>
                    {animals.map((animal) => (
                      <tr key={animal.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                        <td className="py-3 px-4">
                          {animal.photo ? (
                            <img
                              src={animal.photo}
                              alt={animal.name}
                              className="w-12 h-12 rounded-full object-cover"
                            />
                          ) : (
                            <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-xl">
                              🐾
                            </div>
                          )}
                        </td>
                        <td className="py-3 px-4 text-sm font-medium text-gray-900">{animal.name}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{animal.species}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{animal.breed || '-'}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{animal.age || '-'}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{animal.status}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">
                          {new Date(animal.arrival_date).toLocaleDateString('ru-RU')}
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

      {activeTab === 'adoption' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Пристройство</h2>
            <p className="text-gray-600">Заявки на пристройство животных</p>
          </div>

          <TableWidget
            title="Заявки (0)"
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Создать объявление
              </button>
            }
          >
            <div className="text-center py-12 text-gray-400">
              <HeartIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>Заявок пока нет</p>
              <p className="text-sm mt-2">Создайте объявление о поиске дома для питомца</p>
            </div>
          </TableWidget>
        </div>
      )}

      {activeTab === 'volunteers' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Волонтеры</h2>
            <p className="text-gray-600">Управление волонтерами приюта</p>
          </div>

          <TableWidget
            title="Волонтеры (0)"
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Добавить волонтера
              </button>
            }
          >
            <div className="text-center py-12 text-gray-400">
              <UserGroupIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>Волонтеры не найдены</p>
              <p className="text-sm mt-2">Добавьте первого волонтера</p>
            </div>
          </TableWidget>
        </div>
      )}
    </AdminLayout>
  );
}
