'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { friendsApi, Friendship } from '@/lib/api';
import { getMediaUrl, getFullName } from '@/lib/utils';
import { UserIcon } from '@heroicons/react/24/outline';

interface FriendsListProps {
  userId: number;
  limit?: number;
  showViewAll?: boolean;
}

export default function FriendsList({ userId, limit, showViewAll = true }: FriendsListProps) {
  const [friends, setFriends] = useState<Friendship[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    loadFriends();
  }, [userId]);

  const loadFriends = async () => {
    setLoading(true);
    try {
      const response = await friendsApi.getFriends();
      if (response.success && response.data) {
        // Фильтруем друзей для конкретного пользователя
        const userFriends = response.data.filter(
          f => f.status === 'accepted' && (f.user_id === userId || f.friend_id === userId)
        );
        setFriends(limit ? userFriends.slice(0, limit) : userFriends);
      }
    } catch (error) {
      console.error('Ошибка загрузки друзей:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFriendClick = (friendId: number) => {
    router.push(`/id${friendId}`);
  };

  const handleViewAll = () => {
    router.push(`/id${userId}/friends`);
  };

  if (loading) {
    return (
      <div className="flex justify-center py-4">
        <div className="animate-spin rounded-full h-6 w-6 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
      </div>
    );
  }

  if (friends.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-gray-100 mb-3">
          <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
        <p className="text-sm text-gray-600">Пока нет друзей</p>
      </div>
    );
  }

  return (
    <div>
      {showViewAll && (
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">
            Друзья {friends.length > 0 && `(${friends.length})`}
          </h2>
          <button 
            onClick={handleViewAll}
            className="text-sm font-medium" 
            style={{ color: '#1B76FF' }}
          >
            Все друзья
          </button>
        </div>
      )}

      <div className="grid grid-cols-3 gap-3">
        {friends.map((friendship) => {
          const friend = friendship.friend;
          return (
            <button
              key={friendship.id}
              onClick={() => handleFriendClick(friend.id)}
              className="flex flex-col items-center p-2 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden mb-2">
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
              <span className="text-xs text-gray-900 text-center line-clamp-2">
                {getFullName(friend.name, friend.last_name)}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
