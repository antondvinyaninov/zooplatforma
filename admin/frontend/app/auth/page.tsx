'use client';

import { useRouter } from 'next/navigation';
import AuthForm from '../components/AuthForm';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9000';

export default function AdminAuth() {
  const router = useRouter();

  const handleSubmit = async (data: { email: string; password: string }) => {
    try {
      // Сначала логинимся через главный backend
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

      // Проверяем права администратора
      const meResponse = await fetch(`${API_URL}/api/admin/auth/me`, {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        return { success: false, error: 'У вас нет прав администратора' };
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
      logoText="ЗооАдминка"
      logoAlt="ЗооАдминка"
      subtitle="Войдите в панель администратора"
      infoTitle="🔒 Доступ ограничен"
      infoText="Доступ только для администраторов платформы"
    />
  );
}
