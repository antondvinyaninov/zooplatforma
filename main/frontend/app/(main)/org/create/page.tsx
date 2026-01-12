'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { 
  BuildingOfficeIcon,
  MagnifyingGlassIcon,
  CheckCircleIcon,
  XCircleIcon
} from '@heroicons/react/24/outline';
import { organizationsApi } from '@/lib/api';

interface DaDataSuggestion {
  value: string;
  unrestricted_value: string;
  data: {
    inn: string;
    ogrn: string;
    kpp?: string;
    name: {
      full: string;
      short?: string;
    };
    address: {
      value: string;
      data: {
        postal_code?: string;
        region?: string;
        city?: string;
        street?: string;
        house?: string;
        flat?: string;
        geo_lat?: string;
        geo_lon?: string;
      };
    };
    management?: {
      name?: string;
      post?: string;
    };
    emails?: Array<{ value: string }>;
    phones?: Array<{ value: string }>;
    opf?: {
      full?: string;
    };
  };
}

export default function CreateOrganizationPage() {
  const router = useRouter();
  const [step, setStep] = useState<'search' | 'form'>('search');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  // Поиск по ИНН
  const [inn, setInn] = useState('');
  const [suggestions, setSuggestions] = useState<DaDataSuggestion[]>([]);
  const [searching, setSearching] = useState(false);
  
  // Данные формы
  const [formData, setFormData] = useState({
    name: '',
    short_name: '',
    legal_form: '',
    type: 'other',
    inn: '',
    ogrn: '',
    kpp: '',
    email: '',
    phone: '',
    website: '',
    address_full: '',
    address_postal_code: '',
    address_region: '',
    address_city: '',
    address_street: '',
    address_house: '',
    address_office: '',
    geo_lat: null as number | null,
    geo_lon: null as number | null,
    description: '',
    bio: '',
    director_name: '',
    director_position: '',
  });

  // Поиск организации по ИНН через DaData
  const searchByInn = async () => {
    if (!inn || inn.length < 10) {
      setError('ИНН должен содержать минимум 10 цифр');
      return;
    }

    setSearching(true);
    setError('');
    
    try {
      const response = await fetch('https://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Token ${process.env.NEXT_PUBLIC_DADATA_API_KEY}`,
        },
        body: JSON.stringify({ query: inn }),
      });

      const data = await response.json();
      
      if (data.suggestions && data.suggestions.length > 0) {
        setSuggestions(data.suggestions);
      } else {
        setError('Организация с таким ИНН не найдена');
        setSuggestions([]);
      }
    } catch (err) {
      setError('Ошибка при поиске организации');
      console.error(err);
    } finally {
      setSearching(false);
    }
  };

  // Выбор организации из списка
  const selectOrganization = (suggestion: DaDataSuggestion) => {
    const data = suggestion.data;
    
    setFormData({
      name: data.name.full,
      short_name: data.name.short || '',
      legal_form: data.opf?.full || '',
      type: 'other',
      inn: data.inn,
      ogrn: data.ogrn,
      kpp: data.kpp || '',
      email: data.emails?.[0]?.value || '',
      phone: data.phones?.[0]?.value || '',
      website: '',
      address_full: data.address.value,
      address_postal_code: data.address.data.postal_code || '',
      address_region: data.address.data.region || '',
      address_city: data.address.data.city || '',
      address_street: data.address.data.street || '',
      address_house: data.address.data.house || '',
      address_office: data.address.data.flat || '',
      geo_lat: data.address.data.geo_lat ? parseFloat(data.address.data.geo_lat) : null,
      geo_lon: data.address.data.geo_lon ? parseFloat(data.address.data.geo_lon) : null,
      description: '',
      bio: '',
      director_name: data.management?.name || '',
      director_position: data.management?.post || '',
    });
    
    setStep('form');
  };

  // Создание организации вручную
  const createManually = () => {
    setStep('form');
  };

  // Отправка формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name || !formData.inn) {
      setError('Заполните обязательные поля: название и ИНН');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await organizationsApi.create(formData);
      
      if (response.success && response.data) {
        // Перенаправляем на страницу созданной организации
        router.push(`/org/${response.data.id}`);
      } else {
        setError(response.error || 'Ошибка при создании организации');
      }
    } catch (err) {
      setError('Ошибка при создании организации');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (step === 'search') {
    return (
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <div className="flex items-center gap-3 mb-6">
            <BuildingOfficeIcon className="w-8 h-8 text-blue-500" />
            <h1 className="text-2xl font-bold text-gray-900">Создать организацию</h1>
          </div>

          <p className="text-gray-600 mb-6">
            Введите ИНН организации для автоматического заполнения данных из ЕГРЮЛ
          </p>

          {/* Поиск по ИНН */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              ИНН организации
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={inn}
                onChange={(e) => setInn(e.target.value.replace(/\D/g, ''))}
                placeholder="1234567890"
                maxLength={12}
                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                onClick={searchByInn}
                disabled={searching || inn.length < 10}
                className="px-6 py-2 bg-blue-500 text-white rounded-lg font-medium hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed flex items-center gap-2"
              >
                {searching ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    Поиск...
                  </>
                ) : (
                  <>
                    <MagnifyingGlassIcon className="w-5 h-5" />
                    Найти
                  </>
                )}
              </button>
            </div>
            {error && (
              <div className="mt-2 flex items-center gap-2 text-red-600 text-sm">
                <XCircleIcon className="w-4 h-4" />
                {error}
              </div>
            )}
          </div>

          {/* Результаты поиска */}
          {suggestions.length > 0 && (
            <div className="mb-6">
              <h3 className="text-sm font-medium text-gray-700 mb-3">Найденные организации:</h3>
              <div className="space-y-2">
                {suggestions.map((suggestion, index) => (
                  <button
                    key={index}
                    onClick={() => selectOrganization(suggestion)}
                    className="w-full text-left p-4 border border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
                  >
                    <div className="flex items-start gap-3">
                      <CheckCircleIcon className="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" />
                      <div className="flex-1 min-w-0">
                        <div className="font-medium text-gray-900 mb-1">
                          {suggestion.data.name.full}
                        </div>
                        <div className="text-sm text-gray-600">
                          ИНН: {suggestion.data.inn} • {suggestion.data.address.value}
                        </div>
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Разделитель */}
          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-200"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-white text-gray-500">или</span>
            </div>
          </div>

          {/* Создать вручную */}
          <button
            onClick={createManually}
            className="w-full px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
          >
            Заполнить данные вручную
          </button>
        </div>
      </div>
    );
  }

  // Форма создания
  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
        <div className="flex items-center gap-3 mb-6">
          <BuildingOfficeIcon className="w-8 h-8 text-blue-500" />
          <h1 className="text-2xl font-bold text-gray-900">Создать организацию</h1>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Основная информация */}
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Основная информация</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Полное название <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Краткое название
                </label>
                <input
                  type="text"
                  value={formData.short_name}
                  onChange={(e) => setFormData({ ...formData, short_name: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Тип организации <span className="text-red-500">*</span>
                </label>
                <select
                  value={formData.type}
                  onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                  required
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="shelter">Приют</option>
                  <option value="vet_clinic">Ветеринарная клиника</option>
                  <option value="pet_shop">Зоомагазин</option>
                  <option value="foundation">Благотворительный фонд</option>
                  <option value="kennel">Кинологический центр</option>
                  <option value="other">Другое</option>
                </select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    ИНН <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={formData.inn}
                    onChange={(e) => setFormData({ ...formData, inn: e.target.value.replace(/\D/g, '') })}
                    required
                    maxLength={12}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    ОГРН
                  </label>
                  <input
                    type="text"
                    value={formData.ogrn}
                    onChange={(e) => setFormData({ ...formData, ogrn: e.target.value.replace(/\D/g, '') })}
                    maxLength={15}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Контакты */}
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Контактная информация</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Телефон
                </label>
                <input
                  type="tel"
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Веб-сайт
                </label>
                <input
                  type="url"
                  value={formData.website}
                  onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                  placeholder="https://example.com"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Адрес */}
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Адрес</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Полный адрес
                </label>
                <input
                  type="text"
                  value={formData.address_full}
                  onChange={(e) => setFormData({ ...formData, address_full: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Регион
                  </label>
                  <input
                    type="text"
                    value={formData.address_region}
                    onChange={(e) => setFormData({ ...formData, address_region: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Город
                  </label>
                  <input
                    type="text"
                    value={formData.address_city}
                    onChange={(e) => setFormData({ ...formData, address_city: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Описание */}
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">О организации</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Краткое описание
                </label>
                <textarea
                  value={formData.bio}
                  onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                  rows={2}
                  maxLength={200}
                  placeholder="Краткое описание для карточки организации (до 200 символов)"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                />
                <div className="text-xs text-gray-500 mt-1">
                  {formData.bio.length}/200
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Полное описание
                </label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  rows={4}
                  placeholder="Подробное описание деятельности организации"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                />
              </div>
            </div>
          </div>

          {error && (
            <div className="p-4 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2 text-red-700">
              <XCircleIcon className="w-5 h-5 flex-shrink-0" />
              {error}
            </div>
          )}

          {/* Кнопки */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={() => router.back()}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg font-medium hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? 'Создание...' : 'Создать организацию'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
