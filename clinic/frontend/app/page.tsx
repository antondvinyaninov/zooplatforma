'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    // Проверяем есть ли выбранная клиника
    const clinicId = localStorage.getItem('selectedClinicId');
    if (clinicId) {
      // Если клиника выбрана - идём в dashboard
      router.push('/overview');
    } else {
      // Если нет - на выбор клиники
      router.push('/select');
    }
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
      <div className="text-center">
        <div className="text-lg text-gray-600">Перенаправление...</div>
      </div>
    </div>
  );
}
