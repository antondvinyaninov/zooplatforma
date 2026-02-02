'use client';

import { useRouter } from 'next/navigation';
import AuthForm from '../components/AuthForm';

export default function ShelterAuth() {
  const router = useRouter();

  const handleSubmit = async (data: { email: string; password: string }) => {
    try {
      // Логинимся через главный backend (SSO)
      const loginResponse = await fetch('http://localhost:7100/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email: data.email, password: data.password }),
      });

      const loginResult = await loginResponse.json();

      if (!loginResult.success) {
        return { success: false, error: loginResult.error || 'Неверный email или пароль' };
      }

      // Проверяем авторизацию через Main API (SSO)
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        return { success: false, error: 'Ошибка авторизации' };
      }

      // Успешный вход - переходим на главную страницу
      // Там будет проверка наличия приютов
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
