# Specifications Index

> **Main Service Specifications**  
> **Last Updated:** 03.02.2026

## Overview

This directory contains all specifications for the Main Service. Each specification documents a specific feature, architecture, or system component.

---

## Available Specifications

### 1. Auth & CORS Architecture ✅

**Status:** Complete  
**Version:** 1.0.0  
**Path:** [auth-cors-architecture/](./auth-cors-architecture/)

**Description:**
Complete specification for authentication and CORS architecture in Main Service. Documents how authentication works in development vs production, why Gateway manages CORS, and how to use apiClient correctly.

**Quick Links:**
- [Overview](./auth-cors-architecture/README.md)
- [Requirements](./auth-cors-architecture/requirements.md)
- [Design](./auth-cors-architecture/design.md)
- [Troubleshooting](./auth-cors-architecture/troubleshooting.md)
- [Tasks](./auth-cors-architecture/tasks.md)
- [Changelog](./auth-cors-architecture/CHANGELOG.md)

**Key Topics:**
- Authentication flow (development vs production)
- CORS management (Gateway as single source)
- API Client implementation
- JWT token validation
- Cookie configuration
- Common issues and solutions

**When to Read:**
- Learning the authentication system
- Debugging 401 or CORS errors
- Implementing new authenticated features
- Deploying to production

---

## Specification Template

When creating a new specification, follow this structure:

```
.kiro/specs/
└── feature-name/
    ├── README.md           # Overview and quick reference
    ├── requirements.md     # User stories and acceptance criteria
    ├── design.md          # Detailed design and architecture
    ├── troubleshooting.md # Common issues and solutions
    ├── tasks.md           # Implementation tasks and checklist
    └── CHANGELOG.md       # Version history and changes
```

### Document Purposes

**README.md:**
- Overview of the specification
- Quick links to other documents
- Key concepts and principles
- Quick start guide

**requirements.md:**
- User stories (As a developer, I want...)
- Acceptance criteria (THE system SHALL...)
- Requirements for implementation

**design.md:**
- Architecture diagrams
- Detailed design decisions
- Component interactions
- Code examples
- Security considerations

**troubleshooting.md:**
- Common errors and symptoms
- Diagnostic steps
- Solutions with code examples
- Verification checklist
- Prevention best practices

**tasks.md:**
- Implementation tasks
- Testing checklist
- Maintenance tasks
- Future improvements
- Rollback plan

**CHANGELOG.md:**
- Version history
- What was fixed
- Breaking changes
- Migration guide

---

## Creating a New Specification

### Step 1: Create Directory

```bash
mkdir -p .kiro/specs/feature-name
```

### Step 2: Create Documents

```bash
cd .kiro/specs/feature-name

# Create all required documents
touch README.md
touch requirements.md
touch design.md
touch troubleshooting.md
touch tasks.md
touch CHANGELOG.md
```

### Step 3: Fill in Content

Use the Auth & CORS Architecture specification as a template:
- Copy structure from existing documents
- Adapt content to your feature
- Include diagrams where helpful
- Add code examples
- Document common issues

### Step 4: Update Index

Add your specification to this README.md:

```markdown
### 2. Your Feature Name

**Status:** In Progress / Complete  
**Version:** 1.0.0  
**Path:** [feature-name/](./feature-name/)

**Description:**
Brief description of what this specification covers.

**Quick Links:**
- [Overview](./feature-name/README.md)
- [Requirements](./feature-name/requirements.md)
- [Design](./feature-name/design.md)
- [Troubleshooting](./feature-name/troubleshooting.md)
- [Tasks](./feature-name/tasks.md)

**Key Topics:**
- Topic 1
- Topic 2
- Topic 3

**When to Read:**
- Scenario 1
- Scenario 2
- Scenario 3
```

---

## Specification Guidelines

### Writing Style

**Be Clear:**
- Use simple language
- Avoid jargon when possible
- Define technical terms
- Use examples

**Be Concise:**
- Get to the point quickly
- Use bullet points
- Break up long paragraphs
- Use headings

**Be Practical:**
- Include code examples
- Show before/after comparisons
- Provide diagnostic commands
- Give step-by-step solutions

### Diagrams

**Use ASCII diagrams for:**
- Architecture flows
- Request/response cycles
- Component interactions

**Example:**
```
Frontend → Gateway → Backend → Database
           ↑ Validates JWT
           ↑ Adds headers
```

### Code Examples

**Always show:**
- ✅ Correct way
- ❌ Wrong way
- Why it matters

**Example:**
```typescript
// ❌ Wrong
const response = await fetch('http://localhost:8000/api/posts');

// ✅ Correct
import { apiClient } from '@/lib/api';
const response = await apiClient.get('/api/posts');

// Why: apiClient handles URL configuration, credentials, and errors
```

---

## Maintenance

### Regular Updates

**When to Update:**
- After implementing new features
- After fixing bugs
- After architecture changes
- After deployment issues
- Quarterly reviews

**What to Update:**
- Add new issues to troubleshooting
- Update diagrams if architecture changed
- Add new best practices learned
- Update version numbers
- Update last updated dates

### Version Control

**Version Format:** MAJOR.MINOR.PATCH

- **MAJOR:** Breaking changes, major architecture changes
- **MINOR:** New features, new sections added
- **PATCH:** Bug fixes, clarifications, typos

**Example:**
- 1.0.0 → Initial release
- 1.1.0 → Added new troubleshooting section
- 1.1.1 → Fixed typos, clarified examples
- 2.0.0 → Changed authentication architecture

---

## Best Practices

### For Specification Authors

1. **Start with requirements** - Define what needs to be documented
2. **Include diagrams** - Visual aids help understanding
3. **Provide examples** - Show don't just tell
4. **Document common issues** - Save others time debugging
5. **Keep it updated** - Outdated docs are worse than no docs

### For Specification Readers

1. **Start with README** - Get overview first
2. **Read design for understanding** - Learn how it works
3. **Use troubleshooting when stuck** - Find solutions quickly
4. **Check tasks for implementation** - See what's done/todo
5. **Refer back often** - Specs are living documents

---

## Related Documentation

### Project Documentation

- [ARCHITECTURE.md](../../ARCHITECTURE.md) - Quick architecture reference
- [DEPLOYMENT.md](../../DEPLOYMENT.md) - Deployment guide
- [STARTUP.md](../../STARTUP.md) - Local development setup
- [README_API.md](../../README_API.md) - API documentation

### Workspace Rules

- [workspace-rules.md](../steering/workspace-rules.md) - Main workspace rules
- [auth-flow.md](../steering/auth-flow.md) - Authentication flow rules
- [docs-rules.md](../steering/docs-rules.md) - Documentation rules

---

## Statistics

### Current Specifications

- **Total:** 1
- **Complete:** 1
- **In Progress:** 0
- **Planned:** 0

### Documentation Coverage

- **Authentication:** ✅ Complete
- **CORS:** ✅ Complete
- **API Client:** ✅ Complete
- **Deployment:** ✅ Complete (DEPLOYMENT.md)
- **Gateway:** ✅ Complete (gateway.md)

### Planned Specifications

None currently planned. New specifications will be created as needed when:
- Implementing major new features
- Solving complex architectural problems
- Documenting critical systems
- After significant debugging sessions

---

## Contributing

### Adding a Specification

1. Create feature branch
2. Create specification directory
3. Write all required documents
4. Update this index
5. Submit for review
6. Merge after approval

### Updating a Specification

1. Create feature branch
2. Update relevant documents
3. Update version number
4. Update CHANGELOG.md
5. Update last updated date
6. Submit for review
7. Merge after approval

---

## Support

### Questions?

- Check the specification's README first
- Check the troubleshooting guide
- Check related documentation
- Ask the development team

### Found an Issue?

- Check if it's already documented
- If not, create a bug report
- Include diagnostic information
- Suggest a solution if possible

---

**Last Updated:** 03.02.2026  
**Maintained By:** Development Team
