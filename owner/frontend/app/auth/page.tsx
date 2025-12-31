'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export default function AuthPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      // Проверяем, авторизован ли уже пользователь
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      // Проверяем статус ответа
      if (!meResponse.ok) {
        console.error('Auth check failed with status:', meResponse.status);
        // Не авторизован - редирект на SSO
        const from = searchParams.get('from');
        if (from !== 'main') {
          window.location.href = 'http://localhost:3000/auth?redirect=http://localhost:6100/auth?from=main';
        } else {
          setChecking(false);
        }
        return;
      }

      const meResult = await meResponse.json();

      if (meResult.success) {
        // Уже авторизован - переходим в дашборд
        router.push('/pets');
        return;
      }

      // Не авторизован - редирект на SSO только если не пришли оттуда
      const from = searchParams.get('from');
      if (from !== 'main') {
        window.location.href = 'http://localhost:3000/auth?redirect=http://localhost:6100/auth?from=main';
      } else {
        // Пришли с основного сайта, но всё ещё не авторизованы
        setChecking(false);
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      // При ошибке редиректим на SSO
      const from = searchParams.get('from');
      if (from !== 'main') {
        window.location.href = 'http://localhost:3000/auth?redirect=http://localhost:6100/auth?from=main';
      } else {
        setChecking(false);
      }
    }
  };

  if (!checking && searchParams.get('from') === 'main') {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
        <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Ошибка авторизации</h1>
          <p className="text-gray-600 mb-6">
            Не удалось войти в систему. Попробуйте ещё раз.
          </p>
          <button
            onClick={() => window.location.href = 'http://localhost:3000/auth?redirect=http://localhost:6100'}
            className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors font-medium"
          >
            Войти заново
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
      <div className="text-center">
        <div className="text-lg text-gray-600">
          {checking ? 'Проверка авторизации...' : 'Перенаправление...'}
        </div>
      </div>
    </div>
  );
}
