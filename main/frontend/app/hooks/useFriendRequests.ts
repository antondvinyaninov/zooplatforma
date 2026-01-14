'use client';

import { useState, useEffect } from 'react';
import { friendsApi, Friendship } from '@/lib/api';

export function useFriendRequests() {
  const [requests, setRequests] = useState<Friendship[]>([]);
  const [count, setCount] = useState(0);
  const [loading, setLoading] = useState(true);

  const loadRequests = async () => {
    setLoading(true);
    const result = await friendsApi.getRequests();
    if (result.success && result.data) {
      setRequests(result.data);
      setCount(result.data.length);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadRequests();
    
    // Обновляем каждые 30 секунд
    const interval = setInterval(loadRequests, 30000);
    
    return () => clearInterval(interval);
  }, []);

  return { requests, count, loading, refresh: loadRequests };
}
