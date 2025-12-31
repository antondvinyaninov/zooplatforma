const API_URL = 'http://localhost:8400';

export const api = {
  async getMyPets() {
    const response = await fetch(`${API_URL}/api/my-pets`, {
      method: 'GET',
      credentials: 'include',
    });
    return response.json();
  },

  async getPetEvents(petId: number) {
    const response = await fetch(`${API_URL}/api/pets/events?petId=${petId}`, {
      method: 'GET',
      credentials: 'include',
    });
    return response.json();
  },

  async getProfile() {
    const response = await fetch(`${API_URL}/api/profile`, {
      method: 'GET',
      credentials: 'include',
    });
    return response.json();
  },
};
