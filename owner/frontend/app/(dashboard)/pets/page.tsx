'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import PetList from '../../../components/PetList';
import AddPetModal from '../../../components/AddPetModal';
import petidApi from '../../../lib/petid-api';
import { Pet, CreatePetData } from '../../../types/owner';

export default function PetsListPage() {
  const router = useRouter();
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  useEffect(() => {
    loadPets();
  }, []);

  const loadPets = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await petidApi.getPets();
      
      if (response.error) {
        setError(response.error);
      } else if (response.data) {
        setPets(response.data);
      }
    } catch (err) {
      console.error('Error loading pets:', err);
      setError('Не удалось загрузить список питомцев');
    } finally {
      setLoading(false);
    }
  };

  const handlePetClick = (petId: number) => {
    router.push(`/pets/${petId}`);
  };

  const handleAddPet = () => {
    setIsAddModalOpen(true);
  };

  const handleSubmitPet = async (data: CreatePetData) => {
    const response = await petidApi.createPet(data);
    
    if (response.error) {
      throw new Error(response.error);
    }
    
    // Перезагружаем список питомцев
    await loadPets();
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Загрузка питомцев...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
        <div className="text-4xl mb-3">⚠️</div>
        <h3 className="text-lg font-semibold text-red-900 mb-2">Ошибка загрузки</h3>
        <p className="text-red-700 mb-4">{error}</p>
        <button
          onClick={loadPets}
          className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
        >
          Попробовать снова
        </button>
      </div>
    );
  }

  return (
    <>
      <PetList
        pets={pets}
        onPetClick={handlePetClick}
        onAddPet={handleAddPet}
      />
      
      <AddPetModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onSubmit={handleSubmitPet}
      />
    </>
  );
}
