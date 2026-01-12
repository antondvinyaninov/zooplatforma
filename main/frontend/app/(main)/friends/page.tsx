'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../../contexts/AuthContext';
import { friendsApi, Friendship } from '../../../lib/api';
import { getMediaUrl, getFullName } from '../../../lib/utils';
import { UserIcon, MagnifyingGlassIcon } from '@heroicons/react/24/outline';

export default function FriendsPage() {
  const [friends, setFriends] = useState<Friendship[]>([]);
  const [filteredFriends, setFilteredFriends] = useState<Friendship[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const { user: currentUser, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (typeof window === 'undefined') return;

    if (!isLoading && !isAuthenticated) {
      router.push('/auth');
      return;
    }

    if (currentUser) {
      loadFriends();
    }
  }, [isLoading, isAuthenticated, currentUser]);

  useEffect(() => {
    // Фильтрация друзей по поисковому запросу
    if (searchQuery.trim() === '') {
      setFilteredFriends(friends);
    } else {
      const query = searchQuery.toLowerCase();
      const filtered = friends.filter(friendship => {
        const friend = friendship.friend;
        const fullName = getFullName(friend.name, friend.last_name).toLowerCase();
        return fullName.includes(query);
      });
      setFilteredFriends(filtered);
    }
  }, [searchQuery, friends]);

  const loadFriends = async () => {
    setLoading(true);
    try {
      const response = await friendsApi.getFriends();
      if (response.success && response.data) {
        const acceptedFriends = response.data.filter(f => f.status === 'accepted');
        setFriends(acceptedFriends);
        setFilteredFriends(acceptedFriends);
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
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Друзья</h1>
        
        {/* Search */}
        <div className="relative">
          <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Поиск друзей..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Stats */}
        <div className="mt-4 text-sm text-gray-600">
          {filteredFriends.length === friends.length ? (
            <span>Всего друзей: {friends.length}</span>
          ) : (
            <span>Найдено: {filteredFriends.length} из {friends.length}</span>
          )}
        </div>
      </div>

      {/* Friends Grid */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
        {filteredFriends.length === 0 ? (
          <div className="text-center py-12">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
              <UserIcon className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-600 mb-2">
              {searchQuery ? 'Друзья не найдены' : 'У вас пока нет друзей'}
            </p>
            {!searchQuery && (
              <p className="text-sm text-gray-500">
                Найдите интересных людей и добавьте их в друзья
              </p>
            )}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            {filteredFriends.map((friendship) => {
              const friend = friendship.friend;
              return (
                <button
                  key={friendship.id}
                  onClick={() => handleFriendClick(friend.id)}
                  className="flex items-center gap-3 p-4 rounded-lg border border-gray-200 hover:border-blue-500 hover:bg-blue-50 transition-all"
                >
                  <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden flex-shrink-0">
                    {friend.avatar ? (
                      <img 
                        src={getMediaUrl(friend.avatar) || ''} 
                        alt={getFullName(friend.name, friend.last_name)} 
                        className="w-full h-full object-cover" 
                      />
                    ) : (
                      <UserIcon className="w-6 h-6 text-gray-400" />
                    )}
                  </div>
                  <div className="flex-1 text-left min-w-0">
                    <p className="font-medium text-gray-900 truncate">
                      {getFullName(friend.name, friend.last_name)}
                    </p>
                    {friend.location && (
                      <p className="text-xs text-gray-500 truncate">{friend.location}</p>
                    )}
                  </div>
                </button>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
