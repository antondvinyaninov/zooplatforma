'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { PlusIcon } from '@heroicons/react/24/outline';

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  photo?: string;
  created_at: string;
}

export default function PetsListPage() {
  const router = useRouter();
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPets();
  }, []);

  const loadPets = async () => {
    try {
      const meResponse = await fetch('http://localhost:8000/api/auth/me', {
        credentials: 'include',
      });

      if (meResponse.ok) {
        const meResult = await meResponse.json();
        const userId = meResult.data.id;

        const petsResponse = await fetch(`http://localhost:8100/api/pets?user_id=${userId}`, {
          headers: {
            'X-User-ID': userId.toString(),
          },
          credentials: 'include',
        });

        if (petsResponse.ok) {
          const petsResult = await petsResponse.json();
          if (petsResult.success && petsResult.data) {
            setPets(petsResult.data);
          }
        }
      }
    } catch (error) {
      console.error('Error loading pets:', error);
    } finally {
      setLoading(false);
    }
  };

  const getSpeciesLabel = (species: string) => {
    const labels: Record<string, string> = {
      dog: 'Собака',
      cat: 'Кошка',
      bird: 'Птица',
      rodent: 'Грызун',
      reptile: 'Рептилия',
      fish: 'Рыба',
      other: 'Другое',
    };
    return labels[species] || species;
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Мои питомцы</h2>
          <p className="text-gray-600">Управление вашими любимцами</p>
        </div>
        <button
          onClick={() => window.open('http://localhost:3000/pets', '_blank')}
          className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
        >
          <PlusIcon className="w-4 h-4" />
          Добавить питомца
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        {loading ? (
          <div className="text-center py-12 text-gray-400">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4"></div>
            Загрузка...
          </div>
        ) : pets.length === 0 ? (
          <div className="text-center py-12 text-gray-400">
            <div className="text-4xl mb-3">🐾</div>
            <p>Питомцы не найдены</p>
            <p className="text-sm mt-2">Добавьте первого питомца на основном сайте</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="text-left py-3 px-6 text-sm font-semibold text-gray-700">Фото</th>
                  <th className="text-left py-3 px-6 text-sm font-semibold text-gray-700">Кличка</th>
                  <th className="text-left py-3 px-6 text-sm font-semibold text-gray-700">Вид</th>
                  <th className="text-left py-3 px-6 text-sm font-semibold text-gray-700">Порода</th>
                  <th className="text-left py-3 px-6 text-sm font-semibold text-gray-700">Добавлен</th>
                  <th className="text-left py-3 px-6 text-sm font-semibold text-gray-700">Действия</th>
                </tr>
              </thead>
              <tbody>
                {pets.map((pet) => (
                  <tr key={pet.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                    <td className="py-4 px-6">
                      {pet.photo ? (
                        <img
                          src={pet.photo}
                          alt={pet.name}
                          className="w-12 h-12 rounded-full object-cover"
                        />
                      ) : (
                        <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-xl">
                          🐾
                        </div>
                      )}
                    </td>
                    <td className="py-4 px-6 text-sm font-medium text-gray-900">{pet.name}</td>
                    <td className="py-4 px-6 text-sm text-gray-600">{getSpeciesLabel(pet.species)}</td>
                    <td className="py-4 px-6 text-sm text-gray-600">{pet.breed || '-'}</td>
                    <td className="py-4 px-6 text-sm text-gray-600">
                      {new Date(pet.created_at).toLocaleDateString('ru-RU')}
                    </td>
                    <td className="py-4 px-6 text-sm">
                      <div className="flex gap-3">
                        <button
                          onClick={() => router.push(`/pets/${pet.id}`)}
                          className="text-blue-600 hover:text-blue-700 font-medium"
                        >
                          Просмотр
                        </button>
                        <button
                          onClick={() => router.push(`/pets/${pet.id}/edit`)}
                          className="text-gray-600 hover:text-gray-700"
                        >
                          Редактировать
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
