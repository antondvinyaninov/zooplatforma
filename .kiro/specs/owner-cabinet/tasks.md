# Implementation Plan: Owner Cabinet

## Overview

This implementation plan breaks down the Owner Cabinet feature into discrete, incremental tasks. The microservice already exists with basic infrastructure (ports 6100/8400, SSO integration, AdminLayout with 3 tabs). We will build upon this foundation to implement pet management, medical tracking, and status management features.

## Tasks

- [x] 1. Set up database and core data models
  - Create treatments table migration in Owner backend
  - Define Go structs for Pet, Treatment, MedicalEvent, StatusChangeRequest
  - Define TypeScript interfaces in frontend types file
  - Set up database connection and migration runner
  - _Requirements: 1.1, 1.2, 3.3, 10.1_

- [ ]* 1.1 Write property test for database schema
  - **Property 1: Required Field Validation**
  - **Validates: Requirements 1.2**

- [ ] 2. Implement PetID service integration
  - [x] 2.1 Create PetID API client (lib/petid-api.ts)
    - Implement getPets, getPet, createPet, updatePet methods
    - Implement uploadPhoto, getMedicalHistory, createEvent methods
    - Add JWT token to all requests (Authorization header)
    - Handle error responses (401, 403, 404, 500)
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 8.5_

- [ ]* 2.2 Write property tests for PetID client
  - **Property 34: JWT Token Inclusion**
  - **Validates: Requirements 8.5**

- [ ]* 2.3 Write unit tests for PetID client error handling
  - Test 401 redirect to login
  - Test 403 access denied message
  - Test 404 not found message
  - Test 500 service unavailable message
  - _Requirements: 7.2_

- [x] 3. Implement pet list and filtering
  - [x] 3.1 Create PetList component
    - Display grid of pet cards (3 columns)
    - Implement status filter dropdown (all, home, lost, at_vet)
    - Show "Add Pet" button
    - Show placeholder when no pets exist
    - _Requirements: 2.1, 2.2, 2.3, 2.5_

- [ ]* 3.2 Write property test for pet list filtering
  - **Property 8: Status Filter Correctness**
  - **Validates: Requirements 2.5**

- [ ]* 3.3 Write property test for ownership filtering
  - **Property 5: Pet List Ownership Filter**
  - **Validates: Requirements 2.1, 6.4**

- [x] 3.4 Create PetCard component
  - Display pet photo or placeholder
  - Show name, species, breed, age, status badge
  - Add hover effect (scale-105)
  - Add click handler for navigation
  - _Requirements: 2.2, 9.5_

- [ ]* 3.5 Write property test for pet card information
  - **Property 6: Pet Card Information Completeness**
  - **Validates: Requirements 2.2**

- [ ]* 3.6 Write property test for placeholder display
  - **Property 39: Placeholder Image Display**
  - **Validates: Requirements 9.5**

- [ ] 4. Implement pet registration
  - [ ] 4.1 Create PetRegistrationForm component
    - Add form fields: name, species, breed, age, sex, color
    - Add optional fields: weight, special_marks, character, health_notes
    - Implement client-side validation (required fields)
    - Add photo upload with preview
    - Show loading state during submission
    - Display success notification with pending_verification message
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6_

- [ ]* 4.2 Write property test for required field validation
  - **Property 1: Required Field Validation**
  - **Validates: Requirements 1.2**

- [ ]* 4.3 Write property test for owner ID assignment
  - **Property 3: Owner ID Assignment**
  - **Validates: Requirements 1.4**

- [ ]* 4.4 Write property test for pending verification status
  - **Property 2: Pending Verification Status**
  - **Validates: Requirements 1.3**

- [ ] 4.5 Create backend POST /api/owner/pets endpoint
  - Validate required fields (name, species, age, sex, color)
  - Set owner_id from JWT token context
  - Set verification_status to "pending_verification"
  - Call PetID service to create pet
  - Return created pet with success message
  - _Requirements: 1.2, 1.3, 1.4_

- [ ]* 4.6 Write unit test for backend pet creation
  - Test required field validation
  - Test owner_id assignment from context
  - Test pending_verification status
  - _Requirements: 1.2, 1.3, 1.4_

- [ ] 5. Implement photo upload
  - [ ] 5.1 Add photo upload to registration form
    - File input with drag-and-drop support
    - Validate file format (jpg, png, webp)
    - Validate file size (max 5MB)
    - Show preview after selection
    - Display upload progress
    - _Requirements: 1.6, 9.1, 9.2, 9.3_

- [ ]* 5.2 Write property test for image format validation
  - **Property 35: Image Format Validation**
  - **Validates: Requirements 9.1**

- [ ]* 5.3 Write property test for image size validation
  - **Property 36: Image Size Validation**
  - **Validates: Requirements 9.2**

- [ ] 5.4 Create backend POST /api/owner/pets/:id/photo endpoint
  - Validate ownership (pet.owner_id == user_id)
  - Validate file format and size
  - Upload to PetID service
    - Return photo URL
    - Handle upload errors with retry option
    - _Requirements: 1.6, 8.4, 9.1, 9.2, 9.3, 9.4_

- [ ]* 5.5 Write property test for photo upload round trip
  - **Property 4: Photo Upload Round Trip**
  - **Validates: Requirements 1.6**

- [ ]* 5.6 Write property test for photo replacement
  - **Property 38: Photo Replacement**
  - **Validates: Requirements 9.4**

- [ ] 6. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. Implement pet editing
  - [ ] 7.1 Create PetEditForm component
    - Pre-fill form with existing pet data
    - Allow editing all fields except owner_id
    - Track changed fields only
    - Confirm before discarding changes
    - Show loading state during save
    - _Requirements: 2.4, 2.6_

- [ ]* 7.2 Write property test for pet edit round trip
  - **Property 7: Pet Edit Round Trip**
  - **Validates: Requirements 2.6**

- [ ] 7.3 Create backend PUT /api/owner/pets/:id endpoint
  - Validate ownership
  - Update only provided fields
  - Call PetID service to update pet
  - Return updated pet
  - _Requirements: 2.6, 8.4_

- [ ]* 7.4 Write property test for unauthorized pet access
  - **Property 33: Unauthorized Pet Access**
  - **Validates: Requirements 8.4**

- [ ] 8. Implement medical history display
  - [ ] 8.1 Create MedicalHistory component
    - Fetch medical events from PetID service
    - Fetch owner treatments from Owner backend
    - Separate clinic events from owner treatments
    - Display events chronologically
    - Show clinic events as read-only
    - Allow editing/deleting owner treatments
    - _Requirements: 3.1, 3.2, 3.6, 3.8_

- [ ]* 8.2 Write property test for medical history completeness
  - **Property 9: Medical History Completeness**
  - **Validates: Requirements 3.1**

- [ ]* 8.3 Write property test for medical event separation
  - **Property 10: Medical Event Separation**
  - **Validates: Requirements 3.2**

- [ ]* 8.4 Write property test for chronological ordering
  - **Property 14: Chronological Treatment Ordering**
  - **Validates: Requirements 3.6, 10.3**

- [ ]* 8.5 Write property test for clinic event read-only
  - **Property 16: Clinic Event Read-Only**
  - **Validates: Requirements 3.8**

- [ ] 9. Implement treatment tracking
  - [ ] 9.1 Create TreatmentForm component
    - Add fields: date, medication, dosage, next_date (optional), notes (optional)
    - Validate required fields (date, medication, dosage)
    - Validate date not in future
    - Validate next_date after date if provided
    - Show success notification after save
    - _Requirements: 3.3, 3.4, 3.5_

- [ ]* 9.2 Write property test for treatment required fields
  - **Property 11: Treatment Required Fields**
  - **Validates: Requirements 3.3**

- [ ]* 9.3 Write property test for optional next date
  - **Property 12: Treatment Optional Next Date**
  - **Validates: Requirements 3.4**

- [ ] 9.4 Create backend POST /api/owner/pets/:id/treatments endpoint
  - Validate ownership
  - Validate required fields
  - Save treatment to Owner backend database
    - Create event in PetID service with event_type "treatment"
    - Return created treatment
    - _Requirements: 3.3, 3.4, 3.5, 8.4_

- [ ]* 9.5 Write property test for treatment event type
  - **Property 13: Treatment Event Type**
  - **Validates: Requirements 3.5**

- [ ] 9.6 Create backend GET /api/owner/pets/:id/treatments endpoint
  - Validate ownership
  - Return all treatments for pet
  - Sort by date descending
  - _Requirements: 3.6_

- [ ] 9.7 Create backend PUT /api/owner/pets/:id/treatments/:treatmentId endpoint
  - Validate ownership
  - Update treatment
  - Return updated treatment
  - _Requirements: 3.3, 8.4_

- [ ] 9.8 Create backend DELETE /api/owner/pets/:id/treatments/:treatmentId endpoint
  - Validate ownership
  - Delete treatment
  - Return success response
  - _Requirements: 8.4_

- [ ] 10. Implement treatment reminders
  - [ ] 10.1 Create TreatmentReminders component
    - Fetch all user's pets
    - Calculate days until next treatment for each
    - Filter treatments within 7 days or overdue
    - Display reminder badge on Health tab
    - Highlight overdue treatments in red
    - Sort by date ascending
    - _Requirements: 3.7, 10.1, 10.2, 10.3, 10.4, 10.5_

- [ ]* 10.2 Write property test for days until calculation
  - **Property 40: Days Until Treatment Calculation**
  - **Validates: Requirements 10.1**

- [ ]* 10.3 Write property test for reminder threshold
  - **Property 15: Treatment Reminder Threshold**
  - **Validates: Requirements 3.7, 10.2**

- [ ]* 10.4 Write property test for overdue highlighting
  - **Property 41: Overdue Treatment Highlighting**
  - **Validates: Requirements 10.4**

- [ ] 10.5 Create backend GET /api/owner/reminders endpoint
  - Get all pets for authenticated user
  - Get all treatments with next_date not null
  - Calculate days_until for each
    - Filter treatments within 7 days or overdue
    - Return reminders sorted by date
    - _Requirements: 3.7, 10.1, 10.2, 10.3_

- [ ] 11. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 12. Implement status management
  - [ ] 12.1 Create StatusChangeModal component
    - Show different fields based on target status
    - For "lost": require date, location, circumstances, contact_info
    - Prevent changing to "looking_for_home" (show error)
    - Confirm before submitting
    - Show loading state during submission
    - Display confirmation notification on success
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

- [ ]* 12.2 Write property test for looking_for_home prevention
  - **Property 18: Looking For Home Prevention**
  - **Validates: Requirements 4.4**

- [ ]* 12.3 Write property test for status change confirmation
  - **Property 20: Status Change Confirmation**
  - **Validates: Requirements 4.5**

- [ ] 12.4 Create backend POST /api/owner/pets/:id/status endpoint
  - Validate ownership
  - Validate status (prevent "looking_for_home")
  - Update pet status in PetID service
  - Create status change event in PetID service
    - If status is "lost", create post in Main service
    - Return updated pet and optional post_id
    - _Requirements: 4.2, 4.3, 4.4, 4.6_

- [ ]* 12.5 Write property test for lost status post creation
  - **Property 17: Lost Status Post Creation**
  - **Validates: Requirements 4.3**

- [ ]* 12.6 Write property test for status change event creation
  - **Property 19: Status Change Event Creation**
  - **Validates: Requirements 4.2, 4.6**

- [ ] 13. Implement Main service integration
  - [ ] 13.1 Create Main API client (lib/main-api.ts)
    - Implement createPost method
    - Implement getPostsByPet method
    - Add JWT token to all requests
    - Handle error responses
    - _Requirements: 5.1, 5.2, 5.3_

- [ ]* 13.2 Write property test for pet ID parameter passing
  - **Property 21: Pet ID Parameter Passing**
  - **Validates: Requirements 5.2**

- [ ]* 13.3 Write property test for pet posts retrieval
  - **Property 22: Pet Posts Retrieval**
  - **Validates: Requirements 5.3**

- [ ] 13.4 Add "Create Post" button to pet profile
  - Navigate to Main service with pet_id parameter
  - Open in new tab or same tab based on user preference
    - _Requirements: 5.1, 5.2_

- [ ] 13.5 Add "View Posts" section to pet profile
  - Fetch posts from Main service
  - Display post cards with images and text
  - Link to full post in Main service
  - _Requirements: 5.3_

- [ ] 14. Implement veterinary visit history
  - [ ] 14.1 Add veterinary visits section to medical history
    - Filter medical events where clinic_id is not null
    - Display clinic name, date, event type, description
    - Sort by date descending
    - Link to clinic profile (if available)
    - _Requirements: 5.4, 5.5_

- [ ]* 14.2 Write property test for veterinary visit filtering
  - **Property 23: Veterinary Visit Filtering**
  - **Validates: Requirements 5.4**

- [ ]* 14.3 Write property test for clinic record information
  - **Property 24: Clinic Record Information**
  - **Validates: Requirements 5.5**

- [ ] 15. Implement Overview tab
  - [ ] 15.1 Create Overview page component
    - Display summary statistics (total pets, pets at home, lost pets, upcoming treatments)
    - Show recent activity (recent treatments, status changes)
    - Display quick actions (add pet, view reminders)
    - Use gradient cards for statistics
    - _Requirements: 6.3_

- [ ]* 15.2 Write property test for overview statistics
  - **Property 26: Overview Statistics Display**
  - **Validates: Requirements 6.3**

- [ ] 16. Implement Health tab
  - [ ] 16.1 Create Health page component
    - Display TreatmentReminders component
    - Show recent medical events across all pets
    - Add quick action to add treatment
    - Show badge count on tab if reminders exist
    - _Requirements: 6.5_

- [ ]* 16.2 Write property test for health tab content
  - **Property 27: Health Tab Content**
  - **Validates: Requirements 6.5**

- [ ] 17. Implement authentication and authorization
  - [ ] 17.1 Add authentication middleware to backend
    - Verify JWT token on all protected routes
    - Extract user_id from token and add to context
    - Return 401 for missing/invalid token
    - Return 403 for insufficient permissions
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ]* 17.2 Write property test for JWT token verification
  - **Property 31: JWT Token Verification**
  - **Validates: Requirements 8.2**

- [ ]* 17.3 Write property test for invalid token redirect
  - **Property 32: Invalid Token Redirect**
  - **Validates: Requirements 8.3**

- [ ] 17.4 Add authentication check to frontend layout
  - Check for valid session on mount
  - Redirect to SSO login if not authenticated
    - Show loading state during check
    - _Requirements: 6.2_

- [ ]* 17.5 Write property test for unauthenticated redirect
  - **Property 25: Unauthenticated Redirect**
  - **Validates: Requirements 6.2**

- [ ] 18. Implement error handling and validation
  - [ ] 18.1 Add form validation error display
    - Show inline errors below invalid fields
    - Highlight invalid fields with red border
    - Show error summary at top of form
    - Prevent submission until errors resolved
    - _Requirements: 7.1, 7.5_

- [ ]* 18.2 Write property test for validation error display
  - **Property 28: Validation Error Display**
  - **Validates: Requirements 7.1, 7.5**

- [ ] 18.3 Add API error handling
    - Display toast notifications for errors
    - Provide retry button for failed requests
    - Show specific error messages (401, 403, 404, 500)
    - Log errors to console
    - _Requirements: 7.2, 7.3_

- [ ]* 18.4 Write property test for API error handling
  - **Property 29: API Error Handling**
  - **Validates: Requirements 7.2**

- [ ]* 18.5 Write property test for photo upload error recovery
  - **Property 30: Photo Upload Error Recovery**
  - **Validates: Requirements 7.3**

- [ ] 19. Final integration and polish
  - [ ] 19.1 Test complete user flows
    - Register new pet → upload photo → view in list
    - Add treatment → view in medical history → see reminder
    - Change status to lost → verify post created in Main
    - Edit pet information → verify changes saved
    - _Requirements: All_

- [ ] 19.2 Add loading states and transitions
    - Skeleton loaders for pet list
    - Spinner for form submissions
    - Progress bar for photo uploads
    - Smooth transitions between tabs

- [ ] 19.3 Optimize performance
    - Implement pagination for pet list (if >20 pets)
    - Cache PetID service responses
    - Lazy load images
    - Debounce search/filter inputs

- [ ] 19.4 Update documentation
    - Add API documentation for all endpoints
    - Document component props and usage
    - Add README for Owner Cabinet
    - Update main project README

- [ ] 20. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional property-based tests and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
- The Owner Cabinet microservice infrastructure already exists (v0.11.0), so we're building on that foundation
- All pet data is stored in PetID service; Owner Cabinet only stores treatments locally
- Integration with Main service is required for "lost" posts and viewing pet posts
