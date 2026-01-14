'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import { UserIcon } from '@heroicons/react/24/outline';
import { Pencil, Trash2 } from 'lucide-react';
import { getMediaUrl, getFullName } from '@/lib/utils';
import { useAuth } from '@/contexts/AuthContext';
import PostComments from '../shared/PostComments';
import PostModal from './PostModal';
import PollDisplay from '../polls/PollDisplay';
import PhotoGrid from './PhotoGrid';
import LikersModal from './LikersModal';
import PetCard from './PetCard';
import ReportButton from './ReportButton';

interface User {
  id: number;
  name: string;
  last_name?: string;
  email: string;
  avatar?: string;
}

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  gender?: string;
  birth_date?: string;
  color?: string;
  size?: string;
  photo?: string;
  status?: string;
  city?: string;
  region?: string;
  urgent?: boolean;
  story?: string;
  organization_name?: string;
  organization_type?: string;
}

interface PollOption {
  id: number;
  option_text: string;
  votes_count: number;
  percentage?: number;
  voters?: PollVoter[];
}

interface PollVoter {
  user_id: number;
  user_name: string;
  avatar?: string;
}

interface Poll {
  id: number;
  question: string;
  multiple_choice: boolean;
  allow_vote_changes: boolean;
  anonymous_voting: boolean;
  expires_at?: string;
  options: PollOption[];
  total_voters: number;
  user_voted: boolean;
  user_votes?: number[];
  is_expired: boolean;
  voters?: PollVoter[];
}

interface Organization {
  id: number;
  name: string;
  short_name?: string;
  logo?: string;
  type?: string;
}

interface Post {
  id: number;
  author_id: number;
  author_type: string;
  content: string;
  attached_pets: number[];
  attachments: any[];
  tags: string[];
  created_at: string;
  updated_at: string;
  user?: User;
  organization?: Organization;
  pets?: Pet[];
  poll?: Poll;
  comments_count?: number;
}

interface PostCardProps {
  post: Post;
  onDelete?: (postId: number) => void;
}

export default function PostCard({ post, onDelete }: PostCardProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user } = useAuth();
  const [showMenu, setShowMenu] = useState(false);
  const [isLiked, setIsLiked] = useState(false);
  const [likesCount, setLikesCount] = useState(0);
  const [commentsCount, setCommentsCount] = useState(post.comments_count || 0);
  const [showModal, setShowModal] = useState(false);
  const [showLikersModal, setShowLikersModal] = useState(false);
  const [showShareMenu, setShowShareMenu] = useState(false);
  const [shareMessage, setShareMessage] = useState('');
  const [showLikeAnimation, setShowLikeAnimation] = useState(false);
  const [showReportModal, setShowReportModal] = useState(false);
  const [canEditPost, setCanEditPost] = useState(false);

  // Проверка прав на редактирование поста
  useEffect(() => {
    const checkEditPermission = async () => {
      if (!user) {
        setCanEditPost(false);
        return;
      }

      // Если пост от пользователя - проверяем ID
      if (post.author_type === 'user' && user.id === post.author_id) {
        setCanEditPost(true);
        return;
      }

      // Если пост от организации - проверяем членство
      if (post.author_type === 'organization' && post.organization) {
        try {
          const response = await fetch(
            `http://localhost:8000/api/organizations/${post.organization.id}/members`,
            { credentials: 'include' }
          );
          
          if (response.ok) {
            const data = await response.json();
            // Проверяем, есть ли текущий пользователь в списке участников
            const member = data.members?.find((m: any) => m.user_id === user.id);
            // Разрешаем редактирование owner, admin, moderator
            if (member && ['owner', 'admin', 'moderator'].includes(member.role)) {
              setCanEditPost(true);
              return;
            }
          }
        } catch (error) {
          console.error('Ошибка проверки прав:', error);
        }
      }

      setCanEditPost(false);
    };

    checkEditPermission();
  }, [user, post.author_id, post.author_type, post.organization]);

  // Проверяем, открыт ли этот пост через URL
  useEffect(() => {
    const postId = searchParams.get('metka');
    if (postId === String(post.id)) {
      setShowModal(true);
    }
  }, [searchParams, post.id]);

  useEffect(() => {
    loadLikeStatus();
  }, [post.id]);

  const loadLikeStatus = async () => {
    try {
      const response = await fetch(`http://localhost:8000/api/posts/${post.id}/like`, {
        method: 'GET',
        credentials: 'include',
      });
      
      const result = await response.json();
      if (result.success && result.data) {
        setIsLiked(result.data.liked);
        setLikesCount(result.data.likes_count);
      }
    } catch (error) {
      console.error('Ошибка загрузки статуса лайка:', error);
    }
  };

  const handleLike = async () => {
    try {
      const response = await fetch(`http://localhost:8000/api/posts/${post.id}/like`, {
        method: 'POST',
        credentials: 'include',
      });
      
      const result = await response.json();
      if (result.success && result.data) {
        setIsLiked(result.data.liked);
        setLikesCount(result.data.likes_count);
        
        // Показываем анимацию только при лайке (не при анлайке)
        if (result.data.liked) {
          setShowLikeAnimation(true);
          setTimeout(() => setShowLikeAnimation(false), 1000);
        }
      }
    } catch (error) {
      console.error('Ошибка лайка:', error);
    }
  };

  const handleDoubleClick = () => {
    if (!isLiked) {
      handleLike();
    }
  };

  const handleShare = async (type: 'copy' | 'vk' | 'telegram' | 'whatsapp') => {
    const postUrl = `${window.location.origin}/?metka=${post.id}`;
    const text = post.content.substring(0, 100) + (post.content.length > 100 ? '...' : '');

    switch (type) {
      case 'copy':
        try {
          await navigator.clipboard.writeText(postUrl);
          setShareMessage('Ссылка скопирована!');
          setTimeout(() => setShareMessage(''), 2000);
        } catch (error) {
          console.error('Ошибка копирования:', error);
        }
        break;
      case 'vk':
        window.open(`https://vk.com/share.php?url=${encodeURIComponent(postUrl)}&title=${encodeURIComponent(text)}`, '_blank');
        break;
      case 'telegram':
        window.open(`https://t.me/share/url?url=${encodeURIComponent(postUrl)}&text=${encodeURIComponent(text)}`, '_blank');
        break;
      case 'whatsapp':
        window.open(`https://wa.me/?text=${encodeURIComponent(text + ' ' + postUrl)}`, '_blank');
        break;
    }
    setShowShareMenu(false);
  };

  const handleOpenModal = () => {
    setShowModal(true);
    // Добавляем параметр metka в URL
    router.push(`/?metka=${post.id}`, { scroll: false });
  };

  const handleCloseModal = () => {
    setShowModal(false);
    // Убираем параметр metka из URL
    router.push('/', { scroll: false });
  };

  const handleDeletePost = async () => {
    if (!confirm('Вы уверены, что хотите удалить этот пост?')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8000/api/posts/${post.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (response.ok) {
        if (onDelete) {
          onDelete(post.id);
        }
        // Можно добавить уведомление об успешном удалении
      } else {
        alert('Ошибка при удалении поста');
      }
    } catch (error) {
      console.error('Ошибка удаления поста:', error);
      alert('Ошибка при удалении поста');
    }
    setShowMenu(false);
  };

  const handleEditPost = () => {
    // TODO: Открыть модальное окно редактирования или перейти на страницу редактирования
    alert('Функция редактирования в разработке');
    setShowMenu(false);
  };

  const getTimeAgo = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true, locale: ru });
    } catch {
      return 'недавно';
    }
  };

  const getTagTextColor = (tag: string) => {
    switch (tag) {
      case 'потерян':
        return 'text-red-700';
      case 'найден':
        return 'text-green-700';
      case 'ищет дом':
        return 'text-orange-700';
      default:
        return 'text-gray-700';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow">
      {/* Header */}
      <div className="flex items-start justify-between p-4 pb-3">
        <div className="flex items-center gap-3">
          {/* Avatar */}
          <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center text-white font-semibold overflow-hidden">
            {post.author_type === 'organization' ? (
              post.organization?.logo ? (
                <img src={getMediaUrl(post.organization.logo) || ''} alt={post.organization.name} className="w-full h-full object-cover" />
              ) : (
                <UserIcon className="w-6 h-6 text-gray-500" />
              )
            ) : (
              post.user?.avatar ? (
                <img src={getMediaUrl(post.user.avatar) || ''} alt={post.user.name} className="w-full h-full object-cover" />
              ) : (
                <UserIcon className="w-6 h-6 text-gray-500" />
              )
            )}
          </div>

          {/* User Info */}
          <div>
            <div className="font-semibold text-gray-900">
              {post.author_type === 'organization' 
                ? (post.organization?.short_name || post.organization?.name || 'Организация')
                : getFullName(post.user?.name || 'Пользователь', post.user?.last_name)
              }
            </div>
            <div className="text-sm text-gray-500">{getTimeAgo(post.created_at)}</div>
          </div>
        </div>

        {/* Menu Button */}
        <div className="relative">
          <button
            onClick={() => setShowMenu(!showMenu)}
            className="p-2 hover:bg-gray-100 rounded-full transition-colors"
          >
            <svg className="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z" />
            </svg>
          </button>

          {/* Dropdown Menu */}
          {showMenu && (
            <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-10">
              {/* Действия автора */}
              {canEditPost && (
                <>
                  <button
                    onClick={handleEditPost}
                    className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-gray-700 transition-colors"
                  >
                    <Pencil className="w-5 h-5 text-blue-500" />
                    <span>Редактировать</span>
                  </button>
                  <button
                    onClick={handleDeletePost}
                    className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-red-600 transition-colors"
                  >
                    <Trash2 className="w-5 h-5" />
                    <span>Удалить</span>
                  </button>
                  <div className="border-t border-gray-200 my-2"></div>
                </>
              )}
              
              {/* Пожаловаться */}
              <button
                onClick={() => {
                  setShowReportModal(true);
                  setShowMenu(false);
                }}
                className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-gray-700 transition-colors"
              >
                <svg className="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9" />
                </svg>
                <span>Пожаловаться</span>
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Tags */}
      {post.tags.length > 0 && (
        <div className="flex gap-2 px-4 pb-3">
          {post.tags.map((tag, index) => (
            <span
              key={index}
              className={`px-3 py-1 rounded-full text-sm font-medium ${getTagTextColor(tag)}`}
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* Content */}
      <div 
        className="text-gray-900 px-4 pb-3 whitespace-pre-wrap relative"
        onDoubleClick={handleDoubleClick}
      >
        {post.content}
        
        {/* Like Animation */}
        {showLikeAnimation && (
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <svg 
              className="w-24 h-24 text-red-500 fill-current animate-ping"
              viewBox="0 0 24 24"
            >
              <path d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
            </svg>
          </div>
        )}
      </div>

      {/* Photos/Videos */}
      {post.attachments && post.attachments.length > 0 && (
        <div className="pb-3">
          <PhotoGrid 
            photos={post.attachments.filter(a => a.type === 'image' || a.type === 'video' || a.media_type === 'image' || a.media_type === 'video')} 
            onClick={handleOpenModal}
          />
        </div>
      )}

      {/* Poll */}
      {post.poll && (
        <div className="px-4 pb-3">
          <PollDisplay poll={post.poll} />
        </div>
      )}

      {/* Attached Pets */}
      {post.pets && post.pets.length > 0 && (
        <div className="space-y-3 px-4 pb-3">
          {post.pets.map((pet) => (
            <PetCard key={pet.id} pet={pet} />
          ))}
        </div>
      )}

      {/* Actions */}
      <div className="flex items-center gap-6 text-gray-600 px-4 py-3 border-t border-gray-100">
        <button 
          onClick={handleLike}
          className={`flex items-center gap-2 transition-colors ${isLiked ? 'text-red-500' : 'hover:text-red-500'}`}
        >
          <svg className={`w-5 h-5 ${isLiked ? 'fill-current' : ''}`} fill={isLiked ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          {likesCount > 0 && (
            <span 
              className="text-sm hover:underline cursor-pointer"
              onClick={(e) => {
                e.stopPropagation();
                setShowLikersModal(true);
              }}
            >
              {likesCount}
            </span>
          )}
        </button>

        <button 
          onClick={handleOpenModal}
          className="flex items-center gap-2 hover:text-blue-500 transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          {commentsCount > 0 && <span className="text-sm">{commentsCount}</span>}
        </button>

        <div className="relative">
          <button 
            onClick={() => setShowShareMenu(!showShareMenu)}
            className="flex items-center gap-2 hover:text-purple-500 transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
            </svg>
          </button>

          {/* Share Menu */}
          {showShareMenu && (
            <div className="absolute bottom-full right-0 mb-2 bg-white rounded-lg shadow-lg border border-gray-200 py-2 min-w-[200px] z-10">
              <button
                onClick={() => handleShare('copy')}
                className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-gray-700"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                Копировать ссылку
              </button>
              <button
                onClick={() => handleShare('vk')}
                className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-gray-700"
              >
                <span className="text-blue-500 font-bold">VK</span>
                ВКонтакте
              </button>
              <button
                onClick={() => handleShare('telegram')}
                className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-gray-700"
              >
                <span className="text-blue-400">✈️</span>
                Telegram
              </button>
              <button
                onClick={() => handleShare('whatsapp')}
                className="w-full px-4 py-2 text-left hover:bg-gray-50 flex items-center gap-3 text-gray-700"
              >
                <span className="text-green-500">💬</span>
                WhatsApp
              </button>
            </div>
          )}

          {/* Share Message */}
          {shareMessage && (
            <div className="absolute bottom-full right-0 mb-2 bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg text-sm whitespace-nowrap">
              {shareMessage}
            </div>
          )}
        </div>
      </div>
      
      {/* Post Modal */}
      <PostModal
        post={post}
        isOpen={showModal}
        onClose={handleCloseModal}
        onCountChange={(count) => setCommentsCount(count)}
      />

      {/* Likers Modal */}
      <LikersModal
        postId={post.id}
        isOpen={showLikersModal}
        onClose={() => setShowLikersModal(false)}
      />

      {/* Report Modal */}
      {showReportModal && (
        <ReportButton 
          targetType="post" 
          targetId={post.id} 
          targetName={`Пост от ${post.author_type === 'organization' ? (post.organization?.short_name || post.organization?.name) : getFullName(post.user?.name || '', post.user?.last_name)}`}
          isOpen={showReportModal}
          onClose={() => setShowReportModal(false)}
        />
      )}
    </div>
  );
}
