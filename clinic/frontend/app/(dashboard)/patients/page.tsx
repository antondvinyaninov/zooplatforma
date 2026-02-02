'use client';

import { useState, useEffect } from 'react';
import TableWidget from '../../components/admin/widgets/TableWidget';
import {
  UserGroupIcon,
} from '@heroicons/react/24/outline';

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string;
  photo?: string;
  status: string;
  created_at: string;
}

export default function PatientsPage() {
  const [patients, setPatients] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPatients();
  }, []);

  const loadPatients = async () => {
    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (clinicId) {
        const response = await fetch('http://localhost:8600/api/my-patients', {
          method: 'GET',
          credentials: 'include',
          headers: {
            'X-Clinic-ID': clinicId,
          },
        });

        if (response.ok) {
          const result = await response.json();
          if (result.success && result.data) {
            setPatients(result.data);
          }
        }
      }
    } catch (error) {
      console.error('Failed to load patients:', error);
    }

    setLoading(false);
  };

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Пациенты</h2>
        <p className="text-gray-600">Все животные, которые были на приёме</p>
      </div>

      <TableWidget
        title={`Пациентов (${patients.length})`}
      >
        {loading ? (
          <div className="text-center py-12 text-gray-400">Загрузка...</div>
        ) : patients.length === 0 ? (
          <div className="text-center py-12 text-gray-400">
            <UserGroupIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
            <p>Пациентов пока нет</p>
            <p className="text-sm mt-2">Добавьте первого пациента</p>
          </div>
        ) : (
          <div className="overflow-x-auto p-6">
            <table className="w-full min-w-[800px]">
              <thead>
                <tr className="border-b border-gray-100">
                  <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Фото</th>
                  <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Кличка</th>
                  <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Вид</th>
                  <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Порода</th>
                  <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Статус</th>
                  <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Действия</th>
                </tr>
              </thead>
              <tbody>
                {patients.map((pet) => (
                  <tr key={pet.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                    <td className="py-3 px-4">
                      {pet.photo ? (
                        <img
                          src={`http://localhost:8100${pet.photo}`}
                          alt={pet.name}
                          className="w-12 h-12 rounded-full object-cover"
                        />
                      ) : (
                        <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-xl">
                          🐾
                        </div>
                      )}
                    </td>
                    <td className="py-3 px-4 text-sm font-medium text-gray-900">{pet.name}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{pet.species}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{pet.breed || '-'}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{pet.status}</td>
                    <td className="py-3 px-4 text-sm">
                      <button
                        onClick={() => window.location.href = `http://localhost:3000/pets/${pet.id}`}
                        className="text-blue-600 hover:text-blue-700"
                      >
                        Открыть
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </TableWidget>
    </div>
  );
}
