'use client';

import { HeartIcon, PlusIcon } from '@heroicons/react/24/outline';

export default function HealthPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Здоровье питомцев</h2>
        <p className="text-gray-600">Медицинская история и визиты к ветеринару</p>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b border-gray-200">
          <div className="flex justify-between items-center">
            <h3 className="text-xl font-bold text-gray-900">Визиты к ветеринару (0)</h3>
            <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2">
              <PlusIcon className="w-4 h-4" />
              Добавить визит
            </button>
          </div>
        </div>

        <div className="text-center py-12 text-gray-400">
          <HeartIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>Визитов пока нет</p>
          <p className="text-sm mt-2">Добавьте информацию о визите к ветеринару</p>
        </div>
      </div>
    </div>
  );
}
