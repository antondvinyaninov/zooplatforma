'use client';

import { useState } from 'react';
import { MagnifyingGlassIcon, PlusIcon } from '@heroicons/react/24/outline';
import SearchExistingCard from '@/app/components/registration/SearchExistingCard';
import CreateNewCardForm from '@/app/components/registration/CreateNewCardForm';

type Mode = 'select' | 'search' | 'create';

export default function RegistrationPage() {
  const [mode, setMode] = useState<Mode>('select');

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-4xl font-bold text-gray-900 mb-2">Регистрация животных</h2>
        <p className="text-base text-gray-600">
          Дополните существующую карточку или создайте новую
        </p>
      </div>

      {mode === 'select' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Сценарий 1: Дополнить существующую карточку */}
          <button
            onClick={() => setMode('search')}
            className="bg-white rounded-2xl p-8 shadow-sm border-2 border-gray-200 hover:border-blue-500 hover:shadow-lg transition-all text-left group"
          >
            <div className="flex items-center justify-between mb-6">
              <div className="w-16 h-16 bg-blue-100 rounded-xl flex items-center justify-center group-hover:bg-blue-500 transition-colors">
                <MagnifyingGlassIcon className="w-8 h-8 text-blue-600 group-hover:text-white transition-colors" />
              </div>
              <div className="text-4xl">📋</div>
            </div>
            
            <h3 className="text-2xl font-bold text-gray-900 mb-3">
              Дополнить существующую карточку
            </h3>
            
            <p className="text-gray-600 mb-4">
              Владелец уже создал карточку питомца в своем кабинете. 
              Найдите карточку и дополните медицинской информацией.
            </p>

            <div className="space-y-2 text-sm text-gray-500">
              <div className="flex items-center space-x-2">
                <span className="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>
                <span>Поиск по телефону владельца</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>
                <span>Статус: pending_verification → verified</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>
                <span>Добавление медицинской информации</span>
              </div>
            </div>

            <div className="mt-6 text-blue-600 font-medium group-hover:text-blue-700">
              Найти карточку →
            </div>
          </button>

          {/* Сценарий 2: Создать новую карточку */}
          <button
            onClick={() => setMode('create')}
            className="bg-white rounded-2xl p-8 shadow-sm border-2 border-gray-200 hover:border-green-500 hover:shadow-lg transition-all text-left group"
          >
            <div className="flex items-center justify-between mb-6">
              <div className="w-16 h-16 bg-green-100 rounded-xl flex items-center justify-center group-hover:bg-green-500 transition-colors">
                <PlusIcon className="w-8 h-8 text-green-600 group-hover:text-white transition-colors" />
              </div>
              <div className="text-4xl">➕</div>
            </div>
            
            <h3 className="text-2xl font-bold text-gray-900 mb-3">
              Создать новую карточку
            </h3>
            
            <p className="text-gray-600 mb-4">
              Владелец пришел без карточки. Создайте мини-карточку питомца 
              с нуля и привяжите к владельцу.
            </p>

            <div className="space-y-2 text-sm text-gray-500">
              <div className="flex items-center space-x-2">
                <span className="w-1.5 h-1.5 bg-green-500 rounded-full"></span>
                <span>Поиск или создание владельца</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="w-1.5 h-1.5 bg-green-500 rounded-full"></span>
                <span>Статус: verified (сразу подтверждено)</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="w-1.5 h-1.5 bg-green-500 rounded-full"></span>
                <span>Генерация chip number</span>
              </div>
            </div>

            <div className="mt-6 text-green-600 font-medium group-hover:text-green-700">
              Создать карточку →
            </div>
          </button>
        </div>
      )}

      {mode === 'search' && (
        <div>
          <button
            onClick={() => setMode('select')}
            className="mb-6 text-gray-600 hover:text-gray-900 transition-colors flex items-center space-x-2"
          >
            <span>←</span>
            <span>Назад к выбору</span>
          </button>
          <SearchExistingCard />
        </div>
      )}

      {mode === 'create' && (
        <div>
          <button
            onClick={() => setMode('select')}
            className="mb-6 text-gray-600 hover:text-gray-900 transition-colors flex items-center space-x-2"
          >
            <span>←</span>
            <span>Назад к выбору</span>
          </button>
          <CreateNewCardForm />
        </div>
      )}
    </div>
  );
}
