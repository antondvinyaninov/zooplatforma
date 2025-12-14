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
  { id: 'cards', label: 'Карточки', icon: <DocumentTextIcon className="w-5 h-5" /> },
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

export default function PetBaseDashboard() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('stats');
  const [isClient, setIsClient] = useState(false);
  const [species, setSpecies] = useState<Species[]>([]);
  const [breeds, setBreeds] = useState<Breed[]>([]);
  const [cards, setCards] = useState<PetCard[]>([]);
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
    } else if (activeTab === 'cards') {
      const result = await petbaseApi.getCards();
      if (result.success) {
        setCards(result.data || []);
      }
    } else if (activeTab === 'stats') {
      // Загружаем все для статистики
      const [speciesResult, breedsResult, cardsResult] = await Promise.all([
        petbaseApi.getSpecies(),
        petbaseApi.getBreeds(),
        petbaseApi.getCards(),
      ]);
      if (speciesResult.success) setSpecies(speciesResult.data || []);
      if (breedsResult.success) setBreeds(breedsResult.data || []);
      if (cardsResult.success) setCards(cardsResult.data || []);
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
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Статистика ЗооБазы</h2>
            <p className="text-gray-600">Общая информация о справочнике</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatsWidget
              title="Видов животных"
              value={species.length}
              icon={<Square3Stack3DIcon className="w-6 h-6" />}
              color="blue"
            />
            <StatsWidget
              title="Пород"
              value={breeds.length}
              icon={<RectangleStackIcon className="w-6 h-6" />}
              color="green"
            />
            <StatsWidget
              title="Карточек"
              value={cards.length}
              icon={<DocumentTextIcon className="w-6 h-6" />}
              color="orange"
            />
            <StatsWidget
              title="Опубликовано"
              value={cards.filter((c) => c.is_published).length}
              icon={<ChartBarIcon className="w-6 h-6" />}
              color="purple"
            />
          </div>

          <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
            <h3 className="text-xl font-bold text-gray-900 mb-4">О ЗооБазе</h3>
            <div className="space-y-4 text-gray-600">
              <p>
                ЗооБаза - это централизованный справочник домашних животных, который содержит
                подробную информацию о видах, породах и особенностях содержания.
              </p>
              <p>
                Данные из ЗооБазы используются другими сервисами платформы для автозаполнения
                информации о питомцах пользователей и предоставления справочной информации.
              </p>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
                <div className="p-4 bg-blue-50 rounded-lg">
                  <div className="text-2xl font-bold text-blue-600 mb-1">API</div>
                  <div className="text-sm text-gray-600">Доступ через REST API</div>
                </div>
                <div className="p-4 bg-green-50 rounded-lg">
                  <div className="text-2xl font-bold text-green-600 mb-1">Поиск</div>
                  <div className="text-sm text-gray-600">Быстрый поиск по базе</div>
                </div>
                <div className="p-4 bg-purple-50 rounded-lg">
                  <div className="text-2xl font-bold text-purple-600 mb-1">Фото</div>
                  <div className="text-sm text-gray-600">Галерея изображений</div>
                </div>
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

      {activeTab === 'cards' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Карточки животных</h2>
            <p className="text-gray-600">Подробные карточки с информацией</p>
          </div>

          <TableWidget
            title={`Карточки (${cards.length})`}
            actions={
              <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                Добавить карточку
              </button>
            }
          >
            {loading ? (
              <div className="text-center py-12 text-gray-400">Загрузка...</div>
            ) : cards.length === 0 ? (
              <div className="text-center py-12 text-gray-400">
                <DocumentTextIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Карточки не найдены</p>
                <p className="text-sm mt-2">Добавьте первую карточку, чтобы начать работу</p>
              </div>
            ) : (
              <div className="overflow-x-auto p-6">
                <table className="w-full min-w-[900px]">
                  <thead>
                    <tr className="border-b border-gray-100">
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Заголовок</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Порода</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Описание</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Статус</th>
                      <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Дата создания</th>
                    </tr>
                  </thead>
                  <tbody>
                    {cards.map((card) => (
                      <tr key={card.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                        <td className="py-3 px-4 text-sm font-medium text-gray-900">{card.title}</td>
                        <td className="py-3 px-4 text-sm text-gray-600">{card.breed_name}</td>
                        <td className="py-3 px-4 text-sm text-gray-600 max-w-md truncate">
                          {card.description}
                        </td>
                        <td className="py-3 px-4">
                          <span
                            className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${
                              card.is_published
                                ? 'bg-green-100 text-green-700'
                                : 'bg-gray-100 text-gray-600'
                            }`}
                          >
                            {card.is_published ? 'Опубликовано' : 'Черновик'}
                          </span>
                        </td>
                        <td className="py-3 px-4 text-sm text-gray-600">
                          {new Date(card.created_at).toLocaleDateString('ru-RU')}
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
    </AdminLayout>
  );
}
