import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';

export function useUnreadMessages() {
  const { user } = useAuth();
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    if (!user) {
      setUnreadCount(0);
      return;
    }

    const fetchUnreadCount = async () => {
      try {
        const token = localStorage.getItem('auth_token');
        const headers: HeadersInit = {
          'Content-Type': 'application/json',
        };
        
        if (token) {
          headers['Authorization'] = `Bearer ${token}`;
        }
        
        const response = await fetch('/api/messages/unread', {
          credentials: 'include',
          headers,
        });

        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data) {
            setUnreadCount(data.data.count || 0);
          } else {
            setUnreadCount(0);
          }
        } else if (response.status === 401) {
          // Пользователь не авторизован - это нормально
          setUnreadCount(0);
        }
      } catch (error) {
        // Тихо игнорируем ошибки сети - не критично для UI
        // Ошибка может возникать когда пользователь не авторизован
        setUnreadCount(0);
      }
    };

    // Небольшая задержка перед первым запросом
    const initialTimeout = setTimeout(fetchUnreadCount, 100);

    // Обновляем каждые 30 секунд
    const interval = setInterval(fetchUnreadCount, 30000);

    return () => {
      clearTimeout(initialTimeout);
      clearInterval(interval);
    };
  }, [user]);

  return unreadCount;
}
