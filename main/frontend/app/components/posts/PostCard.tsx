'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import { UserIcon } from '@heroicons/react/24/outline';
import PostComments from '../shared/PostComments';
import PostModal from './PostModal';
import PollDisplay from '../polls/PollDisplay';

interface User {
  id: number;
  name: string;
  email: string;
  avatar?: string;
}

interface Pet {
  id: number;
  name: string;
  species: string;
  photo?: string;
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
  pets?: Pet[];
  poll?: Poll;
  comments_count: number;
}

interface PostCardProps {
  post: Post;
  onDelete?: (postId: number) => void;
}

export default function PostCard({ post, onDelete }: PostCardProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [showMenu, setShowMenu] = useState(false);
  const [isLiked, setIsLiked] = useState(false);
  const [likesCount, setLikesCount] = useState(0);
  const [commentsCount, setCommentsCount] = useState(post.comments_count);
  const [showModal, setShowModal] = useState(false);

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
      }
    } catch (error) {
      console.error('Ошибка лайка:', error);
    }
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
            {post.user?.avatar ? (
              <img src={post.user.avatar} alt={post.user.name} className="w-full h-full object-cover" />
            ) : (
              <UserIcon className="w-6 h-6 text-gray-500" />
            )}
          </div>

          {/* User Info */}
          <div>
            <div className="font-semibold text-gray-900">{post.user?.name || 'Пользователь'}</div>
            <div className="text-sm text-gray-500">{getTimeAgo(post.created_at)}</div>
          </div>
        </div>

        {/* Menu Button */}
        <button
          onClick={() => setShowMenu(!showMenu)}
          className="p-2 hover:bg-gray-100 rounded-full transition-colors"
        >
          <svg className="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z" />
          </svg>
        </button>
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
      <div className="text-gray-900 px-4 pb-3 whitespace-pre-wrap">{post.content}</div>

      {/* Poll */}
      {post.poll && (
        <div className="px-4 pb-3">
          <PollDisplay poll={post.poll} />
        </div>
      )}

      {/* Attached Pets */}
      {post.pets && post.pets.length > 0 && (
        <div className="space-y-2 px-4 pb-3">
          {post.pets.map((pet) => (
            <div
              key={pet.id}
              className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg border border-gray-200 hover:border-gray-300 transition-colors cursor-pointer"
            >
              {/* Pet Photo */}
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-400 to-cyan-400 flex items-center justify-center text-white font-semibold overflow-hidden">
                {pet.photo ? (
                  <img src={pet.photo} alt={pet.name} className="w-full h-full object-cover" />
                ) : (
                  pet.name[0]?.toUpperCase()
                )}
              </div>

              {/* Pet Info */}
              <div className="flex-1">
                <div className="font-semibold text-gray-900">{pet.name}</div>
                <div className="text-sm text-gray-600">{pet.species}</div>
              </div>

              {/* Arrow */}
              <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
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
          {likesCount > 0 && <span className="text-sm">{likesCount}</span>}
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

        <button className="flex items-center gap-2 hover:text-green-500 transition-colors">
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>

        <button className="flex items-center gap-2 hover:text-purple-500 transition-colors">
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
        </button>
      </div>
      
      {/* Post Modal */}
      <PostModal
        post={post}
        isOpen={showModal}
        onClose={handleCloseModal}
        onCountChange={(count) => setCommentsCount(count)}
      />
    </div>
  );
}
