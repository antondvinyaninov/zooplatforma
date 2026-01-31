'use client';

import { useState } from 'react';
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import CompleteCardForm from './CompleteCardForm';

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  birth_date?: string;
  gender?: string;
  photo?: string;
  chip_number?: string;
  status: string;
  owner_name?: string;
  owner_phone?: string;
  owner_email?: string;
  owner_id?: number;
  created_at: string;
}

interface Owner {
  id: number;
  name: string;
  last_name?: string;
  email?: string;
  phone?: string;
  avatar?: string;
}

interface PetsWithOwnerResponse {
  owner: Owner;
  pets: Pet[];
}

export default function SearchExistingCard() {
  const [searchType, setSearchType] = useState<'owner-id'>('owner-id');
  const [searchQuery, setSearchQuery] = useState('');
  const [ownerData, setOwnerData] = useState<PetsWithOwnerResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [selectedPet, setSelectedPet] = useState<Pet | null>(null);

  const handleSearch = async () => {
    if (!searchQuery.trim()) return;

    setLoading(true);
    setErrorMessage(null);
    setOwnerData(null);
    
    try {
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        credentials: 'include',
      });

      if (!meResponse.ok) {
        throw new Error('Failed to load auth data');
      }

      const meData = await meResponse.json();
      const token = meData.data?.token;

      if (!meData.success || !token) {
        throw new Error('Missing auth token');
      }

      // Парсим ID владельца
      const rawValue = searchQuery.trim();
      const ownerId = rawValue.startsWith('id') ? rawValue.slice(2) : rawValue;

      // Один запрос к PetBase - получаем владельца + питомцев
      const response = await fetch(
        `http://localhost:8100/api/pets/user/${ownerId}?include_owner=true`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to load data');
      }

      const data = await response.json();
      if (data.success && data.data) {
        setOwnerData(data.data);
      } else {
        setErrorMessage('Владелец не найден');
      }
    } catch (error) {
      console.error('Search failed:', error);
      setErrorMessage('Не удалось выполнить поиск. Проверьте доступность сервиса.');
    }
    setLoading(false);
  };

  if (selectedPet) {
    return (
      <div>
        <button
          onClick={() => {
            setSelectedPet(null);
            setOwnerData(null);
          }}
          className="mb-6 text-gray-600 hover:text-gray-900 transition-colors flex items-center space-x-2"
        >
          <span>←</span>
          <span>Назад к результатам</span>
        </button>
        <CompleteCardForm pet={selectedPet} onComplete={() => {
          setSelectedPet(null);
          setOwnerData(null);
          setSearchQuery('');
        }} />
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">
        Поиск существующей карточки
      </h3>

      {/* Шаг 1: Поиск владельца по ID */}
      {!ownerData && (
        <>
          {/* Поле поиска */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              ID владельца
            </label>
            <div className="flex space-x-4">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="id123 или 123"
                className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <button
                onClick={handleSearch}
                disabled={loading || !searchQuery.trim()}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <MagnifyingGlassIcon className="w-5 h-5" />
                <span>{loading ? 'Поиск...' : 'Найти'}</span>
              </button>
            </div>
          </div>

          {errorMessage && (
            <div className="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
              {errorMessage}
            </div>
          )}
        </>
      )}

      {/* Шаг 2: Информация о владельце + список питомцев */}
      {ownerData && (
        <div>
          <div className="flex items-center justify-between mb-4">
            <h4 className="text-lg font-semibold text-gray-900">Владелец и его питомцы</h4>
            <button
              onClick={() => {
                setOwnerData(null);
                setSearchQuery('');
                setErrorMessage(null);
              }}
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              ← Новый поиск
            </button>
          </div>

          {/* Карточка владельца */}
          <div className="border border-gray-200 rounded-xl p-6 bg-gray-50 mb-6">
            <div className="flex items-start space-x-4">
              {/* Аватар */}
              <div className="w-20 h-20 bg-gray-200 rounded-full overflow-hidden flex-shrink-0">
                {ownerData.owner.avatar ? (
                  <img
                    src={`http://localhost:8000${ownerData.owner.avatar}`}
                    alt={ownerData.owner.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-3xl">
                    👤
                  </div>
                )}
              </div>

              {/* Информация */}
              <div className="flex-1">
                <h5 className="text-xl font-bold text-gray-900 mb-2">
                  {ownerData.owner.name} {ownerData.owner.last_name}
                </h5>
                {ownerData.owner.email && (
                  <p className="text-gray-600 mb-1">
                    📧 {ownerData.owner.email}
                  </p>
                )}
                {ownerData.owner.phone && (
                  <p className="text-gray-600 mb-1">
                    📱 {ownerData.owner.phone}
                  </p>
                )}
                <p className="text-sm text-gray-500 mt-2">
                  ID: {ownerData.owner.id}
                </p>
              </div>
            </div>
          </div>

          {/* Список питомцев */}
          {ownerData.pets.length > 0 ? (
            <div>
              <h5 className="text-lg font-semibold text-gray-900 mb-4">
                Питомцы ({ownerData.pets.length})
              </h5>
              <div className="space-y-4">
                {ownerData.pets.map((pet) => (
                  <div
                    key={pet.id}
                    className="border border-gray-200 rounded-xl p-6 hover:border-blue-300 hover:shadow-md transition-all cursor-pointer"
                    onClick={() => setSelectedPet(pet)}
                  >
                    <div className="flex items-start space-x-4">
                      {/* Фото питомца */}
                      <div className="w-20 h-20 bg-gray-100 rounded-xl overflow-hidden flex-shrink-0">
                        {pet.photo ? (
                          <img
                            src={`http://localhost:8000${pet.photo}`}
                            alt={pet.name}
                            className="w-full h-full object-cover"
                          />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center text-3xl">
                            🐾
                          </div>
                        )}
                      </div>

                      {/* Информация */}
                      <div className="flex-1">
                        <div className="flex items-start justify-between mb-2">
                          <div>
                            <h5 className="text-xl font-bold text-gray-900">{pet.name}</h5>
                            <p className="text-gray-600">
                              {pet.species}
                              {pet.breed && ` • ${pet.breed}`}
                              {pet.gender && ` • ${pet.gender === 'male' ? 'Самец' : 'Самка'}`}
                            </p>
                          </div>
                          <span className={`px-3 py-1 text-sm font-medium rounded-full ${
                            pet.status === 'pending_verification'
                              ? 'bg-yellow-100 text-yellow-800'
                              : 'bg-green-100 text-green-800'
                          }`}>
                            {pet.status === 'pending_verification' ? 'Ожидает подтверждения' : 'Подтверждено'}
                          </span>
                        </div>

                        {/* Chip number */}
                        {pet.chip_number && (
                          <div className="mb-2">
                            <p className="text-sm text-gray-500">Chip number:</p>
                            <p className="text-gray-900 font-mono">{pet.chip_number}</p>
                          </div>
                        )}

                        {/* Дата создания */}
                        <p className="text-sm text-gray-500">
                          Создано: {new Date(pet.created_at).toLocaleDateString('ru-RU')}
                        </p>
                      </div>

                      {/* Кнопка */}
                      <div className="flex-shrink-0">
                        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                          Дополнить →
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="text-center py-12">
              <div className="text-6xl mb-4">🐾</div>
              <h4 className="text-xl font-semibold text-gray-900 mb-2">
                У этого владельца пока нет питомцев
              </h4>
              <p className="text-gray-600 mb-6">
                Вы можете создать новую карточку для этого владельца
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
