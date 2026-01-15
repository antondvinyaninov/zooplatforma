/**
 * PetList Component
 * 
 * Отображает сетку карточек питомцев владельца
 */

'use client';

import { PlusIcon } from '@heroicons/react/24/outline';
import { Pet } from '../types/owner';
import PetCard from './PetCard';

interface PetListProps {
  pets: Pet[];
  onPetClick: (petId: number) => void;
  onAddPet: () => void;
}

export default function PetList({ pets, onPetClick, onAddPet }: PetListProps) {
  return (
    <div className="space-y-6">
      {/* Заголовок и кнопка добавления */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Мои питомцы</h2>
          <p className="text-gray-600">
            {pets.length === 0 
              ? 'У вас пока нет питомцев' 
              : `Всего питомцев: ${pets.length}`
            }
          </p>
        </div>
        <button
          onClick={onAddPet}
          className="px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white text-sm rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all duration-200 flex items-center gap-2 shadow-md hover:shadow-lg"
        >
          <PlusIcon className="w-5 h-5" />
          Добавить питомца
        </button>
      </div>

      {/* Список питомцев */}
      {pets.length === 0 ? (
        <div className="bg-white rounded-xl shadow-sm p-12 text-center">
          <div className="text-6xl mb-4">🐾</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            У вас пока нет питомцев
          </h3>
          <p className="text-gray-600 mb-6">
            Добавьте первого питомца, чтобы начать управлять их информацией
          </p>
          <button
            onClick={onAddPet}
            className="px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all duration-200 inline-flex items-center gap-2 shadow-md hover:shadow-lg"
          >
            <PlusIcon className="w-5 h-5" />
            Добавить первого питомца
          </button>
        </div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
          <div className="divide-y divide-gray-200">
            {pets.map(pet => (
              <div
                key={pet.id}
                onClick={() => onPetClick(pet.id)}
                className="p-4 hover:bg-gray-50 transition-colors cursor-pointer flex items-center gap-4"
              >
                {/* Фото питомца */}
                <div className="flex-shrink-0">
                  {pet.photo_url ? (
                    <img
                      src={`http://localhost:8400${pet.photo_url}`}
                      alt={pet.name}
                      className="w-16 h-16 rounded-full object-cover"
                    />
                  ) : (
                    <div className="w-16 h-16 rounded-full bg-gradient-to-br from-blue-50 to-purple-50 flex items-center justify-center text-3xl">
                      {pet.species === 'dog' ? '🐕' : pet.species === 'cat' ? '🐈' : '🐾'}
                    </div>
                  )}
                </div>

                {/* Информация о питомце */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="text-lg font-semibold text-gray-900 truncate">
                      {pet.name}
                    </h3>
                    {pet.verification_status === 'pending_verification' && (
                      <span className="px-2 py-0.5 bg-yellow-100 text-yellow-800 text-xs rounded-full">
                        На проверке
                      </span>
                    )}
                  </div>
                  <div className="flex items-center gap-3 text-sm text-gray-600">
                    <span>{pet.breed || 'Порода не указана'}</span>
                    <span>•</span>
                    <span>{pet.age} {pet.age === 1 ? 'год' : pet.age < 5 ? 'года' : 'лет'}</span>
                    <span>•</span>
                    <span>{pet.sex === 'male' ? '♂️ Самец' : '♀️ Самка'}</span>
                  </div>
                </div>

                {/* Badges */}
                <div className="flex items-center gap-2">
                  {pet.sterilized && (
                    <span className="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full" title="Стерилизован">
                      ✂️
                    </span>
                  )}
                  {pet.chip_number && (
                    <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full" title="Чипирован">
                      💎
                    </span>
                  )}
                </div>

                {/* Стрелка */}
                <div className="flex-shrink-0 text-gray-400">
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
