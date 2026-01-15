/**
 * PetCard Component
 * 
 * Отображает карточку питомца с фото, основной информацией и статусом
 */

import { Pet } from '../types/owner';

interface PetCardProps {
  pet: Pet;
  onClick: () => void;
}

const SPECIES_EMOJI: Record<string, string> = {
  dog: '🐕',
  cat: '🐈',
  bird: '🦜',
  rodent: '🐹',
  reptile: '🦎',
  fish: '🐠',
  other: '🐾',
};

const SPECIES_LABELS: Record<string, string> = {
  dog: 'Собака',
  cat: 'Кошка',
  bird: 'Птица',
  rodent: 'Грызун',
  reptile: 'Рептилия',
  fish: 'Рыба',
  other: 'Другое',
};

const STATUS_LABELS: Record<string, string> = {
  home: 'Дома',
  lost: 'Потерялся',
  found: 'Найден',
  looking_for_home: 'Ищет дом',
  needs_help: 'Нужна помощь',
  at_vet: 'У ветеринара',
  died: 'Умер',
};

const STATUS_COLORS: Record<string, string> = {
  home: 'bg-green-100 text-green-800',
  lost: 'bg-red-100 text-red-800',
  found: 'bg-blue-100 text-blue-800',
  looking_for_home: 'bg-yellow-100 text-yellow-800',
  needs_help: 'bg-orange-100 text-orange-800',
  at_vet: 'bg-purple-100 text-purple-800',
  died: 'bg-gray-100 text-gray-800',
};

export default function PetCard({ pet, onClick }: PetCardProps) {
  const speciesEmoji = SPECIES_EMOJI[pet.species] || '🐾';
  const speciesLabel = SPECIES_LABELS[pet.species] || pet.species;
  const statusLabel = STATUS_LABELS[pet.status] || pet.status;
  const statusColor = STATUS_COLORS[pet.status] || 'bg-gray-100 text-gray-800';

  return (
    <div
      onClick={onClick}
      className="bg-white rounded-xl shadow-sm hover:shadow-md transition-all duration-200 cursor-pointer overflow-hidden group hover:scale-105"
    >
      {/* Фото питомца */}
      <div className="relative h-48 bg-gradient-to-br from-blue-50 to-purple-50">
        {pet.photo_url ? (
          <img
            src={`http://localhost:8400${pet.photo_url}`}
            alt={pet.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-6xl">
            {speciesEmoji}
          </div>
        )}
        
        {/* Статус badge */}
        <div className="absolute top-3 right-3">
          <span className={`px-3 py-1 rounded-full text-xs font-medium ${statusColor}`}>
            {statusLabel}
          </span>
        </div>

        {/* Verification status */}
        {pet.verification_status === 'pending_verification' && (
          <div className="absolute top-3 left-3">
            <span className="px-3 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
              ⏳ Ожидает подтверждения
            </span>
          </div>
        )}
      </div>

      {/* Информация о питомце */}
      <div className="p-4">
        <h3 className="text-xl font-bold text-gray-900 mb-2 group-hover:text-blue-600 transition-colors">
          {pet.name}
        </h3>
        
        <div className="space-y-1 text-sm text-gray-600">
          <div className="flex items-center gap-2">
            <span className="text-lg">{speciesEmoji}</span>
            <span>{speciesLabel}</span>
            {pet.breed && <span className="text-gray-400">• {pet.breed}</span>}
          </div>
          
          <div className="flex items-center gap-2">
            <span>🎂</span>
            <span>{pet.age} {pet.age === 1 ? 'год' : pet.age < 5 ? 'года' : 'лет'}</span>
            <span className="text-gray-400">•</span>
            <span>{pet.sex === 'male' ? '♂️ Самец' : '♀️ Самка'}</span>
          </div>

          {pet.color && (
            <div className="flex items-center gap-2">
              <span>🎨</span>
              <span>{pet.color}</span>
            </div>
          )}
        </div>

        {/* Дополнительные бейджи */}
        <div className="flex flex-wrap gap-2 mt-3">
          {pet.sterilized && (
            <span className="px-2 py-1 bg-blue-50 text-blue-700 rounded-md text-xs">
              ✂️ Стерилизован
            </span>
          )}
          {pet.chip_number && (
            <span className="px-2 py-1 bg-purple-50 text-purple-700 rounded-md text-xs">
              💎 Чипирован
            </span>
          )}
          {pet.urgent && (
            <span className="px-2 py-1 bg-red-50 text-red-700 rounded-md text-xs font-medium">
              🚨 СРОЧНО
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
