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
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        credentials: 'include',
      });

      if (meResponse.ok) {
        const meResult = await meResponse.json();
        const userId = meResult.data?.user?.id || meResult.data?.id;

        if (!userId) {
          console.error('User ID not found in response:', meResult);
          return;
        }

        const petsResponse = await fetch(`http://localhost:8100/api/pets/user/${userId}`, {
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

      {/* Красивые метки для демонстрации */}
      <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
        <h3 className="text-2xl font-bold text-gray-900 mb-6">Статусы питомцев</h3>
        <div className="space-y-6">
          {/* Статусы */}
          <div>
            <h4 className="text-sm font-semibold text-gray-700 mb-3">Основные статусы</h4>
            <div className="flex flex-wrap gap-3">
              <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-green-100 text-green-800 border border-green-200">
                <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                Дома
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-blue-100 text-blue-800 border border-blue-200">
                <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
                Здоров
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-purple-100 text-purple-800 border border-purple-200">
                <span className="w-2 h-2 bg-purple-500 rounded-full mr-2"></span>
                Вакцинирован
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-pink-100 text-pink-800 border border-pink-200">
                <span className="w-2 h-2 bg-pink-500 rounded-full mr-2"></span>
                Стерилизован
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800 border border-yellow-200">
                <span className="w-2 h-2 bg-yellow-500 rounded-full mr-2"></span>
                Требует внимания
              </span>
            </div>
          </div>

          {/* Типы животных */}
          <div>
            <h4 className="text-sm font-semibold text-gray-700 mb-3">Виды животных</h4>
            <div className="flex flex-wrap gap-3">
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium bg-gradient-to-r from-orange-100 to-orange-50 text-orange-800 border border-orange-200 shadow-sm">
                🐕 Собака
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium bg-gradient-to-r from-gray-100 to-gray-50 text-gray-800 border border-gray-200 shadow-sm">
                🐈 Кошка
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium bg-gradient-to-r from-sky-100 to-sky-50 text-sky-800 border border-sky-200 shadow-sm">
                🐦 Птица
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium bg-gradient-to-r from-emerald-100 to-emerald-50 text-emerald-800 border border-emerald-200 shadow-sm">
                🐰 Грызун
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium bg-gradient-to-r from-teal-100 to-teal-50 text-teal-800 border border-teal-200 shadow-sm">
                🐢 Рептилия
              </span>
            </div>
          </div>

          {/* Медицинские метки */}
          <div>
            <h4 className="text-sm font-semibold text-gray-700 mb-3">Медицинские показатели</h4>
            <div className="flex flex-wrap gap-3">
              <span className="inline-flex items-center px-3 py-1.5 rounded-md text-xs font-semibold bg-red-50 text-red-700 border-2 border-red-200">
                ⚠️ СРОЧНО
              </span>
              <span className="inline-flex items-center px-3 py-1.5 rounded-md text-xs font-semibold bg-green-50 text-green-700 border-2 border-green-200">
                ✓ Вакцинация актуальна
              </span>
              <span className="inline-flex items-center px-3 py-1.5 rounded-md text-xs font-semibold bg-blue-50 text-blue-700 border-2 border-blue-200">
                💉 Прививка через 2 дня
              </span>
              <span className="inline-flex items-center px-3 py-1.5 rounded-md text-xs font-semibold bg-purple-50 text-purple-700 border-2 border-purple-200">
                🏥 На лечении
              </span>
              <span className="inline-flex items-center px-3 py-1.5 rounded-md text-xs font-semibold bg-orange-50 text-orange-700 border-2 border-orange-200">
                📋 Плановый осмотр
              </span>
            </div>
          </div>

          {/* Размеры */}
          <div>
            <h4 className="text-sm font-semibold text-gray-700 mb-3">Размеры</h4>
            <div className="flex flex-wrap gap-3">
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
                Миниатюрный
              </span>
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                Маленький
              </span>
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-cyan-100 text-cyan-800">
                Средний
              </span>
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-teal-100 text-teal-800">
                Крупный
              </span>
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800">
                Гигантский
              </span>
            </div>
          </div>

          {/* Возрастные группы */}
          <div>
            <h4 className="text-sm font-semibold text-gray-700 mb-3">Возраст</h4>
            <div className="flex flex-wrap gap-3">
              <span className="inline-flex items-center px-4 py-2 rounded-xl text-sm font-medium bg-pink-50 text-pink-700 border border-pink-200 shadow-sm hover:shadow-md transition-shadow">
                🍼 Щенок/Котёнок
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-xl text-sm font-medium bg-blue-50 text-blue-700 border border-blue-200 shadow-sm hover:shadow-md transition-shadow">
                🎾 Молодой
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-xl text-sm font-medium bg-green-50 text-green-700 border border-green-200 shadow-sm hover:shadow-md transition-shadow">
                ⭐ Взрослый
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-xl text-sm font-medium bg-amber-50 text-amber-700 border border-amber-200 shadow-sm hover:shadow-md transition-shadow">
                👴 Пожилой
              </span>
            </div>
          </div>

          {/* Специальные метки */}
          <div>
            <h4 className="text-sm font-semibold text-gray-700 mb-3">Особые отметки</h4>
            <div className="flex flex-wrap gap-3">
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-bold bg-gradient-to-r from-yellow-400 to-yellow-500 text-white shadow-md">
                ⭐ VIP
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-bold bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-md">
                🏆 Чемпион
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-bold bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md">
                🎖️ Родословная
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-bold bg-gradient-to-r from-red-500 to-orange-500 text-white shadow-md">
                🔥 Популярный
              </span>
              <span className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-bold bg-gradient-to-r from-green-500 to-emerald-500 text-white shadow-md">
                ✨ Новый
              </span>
            </div>
          </div>
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
