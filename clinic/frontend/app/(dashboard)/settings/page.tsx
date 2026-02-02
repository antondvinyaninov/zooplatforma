'use client';

import { useState, useEffect } from 'react';
import {
  BuildingOffice2Icon,
  MapPinIcon,
  PhoneIcon,
  ClockIcon,
  UserGroupIcon,
} from '@heroicons/react/24/outline';
import MembersManager from '../../components/MembersManager';

interface Clinic {
  id: number;
  name: string;
  type: string;
  description?: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  working_hours?: string;
  logo?: string;
}

type SettingsTab = 'info' | 'roles';

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<SettingsTab>('info');
  const [clinic, setClinic] = useState<Clinic | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  // Helper функция для получения токена из cookie
  const getAuthToken = () => {
    return document.cookie
      .split('; ')
      .find(row => row.startsWith('auth_token='))
      ?.split('=')[1];
  };

  // Helper функция для создания headers с авторизацией
  const getAuthHeaders = (clinicId: string) => {
    const token = getAuthToken();
    return {
      'Content-Type': 'application/json',
      'X-Clinic-ID': clinicId,
      ...(token && { 'Authorization': `Bearer ${token}` }),
    };
  };

  useEffect(() => {
    loadClinic();
  }, []);

  const loadClinic = async () => {
    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) {
        console.error('No clinic selected');
        setLoading(false);
        return;
      }

      // Получаем токен: сначала из cookie, если нет - через Main API
      let token = getAuthToken();
      
      if (!token) {
        // Получаем токен через Main API
        try {
          const meResponse = await fetch('http://localhost:7100/api/auth/me', {
            method: 'GET',
            credentials: 'include',
          });
          
          if (meResponse.ok) {
            const meData = await meResponse.json();
            if (meData.success && meData.data && meData.data.token) {
              token = meData.data.token;
            }
          }
        } catch (error) {
          console.error('Failed to get token from Main API:', error);
        }
      }

      if (!token) {
        setLoading(false);
        return;
      }

      const response = await fetch(`http://localhost:8600/api/organization`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'X-Clinic-ID': clinicId,
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const result = await response.json();
        
        if (result.success && result.data) {
          setClinic(result.data);
        }
      }
    } catch (error) {
      console.error('Failed to load clinic:', error);
    }

    setLoading(false);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!clinic) return;

    setSaving(true);
    setMessage(null);

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) return;

      // Получаем токен
      let token = getAuthToken();
      if (!token) {
        const meResponse = await fetch('http://localhost:7100/api/auth/me', {
          method: 'GET',
          credentials: 'include',
        });
        if (meResponse.ok) {
          const meData = await meResponse.json();
          if (meData.success && meData.data && meData.data.token) {
            token = meData.data.token;
          }
        }
      }

      if (!token) {
        setMessage({ type: 'error', text: 'Токен авторизации не найден' });
        setSaving(false);
        return;
      }

      const response = await fetch(`http://localhost:8600/api/organization`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'X-Clinic-ID': clinicId,
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(clinic),
      });

      if (response.ok) {
        setMessage({ type: 'success', text: 'Настройки успешно сохранены' });
      } else {
        setMessage({ type: 'error', text: 'Ошибка при сохранении настроек' });
      }
    } catch (error) {
      console.error('Failed to save clinic:', error);
      setMessage({ type: 'error', text: 'Ошибка при сохранении настроек' });
    }

    setSaving(false);
  };

  const handleChange = (field: keyof Clinic, value: string) => {
    if (!clinic) return;
    setClinic({ ...clinic, [field]: value });
  };

  if (loading) {
    return (
      <div className="text-center py-12 text-gray-400">Загрузка...</div>
    );
  }

  if (!clinic) {
    return (
      <div className="text-center py-12 text-gray-400">Клиника не найдена</div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Настройки клиники</h2>
        <p className="text-gray-600">Управление информацией и командой ветеринарной клиники</p>
      </div>

      {/* Вкладки */}
      <div className="border-b border-gray-200">
        <nav className="flex gap-8">
          <button
            onClick={() => setActiveTab('info')}
            className={`pb-4 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === 'info'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <BuildingOffice2Icon className="w-5 h-5" />
              О клинике
            </div>
          </button>
          <button
            onClick={() => setActiveTab('roles')}
            className={`pb-4 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === 'roles'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <UserGroupIcon className="w-5 h-5" />
              Роли и пользователи
            </div>
          </button>
        </nav>
      </div>

      {message && (
        <div
          className={`p-4 rounded-lg ${
            message.type === 'success'
              ? 'bg-green-50 text-green-800 border border-green-200'
              : 'bg-red-50 text-red-800 border border-red-200'
          }`}
        >
          {message.text}
        </div>
      )}

      {/* Вкладка "О клинике" */}
      {activeTab === 'info' && (
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Основная информация */}
          <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
            <div className="flex items-center gap-3 mb-6">
              <BuildingOffice2Icon className="w-6 h-6 text-blue-600" />
              <h3 className="text-xl font-bold text-gray-900">Основная информация</h3>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Название клиники *
                </label>
                <input
                  type="text"
                  value={clinic.name}
                  onChange={(e) => handleChange('name', e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Тип организации
                </label>
                <input
                  type="text"
                  value={clinic.type}
                  disabled
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Описание
                </label>
                <textarea
                  value={clinic.description || ''}
                  onChange={(e) => handleChange('description', e.target.value)}
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Краткое описание клиники, услуги, специализация..."
                />
              </div>
            </div>
          </div>

          {/* Контактная информация */}
          <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
            <div className="flex items-center gap-3 mb-6">
              <PhoneIcon className="w-6 h-6 text-green-600" />
              <h3 className="text-xl font-bold text-gray-900">Контактная информация</h3>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Адрес
                </label>
                <div className="relative">
                  <MapPinIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type="text"
                    value={clinic.address || ''}
                    onChange={(e) => handleChange('address', e.target.value)}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="г. Москва, ул. Примерная, д. 1"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Телефон
                </label>
                <div className="relative">
                  <PhoneIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type="tel"
                    value={clinic.phone || ''}
                    onChange={(e) => handleChange('phone', e.target.value)}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="+7 (999) 123-45-67"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <input
                  type="email"
                  value={clinic.email || ''}
                  onChange={(e) => handleChange('email', e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="clinic@example.com"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Веб-сайт
                </label>
                <input
                  type="url"
                  value={clinic.website || ''}
                  onChange={(e) => handleChange('website', e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="https://example.com"
                />
              </div>
            </div>
          </div>

          {/* Режим работы */}
          <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
            <div className="flex items-center gap-3 mb-6">
              <ClockIcon className="w-6 h-6 text-purple-600" />
              <h3 className="text-xl font-bold text-gray-900">Режим работы</h3>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Часы работы
              </label>
              <textarea
                value={clinic.working_hours || ''}
                onChange={(e) => handleChange('working_hours', e.target.value)}
                rows={3}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Пн-Пт: 9:00-20:00&#10;Сб-Вс: 10:00-18:00"
              />
              <p className="text-sm text-gray-500 mt-2">
                Укажите режим работы клиники (каждый день с новой строки)
              </p>
            </div>
          </div>

          {/* Кнопки */}
          <div className="flex gap-4">
            <button
              type="submit"
              disabled={saving}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving ? 'Сохранение...' : 'Сохранить изменения'}
            </button>
            <button
              type="button"
              onClick={loadClinic}
              className="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
            >
              Отменить
            </button>
          </div>
        </form>
      )}

      {/* Вкладка "Роли и пользователи" */}
      {activeTab === 'roles' && (
        <MembersManager />
      )}
    </div>
  );
}
