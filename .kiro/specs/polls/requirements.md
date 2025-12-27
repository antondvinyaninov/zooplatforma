# Requirements Document: Система голосования (Polls)

## Introduction

Система голосования позволяет пользователям создавать опросы с несколькими вариантами ответов, голосовать и просматривать результаты. Голосования интегрируются в систему меток (posts) и отображаются в ленте.

## Glossary

- **Poll**: Опрос с вопросом и вариантами ответов
- **Poll_Option**: Вариант ответа в опросе
- **Vote**: Голос пользователя за определенный вариант
- **Post**: Метка (публикация), к которой может быть прикреплен опрос
- **User**: Пользователь платформы
- **Multiple_Choice**: Режим опроса, позволяющий выбрать несколько вариантов
- **Single_Choice**: Режим опроса, позволяющий выбрать только один вариант

## Requirements

### Requirement 1: Создание опроса

**User Story:** As a user, I want to create polls with multiple answer options, so that I can gather opinions from the community.

#### Acceptance Criteria

1. WHEN a user clicks the poll icon in the create post modal, THE System SHALL display a poll creation form
2. WHEN creating a poll, THE System SHALL require a question text with minimum 1 character and maximum 500 characters
3. WHEN creating a poll, THE System SHALL allow adding between 2 and 10 answer options
4. WHEN creating a poll, THE System SHALL require each option text to be between 1 and 200 characters
5. WHEN creating a poll, THE System SHALL allow toggling multiple choice mode (on/off)
6. WHEN creating a poll, THE System SHALL allow toggling vote changes permission (on/off)
7. WHEN creating a poll, THE System SHALL allow setting an expiration date and time (optional)
8. WHEN a poll is created without expiration, THE System SHALL set expiration to null
9. WHEN a poll is created with a post, THE System SHALL save the poll data linked to the post ID

### Requirement 2: Отображение опроса до голосования

**User Story:** As a user, I want to see poll options before voting, so that I can make an informed choice.

#### Acceptance Criteria

1. WHEN a user views a post with a poll and has not voted, THE System SHALL display the poll question
2. WHEN displaying a poll before voting, THE System SHALL show all answer options as clickable buttons
3. WHEN displaying a poll before voting, THE System SHALL NOT show vote percentages
4. WHEN displaying a poll before voting, THE System SHALL NOT show vote counts
5. WHEN displaying a poll before voting, THE System SHALL show total number of voters as 0 or actual count
6. WHEN displaying a multiple choice poll, THE System SHALL show checkboxes for options
7. WHEN displaying a single choice poll, THE System SHALL show radio buttons for options
8. WHEN displaying an expired poll, THE System SHALL show results to all users regardless of voting status

### Requirement 3: Голосование в опросе

**User Story:** As a user, I want to vote in polls, so that I can express my opinion.

#### Acceptance Criteria

1. WHEN a user selects an option in a single choice poll, THE System SHALL record the vote for that option
2. WHEN a user selects multiple options in a multiple choice poll, THE System SHALL record votes for all selected options
3. WHEN a user submits a vote, THE System SHALL increment the vote count for selected options
4. WHEN a user submits a vote, THE System SHALL increment the total voters count by 1
5. WHEN a user has already voted and vote changes are disabled, THE System SHALL prevent voting again
6. WHEN a user has already voted and vote changes are enabled, THE System SHALL allow changing the vote
7. WHEN a user changes their vote, THE System SHALL decrement previous option votes and increment new option votes
8. WHEN a user votes in an expired poll, THE System SHALL reject the vote with an error message
9. WHEN an unauthenticated user attempts to vote, THE System SHALL require authentication

### Requirement 4: Отображение результатов опроса

**User Story:** As a user, I want to see poll results after voting, so that I can understand community opinion.

#### Acceptance Criteria

1. WHEN a user has voted in a poll, THE System SHALL display the poll results
2. WHEN displaying poll results, THE System SHALL show vote percentages for each option
3. WHEN displaying poll results, THE System SHALL show vote counts for each option
4. WHEN displaying poll results, THE System SHALL show total number of voters
5. WHEN displaying poll results, THE System SHALL highlight the option(s) the user voted for
6. WHEN calculating percentages, THE System SHALL round to 1 decimal place
7. WHEN a poll has 0 votes, THE System SHALL display 0% for all options
8. WHEN displaying results, THE System SHALL order options by their original creation order

### Requirement 5: Управление опросом

**User Story:** As a poll creator, I want to manage my polls, so that I can control the voting process.

#### Acceptance Criteria

1. WHEN a poll reaches its expiration date and time, THE System SHALL automatically close the poll
2. WHEN a poll is closed, THE System SHALL prevent new votes
3. WHEN a poll is closed, THE System SHALL display results to all users
4. WHEN viewing a closed poll, THE System SHALL show an "Опрос завершен" indicator
5. WHEN a poll creator deletes a post, THE System SHALL also delete the associated poll and votes

### Requirement 6: Интеграция с постами

**User Story:** As a user, I want polls to be seamlessly integrated with posts, so that I can share polls in my feed.

#### Acceptance Criteria

1. WHEN a post contains a poll, THE System SHALL display the poll within the post card
2. WHEN a post with a poll is displayed in the feed, THE System SHALL show the poll question and options
3. WHEN a post with a poll is opened in modal view, THE System SHALL display the full poll with voting interface
4. WHEN a user creates a post with a poll, THE System SHALL save both post and poll data atomically
5. WHEN a post with a poll is deleted, THE System SHALL cascade delete the poll and all votes

### Requirement 7: Валидация и ограничения

**User Story:** As a system, I want to validate poll data, so that data integrity is maintained.

#### Acceptance Criteria

1. WHEN creating a poll, THE System SHALL validate that question text is not empty
2. WHEN creating a poll, THE System SHALL validate that there are at least 2 options
3. WHEN creating a poll, THE System SHALL validate that there are no more than 10 options
4. WHEN creating a poll, THE System SHALL validate that all option texts are not empty
5. WHEN creating a poll with expiration, THE System SHALL validate that expiration is in the future
6. WHEN a user votes, THE System SHALL validate that selected options exist in the poll
7. WHEN a user votes in a single choice poll, THE System SHALL validate that only one option is selected

### Requirement 8: Производительность и масштабируемость

**User Story:** As a system, I want polls to perform efficiently, so that user experience is smooth.

#### Acceptance Criteria

1. WHEN loading a post with a poll, THE System SHALL fetch poll data in a single database query
2. WHEN displaying poll results, THE System SHALL calculate percentages on the backend
3. WHEN a user votes, THE System SHALL update vote counts using atomic database operations
4. WHEN loading multiple posts with polls, THE System SHALL batch load poll data efficiently
