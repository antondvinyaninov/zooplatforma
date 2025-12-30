'use client';

import { useRouter } from 'next/navigation';
import AuthForm from '../components/AuthForm';

const ADMIN_API_URL = 'http://localhost:9000';

export default function ShelterAuth() {
  const router = useRouter();

  const handleSubmit = async (data: { email: string; password: string }) => {
    try {
      // Логинимся через главный backend (SSO)
      const loginResponse = await fetch('http://localhost:8000/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email: data.email, password: data.password }),
      });

      const loginResult = await loginResponse.json();

      if (!loginResult.success) {
        return { success: false, error: loginResult.error || 'Неверный email или пароль' };
      }

      // Проверяем права через Admin API
      const meResponse = await fetch(`${ADMIN_API_URL}/api/admin/auth/me`, {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        return { success: false, error: 'У вас нет прав доступа к кабинету приюта' };
      }

      // Проверяем роль (shelter_admin или выше)
      const allowedRoles = ['shelter_admin', 'moderator', 'admin', 'superadmin'];
      if (!allowedRoles.includes(meResult.data?.role)) {
        return { success: false, error: 'Доступ только для администраторов приюта' };
      }

      // Успешный вход
      router.push('/');
      return { success: true };
    } catch (err) {
      return { success: false, error: 'Ошибка подключения к серверу' };
    }
  };

  return (
    <AuthForm
      mode="login"
      showTabs={false}
      onSubmit={handleSubmit}
      logoText="Приют"
      logoAlt="Кабинет приюта"
      subtitle="Войдите в кабинет приюта"
      infoTitle="🏠 Кабинет приюта"
      infoText="Доступ для администраторов приюта и волонтеров"
    />
  );
}
