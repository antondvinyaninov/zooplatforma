'use client';

import {
  UserGroupIcon,
  CalendarIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';
import { useState, useEffect } from 'react';

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  photo?: string;
  status: string;
  created_at: string;
}

export default function OverviewPage() {
  const [patients, setPatients] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<{ email: string; name?: string } | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      // Получаем данные пользователя
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      if (meResponse.ok) {
        const meResult = await meResponse.json();
        if (meResult.success) {
          setUser(meResult.data);
        }
      }

      // Получаем выбранную клинику из localStorage
      const clinicId = localStorage.getItem('selectedClinicId');
      if (clinicId) {
        // Загружаем пациентов клиники
        const patientsResponse = await fetch('http://localhost:8600/api/my-patients', {
          method: 'GET',
          credentials: 'include',
          headers: {
            'X-Clinic-ID': clinicId,
          },
        });

        if (patientsResponse.ok) {
          const patientsResult = await patientsResponse.json();
          if (patientsResult.success && patientsResult.data) {
            setPatients(patientsResult.data);
          }
        }
      }
    } catch (error) {
      console.error('Failed to load data:', error);
    }

    setLoading(false);
  };

  if (loading) {
    return (
      <div className="text-center py-12 text-gray-400">Загрузка...</div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-4xl font-bold text-gray-900 mb-2">Добро пожаловать!</h2>
        <p className="text-base text-gray-600">
          {user?.name || user?.email}, управление ветеринарной клиникой
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <UserGroupIcon className="w-8 h-8 opacity-80" />
            <div className="text-5xl font-bold">{patients.length}</div>
          </div>
          <div className="text-base font-medium opacity-90">Пациентов</div>
        </div>
        
        <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <CalendarIcon className="w-8 h-8 opacity-80" />
            <div className="text-5xl font-bold">0</div>
          </div>
          <div className="text-base font-medium opacity-90">Записей на сегодня</div>
        </div>
        
        <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl p-6 text-white shadow-lg hover:shadow-xl transition-all hover:scale-105 cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <ChartBarIcon className="w-8 h-8 opacity-80" />
            <div className="text-5xl font-bold">0</div>
          </div>
          <div className="text-base font-medium opacity-90">Приёмов за месяц</div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <h3 className="text-2xl font-bold text-gray-900 mb-4">О кабинете</h3>
          <div className="space-y-4 text-gray-600 text-base leading-relaxed">
            <p>
              Кабинет ветеринарной клиники - это ваш инструмент для управления приёмами,
              ведения медицинских карт и координации работы с владельцами животных.
            </p>
            <p>
              Все данные синхронизируются с основной платформой и доступны владельцам питомцев.
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
                <div className="font-semibold text-gray-900 mb-1 text-base">Пациенты</div>
                <div className="text-sm text-gray-600">База данных всех пациентов</div>
              </div>
            </div>
            
            <div className="flex items-start space-x-4">
              <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <span className="text-2xl">📅</span>
              </div>
              <div>
                <div className="font-semibold text-gray-900 mb-1 text-base">Записи</div>
                <div className="text-sm text-gray-600">Управление записями на приём</div>
              </div>
            </div>
            
            <div className="flex items-start space-x-4">
              <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <span className="text-2xl">📋</span>
              </div>
              <div>
                <div className="font-semibold text-gray-900 mb-1 text-base">Медкарты</div>
                <div className="text-sm text-gray-600">Ведение медицинских карт</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
