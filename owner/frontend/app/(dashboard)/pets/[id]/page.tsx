'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { 
  PencilIcon, 
  CalendarIcon, 
  HeartIcon,
  ScaleIcon,
  BeakerIcon,
  ShieldCheckIcon,
  SparklesIcon,
  ClockIcon,
  ArrowTopRightOnSquareIcon,
} from '@heroicons/react/24/outline';

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  gender?: string;
  birth_date?: string;
  color?: string;
  size?: string;
  weight?: number;
  chip_number?: string;
  tattoo_number?: string;
  ear_tag_number?: string;
  passport_number?: string;
  is_sterilized: boolean;
  sterilization_date?: string;
  is_vaccinated: boolean;
  health_notes?: string;
  character_traits?: string;
  special_needs?: string;
  photo?: string;
  photo_url?: string;
  story?: string;
  // Паспорт
  distinctive_marks?: string;
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
  // Куратор
  curator_name?: string;
  curator_phone?: string;
  location?: string;
  foster_address?: string;
  shelter_name?: string;
  // Экстренные контакты
  emergency_contact_name?: string;
  emergency_contact_phone?: string;
  emergency_contact_relation?: string;
  vet_clinic_name?: string;
  vet_clinic_phone?: string;
  vet_clinic_address?: string;
  insurance_company?: string;
  insurance_policy_number?: string;
  insurance_expiry_date?: string;
  // Каталог
  city?: string;
  region?: string;
  urgent?: boolean;
  contact_name?: string;
  contact_phone?: string;
  organization_name?: string;
}

export default function ViewPetPage() {
  const params = useParams();
  const router = useRouter();
  const [pet, setPet] = useState<Pet | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');

  useEffect(() => {
    loadPet();
  }, [params.id]);

  const loadPet = async () => {
    try {
      console.log('🔍 Loading pet with ID:', params.id);
      
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        credentials: 'include',
      });

      if (!meResponse.ok) {
        console.error('❌ Auth check failed:', meResponse.status);
        window.location.href = 'http://localhost:3000';
        return;
      }

      const meResult = await meResponse.json();
      console.log('✅ Auth result:', meResult);
      
      const userId = meResult.data?.user?.id || meResult.data?.id;

      if (!userId) {
        console.error('❌ User ID not found in response:', meResult);
        setLoading(false);
        return;
      }

      console.log('👤 User ID:', userId);

      // Получаем токен для запроса к Owner Backend
      const token = meResult.data?.token;
      
      if (!token) {
        console.error('❌ No token in auth response');
        setLoading(false);
        return;
      }

      console.log('🔑 Using token for Owner Backend request');

      const petResponse = await fetch(`http://localhost:8400/api/pets/${params.id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        credentials: 'include',
      });

      console.log('📡 Pet response status:', petResponse.status);

      if (!petResponse.ok) {
        console.error('❌ Failed to load pet:', petResponse.status);
        setLoading(false);
        return;
      }

      const petResult = await petResponse.json();
      console.log('📦 Pet result:', petResult);

      if (petResult.pet) {
        console.log('✅ Pet loaded:', petResult.pet);
        setPet(petResult.pet);
      } else {
        console.error('❌ No pet in response:', petResult);
      }
    } catch (error) {
      console.error('❌ Error loading pet:', error);
    } finally {
      setLoading(false);
    }
  };

  const getSpeciesEmoji = (species: string) => {
    const emojis: Record<string, string> = {
      dog: '🐕',
      cat: '🐈',
      bird: '🦜',
      rodent: '🐹',
      reptile: '🦎',
      fish: '🐠',
      other: '🐾',
    };
    return emojis[species] || '🐾';
  };

  const getSpeciesLabel = (species: string) => {
    const labels: Record<string, string> = {
      dog: 'Собака', cat: 'Кошка', bird: 'Птица',
      rodent: 'Грызун', reptile: 'Рептилия', fish: 'Рыба', other: 'Другое',
    };
    return labels[species] || species;
  };

  const getGenderLabel = (gender?: string) => {
    if (!gender) return 'Не указан';
    return gender === 'male' ? 'Самец' : 'Самка';
  };

  const getSizeLabel = (size?: string) => {
    if (!size) return 'Не указан';
    const labels: Record<string, string> = {
      small: 'Маленький', medium: 'Средний', large: 'Большой',
    };
    return labels[size] || size;
  };

  const calculateAge = (birthDate?: string) => {
    if (!birthDate) return null;
    const birth = new Date(birthDate);
    const now = new Date();
    const years = now.getFullYear() - birth.getFullYear();
    const months = now.getMonth() - birth.getMonth();
    
    if (years === 0) {
      return `${months} мес.`;
    } else if (months < 0) {
      return `${years - 1} лет ${12 + months} мес.`;
    } else {
      return `${years} лет ${months} мес.`;
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!pet) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600">Питомец не найден</p>
        <button
          onClick={() => router.push('/pets')}
          className="mt-4 text-blue-600 hover:text-blue-700"
        >
          Вернуться к списку
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Hero секция с фото */}
      <div className="relative bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 rounded-3xl overflow-hidden shadow-2xl">
        <div className="absolute inset-0 bg-black/10"></div>
        <div className="relative p-8">
          <div className="flex flex-col md:flex-row items-start md:items-center gap-6">
            {/* Фото питомца */}
            <div className="relative group">
              {pet.photo_url ? (
                <img
                  src={`http://localhost:8000${pet.photo_url}`}
                  alt={pet.name}
                  className="w-40 h-40 rounded-3xl object-cover border-4 border-white shadow-2xl group-hover:scale-105 transition-transform duration-300"
                />
              ) : (
                <div className="w-40 h-40 rounded-3xl bg-white/20 backdrop-blur-sm flex items-center justify-center text-7xl border-4 border-white shadow-2xl group-hover:scale-105 transition-transform duration-300">
                  {getSpeciesEmoji(pet.species)}
                </div>
              )}
              <div className="absolute -bottom-2 -right-2 bg-white rounded-full p-3 shadow-lg">
                <span className="text-3xl">{getSpeciesEmoji(pet.species)}</span>
              </div>
            </div>

            {/* Информация */}
            <div className="flex-1 text-white">
              <div className="flex items-center gap-3 mb-2">
                <h1 className="text-5xl font-bold">{pet.name}</h1>
                {pet.gender && (
                  <span className="text-3xl">{pet.gender === 'male' ? '♂️' : '♀️'}</span>
                )}
              </div>
              
              <div className="flex flex-wrap gap-2 mb-4">
                <span className="px-4 py-1.5 bg-white/20 backdrop-blur-sm rounded-full text-sm font-medium flex items-center gap-2">
                  <SparklesIcon className="w-4 h-4" />
                  {getSpeciesLabel(pet.species)}
                </span>
                {pet.breed && (
                  <span className="px-4 py-1.5 bg-white/20 backdrop-blur-sm rounded-full text-sm font-medium">
                    {pet.breed}
                  </span>
                )}
                {pet.birth_date && (
                  <span className="px-4 py-1.5 bg-white/20 backdrop-blur-sm rounded-full text-sm font-medium flex items-center gap-2">
                    <CalendarIcon className="w-4 h-4" />
                    {calculateAge(pet.birth_date)}
                  </span>
                )}
              </div>

              {/* Быстрые действия */}
              <div className="flex flex-wrap gap-3">
                <button
                  onClick={() => router.push(`/pets/${pet.id}/edit`)}
                  className="px-6 py-3 bg-white text-blue-600 rounded-xl hover:bg-gray-100 transition-all duration-200 font-semibold flex items-center gap-2 shadow-lg hover:shadow-xl hover:scale-105"
                >
                  <PencilIcon className="w-5 h-5" />
                  Редактировать
                </button>
                <button
                  onClick={() => window.open(`http://localhost:3000/pets/${pet.id}`, '_blank')}
                  className="px-6 py-3 bg-white/20 backdrop-blur-sm text-white rounded-xl hover:bg-white/30 transition-all duration-200 font-semibold flex items-center gap-2 hover:scale-105"
                >
                  <ArrowTopRightOnSquareIcon className="w-5 h-5" />
                  Открыть на сайте
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Статистика */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105">
          <div className="flex items-center justify-between mb-2">
            <ClockIcon className="w-8 h-8 opacity-80" />
            <span className="text-3xl font-bold">
              {pet.birth_date ? calculateAge(pet.birth_date) : '—'}
            </span>
          </div>
          <div className="text-sm opacity-90 font-medium">Возраст</div>
        </div>

        <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105">
          <div className="flex items-center justify-between mb-2">
            <ScaleIcon className="w-8 h-8 opacity-80" />
            <span className="text-3xl font-bold">{pet.weight || '—'}</span>
          </div>
          <div className="text-sm opacity-90 font-medium">Вес (кг)</div>
        </div>

        <div className={`bg-gradient-to-br ${pet.is_sterilized ? 'from-purple-500 to-purple-600' : 'from-gray-400 to-gray-500'} rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105`}>
          <div className="flex items-center justify-between mb-2">
            <ShieldCheckIcon className="w-8 h-8 opacity-80" />
            <span className="text-3xl">{pet.is_sterilized ? '✓' : '✗'}</span>
          </div>
          <div className="text-sm opacity-90 font-medium">Стерилизация</div>
        </div>

        <div className={`bg-gradient-to-br ${pet.is_vaccinated ? 'from-pink-500 to-pink-600' : 'from-gray-400 to-gray-500'} rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105`}>
          <div className="flex items-center justify-between mb-2">
            <BeakerIcon className="w-8 h-8 opacity-80" />
            <span className="text-3xl">{pet.is_vaccinated ? '✓' : '✗'}</span>
          </div>
          <div className="text-sm opacity-90 font-medium">Вакцинация</div>
        </div>
      </div>

      {/* Вкладки */}
      <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
        <div className="flex border-b border-gray-200 overflow-x-auto">
          <button
            onClick={() => setActiveTab('overview')}
            className={`flex-1 min-w-fit px-6 py-4 font-semibold text-sm transition-all ${
              activeTab === 'overview'
                ? 'bg-gradient-to-r from-blue-500 to-purple-500 text-white shadow-lg'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            }`}
          >
            <div className="flex items-center justify-center gap-2">
              <span className="text-xl">📋</span>
              <span>Обзор</span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('medical')}
            className={`flex-1 min-w-fit px-6 py-4 font-semibold text-sm transition-all ${
              activeTab === 'medical'
                ? 'bg-gradient-to-r from-green-500 to-blue-500 text-white shadow-lg'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            }`}
          >
            <div className="flex items-center justify-center gap-2">
              <span className="text-xl">💉</span>
              <span>Медицина</span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('identification')}
            className={`flex-1 min-w-fit px-6 py-4 font-semibold text-sm transition-all ${
              activeTab === 'identification'
                ? 'bg-gradient-to-r from-indigo-500 to-purple-500 text-white shadow-lg'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            }`}
          >
            <div className="flex items-center justify-center gap-2">
              <span className="text-xl">🆔</span>
              <span>Идентификация</span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('contacts')}
            className={`flex-1 min-w-fit px-6 py-4 font-semibold text-sm transition-all ${
              activeTab === 'contacts'
                ? 'bg-gradient-to-r from-orange-500 to-red-500 text-white shadow-lg'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            }`}
          >
            <div className="flex items-center justify-center gap-2">
              <span className="text-xl">🚨</span>
              <span>Контакты</span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('story')}
            className={`flex-1 min-w-fit px-6 py-4 font-semibold text-sm transition-all ${
              activeTab === 'story'
                ? 'bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-lg'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            }`}
          >
            <div className="flex items-center justify-center gap-2">
              <span className="text-xl">📖</span>
              <span>История</span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('photos')}
            className={`flex-1 min-w-fit px-6 py-4 font-semibold text-sm transition-all ${
              activeTab === 'photos'
                ? 'bg-gradient-to-r from-pink-500 to-rose-500 text-white shadow-lg'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            }`}
          >
            <div className="flex items-center justify-center gap-2">
              <span className="text-xl">📸</span>
              <span>Фото</span>
            </div>
          </button>
        </div>
      </div>

      {/* Контент вкладок */}
      {activeTab === 'overview' && (
        <>
          {/* Основная информация */}
      <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
        <div className="bg-gradient-to-r from-blue-50 to-purple-50 px-6 py-4 border-b border-gray-100">
          <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <SparklesIcon className="w-6 h-6 text-blue-600" />
            Основная информация
          </h2>
        </div>
        <div className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="group">
              <div className="text-sm text-gray-500 mb-1 font-medium">Кличка</div>
              <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{pet.name}</div>
            </div>
            <div className="group">
              <div className="text-sm text-gray-500 mb-1 font-medium">Вид</div>
              <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors flex items-center gap-2">
                <span>{getSpeciesEmoji(pet.species)}</span>
                {getSpeciesLabel(pet.species)}
              </div>
            </div>
            {pet.breed && (
              <div className="group">
                <div className="text-sm text-gray-500 mb-1 font-medium">Порода</div>
                <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{pet.breed}</div>
              </div>
            )}
            <div className="group">
              <div className="text-sm text-gray-500 mb-1 font-medium">Пол</div>
              <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors flex items-center gap-2">
                <span>{pet.gender === 'male' ? '♂️' : pet.gender === 'female' ? '♀️' : '—'}</span>
                {getGenderLabel(pet.gender)}
              </div>
            </div>
            {pet.birth_date && (
              <div className="group">
                <div className="text-sm text-gray-500 mb-1 font-medium">Дата рождения</div>
                <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">
                  {new Date(pet.birth_date).toLocaleDateString('ru-RU')}
                  <span className="text-sm text-gray-500 ml-2">({calculateAge(pet.birth_date)})</span>
                </div>
              </div>
            )}
            {pet.color && (
              <div className="group">
                <div className="text-sm text-gray-500 mb-1 font-medium">Окрас</div>
                <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{pet.color}</div>
              </div>
            )}
            {pet.size && (
              <div className="group">
                <div className="text-sm text-gray-500 mb-1 font-medium">Размер</div>
                <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{getSizeLabel(pet.size)}</div>
              </div>
            )}
            {pet.weight && (
              <div className="group">
                <div className="text-sm text-gray-500 mb-1 font-medium">Вес</div>
                <div className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{pet.weight} кг</div>
              </div>
            )}
            {pet.chip_number && (
              <div className="group col-span-2">
                <div className="text-sm text-gray-500 mb-1 font-medium">Номер чипа</div>
                <div className="text-lg font-semibold text-gray-900 font-mono bg-gray-50 px-4 py-2 rounded-lg group-hover:bg-blue-50 transition-colors">{pet.chip_number}</div>
              </div>
            )}
          </div>
        </div>
      </div>

          {/* Медицинская информация */}
          <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
            <div className="bg-gradient-to-r from-green-50 to-blue-50 px-6 py-4 border-b border-gray-100">
              <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                <HeartIcon className="w-6 h-6 text-green-600" />
                Характер и особенности
              </h2>
            </div>
            <div className="p-6">
              {(pet.health_notes || pet.character_traits || pet.special_needs) ? (
                <div className="space-y-4">
                  {pet.health_notes && (
                    <div className="bg-red-50 border-l-4 border-red-400 p-4 rounded-r-lg">
                      <div className="text-sm text-red-700 font-semibold mb-2 flex items-center gap-2">
                        <BeakerIcon className="w-5 h-5" />
                        Заметки о здоровье
                      </div>
                      <div className="text-base text-gray-900">{pet.health_notes}</div>
                    </div>
                  )}

                  {pet.character_traits && (
                    <div className="bg-purple-50 border-l-4 border-purple-400 p-4 rounded-r-lg">
                      <div className="text-sm text-purple-700 font-semibold mb-2 flex items-center gap-2">
                        <SparklesIcon className="w-5 h-5" />
                        Особенности характера
                      </div>
                      <div className="text-base text-gray-900">{pet.character_traits}</div>
                    </div>
                  )}

                  {pet.special_needs && (
                    <div className="bg-orange-50 border-l-4 border-orange-400 p-4 rounded-r-lg">
                      <div className="text-sm text-orange-700 font-semibold mb-2 flex items-center gap-2">
                        <ShieldCheckIcon className="w-5 h-5" />
                        Особые потребности
                      </div>
                      <div className="text-base text-gray-900">{pet.special_needs}</div>
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">💚</div>
                  <p className="text-gray-500 text-lg">Дополнительная информация не указана</p>
                  <button
                    onClick={() => router.push(`/pets/${pet.id}/edit`)}
                    className="mt-4 text-blue-600 hover:text-blue-700 font-medium"
                  >
                    Добавить информацию
                  </button>
                </div>
              )}
            </div>
          </div>
        </>
      )}

      {activeTab === 'medical' && (
        <>
          {/* Медицинские детали */}
          {(pet.blood_type || pet.allergies || pet.chronic_diseases || pet.current_medications || pet.distinctive_marks || pet.sterilization_date) ? (
        <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
          <div className="bg-gradient-to-r from-red-50 to-pink-50 px-6 py-4 border-b border-gray-100">
            <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
              <BeakerIcon className="w-6 h-6 text-red-600" />
              Медицинские детали
            </h2>
          </div>
          <div className="p-6 space-y-4">
            {pet.blood_type && (
              <div className="flex items-start gap-3">
                <div className="w-32 text-sm text-gray-600 font-medium pt-1">Группа крови:</div>
                <div className="flex-1 text-base text-gray-900 font-semibold">{pet.blood_type}</div>
              </div>
            )}
            {pet.allergies && (
              <div className="flex items-start gap-3">
                <div className="w-32 text-sm text-gray-600 font-medium pt-1">Аллергии:</div>
                <div className="flex-1 text-base text-gray-900 bg-red-50 p-3 rounded-lg">{pet.allergies}</div>
              </div>
            )}
            {pet.chronic_diseases && (
              <div className="flex items-start gap-3">
                <div className="w-32 text-sm text-gray-600 font-medium pt-1">Хронические заболевания:</div>
                <div className="flex-1 text-base text-gray-900 bg-orange-50 p-3 rounded-lg">{pet.chronic_diseases}</div>
              </div>
            )}
            {pet.current_medications && (
              <div className="flex items-start gap-3">
                <div className="w-32 text-sm text-gray-600 font-medium pt-1">Текущие препараты:</div>
                <div className="flex-1 text-base text-gray-900 bg-blue-50 p-3 rounded-lg">{pet.current_medications}</div>
              </div>
            )}
            {pet.distinctive_marks && (
              <div className="flex items-start gap-3">
                <div className="w-32 text-sm text-gray-600 font-medium pt-1">Особые приметы:</div>
                <div className="flex-1 text-base text-gray-900 bg-purple-50 p-3 rounded-lg">{pet.distinctive_marks}</div>
              </div>
            )}
            {pet.sterilization_date && (
              <div className="flex items-start gap-3">
                <div className="w-32 text-sm text-gray-600 font-medium pt-1">Дата стерилизации:</div>
                <div className="flex-1 text-base text-gray-900 font-semibold">
                  {new Date(pet.sterilization_date).toLocaleDateString('ru-RU')}
                </div>
              </div>
            )}
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div className="text-center py-12">
            <div className="text-6xl mb-4">💉</div>
            <p className="text-gray-500 text-lg">Медицинская информация не указана</p>
            <button
              onClick={() => router.push(`/pets/${pet.id}/edit`)}
              className="mt-4 text-blue-600 hover:text-blue-700 font-medium"
            >
              Добавить информацию
            </button>
          </div>
        </div>
      )}
        </>
      )}

      {activeTab === 'identification' && (
        <>
          {/* Идентификация */}
          {(pet.chip_number || pet.tattoo_number || pet.ear_tag_number || pet.passport_number || pet.pedigree_number || pet.registration_org) ? (
            <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
              <div className="bg-gradient-to-r from-indigo-50 to-blue-50 px-6 py-4 border-b border-gray-100">
                <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                  <ShieldCheckIcon className="w-6 h-6 text-indigo-600" />
                  Идентификация
                </h2>
              </div>
              <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {pet.chip_number && (
                    <div className="bg-indigo-50 p-4 rounded-lg">
                      <div className="text-sm text-indigo-700 font-semibold mb-1">Номер чипа</div>
                      <div className="text-lg font-mono text-gray-900">{pet.chip_number}</div>
                    </div>
                  )}
                  {pet.tattoo_number && (
                    <div className="bg-purple-50 p-4 rounded-lg">
                      <div className="text-sm text-purple-700 font-semibold mb-1">Клеймо</div>
                      <div className="text-lg font-mono text-gray-900">{pet.tattoo_number}</div>
                    </div>
                  )}
                  {pet.ear_tag_number && (
                    <div className="bg-pink-50 p-4 rounded-lg">
                      <div className="text-sm text-pink-700 font-semibold mb-1">Номер бирки</div>
                      <div className="text-lg font-mono text-gray-900">{pet.ear_tag_number}</div>
                    </div>
                  )}
                  {pet.passport_number && (
                    <div className="bg-blue-50 p-4 rounded-lg">
                      <div className="text-sm text-blue-700 font-semibold mb-1">Номер паспорта</div>
                      <div className="text-lg font-mono text-gray-900">{pet.passport_number}</div>
                    </div>
                  )}
                  {pet.pedigree_number && (
                    <div className="bg-amber-50 p-4 rounded-lg">
                      <div className="text-sm text-amber-700 font-semibold mb-1">Номер родословной</div>
                      <div className="text-lg font-mono text-gray-900">{pet.pedigree_number}</div>
                    </div>
                  )}
                  {pet.registration_org && (
                    <div className="bg-green-50 p-4 rounded-lg">
                      <div className="text-sm text-green-700 font-semibold mb-1">Организация регистрации</div>
                      <div className="text-base text-gray-900">{pet.registration_org}</div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
              <div className="text-center py-12">
                <div className="text-6xl mb-4">🆔</div>
                <p className="text-gray-500 text-lg">Информация об идентификации не указана</p>
                <button
                  onClick={() => router.push(`/pets/${pet.id}/edit`)}
                  className="mt-4 text-blue-600 hover:text-blue-700 font-medium"
                >
                  Добавить информацию
                </button>
              </div>
            </div>
          )}
        </>
      )}

      {activeTab === 'contacts' && (
        <>
          {/* Экстренные контакты */}
          {(pet.emergency_contact_name || pet.vet_clinic_name || pet.insurance_company) ? (
            <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
              <div className="bg-gradient-to-r from-orange-50 to-red-50 px-6 py-4 border-b border-gray-100">
                <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                  <span className="text-2xl">🚨</span>
                  Экстренные контакты
                </h2>
              </div>
              <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {pet.emergency_contact_name && (
                    <div className="bg-red-50 p-4 rounded-lg border-l-4 border-red-400">
                      <div className="text-sm text-red-700 font-semibold mb-2">Экстренный контакт</div>
                      <div className="space-y-1">
                        <div className="text-base font-semibold text-gray-900">{pet.emergency_contact_name}</div>
                        {pet.emergency_contact_phone && (
                          <div className="text-base text-gray-700">📞 {pet.emergency_contact_phone}</div>
                        )}
                        {pet.emergency_contact_relation && (
                          <div className="text-sm text-gray-600">{pet.emergency_contact_relation}</div>
                        )}
                      </div>
                    </div>
                  )}
                  
                  {pet.vet_clinic_name && (
                    <div className="bg-blue-50 p-4 rounded-lg border-l-4 border-blue-400">
                      <div className="text-sm text-blue-700 font-semibold mb-2">Ветеринарная клиника</div>
                      <div className="space-y-1">
                        <div className="text-base font-semibold text-gray-900">{pet.vet_clinic_name}</div>
                        {pet.vet_clinic_phone && (
                          <div className="text-base text-gray-700">📞 {pet.vet_clinic_phone}</div>
                        )}
                        {pet.vet_clinic_address && (
                          <div className="text-sm text-gray-600">📍 {pet.vet_clinic_address}</div>
                        )}
                      </div>
                    </div>
                  )}
                  
                  {pet.insurance_company && (
                    <div className="bg-green-50 p-4 rounded-lg border-l-4 border-green-400 md:col-span-2">
                      <div className="text-sm text-green-700 font-semibold mb-2">Страхование</div>
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div>
                          <div className="text-xs text-gray-600 mb-1">Компания</div>
                          <div className="text-base font-semibold text-gray-900">{pet.insurance_company}</div>
                        </div>
                        {pet.insurance_policy_number && (
                          <div>
                            <div className="text-xs text-gray-600 mb-1">Номер полиса</div>
                            <div className="text-base font-mono text-gray-900">{pet.insurance_policy_number}</div>
                          </div>
                        )}
                        {pet.insurance_expiry_date && (
                          <div>
                            <div className="text-xs text-gray-600 mb-1">Действует до</div>
                            <div className="text-base font-semibold text-gray-900">
                              {new Date(pet.insurance_expiry_date).toLocaleDateString('ru-RU')}
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
              <div className="text-center py-12">
                <div className="text-6xl mb-4">🚨</div>
                <p className="text-gray-500 text-lg">Экстренные контакты не указаны</p>
                <button
                  onClick={() => router.push(`/pets/${pet.id}/edit`)}
                  className="mt-4 text-blue-600 hover:text-blue-700 font-medium"
                >
                  Добавить информацию
                </button>
              </div>
            </div>
          )}
        </>
      )}

      {activeTab === 'story' && (
        <>
          {/* История питомца */}
          {pet.story ? (
            <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
              <div className="bg-gradient-to-r from-amber-50 to-orange-50 px-6 py-4 border-b border-gray-100">
                <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                  <span className="text-2xl">📖</span>
                  История питомца
                </h2>
              </div>
              <div className="p-6">
                <div className="prose prose-lg max-w-none text-gray-700 leading-relaxed">
                  {pet.story}
                </div>
              </div>
            </div>
          ) : (
            <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
              <div className="text-center py-12">
                <div className="text-6xl mb-4">📖</div>
                <p className="text-gray-500 text-lg">История питомца не указана</p>
                <button
                  onClick={() => router.push(`/pets/${pet.id}/edit`)}
                  className="mt-4 text-blue-600 hover:text-blue-700 font-medium"
                >
                  Добавить историю
                </button>
              </div>
            </div>
          )}
        </>
      )}

      {activeTab === 'photos' && (
        <PhotoUploadTab petId={pet.id} currentPhoto={pet.photo} onPhotoUpdate={loadPet} />
      )}
    </div>
  );
}

// Компонент для загрузки фото
function PhotoUploadTab({ petId, currentPhoto, onPhotoUpdate }: { petId: number; currentPhoto?: string; onPhotoUpdate: () => void }) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [preview, setPreview] = useState<string | null>(null);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    console.log('📸 [Frontend] File selected:', file.name, 'size:', file.size, 'type:', file.type);

    // Валидация
    if (!file.type.startsWith('image/')) {
      console.error('❌ [Frontend] Invalid file type:', file.type);
      setError('Пожалуйста, выберите изображение');
      return;
    }

    if (file.size > 10 * 1024 * 1024) {
      console.error('❌ [Frontend] File too large:', file.size);
      setError('Размер файла не должен превышать 10 МБ');
      return;
    }

    setError(null);
    setUploading(true);

    try {
      console.log('📸 [Frontend] Creating preview...');
      // Создаём preview
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
        console.log('✅ [Frontend] Preview created');
      };
      reader.readAsDataURL(file);

      // Загружаем файл
      console.log('📸 [Frontend] Creating FormData...');
      const formData = new FormData();
      formData.append('photo', file);
      console.log('📸 [Frontend] FormData created, uploading to:', `http://localhost:8400/api/pets/${petId}/photo`);

      const response = await fetch(`http://localhost:8400/api/pets/${petId}/photo`, {
        method: 'POST',
        body: formData,
        credentials: 'include',
      });

      console.log('📸 [Frontend] Response status:', response.status, response.statusText);
      console.log('📸 [Frontend] Response headers:', Object.fromEntries(response.headers.entries()));

      if (response.ok) {
        const data = await response.json();
        console.log('✅ [Frontend] Upload successful:', data);
        // Обновляем данные питомца
        onPhotoUpdate();
        setError(null);
      } else {
        const errorText = await response.text();
        console.error('❌ [Frontend] Upload failed:', response.status, errorText);
        const errorData = JSON.parse(errorText || '{}');
        setError(errorData.error || 'Не удалось загрузить фото');
        setPreview(null);
      }
    } catch (err) {
      console.error('❌ [Frontend] Upload error:', err);
      setError('Ошибка при загрузке фото');
      setPreview(null);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
      <div className="bg-gradient-to-r from-pink-50 to-rose-50 px-6 py-4 border-b border-gray-100">
        <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <span className="text-2xl">📸</span>
          Фотография питомца
        </h2>
      </div>
      
      <div className="p-6">
        {/* Текущее фото */}
        <div className="mb-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Текущее фото</h3>
          <div className="flex justify-center">
            {currentPhoto || preview ? (
              <img
                src={preview || (currentPhoto?.startsWith('http') ? currentPhoto : `http://localhost:8400${currentPhoto}`)}
                alt="Фото питомца"
                className="w-64 h-64 rounded-2xl object-cover shadow-lg"
              />
            ) : (
              <div className="w-64 h-64 rounded-2xl bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
                <div className="text-center">
                  <div className="text-6xl mb-2">📷</div>
                  <p className="text-gray-500">Фото не загружено</p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Форма загрузки */}
        <div className="border-2 border-dashed border-gray-300 rounded-xl p-8 text-center hover:border-blue-400 transition-colors">
          <input
            type="file"
            id="photo-upload"
            accept="image/*"
            onChange={handleFileSelect}
            disabled={uploading}
            className="hidden"
          />
          <label
            htmlFor="photo-upload"
            className={`cursor-pointer ${uploading ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            <div className="text-6xl mb-4">📤</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              {uploading ? 'Загрузка...' : 'Загрузить новое фото'}
            </h3>
            <p className="text-gray-600 mb-4">
              Нажмите или перетащите файл сюда
            </p>
            <p className="text-sm text-gray-500">
              Поддерживаются: JPG, PNG, WEBP (макс. 10 МБ)
            </p>
          </label>
        </div>

        {/* Ошибка */}
        {error && (
          <div className="mt-4 bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
            {error}
          </div>
        )}

        {/* Подсказки */}
        <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h4 className="font-semibold text-blue-900 mb-2">💡 Советы для лучшего фото:</h4>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>• Фотографируйте при хорошем освещении</li>
            <li>• Питомец должен быть в фокусе</li>
            <li>• Избегайте размытых изображений</li>
            <li>• Лучше всего подходят портретные фото</li>
          </ul>
        </div>

        {/* Галерея фотографий */}
        <div className="mt-6 bg-gray-50 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <span className="text-2xl">🖼️</span>
            Галерея фотографий
          </h3>
          <div className="text-sm text-gray-600 text-center py-8">
            <div className="text-4xl mb-2">📷</div>
            <div>Галерея фотографий будет доступна в следующих версиях</div>
          </div>
        </div>
      </div>
    </div>
  );
}
