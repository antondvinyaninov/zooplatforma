'use client';

import { useState } from 'react';
import { CheckCircleIcon, MagnifyingGlassIcon } from '@heroicons/react/24/outline';

interface Owner {
  id?: number;
  name: string;
  phone: string;
  email?: string;
}

export default function CreateNewCardForm() {
  const [step, setStep] = useState<'owner' | 'pet'>('owner');
  const [owner, setOwner] = useState<Owner | null>(null);
  const [ownerSearch, setOwnerSearch] = useState('');
  const [searchResults, setSearchResults] = useState<Owner[]>([]);
  
  const [newOwner, setNewOwner] = useState({
    name: '',
    phone: '',
    email: '',
  });

  const [petData, setPetData] = useState({
    name: '',
    species: 'dog',
    breed: '',
    birth_date: '',
    gender: 'male',
    color: '',
    chip_number: '',
    is_sterilized: false,
    sterilization_date: '',
    last_vaccination_date: '',
    blood_type: '',
    notes: '',
  });

  const [photoFile, setPhotoFile] = useState<File | null>(null);
  const [photoPreview, setPhotoPreview] = useState<string | null>(null);

  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  // Поиск владельца
  const handleSearchOwner = async () => {
    if (!ownerSearch.trim()) return;

    try {
      const response = await fetch(
        `http://localhost:8000/api/users/search?phone=${encodeURIComponent(ownerSearch)}`,
        { credentials: 'include' }
      );

      if (response.ok) {
        const data = await response.json();
        if (data.success) {
          setSearchResults(data.data || []);
        }
      }
    } catch (error) {
      console.error('Search failed:', error);
    }
  };

  // Выбор существующего владельца
  const handleSelectOwner = (selectedOwner: Owner) => {
    setOwner(selectedOwner);
    setStep('pet');
  };

  // Создание нового владельца
  const handleCreateNewOwner = () => {
    if (!newOwner.name || !newOwner.phone) {
      alert('Заполните имя и телефон владельца');
      return;
    }
    setOwner({ ...newOwner });
    setStep('pet');
  };

  // Генерация chip number
  const generateChipNumber = () => {
    const prefix = '643094'; // Россия
    const random = Math.floor(Math.random() * 1000000000).toString().padStart(9, '0');
    setPetData({ ...petData, chip_number: prefix + random });
  };

  // Обработка выбора фото
  const handlePhotoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setPhotoFile(file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setPhotoPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  // Создание карточки
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      let ownerId = owner?.id;

      // 1. Если владелец новый - создаем базовый профиль
      if (!ownerId) {
        const createUserResponse = await fetch('http://localhost:8000/api/users', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({
            name: owner?.name,
            phone: owner?.phone,
            email: owner?.email || null,
            password: Math.random().toString(36).slice(-8), // Временный пароль
          }),
        });

        if (createUserResponse.ok) {
          const userData = await createUserResponse.json();
          ownerId = userData.data?.id;
        }
      }

      if (!ownerId) {
        throw new Error('Failed to get owner ID');
      }

      // 2. Загружаем фото (если есть)
      let photoUrl = null;
      if (photoFile) {
        const formData = new FormData();
        formData.append('file', photoFile);
        formData.append('user_id', ownerId.toString());

        const uploadResponse = await fetch('http://localhost:8000/api/media/upload', {
          method: 'POST',
          credentials: 'include',
          body: formData,
        });

        if (uploadResponse.ok) {
          const uploadData = await uploadResponse.json();
          photoUrl = uploadData.data?.file_path;
        }
      }

      // 3. Создаем карточку питомца
      const createPetResponse = await fetch('http://localhost:8100/api/pets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Clinic-ID': clinicId || '',
        },
        credentials: 'include',
        body: JSON.stringify({
          user_id: ownerId,
          name: petData.name,
          species: petData.species,
          breed: petData.breed || null,
          birth_date: petData.birth_date || null,
          gender: petData.gender,
          color: petData.color || null,
          chip_number: petData.chip_number,
          is_sterilized: petData.is_sterilized,
          sterilization_date: petData.sterilization_date || null,
          last_vaccination_date: petData.last_vaccination_date || null,
          blood_type: petData.blood_type || null,
          photo: photoUrl,
          status: 'verified', // Сразу verified
        }),
      });

      if (!createPetResponse.ok) {
        throw new Error('Failed to create pet');
      }

      const petResult = await createPetResponse.json();
      const petId = petResult.data?.id;

      // 4. Создаем событие "Регистрация в клинике"
      if (petId) {
        await fetch('http://localhost:8100/api/pet-events', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Clinic-ID': clinicId || '',
          },
          credentials: 'include',
          body: JSON.stringify({
            pet_id: petId,
            event_type: 'registration',
            event_date: new Date().toISOString().split('T')[0],
            description: 'Регистрация в клинике',
            clinic_id: parseInt(clinicId || '0'),
            notes: petData.notes || null,
          }),
        });
      }

      // 5. Отправляем уведомление владельцу (если есть телефон)
      if (owner?.phone && petId) {
        await fetch('http://localhost:8000/api/notifications/send', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({
            user_id: ownerId,
            type: 'pet_registered',
            title: 'Питомец зарегистрирован',
            message: `Ваш питомец ${petData.name} успешно зарегистрирован в системе PetID`,
            link: `/pets/${petId}`,
          }),
        }).catch(err => console.error('Failed to send notification:', err));
      }

      setSuccess(true);
    } catch (error) {
      console.error('Failed to create card:', error);
      alert('Ошибка при создании карточки');
    }

    setLoading(false);
  };

  if (success) {
    return (
      <div className="bg-white rounded-2xl p-12 shadow-sm border border-gray-100 text-center">
        <CheckCircleIcon className="w-20 h-20 text-green-500 mx-auto mb-6" />
        <h3 className="text-3xl font-bold text-gray-900 mb-4">
          Карточка создана!
        </h3>
        <p className="text-lg text-gray-600 mb-6">
          Питомец зарегистрирован в системе. Владелец получит уведомление.
        </p>
        <button
          onClick={() => window.location.reload()}
          className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
        >
          Создать еще одну карточку
        </button>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">
        Создание новой карточки
      </h3>

      {/* Шаг 1: Владелец */}
      {step === 'owner' && (
        <div className="space-y-6">
          <div className="bg-blue-50 border border-blue-200 rounded-xl p-4">
            <p className="text-sm text-blue-800">
              <strong>Шаг 1 из 2:</strong> Найдите владельца в системе или создайте нового
            </p>
          </div>

          {/* Поиск существующего владельца */}
          <div>
            <h4 className="font-semibold text-gray-900 mb-4">Поиск владельца</h4>
            <div className="flex space-x-4 mb-4">
              <input
                type="text"
                value={ownerSearch}
                onChange={(e) => setOwnerSearch(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearchOwner()}
                placeholder="Телефон или email"
                className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <button
                onClick={handleSearchOwner}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2"
              >
                <MagnifyingGlassIcon className="w-5 h-5" />
                <span>Найти</span>
              </button>
            </div>

            {searchResults.length > 0 && (
              <div className="space-y-2">
                {searchResults.map((result, index) => (
                  <button
                    key={index}
                    onClick={() => handleSelectOwner(result)}
                    className="w-full text-left p-4 border border-gray-200 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-all"
                  >
                    <div className="font-medium text-gray-900">{result.name}</div>
                    <div className="text-sm text-gray-600">
                      {result.phone} {result.email && `• ${result.email}`}
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-white text-gray-500">или</span>
            </div>
          </div>

          {/* Создание нового владельца */}
          <div>
            <h4 className="font-semibold text-gray-900 mb-4">Создать нового владельца</h4>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Имя <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={newOwner.name}
                  onChange={(e) => setNewOwner({ ...newOwner, name: e.target.value })}
                  placeholder="Иван Иванов"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Телефон <span className="text-red-500">*</span>
                </label>
                <input
                  type="tel"
                  value={newOwner.phone}
                  onChange={(e) => setNewOwner({ ...newOwner, phone: e.target.value })}
                  placeholder="+7 (999) 123-45-67"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <input
                  type="email"
                  value={newOwner.email}
                  onChange={(e) => setNewOwner({ ...newOwner, email: e.target.value })}
                  placeholder="ivan@example.com"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <button
                onClick={handleCreateNewOwner}
                className="w-full px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
              >
                Продолжить с новым владельцем →
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Шаг 2: Питомец */}
      {step === 'pet' && owner && (
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="bg-green-50 border border-green-200 rounded-xl p-4">
            <p className="text-sm text-green-800">
              <strong>Шаг 2 из 2:</strong> Заполните информацию о питомце
            </p>
            <p className="text-sm text-green-700 mt-1">
              Владелец: <strong>{owner.name}</strong> • {owner.phone}
            </p>
          </div>

          {/* Основная информация */}
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Кличка <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                value={petData.name}
                onChange={(e) => setPetData({ ...petData, name: e.target.value })}
                required
                placeholder="Барсик"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Вид <span className="text-red-500">*</span>
              </label>
              <select
                value={petData.species}
                onChange={(e) => setPetData({ ...petData, species: e.target.value })}
                required
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="dog">Собака</option>
                <option value="cat">Кошка</option>
                <option value="bird">Птица</option>
                <option value="rodent">Грызун</option>
                <option value="reptile">Рептилия</option>
                <option value="other">Другое</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Порода
              </label>
              <input
                type="text"
                value={petData.breed}
                onChange={(e) => setPetData({ ...petData, breed: e.target.value })}
                placeholder="Лабрадор"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Дата рождения
              </label>
              <input
                type="date"
                value={petData.birth_date}
                onChange={(e) => setPetData({ ...petData, birth_date: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Пол <span className="text-red-500">*</span>
              </label>
              <select
                value={petData.gender}
                onChange={(e) => setPetData({ ...petData, gender: e.target.value })}
                required
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="male">Самец</option>
                <option value="female">Самка</option>
              </select>
            </div>

            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Окрас
              </label>
              <input
                type="text"
                value={petData.color}
                onChange={(e) => setPetData({ ...petData, color: e.target.value })}
                placeholder="Рыжий"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            {/* Фото питомца */}
            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Фото питомца (опционально)
              </label>
              <div className="flex items-start space-x-4">
                {photoPreview ? (
                  <div className="relative">
                    <img
                      src={photoPreview}
                      alt="Preview"
                      className="w-32 h-32 object-cover rounded-xl border-2 border-gray-200"
                    />
                    <button
                      type="button"
                      onClick={() => {
                        setPhotoFile(null);
                        setPhotoPreview(null);
                      }}
                      className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors"
                    >
                      ×
                    </button>
                  </div>
                ) : (
                  <div className="w-32 h-32 bg-gray-100 rounded-xl border-2 border-dashed border-gray-300 flex items-center justify-center">
                    <span className="text-4xl">📷</span>
                  </div>
                )}
                <div className="flex-1">
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handlePhotoChange}
                    className="hidden"
                    id="photo-upload"
                  />
                  <label
                    htmlFor="photo-upload"
                    className="inline-block px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors cursor-pointer"
                  >
                    Выбрать фото
                  </label>
                  <p className="text-sm text-gray-500 mt-2">
                    Рекомендуемый размер: 800x800px. Форматы: JPG, PNG
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Chip number */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Chip number <span className="text-red-500">*</span>
            </label>
            <div className="flex space-x-4">
              <input
                type="text"
                value={petData.chip_number}
                onChange={(e) => setPetData({ ...petData, chip_number: e.target.value })}
                required
                placeholder="643094100123456"
                className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono"
              />
              <button
                type="button"
                onClick={generateChipNumber}
                className="px-4 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
              >
                Генерировать
              </button>
            </div>
          </div>

          {/* Медицинская информация */}
          <div className="border-t border-gray-200 pt-6">
            <h4 className="font-semibold text-gray-900 mb-4">Медицинская информация</h4>
            
            <div className="space-y-4">
              <label className="flex items-center space-x-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={petData.is_sterilized}
                  onChange={(e) => setPetData({ ...petData, is_sterilized: e.target.checked })}
                  className="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span className="text-sm font-medium text-gray-700">
                  Стерилизовано/кастрировано
                </span>
              </label>

              {petData.is_sterilized && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Дата стерилизации
                  </label>
                  <input
                    type="date"
                    value={petData.sterilization_date}
                    onChange={(e) => setPetData({ ...petData, sterilization_date: e.target.value })}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Дата последней вакцинации
                </label>
                <input
                  type="date"
                  value={petData.last_vaccination_date}
                  onChange={(e) => setPetData({ ...petData, last_vaccination_date: e.target.value })}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Группа крови
                </label>
                <input
                  type="text"
                  value={petData.blood_type}
                  onChange={(e) => setPetData({ ...petData, blood_type: e.target.value })}
                  placeholder="A, B, AB и т.д."
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Примечания
                </label>
                <textarea
                  value={petData.notes}
                  onChange={(e) => setPetData({ ...petData, notes: e.target.value })}
                  rows={3}
                  placeholder="Дополнительная информация..."
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                />
              </div>
            </div>
          </div>

          {/* Кнопки */}
          <div className="flex space-x-4 pt-4">
            <button
              type="button"
              onClick={() => setStep('owner')}
              className="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
            >
              ← Назад
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed font-medium"
            >
              {loading ? 'Создание...' : 'Создать карточку'}
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
