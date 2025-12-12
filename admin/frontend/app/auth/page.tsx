'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import styles from './auth.module.css';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9000';

export default function AdminAuth() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // Сначала логинимся через главный backend
      const loginResponse = await fetch('http://localhost:8000/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email, password }),
      });

      const loginResult = await loginResponse.json();

      if (!loginResult.success) {
        setError(loginResult.error || 'Неверный email или пароль');
        setLoading(false);
        return;
      }

      // Проверяем права администратора
      const meResponse = await fetch(`${API_URL}/api/admin/auth/me`, {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        setError('У вас нет прав администратора');
        setLoading(false);
        return;
      }

      // Успешный вход
      router.push('/');
    } catch (err) {
      setError('Ошибка подключения к серверу');
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <h1 className={styles.title}>ЗооПлатформа</h1>
        <p className={styles.subtitle}>Панель администратора</p>

        <form onSubmit={handleSubmit} className={styles.form}>
          {error && <div className={styles.error}>{error}</div>}

          <div className={styles.field}>
            <label className={styles.label}>Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={styles.input}
              placeholder="admin@example.com"
              required
              disabled={loading}
            />
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Пароль</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={styles.input}
              placeholder="••••••••"
              required
              disabled={loading}
            />
          </div>

          <button type="submit" className={styles.button} disabled={loading}>
            {loading ? 'Вход...' : 'Войти'}
          </button>
        </form>

        <div className={styles.footer}>
          <p className={styles.footerText}>
            Доступ только для администраторов платформы
          </p>
        </div>
      </div>
    </div>
  );
}
