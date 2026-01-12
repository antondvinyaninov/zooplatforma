'use client';

import { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { UserPlusIcon, CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { friendsApi, Friendship } from '@/lib/api';
import { getMediaUrl, getFullName } from '@/lib/utils';
import { useToast } from '@/contexts/ToastContext';

export default function FriendRequestsDropdown() {
  const [requests, setRequests] = useState<Friendship[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [processingIds, setProcessingIds] = useState<Set<number>>(new Set());
  const dropdownRef = useRef<HTMLDivElement>(null);
  const { showToast } = useToast();
  const router = useRouter();

  useEffect(() => {
    loadRequests();
    
    // Обновляем каждые 30 секунд
    const interval = setInterval(loadRequests, 30000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

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

  const handleAccept = async (friendshipId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    setProcessingIds(prev => new Set(prev).add(friendshipId));
    try {
      const response = await friendsApi.acceptRequest(friendshipId);
      if (response.success) {
        showToast('success', 'Запрос принят');
        setRequests(prev => prev.filter(r => r.id !== friendshipId));
      } else {
        showToast('error', response.error || 'Ошибка');
      }
    } catch (error) {
      showToast('error', 'Ошибка');
    } finally {
      setProcessingIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(friendshipId);
        return newSet;
      });
    }
  };

  const handleReject = async (friendshipId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    setProcessingIds(prev => new Set(prev).add(friendshipId));
    try {
      const response = await friendsApi.rejectRequest(friendshipId);
      if (response.success) {
        showToast('success', 'Запрос отклонен');
        setRequests(prev => prev.filter(r => r.id !== friendshipId));
      } else {
        showToast('error', response.error || 'Ошибка');
      }
    } catch (error) {
      showToast('error', 'Ошибка');
    } finally {
      setProcessingIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(friendshipId);
        return newSet;
      });
    }
  };

  const handleViewAll = () => {
    setIsOpen(false);
    router.push('/friends/requests');
  };

  const handleProfileClick = (friendId: number) => {
    setIsOpen(false);
    router.push(`/id${friendId}`);
  };

  return (
    <div className="relative" ref={dropdownRef}>
      {/* Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="relative p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
      >
        <UserPlusIcon className="w-6 h-6" strokeWidth={2} />
        {requests.length > 0 && (
          <span 
            className="absolute top-0 right-0 text-white text-[10px] rounded-full w-4 h-4 flex items-center justify-center font-bold" 
            style={{ backgroundColor: '#FC2B2B' }}
          >
            {requests.length > 9 ? '9+' : requests.length}
          </span>
        )}
      </button>

      {/* Dropdown */}
      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg border border-gray-200 overflow-hidden z-50">
          {/* Header */}
          <div className="px-4 py-3 border-b border-gray-200">
            <h3 className="font-semibold text-gray-900">Запросы в друзья</h3>
          </div>

          {/* Content */}
          <div className="max-h-96 overflow-y-auto">
            {loading && requests.length === 0 ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
              </div>
            ) : requests.length === 0 ? (
              <div className="text-center py-8 px-4">
                <UserPlusIcon className="w-12 h-12 text-gray-300 mx-auto mb-2" />
                <p className="text-sm text-gray-500">Нет новых запросов</p>
              </div>
            ) : (
              <div className="divide-y divide-gray-100">
                {requests.map((request) => {
                  const friend = request.friend;
                  const isProcessing = processingIds.has(request.id);

                  return (
                    <div
                      key={request.id}
                      className="p-3 hover:bg-gray-50 transition-colors"
                    >
                      <div className="flex items-start gap-3">
                        {/* Avatar */}
                        <button
                          onClick={() => handleProfileClick(friend.id)}
                          className="flex-shrink-0"
                        >
                          <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                            {friend.avatar ? (
                              <img 
                                src={getMediaUrl(friend.avatar) || ''} 
                                alt={getFullName(friend.name, friend.last_name)} 
                                className="w-full h-full object-cover" 
                              />
                            ) : (
                              <span className="text-lg font-semibold text-gray-500">
                                {friend.name.charAt(0).toUpperCase()}
                              </span>
                            )}
                          </div>
                        </button>

                        {/* Info */}
                        <div className="flex-1 min-w-0">
                          <button
                            onClick={() => handleProfileClick(friend.id)}
                            className="text-left"
                          >
                            <p className="font-medium text-sm text-gray-900 hover:text-blue-600 transition-colors truncate">
                              {getFullName(friend.name, friend.last_name)}
                            </p>
                          </button>
                          <p className="text-xs text-gray-500 mt-0.5">
                            {new Date(request.created_at).toLocaleDateString('ru-RU', {
                              day: 'numeric',
                              month: 'short'
                            })}
                          </p>

                          {/* Actions */}
                          <div className="flex gap-2 mt-2">
                            <button
                              onClick={(e) => handleAccept(request.id, e)}
                              disabled={isProcessing}
                              className="flex-1 px-3 py-1.5 bg-blue-500 text-white text-xs font-medium rounded hover:bg-blue-600 transition-colors disabled:opacity-50"
                            >
                              Принять
                            </button>
                            <button
                              onClick={(e) => handleReject(request.id, e)}
                              disabled={isProcessing}
                              className="flex-1 px-3 py-1.5 bg-gray-100 text-gray-700 text-xs font-medium rounded hover:bg-gray-200 transition-colors disabled:opacity-50"
                            >
                              Отклонить
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          {/* Footer */}
          {requests.length > 0 && (
            <div className="px-4 py-3 border-t border-gray-200">
              <button
                onClick={handleViewAll}
                className="w-full text-center text-sm font-medium text-blue-600 hover:text-blue-700"
              >
                Посмотреть все запросы
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
