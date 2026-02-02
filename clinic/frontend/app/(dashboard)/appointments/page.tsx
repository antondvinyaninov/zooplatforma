'use client';

import TableWidget from '../../components/admin/widgets/TableWidget';
import {
  CalendarIcon,
  PlusIcon,
} from '@heroicons/react/24/outline';

export default function AppointmentsPage() {
  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Записи на приём</h2>
        <p className="text-gray-600">Управление записями пациентов</p>
      </div>

      <TableWidget
        title="Записи (0)"
        actions={
          <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2">
            <PlusIcon className="w-4 h-4" />
            Создать запись
          </button>
        }
      >
        <div className="text-center py-12 text-gray-400">
          <CalendarIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>Записей пока нет</p>
          <p className="text-sm mt-2">Создайте первую запись на приём</p>
        </div>
      </TableWidget>
    </div>
  );
}
