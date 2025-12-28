'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useAuth } from '../../../contexts/AuthContext';
import { postsApi, petsApi, usersApi, Post, Pet, User } from '../../../lib/api';
import { getMediaUrl } from '../../../lib/utils';
import { 
  UserIcon, 
  MapPinIcon, 
  CalendarIcon, 
  PhoneIcon, 
  EnvelopeIcon,
  PencilIcon,
  CameraIcon,
  HeartIcon,
  ChatBubbleLeftIcon
} from '@heroicons/react/24/outline';
import Image from 'next/image';
import PostComments from '../../components/shared/PostComments';
import PostCard from '../../components/posts/PostCard';
import MediaGallery from '../../components/profile/MediaGallery';
import MediaStats from '../../components/profile/MediaStats';

type TabType = 'posts' | 'media';

export default function UserProfilePage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [pets, setPets] = useState<Pet[]>([]);
  const [profileUser, setProfileUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabType>('posts');
  const { user: currentUser, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const params = useParams();

  // Извлекаем ID из параметра (id1 -> 1)
  const userIdParam = params.userId as string;
  const userId = userIdParam?.startsWith('id') ? parseInt(userIdParam.slice(2)) : null;

  // Проверяем, это профиль текущего пользователя или чужой
  const isOwnProfile = currentUser && userId === currentUser.id;

  useEffect(() => {
    // Проверяем только на клиенте
    if (typeof window === 'undefined') return;

    if (!isLoading && !isAuthenticated) {
      router.push('/auth');
      return;
    }

    if (userId) {
      loadUserProfile();
    }
  }, [userId, isLoading, isAuthenticated]);

  const loadUserProfile = async () => {
    if (!userId) return;

    setLoading(true);
    try {
      // Если это свой профиль, используем данные из контекста
      if (isOwnProfile && currentUser) {
        setProfileUser(currentUser);
      } else {
        // Загружаем данные другого пользователя
        const userResponse = await usersApi.getById(userId);
        if (userResponse.success && userResponse.data) {
          setProfileUser(userResponse.data);
        } else {
          router.push('/');
          return;
        }
      }

      // Загружаем посты и питомцев
      const [postsResponse, petsResponse] = await Promise.all([
        postsApi.getUserPosts(userId),
        petsApi.getUserPets(userId),
      ]);

      if (postsResponse.success && postsResponse.data) {
        setPosts(postsResponse.data);
      }

      if (petsResponse.success && petsResponse.data) {
        setPets(petsResponse.data);
      }
    } catch (error) {
      console.error('Ошибка загрузки профиля:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleEditClick = () => {
    router.push('/profile/edit');
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'только что';
    if (diffMins < 60) return `${diffMins} мин назад`;
    if (diffHours < 24) return `${diffHours} ч назад`;
    if (diffDays < 7) return `${diffDays} дн назад`;
    
    return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' });
  };

  if (isLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
      </div>
    );
  }

  if (!profileUser) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <p className="text-gray-500">Пользователь не найден</p>
        </div>
      </div>
    );
  }

  const profile = {
    name: profileUser.name,
    avatar: profileUser.avatar || null,
    coverPhoto: profileUser.cover_photo || null,
    location: profileUser.location || 'Не указано',
    joinDate: profileUser.created_at ? new Date(profileUser.created_at).toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' }) : 'Недавно',
    phone: profileUser.phone || 'Не указан',
    email: profileUser.email,
    bio: profileUser.bio || 'Информация не заполнена',
  };

  return (
    <div>
      {/* Обложка профиля */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden mb-2.5">
        {/* Cover Photo */}
        <div className="relative h-48 bg-gradient-to-br from-blue-400 via-purple-400 to-pink-400">
          {profile.coverPhoto ? (
            <img src={getMediaUrl(profile.coverPhoto) || ''} alt="Cover" className="w-full h-full object-cover" />
          ) : (
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="text-center">
                <div className="text-6xl font-bold text-white/30 mb-2">
                  {profile.name.charAt(0).toUpperCase()}
                </div>
                <div className="text-white/20 text-sm font-medium">
                  {profile.name}
                </div>
              </div>
            </div>
          )}
          {isOwnProfile && (
            <button 
              className="absolute top-4 right-4 p-2 bg-white rounded-lg shadow-md hover:bg-gray-50 transition-colors"
              title="Изменить обложку"
            >
              <CameraIcon className="w-5 h-5 text-gray-600" />
            </button>
          )}
        </div>

        {/* Profile Info */}
        <div className="p-4 sm:p-6">
          <div className="flex flex-col sm:flex-row items-start gap-4 sm:gap-6">
            {/* Avatar */}
            <div className="relative -mt-16 sm:-mt-20">
              <div className="w-24 h-24 sm:w-32 sm:h-32 rounded-full bg-gray-300 border-4 border-white shadow-lg flex items-center justify-center overflow-hidden">
                {profile.avatar ? (
                  <img src={getMediaUrl(profile.avatar) || ''} alt={profile.name} className="w-full h-full object-cover" />
                ) : (
                  <UserIcon className="w-12 h-12 sm:w-16 sm:h-16 text-gray-500" />
                )}
              </div>
              {isOwnProfile && (
                <button 
                  className="absolute bottom-0 right-0 p-2 bg-white rounded-full shadow-md hover:bg-gray-50 transition-colors border border-gray-200"
                  title="Изменить фото"
                >
                  <CameraIcon className="w-4 h-4 text-gray-600" />
                </button>
              )}
            </div>

            {/* Name and Actions */}
            <div className="flex-1 w-full min-w-0">
              <div className="flex flex-col sm:flex-row items-start justify-between gap-3">
                <div className="min-w-0 flex-1">
                  <h1 className="text-xl sm:text-2xl font-bold text-gray-900 truncate">{profile.name}</h1>
                  <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4 mt-1 text-xs sm:text-sm text-gray-600">
                    <div className="flex items-center gap-1">
                      <MapPinIcon className="w-4 h-4 flex-shrink-0" />
                      <span className="truncate">{profile.location}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <CalendarIcon className="w-4 h-4 flex-shrink-0" />
                      <span className="truncate">На платформе с {profile.joinDate}</span>
                    </div>
                  </div>
                </div>

                {isOwnProfile && (
                  <button 
                    onClick={handleEditClick}
                    className="flex items-center gap-2 px-3 sm:px-4 py-2 rounded-lg text-white text-sm font-medium transition-colors whitespace-nowrap"
                    style={{ backgroundColor: '#1B76FF' }}
                  >
                    <PencilIcon className="w-4 h-4" />
                    <span className="hidden sm:inline">Редактировать профиль</span>
                    <span className="sm:hidden">Редактировать</span>
                  </button>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-2.5">
        {/* Left Column - Activity */}
        <div className="lg:col-span-2 space-y-2.5">
          {/* Tabs */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-100">
            <div className="flex border-b border-gray-200 overflow-x-auto">
              <button 
                onClick={() => setActiveTab('posts')}
                className={`flex-1 px-3 sm:px-6 py-3 text-xs sm:text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${
                  activeTab === 'posts'
                    ? 'border-[#1B76FF] text-[#1B76FF]'
                    : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
              >
                Публикации
              </button>
              <button 
                onClick={() => setActiveTab('media')}
                className={`flex-1 px-3 sm:px-6 py-3 text-xs sm:text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${
                  activeTab === 'media'
                    ? 'border-[#1B76FF] text-[#1B76FF]'
                    : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
              >
                Медиа
              </button>
            </div>

            {/* Content */}
            <div className="p-4">
              {activeTab === 'posts' && (
                <div className="space-y-2.5">
                  {loading ? (
                    <div className="flex justify-center py-8">
                      <div className="animate-spin rounded-full h-8 w-8 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
                    </div>
                  ) : posts.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                      <p>Пока нет публикаций</p>
                    </div>
                  ) : (
                    <>
                      {posts.map((post) => (
                        <PostCard key={post.id} post={post} />
                      ))}
                    </>
                  )}
                </div>
              )}

              {activeTab === 'media' && userId && (
                <MediaGallery userId={userId} mediaType="all" />
              )}
            </div>
          </div>
        </div>

        {/* Right Column - Info */}
        <div className="lg:col-span-1 space-y-2.5">
          {/* О себе */}
          <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">О себе</h2>
            <p className="text-sm text-gray-700 leading-relaxed">{profile.bio}</p>
          </div>

          {/* Контактная информация */}
          <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Контакты</h2>
            <div className="space-y-3">
              <div className="flex items-center gap-3 text-sm">
                <PhoneIcon className="w-5 h-5 text-gray-400" />
                <span className="text-gray-700">{profile.phone}</span>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <EnvelopeIcon className="w-5 h-5 text-gray-400" />
                <span className="text-gray-700">{profile.email}</span>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <MapPinIcon className="w-5 h-5 text-gray-400" />
                <span className="text-gray-700">{profile.location}</span>
              </div>
            </div>
          </div>

          {/* Мои питомцы */}
          <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Питомцы</h2>
              {isOwnProfile && (
                <button className="text-sm font-medium" style={{ color: '#1B76FF' }}>
                  Добавить
                </button>
              )}
            </div>
            {loading ? (
              <div className="flex justify-center py-4">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
              </div>
            ) : pets.length === 0 ? (
              <div className="text-center py-4 text-sm text-gray-500">
                <p>Питомцы не добавлены</p>
              </div>
            ) : (
              <div className="grid grid-cols-3 gap-3">
                {pets.map((pet) => (
                  <div key={pet.id} className="aspect-square rounded-lg bg-gray-200 flex flex-col items-center justify-center overflow-hidden p-2">
                    {pet.photo ? (
                      <Image src={pet.photo} alt={pet.name} fill className="object-cover" />
                    ) : (
                      <>
                        <span className="text-3xl mb-1">🐕</span>
                        <span className="text-xs text-gray-600 text-center">{pet.name}</span>
                      </>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Статистика медиа */}
          {isOwnProfile && <MediaStats />}
        </div>
      </div>
    </div>
  );
}
