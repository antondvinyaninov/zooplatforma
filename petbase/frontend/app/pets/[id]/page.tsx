'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import AdminLayout, { AdminTab } from '@/app/components/admin/AdminLayout';
import {
  ArrowLeftIcon,
  CalendarIcon,
  ScaleIcon,
  HeartIcon,
  ShieldCheckIcon,
  DocumentTextIcon,
  Square3Stack3DIcon,
  RectangleStackIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

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
  sterilization_date?: string;
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
  // Новые поля паспорта
  distinctive_marks?: string;
  tattoo_number?: string;
  owner_name?: string;
  owner_address?: string;
  owner_phone?: string;
  owner_email?: string;
  blood_type?: string;
  allergies?: string;
  chronic_diseases?: string;
  current_medications?: string;
  pedigree_number?: string;
  registration_org?: string;
  ear_tag_number?: string;
  // Куратор и локация
  curator_id?: number;
  curator_name?: string;
  curator_phone?: string;
  location?: string;
  foster_address?: string;
  shelter_name?: string;
}

export default function PetDetailPage() {
  const router = useRouter();
  const params = useParams();
  const [pet, setPet] = useState<Pet | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeSection, setActiveSection] = useState<'info' | 'owner' | 'identification' | 'pedigree' | 'medical' | 'history'>('info');

  const tabs: AdminTab[] = [
    { id: 'stats', label: 'Статистика', icon: <ChartBarIcon className="w-5 h-5" /> },
    { id: 'species', label: 'Виды', icon: <Square3Stack3DIcon className="w-5 h-5" /> },
    { id: 'breeds', label: 'Породы', icon: <RectangleStackIcon className="w-5 h-5" /> },
    { id: 'pets', label: 'Питомцы', icon: <Square3Stack3DIcon className="w-5 h-5" /> },
  ];

  useEffect(() => {
    const loadPet = async () => {
      try {
        // Временно для разработки: используем X-User-ID заголовок
        // TODO: Заменить на JWT токен после внедрения полной аутентификации
        const response = await fetch(`http://localhost:8100/api/pets/${params.id}`, {
          headers: {
            'X-User-ID': '1', // Временно: ID администратора
          },
        });
        const result = await response.json();
        
        if (result.success) {
          setPet(result.data);
        } else {
          console.error('Failed to load pet:', result.error);
        }
      } catch (error) {
        console.error('Error loading pet:', error);
      } finally {
        setLoading(false);
      }
    };

    if (params.id) {
      loadPet();
    }
  }, [params.id]);

  if (loading) {
    return (
      <AdminLayout
        logoSrc="/favicon.svg"
        logoText="ЗооБаза"
        logoAlt="ЗооБаза"
        tabs={tabs}
        activeTab="pets"
        onTabChange={(tabId) => {
          if (tabId !== 'pets') {
            router.push('/');
          }
        }}
        adminUser={{ email: 'admin', role: 'admin' }}
        onLogout={() => {
          router.push('/auth');
        }}
        mainSiteUrl="http://localhost:3000"
      >
        <div className="flex items-center justify-center py-12">
          <div className="text-gray-600">Загрузка...</div>
        </div>
      </AdminLayout>
    );
  }

  if (!pet) {
    return (
      <AdminLayout
        logoSrc="/favicon.svg"
        logoText="ЗооБаза"
        logoAlt="ЗооБаза"
        tabs={tabs}
        activeTab="pets"
        onTabChange={(tabId) => {
          if (tabId !== 'pets') {
            router.push('/');
          }
        }}
        adminUser={{ email: 'admin', role: 'admin' }}
        onLogout={() => {
          router.push('/auth');
        }}
        mainSiteUrl="http://localhost:3000"
      >
        <div className="flex items-center justify-center py-12">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Питомец не найден</h2>
            <button
              onClick={() => router.push('/')}
              className="text-blue-600 hover:text-blue-700"
            >
              Вернуться назад
            </button>
          </div>
        </div>
      </AdminLayout>
    );
  }

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
  const statusMap: Record<string, { label: string; color: string }> = {
    home: { label: 'Дома', color: 'bg-green-100 text-green-700' },
    looking_for_home: { label: 'Ищет дом', color: 'bg-yellow-100 text-yellow-700' },
    lost: { label: 'Потерялся', color: 'bg-red-100 text-red-700' },
    found: { label: 'Найден', color: 'bg-blue-100 text-blue-700' },
    adopted: { label: 'Пристроен', color: 'bg-purple-100 text-purple-700' },
  };

  const status = statusMap[pet.status] || { label: pet.status, color: 'bg-gray-100 text-gray-700' };

  // Пол на русском
  const genderMap: Record<string, string> = {
    male: 'Самец',
    female: 'Самка',
  };

  // Размер на русском
  const sizeMap: Record<string, string> = {
    small: 'Маленький',
    medium: 'Средний',
    large: 'Большой',
  };

  return (
    <AdminLayout
      logoSrc="/favicon.svg"
      logoText="ЗооБаза"
      logoAlt="ЗооБаза"
      tabs={tabs}
      activeTab="pets"
      onTabChange={(tabId) => {
        if (tabId !== 'pets') {
          router.push('/');
        }
      }}
      adminUser={{ email: 'admin', role: 'admin' }}
      onLogout={() => {
        router.push('/auth');
      }}
      mainSiteUrl="http://localhost:3000"
    >
      <div className="space-y-6">
        {/* Back Button */}
        <button
          onClick={() => router.push('/')}
          className="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
        >
          <ArrowLeftIcon className="w-5 h-5 mr-2" />
          Назад к списку
        </button>

        {/* Pet Card */}
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
          {/* Hero Section */}
          <div className="bg-gradient-to-r from-blue-500 to-purple-600 p-8 text-white">
            <div className="flex items-start justify-between">
              <div>
                <div className="flex items-center gap-3 mb-2">
                  <h1 className="text-4xl font-bold">{pet.name}</h1>
                  <span className="text-sm text-blue-100 bg-blue-600/30 px-3 py-1 rounded-full">
                    ID: #{pet.id}
                  </span>
                </div>
                <p className="text-xl text-blue-100 mb-2">
                  {pet.species} {pet.breed && `• ${pet.breed}`}
                </p>
                <p className="text-sm text-blue-200">
                  Создано: {new Date(pet.created_at).toLocaleDateString('ru-RU', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </p>
              </div>
              <span className={`px-4 py-2 rounded-full text-sm font-medium ${status.color}`}>
                {status.label}
              </span>
            </div>
          </div>

          {/* Section Tabs */}
          <div className="border-b border-gray-200 bg-gray-50">
            <div className="flex space-x-1 p-2 overflow-x-auto">
              <button
                onClick={() => setActiveSection('info')}
                className={`px-6 py-3 rounded-lg font-medium transition-colors whitespace-nowrap ${
                  activeSection === 'info'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
              >
                📋 Основная информация
              </button>
              <button
                onClick={() => setActiveSection('owner')}
                className={`px-6 py-3 rounded-lg font-medium transition-colors whitespace-nowrap ${
                  activeSection === 'owner'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
              >
                👤 Владелец и куратор
              </button>
              <button
                onClick={() => setActiveSection('identification')}
                className={`px-6 py-3 rounded-lg font-medium transition-colors whitespace-nowrap ${
                  activeSection === 'identification'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
              >
                🔍 Идентификация
              </button>
              <button
                onClick={() => setActiveSection('pedigree')}
                className={`px-6 py-3 rounded-lg font-medium transition-colors whitespace-nowrap ${
                  activeSection === 'pedigree'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
              >
                👑 Родословная
              </button>
              <button
                onClick={() => setActiveSection('medical')}
                className={`px-6 py-3 rounded-lg font-medium transition-colors whitespace-nowrap ${
                  activeSection === 'medical'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
              >
                🏥 Медицинская информация
              </button>
              <button
                onClick={() => setActiveSection('history')}
                className={`px-6 py-3 rounded-lg font-medium transition-colors whitespace-nowrap ${
                  activeSection === 'history'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
              >
                📜 История
              </button>
            </div>
          </div>

          {/* Content Sections */}
          <div className="p-8">
            {/* Основная информация */}
            {activeSection === 'info' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Основная информация</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {pet.gender && (
                    <div>
                      <div className="text-sm text-gray-500">Пол</div>
                      <div className="text-base font-medium text-gray-900">
                        {genderMap[pet.gender] || pet.gender}
                      </div>
                    </div>
                  )}

                  {pet.birth_date && (
                    <div>
                      <div className="text-sm text-gray-500 flex items-center">
                        <CalendarIcon className="w-4 h-4 mr-1" />
                        Возраст
                      </div>
                      <div className="text-base font-medium text-gray-900">{age}</div>
                      <div className="text-xs text-gray-400">
                        Дата рождения: {new Date(pet.birth_date).toLocaleDateString('ru-RU')}
                      </div>
                    </div>
                  )}

                  {pet.color && (
                    <div>
                      <div className="text-sm text-gray-500">Окрас</div>
                      <div className="text-base font-medium text-gray-900">{pet.color}</div>
                    </div>
                  )}

                  {pet.size && (
                    <div>
                      <div className="text-sm text-gray-500">Размер</div>
                      <div className="text-base font-medium text-gray-900">
                        {sizeMap[pet.size] || pet.size}
                      </div>
                    </div>
                  )}

                  {pet.weight && (
                    <div>
                      <div className="text-sm text-gray-500 flex items-center">
                        <ScaleIcon className="w-4 h-4 mr-1" />
                        Вес
                      </div>
                      <div className="text-base font-medium text-gray-900">{pet.weight} кг</div>
                    </div>
                  )}

                  {pet.distinctive_marks && (
                    <div className="md:col-span-2 lg:col-span-3">
                      <div className="text-sm text-gray-500">Особые приметы</div>
                      <div className="text-sm text-gray-700 mt-1">{pet.distinctive_marks}</div>
                    </div>
                  )}

                  {pet.character_traits && (
                    <div className="md:col-span-2 lg:col-span-3">
                      <div className="text-sm text-gray-500">Характер</div>
                      <div className="text-sm text-gray-700 mt-1">{pet.character_traits}</div>
                    </div>
                  )}

                  {pet.special_needs && (
                    <div className="md:col-span-2 lg:col-span-3">
                      <div className="text-sm text-gray-500">Особые потребности</div>
                      <div className="text-sm text-gray-700 mt-1">{pet.special_needs}</div>
                    </div>
                  )}

                  {pet.story && (
                    <div className="md:col-span-2 lg:col-span-3">
                      <div className="text-sm text-gray-500">История</div>
                      <div className="text-sm text-gray-700 mt-1">{pet.story}</div>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Владелец и куратор */}
            {activeSection === 'owner' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Владелец и куратор</h2>
                
                {/* Кто отвечает за питомца */}
                <div className="bg-gray-50 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">👥</span>
                    Кто отвечает за питомца
                  </h3>
                  
                  {/* Есть владелец */}
                  {pet.user_id && pet.user_id > 0 && (
                    <div className="bg-green-50 border-2 border-green-200 rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-3">
                        <span className="text-2xl">👤</span>
                        <span className="font-semibold text-green-900 text-lg">Есть владелец</span>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <div className="text-xs text-green-600">ID владельца в системе</div>
                          <div className="text-sm font-medium text-green-900">#{pet.user_id}</div>
                        </div>

                        {pet.owner_name && (
                          <div>
                            <div className="text-xs text-green-600">ФИО владельца</div>
                            <div className="text-sm font-medium text-green-900">{pet.owner_name}</div>
                          </div>
                        )}

                        {pet.owner_phone && (
                          <div>
                            <div className="text-xs text-green-600">Телефон</div>
                            <div className="text-sm font-medium text-green-900">
                              <a href={`tel:${pet.owner_phone}`} className="text-green-600 hover:text-green-700 underline">
                                {pet.owner_phone}
                              </a>
                            </div>
                          </div>
                        )}

                        {pet.owner_email && (
                          <div>
                            <div className="text-xs text-green-600">Email</div>
                            <div className="text-sm font-medium text-green-900">
                              <a href={`mailto:${pet.owner_email}`} className="text-green-600 hover:text-green-700 underline">
                                {pet.owner_email}
                              </a>
                            </div>
                          </div>
                        )}

                        {pet.owner_address && (
                          <div className="md:col-span-2">
                            <div className="text-xs text-green-600">Адрес</div>
                            <div className="text-sm text-green-900">{pet.owner_address}</div>
                          </div>
                        )}
                      </div>
                    </div>
                  )}

                  {/* Есть куратор (но нет владельца) */}
                  {(!pet.user_id || pet.user_id === 0) && (pet.curator_id || pet.curator_name || pet.curator_phone) && (
                    <div className="bg-blue-50 border-2 border-blue-200 rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-3">
                        <span className="text-2xl">🤝</span>
                        <span className="font-semibold text-blue-900 text-lg">Есть куратор</span>
                      </div>
                      <div className="text-sm text-blue-800 mb-3">
                        За питомцем следит волонтёр-куратор
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {pet.curator_id && (
                          <div>
                            <div className="text-xs text-blue-600">ID куратора</div>
                            <div className="text-sm font-medium text-blue-900">#{pet.curator_id}</div>
                          </div>
                        )}

                        {pet.curator_name && (
                          <div>
                            <div className="text-xs text-blue-600">ФИО куратора</div>
                            <div className="text-sm font-medium text-blue-900">{pet.curator_name}</div>
                          </div>
                        )}

                        {pet.curator_phone && (
                          <div className="md:col-span-2">
                            <div className="text-xs text-blue-600">Телефон куратора</div>
                            <div className="text-sm font-medium text-blue-900">
                              <a href={`tel:${pet.curator_phone}`} className="text-blue-600 hover:text-blue-700 underline">
                                {pet.curator_phone}
                              </a>
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  )}

                  {/* Нет ни владельца, ни куратора */}
                  {(!pet.user_id || pet.user_id === 0) && !pet.curator_id && !pet.curator_name && !pet.curator_phone && (
                    <div className="bg-gray-100 border-2 border-gray-300 rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-2xl">🐾</span>
                        <span className="font-semibold text-gray-900 text-lg">Бездомное животное</span>
                      </div>
                      <div className="text-sm text-gray-700">
                        У питомца нет ни владельца, ни куратора. 
                        {pet.location === 'street' && ' Живёт на улице.'}
                        {pet.location === 'shelter' && ' Находится в приюте.'}
                        {pet.location === 'foster' && ' Находится на передержке.'}
                      </div>
                    </div>
                  )}
                </div>

                {/* Текущая локация */}
                <div className="bg-blue-50 rounded-lg p-6 border-2 border-blue-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">📍</span>
                    Текущая локация
                  </h3>
                  <div className="space-y-4">
                    <div>
                      <div className="text-sm text-gray-500 mb-2">Где находится питомец</div>
                      <div className="flex items-center gap-2">
                        {pet.location === 'home' && (
                          <span className="px-4 py-2 bg-green-100 text-green-800 rounded-lg font-medium">
                            🏠 Дома у владельца
                          </span>
                        )}
                        {pet.location === 'street' && (
                          <span className="px-4 py-2 bg-orange-100 text-orange-800 rounded-lg font-medium">
                            🌆 На улице
                          </span>
                        )}
                        {pet.location === 'foster' && (
                          <span className="px-4 py-2 bg-purple-100 text-purple-800 rounded-lg font-medium">
                            🏡 На передержке
                          </span>
                        )}
                        {pet.location === 'shelter' && (
                          <span className="px-4 py-2 bg-blue-100 text-blue-800 rounded-lg font-medium">
                            🏥 В приюте
                          </span>
                        )}
                      </div>
                    </div>

                    {pet.location === 'foster' && pet.foster_address && (
                      <div>
                        <div className="text-sm text-gray-500">Адрес передержки</div>
                        <div className="text-base text-gray-900 mt-1">{pet.foster_address}</div>
                      </div>
                    )}

                    {pet.location === 'shelter' && pet.shelter_name && (
                      <div>
                        <div className="text-sm text-gray-500">Название приюта</div>
                        <div className="text-base text-gray-900 mt-1">{pet.shelter_name}</div>
                      </div>
                    )}
                  </div>
                </div>

                {/* Пояснение */}
                <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                  <div className="text-sm text-yellow-800">
                    <strong>ℹ️ Система ответственности:</strong>
                    <ul className="mt-2 space-y-1 ml-4 list-disc">
                      <li><strong>Владелец</strong> - постоянный хозяин питомца. Животное обычно живёт дома.</li>
                      <li><strong>Куратор</strong> - волонтёр, который следит за бездомным животным (на улице, передержке или в приюте).</li>
                      <li><strong>Без владельца и куратора</strong> - бездомное животное, за которым никто не следит.</li>
                      <li>У животного может быть <strong>либо</strong> владелец, <strong>либо</strong> куратор, <strong>либо</strong> никого.</li>
                    </ul>
                  </div>
                </div>
              </div>
            )}

            {/* Идентификация */}
            {activeSection === 'identification' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Идентификация питомца</h2>
                
                {/* Микрочип */}
                <div className="bg-blue-50 rounded-lg p-6 border-2 border-blue-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">💎</span>
                    Микрочип
                  </h3>
                  {pet.chip_number ? (
                    <div>
                      <div className="text-sm text-gray-500 mb-1">Номер микрочипа</div>
                      <div className="text-2xl font-mono font-bold text-blue-900 mb-2">{pet.chip_number}</div>
                      <div className="text-xs text-green-600 flex items-center gap-1">
                        <span className="text-lg">✓</span>
                        <span>Чипирован</span>
                      </div>
                    </div>
                  ) : (
                    <div className="text-sm text-gray-500 italic">
                      Микрочип не установлен
                    </div>
                  )}
                </div>

                {/* Клеймо */}
                <div className="bg-purple-50 rounded-lg p-6 border-2 border-purple-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">🏷️</span>
                    Клеймо (татуировка)
                  </h3>
                  {pet.tattoo_number ? (
                    <div>
                      <div className="text-sm text-gray-500 mb-1">Номер клейма</div>
                      <div className="text-2xl font-mono font-bold text-purple-900 mb-2">{pet.tattoo_number}</div>
                      <div className="text-xs text-green-600 flex items-center gap-1">
                        <span className="text-lg">✓</span>
                        <span>Клеймо поставлено</span>
                      </div>
                    </div>
                  ) : (
                    <div className="text-sm text-gray-500 italic">
                      Клеймо не поставлено
                    </div>
                  )}
                </div>

                {/* Ушная бирка ОСВВ */}
                <div className="bg-orange-50 rounded-lg p-6 border-2 border-orange-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">🏷️</span>
                    Ушная бирка (ОСВВ)
                  </h3>
                  {pet.ear_tag_number ? (
                    <div>
                      <div className="text-sm text-gray-500 mb-1">Номер ушной бирки</div>
                      <div className="text-2xl font-mono font-bold text-orange-900 mb-2">{pet.ear_tag_number}</div>
                      <div className="text-xs text-green-600 flex items-center gap-1">
                        <span className="text-lg">✓</span>
                        <span>Бирка установлена в приюте</span>
                      </div>
                      <div className="text-xs text-orange-700 mt-2">
                        Программа ОСВВ (Отлов-Стерилизация-Вакцинация-Возврат)
                      </div>
                    </div>
                  ) : (
                    <div className="text-sm text-gray-500 italic">
                      Ушная бирка не установлена
                    </div>
                  )}
                </div>

                {/* Ветеринарный паспорт */}
                <div className="bg-green-50 rounded-lg p-6 border-2 border-green-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">📗</span>
                    Ветеринарный паспорт
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <div className="text-sm text-gray-500 mb-1">Номер паспорта</div>
                      {pet.passport_number ? (
                        <div className="text-lg font-mono font-semibold text-green-900">{pet.passport_number}</div>
                      ) : (
                        <div className="text-sm text-gray-400 italic">Не указан</div>
                      )}
                    </div>

                    {pet.blood_type && (
                      <div>
                        <div className="text-sm text-gray-500 mb-1">Группа крови</div>
                        <div className="text-lg font-semibold text-green-900">{pet.blood_type}</div>
                      </div>
                    )}
                  </div>
                </div>

                {/* Сводка по идентификации */}
                <div className="bg-gray-50 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Сводка по идентификации</h3>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="text-center">
                      <div className={`text-4xl mb-2 ${pet.chip_number ? 'text-green-500' : 'text-gray-300'}`}>
                        {pet.chip_number ? '✓' : '○'}
                      </div>
                      <div className="text-xs text-gray-600">Микрочип</div>
                    </div>
                    <div className="text-center">
                      <div className={`text-4xl mb-2 ${pet.tattoo_number ? 'text-green-500' : 'text-gray-300'}`}>
                        {pet.tattoo_number ? '✓' : '○'}
                      </div>
                      <div className="text-xs text-gray-600">Клеймо</div>
                    </div>
                    <div className="text-center">
                      <div className={`text-4xl mb-2 ${pet.ear_tag_number ? 'text-green-500' : 'text-gray-300'}`}>
                        {pet.ear_tag_number ? '✓' : '○'}
                      </div>
                      <div className="text-xs text-gray-600">Бирка</div>
                    </div>
                    <div className="text-center">
                      <div className={`text-4xl mb-2 ${pet.passport_number ? 'text-green-500' : 'text-gray-300'}`}>
                        {pet.passport_number ? '✓' : '○'}
                      </div>
                      <div className="text-xs text-gray-600">Паспорт</div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Родословная */}
            {activeSection === 'pedigree' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Родословная и регистрация</h2>
                
                {/* Родословная */}
                <div className="bg-purple-50 rounded-lg p-6 border-2 border-purple-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">👑</span>
                    Родословная
                  </h3>
                  {pet.pedigree_number ? (
                    <div>
                      <div className="text-sm text-gray-500 mb-1">Номер родословной</div>
                      <div className="text-2xl font-mono font-bold text-purple-900 mb-2">{pet.pedigree_number}</div>
                      <div className="text-xs text-green-600 flex items-center gap-1">
                        <span className="text-lg">✓</span>
                        <span>Родословная оформлена</span>
                      </div>
                    </div>
                  ) : (
                    <div className="text-sm text-gray-500 italic">
                      Родословная не оформлена
                    </div>
                  )}
                </div>

                {/* Организация регистрации */}
                <div className="bg-blue-50 rounded-lg p-6 border-2 border-blue-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">🏛️</span>
                    Организация регистрации
                  </h3>
                  {pet.registration_org ? (
                    <div>
                      <div className="text-sm text-gray-500 mb-1">Зарегистрирован в</div>
                      <div className="text-xl font-semibold text-blue-900 mb-2">{pet.registration_org}</div>
                      <div className="text-xs text-blue-700 mt-2">
                        {pet.registration_org.includes('РКФ') && '🇷🇺 Российская Кинологическая Федерация'}
                        {pet.registration_org.includes('FCI') && '🌍 Международная Кинологическая Федерация'}
                        {pet.registration_org.includes('WCF') && '🌍 World Cat Federation'}
                        {pet.registration_org.includes('TICA') && '🌍 The International Cat Association'}
                      </div>
                    </div>
                  ) : (
                    <div className="text-sm text-gray-500 italic">
                      Не зарегистрирован в кинологических/фелинологических организациях
                    </div>
                  )}
                </div>

                {/* Информация о породе */}
                {pet.breed && (
                  <div className="bg-gray-50 rounded-lg p-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                      <span className="text-2xl mr-2">🐾</span>
                      Информация о породе
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <div className="text-sm text-gray-500">Вид</div>
                        <div className="text-base font-medium text-gray-900">{pet.species}</div>
                      </div>
                      <div>
                        <div className="text-sm text-gray-500">Порода</div>
                        <div className="text-base font-medium text-gray-900">{pet.breed}</div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Пояснение */}
                <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                  <div className="text-sm text-yellow-800">
                    <strong>ℹ️ О родословной:</strong>
                    <ul className="mt-2 space-y-1 ml-4 list-disc">
                      <li><strong>Родословная</strong> - официальный документ о происхождении животного</li>
                      <li><strong>РКФ</strong> - Российская Кинологическая Федерация (для собак)</li>
                      <li><strong>FCI</strong> - Международная Кинологическая Федерация (для собак)</li>
                      <li><strong>WCF/TICA</strong> - Международные фелинологические организации (для кошек)</li>
                      <li>Родословная нужна для участия в выставках и племенного разведения</li>
                    </ul>
                  </div>
                </div>
              </div>
            )}

            {/* Медицинская информация */}
            {activeSection === 'medical' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Медицинская информация</h2>
                
                {/* Базовые данные */}
                <div className="bg-gray-50 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">📊</span>
                    Базовые данные
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    {pet.gender && (
                      <div>
                        <div className="text-sm text-gray-500">Пол</div>
                        <div className="text-base font-medium text-gray-900">
                          {genderMap[pet.gender] || pet.gender}
                        </div>
                      </div>
                    )}

                    {pet.birth_date && (
                      <div>
                        <div className="text-sm text-gray-500 flex items-center">
                          <CalendarIcon className="w-4 h-4 mr-1" />
                          Возраст
                        </div>
                        <div className="text-base font-medium text-gray-900">{age}</div>
                        <div className="text-xs text-gray-400">
                          Дата рождения: {new Date(pet.birth_date).toLocaleDateString('ru-RU')}
                        </div>
                      </div>
                    )}

                    {pet.weight && (
                      <div>
                        <div className="text-sm text-gray-500 flex items-center">
                          <ScaleIcon className="w-4 h-4 mr-1" />
                          Вес
                        </div>
                        <div className="text-base font-medium text-gray-900">{pet.weight} кг</div>
                      </div>
                    )}

                    {pet.blood_type && (
                      <div>
                        <div className="text-sm text-gray-500">Группа крови</div>
                        <div className="text-base font-medium text-gray-900">{pet.blood_type}</div>
                      </div>
                    )}
                  </div>
                </div>

                {/* Стерилизация и вакцинация */}
                <div className="bg-blue-50 rounded-lg p-6 border-2 border-blue-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">💉</span>
                    Стерилизация и вакцинация
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className={`p-4 rounded-lg ${pet.is_sterilized ? 'bg-green-100 border-2 border-green-300' : 'bg-gray-100 border-2 border-gray-300'}`}>
                      <div className="flex items-center gap-2 mb-2">
                        <div className={`w-4 h-4 rounded-full ${pet.is_sterilized ? 'bg-green-500' : 'bg-gray-400'}`} />
                        <span className="font-semibold text-gray-900">
                          {pet.is_sterilized ? 'Стерилизован/кастрирован' : 'Не стерилизован'}
                        </span>
                      </div>
                      {pet.is_sterilized && pet.sterilization_date && (
                        <div className="text-sm text-gray-700 mt-2">
                          Дата операции: {new Date(pet.sterilization_date).toLocaleDateString('ru-RU', {
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                          })}
                        </div>
                      )}
                    </div>

                    <div className={`p-4 rounded-lg ${pet.is_vaccinated ? 'bg-green-100 border-2 border-green-300' : 'bg-gray-100 border-2 border-gray-300'}`}>
                      <div className="flex items-center gap-2 mb-2">
                        <div className={`w-4 h-4 rounded-full ${pet.is_vaccinated ? 'bg-green-500' : 'bg-gray-400'}`} />
                        <span className="font-semibold text-gray-900">
                          {pet.is_vaccinated ? 'Вакцинирован' : 'Не вакцинирован'}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Медицинские особенности */}
                {(pet.allergies || pet.chronic_diseases || pet.current_medications || pet.health_notes) && (
                  <div className="space-y-4">
                    <h3 className="text-lg font-semibold text-gray-900 flex items-center">
                      <span className="text-2xl mr-2">⚠️</span>
                      Медицинские особенности
                    </h3>
                    
                    {pet.allergies && (
                      <div className="bg-red-50 border-2 border-red-200 rounded-lg p-4">
                        <div className="text-sm font-semibold text-red-900 mb-1 flex items-center gap-2">
                          <span className="text-xl">🚨</span>
                          Аллергии
                        </div>
                        <p className="text-sm text-red-800">{pet.allergies}</p>
                      </div>
                    )}

                    {pet.chronic_diseases && (
                      <div className="bg-orange-50 border-2 border-orange-200 rounded-lg p-4">
                        <div className="text-sm font-semibold text-orange-900 mb-1 flex items-center gap-2">
                          <span className="text-xl">🏥</span>
                          Хронические заболевания
                        </div>
                        <p className="text-sm text-orange-800">{pet.chronic_diseases}</p>
                      </div>
                    )}

                    {pet.current_medications && (
                      <div className="bg-blue-50 border-2 border-blue-200 rounded-lg p-4">
                        <div className="text-sm font-semibold text-blue-900 mb-1 flex items-center gap-2">
                          <span className="text-xl">💊</span>
                          Текущие лекарства
                        </div>
                        <p className="text-sm text-blue-800">{pet.current_medications}</p>
                      </div>
                    )}

                    {pet.health_notes && (
                      <div className="bg-yellow-50 border-2 border-yellow-200 rounded-lg p-4">
                        <div className="text-sm font-semibold text-yellow-900 mb-1 flex items-center gap-2">
                          <span className="text-xl">📝</span>
                          Заметки о здоровье
                        </div>
                        <p className="text-sm text-yellow-800">{pet.health_notes}</p>
                      </div>
                    )}
                  </div>
                )}

                {/* Медицинские записи */}
                <div className="border-t pt-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">📋</span>
                    Медицинские записи
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <div className="bg-green-50 border-2 border-green-200 rounded-lg p-4 text-center hover:bg-green-100 transition-colors cursor-pointer">
                      <div className="text-3xl mb-2">💉</div>
                      <div className="text-sm font-semibold text-gray-900">Прививки</div>
                      <div className="text-xs text-gray-500 mt-1">Скоро будет доступно</div>
                    </div>
                    <div className="bg-purple-50 border-2 border-purple-200 rounded-lg p-4 text-center hover:bg-purple-100 transition-colors cursor-pointer">
                      <div className="text-3xl mb-2">🐛</div>
                      <div className="text-sm font-semibold text-gray-900">Обработки</div>
                      <div className="text-xs text-gray-500 mt-1">От паразитов</div>
                    </div>
                    <div className="bg-pink-50 border-2 border-pink-200 rounded-lg p-4 text-center hover:bg-pink-100 transition-colors cursor-pointer">
                      <div className="text-3xl mb-2">🏥</div>
                      <div className="text-sm font-semibold text-gray-900">Операции</div>
                      <div className="text-xs text-gray-500 mt-1">История операций</div>
                    </div>
                    <div className="bg-blue-50 border-2 border-blue-200 rounded-lg p-4 text-center hover:bg-blue-100 transition-colors cursor-pointer">
                      <div className="text-3xl mb-2">🩺</div>
                      <div className="text-sm font-semibold text-gray-900">Визиты к врачу</div>
                      <div className="text-xs text-gray-500 mt-1">История обращений</div>
                    </div>
                  </div>
                </div>

                {/* Анализы и диагнозы */}
                <div className="bg-gray-50 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">🔬</span>
                    Анализы и диагнозы
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="bg-white border-2 border-gray-200 rounded-lg p-4 text-center hover:bg-gray-50 transition-colors cursor-pointer">
                      <div className="text-3xl mb-2">🧪</div>
                      <div className="text-sm font-semibold text-gray-900">Анализы</div>
                      <div className="text-xs text-gray-500 mt-1">Результаты анализов</div>
                    </div>
                    <div className="bg-white border-2 border-gray-200 rounded-lg p-4 text-center hover:bg-gray-50 transition-colors cursor-pointer">
                      <div className="text-3xl mb-2">📊</div>
                      <div className="text-sm font-semibold text-gray-900">Диагнозы</div>
                      <div className="text-xs text-gray-500 mt-1">История диагнозов</div>
                    </div>
                  </div>
                </div>

                {/* План лечения */}
                <div className="bg-indigo-50 border-2 border-indigo-200 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <span className="text-2xl mr-2">📅</span>
                    План лечения
                  </h3>
                  <div className="text-sm text-gray-600 text-center py-4">
                    <div className="text-4xl mb-2">📋</div>
                    <div>Функция планирования лечения будет доступна в следующих версиях</div>
                  </div>
                </div>
              </div>
            )}

            {/* История */}
            {activeSection === 'history' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">История действий</h2>
                <div className="space-y-4">
                  {/* Timeline */}
                  <div className="relative border-l-2 border-gray-200 pl-6 space-y-6">
                    <div className="relative">
                      <div className="absolute -left-8 w-4 h-4 bg-blue-500 rounded-full border-4 border-white"></div>
                      <div className="bg-gray-50 rounded-lg p-4">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-semibold text-gray-900">Создание карточки</span>
                          <span className="text-xs text-gray-500">
                            {new Date(pet.created_at).toLocaleDateString('ru-RU', {
                              year: 'numeric',
                              month: 'long',
                              day: 'numeric',
                              hour: '2-digit',
                              minute: '2-digit',
                            })}
                          </span>
                        </div>
                        <p className="text-sm text-gray-600">
                          Карточка питомца была создана в системе
                        </p>
                      </div>
                    </div>

                    {pet.updated_at !== pet.created_at && (
                      <div className="relative">
                        <div className="absolute -left-8 w-4 h-4 bg-green-500 rounded-full border-4 border-white"></div>
                        <div className="bg-gray-50 rounded-lg p-4">
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm font-semibold text-gray-900">Обновление информации</span>
                            <span className="text-xs text-gray-500">
                              {new Date(pet.updated_at).toLocaleDateString('ru-RU', {
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit',
                              })}
                            </span>
                          </div>
                          <p className="text-sm text-gray-600">
                            Информация о питомце была обновлена
                          </p>
                        </div>
                      </div>
                    )}

                    <div className="relative">
                      <div className="absolute -left-8 w-4 h-4 bg-gray-300 rounded-full border-4 border-white"></div>
                      <div className="bg-gray-50 rounded-lg p-4">
                        <p className="text-sm text-gray-500 italic">
                          Полная история действий будет доступна в следующих версиях
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </AdminLayout>
  );
}
