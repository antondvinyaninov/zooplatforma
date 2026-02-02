'use client';

import { useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { 
  PetEventType, 
  eventTypeLabels, 
  eventTypeIcons,
  deathReasonLabels,
  petEventsApi 
} from '../../../lib/api';

interface AddEventModalProps {
  petId: number;
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export default function AddEventModal({ petId, isOpen, onClose, onSuccess }: AddEventModalProps) {
  const [loading, setLoading] = useState(false);
  const [eventType, setEventType] = useState<PetEventType>('general');
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    event_date: new Date().toISOString().split('T')[0],
    // Vaccination fields
    vaccine_name: '',
    vaccine_batch: '',
    next_date: '',
    // Treatment fields
    medication_name: '',
    dosage: '',
    // Ownership change fields
    previous_owner_name: '',
    new_owner_name: '',
    transfer_reason: '',
    // Lost/Found fields
    location: '',
    circumstances: '',
    contact_name: '',
    contact_phone: '',
    // Death fields
    death_reason: '',
    // Shelter fields
    shelter_name: '',
    adoption_contract: '',
  });

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      
      const eventData: any = {
        event_type: eventType,
        event_date: formData.event_date,
        title: formData.title || undefined,
        description: formData.description || undefined,
      };

      // Add type-specific fields
      if (eventType === 'vaccination') {
        if (formData.vaccine_name) eventData.vaccine_name = formData.vaccine_name;
        if (formData.vaccine_batch) eventData.vaccine_batch = formData.vaccine_batch;
        if (formData.next_date) eventData.next_date = formData.next_date;
      } else if (eventType === 'treatment') {
        if (formData.medication_name) eventData.medication_name = formData.medication_name;
        if (formData.dosage) eventData.dosage = formData.dosage;
      } else if (eventType === 'ownership_change') {
        if (formData.previous_owner_name) eventData.previous_owner_name = formData.previous_owner_name;
        if (formData.new_owner_name) eventData.new_owner_name = formData.new_owner_name;
        if (formData.transfer_reason) eventData.transfer_reason = formData.transfer_reason;
      } else if (eventType === 'lost' || eventType === 'found') {
        if (formData.location) eventData.location = formData.location;
        if (formData.circumstances) eventData.circumstances = formData.circumstances;
        if (formData.contact_name) eventData.contact_name = formData.contact_name;
        if (formData.contact_phone) eventData.contact_phone = formData.contact_phone;
      } else if (eventType === 'death') {
        if (formData.death_reason) eventData.death_reason = formData.death_reason;
      } else if (eventType === 'shelter_intake' || eventType === 'adoption') {
        if (formData.shelter_name) eventData.shelter_name = formData.shelter_name;
        if (formData.adoption_contract) eventData.adoption_contract = formData.adoption_contract;
      }

      const response = await petEventsApi.createEvent(petId, eventData);
      
      if (response.success) {
        onSuccess();
        onClose();
        // Reset form
        setFormData({
          title: '',
          description: '',
          event_date: new Date().toISOString().split('T')[0],
          vaccine_name: '',
          vaccine_batch: '',
          next_date: '',
          medication_name: '',
          dosage: '',
          previous_owner_name: '',
          new_owner_name: '',
          transfer_reason: '',
          location: '',
          circumstances: '',
          contact_name: '',
          contact_phone: '',
          death_reason: '',
          shelter_name: '',
          adoption_contract: '',
        });
        setEventType('general');
      }
    } catch (error) {
      console.error('Error creating event:', error);
      alert('Ошибка при создании события');
    } finally {
      setLoading(false);
    }
  };

  const renderTypeSpecificFields = () => {
    switch (eventType) {
      case 'vaccination':
        return (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Название вакцины
              </label>
              <input
                type="text"
                value={formData.vaccine_name}
                onChange={(e) => setFormData({ ...formData, vaccine_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Например: Нобивак DHPPi"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Серия вакцины
              </label>
              <input
                type="text"
                value={formData.vaccine_batch}
                onChange={(e) => setFormData({ ...formData, vaccine_batch: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Дата следующей вакцинации
              </label>
              <input
                type="date"
                value={formData.next_date}
                onChange={(e) => setFormData({ ...formData, next_date: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </>
        );

      case 'treatment':
        return (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Название препарата
              </label>
              <input
                type="text"
                value={formData.medication_name}
                onChange={(e) => setFormData({ ...formData, medication_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Дозировка
              </label>
              <input
                type="text"
                value={formData.dosage}
                onChange={(e) => setFormData({ ...formData, dosage: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Например: 1 таблетка 2 раза в день"
              />
            </div>
          </>
        );

      case 'ownership_change':
        return (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Предыдущий владелец
              </label>
              <input
                type="text"
                value={formData.previous_owner_name}
                onChange={(e) => setFormData({ ...formData, previous_owner_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Новый владелец
              </label>
              <input
                type="text"
                value={formData.new_owner_name}
                onChange={(e) => setFormData({ ...formData, new_owner_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Причина передачи
              </label>
              <input
                type="text"
                value={formData.transfer_reason}
                onChange={(e) => setFormData({ ...formData, transfer_reason: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </>
        );

      case 'lost':
      case 'found':
        return (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Место
              </label>
              <input
                type="text"
                value={formData.location}
                onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Например: ул. Пушкина, д. 10"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Обстоятельства
              </label>
              <textarea
                value={formData.circumstances}
                onChange={(e) => setFormData({ ...formData, circumstances: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                rows={2}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Контактное лицо
              </label>
              <input
                type="text"
                value={formData.contact_name}
                onChange={(e) => setFormData({ ...formData, contact_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Телефон
              </label>
              <input
                type="tel"
                value={formData.contact_phone}
                onChange={(e) => setFormData({ ...formData, contact_phone: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </>
        );

      case 'death':
        return (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Причина смерти
            </label>
            <select
              value={formData.death_reason}
              onChange={(e) => setFormData({ ...formData, death_reason: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Выберите причину</option>
              {Object.entries(deathReasonLabels).map(([value, label]) => (
                <option key={value} value={value}>{label}</option>
              ))}
            </select>
          </div>
        );

      case 'shelter_intake':
      case 'adoption':
        return (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Название приюта
              </label>
              <input
                type="text"
                value={formData.shelter_name}
                onChange={(e) => setFormData({ ...formData, shelter_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            {eventType === 'adoption' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Номер договора
                </label>
                <input
                  type="text"
                  value={formData.adoption_contract}
                  onChange={(e) => setFormData({ ...formData, adoption_contract: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            )}
          </>
        );

      default:
        return null;
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <h2 className="text-xl font-bold text-gray-900">Добавить событие</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          {/* Event Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Тип события *
            </label>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
              {(Object.keys(eventTypeLabels) as PetEventType[]).map(type => (
                <button
                  key={type}
                  type="button"
                  onClick={() => setEventType(type)}
                  className={`flex items-center gap-2 px-3 py-2 rounded-lg border-2 text-sm font-medium transition-colors ${
                    eventType === type
                      ? 'bg-blue-50 border-blue-500 text-blue-700'
                      : 'bg-white border-gray-200 text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <span>{eventTypeIcons[type]}</span>
                  <span>{eventTypeLabels[type]}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Event Date */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Дата события *
            </label>
            <input
              type="date"
              required
              value={formData.event_date}
              onChange={(e) => setFormData({ ...formData, event_date: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {/* Title */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Заголовок (необязательно)
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Например: Первая прививка"
            />
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Описание
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              rows={3}
              placeholder="Дополнительная информация о событии"
            />
          </div>

          {/* Type-specific fields */}
          {renderTypeSpecificFields()}

          {/* Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border-2 border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Сохранение...' : 'Добавить событие'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
