'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { BuildingOffice2Icon, ArrowRightIcon } from '@heroicons/react/24/outline';

interface Clinic {
  id: number;
  name: string;
  type: string;
  description?: string;
  address?: string;
  logo?: string;
}

export default function SelectClinic() {
  const router = useRouter();
  const [clinics, setClinics] = useState<Clinic[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<{ email: string; name?: string } | null>(null);
  const [hasChecked, setHasChecked] = useState(false);

  useEffect(() => {
    if (!hasChecked) {
      checkAuthAndLoadClinics();
    }
  }, [hasChecked]);

  const checkAuthAndLoadClinics = async () => {
    if (hasChecked) return;
    setHasChecked(true);
    
    try {
      // Проверяем авторизацию через Main API (SSO)
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      if (!meResponse.ok) {
        console.error('Auth check failed with status:', meResponse.status);
        router.push('/auth');
        return;
      }

      const meResult = await meResponse.json();

      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      setUser(meResult.data);

      // Загружаем список клиник пользователя
      const clinicsResponse = await fetch('http://localhost:8600/api/my-clinics', {
        method: 'GET',
        credentials: 'include',
      });

      if (!clinicsResponse.ok) {
        console.error('Failed to load clinics, status:', clinicsResponse.status);
        // Просто показываем пустой список, не редиректим
        setClinics([]);
        setLoading(false);
        return;
      }

      const clinicsResult = await clinicsResponse.json();

      if (clinicsResult.success && clinicsResult.data) {
        setClinics(clinicsResult.data);
        
        // Убрал автоматический редирект - пусть пользователь сам выберет
        // if (clinicsResult.data.length === 1) {
        //   selectClinic(clinicsResult.data[0].id);
        //   return;
        // }
      }
    } catch (error) {
      console.error('Failed to load clinics:', error);
      router.push('/auth');
      return;
    }

    setLoading(false);
  };

  const selectClinic = (clinicId: number) => {
    // Сохраняем выбранную клинику в localStorage
    localStorage.setItem('selectedClinicId', clinicId.toString());
    // Переходим в дашборд
    router.push('/dashboard');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
        <div className="text-center">
          <div className="text-lg text-gray-600">Загрузка клиник...</div>
        </div>
      </div>
    );
  }

  if (!user) {
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
            Выберите клинику
          </h1>
          <p className="text-lg text-gray-600">
            Вы являетесь сотрудником {clinics.length} {clinics.length === 1 ? 'клиники' : 'клиник'}
          </p>
        </div>

        {clinics.length === 0 ? (
          <div className="bg-white rounded-2xl p-12 shadow-sm border border-gray-100 text-center">
            <BuildingOffice2Icon className="w-16 h-16 mx-auto mb-4 text-gray-300" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              У вас пока нет клиники
            </h3>
            <p className="text-gray-600 mb-6">
              Создайте свою клинику, чтобы начать управлять пациентами, записями и медицинскими картами
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button
                onClick={() => router.push('/create')}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
              >
                Создать клинику
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
            {clinics.map((clinic) => (
              <button
                key={clinic.id}
                onClick={() => selectClinic(clinic.id)}
                className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-lg hover:border-blue-200 transition-all text-left group"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    {clinic.logo ? (
                      <img
                        src={clinic.logo}
                        alt={clinic.name}
                        className="w-16 h-16 rounded-xl object-cover"
                      />
                    ) : (
                      <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center">
                        <span className="text-3xl">🏥</span>
                      </div>
                    )}
                    <div>
                      <h3 className="text-xl font-bold text-gray-900 mb-1">
                        {clinic.name}
                      </h3>
                      {clinic.address && (
                        <p className="text-sm text-gray-500">{clinic.address}</p>
                      )}
                    </div>
                  </div>
                  <ArrowRightIcon className="w-6 h-6 text-gray-400 group-hover:text-blue-600 group-hover:translate-x-1 transition-all" />
                </div>
                
                {clinic.description && (
                  <p className="text-gray-600 text-sm line-clamp-2 mb-4">
                    {clinic.description}
                  </p>
                )}

                <div className="flex items-center justify-between pt-4 border-t border-gray-100">
                  <span className="text-sm text-gray-500">Тип: {clinic.type}</span>
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
