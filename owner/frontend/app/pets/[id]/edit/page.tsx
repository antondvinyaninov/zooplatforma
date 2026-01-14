'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import AdminLayout, { AdminTab } from '../../../components/admin/AdminLayout';
import { PencilIcon } from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'edit', label: 'Редактирование', icon: <PencilIcon className="w-5 h-5" /> },
];

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  gender?: string;
  birth_date?: string;
  color?: string;
  size?: string;
  weight?: number;
  chip_number?: string;
  is_sterilized: boolean;
  is_vaccinated: boolean;
  health_notes?: string;
  character_traits?: string;
  special_needs?: string;
  photo?: string;
}

export default function EditPetPage() {
  const params = useParams();
  const router = useRouter();
  const [pet, setPet] = useState<Pet | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [user, setUser] = useState<{ email: string; name?: string; role: string } | null>(null);
  const [activeTab, setActiveTab] = useState('edit');

  useEffect(() => {
    loadData();
  }, [params.id]);

  const loadData = async () => {
    try {
      // Проверяем авторизацию
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        credentials: 'include',
      });

      if (!meResponse.ok) {
        router.push('/auth');
        return;
      }

      const meResult = await meResponse.json();
      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      setUser({
        email: meResult.data.email,
        name: meResult.data.name,
        role: 'owner',
      });

      // Загружаем данные питомца из PetBase
      const petResponse = await fetch(`http://localhost:8100/api/pets/${params.id}`, {
        credentials: 'include',
      });

      if (petResponse.ok) {
        const petResult = await petResponse.json();
        if (petResult.success && petResult.data) {
          setPet(petResult.data);
        }
      }
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!pet) return;

    setSaving(true);

    try {
      const response = await fetch(`http://localhost:8100/api/pets/${params.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(pet),
      });

      if (response.ok) {
        alert('Питомец успешно обновлён!');
        router.push('/pets');
      } else {
        alert('Ошибка при сохранении');
      }
    } catch (error) {
      console.error('Error saving pet:', error);
      alert('Ошибка при сохранении');
    } finally {
      setSaving(false);
    }
  };

  const handleLogout = async () => {
    await fetch('http://localhost:8000/api/auth/logout', {
      method: 'POST',
      credentials: 'include',
    });
    router.push('/auth');
  };

  if (loading) {
    return (
      <AdminLayout
        logoSrc="/favicon.svg"
        logoText="Кабинет владельца"
        logoAlt="Кабинет владельца"
        tabs={tabs}
        activeTab={activeTab}
        onTabChange={setActiveTab}
        adminUser={user}
        onLogout={handleLogout}
        mainSiteUrl="http://localhost:3000"
      >
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </AdminLayout>
    );
  }

  if (!pet) {
    return (
      <AdminLayout
        logoSrc="/favicon.svg"
        logoText="Кабинет владельца"
        logoAlt="Кабинет владельца"
        tabs={tabs}
        activeTab={activeTab}
        onTabChange={setActiveTab}
        adminUser={user}
        onLogout={handleLogout}
        mainSiteUrl="http://localhost:3000"
      >
        <div className="text-center py-12">
          <p className="text-gray-600">Питомец не найден</p>
          <button
            onClick={() => router.push('/pets')}
            className="mt-4 text-blue-600 hover:text-blue-700"
          >
            Вернуться к списку
          </button>
        </div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout
      logoSrc="/favicon.svg"
      logoText="Кабинет владельца"
      logoAlt="Кабинет владельца"
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      adminUser={user}
      onLogout={handleLogout}
      mainSiteUrl="http://localhost:3000"
    >
      <div className="space-y-6">
        <div>
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Редактирование: {pet.name}</h2>
          <p className="text-gray-600">Обновите информацию о вашем питомце</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
        {/* Основная информация */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Основная информация</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Кличка *
              </label>
              <input
                type="text"
                required
                value={pet.name}
                onChange={(e) => setPet({ ...pet, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Вид *
              </label>
              <select
                required
                value={pet.species}
                onChange={(e) => setPet({ ...pet, species: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="dog">Собака</option>
                <option value="cat">Кошка</option>
                <option value="bird">Птица</option>
                <option value="rodent">Грызун</option>
                <option value="reptile">Рептилия</option>
                <option value="fish">Рыба</option>
                <option value="other">Другое</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Порода
              </label>
              <input
                type="text"
                value={pet.breed || ''}
                onChange={(e) => setPet({ ...pet, breed: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Пол
              </label>
              <select
                value={pet.gender || ''}
                onChange={(e) => setPet({ ...pet, gender: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Не указан</option>
                <option value="male">Самец</option>
                <option value="female">Самка</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Дата рождения
              </label>
              <input
                type="date"
                value={pet.birth_date || ''}
                onChange={(e) => setPet({ ...pet, birth_date: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Окрас
              </label>
              <input
                type="text"
                value={pet.color || ''}
                onChange={(e) => setPet({ ...pet, color: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Размер
              </label>
              <select
                value={pet.size || ''}
                onChange={(e) => setPet({ ...pet, size: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Не указан</option>
                <option value="small">Маленький</option>
                <option value="medium">Средний</option>
                <option value="large">Большой</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Вес (кг)
              </label>
              <input
                type="number"
                step="0.1"
                value={pet.weight || ''}
                onChange={(e) => setPet({ ...pet, weight: parseFloat(e.target.value) || undefined })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Номер чипа
              </label>
              <input
                type="text"
                value={pet.chip_number || ''}
                onChange={(e) => setPet({ ...pet, chip_number: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        {/* Медицинская информация */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Медицинская информация</h2>
          
          <div className="space-y-4">
            <div className="flex items-center gap-4">
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={pet.is_sterilized}
                  onChange={(e) => setPet({ ...pet, is_sterilized: e.target.checked })}
                  className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span className="text-sm font-medium text-gray-700">Стерилизован</span>
              </label>

              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={pet.is_vaccinated}
                  onChange={(e) => setPet({ ...pet, is_vaccinated: e.target.checked })}
                  className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span className="text-sm font-medium text-gray-700">Привит</span>
              </label>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Заметки о здоровье
              </label>
              <textarea
                value={pet.health_notes || ''}
                onChange={(e) => setPet({ ...pet, health_notes: e.target.value })}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Особенности характера
              </label>
              <textarea
                value={pet.character_traits || ''}
                onChange={(e) => setPet({ ...pet, character_traits: e.target.value })}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Особые потребности
              </label>
              <textarea
                value={pet.special_needs || ''}
                onChange={(e) => setPet({ ...pet, special_needs: e.target.value })}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        {/* Кнопки */}
        <div className="flex gap-3">
          <button
            type="button"
            onClick={() => router.push('/pets')}
            className="flex-1 px-4 py-2 border-2 border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
          >
            Отмена
          </button>
          <button
            type="submit"
            disabled={saving}
            className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {saving ? 'Сохранение...' : 'Сохранить изменения'}
          </button>
        </div>
      </form>
      </div>
    </AdminLayout>
  );
}
