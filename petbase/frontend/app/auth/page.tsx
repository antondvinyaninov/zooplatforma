'use client';

import { useRouter } from 'next/navigation';
import AuthForm from '../components/AuthForm';

const API_URL = 'http://localhost:9000';

export default function PetBaseAuth() {
  const router = useRouter();

  const handleSubmit = async (data: { email: string; password: string }) => {
    try {
      // Логинимся через главный backend
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

      // Проверяем права суперадмина
      const meResponse = await fetch(`${API_URL}/api/admin/auth/me`, {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success || meResult.data?.role !== 'superadmin') {
        return { success: false, error: 'Доступ только для суперадминистраторов' };
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
      logoText="ЗооБаза"
      logoAlt="ЗооБаза"
      subtitle="Войдите в справочник животных"
      infoTitle="🔒 Доступ ограничен"
      infoText="Доступ только для суперадминистраторов платформы"
    />
  );
}
