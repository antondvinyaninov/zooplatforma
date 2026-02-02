'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';

export default function CreatePetPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const [formData, setFormData] = useState({
    name: '',
    species: 'dog',
    breed: '',
    gender: '',
    birth_date: '',
    color: '',
    size: '',
    location: '',
    city: '',
    region: '',
    contact_name: '',
    contact_phone: '',
    story: '',
    character_traits: '',
    health_notes: '',
    urgent: false,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      console.log('🔵 Step 1: Getting auth token from Auth Service...');
      
      // Получаем токен через Auth Service
      const authResponse = await fetch('http://localhost:7100/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      console.log('🔵 Step 2: Auth Service response status:', authResponse.status);

      if (!authResponse.ok) {
        console.error('❌ Auth Service returned error:', authResponse.status);
        throw new Error('Необходима авторизация');
      }

      const authResult = await authResponse.json();
      console.log('🔵 Step 3: Auth result:', { 
        success: authResult.success, 
        hasToken: !!authResult.data?.token,
        hasUser: !!authResult.data?.user 
      });
      
      if (!authResult.success || !authResult.data?.token) {
        console.error('❌ No token in auth response');
        throw new Error('Не удалось получить токен авторизации');
      }

      const token = authResult.data.token;
      console.log('✅ Got token:', token.substring(0, 20) + '...');

      console.log('🔵 Step 4: Creating pet via Volunteer API...');
      console.log('🔵 Request data:', formData);

      // Создаем животное с токеном
      const response = await fetch('http://localhost:8500/api/pets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        credentials: 'include',
        body: JSON.stringify(formData),
      });

      console.log('🔵 Step 5: Volunteer API response status:', response.status);

      const result = await response.json();
      console.log('🔵 Step 6: Volunteer API result:', result);

      if (!response.ok || !result.success) {
        console.error('❌ Volunteer API error:', result);
        throw new Error(result.error || 'Не удалось создать животное');
      }

      console.log('✅ Pet created successfully!');

      // Успешно создано - возвращаемся к списку
      router.push('/dashboard');
      // Программно переключаем на вкладку pets
      if (typeof window !== 'undefined') {
        localStorage.setItem('volunteerActiveTab', 'pets');
      }
    } catch (err) {
      console.error('❌ Error in handleSubmit:', err);
      setError(err instanceof Error ? err.message : 'Произошла ошибка');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    
    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData(prev => ({ ...prev, [name]: checked }));
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <button
          onClick={() => router.push('/dashboard')}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors mb-4"
        >
          <ArrowLeftIcon className="w-5 h-5" />
          Назад к списку
        </button>
        <h2 className="text-3xl font-bold text-gray-900">Добавить подопечного</h2>
        <p className="text-gray-600 mt-1">Заполните информацию о животном</p>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 space-y-8">
        {error && (
          <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-800">
            {error}
          </div>
        )}

        {/* Основная информация */}
        <div className="space-y-6">
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-1">Основная информация</h3>
            <p className="text-sm text-gray-600">Базовые данные о животном</p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Кличка <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Например: Барсик"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Вид <span className="text-red-500">*</span>
              </label>
              <select
                name="species"
                value={formData.species}
                onChange={handleChange}
                required
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              >
                <option value="dog">Собака</option>
                <option value="cat">Кошка</option>
                <option value="bird">Птица</option>
                <option value="rodent">Грызун</option>
                <option value="other">Другое</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Порода
              </label>
              <input
                type="text"
                name="breed"
                value={formData.breed}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Например: Дворняга"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Пол
              </label>
              <select
                name="gender"
                value={formData.gender}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              >
                <option value="">Не указан</option>
                <option value="male">Самец</option>
                <option value="female">Самка</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Дата рождения (примерная)
              </label>
              <input
                type="date"
                name="birth_date"
                value={formData.birth_date}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Окрас
              </label>
              <input
                type="text"
                name="color"
                value={formData.color}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Например: Рыжий"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Размер
              </label>
              <select
                name="size"
                value={formData.size}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              >
                <option value="">Не указан</option>
                <option value="small">Маленький</option>
                <option value="medium">Средний</option>
                <option value="large">Крупный</option>
              </select>
            </div>
          </div>
        </div>

        {/* Местоположение */}
        <div className="space-y-6 pt-6 border-t border-gray-200">
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-1">Местоположение</h3>
            <p className="text-sm text-gray-600">Где находится животное</p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Город
              </label>
              <input
                type="text"
                name="city"
                value={formData.city}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Например: Ижевск"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Регион
              </label>
              <input
                type="text"
                name="region"
                value={formData.region}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Например: Удмуртская Республика"
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Адрес передержки
              </label>
              <input
                type="text"
                name="location"
                value={formData.location}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Например: ул. Ленина, 10"
              />
            </div>
          </div>
        </div>

        {/* Контакты */}
        <div className="space-y-6 pt-6 border-t border-gray-200">
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-1">Контакты для связи</h3>
            <p className="text-sm text-gray-600">Как с вами связаться по поводу животного</p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Имя для связи
              </label>
              <input
                type="text"
                name="contact_name"
                value={formData.contact_name}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Ваше имя"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Телефон для связи
              </label>
              <input
                type="tel"
                name="contact_phone"
                value={formData.contact_phone}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="+7 (999) 123-45-67"
              />
            </div>
          </div>
        </div>

        {/* Описание */}
        <div className="space-y-6 pt-6 border-t border-gray-200">
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-1">Описание</h3>
            <p className="text-sm text-gray-600">Расскажите о животном</p>
          </div>
          
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                История животного
              </label>
              <textarea
                name="story"
                value={formData.story}
                onChange={handleChange}
                rows={4}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Расскажите историю животного..."
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Характер
              </label>
              <textarea
                name="character_traits"
                value={formData.character_traits}
                onChange={handleChange}
                rows={3}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Опишите характер животного..."
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Здоровье
              </label>
              <textarea
                name="health_notes"
                value={formData.health_notes}
                onChange={handleChange}
                rows={3}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Информация о здоровье..."
              />
            </div>
          </div>
        </div>

        {/* Срочность */}
        <div className="pt-6 border-t border-gray-200">
          <label className="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              name="urgent"
              checked={formData.urgent}
              onChange={handleChange}
              className="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <div>
              <div className="text-sm font-medium text-gray-900">Срочно нужна помощь</div>
              <div className="text-sm text-gray-600">Животное будет отмечено как срочное в каталоге</div>
            </div>
          </label>
        </div>

        {/* Buttons */}
        <div className="flex gap-4 pt-6 border-t border-gray-200">
          <button
            type="button"
            onClick={() => router.push('/dashboard')}
            className="flex-1 px-6 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
          >
            Отмена
          </button>
          <button
            type="submit"
            disabled={loading}
            className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed font-medium"
          >
            {loading ? 'Создание...' : 'Создать подопечного'}
          </button>
        </div>
      </form>
    </div>
  );
}
