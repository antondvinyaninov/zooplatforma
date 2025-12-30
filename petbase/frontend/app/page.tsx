'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import AdminLayout, { AdminTab } from './components/admin/AdminLayout';
import StatsWidget from './components/admin/widgets/StatsWidget';
import TableWidget from './components/admin/widgets/TableWidget';
import {
  Square3Stack3DIcon,
  RectangleStackIcon,
  DocumentTextIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';
import { petbaseApi } from '@/lib/api';

const tabs: AdminTab[] = [
  { id: 'stats', label: 'Статистика', icon: <ChartBarIcon className="w-5 h-5" /> },
  { id: 'species', label: 'Виды', icon: <Square3Stack3DIcon className="w-5 h-5" /> },
  { id: 'breeds', label: 'Породы', icon: <RectangleStackIcon className="w-5 h-5" /> },
  { id: 'pets', label: 'Питомцы', icon: <Square3Stack3DIcon className="w-5 h-5" /> },
];

interface Species {
  id: number;
  name: string;
  name_en: string;
  description: string;
  icon: string;
  created_at: string;
}

interface Breed {
  id: number;
  species_id: number;
  species_name: string;
  name: string;
  name_en: string;
  description: string;
  origin: string;
  size: string;
  weight_min: number;
  weight_max: number;
  lifespan_min: number;
  lifespan_max: number;
  temperament: string;
  care_level: string;
  photo: string;
  created_at: string;
}

interface PetCard {
  id: number;
  breed_id: number;
  breed_name: string;
  title: string;
  description: string;
  characteristics: string;
  care_tips: string;
  health_info: string;
  nutrition: string;
  photos: string;
  is_published: boolean;
  created_at: string;
  updated_at: string;
}

interface Pet {
  id: number;
  user_id: number;
  name: string;
  species: string;
  breed?: string;
  gender?: string;
  birth_date?: string;
  color?: string;
  size?: string;
  weight?: number;
  chip_number?: string;
  passport_number?: string;
  is_sterilized: boolean;
  is_vaccinated: boolean;
  health_notes?: string;
  character_traits?: string;
  special_needs?: string;
  status: string;
  status_updated_at?: string;
  photo?: string;
  photos?: string;
  story?: string;
  created_at: string;
  updated_at: string;
}

export default function PetBaseDashboard() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('stats');
  const [isClient, setIsClient] = useState(false);
  const [species, setSpecies] = useState<Species[]>([]);
  const [breeds, setBreeds] = useState<Breed[]>([]);
  const [cards, setCards] = useState<PetCard[]>([]);
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  const [adminUser, setAdminUser] = useState<{ email: string; role: string } | null>(null);

  useEffect(() => {
    setIsClient(true);
    const savedTab = localStorage.getItem('petbaseActiveTab');
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

    // Проверяем авторизацию
    const authResult = await petbaseApi.checkAuth();
    if (!authResult.success) {
      router.push('/auth');
      return;
    }

    setAdminUser(authResult.data);

    // Загружаем данные в зависимости от вкладки
    if (activeTab === 'species') {
      const result = await petbaseApi.getSpecies();
      if (result.success) {
        setSpecies(result.data || []);
      }
    } else if (activeTab === 'breeds') {
      const result = await petbaseApi.getBreeds();
      if (result.success) {
        setBreeds(result.data || []);
      }
    } else if (activeTab === 'pets') {
      // Загружаем питомцев
      // Временно для разработки: используем X-User-ID заголовок
      // TODO: Заменить на JWT токен после внедрения полной аутентификации
      const response = await fetch('http://localhost:8100/api/pets', {
        headers: {
          'X-User-ID': '1', // Временно: ID администратора
        },
      });
      const result = await response.json();
      if (result.success) {
        setPets(result.data || []);
      }
    } else if (activeTab === 'stats') {
      // Загружаем все для статистики
      const [speciesResult, breedsResult] = await Promise.all([
        petbaseApi.getSpecies(),
        petbaseApi.getBreeds(),
      ]);
      if (speciesResult.success) setSpecies(speciesResult.data || []);
      if (breedsResult.success) setBreeds(breedsResult.data || []);
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
    localStorage.setItem('petbaseActiveTab', tabId);
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
      logoText="ЗооБаза"
      logoAlt="ЗооБаза"
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
            <h2 className="text-4xl font-bold text-gray-900 mb-2">Статистика ЗооБазы</h2>
            <p className="text-base text-gray-600">Общая информация о справочнике</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <Square3Stack3DIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{species.length}</div>
              </div>
              <div className="text-base font-medium opacity-90">Видов животных</div>
            </div>
            
            <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <RectangleStackIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{breeds.length}</div>
              </div>
              <div className="text-base font-medium opacity-90">Пород</div>
            </div>
            
            <div className="bg-gradient-to-br from-orange-500 to-orange-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <DocumentTextIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{pets.length}</div>
              </div>
              <div className="text-base font-medium opacity-90">Питомцев</div>
            </div>
            
            <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
              <div className="flex items-center justify-between mb-4">
                <ChartBarIcon className="w-8 h-8 opacity-80" />
                <div className="text-5xl font-bold">{cards.filter((c) => c.is_published).length}</div>
              </div>
              <div className="text-base font-medium opacity-90">Опубликовано</div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <h3 className="text-2xl font-bold text-gray-900 mb-4">О ЗооБазе</h3>
              <div className="space-y-4 text-gray-600 text-base leading-relaxed">
                <p>
                  ЗооБаза - это централизованный справочник домашних животных, который содержит
                  подробную информацию о видах, породах и особенностях содержания.
                </p>
                <p>
                  Данные из ЗооБазы используются другими сервисами платформы для автозаполнения
                  информации о питомцах пользователей и предоставления справочной информации.
                </p>
              </div>
            </div>

            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <h3 className="text-2xl font-bold text-gray-900 mb-6">Возможности</h3>
              <div className="space-y-4">
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">🔌</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">REST API</div>
                    <div className="text-sm text-gray-600">Доступ к данным через API</div>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">🔍</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Быстрый поиск</div>
                    <div className="text-sm text-gray-600">Поиск по видам и породам</div>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                    <span className="text-2xl">📸</span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1 text-base">Галерея</div>
                    <div className="text-sm text-gray-600">Фотографии и изображения</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl p-8 text-white shadow-lg">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-3xl font-bold mb-2">Подписка на API</h3>
                <p className="text-indigo-100 mb-4 text-base">
                  Для создания ключей и использования API требуется подписка
                </p>
                <button className="px-6 py-3 bg-white text-indigo-600 font-semibold rounded-xl hover:bg-indigo-50 transition-colors shadow-md text-base">
                  Купить подписку
                </button>
              </div>
              <div className="hidden md:block text-8xl opacity-20">
                🔑
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'species' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Виды животных</h2>
            <p className="text-gray-600">Управление видами домашних животных</p>
          </div>

          <TableWidget
            title={`Виды (${species.length})`}
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Добавить вид
              </button>
            }
          >
            {loading ? (
              <div className="text-center py-12 text-gray-400">Загрузка...</div>
            ) : species.length === 0 ? (
              <div className="text-center py-12 text-gray-400">
                <Square3Stack3DIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Виды животных не найдены</p>
                <p className="text-sm mt-2">Добавьте первый вид, чтобы начать работу</p>
              </div>
            ) : (
              <div className="overflow-x-auto p-6">
                <table className="w-full min-w-[800px]">
                  <thead>
                    <tr className="border-b border-gray-100">
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Иконка</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Название</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">English</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Описание</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Дата создания</th>
                    </tr>
                  </thead>
                  <tbody>
                    {species.map((item) => (
                      <tr key={item.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                        <td className="py-3 px-4 text-2xl">{item.icon}</td>
                        <td className="py-3 px-4 text-sm font-medium text-gray-900">{item.name}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{item.name_en}</td>
                        <td className="py-3 px-4 text-sm text-gray-600 max-w-md truncate">{item.description}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">
                          {new Date(item.created_at).toLocaleDateString('ru-RU')}
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

      {activeTab === 'breeds' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Породы</h2>
            <p className="text-gray-600">Управление породами животных</p>
          </div>

          <TableWidget
            title={`Породы (${breeds.length})`}
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Добавить породу
              </button>
            }
          >
            {loading ? (
              <div className="text-center py-12 text-gray-400">Загрузка...</div>
            ) : breeds.length === 0 ? (
              <div className="text-center py-12 text-gray-400">
                <RectangleStackIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Породы не найдены</p>
                <p className="text-sm mt-2">Добавьте первую породу, чтобы начать работу</p>
              </div>
            ) : (
              <div className="overflow-x-auto p-6">
                <table className="w-full min-w-[900px]">
                  <thead>
                    <tr className="border-b border-gray-100">
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Название</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Вид</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Размер</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Вес (кг)</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Продолжительность жизни</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Уход</th>
                    </tr>
                  </thead>
                  <tbody>
                    {breeds.map((breed) => (
                      <tr key={breed.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                        <td className="py-3 px-4">
                          <div className="font-medium text-gray-900">{breed.name}</div>
                          <div className="text-xs text-gray-500">{breed.name_en}</div>
                        </td>
                        <td className="py-3 px-4 text-sm text-gray-600">{breed.species_name}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{breed.size || '-'}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">
                          {breed.weight_min && breed.weight_max
                            ? `${breed.weight_min}-${breed.weight_max}`
                            : '-'}
                        </td>
                        <td className="py-3 px-4 text-sm text-gray-600">
                          {breed.lifespan_min && breed.lifespan_max
                            ? `${breed.lifespan_min}-${breed.lifespan_max} лет`
                            : '-'}
                        </td>
                        <td className="py-3 px-4 text-sm text-gray-600">{breed.care_level || '-'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </TableWidget>
        </div>
      )}
      {activeTab === 'pets' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Питомцы пользователей</h2>
            <p className="text-gray-600">Реальные карточки питомцев из PetID реестра</p>
          </div>

          <TableWidget
            title={`Питомцев (${pets.length})`}
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Добавить питомца
              </button>
            }
          >
            {loading ? (
              <div className="text-center py-12 text-gray-400">Загрузка...</div>
            ) : pets.length === 0 ? (
              <div className="text-center py-12 text-gray-400">
                <Square3Stack3DIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Питомцы не найдены</p>
                <p className="text-sm mt-2">Добавьте первого питомца, чтобы начать работу</p>
              </div>
            ) : (
              <div className="overflow-x-auto p-6">
                <table className="w-full min-w-[1000px]">
                  <thead>
                    <tr className="border-b border-gray-100">
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Фото</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Имя</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Вид</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Порода</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Пол</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Возраст</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Чип</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Статус</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Владелец</th>
                    </tr>
                  </thead>
                  <tbody>
                    {pets.map((pet) => {
                      // Вычисляем возраст
                      let age = '-';
                      if (pet.birth_date) {
                        const birthDate = new Date(pet.birth_date);
                        const today = new Date();
                        const years = today.getFullYear() - birthDate.getFullYear();
                        const months = today.getMonth() - birthDate.getMonth();
                        if (years > 0) {
                          age = `${years} ${years === 1 ? 'год' : years < 5 ? 'года' : 'лет'}`;
                        } else if (months > 0) {
                          age = `${months} ${months === 1 ? 'месяц' : months < 5 ? 'месяца' : 'месяцев'}`;
                        } else {
                          age = 'Новорождённый';
                        }
                      }

                      // Статус на русском
                      const statusMap: Record<string, string> = {
                        home: 'Дома',
                        looking_for_home: 'Ищет дом',
                        lost: 'Потерялся',
                        found: 'Найден',
                        deceased: 'Умер',
                      };

                      // Цвет статуса
                      const statusColorMap: Record<string, string> = {
                        home: 'bg-green-100 text-green-700',
                        looking_for_home: 'bg-orange-100 text-orange-700',
                        lost: 'bg-red-100 text-red-700',
                        found: 'bg-blue-100 text-blue-700',
                        deceased: 'bg-gray-100 text-gray-700',
                      };

                      return (
                        <tr 
                          key={pet.id} 
                          onClick={() => router.push(`/pets/${pet.id}`)}
                          className="border-b border-gray-50 hover:bg-blue-50 transition-colors cursor-pointer"
                        >
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
                          <td className="py-3 px-4">
                            <div className="font-medium text-gray-900">{pet.name}</div>
                            <div className="text-xs text-gray-500">ID: {pet.id}</div>
                          </td>
                          <td className="py-3 px-4 text-sm text-gray-600">
                            {pet.species === 'dog' && '🐕 Собака'}
                            {pet.species === 'cat' && '🐈 Кошка'}
                            {pet.species === 'bird' && '🐦 Птица'}
                            {pet.species === 'other' && '🐾 Другое'}
                            {!pet.species && '-'}
                          </td>
                          <td className="py-3 px-4 text-sm text-gray-600">{pet.breed || '-'}</td>
                          <td className="py-3 px-4 text-sm text-gray-600">
                            {pet.gender === 'male' && '♂️ Самец'}
                            {pet.gender === 'female' && '♀️ Самка'}
                            {!pet.gender && '-'}
                          </td>
                          <td className="py-3 px-4 text-sm text-gray-600">{age}</td>
                          <td className="py-3 px-4 text-sm text-gray-600">
                            {pet.chip_number ? (
                              <span className="font-mono text-xs bg-gray-100 px-2 py-1 rounded">
                                {pet.chip_number}
                              </span>
                            ) : (
                              '-'
                            )}
                          </td>
                          <td className="py-3 px-4">
                            <span
                              className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${
                                statusColorMap[pet.status] || 'bg-gray-100 text-gray-600'
                              }`}
                            >
                              {statusMap[pet.status] || pet.status}
                            </span>
                          </td>
                          <td className="py-3 px-4 text-sm text-gray-600">
                            User #{pet.user_id}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </TableWidget>
        </div>
      )}
    </AdminLayout>
  );
}
