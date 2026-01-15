/**
 * AddPetModal Component
 * 
 * Упрощённая форма добавления питомца (только кошки и собаки)
 */

'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { CreatePetData } from '../types/owner';
import { petidApi, Breed } from '../lib/petid-api';

interface AddPetModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreatePetData) => Promise<void>;
}

export default function AddPetModal({ isOpen, onClose, onSubmit }: AddPetModalProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [breeds, setBreeds] = useState<Breed[]>([]);
  const [filteredBreeds, setFilteredBreeds] = useState<Breed[]>([]);
  const [showBreedDropdown, setShowBreedDropdown] = useState(false);
  const [ageInputType, setAgeInputType] = useState<'age' | 'birthdate'>('age'); // Тип ввода возраста
  const [birthdate, setBirthdate] = useState<string>(''); // Дата рождения
  
  const [formData, setFormData] = useState<CreatePetData>({
    name: '',
    species: 'dog',
    breed: '',
    age: 0,
    sex: 'male',
    color: '',
    weight: undefined,
    chip_number: '',
    sterilized: false,
  });

  // Загружаем породы при изменении вида животного
  useEffect(() => {
    const loadBreeds = async () => {
      const response = await petidApi.getBreeds(formData.species as 'dog' | 'cat');
      if (response.data) {
        setBreeds(response.data);
        setFilteredBreeds(response.data);
      }
    };
    
    if (isOpen) {
      loadBreeds();
    }
  }, [formData.species, isOpen]);

  // Фильтруем породы при вводе
  const handleBreedInput = (value: string) => {
    setFormData({ ...formData, breed: value });
    
    if (value.trim() === '') {
      setFilteredBreeds(breeds);
      setShowBreedDropdown(false);
    } else {
      const filtered = breeds.filter(breed =>
        breed.name.toLowerCase().includes(value.toLowerCase()) ||
        breed.name_en.toLowerCase().includes(value.toLowerCase())
      );
      setFilteredBreeds(filtered);
      setShowBreedDropdown(true);
    }
  };

  // Выбор породы из списка
  const selectBreed = (breedName: string) => {
    setFormData({ ...formData, breed: breedName });
    setShowBreedDropdown(false);
  };

  // Вычисление возраста из даты рождения
  const calculateAgeFromBirthdate = (birthdate: string): number => {
    const birth = new Date(birthdate);
    const today = new Date();
    let age = today.getFullYear() - birth.getFullYear();
    const monthDiff = today.getMonth() - birth.getMonth();
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
      age--;
    }
    return age;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await onSubmit(formData);
      // Сбрасываем форму после успешного добавления
      setFormData({
        name: '',
        species: 'dog',
        breed: '',
        age: 0,
        sex: 'male',
        color: '',
        weight: undefined,
        chip_number: '',
        sterilized: false,
      });
      setBirthdate('');
      setShowBreedDropdown(false);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось добавить питомца');
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={onClose}
      />
      
      {/* Modal */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div className="relative bg-white rounded-xl shadow-xl max-w-lg w-full">
          {/* Header */}
          <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4 rounded-t-xl flex items-center justify-between">
            <h2 className="text-2xl font-bold text-white">Добавить питомца</h2>
            <button
              onClick={onClose}
              className="text-white hover:text-gray-200 transition-colors"
            >
              <XMarkIcon className="w-6 h-6" />
            </button>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="p-6 space-y-4">
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-3 text-red-700 text-sm">
                {error}
              </div>
            )}

            {/* Имя */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Имя питомца *
              </label>
              <input
                type="text"
                required
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Например: Барсик"
              />
            </div>

            {/* Вид (кошка/собака) */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Вид *
              </label>
              <div className="grid grid-cols-2 gap-3">
                <button
                  type="button"
                  onClick={() => {
                    setFormData({ ...formData, species: 'dog', breed: '' });
                    setShowBreedDropdown(false);
                  }}
                  className={`px-4 py-3 rounded-lg border-2 transition-all ${
                    formData.species === 'dog'
                      ? 'border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className="text-3xl mb-1">🐕</div>
                  <div className="text-sm font-medium">Собака</div>
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setFormData({ ...formData, species: 'cat', breed: '' });
                    setShowBreedDropdown(false);
                  }}
                  className={`px-4 py-3 rounded-lg border-2 transition-all ${
                    formData.species === 'cat'
                      ? 'border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className="text-3xl mb-1">🐈</div>
                  <div className="text-sm font-medium">Кошка</div>
                </button>
              </div>
            </div>

            {/* Порода с автодополнением */}
            <div className="relative">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Порода
              </label>
              <input
                type="text"
                value={formData.breed || ''}
                onChange={(e) => handleBreedInput(e.target.value)}
                onFocus={() => {
                  if (formData.breed && filteredBreeds.length > 0) {
                    setShowBreedDropdown(true);
                  }
                }}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Начните вводить породу..."
                autoComplete="off"
              />
              
              {/* Dropdown с породами */}
              {showBreedDropdown && filteredBreeds.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                  {filteredBreeds.map((breed) => (
                    <button
                      key={breed.id}
                      type="button"
                      onClick={() => selectBreed(breed.name)}
                      className="w-full px-3 py-2 text-left hover:bg-blue-50 transition-colors border-b border-gray-100 last:border-b-0"
                    >
                      <div className="font-medium text-gray-900">{breed.name}</div>
                      <div className="text-xs text-gray-500">{breed.name_en}</div>
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Возраст или дата рождения */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-gray-700">
                  {ageInputType === 'age' ? 'Примерный возраст (лет) *' : 'Дата рождения *'}
                </label>
                <button
                  type="button"
                  onClick={() => setAgeInputType(ageInputType === 'age' ? 'birthdate' : 'age')}
                  className="text-xs text-blue-600 hover:text-blue-700 underline"
                >
                  {ageInputType === 'age' ? 'Указать точную дату' : 'Указать примерный возраст'}
                </button>
              </div>
              
              {ageInputType === 'age' ? (
                <input
                  type="number"
                  required
                  min="0"
                  max="30"
                  value={formData.age}
                  onChange={(e) => setFormData({ ...formData, age: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Например: 3"
                />
              ) : (
                <input
                  type="date"
                  required
                  value={birthdate}
                  max={new Date().toISOString().split('T')[0]}
                  onChange={(e) => {
                    setBirthdate(e.target.value);
                    const age = calculateAgeFromBirthdate(e.target.value);
                    setFormData({ ...formData, age });
                  }}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              )}
              <p className="text-xs text-gray-500 mt-1">
                {ageInputType === 'age' 
                  ? 'Укажите примерный возраст, если точная дата неизвестна' 
                  : 'Возраст будет рассчитан автоматически'}
              </p>
            </div>

            {/* Пол */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Пол *
              </label>
              <select
                value={formData.sex}
                onChange={(e) => setFormData({ ...formData, sex: e.target.value as 'male' | 'female' })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="male">♂️ Самец</option>
                <option value="female">♀️ Самка</option>
              </select>
            </div>

            {/* Buttons */}
            <div className="flex gap-3 pt-4">
              <button
                type="button"
                onClick={onClose}
                disabled={loading}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={loading}
                className="flex-1 px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed font-medium"
              >
                {loading ? 'Добавление...' : 'Добавить'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
