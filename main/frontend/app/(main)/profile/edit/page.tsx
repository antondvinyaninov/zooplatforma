'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../../../contexts/AuthContext';
import { usersApi } from '../../../../lib/api';
import { UserIcon, CameraIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';

export default function EditProfilePage() {
  const [isSaving, setIsSaving] = useState(false);
  const [editForm, setEditForm] = useState({
    name: '',
    bio: '',
    phone: '',
    location: '',
  });
  const { user, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/auth');
      return;
    }

    if (user) {
      setEditForm({
        name: user.name,
        bio: user.bio || '',
        phone: user.phone || '',
        location: user.location || '',
      });
    }
  }, [user, isLoading, isAuthenticated, router]);

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || isSaving) return;

    setIsSaving(true);
    try {
      const response = await usersApi.updateProfile(editForm);
      
      if (response.success) {
        alert('Профиль успешно обновлен!');
        router.push(`/id${user.id}`);
      } else {
        alert(response.error || 'Ошибка обновления профиля');
      }
    } catch (error) {
      console.error('Ошибка обновления профиля:', error);
      alert('Ошибка обновления профиля');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2" style={{ borderColor: '#1B76FF' }}></div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        {/* Header */}
        <div className="border-b border-gray-200 p-6">
          <h1 className="text-2xl font-bold text-gray-900">Редактирование профиля</h1>
          <p className="text-sm text-gray-600 mt-1">Обновите информацию о себе</p>
        </div>

        {/* Form */}
        <form onSubmit={handleSaveProfile} className="p-6 space-y-6">
          {/* Avatar Section */}
          <div className="flex items-center gap-6">
            <div className="relative">
              <div className="w-24 h-24 rounded-full bg-gray-300 flex items-center justify-center overflow-hidden">
                {user.avatar ? (
                  <Image src={user.avatar} alt={user.name} fill className="object-cover" />
                ) : (
                  <UserIcon className="w-12 h-12 text-gray-500" />
                )}
              </div>
              <button 
                type="button"
                className="absolute bottom-0 right-0 p-2 bg-white rounded-full shadow-md hover:bg-gray-50 transition-colors border border-gray-200"
                title="Изменить фото"
              >
                <CameraIcon className="w-4 h-4 text-gray-600" />
              </button>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900">{user.name}</h3>
              <p className="text-sm text-gray-600">{user.email}</p>
              <button 
                type="button"
                className="mt-2 text-sm font-medium"
                style={{ color: '#1B76FF' }}
              >
                Изменить фото профиля
              </button>
            </div>
          </div>

          {/* Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Имя <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={editForm.name}
              onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
              required
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:border-transparent text-sm"
              style={{ '--tw-ring-color': '#1B76FF' } as React.CSSProperties}
              placeholder="Введите ваше имя"
            />
          </div>

          {/* Bio */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              О себе
            </label>
            <textarea
              value={editForm.bio}
              onChange={(e) => setEditForm({ ...editForm, bio: e.target.value })}
              rows={4}
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:border-transparent resize-none text-sm"
              style={{ '--tw-ring-color': '#1B76FF' } as React.CSSProperties}
              placeholder="Расскажите о себе..."
            />
            <p className="text-xs text-gray-500 mt-1">
              {editForm.bio.length} / 500 символов
            </p>
          </div>

          {/* Phone */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Телефон
            </label>
            <input
              type="tel"
              value={editForm.phone}
              onChange={(e) => setEditForm({ ...editForm, phone: e.target.value })}
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:border-transparent text-sm"
              style={{ '--tw-ring-color': '#1B76FF' } as React.CSSProperties}
              placeholder="+7 (999) 123-45-67"
            />
          </div>

          {/* Location */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Город
            </label>
            <input
              type="text"
              value={editForm.location}
              onChange={(e) => setEditForm({ ...editForm, location: e.target.value })}
              className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:border-transparent text-sm"
              style={{ '--tw-ring-color': '#1B76FF' } as React.CSSProperties}
              placeholder="Москва"
            />
          </div>

          {/* Actions */}
          <div className="flex gap-3 pt-4 border-t border-gray-200">
            <button
              type="button"
              onClick={() => router.back()}
              className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-50 transition-colors text-sm"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isSaving || !editForm.name.trim()}
              className="flex-1 px-4 py-2.5 rounded-lg text-white font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
              style={{ backgroundColor: '#1B76FF' }}
            >
              {isSaving ? 'Сохранение...' : 'Сохранить изменения'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
