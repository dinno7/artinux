# Welcome

Dear Candidate,

Thank you for your interest in joining our team and for taking the time to complete this technical challenge.

We understand that participating in a technical assessment requires a meaningful investment of your time and effort, and we genuinely appreciate your willingness to do so. This challenge has been designed to provide an opportunity for you to demonstrate your technical skills, engineering practices, software design decisions, testing approach, and problem-solving abilities in a realistic scenario.

Please read the requirements carefully before starting your implementation.

## Submission Deadline

You have **7 calendar days** from the date this challenge is assigned to complete and submit your solution.

**Submissions received after the 7-day deadline will not be considered unless otherwise agreed in advance.**

## Submission Format

You may submit your solution using either of the following methods:

- A link to a source code repository (GitHub, GitLab, Bitbucket, or equivalent)
- A compressed archive containing the complete project source code and all related files

Please ensure that all required deliverables are included and that the project can be built, tested, and executed using only the instructions provided in the README.

We wish you the best of luck and look forward to reviewing your solution.

Thank you again for your time and interest.

---

# Golang Technical Challenge

## Linux Artifact Storage Service

## Overview

This assignment evaluates your ability to design, implement, test, document, and containerize a production-grade application using Golang 1.23.7 and an S3-compatible object storage system.

The focus of this challenge is on:

- Software architecture
- Code quality
- Testing strategy
- Object storage integration
- Validation and reliability
- Maintainability
- Documentation
- Operational readiness

---

# Scenario

Your organization operates an internal CI/CD platform that produces Linux artifacts such as binaries, installation packages, archives, and log files.

A new service must be developed to manage these artifacts using an S3-compatible object storage platform powered by MinIO.

The service should allow authorized users and systems to store, retrieve, inspect, and manage Linux artifacts while ensuring validation, traceability, and operational reliability.

A local development and testing environment must also be provided using Docker.

---

# Technical Constraints

The implementation must use:

- Golang 1.23.7
- MinIO Go SDK (`github.com/minio/minio-go/v7`)
- Docker
- Docker Compose
- YAML-based configuration

---

# Infrastructure Requirements

A local MinIO environment must be provisioned using Docker Compose.

Minimum requirements:

- MinIO S3-compatible storage
- Dedicated bucket for artifact storage
- Configurable credentials
- Reproducible local setup

The solution should be executable by a reviewer with minimal setup effort.

---

# Functional Requirements

The solution must provide capabilities that allow:

- Verification of connectivity to the object storage backend
- Storage of Linux artifacts
- Retrieval and inspection of stored artifacts
- Management of previously stored artifacts
- Processing of multiple artifacts in a single operation
- Validation of incoming files before storage
- Storage and retrieval of artifact metadata
- Integrity verification of uploaded artifacts

The implementation details, API design, command structure, project layout, and user interaction model are intentionally left to the candidate.

---

# Artifact Requirements

Supported artifact types must be restricted to Linux-related distributable or operational files.

The implementation must validate:

- File extension
- File size
- File existence
- File accessibility

The allowed file types and size limits must be configurable through external configuration.

---

# Metadata Requirements

Every stored artifact must include metadata that enables auditing and traceability.

At minimum, metadata should include:

- Upload timestamp
- Operating system information
- Architecture information
- Host information
- User information
- SHA256 checksum

Additional metadata may be included if justified.

---

# Object Organization

Artifacts must be stored using a predictable and scalable naming strategy.

The chosen structure should:

- Prevent collisions
- Support long-term scalability
- Simplify artifact discovery
- Be clearly documented

The naming strategy and rationale should be explained in the README.

---

# Configuration Requirements

Application configuration must be externalized.

The implementation should avoid hardcoded values whenever possible.

Configuration should include, at minimum:

- Storage endpoint
- Authentication credentials
- Bucket configuration
- Upload restrictions
- Runtime settings

The application should be configurable without requiring source code changes.

---

# Error Handling Requirements

The system must provide clear and actionable error reporting.

Typical failure scenarios should be handled gracefully, including but not limited to:

- Authentication failures
- Connectivity issues
- Missing files
- Invalid files
- Storage failures
- Configuration errors
- Permission issues

The application must fail safely and provide meaningful diagnostics.

---

# Testing Requirements

Automated testing is mandatory.

The submission must include unit tests covering critical business logic.

At minimum, tests should cover:

- Validation logic
- Configuration loading
- Checksum generation
- Object naming/path generation
- Metadata generation
- Error handling paths

Requirements:

- Tests must be executable using standard Go tooling.
- Mocking or abstraction layers should be used where appropriate.
- Business logic should be testable independently from external services.
- Critical business logic should not depend directly on MinIO SDK implementations and should be testable through interfaces or abstractions.

Additional consideration will be given for:

- Table-driven tests
- Integration tests
- Storage client abstractions
- Meaningful test coverage
- Well-structured test suites

---

# MinIO Integration Requirements

The implementation must demonstrate proper usage of MinIO as an object storage backend.

The reviewer should be able to verify:

- Bucket management
- Object upload
- Object retrieval
- Object listing
- Object deletion
- Metadata handling

The solution should make appropriate use of the MinIO Go SDK.

---

# Containerization Requirements

The solution must be fully containerized.

The repository must contain:

- Dockerfile
- docker-compose.yml

A reviewer should be able to build and run the solution using standard Docker commands.

---

# Code Quality Expectations

The implementation should demonstrate:

- Clean architecture principles
- Separation of concerns
- Dependency management
- Maintainability
- Readability
- Proper error propagation
- Appropriate logging
- Sensible abstractions

There is no requirement to follow a specific architectural pattern, provided that design decisions are justified.

---

# README Requirements

A comprehensive README.md file is required and will be considered part of the evaluation.

The README should be written in English and must contain sufficient information for a reviewer to understand, build, test, and run the solution without additional assistance.

At minimum, the README should include:

## Project Overview

- Purpose of the application
- High-level architecture
- Design decisions

## Architecture

- Project structure
- Package organization
- Component responsibilities

## Prerequisites

- Required tools
- Required software versions

## Local Environment Setup

- Starting the MinIO environment
- Required configuration
- Initialization steps

## Build Instructions

- How to build the project

## Running the Application

- How to run the application
- Typical usage scenarios

## Configuration

- Available configuration options
- Example configuration file

## Testing

- How to execute unit tests
- How to execute integration tests (if provided)

## Assumptions

- Any assumptions made during implementation

## Limitations

- Known limitations
- Unsupported scenarios

## Future Improvements

- Potential production enhancements

### Important

The project should be reviewable without direct assistance from the candidate.

All setup, build, test, and execution steps must be documented in the README.

---

# AI-Assisted Development

The use of AI-assisted development tools (e.g. GitHub Copilot, ChatGPT, Claude, Cursor, or similar tools) is permitted.

However, candidates must be able to explain and justify all implementation decisions, architectural choices, and code submitted as part of the solution.

---

# Optional Enhancements

The following items are optional and may positively influence the evaluation:

- Concurrent artifact processing
- Retry mechanisms
- Progress reporting
- Presigned URL generation
- Structured logging
- Metrics collection
- Integration tests
- CI/CD pipeline configuration
- Observability improvements
- Additional security considerations

---

# Minimum Acceptance Criteria

To be considered for review, the submission must:

- Build successfully using Go 1.23.7
- Run successfully using the provided documentation
- Include automated tests
- Successfully connect to a MinIO instance
- Demonstrate artifact storage capabilities
- Include a README covering setup, usage, and design decisions
- Be submitted within the specified deadline

Submissions that do not meet these minimum requirements may be rejected without further review.

---

# Deliverables

The submission must include:

1. Source code
2. Dockerfile
3. docker-compose.yml
4. Example configuration file
5. Unit tests
6. README.md

Optional artifacts may also be included if relevant.

---

# Evaluation Criteria

| Category                 | Weight     |
| ------------------------ | ---------- |
| Architecture & Design    | 20%        |
| Testing Strategy         | 20%        |
| MinIO Integration        | 15%        |
| Code Quality             | 15%        |
| Validation & Reliability | 10%        |
| Error Handling           | 10%        |
| Documentation            | 5%         |
| Containerization         | 5%         |
| Bonus Features           | Up to +10% |

---

# Additional Notes

Candidates are encouraged to focus on:

- Correctness
- Maintainability
- Reliability
- Testing quality
- Documentation quality

Rather than implementing every possible enhancement.

Thoughtful trade-offs, documented assumptions, and clear design decisions are valued more highly than feature quantity.

---

# Closing Note

Thank you for completing this technical challenge.

We sincerely appreciate the time, effort, and thought you have invested in preparing your solution. We recognize that every submission represents a significant commitment, and we value the opportunity to learn more about your engineering approach through your work.

Regardless of the outcome, we would like to thank you for your interest in our team and for participating in our hiring process.

We hope you found the challenge engaging and enjoyable, and we genuinely hope that your solution demonstrates a strong fit for our engineering culture and technical expectations.

Most importantly, we hope this challenge is the beginning of a successful journey that leads to you joining our team.

We wish you continued success in your professional career and look forward to reviewing your submission.

Best regards,

Engineering Team
