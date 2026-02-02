'use client';

import { useState } from 'react';
import { CheckCircleIcon } from '@heroicons/react/24/outline';

interface Pet {
  id: number;
  user_id?: number;
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
}

interface Props {
  pet: Pet;
  onComplete: () => void;
}

export default function CompleteCardForm({ pet, onComplete }: Props) {
  const [formData, setFormData] = useState({
    chip_number: pet.chip_number || '',
    is_sterilized: false,
    sterilization_date: '',
    last_vaccination_date: '',
    blood_type: '',
    notes: '',
  });
  const [photoFile, setPhotoFile] = useState<File | null>(null);
  const [photoPreview, setPhotoPreview] = useState<string | null>(
    pet.photo ? `http://localhost:8000${pet.photo}` : null
  );
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      console.log('🏥 Clinic ID:', clinicId);
      
      // Получаем токен из Main API
      console.log('🔐 Step 1: Fetching token from Main API...');
      const meResponse = await fetch('http://localhost:7100/api/auth/me', {
        credentials: 'include',
      });

      console.log('📡 Main API response status:', meResponse.status);
      
      if (!meResponse.ok) {
        console.error('❌ Main API returned error:', meResponse.status);
        throw new Error('Unauthorized');
      }

      const meData = await meResponse.json();
      console.log('📦 Main API response data:', meData);
      
      // Токен может быть либо в meData.token, либо в meData.data.token
      const token = meData.token || meData.data?.token;
      console.log('🎫 Token extracted:', token ? `${token.substring(0, 20)}...` : 'NULL');

      if (!token) {
        console.error('❌ No token in response');
        throw new Error('No auth token');
      }

      let photoUrl = pet.photo ? `http://localhost:8000${pet.photo}` : null;

      // 1. Загружаем новое фото (если выбрано)
      if (photoFile) {
        console.log('📸 Step 2: Uploading new photo...');
        const formData = new FormData();
        formData.append('file', photoFile);
        formData.append('user_id', pet.user_id?.toString() || '');

        const uploadResponse = await fetch('http://localhost:8000/api/media/upload', {
          method: 'POST',
          credentials: 'include',
          body: formData,
        });

        console.log('📡 Upload response status:', uploadResponse.status);

        if (uploadResponse.ok) {
          const uploadData = await uploadResponse.json();
          photoUrl = uploadData.data?.file_path;
          console.log('✅ Photo uploaded:', photoUrl);
        } else {
          console.error('❌ Photo upload failed');
        }
      } else {
        console.log('⏭️ Step 2: Skipping photo upload (no new photo)');
      }

      // 2. Обновляем карточку питомца
      console.log('🐾 Step 3: Updating pet card...');
      console.log('📋 Request URL:', `http://localhost:8100/api/pets/${pet.id}`);
      console.log('🔑 Authorization header:', `Bearer ${token.substring(0, 20)}...`);
      console.log('🏥 X-Clinic-ID header:', clinicId);
      
      const updatePayload = {
        chip_number: formData.chip_number,
        is_sterilized: formData.is_sterilized,
        sterilization_date: formData.sterilization_date || null,
        last_vaccination_date: formData.last_vaccination_date || null,
        blood_type: formData.blood_type || null,
        photo: photoUrl,
        status: 'verified',
      };
      console.log('📦 Update payload:', updatePayload);

      const updateResponse = await fetch(`http://localhost:8100/api/pets/${pet.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
          'X-Clinic-ID': clinicId || '',
        },
        body: JSON.stringify(updatePayload),
      });

      console.log('📡 Update response status:', updateResponse.status);
      console.log('📡 Update response headers:', Object.fromEntries(updateResponse.headers.entries()));

      if (!updateResponse.ok) {
        const errorData = await updateResponse.json();
        console.error('❌ Update failed:', errorData);
        throw new Error(errorData.error || 'Failed to update pet');
      }

      const updateData = await updateResponse.json();
      console.log('✅ Pet updated:', updateData);

      // 3. Создаем событие "Подтверждение PetID клиникой"
      console.log('📝 Step 4: Creating pet event...');
      const eventResponse = await fetch('http://localhost:8100/api/pet-events', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
          'X-Clinic-ID': clinicId || '',
        },
        body: JSON.stringify({
          pet_id: pet.id,
          event_type: 'verification',
          event_date: new Date().toISOString().split('T')[0],
          description: 'Подтверждение PetID клиникой',
          clinic_id: parseInt(clinicId || '0'),
          notes: formData.notes || null,
        }),
      });

      console.log('📡 Event response status:', eventResponse.status);

      if (!eventResponse.ok) {
        console.error('❌ Failed to create event');
      } else {
        console.log('✅ Event created');
      }

      // 4. Отправляем уведомление владельцу
      if (pet.user_id) {
        console.log('📬 Step 5: Sending notification to owner...');
        await fetch('http://localhost:8000/api/notifications/send', {
          method: 'POST',
          headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify({
            user_id: pet.user_id,
            type: 'pet_verified',
            title: 'Карточка питомца подтверждена',
            message: `Карточка вашего питомца ${pet.name} подтверждена ветеринарной клиникой`,
            link: `/pets/${pet.id}`,
          }),
        }).then(res => {
          console.log('📡 Notification response status:', res.status);
          if (res.ok) {
            console.log('✅ Notification sent');
          } else {
            console.error('❌ Failed to send notification');
          }
        }).catch(err => console.error('❌ Notification error:', err));
      } else {
        console.log('⏭️ Step 5: Skipping notification (no user_id)');
      }

      console.log('🎉 All steps completed successfully!');
      setSuccess(true);
      setTimeout(() => {
        onComplete();
      }, 2000);
    } catch (error) {
      console.error('💥 Error in handleSubmit:', error);
      console.error('💥 Error stack:', error instanceof Error ? error.stack : 'No stack');
      alert('Ошибка при сохранении данных: ' + (error instanceof Error ? error.message : 'Unknown error'));
    }

    setLoading(false);
  };

  if (success) {
    return (
      <div className="bg-white rounded-2xl p-12 shadow-sm border border-gray-100 text-center">
        <CheckCircleIcon className="w-20 h-20 text-green-500 mx-auto mb-6" />
        <h3 className="text-3xl font-bold text-gray-900 mb-4">
          Карточка подтверждена!
        </h3>
        <p className="text-lg text-gray-600 mb-6">
          Статус изменен на "Verified". Владелец получит уведомление.
        </p>
        <div className="text-sm text-gray-500">
          Перенаправление...
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">
        Дополнение карточки: {pet.name}
      </h3>

      {/* Информация о питомце */}
      <div className="bg-gray-50 rounded-xl p-6 mb-8">
        <h4 className="font-semibold text-gray-900 mb-4">Текущая информация</h4>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">Кличка:</span>
            <span className="ml-2 text-gray-900 font-medium">{pet.name}</span>
          </div>
          <div>
            <span className="text-gray-500">Вид:</span>
            <span className="ml-2 text-gray-900 font-medium">{pet.species}</span>
          </div>
          {pet.breed && (
            <div>
              <span className="text-gray-500">Порода:</span>
              <span className="ml-2 text-gray-900 font-medium">{pet.breed}</span>
            </div>
          )}
          {pet.gender && (
            <div>
              <span className="text-gray-500">Пол:</span>
              <span className="ml-2 text-gray-900 font-medium">
                {pet.gender === 'male' ? 'Самец' : 'Самка'}
              </span>
            </div>
          )}
          {pet.owner_name && (
            <div>
              <span className="text-gray-500">Владелец:</span>
              <span className="ml-2 text-gray-900 font-medium">{pet.owner_name}</span>
            </div>
          )}
          {pet.owner_phone && (
            <div>
              <span className="text-gray-500">Телефон:</span>
              <span className="ml-2 text-gray-900 font-medium">{pet.owner_phone}</span>
            </div>
          )}
        </div>
      </div>

      {/* Форма дополнения */}
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Фото питомца */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Фото питомца
          </label>
          <div className="flex items-start space-x-4">
            {photoPreview ? (
              <div className="relative">
                <img
                  src={photoPreview}
                  alt="Preview"
                  className="w-32 h-32 object-cover rounded-xl border-2 border-gray-200"
                />
                {photoFile && (
                  <button
                    type="button"
                    onClick={() => {
                      setPhotoFile(null);
                      setPhotoPreview(pet.photo ? `http://localhost:8000${pet.photo}` : null);
                    }}
                    className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors"
                  >
                    ×
                  </button>
                )}
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
                {pet.photo ? 'Изменить фото' : 'Добавить фото'}
              </label>
              <p className="text-sm text-gray-500 mt-2">
                Рекомендуемый размер: 800x800px. Форматы: JPG, PNG
              </p>
            </div>
          </div>
        </div>

        {/* Chip number */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Chip number {!pet.chip_number && <span className="text-red-500">*</span>}
          </label>
          <input
            type="text"
            value={formData.chip_number}
            onChange={(e) => setFormData({ ...formData, chip_number: e.target.value })}
            placeholder="643094100123456"
            required={!pet.chip_number}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono"
          />
          <p className="mt-1 text-sm text-gray-500">
            15-значный номер микрочипа
          </p>
        </div>

        {/* Стерилизация */}
        <div>
          <label className="flex items-center space-x-3 cursor-pointer">
            <input
              type="checkbox"
              checked={formData.is_sterilized}
              onChange={(e) => setFormData({ ...formData, is_sterilized: e.target.checked })}
              className="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <span className="text-sm font-medium text-gray-700">
              Животное стерилизовано/кастрировано
            </span>
          </label>
        </div>

        {formData.is_sterilized && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Дата стерилизации
            </label>
            <input
              type="date"
              value={formData.sterilization_date}
              onChange={(e) => setFormData({ ...formData, sterilization_date: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        )}

        {/* Вакцинация */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Дата последней вакцинации
          </label>
          <input
            type="date"
            value={formData.last_vaccination_date}
            onChange={(e) => setFormData({ ...formData, last_vaccination_date: e.target.value })}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Группа крови */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Группа крови
          </label>
          <input
            type="text"
            value={formData.blood_type}
            onChange={(e) => setFormData({ ...formData, blood_type: e.target.value })}
            placeholder="A, B, AB и т.д."
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Примечания */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Примечания
          </label>
          <textarea
            value={formData.notes}
            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
            rows={4}
            placeholder="Дополнительная информация о питомце..."
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
          />
        </div>

        {/* Кнопки */}
        <div className="flex space-x-4 pt-4">
          <button
            type="submit"
            disabled={loading}
            className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed font-medium"
          >
            {loading ? 'Сохранение...' : 'Подтвердить карточку'}
          </button>
          <button
            type="button"
            onClick={onComplete}
            className="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
          >
            Отмена
          </button>
        </div>
      </form>
    </div>
  );
}
