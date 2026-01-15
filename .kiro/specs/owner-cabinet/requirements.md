# Requirements Document: Owner Cabinet

## Introduction

The Owner Cabinet is a dedicated microservice for pet owners to manage their pets' information, medical records, and status updates. It provides a centralized interface for owners to register new pets, track medical history, manage pet statuses, and integrate with other platform services.

## Glossary

- **Owner**: A registered user who owns one or more pets
- **PetID**: The unique identifier and digital passport for each pet in the system
- **Owner_Cabinet**: The microservice running on localhost:6100 (frontend) and localhost:8400 (backend)
- **PetID_Service**: The centralized pet database service (localhost:4100/8100)
- **Main_Service**: The main social network service (localhost:3000/8000)
- **Clinic_Service**: The veterinary clinic service (localhost:6300/8600)
- **Treatment**: Medical procedures performed by owners (flea/tick/worm treatments)
- **Medical_Event**: Medical procedures performed by veterinary clinics (vaccinations, surgeries, examinations)
- **Pet_Status**: The current state of a pet (home, lost, at_vet, needs_help)
- **Verification_Status**: The confirmation state of a PetID (pending_verification, verified)

## Requirements

### Requirement 1: Pet Registration

**User Story:** As an owner, I want to register my pet in the system, so that I can create a digital passport and track their information.

#### Acceptance Criteria

1. WHEN an owner accesses the registration form, THE Owner_Cabinet SHALL display all available pet fields for input
2. WHEN an owner submits the registration form, THE Owner_Cabinet SHALL validate that required fields (name, species, age, sex, color) are filled
3. WHEN validation passes, THE Owner_Cabinet SHALL create a new PetID with status pending_verification
4. WHEN a PetID is created, THE Owner_Cabinet SHALL set the owner_id to the current user's ID
5. WHEN a PetID is created, THE Owner_Cabinet SHALL display a notification stating "PetID created, awaiting veterinary clinic confirmation"
6. WHEN an owner uploads a pet photo, THE Owner_Cabinet SHALL store the photo and associate it with the PetID

### Requirement 2: Pet Management

**User Story:** As an owner, I want to view and manage all my pets, so that I can keep their information up to date.

#### Acceptance Criteria

1. WHEN an owner views the pet list, THE Owner_Cabinet SHALL retrieve all pets where owner_id equals the current user's ID from PetID_Service
2. WHEN displaying pet cards, THE Owner_Cabinet SHALL show photo, name, species, breed, age, and current status for each pet
3. WHEN an owner clicks "Add Pet", THE Owner_Cabinet SHALL navigate to the registration form
4. WHEN an owner clicks "Edit" on a pet card, THE Owner_Cabinet SHALL navigate to the pet editing form with pre-filled data
5. WHEN an owner applies status filters, THE Owner_Cabinet SHALL display only pets matching the selected status (home, lost, at_vet)
6. WHEN an owner saves pet edits, THE Owner_Cabinet SHALL update the pet information in PetID_Service

### Requirement 3: Medical Information Management

**User Story:** As an owner, I want to manage my pet's medical information, so that I can track treatments and view veterinary records.

#### Acceptance Criteria

1. WHEN an owner views medical history, THE Owner_Cabinet SHALL retrieve all medical events for the pet from PetID_Service
2. WHEN displaying medical history, THE Owner_Cabinet SHALL show events created by veterinary clinics separately from owner-entered treatments
3. WHEN an owner adds a treatment, THE Owner_Cabinet SHALL require date, medication name, and dosage
4. WHEN an owner adds a treatment, THE Owner_Cabinet SHALL optionally accept a next treatment date
5. WHEN a treatment is saved, THE Owner_Cabinet SHALL create a treatment event in PetID_Service with event_type "treatment"
6. WHEN displaying treatment history, THE Owner_Cabinet SHALL show all treatments in chronological order with dates
7. WHEN a next treatment date is approaching, THE Owner_Cabinet SHALL display a reminder notification
8. WHEN viewing medical events from clinics, THE Owner_Cabinet SHALL display them as read-only

### Requirement 4: Pet Status Management

**User Story:** As an owner, I want to change my pet's status, so that I can indicate when my pet is lost or found.

#### Acceptance Criteria

1. WHEN an owner changes status to "lost", THE Owner_Cabinet SHALL display a form requesting date, location, circumstances, and contact information
2. WHEN an owner submits a "lost" status change, THE Owner_Cabinet SHALL update the pet status in PetID_Service
3. WHEN a pet status changes to "lost", THE Owner_Cabinet SHALL create an automatic post in Main_Service with the "Lost" tag
4. WHEN an owner attempts to change status to "looking_for_home", THE Owner_Cabinet SHALL prevent the action and display an error message
5. WHEN a status change is completed, THE Owner_Cabinet SHALL display a confirmation notification
6. WHEN a pet status changes, THE Owner_Cabinet SHALL create a status change event in PetID_Service

### Requirement 5: Integration with Other Services

**User Story:** As an owner, I want to interact with other platform services, so that I can share posts about my pets and view their complete history.

#### Acceptance Criteria

1. WHEN an owner clicks "Create Post", THE Owner_Cabinet SHALL provide a button that navigates to Main_Service post creation
2. WHEN navigating to post creation, THE Owner_Cabinet SHALL pass the pet_id as a parameter to Main_Service
3. WHEN an owner views a pet's posts, THE Owner_Cabinet SHALL retrieve all posts tagged with the pet_id from Main_Service
4. WHEN displaying veterinary visit history, THE Owner_Cabinet SHALL retrieve all medical events where clinic_id is not null from PetID_Service
5. WHEN an owner views clinic records, THE Owner_Cabinet SHALL display the clinic name, date, event type, and description

### Requirement 6: User Interface and Navigation

**User Story:** As an owner, I want an intuitive interface, so that I can easily navigate and manage my pets.

#### Acceptance Criteria

1. THE Owner_Cabinet SHALL display a sidebar with three tabs: Overview, My Pets, Health
2. WHEN an owner is not authenticated, THE Owner_Cabinet SHALL redirect to the SSO login page
3. WHEN an owner accesses the Overview tab, THE Owner_Cabinet SHALL display summary statistics of their pets
4. WHEN an owner accesses the My Pets tab, THE Owner_Cabinet SHALL display the list of all their pets
5. WHEN an owner accesses the Health tab, THE Owner_Cabinet SHALL display upcoming treatment reminders and recent medical events
6. THE Owner_Cabinet SHALL maintain consistent styling with other platform services using AdminLayout

### Requirement 7: Data Validation and Error Handling

**User Story:** As an owner, I want clear error messages, so that I can correct any mistakes when entering pet information.

#### Acceptance Criteria

1. WHEN required fields are missing, THE Owner_Cabinet SHALL display field-specific error messages
2. WHEN an API request fails, THE Owner_Cabinet SHALL display a user-friendly error message
3. WHEN a pet photo upload fails, THE Owner_Cabinet SHALL display an error and allow retry
4. WHEN network connectivity is lost, THE Owner_Cabinet SHALL display a connection error message
5. WHEN invalid data is entered, THE Owner_Cabinet SHALL prevent form submission and highlight invalid fields

### Requirement 8: Authentication and Authorization

**User Story:** As an owner, I want secure access to my pet information, so that only I can manage my pets.

#### Acceptance Criteria

1. THE Owner_Cabinet SHALL use SSO authentication through Main_Service
2. WHEN an owner accesses any protected route, THE Owner_Cabinet SHALL verify the JWT token
3. WHEN a JWT token is invalid or expired, THE Owner_Cabinet SHALL redirect to the login page
4. WHEN an owner attempts to edit another user's pet, THE Owner_Cabinet SHALL return a 403 Forbidden error
5. THE Owner_Cabinet SHALL include the JWT token in all API requests to PetID_Service

### Requirement 9: Pet Photo Management

**User Story:** As an owner, I want to upload and manage photos of my pets, so that their profiles are complete and visually appealing.

#### Acceptance Criteria

1. WHEN an owner uploads a pet photo, THE Owner_Cabinet SHALL validate that the file is an image format (jpg, png, webp)
2. WHEN an owner uploads a pet photo, THE Owner_Cabinet SHALL validate that the file size is less than 5MB
3. WHEN a photo upload is successful, THE Owner_Cabinet SHALL display a preview of the uploaded photo
4. WHEN an owner changes a pet photo, THE Owner_Cabinet SHALL replace the existing photo with the new one
5. WHEN displaying pet cards, THE Owner_Cabinet SHALL show a placeholder image if no photo exists

### Requirement 10: Treatment Reminders

**User Story:** As an owner, I want to receive reminders for upcoming treatments, so that I don't forget to treat my pet.

#### Acceptance Criteria

1. WHEN a treatment has a next_treatment_date, THE Owner_Cabinet SHALL calculate days until the next treatment
2. WHEN the next treatment date is within 7 days, THE Owner_Cabinet SHALL display a reminder badge on the Health tab
3. WHEN an owner views the Health tab, THE Owner_Cabinet SHALL list all upcoming treatments sorted by date
4. WHEN a treatment is overdue, THE Owner_Cabinet SHALL highlight it in red
5. WHEN an owner completes a treatment, THE Owner_Cabinet SHALL allow marking it as done and optionally scheduling the next one
