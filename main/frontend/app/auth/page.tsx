'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '../../contexts/AuthContext';
import AuthForm from '../components/AuthForm';

export default function AuthPage() {
  const { login, register } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();

  const handleSubmit = async (data: { name?: string; email: string; password: string; mode: 'login' | 'register' }) => {
    const result = data.mode === 'login' 
      ? await login(data.email, data.password)
      : await register(data.name!, data.email, data.password);

    if (result.success) {
      // Проверяем есть ли параметр redirect
      const redirectUrl = searchParams.get('redirect');
      if (redirectUrl) {
        // Редиректим на внешний URL (Owner, Volunteer и т.д.)
        window.location.href = redirectUrl;
      } else {
        // Редиректим на главную страницу Main
        router.push('/');
      }
    }

    return result;
  };

  return (
    <AuthForm
      showTabs={true}
      onSubmit={handleSubmit}
      logoText="Зооплатформа"
      logoAlt="ЗооПлатформа"
    />
  );
}
