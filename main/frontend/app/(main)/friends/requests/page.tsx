'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../../../contexts/AuthContext';
import { friendsApi, Friendship } from '../../../../lib/api';
import { getMediaUrl, getFullName } from '../../../../lib/utils';
import { UserIcon, CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { useToast } from '@/contexts/ToastContext';

export default function FriendRequestsPage() {
  const [requests, setRequests] = useState<Friendship[]>([]);
  const [loading, setLoading] = useState(true);
  const [processingIds, setProcessingIds] = useState<Set<number>>(new Set());
  const { user: currentUser, isAuthenticated, isLoading } = useAuth();
  const { showToast } = useToast();
  const router = useRouter();

  useEffect(() => {
    if (typeof window === 'undefined') return;

    if (!isLoading && !isAuthenticated) {
      router.push('/auth');
      return;
    }

    if (currentUser) {
      loadRequests();
    }
  }, [isLoading, isAuthenticated, currentUser]);

  const loadRequests = async () => {
    setLoading(true);
    try {
      const response = await friendsApi.getRequests();
      if (response.success && response.data) {
        setRequests(response.data);
      }
    } catch (error) {
      console.error('Ошибка загрузки запросов:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAccept = async (friendshipId: number) => {
    setProcessingIds(prev => new Set(prev).add(friendshipId));
    try {
      const response = await friendsApi.acceptRequest(friendshipId);
      if (response.success) {
        showToast('success', 'Запрос принят');
        setRequests(prev => prev.filter(r => r.id !== friendshipId));
      } else {
        showToast('error', response.error || 'Ошибка принятия запроса');
      }
    } catch (error) {
      showToast('error', 'Ошибка принятия запроса');
    } finally {
      setProcessingIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(friendshipId);
        return newSet;
      });
    }
  };

  const handleReject = async (friendshipId: number) => {
    setProcessingIds(prev => new Set(prev).add(friendshipId));
    try {
      const response = await friendsApi.rejectRequest(friendshipId);
      if (response.success) {
        showToast('success', 'Запрос отклонен');
        setRequests(prev => prev.filter(r => r.id !== friendshipId));
      } else {
        showToast('error', response.error || 'Ошибка отклонения запроса');
      }
    } catch (error) {
      showToast('error', 'Ошибка отклонения запроса');
    } finally {
      setProcessingIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(friendshipId);
        return newSet;
      });
    }
  };

  const handleProfileClick = (friendId: number) => {
    router.push(`/id${friendId}`);
  };

  if (isLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 mb-4">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Запросы в друзья</h1>
        <p className="text-sm text-gray-600">
          {requests.length === 0 
            ? 'У вас нет новых запросов' 
            : `Новых запросов: ${requests.length}`}
        </p>
      </div>

      {/* Requests List */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
        {requests.length === 0 ? (
          <div className="text-center py-12">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
              <UserIcon className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-600 mb-2">Нет новых запросов</p>
            <p className="text-sm text-gray-500">
              Когда кто-то отправит вам запрос в друзья, он появится здесь
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {requests.map((request) => {
              const friend = request.friend;
              const isProcessing = processingIds.has(request.id);
              
              return (
                <div
                  key={request.id}
                  className="flex items-center gap-4 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors"
                >
                  {/* Avatar */}
                  <button
                    onClick={() => handleProfileClick(friend.id)}
                    className="flex-shrink-0"
                  >
                    <div className="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                      {friend.avatar ? (
                        <img 
                          src={getMediaUrl(friend.avatar) || ''} 
                          alt={getFullName(friend.name, friend.last_name)} 
                          className="w-full h-full object-cover" 
                        />
                      ) : (
                        <UserIcon className="w-8 h-8 text-gray-400" />
                      )}
                    </div>
                  </button>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <button
                      onClick={() => handleProfileClick(friend.id)}
                      className="text-left"
                    >
                      <p className="font-semibold text-gray-900 hover:text-blue-600 transition-colors">
                        {getFullName(friend.name, friend.last_name)}
                      </p>
                      {friend.location && (
                        <p className="text-sm text-gray-500">{friend.location}</p>
                      )}
                    </button>
                    <p className="text-xs text-gray-400 mt-1">
                      {new Date(request.created_at).toLocaleDateString('ru-RU', {
                        day: 'numeric',
                        month: 'long',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}
                    </p>
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-2 flex-shrink-0">
                    <button
                      onClick={() => handleAccept(request.id)}
                      disabled={isProcessing}
                      className="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <CheckIcon className="w-5 h-5" />
                      <span className="hidden sm:inline">Принять</span>
                    </button>
                    <button
                      onClick={() => handleReject(request.id)}
                      disabled={isProcessing}
                      className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <XMarkIcon className="w-5 h-5" />
                      <span className="hidden sm:inline">Отклонить</span>
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
