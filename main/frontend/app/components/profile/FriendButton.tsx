'use client';

import { useState, useEffect } from 'react';
import { UserPlusIcon, UserMinusIcon, CheckIcon, XMarkIcon, ClockIcon } from '@heroicons/react/24/outline';
import { friendsApi, FriendshipStatus } from '@/lib/api';
import { useToast } from '@/contexts/ToastContext';

interface FriendButtonProps {
  userId: number;
  currentUserId: number;
}

export default function FriendButton({ userId, currentUserId }: FriendButtonProps) {
  const [status, setStatus] = useState<FriendshipStatus>({ status: 'none' });
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  useEffect(() => {
    loadFriendshipStatus();
  }, [userId]);

  const loadFriendshipStatus = async () => {
    try {
      const response = await friendsApi.getStatus(userId);
      if (response.success && response.data) {
        setStatus(response.data);
      }
    } catch (error) {
      console.error('Ошибка загрузки статуса дружбы:', error);
    }
  };

  const handleSendRequest = async () => {
    setLoading(true);
    try {
      const response = await friendsApi.sendRequest(userId);
      if (response.success) {
        showToast('success', 'Запрос в друзья отправлен');
        await loadFriendshipStatus();
      } else {
        showToast('error', response.error || 'Ошибка отправки запроса');
      }
    } catch (error) {
      showToast('error', 'Ошибка отправки запроса');
    } finally {
      setLoading(false);
    }
  };

  const handleAcceptRequest = async () => {
    if (!status.id) return;
    setLoading(true);
    try {
      const response = await friendsApi.acceptRequest(status.id);
      if (response.success) {
        showToast('success', 'Запрос принят');
        await loadFriendshipStatus();
      } else {
        showToast('error', response.error || 'Ошибка принятия запроса');
      }
    } catch (error) {
      showToast('error', 'Ошибка принятия запроса');
    } finally {
      setLoading(false);
    }
  };

  const handleRejectRequest = async () => {
    if (!status.id) return;
    setLoading(true);
    try {
      const response = await friendsApi.rejectRequest(status.id);
      if (response.success) {
        showToast('success', 'Запрос отклонен');
        await loadFriendshipStatus();
      } else {
        showToast('error', response.error || 'Ошибка отклонения запроса');
      }
    } catch (error) {
      showToast('error', 'Ошибка отклонения запроса');
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveFriend = async () => {
    if (!status.id) return;
    if (!confirm('Вы уверены, что хотите удалить из друзей?')) return;
    
    setLoading(true);
    try {
      const response = await friendsApi.removeFriend(status.id);
      if (response.success) {
        showToast('success', 'Удалено из друзей');
        await loadFriendshipStatus();
      } else {
        showToast('error', response.error || 'Ошибка удаления из друзей');
      }
    } catch (error) {
      showToast('error', 'Ошибка удаления из друзей');
    } finally {
      setLoading(false);
    }
  };

  // Не показываем кнопку для своего профиля
  if (userId === currentUserId) {
    return null;
  }

  // Друзья
  if (status.status === 'accepted') {
    return (
      <button
        onClick={handleRemoveFriend}
        disabled={loading}
        className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
      >
        <CheckIcon className="w-5 h-5" />
        <span>В друзьях</span>
      </button>
    );
  }

  // Исходящий запрос (ожидание)
  if (status.status === 'pending' && status.is_outgoing) {
    return (
      <button
        disabled
        className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-500 rounded-lg cursor-not-allowed"
      >
        <ClockIcon className="w-5 h-5" />
        <span>Запрос отправлен</span>
      </button>
    );
  }

  // Входящий запрос (нужно принять/отклонить)
  if (status.status === 'pending' && !status.is_outgoing) {
    return (
      <div className="flex items-center gap-2">
        <button
          onClick={handleAcceptRequest}
          disabled={loading}
          className="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
        >
          <CheckIcon className="w-5 h-5" />
          <span>Принять</span>
        </button>
        <button
          onClick={handleRejectRequest}
          disabled={loading}
          className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
        >
          <XMarkIcon className="w-5 h-5" />
          <span>Отклонить</span>
        </button>
      </div>
    );
  }

  // Нет дружбы - показываем кнопку добавления
  return (
    <button
      onClick={handleSendRequest}
      disabled={loading}
      className="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
    >
      <UserPlusIcon className="w-5 h-5" />
      <span>Добавить в друзья</span>
    </button>
  );
}
