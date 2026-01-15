# Design Document: Owner Cabinet

## Overview

The Owner Cabinet is a dedicated microservice that enables pet owners to manage their pets' digital passports (PetIDs), track medical information, and interact with other platform services. The service follows a multi-tenant architecture where each user can manage multiple pets, with all pet data centralized in the PetID service.

### Key Design Principles

1. **Single Source of Truth**: PetID service is the authoritative source for all pet data
2. **SSO Integration**: Authentication handled through Main service
3. **Read-Heavy Operations**: Most operations are reads from PetID service
4. **Event-Driven Updates**: Status changes trigger events and posts
5. **Progressive Enhancement**: Core functionality works without JavaScript

### Architecture Goals

- Maintain consistency with existing platform services (AdminLayout, styling)
- Minimize data duplication (rely on PetID service)
- Provide intuitive UX for non-technical users
- Support offline-first capabilities for viewing pet information

## Architecture

### System Context

```
┌─────────────────┐
│   Owner User    │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│      Owner Cabinet Service          │
│  Frontend (Next.js) - Port 6100     │
│  Backend (Go) - Port 8400           │
└──────┬──────────────────────────────┘
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
┌─────────────┐    ┌─────────────┐
│   PetID     │    │    Main     │
│  Service    │    │   Service   │
│ Port 8100   │    │  Port 8000  │
└─────────────┘    └─────────────┘
```

### Service Boundaries

**Owner Cabinet Responsibilities:**
- User interface for pet management
- Treatment tracking (owner-entered data)
- Status change workflows
- Photo upload and management
- Treatment reminders

**PetID Service Responsibilities:**
- Pet data storage (CRUD)
- Medical event storage
- Pet status management
- Photo storage
- Event history

**Main Service Responsibilities:**
- SSO authentication
- Post creation and display
- Social features (likes, comments)

### Communication Patterns

1. **Synchronous REST API**: Owner Cabinet ↔ PetID Service
2. **Synchronous REST API**: Owner Cabinet ↔ Main Service
3. **JWT Token Propagation**: All requests include user authentication

## Components and Interfaces

### Frontend Components

#### 1. Layout Component (`(dashboard)/layout.tsx`)

**Purpose**: Provides consistent navigation and authentication wrapper

**Props**: None (uses session context)

**State**:
- `isLoading: boolean` - authentication check in progress
- `user: User | null` - current authenticated user

**Behavior**:
- Checks authentication on mount
- Redirects to login if not authenticated
- Renders sidebar with 3 tabs: Overview, My Pets, Health
- Maintains active tab state

#### 2. PetList Component (`components/PetList.tsx`)

**Purpose**: Displays grid of pet cards with filtering

**Props**:
```typescript
interface PetListProps {
  pets: Pet[];
  onPetClick: (petId: number) => void;
  onAddPet: () => void;
}
```

**State**:
- `statusFilter: PetStatus | 'all'` - current filter selection
- `filteredPets: Pet[]` - pets matching current filter

**Behavior**:
- Filters pets by status when filter changes
- Displays placeholder when no pets exist
- Shows "Add Pet" button prominently

#### 3. PetCard Component (`components/PetCard.tsx`)

**Purpose**: Displays individual pet summary

**Props**:
```typescript
interface PetCardProps {
  pet: Pet;
  onClick: () => void;
}
```

**Behavior**:
- Shows pet photo (or placeholder)
- Displays name, species, breed, age
- Shows status badge with color coding
- Hover effect for interactivity

#### 4. PetRegistrationForm Component (`components/PetRegistrationForm.tsx`)

**Purpose**: Form for creating new PetID

**Props**:
```typescript
interface PetRegistrationFormProps {
  onSuccess: (petId: number) => void;
  onCancel: () => void;
}
```

**State**:
- `formData: Partial<Pet>` - form field values
- `errors: Record<string, string>` - validation errors
- `photoFile: File | null` - selected photo
- `isSubmitting: boolean` - submission in progress

**Validation Rules**:
- name: required, 1-50 characters
- species: required, one of predefined list
- age: required, 0-30 years
- sex: required, 'male' or 'female'
- color: required, 1-100 characters

**Behavior**:
- Validates on blur and submit
- Uploads photo separately after pet creation
- Shows success notification on completion
- Displays pending_verification status message

#### 5. PetEditForm Component (`components/PetEditForm.tsx`)

**Purpose**: Form for editing existing pet information

**Props**:
```typescript
interface PetEditFormProps {
  pet: Pet;
  onSuccess: () => void;
  onCancel: () => void;
}
```

**State**: Same as PetRegistrationForm plus:
- `originalData: Pet` - for change detection

**Behavior**:
- Pre-fills form with existing data
- Only submits changed fields
- Confirms before discarding changes

#### 6. MedicalHistory Component (`components/MedicalHistory.tsx`)

**Purpose**: Displays medical events and treatments

**Props**:
```typescript
interface MedicalHistoryProps {
  petId: number;
}
```

**State**:
- `events: MedicalEvent[]` - all medical events
- `treatments: Treatment[]` - owner-entered treatments
- `isLoading: boolean`

**Behavior**:
- Separates clinic events from owner treatments
- Displays events chronologically
- Shows read-only clinic events
- Allows editing/deleting owner treatments

#### 7. TreatmentForm Component (`components/TreatmentForm.tsx`)

**Purpose**: Form for adding flea/tick/worm treatments

**Props**:
```typescript
interface TreatmentFormProps {
  petId: number;
  onSuccess: () => void;
  onCancel: () => void;
}
```

**State**:
- `date: Date` - treatment date
- `medication: string` - medication name
- `dosage: string` - dosage amount
- `nextDate: Date | null` - optional next treatment date
- `notes: string` - optional notes

**Validation Rules**:
- date: required, not in future
- medication: required, 1-100 characters
- dosage: required, 1-50 characters
- nextDate: optional, must be after date

#### 8. StatusChangeModal Component (`components/StatusChangeModal.tsx`)

**Purpose**: Modal for changing pet status

**Props**:
```typescript
interface StatusChangeModalProps {
  pet: Pet;
  newStatus: PetStatus;
  onConfirm: (data: StatusChangeData) => void;
  onCancel: () => void;
}
```

**State**:
- `date: Date` - when status changed
- `location: string` - for "lost" status
- `circumstances: string` - description
- `contactInfo: string` - contact details

**Behavior**:
- Shows different fields based on status
- For "lost": requires location, circumstances, contact
- Prevents changing to "looking_for_home"
- Confirms before submitting

#### 9. TreatmentReminders Component (`components/TreatmentReminders.tsx`)

**Purpose**: Displays upcoming treatment reminders

**Props**:
```typescript
interface TreatmentRemindersProps {
  userId: number;
}
```

**State**:
- `reminders: TreatmentReminder[]` - upcoming treatments
- `isLoading: boolean`

**Behavior**:
- Fetches all pets for user
- Calculates days until next treatment
- Highlights overdue treatments in red
- Shows badge count on Health tab

### Backend API Endpoints

#### Pet Management

**GET /api/owner/pets**
- Returns all pets owned by authenticated user
- Query params: `status` (optional filter)
- Response: `{ pets: Pet[] }`

**POST /api/owner/pets**
- Creates new pet with pending_verification status
- Body: Pet data (without ID)
- Response: `{ pet: Pet, message: string }`

**GET /api/owner/pets/:id**
- Returns single pet details
- Validates ownership
- Response: `{ pet: Pet }`

**PUT /api/owner/pets/:id**
- Updates pet information
- Validates ownership
- Body: Partial pet data
- Response: `{ pet: Pet }`

**POST /api/owner/pets/:id/photo**
- Uploads pet photo
- Validates ownership
- Body: multipart/form-data with photo file
- Response: `{ photoUrl: string }`

#### Treatment Management

**GET /api/owner/pets/:id/treatments**
- Returns all owner-entered treatments for pet
- Validates ownership
- Response: `{ treatments: Treatment[] }`

**POST /api/owner/pets/:id/treatments**
- Creates new treatment record
- Validates ownership
- Body: Treatment data
- Response: `{ treatment: Treatment }`

**PUT /api/owner/pets/:id/treatments/:treatmentId**
- Updates treatment record
- Validates ownership
- Body: Partial treatment data
- Response: `{ treatment: Treatment }`

**DELETE /api/owner/pets/:id/treatments/:treatmentId**
- Deletes treatment record
- Validates ownership
- Response: `{ success: boolean }`

#### Medical History

**GET /api/owner/pets/:id/medical-history**
- Returns all medical events from clinics
- Validates ownership
- Response: `{ events: MedicalEvent[] }`

#### Status Management

**POST /api/owner/pets/:id/status**
- Changes pet status
- Validates ownership
- Prevents changing to "looking_for_home"
- Creates event in PetID service
- For "lost" status: creates post in Main service
- Body: `{ status: PetStatus, data: StatusChangeData }`
- Response: `{ pet: Pet, postId?: number }`

#### Reminders

**GET /api/owner/reminders**
- Returns upcoming treatment reminders for all user's pets
- Response: `{ reminders: TreatmentReminder[] }`

### Integration Points

#### PetID Service Integration

**Client**: `lib/petid-api.ts`

```typescript
class PetIDClient {
  async getPets(ownerId: number): Promise<Pet[]>
  async getPet(id: number): Promise<Pet>
  async createPet(data: CreatePetData): Promise<Pet>
  async updatePet(id: number, data: Partial<Pet>): Promise<Pet>
  async uploadPhoto(id: number, file: File): Promise<string>
  async getMedicalHistory(id: number): Promise<MedicalEvent[]>
  async createEvent(petId: number, event: EventData): Promise<MedicalEvent>
}
```

**Authentication**: JWT token in Authorization header

**Error Handling**: 
- 401: Redirect to login
- 403: Show "Access Denied" message
- 404: Show "Pet Not Found" message
- 500: Show "Service Unavailable" message

#### Main Service Integration

**Client**: `lib/main-api.ts`

```typescript
class MainClient {
  async createPost(data: CreatePostData): Promise<Post>
  async getPostsByPet(petId: number): Promise<Post[]>
}
```

**Usage**: 
- Create "lost" posts automatically
- Navigate to Main for manual post creation
- Display posts on pet profile

## Data Models

### Frontend TypeScript Types

```typescript
interface Pet {
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

type PetStatus = 'home' | 'lost' | 'found' | 'looking_for_home' | 'needs_help' | 'at_vet' | 'died';

interface Treatment {
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

interface MedicalEvent {
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

interface StatusChangeData {
  date: string;
  location?: string;
  circumstances?: string;
  contact_info?: string;
}

interface TreatmentReminder {
  pet_id: number;
  pet_name: string;
  treatment_id: number;
  medication: string;
  next_date: string;
  days_until: number;
  is_overdue: boolean;
}

interface User {
  id: number;
  email: string;
  name: string;
  avatar?: string;
}
```

### Backend Go Models

```go
type Pet struct {
    ID                 int       `json:"id"`
    OwnerID            int       `json:"owner_id"`
    CuratorID          *int      `json:"curator_id,omitempty"`
    OrganizationID     *int      `json:"organization_id,omitempty"`
    Name               string    `json:"name"`
    Species            string    `json:"species"`
    Breed              *string   `json:"breed,omitempty"`
    Age                int       `json:"age"`
    Sex                string    `json:"sex"`
    Color              string    `json:"color"`
    Weight             *float64  `json:"weight,omitempty"`
    ChipNumber         *string   `json:"chip_number,omitempty"`
    TattooNumber       *string   `json:"tattoo_number,omitempty"`
    PassportNumber     *string   `json:"passport_number,omitempty"`
    PhotoURL           *string   `json:"photo_url,omitempty"`
    Status             string    `json:"status"`
    VerificationStatus string    `json:"verification_status"`
    Sterilized         bool      `json:"sterilized"`
    SterilizationDate  *string   `json:"sterilization_date,omitempty"`
    SpecialMarks       *string   `json:"special_marks,omitempty"`
    Character          *string   `json:"character,omitempty"`
    HealthNotes        *string   `json:"health_notes,omitempty"`
    Allergies          *string   `json:"allergies,omitempty"`
    ChronicConditions  *string   `json:"chronic_conditions,omitempty"`
    EmergencyContact   *string   `json:"emergency_contact,omitempty"`
    EmergencyPhone     *string   `json:"emergency_phone,omitempty"`
    InsuranceCompany   *string   `json:"insurance_company,omitempty"`
    InsurancePolicy    *string   `json:"insurance_policy,omitempty"`
    City               *string   `json:"city,omitempty"`
    Region             *string   `json:"region,omitempty"`
    Urgent             *bool     `json:"urgent,omitempty"`
    ContactPhone       *string   `json:"contact_phone,omitempty"`
    ContactName        *string   `json:"contact_name,omitempty"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}

type Treatment struct {
    ID         int       `json:"id"`
    PetID      int       `json:"pet_id"`
    UserID     int       `json:"user_id"`
    Date       time.Time `json:"date"`
    Medication string    `json:"medication"`
    Dosage     string    `json:"dosage"`
    NextDate   *time.Time `json:"next_date,omitempty"`
    Notes      *string   `json:"notes,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type StatusChangeRequest struct {
    Status        string  `json:"status"`
    Date          string  `json:"date"`
    Location      *string `json:"location,omitempty"`
    Circumstances *string `json:"circumstances,omitempty"`
    ContactInfo   *string `json:"contact_info,omitempty"`
}
```

### Database Schema

**Note**: Pet data is stored in PetID service. Owner Cabinet only stores treatments.

**treatments table** (Owner Cabinet database):
```sql
CREATE TABLE treatments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    date DATETIME NOT NULL,
    medication TEXT NOT NULL,
    dosage TEXT NOT NULL,
    next_date DATETIME,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_treatments_pet_id ON treatments(pet_id);
CREATE INDEX idx_treatments_user_id ON treatments(user_id);
CREATE INDEX idx_treatments_next_date ON treatments(next_date);
```

**Rationale**: Treatments are owner-specific data that don't need to be in the centralized PetID service. This allows Owner Cabinet to operate independently for treatment tracking while still relying on PetID for core pet data.


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Required Field Validation

*For any* pet registration form submission with missing required fields (name, species, age, sex, or color), the validation should fail and prevent submission.

**Validates: Requirements 1.2**

### Property 2: Pending Verification Status

*For any* newly created PetID, the verification_status should be set to "pending_verification".

**Validates: Requirements 1.3**

### Property 3: Owner ID Assignment

*For any* authenticated user creating a pet, the created pet's owner_id should equal the user's ID.

**Validates: Requirements 1.4**

### Property 4: Photo Upload Round Trip

*For any* valid image file uploaded as a pet photo, the photo should be stored and retrievable via the returned photo_url.

**Validates: Requirements 1.6**

### Property 5: Pet List Ownership Filter

*For any* user viewing their pet list, all returned pets should have owner_id equal to that user's ID.

**Validates: Requirements 2.1, 6.4**

### Property 6: Pet Card Information Completeness

*For any* pet displayed in a pet card, the rendered card should contain the pet's photo (or placeholder), name, species, breed, age, and status.

**Validates: Requirements 2.2**

### Property 7: Pet Edit Round Trip

*For any* pet being edited, saving changes then retrieving the pet should reflect all the changes made.

**Validates: Requirements 2.6**

### Property 8: Status Filter Correctness

*For any* status filter applied to the pet list, all displayed pets should have a status matching the selected filter.

**Validates: Requirements 2.5**

### Property 9: Medical History Completeness

*For any* pet with medical events, viewing the medical history should retrieve all events associated with that pet_id.

**Validates: Requirements 3.1**

### Property 10: Medical Event Separation

*For any* medical history display, events with clinic_id not null should be shown separately from treatments with clinic_id null.

**Validates: Requirements 3.2**

### Property 11: Treatment Required Fields

*For any* treatment submission with missing required fields (date, medication, or dosage), the validation should fail and prevent submission.

**Validates: Requirements 3.3**

### Property 12: Treatment Optional Next Date

*For any* treatment submission, the system should accept treatments both with and without a next_date field.

**Validates: Requirements 3.4**

### Property 13: Treatment Event Type

*For any* saved treatment, the created event in PetID_Service should have event_type set to "treatment".

**Validates: Requirements 3.5**

### Property 14: Chronological Treatment Ordering

*For any* list of treatments or upcoming reminders, they should be sorted in chronological order by date.

**Validates: Requirements 3.6, 10.3**

### Property 15: Treatment Reminder Threshold

*For any* treatment with next_date within 7 days from today, a reminder should be displayed on the Health tab.

**Validates: Requirements 3.7, 10.2**

### Property 16: Clinic Event Read-Only

*For any* medical event with clinic_id not null, the display should not include edit or delete controls.

**Validates: Requirements 3.8**

### Property 17: Lost Status Post Creation

*For any* pet status change to "lost", an automatic post should be created in Main_Service with the "Lost" tag.

**Validates: Requirements 4.3**

### Property 18: Looking For Home Prevention

*For any* attempt to change a pet's status to "looking_for_home", the system should reject the request and return an error.

**Validates: Requirements 4.4**

### Property 19: Status Change Event Creation

*For any* pet status change, an event should be created in PetID_Service recording the status change.

**Validates: Requirements 4.2, 4.6**

### Property 20: Status Change Confirmation

*For any* completed status change, a confirmation notification should be displayed to the user.

**Validates: Requirements 4.5**

### Property 21: Pet ID Parameter Passing

*For any* navigation to Main_Service post creation from a pet profile, the pet_id should be included as a parameter.

**Validates: Requirements 5.2**

### Property 22: Pet Posts Retrieval

*For any* pet, viewing the pet's posts should retrieve all posts from Main_Service where the post is tagged with that pet_id.

**Validates: Requirements 5.3**

### Property 23: Veterinary Visit Filtering

*For any* pet's veterinary visit history, only medical events with clinic_id not null should be displayed.

**Validates: Requirements 5.4**

### Property 24: Clinic Record Information

*For any* clinic event displayed, the rendered record should include clinic name, date, event type, and description.

**Validates: Requirements 5.5**

### Property 25: Unauthenticated Redirect

*For any* request to a protected route without valid authentication, the system should redirect to the SSO login page.

**Validates: Requirements 6.2**

### Property 26: Overview Statistics Display

*For any* user accessing the Overview tab, summary statistics of their pets should be displayed.

**Validates: Requirements 6.3**

### Property 27: Health Tab Content

*For any* user accessing the Health tab, upcoming treatment reminders and recent medical events should be displayed.

**Validates: Requirements 6.5**

### Property 28: Validation Error Display

*For any* form submission with validation errors, field-specific error messages should be displayed for each invalid field.

**Validates: Requirements 7.1, 7.5**

### Property 29: API Error Handling

*For any* failed API request, a user-friendly error message should be displayed to the user.

**Validates: Requirements 7.2**

### Property 30: Photo Upload Error Recovery

*For any* failed photo upload, an error message should be displayed and a retry option should be available.

**Validates: Requirements 7.3**

### Property 31: JWT Token Verification

*For any* access to a protected route, the JWT token should be verified before allowing access.

**Validates: Requirements 8.2**

### Property 32: Invalid Token Redirect

*For any* request with an invalid or expired JWT token, the system should redirect to the login page.

**Validates: Requirements 8.3**

### Property 33: Unauthorized Pet Access

*For any* attempt to edit a pet not owned by the authenticated user, the system should return a 403 Forbidden error.

**Validates: Requirements 8.4**

### Property 34: JWT Token Inclusion

*For any* API request to PetID_Service, the JWT token should be included in the Authorization header.

**Validates: Requirements 8.5**

### Property 35: Image Format Validation

*For any* file upload attempt with a non-image format (not jpg, png, or webp), the upload should be rejected with a validation error.

**Validates: Requirements 9.1**

### Property 36: Image Size Validation

*For any* file upload attempt with size greater than 5MB, the upload should be rejected with a validation error.

**Validates: Requirements 9.2**

### Property 37: Photo Upload Preview

*For any* successful photo upload, a preview of the uploaded photo should be displayed.

**Validates: Requirements 9.3**

### Property 38: Photo Replacement

*For any* pet photo change, the new photo should replace the existing photo in the pet's record.

**Validates: Requirements 9.4**

### Property 39: Placeholder Image Display

*For any* pet without a photo_url, a placeholder image should be displayed in the pet card.

**Validates: Requirements 9.5**

### Property 40: Days Until Treatment Calculation

*For any* treatment with a next_treatment_date, the days_until value should be calculated as the difference between next_treatment_date and today.

**Validates: Requirements 10.1**

### Property 41: Overdue Treatment Highlighting

*For any* treatment with next_treatment_date in the past, the treatment should be highlighted in red.

**Validates: Requirements 10.4**

### Property 42: Treatment Completion Workflow

*For any* treatment marked as done, the system should allow optionally scheduling the next treatment.

**Validates: Requirements 10.5**

## Error Handling

### Client-Side Error Handling

**Form Validation Errors**:
- Display inline error messages below invalid fields
- Prevent form submission until all errors are resolved
- Highlight invalid fields with red border
- Show summary of errors at top of form

**Network Errors**:
- Display toast notification with error message
- Provide retry button for failed requests
- Show loading state during retry
- Log errors to console for debugging

**Authentication Errors**:
- 401 Unauthorized: Redirect to login page
- 403 Forbidden: Show "Access Denied" message
- Clear local storage and session data

**File Upload Errors**:
- Invalid format: "Please upload a JPG, PNG, or WebP image"
- File too large: "Image must be less than 5MB"
- Upload failed: "Upload failed. Please try again"
- Show progress bar during upload

### Server-Side Error Handling

**Validation Errors**:
- Return 400 Bad Request with detailed error messages
- Include field names in error response
- Example: `{ "error": "Validation failed", "fields": { "name": "Name is required" } }`

**Authentication Errors**:
- Return 401 Unauthorized for missing/invalid token
- Return 403 Forbidden for insufficient permissions
- Include WWW-Authenticate header

**Resource Errors**:
- Return 404 Not Found for non-existent pets
- Return 409 Conflict for duplicate operations
- Include resource ID in error message

**Integration Errors**:
- Retry failed requests to PetID_Service (max 3 attempts)
- Return 503 Service Unavailable if PetID_Service is down
- Log integration errors for monitoring
- Provide fallback behavior where possible

**Database Errors**:
- Return 500 Internal Server Error for database failures
- Log full error details server-side
- Return generic message to client
- Implement transaction rollback for data consistency

### Error Logging

**Client-Side**:
- Log errors to browser console
- Send critical errors to monitoring service
- Include user ID, timestamp, and error context

**Server-Side**:
- Log all errors to file with timestamp
- Include request ID for tracing
- Log stack traces for debugging
- Monitor error rates and alert on spikes

## Testing Strategy

### Dual Testing Approach

The Owner Cabinet will use both unit tests and property-based tests to ensure comprehensive coverage:

**Unit Tests**: Verify specific examples, edge cases, and error conditions
- Test specific form validation scenarios
- Test authentication redirect flows
- Test error message display
- Test component rendering with specific props
- Test API integration with mocked responses

**Property Tests**: Verify universal properties across all inputs
- Test validation logic with randomly generated form data
- Test filtering and sorting with random pet lists
- Test ownership verification with random user/pet combinations
- Test date calculations with random dates
- Test file validation with random file types and sizes

Both testing approaches are complementary and necessary for comprehensive coverage. Unit tests catch concrete bugs in specific scenarios, while property tests verify general correctness across a wide range of inputs.

### Property-Based Testing Configuration

**Testing Library**: fast-check (JavaScript/TypeScript property-based testing library)

**Configuration**:
- Minimum 100 iterations per property test
- Each property test must reference its design document property
- Tag format: `Feature: owner-cabinet, Property {number}: {property_text}`

**Example Property Test**:
```typescript
import fc from 'fast-check';

// Feature: owner-cabinet, Property 1: Required Field Validation
test('pet registration requires all mandatory fields', () => {
  fc.assert(
    fc.property(
      fc.record({
        name: fc.option(fc.string(), { nil: undefined }),
        species: fc.option(fc.constantFrom('dog', 'cat', 'bird'), { nil: undefined }),
        age: fc.option(fc.integer({ min: 0, max: 30 }), { nil: undefined }),
        sex: fc.option(fc.constantFrom('male', 'female'), { nil: undefined }),
        color: fc.option(fc.string(), { nil: undefined }),
      }),
      (formData) => {
        const hasAllRequired = 
          formData.name && 
          formData.species && 
          formData.age !== undefined && 
          formData.sex && 
          formData.color;
        
        const validationResult = validatePetForm(formData);
        
        if (hasAllRequired) {
          expect(validationResult.isValid).toBe(true);
        } else {
          expect(validationResult.isValid).toBe(false);
          expect(validationResult.errors).toBeDefined();
        }
      }
    ),
    { numRuns: 100 }
  );
});
```

### Unit Testing Examples

**Component Tests** (React Testing Library):
```typescript
describe('PetCard', () => {
  it('displays placeholder when pet has no photo', () => {
    const pet = { ...mockPet, photo_url: null };
    render(<PetCard pet={pet} onClick={jest.fn()} />);
    expect(screen.getByAltText('Pet placeholder')).toBeInTheDocument();
  });

  it('displays pet photo when available', () => {
    const pet = { ...mockPet, photo_url: 'https://example.com/photo.jpg' };
    render(<PetCard pet={pet} onClick={jest.fn()} />);
    expect(screen.getByAltText(pet.name)).toHaveAttribute('src', pet.photo_url);
  });
});
```

**API Integration Tests**:
```typescript
describe('PetID API Client', () => {
  it('includes JWT token in all requests', async () => {
    const mockFetch = jest.fn().mockResolvedValue({ ok: true, json: async () => ({}) });
    global.fetch = mockFetch;
    
    const client = new PetIDClient('test-token');
    await client.getPets(1);
    
    expect(mockFetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          'Authorization': 'Bearer test-token'
        })
      })
    );
  });
});
```

**Backend Handler Tests** (Go):
```go
func TestGetPets_OnlyReturnsOwnedPets(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    handler := NewPetHandler(db)
    
    // Create pets for different users
    createTestPet(db, 1, 100) // user 1, pet 100
    createTestPet(db, 2, 101) // user 2, pet 101
    
    // Test
    req := httptest.NewRequest("GET", "/api/owner/pets", nil)
    req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
    w := httptest.NewRecorder()
    
    handler.GetPets(w, req)
    
    // Assert
    var response struct {
        Pets []Pet `json:"pets"`
    }
    json.NewDecoder(w.Body).Decode(&response)
    
    assert.Equal(t, 1, len(response.Pets))
    assert.Equal(t, 100, response.Pets[0].ID)
    assert.Equal(t, 1, response.Pets[0].OwnerID)
}
```

### Test Coverage Goals

- **Unit Test Coverage**: >80% line coverage
- **Property Test Coverage**: All correctness properties implemented
- **Integration Test Coverage**: All API endpoints tested
- **E2E Test Coverage**: Critical user flows (register pet, add treatment, change status)

### Testing Priorities

1. **Critical Path**: Pet registration, authentication, ownership verification
2. **Data Integrity**: Status changes, treatment tracking, photo uploads
3. **Integration**: PetID service communication, Main service integration
4. **Error Handling**: Validation errors, network errors, authentication errors
5. **UI/UX**: Form validation, error messages, loading states
