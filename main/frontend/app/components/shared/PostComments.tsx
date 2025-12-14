'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '../../../contexts/AuthContext';
import { commentsApi, Comment } from '../../../lib/api';
import { UserIcon, TrashIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';

interface PostCommentsProps {
  postId: number;
  initialCount?: number;
  onCountChange?: (count: number) => void;
}

export default function PostComments({ postId, initialCount = 0, onCountChange }: PostCommentsProps) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState('');
  const [replyTo, setReplyTo] = useState<{ commentId: number; userName: string; userId: number } | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showComments, setShowComments] = useState(false);
  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    if (showComments && comments.length === 0) {
      loadComments();
    }
  }, [showComments]);

  const loadComments = async () => {
    setIsLoading(true);
    try {
      const response = await commentsApi.getPostComments(postId);
      if (response.success && response.data) {
        setComments(response.data);
        // Обновляем счетчик в родительском компоненте
        const count = response.data.reduce((sum, c) => sum + 1 + (c.replies?.length || 0), 0);
        onCountChange?.(count);
      }
    } catch (error) {
      console.error('Ошибка загрузки комментариев:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim() || !isAuthenticated || isSubmitting) return;

    setIsSubmitting(true);
    try {
      const response = await commentsApi.create(postId, {
        content: newComment.trim(),
        parent_id: replyTo?.commentId,
        reply_to_user_id: replyTo?.userId,
      });

      if (response.success && response.data) {
        let newComments;
        // Если это ответ, добавляем в replies родительского комментария
        if (replyTo) {
          newComments = comments.map(c => {
            if (c.id === replyTo.commentId) {
              return {
                ...c,
                replies: [...(c.replies || []), response.data!]
              };
            }
            return c;
          });
        } else {
          // Иначе добавляем как новый корневой комментарий
          newComments = [...comments, response.data];
        }
        setComments(newComments);
        
        // Обновляем счетчик в родительском компоненте
        const count = newComments.reduce((sum, c) => sum + 1 + (c.replies?.length || 0), 0);
        onCountChange?.(count);
        
        setNewComment('');
        setReplyTo(null);
      } else {
        alert(response.error || 'Ошибка создания комментария');
      }
    } catch (error) {
      console.error('Ошибка создания комментария:', error);
      alert('Ошибка создания комментария');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (commentId: number, parentId?: number) => {
    if (!confirm('Удалить комментарий?')) return;

    try {
      const response = await commentsApi.delete(commentId);
      if (response.success) {
        let newComments;
        if (parentId) {
          // Удаляем ответ из replies
          newComments = comments.map(c => {
            if (c.id === parentId) {
              return {
                ...c,
                replies: c.replies?.filter(r => r.id !== commentId)
              };
            }
            return c;
          });
        } else {
          // Удаляем корневой комментарий
          newComments = comments.filter(c => c.id !== commentId);
        }
        setComments(newComments);
        
        // Обновляем счетчик в родительском компоненте
        const count = newComments.reduce((sum, c) => sum + 1 + (c.replies?.length || 0), 0);
        onCountChange?.(count);
      } else {
        alert(response.error || 'Ошибка удаления комментария');
      }
    } catch (error) {
      console.error('Ошибка удаления комментария:', error);
      alert('Ошибка удаления комментария');
    }
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

  const renderComment = (comment: Comment, isReply = false, parentCommentId?: number) => (
    <div key={comment.id} className={`flex gap-2 ${isReply ? 'ml-10' : ''}`}>
      <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center flex-shrink-0 overflow-hidden">
        {comment.user?.avatar ? (
          <Image src={comment.user.avatar} alt={comment.user.name} width={32} height={32} className="object-cover" />
        ) : (
          <UserIcon className="w-4 h-4 text-gray-500" />
        )}
      </div>
      <div className="flex-1">
        <div className="bg-gray-100 rounded-lg px-3 py-2">
          <div className="flex items-center justify-between mb-1">
            <span className="text-sm font-semibold text-gray-900">{comment.user?.name}</span>
            {user && comment.user_id === user.id && (
              <button
                onClick={() => handleDelete(comment.id, comment.parent_id)}
                className="text-gray-400 hover:text-red-500 transition-colors"
                title="Удалить"
              >
                <TrashIcon className="w-3.5 h-3.5" />
              </button>
            )}
          </div>
          <p className="text-sm text-gray-700 whitespace-pre-wrap">
            {comment.reply_to_user && (
              <span className="font-semibold" style={{ color: '#1B76FF' }}>
                {comment.reply_to_user.name},{' '}
              </span>
            )}
            {comment.content}
          </p>
        </div>
        <div className="flex items-center gap-3 mt-1 text-xs text-gray-500 px-3">
          <span>{formatDate(comment.created_at)}</span>
          {isAuthenticated && (
            <button
              onClick={() => setReplyTo({ 
                commentId: parentCommentId || comment.id, 
                userName: comment.user?.name || '', 
                userId: comment.user_id 
              })}
              className="hover:text-gray-700 transition-colors"
            >
              Ответить
            </button>
          )}
        </div>
      </div>
    </div>
  );

  // Считаем количество комментариев динамически
  const totalComments = showComments 
    ? comments.reduce((sum, c) => sum + 1 + (c.replies?.length || 0), 0)
    : initialCount;

  return (
    <div className="mt-3 border-t border-gray-200 pt-3">
      {/* Кнопка показать/скрыть комментарии */}
      <button
        onClick={() => setShowComments(!showComments)}
        className="text-sm text-gray-600 hover:text-gray-900 transition-colors mb-3"
      >
        {showComments ? 'Скрыть комментарии' : `Показать комментарии (${totalComments})`}
      </button>

      {showComments && (
        <div className="space-y-3">
          {/* Список комментариев */}
          {isLoading ? (
            <div className="flex justify-center py-4">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
            </div>
          ) : comments.length === 0 ? (
            <p className="text-sm text-gray-500 text-center py-2">Комментариев пока нет</p>
          ) : (
            <div className="space-y-3">
              {comments.map((comment) => (
                <div key={comment.id}>
                  {renderComment(comment)}
                  {/* Ответы на комментарий */}
                  {comment.replies && comment.replies.length > 0 && (
                    <div className="mt-2 space-y-2">
                      {comment.replies.map((reply) => renderComment(reply, true, comment.id))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Форма добавления комментария */}
          {isAuthenticated && (
            <form onSubmit={handleSubmit} className="mt-3">
              {replyTo && (
                <div className="mb-2 text-sm text-gray-600 flex items-center gap-2">
                  <span>Ответ для <span className="font-semibold">{replyTo.userName}</span></span>
                  <button
                    type="button"
                    onClick={() => setReplyTo(null)}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    ✕
                  </button>
                </div>
              )}
              <div className="flex gap-2">
                <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center flex-shrink-0 overflow-hidden">
                  {user?.avatar ? (
                    <Image src={user.avatar} alt={user.name} width={32} height={32} className="object-cover" />
                  ) : (
                    <UserIcon className="w-4 h-4 text-gray-500" />
                  )}
                </div>
                <div className="flex-1 flex gap-2">
                  <input
                    type="text"
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                    placeholder={replyTo ? `Ответить ${replyTo.userName}...` : "Написать комментарий..."}
                    className="flex-1 px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:border-transparent bg-white"
                    style={{ '--tw-ring-color': '#1B76FF' } as React.CSSProperties}
                  />
                  <button
                    type="submit"
                    disabled={!newComment.trim() || isSubmitting}
                    className="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    style={{ backgroundColor: '#1B76FF' }}
                  >
                    {isSubmitting ? '...' : 'Отправить'}
                  </button>
                </div>
              </div>
            </form>
          )}
        </div>
      )}
    </div>
  );
}
