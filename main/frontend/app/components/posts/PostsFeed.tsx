'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import PostCard from './PostCard';
import CreatePost from './CreatePost';

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
  user?: any;
  pets?: any[];
  comments_count: number;
}

interface PostsFeedProps {
  activeFilter?: 'for-you' | 'following' | 'city' | 'lost' | 'found' | 'looking-for-home';
}

export default function PostsFeed({ activeFilter = 'for-you' }: PostsFeedProps) {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  const loadPosts = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get('/api/posts');
      setPosts(response.data || []);
    } catch (error) {
      console.error('Ошибка загрузки постов:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPosts();
  }, [activeFilter]);

  return (
    <div className="space-y-2.5">
      {/* Create Post */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <CreatePost onPostCreated={loadPosts} />
      </div>

      {/* Posts */}
      {loading ? (
        <div className="flex items-center justify-center py-12 bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
      ) : posts.length === 0 ? (
        <div className="text-center py-12 text-gray-500 bg-white rounded-lg shadow-sm border border-gray-200">
          <p>Пока нет постов</p>
          <p className="text-sm mt-2">Создайте первый пост!</p>
        </div>
      ) : (
        <>
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </>
      )}
    </div>
  );
}
