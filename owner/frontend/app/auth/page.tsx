'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function AuthPage() {
  const router = useRouter();
  const [checking, setChecking] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      // Проверяем, авторизован ли уже пользователь
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      if (meResponse.ok) {
        const meResult = await meResponse.json();
        if (meResult.success) {
          // Уже авторизован - переходим в дашборд
          router.push('/pets');
          return;
        }
      }

      // Не авторизован - показываем сообщение
      setError('Вы не авторизованы. Пожалуйста, войдите через главный сайт.');
      setChecking(false);
    } catch (error) {
      console.error('Auth check failed:', error);
      setError('Ошибка проверки авторизации.');
      setChecking(false);
    }
  };

  if (checking) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
        <div className="text-center">
          <div className="text-lg text-gray-600">Проверка авторизации...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Требуется авторизация</h1>
        <p className="text-gray-600 mb-6">
          {error || 'Для доступа к кабинету владельца необходимо войти в систему.'}
        </p>
        <a
          href="http://localhost:3000/auth"
          className="block w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors font-medium"
        >
          Войти через главный сайт
        </a>
      </div>
    </div>
  );
}
