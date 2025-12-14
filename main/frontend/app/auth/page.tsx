'use client';

import { useRouter } from 'next/navigation';
import { useAuth } from '../../contexts/AuthContext';
import AuthForm from '../components/AuthForm';

export default function AuthPage() {
  const { login, register } = useAuth();
  const router = useRouter();

  const handleSubmit = async (data: { name?: string; email: string; password: string; mode: 'login' | 'register' }) => {
    const result = data.mode === 'login' 
      ? await login(data.email, data.password)
      : await register(data.name!, data.email, data.password);

    if (result.success) {
      router.push('/');
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
