'use client';

import { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../../../contexts/AuthContext';
import { usersApi } from '../../../../lib/api';
import { UserIcon, CameraIcon, XMarkIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';

export default function EditProfilePage() {
  const [isSaving, setIsSaving] = useState(false);
  const [isUploadingAvatar, setIsUploadingAvatar] = useState(false);
  const [isUploadingCover, setIsUploadingCover] = useState(false);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [coverPreview, setCoverPreview] = useState<string | null>(null);
  const avatarInputRef = useRef<HTMLInputElement>(null);
  const coverInputRef = useRef<HTMLInputElement>(null);
  const [editForm, setEditForm] = useState({
    name: '',
    bio: '',
    phone: '',
    location: '',
  });
  const { user, isAuthenticated, isLoading, refreshUser } = useAuth();
  const router = useRouter();

  useEffect(() => {
    // Проверяем только на клиенте
    if (typeof window === 'undefined') return;

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
      setAvatarPreview(user.avatar ? `http://localhost:8000${user.avatar}` : null);
      setCoverPreview(user.cover_photo ? `http://localhost:8000${user.cover_photo}` : null);
    }
  }, [user, isLoading, isAuthenticated, router]);

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Проверка типа файла
    if (!file.type.startsWith('image/')) {
      alert('Пожалуйста, выберите изображение');
      return;
    }

    // Проверка размера (максимум 10MB)
    if (file.size > 10 * 1024 * 1024) {
      alert('Файл слишком большой (максимум 10MB)');
      return;
    }

    // Показываем превью
    const reader = new FileReader();
    reader.onloadend = () => {
      setAvatarPreview(reader.result as string);
    };
    reader.readAsDataURL(file);

    // Загружаем на сервер
    setIsUploadingAvatar(true);
    try {
      const response = await usersApi.uploadAvatar(file);
      
      if (response.success && response.data) {
        // Обновляем данные пользователя
        await refreshUser();
        // Устанавливаем правильный URL с сервера
        setAvatarPreview(`http://localhost:8000${response.data.avatar_url}`);
        alert('Аватар успешно обновлен!');
      } else {
        alert(response.error || 'Ошибка загрузки аватара');
        // Возвращаем старый аватар
        setAvatarPreview(user?.avatar ? `http://localhost:8000${user.avatar}` : null);
      }
    } catch (error) {
      console.error('Ошибка загрузки аватара:', error);
      alert('Ошибка загрузки аватара');
      setAvatarPreview(user?.avatar ? `http://localhost:8000${user.avatar}` : null);
    } finally {
      setIsUploadingAvatar(false);
    }
  };

  const handleCoverChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Проверка типа файла
    if (!file.type.startsWith('image/')) {
      alert('Пожалуйста, выберите изображение');
      return;
    }

    // Проверка размера (максимум 10MB)
    if (file.size > 10 * 1024 * 1024) {
      alert('Файл слишком большой (максимум 10MB)');
      return;
    }

    // Показываем превью
    const reader = new FileReader();
    reader.onloadend = () => {
      setCoverPreview(reader.result as string);
    };
    reader.readAsDataURL(file);

    // Загружаем на сервер
    setIsUploadingCover(true);
    try {
      const response = await usersApi.uploadCover(file);
      
      if (response.success && response.data) {
        // Обновляем данные пользователя
        await refreshUser();
        // Устанавливаем правильный URL с сервера
        setCoverPreview(`http://localhost:8000${response.data.cover_url}`);
        alert('Обложка успешно обновлена!');
      } else {
        alert(response.error || 'Ошибка загрузки обложки');
        // Возвращаем старую обложку
        setCoverPreview(user?.cover_photo ? `http://localhost:8000${user.cover_photo}` : null);
      }
    } catch (error) {
      console.error('Ошибка загрузки обложки:', error);
      alert('Ошибка загрузки обложки');
      setCoverPreview(user?.cover_photo ? `http://localhost:8000${user.cover_photo}` : null);
    } finally {
      setIsUploadingCover(false);
    }
  };

  const handleDeleteAvatar = async () => {
    if (!confirm('Вы уверены, что хотите удалить аватар?')) return;

    setIsUploadingAvatar(true);
    try {
      const response = await usersApi.deleteAvatar();
      
      if (response.success) {
        await refreshUser();
        setAvatarPreview(null);
        alert('Аватар успешно удален!');
      } else {
        alert(response.error || 'Ошибка удаления аватара');
      }
    } catch (error) {
      console.error('Ошибка удаления аватара:', error);
      alert('Ошибка удаления аватара');
    } finally {
      setIsUploadingAvatar(false);
    }
  };

  const handleDeleteCover = async () => {
    if (!confirm('Вы уверены, что хотите удалить обложку?')) return;

    setIsUploadingCover(true);
    try {
      const response = await usersApi.deleteCover();
      
      if (response.success) {
        await refreshUser();
        setCoverPreview(null);
        alert('Обложка успешно удалена!');
      } else {
        alert(response.error || 'Ошибка удаления обложки');
      }
    } catch (error) {
      console.error('Ошибка удаления обложки:', error);
      alert('Ошибка удаления обложки');
    } finally {
      setIsUploadingCover(false);
    }
  };

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
          {/* Cover Photo Section */}
          {coverPreview ? (
            <div className="relative">
              <div className="relative h-48 bg-gradient-to-br from-blue-100 to-blue-200 rounded-lg overflow-hidden">
                <img src={coverPreview} alt="Cover" className="w-full h-full object-cover" />
                {isUploadingCover && (
                  <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
                  </div>
                )}
              </div>
              <input
                ref={coverInputRef}
                type="file"
                accept="image/*"
                onChange={handleCoverChange}
                className="hidden"
              />
              <button
                type="button"
                onClick={() => coverInputRef.current?.click()}
                disabled={isUploadingCover}
                className="absolute top-4 right-4 p-2 bg-white rounded-lg shadow-md hover:bg-gray-50 transition-colors disabled:opacity-50"
              >
                <CameraIcon className="w-5 h-5 text-gray-600" />
              </button>
              <button
                type="button"
                onClick={handleDeleteCover}
                disabled={isUploadingCover}
                className="absolute top-4 right-16 p-2 bg-red-500 rounded-lg shadow-md hover:bg-red-600 transition-colors disabled:opacity-50"
                title="Удалить обложку"
              >
                <XMarkIcon className="w-5 h-5 text-white" />
              </button>
            </div>
          ) : (
            <div className="relative">
              <input
                ref={coverInputRef}
                type="file"
                accept="image/*"
                onChange={handleCoverChange}
                className="hidden"
              />
              <button
                type="button"
                onClick={() => coverInputRef.current?.click()}
                disabled={isUploadingCover}
                className="w-full py-3 px-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-gray-400 transition-colors disabled:opacity-50 flex items-center justify-center gap-2 text-gray-600"
              >
                <CameraIcon className="w-5 h-5" />
                <span className="text-sm font-medium">
                  {isUploadingCover ? 'Загрузка...' : 'Добавить обложку профиля'}
                </span>
              </button>
            </div>
          )}

          {/* Avatar Section */}
          <div className="flex items-center gap-6">
            <div className="relative">
              <div className="w-24 h-24 rounded-full bg-gray-300 flex items-center justify-center overflow-hidden">
                {avatarPreview ? (
                  <img src={avatarPreview} alt={user?.name} className="w-full h-full object-cover" />
                ) : (
                  <UserIcon className="w-12 h-12 text-gray-500" />
                )}
                {isUploadingAvatar && (
                  <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-white"></div>
                  </div>
                )}
              </div>
              <input
                ref={avatarInputRef}
                type="file"
                accept="image/*"
                onChange={handleAvatarChange}
                className="hidden"
              />
              <button
                type="button"
                onClick={() => avatarInputRef.current?.click()}
                disabled={isUploadingAvatar}
                className="absolute bottom-0 right-0 p-2 bg-white rounded-full shadow-md hover:bg-gray-50 transition-colors border border-gray-200 disabled:opacity-50"
                title="Изменить фото"
              >
                <CameraIcon className="w-4 h-4 text-gray-600" />
              </button>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900">{user?.name}</h3>
              <p className="text-sm text-gray-600">{user?.email}</p>
              <div className="flex gap-2 mt-2">
                <button
                  type="button"
                  onClick={() => avatarInputRef.current?.click()}
                  disabled={isUploadingAvatar}
                  className="text-sm font-medium disabled:opacity-50"
                  style={{ color: '#1B76FF' }}
                >
                  {isUploadingAvatar ? 'Загрузка...' : 'Изменить фото профиля'}
                </button>
                {avatarPreview && (
                  <>
                    <span className="text-gray-300">|</span>
                    <button
                      type="button"
                      onClick={handleDeleteAvatar}
                      disabled={isUploadingAvatar}
                      className="text-sm font-medium text-red-500 hover:text-red-600 disabled:opacity-50"
                    >
                      Удалить
                    </button>
                  </>
                )}
              </div>
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
