'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { BuildingOffice2Icon, ArrowRightIcon } from '@heroicons/react/24/outline';

interface Shelter {
  id: number;
  name: string;
  type: string;
  description?: string;
  address?: string;
  logo?: string;
}

export default function SelectShelter() {
  const router = useRouter();
  const [shelters, setShelters] = useState<Shelter[]>([]);
  const [loading, setLoading] = useState(true);
  const [adminUser, setAdminUser] = useState<{ email: string; role: string } | null>(null);

  useEffect(() => {
    checkAuthAndLoadShelters();
  }, []);

  const checkAuthAndLoadShelters = async () => {
    try {
      // Проверяем авторизацию через Admin API (SSO)
      const meResponse = await fetch('http://localhost:9000/api/admin/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      // Проверяем роль (shelter_admin или выше)
      const allowedRoles = ['shelter_admin', 'moderator', 'admin', 'superadmin'];
      if (!allowedRoles.includes(meResult.data?.role)) {
        router.push('/auth');
        return;
      }

      setAdminUser(meResult.data);

      // Загружаем список приютов пользователя
      const sheltersResponse = await fetch('http://localhost:8200/api/my-shelters', {
        method: 'GET',
        credentials: 'include',
      });

      const sheltersResult = await sheltersResponse.json();

      if (sheltersResult.success && sheltersResult.data) {
        setShelters(sheltersResult.data);
        
        // Если только один приют - сразу переходим в него
        if (sheltersResult.data.length === 1) {
          selectShelter(sheltersResult.data[0].id);
          return;
        }
      }
    } catch (error) {
      console.error('Failed to load shelters:', error);
      router.push('/auth');
      return;
    }

    setLoading(false);
  };

  const selectShelter = (shelterId: number) => {
    // Сохраняем выбранный приют в localStorage
    localStorage.setItem('selectedShelterId', shelterId.toString());
    // Переходим в дашборд
    router.push('/dashboard');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
        <div className="text-center">
          <div className="text-lg text-gray-600">Загрузка приютов...</div>
        </div>
      </div>
    );
  }

  if (!adminUser) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 py-12 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-12">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-blue-600 rounded-2xl mb-6">
            <BuildingOffice2Icon className="w-10 h-10 text-white" />
          </div>
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Выберите приют
          </h1>
          <p className="text-lg text-gray-600">
            Вы являетесь администратором {shelters.length} {shelters.length === 1 ? 'приюта' : 'приютов'}
          </p>
        </div>

        {shelters.length === 0 ? (
          <div className="bg-white rounded-2xl p-12 shadow-sm border border-gray-100 text-center">
            <BuildingOffice2Icon className="w-16 h-16 mx-auto mb-4 text-gray-300" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              У вас пока нет приюта
            </h3>
            <p className="text-gray-600 mb-6">
              Создайте свой приют, чтобы начать управлять животными, волонтерами и заявками на пристройство
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button
                onClick={() => router.push('/create')}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
              >
                Создать приют
              </button>
              <button
                onClick={() => window.open('http://localhost:3000', '_blank')}
                className="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
              >
                Вернуться на главную
              </button>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {shelters.map((shelter) => (
              <button
                key={shelter.id}
                onClick={() => selectShelter(shelter.id)}
                className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-lg hover:border-blue-200 transition-all text-left group"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    {shelter.logo ? (
                      <img
                        src={shelter.logo}
                        alt={shelter.name}
                        className="w-16 h-16 rounded-xl object-cover"
                      />
                    ) : (
                      <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center">
                        <span className="text-3xl">🏠</span>
                      </div>
                    )}
                    <div>
                      <h3 className="text-xl font-bold text-gray-900 mb-1">
                        {shelter.name}
                      </h3>
                      {shelter.address && (
                        <p className="text-sm text-gray-500">{shelter.address}</p>
                      )}
                    </div>
                  </div>
                  <ArrowRightIcon className="w-6 h-6 text-gray-400 group-hover:text-blue-600 group-hover:translate-x-1 transition-all" />
                </div>
                
                {shelter.description && (
                  <p className="text-gray-600 text-sm line-clamp-2 mb-4">
                    {shelter.description}
                  </p>
                )}

                <div className="flex items-center justify-between pt-4 border-t border-gray-100">
                  <span className="text-sm text-gray-500">Тип: {shelter.type}</span>
                  <span className="text-sm font-medium text-blue-600 group-hover:text-blue-700">
                    Открыть →
                  </span>
                </div>
              </button>
            ))}
          </div>
        )}

        <div className="mt-8 text-center">
          <button
            onClick={() => {
              document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; domain=localhost';
              document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
              router.push('/auth');
            }}
            className="text-gray-600 hover:text-gray-900 transition-colors"
          >
            Выйти из аккаунта
          </button>
        </div>
      </div>
    </div>
  );
}
