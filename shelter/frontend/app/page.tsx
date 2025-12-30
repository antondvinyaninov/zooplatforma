'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Home() {
  const router = useRouter();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    checkShelters();
  }, []);

  const checkShelters = async () => {
    try {
      // Проверяем авторизацию
      const meResponse = await fetch('http://localhost:9000/api/admin/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      // Проверяем роль
      const allowedRoles = ['shelter_admin', 'moderator', 'admin', 'superadmin'];
      if (!allowedRoles.includes(meResult.data?.role)) {
        router.push('/auth');
        return;
      }

      // Загружаем список приютов
      const sheltersResponse = await fetch('http://localhost:8200/api/my-shelters', {
        method: 'GET',
        credentials: 'include',
      });

      const sheltersResult = await sheltersResponse.json();

      if (!sheltersResult.success || !sheltersResult.data) {
        router.push('/select');
        return;
      }

      const shelters = sheltersResult.data;

      // Если приютов нет - на страницу выбора (там покажется сообщение)
      if (shelters.length === 0) {
        router.push('/select');
        return;
      }

      // Всегда показываем страницу выбора, если приютов больше одного
      if (shelters.length > 1) {
        router.push('/select');
        return;
      }

      // Если приют один - сразу в дашборд
      if (shelters.length === 1) {
        localStorage.setItem('selectedShelterId', shelters[0].id.toString());
        router.push('/dashboard');
        return;
      }
    } catch (error) {
      console.error('Failed to check shelters:', error);
      router.push('/auth');
    } finally {
      setChecking(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
      <div className="text-center">
        <div className="text-lg text-gray-600">
          {checking ? 'Проверка доступа...' : 'Перенаправление...'}
        </div>
      </div>
    </div>
  );
}
