'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  photo?: string;
  created_at: string;
}

export default function PetsPage() {
  const router = useRouter();
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    checkAuthAndLoadPets();
  }, []);

  const checkAuthAndLoadPets = async () => {
    try {
      // Проверяем авторизацию
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        method: 'GET',
        credentials: 'include',
      });

      const meResult = await meResponse.json();

      if (!meResult.success) {
        router.push('/auth');
        return;
      }

      setUser(meResult.data);

      // Загружаем питомцев
      const petsResponse = await fetch('http://localhost:8400/api/my-pets', {
        method: 'GET',
        credentials: 'include',
      });

      const petsResult = await petsResponse.json();

      if (petsResult.success && petsResult.data) {
        setPets(petsResult.data);
      }
    } catch (error) {
      console.error('Failed to load pets:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg text-gray-600">Загрузка...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                Мои питомцы
              </h1>
              <p className="text-sm text-gray-600 mt-1">
                {user?.name || 'Владелец'}
              </p>
            </div>
            <button
              onClick={() => window.location.href = 'http://localhost:3000'}
              className="text-blue-600 hover:text-blue-700"
            >
              На главную
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {pets.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center">
            <p className="text-gray-600 mb-4">
              У вас пока нет питомцев
            </p>
            <button
              onClick={() => window.location.href = 'http://localhost:3000/pets'}
              className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700"
            >
              Добавить питомца
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {pets.map((pet) => (
              <div
                key={pet.id}
                className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow cursor-pointer"
                onClick={() => router.push(`/pets/${pet.id}`)}
              >
                {pet.photo && (
                  <img
                    src={pet.photo}
                    alt={pet.name}
                    className="w-full h-48 object-cover rounded-t-lg"
                  />
                )}
                <div className="p-4">
                  <h3 className="text-lg font-semibold text-gray-900">
                    {pet.name}
                  </h3>
                  <p className="text-sm text-gray-600">
                    {pet.species} {pet.breed && `• ${pet.breed}`}
                  </p>
                  <p className="text-xs text-gray-400 mt-2">
                    Добавлен: {new Date(pet.created_at).toLocaleDateString('ru-RU')}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
