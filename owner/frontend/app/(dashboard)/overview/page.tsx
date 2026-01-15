'use client';

import { useState, useEffect } from 'react';
import {
  HomeIcon,
  HeartIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

interface Pet {
  id: number;
  name: string;
}

export default function OverviewPage() {
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPets();
  }, []);

  const loadPets = async () => {
    try {
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        credentials: 'include',
      });

      if (meResponse.ok) {
        const meResult = await meResponse.json();
        const userId = meResult.data.id;

        const petsResponse = await fetch(`http://localhost:8100/api/pets?user_id=${userId}`, {
          headers: {
            'X-User-ID': userId.toString(),
          },
          credentials: 'include',
        });

        if (petsResponse.ok) {
          const petsResult = await petsResponse.json();
          if (petsResult.success && petsResult.data) {
            setPets(petsResult.data);
          }
        }
      }
    } catch (error) {
      console.error('Error loading pets:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-4xl font-bold text-gray-900 mb-2">Добро пожаловать!</h2>
        <p className="text-base text-gray-600">
          Управляйте своими питомцами
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <HomeIcon className="w-8 h-8 opacity-80" />
            <div className="text-5xl font-bold">{loading ? '...' : pets.length}</div>
          </div>
          <div className="text-base font-medium opacity-90">Моих питомцев</div>
        </div>
        
        <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <HeartIcon className="w-8 h-8 opacity-80" />
            <div className="text-5xl font-bold">0</div>
          </div>
          <div className="text-base font-medium opacity-90">Визитов к ветеринару</div>
        </div>
        
        <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <ChartBarIcon className="w-8 h-8 opacity-80" />
            <div className="text-5xl font-bold">0</div>
          </div>
          <div className="text-base font-medium opacity-90">Событий за месяц</div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <h3 className="text-2xl font-bold text-gray-900 mb-4">О кабинете</h3>
          <div className="space-y-4 text-gray-600 text-base leading-relaxed">
            <p>
              Кабинет владельца - это ваш личный инструмент для управления питомцами.
              Здесь вы можете отслеживать здоровье, вести историю визитов к ветеринару
              и хранить важную информацию о ваших любимцах.
            </p>
            <p>
              Все данные синхронизируются с основной платформой и доступны в любое время.
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
                <div className="font-semibold text-gray-900 mb-1 text-base">Учет питомцев</div>
                <div className="text-sm text-gray-600">Полная информация о каждом любимце</div>
              </div>
            </div>
            
            <div className="flex items-start space-x-4">
              <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <span className="text-2xl">💚</span>
              </div>
              <div>
                <div className="font-semibold text-gray-900 mb-1 text-base">Медицинская карта</div>
                <div className="text-sm text-gray-600">История визитов и вакцинаций</div>
              </div>
            </div>
            
            <div className="flex items-start space-x-4">
              <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <span className="text-2xl">📊</span>
              </div>
              <div>
                <div className="font-semibold text-gray-900 mb-1 text-base">Статистика</div>
                <div className="text-sm text-gray-600">Отслеживание событий и активности</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
