// TypeScript types for Owner Cabinet

export interface Pet {
  id: number;
  owner_id: number;
  curator_id?: number;
  organization_id?: number;
  name: string;
  species: 'dog' | 'cat' | 'bird' | 'rodent' | 'reptile' | 'fish' | 'other';
  breed?: string;
  age: number;
  sex: 'male' | 'female';
  color: string;
  weight?: number;
  chip_number?: string;
  tattoo_number?: string;
  passport_number?: string;
  photo_url?: string;
  status: PetStatus;
  verification_status: 'pending_verification' | 'verified';
  sterilized: boolean;
  sterilization_date?: string;
  special_marks?: string;
  character?: string;
  health_notes?: string;
  allergies?: string;
  chronic_conditions?: string;
  emergency_contact?: string;
  emergency_phone?: string;
  insurance_company?: string;
  insurance_policy?: string;
  city?: string;
  region?: string;
  urgent?: boolean;
  contact_phone?: string;
  contact_name?: string;
  created_at: string;
  updated_at: string;
}

export type PetStatus = 'home' | 'lost' | 'found' | 'looking_for_home' | 'needs_help' | 'at_vet' | 'died';

export interface Treatment {
  id: number;
  pet_id: number;
  user_id: number;
  date: string;
  medication: string;
  dosage: string;
  next_date?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface MedicalEvent {
  id: number;
  pet_id: number;
  event_type: 'registration' | 'vaccination' | 'sterilization' | 'treatment' | 'surgery' | 'examination';
  user_id: number;
  clinic_id?: number;
  organization_id?: number;
  event_date: string;
  description: string;
  veterinarian?: string;
  next_visit_date?: string;
  created_at: string;
}

export interface StatusChangeData {
  date: string;
  location?: string;
  circumstances?: string;
  contact_info?: string;
}

export interface TreatmentReminder {
  pet_id: number;
  pet_name: string;
  treatment_id: number;
  medication: string;
  next_date: string;
  days_until: number;
  is_overdue: boolean;
}

export interface User {
  id: number;
  email: string;
  name: string;
  avatar?: string;
}

// Form data types
export interface CreatePetData {
  name: string;
  species: Pet['species'];
  breed?: string;
  age: number;
  sex: Pet['sex'];
  color: string;
  weight?: number;
  chip_number?: string;
  tattoo_number?: string;
  passport_number?: string;
  sterilized?: boolean;
  sterilization_date?: string;
  special_marks?: string;
  character?: string;
  health_notes?: string;
  allergies?: string;
  chronic_conditions?: string;
  emergency_contact?: string;
  emergency_phone?: string;
  insurance_company?: string;
  insurance_policy?: string;
  city?: string;
  region?: string;
}

export interface CreateTreatmentData {
  pet_id: number;
  date: string;
  medication: string;
  dosage: string;
  next_date?: string;
  notes?: string;
}

export interface StatusChangeRequest {
  status: PetStatus;
  date: string;
  location?: string;
  circumstances?: string;
  contact_info?: string;
}

// API Response types
export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PetsResponse {
  pets: Pet[];
}

export interface TreatmentsResponse {
  treatments: Treatment[];
}

export interface MedicalHistoryResponse {
  events: MedicalEvent[];
}

export interface RemindersResponse {
  reminders: TreatmentReminder[];
}
