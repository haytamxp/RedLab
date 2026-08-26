
# RedLab — Modular Active Directory Offensive Security Assessment Platform

> **Project:** RedLab
> **Purpose:** Authorized Active Directory security assessment
> **Project type:** PFA / Academic Cybersecurity Project
> **Current phase:** MVP development and validation
> **Target environment:** Active Directory laboratory environments and authorized security assessments

---

## 1. Project Overview

RedLab is a **modular offensive security assessment platform for Microsoft Active Directory**.

The platform is designed to centralize and structure an authorized security assessment workflow, from the creation of an assessment and execution of tasks to the collection of evidence, identification of findings, MITRE ATT&CK mapping, and report generation.

The project is not intended to replace specialized offensive-security tools. Instead, RedLab provides a **centralized orchestration and assessment layer** that connects the different components involved in an Active Directory assessment.

The current MVP focuses on demonstrating a complete and reproducible assessment lifecycle:
```text

Assessment
    |
    v
Task Creation
    |
    v
Windows Agent
    |
    v
Active Directory Enumeration
    |
    v
Evidence Collection
    |
    v
MITRE ATT&CK Mapping
    |
    v
Finding Creation
    |
    v
Assessment Results
    |
    v
HTML Report
````

The platform is being developed with a modular architecture so that additional assessment capabilities can be integrated without redesigning the entire system.

---

# 2. Project Objectives

The main objective of RedLab is to develop a platform capable of **centralizing, automating, and documenting authorized Active Directory security assessments**.

The project focuses on the following objectives:

* centralize the management of Active Directory assessments;
* provide authenticated access to the platform;
* manage assessment targets and tasks;
* register and authenticate Windows assessment agents;
* execute assessment tasks through the agent;
* collect technical evidence from the target environment;
* implement Active Directory enumeration capabilities;
* associate assessment activities with MITRE ATT&CK techniques;
* persist security findings;
* generate structured assessment reports;
* provide a web-based dashboard;
* maintain a modular architecture for future expansion.

---

# 3. Assessment Workflow

RedLab is organized around a structured assessment workflow.

```text
+------------------+
|     Operator     |
+--------+---------+
         |
         v
+------------------+
|    Assessment    |
+--------+---------+
         |
         v
+------------------+
|      Tasks       |
+--------+---------+
         |
         v
+------------------+
|  Windows Agent   |
+--------+---------+
         |
         v
+------------------+
| Active Directory |
|    Enumeration   |
+--------+---------+
         |
         v
+------------------+
|     Evidence     |
+--------+---------+
         |
         v
+------------------+
| MITRE ATT&CK Map |
+--------+---------+
         |
         v
+------------------+
|     Finding      |
+--------+---------+
         |
         v
+------------------+
|     Assessment   |
|     Results      |
+--------+---------+
         |
         v
+------------------+
|   HTML Report    |
+------------------+
```

The workflow separates the different responsibilities of the platform while allowing each stage to contribute information to the final assessment.

---

# 4. Existing MVP Scope

The current MVP contains the following implemented or targeted capabilities:

| Component         | Capability                 | Status      |
| ----------------- | -------------------------- | ----------- |
| Authentication    | JWT authentication         | Implemented |
| Authorization     | Role-Based Access Control  | Implemented |
| Persistence       | PostgreSQL                 | Implemented |
| Agent             | Windows agent registration | Implemented |
| Agent             | Agent token authentication | Implemented |
| Agent             | Heartbeat mechanism        | Implemented |
| Task System       | Task creation              | Implemented |
| Task System       | Task execution             | Implemented |
| Task System       | Task polling               | Implemented |
| Task System       | Result submission          | Implemented |
| Active Directory  | User enumeration           | Implemented |
| Active Directory  | Computer enumeration       | Implemented |
| Active Directory  | SPN enumeration            | Implemented |
| Security Analysis | MITRE ATT&CK mapping       | Implemented |
| Findings          | Finding persistence        | Implemented |
| Assessments       | Assessment management      | Implemented |
| Reporting         | HTML report generation     | Implemented |
| Frontend          | Web dashboard              | Implemented |

> **Note:** The implementation status above reflects the current MVP scope. Features not yet fully validated in the laboratory should be considered subject to testing and verification.

---

# 5. Architecture

RedLab follows a modular architecture composed of a frontend, a centralized backend, a task system, a Windows assessment agent, a PostgreSQL database, and an Active Directory environment.

```text
                         +----------------------+
                         |       Frontend       |
                         |    RedLab Dashboard  |
                         +----------+-----------+
                                    |
                                    | HTTP / REST
                                    v
                         +----------------------+
                         |    RedLab Backend    |
                         |       Go / Gin       |
                         +----------+-----------+
                                    |
                 +------------------+------------------+
                 |                  |                  |
                 v                  v                  v
        +----------------+   +-------------+   +-------------+
        |   PostgreSQL   |   | Task System |   |  Reporting  |
        +----------------+   +------+------+   +-------------+
                                    |
                                    | HTTP / REST
                                    v
                         +----------------------+
                         |    Windows Agent     |
                         |         Go           |
                         +----------+-----------+
                                    |
                                    | LDAP
                                    v
                         +----------------------+
                         |  Active Directory    |
                         +----------------------+
```

### Main Components

#### Frontend

The frontend provides the web interface used by the operator to interact with RedLab.

It is responsible for:

* authentication;
* assessment management;
* task management;
* result visualization;
* finding visualization;
* report access.

#### Backend

The backend is the central component of the platform.

It is responsible for:

* API management;
* authentication and authorization;
* assessment management;
* task management;
* agent management;
* result processing;
* finding persistence;
* database interaction;
* reporting coordination.

#### Task System

The task system provides the execution workflow between the backend and the assessment agent.

It allows the platform to:

* create tasks;
* assign tasks;
* track task status;
* retrieve task results;
* associate results with an assessment.

#### Windows Agent

The Windows agent is a Go-based component responsible for executing authorized assessment operations on Windows systems.

The agent communicates with the RedLab backend and interacts with the Active Directory environment according to the task assigned to it.

#### PostgreSQL

PostgreSQL is used for persistent storage.

The database stores information related to:

* users;
* roles;
* agents;
* assessments;
* tasks;
* modules;
* findings;
* evidence;
* MITRE ATT&CK mappings;
* reports.

---

# 6. Technology Stack

| Component            | Technology                        |
| -------------------- | --------------------------------- |
| Backend              | Go                                |
| HTTP Framework       | Gin                               |
| Database             | PostgreSQL                        |
| Directory            | Microsoft Active Directory / LDAP |
| Agent                | Go / Windows                      |
| Frontend             | JavaScript / Vite                 |
| Authentication       | JWT                               |
| Agent Authentication | Agent Registration Tokens         |
| Security Framework   | MITRE ATT&CK                      |
| Reporting            | HTML                              |
| API                  | REST / JSON                       |
| Version Control      | Git / GitHub                      |

---

# 7. Project Structure

```text
RedLab/
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   ├── migrations/
│   ├── .env
│   └── .env.example
│
├── agent/
│   ├── cmd/
│   │   └── agent/
│   ├── internal/
│   └── scripts/
│
├── frontend/
│   ├── src/
│   └── package.json
│
├── docs/
│   ├── ARCHITECTURE.md
│   ├── INSTALLATION.md
│   ├── API.md
│   ├── DEMO.md
│   └── MITRE-ATTACK.md
│
├── reports/
│
├── database/
│
├── lab/
│
├── scripts/
│
├── assets/
│
├── .gitignore
├── LICENSE
├── SECURITY.md
└── README.md
```

The repository is organized to keep the major components of the platform isolated while maintaining a single project structure.

---

# 8. Backend Architecture

The backend follows a layered structure.

```text
backend/
|
+-- cmd/
|   +-- api/
|       +-- main.go
|
+-- internal/
|   +-- app.go
|   +-- server.go
|
+-- config/
|   +-- config.go
|   +-- env.go
|   +-- loader.go
|
+-- handlers/
|   +-- health.go
|
+-- logger/
|   +-- zap.go
|
+-- middleware/
|
+-- models/
|
+-- repository/
|
+-- router/
|
+-- services/
```

The architecture separates:

* configuration;
* server initialization;
* HTTP handlers;
* middleware;
* business services;
* repositories;
* data models;
* logging.

This separation improves maintainability and provides a foundation for future extension.

---

# 9. Windows Agent

The Windows agent is responsible for executing assessment tasks received from the RedLab backend.

Its main lifecycle is:

```text
Agent Startup
      |
      v
Registration
      |
      v
Authentication
      |
      v
Heartbeat
      |
      v
Task Polling
      |
      v
Task Execution
      |
      v
Result Collection
      |
      v
Result Submission
```

The agent is designed to remain lightweight and focused on execution.

The backend remains responsible for assessment orchestration and persistence.

---

# 10. Supported Active Directory Assessment Modules

The current MVP supports the following assessment modules.

| Module                    | Purpose                           | MITRE ATT&CK |
| ------------------------- | --------------------------------- | ------------ |
| `AD_USER_ENUMERATION`     | Enumerate domain users            | T1087.002    |
| `AD_COMPUTER_ENUMERATION` | Enumerate domain computers        | T1018        |
| `SPN_ENUMERATION`         | Enumerate Service Principal Names | T1558.003    |

These modules demonstrate the modular nature of the platform.

Additional modules can be integrated using the same assessment/task architecture.

---

# 11. MITRE ATT&CK Integration

RedLab uses the **MITRE ATT&CK framework** to provide contextual mapping between assessment activities and known adversary techniques.

Current mappings include:

```text
AD_USER_ENUMERATION
        |
        v
T1087.002
Domain Account Discovery:
Domain Account

AD_COMPUTER_ENUMERATION
        |
        v
T1018
Remote System Discovery

SPN_ENUMERATION
        |
        v
T1558.003
Kerberoasting
```

The mapping is associated with the corresponding security finding.

This allows the final assessment to contain both:

* technical evidence;
* adversary-technique context.

Detailed mapping information is maintained in:

```text
docs/MITRE-ATTACK.md
```

---

# 12. Findings

Findings represent the security-relevant conclusions generated from assessment activities.

A finding can contain information such as:

* title;
* description;
* severity;
* affected target;
* evidence;
* assessment reference;
* module reference;
* MITRE ATT&CK technique;
* recommendation;
* creation timestamp.

The finding model allows raw technical results to be transformed into structured security observations.

---

# 13. Reporting

RedLab generates an HTML assessment report.

The report can contain:

* assessment information;
* target information;
* agent information;
* execution status;
* collected evidence;
* findings;
* severity information;
* MITRE ATT&CK mappings;
* recommendations.

The reporting pipeline can be represented as:

```text
Assessment
    |
    v
Tasks
    |
    v
Results
    |
    v
Evidence
    |
    v
Findings
    |
    v
MITRE ATT&CK
    |
    v
HTML Report
```

The objective is to transform low-level technical assessment data into structured information that can be reviewed by a security professional.

---

# 14. API

RedLab exposes a REST API used by the frontend and the Windows agent.

The API is responsible for resources such as:

```text
/health

/api/v1/
    |
    +-- auth/
    |
    +-- users/
    |
    +-- agents/
    |
    +-- assessments/
    |
    +-- tasks/
    |
    +-- modules/
    |
    +-- findings/
    |
    +-- reports/
```

The API is versioned to allow future evolution while preserving compatibility with existing clients.

Detailed API documentation is maintained in:

```text
docs/API.md
```

---

# 15. Quick Start

## 15.1 Backend

```powershell
cd backend
go mod tidy
go test ./...
go run .\cmd\api
```

Health check:

```powershell
curl.exe http://127.0.0.1:8080/health
```

Expected response:

```json
{
    "success": true,
    "message": "Backend is healthy",
    "data": {
        "service": "RedLab Backend",
        "status": "running",
        "version": "1.0.0"
    }
}
```

---

## 15.2 Windows Agent

```powershell
cd agent
go mod tidy
go test ./...
go build -o redlab-agent.exe .\cmd\agent
.\redlab-agent.exe
```

The agent must be configured with the appropriate backend endpoint and authorized registration credentials.

---

## 15.3 Frontend

```powershell
cd frontend
npm install
npm run dev
```

Production build:

```powershell
npm run build
```

---

# 16. Configuration

Sensitive configuration must remain outside the source code.

The project uses environment-based configuration.

Example:

```text
backend/.env.example
```

The local environment file should contain values specific to the deployment environment.

Do not commit the real `.env` file.

---

# 17. Security Requirements

RedLab is a cybersecurity platform and may handle sensitive assessment information.

The following information must **never** be committed to the repository:

* database passwords;
* LDAP passwords;
* JWT secrets;
* agent registration tokens;
* API credentials;
* private keys;
* certificates containing private material;
* production assessment data;
* sensitive Active Directory exports;
* credentials obtained during an authorized assessment.

Use:

```text
.env.example
```

as the configuration template.

Before pushing changes, verify that secrets are not present in:

* source files;
* Git history;
* logs;
* test fixtures;
* configuration examples;
* generated reports.

---

# 18. Development Principles

## Modular by design

Assessment capabilities should be implemented as independent modules whenever possible.

## Separate orchestration from execution

The backend coordinates assessments while the agent performs the operations assigned to it.

## Security by design

Authentication, authorization, input validation, secret management, logging and error handling are treated as core requirements.

## Evidence-driven assessment

Security findings should be supported by collected evidence rather than generated only from assumptions.

## Reproducibility

Laboratory configuration, API behavior, database structures and assessment procedures should be documented so that results can be reproduced.

## Preserve architectural boundaries

Changes to one component should not unnecessarily introduce dependencies into unrelated components.

## Validate before claiming support

A feature should be considered fully implemented only after functional validation in the intended laboratory environment.

---

# 19. Project Timeline

| Milestone | Objective                                                            |
| --------- | -------------------------------------------------------------------- |
| M0        | Project analysis, architecture definition and repository preparation |
| M1        | Backend foundation and API                                           |
| M2        | Database and authentication                                          |
| M3        | Agent registration and task system                                   |
| M4        | Active Directory assessment modules                                  |
| M5        | Findings and MITRE ATT&CK integration                                |
| M6        | Reporting and dashboard                                              |
| M7        | Laboratory validation and final documentation                        |

The timeline may evolve according to implementation progress and validation results.

---

# 20. Validation Strategy

RedLab is validated in a controlled Active Directory laboratory.

Validation focuses on:

* backend availability;
* API functionality;
* authentication;
* authorization;
* database persistence;
* agent registration;
* agent authentication;
* heartbeat;
* task creation;
* task execution;
* task polling;
* result submission;
* Active Directory enumeration;
* evidence persistence;
* MITRE ATT&CK mapping;
* finding generation;
* HTML report generation;
* frontend integration.

The objective is to verify the complete end-to-end workflow rather than validating isolated components only.

---

# 21. Documentation

Technical documentation is maintained under:

```text
docs/
```

Important documentation includes:

| Document          | Description                                    |
| ----------------- | ---------------------------------------------- |
| `ARCHITECTURE.md` | System architecture and component interactions |
| `INSTALLATION.md` | Installation and environment setup             |
| `API.md`          | REST API documentation                         |
| `DEMO.md`         | Demonstration and validation procedure         |
| `MITRE-ATTACK.md` | ATT&CK technique mapping                       |

Additional documentation may be added as the platform evolves.

---

# 22. Current Project Status

## MVP Status

**Status: Active development and validation**

The current RedLab implementation already demonstrates the core assessment lifecycle:

```text
Task
  |
  v
Agent
  |
  v
Active Directory Enumeration
  |
  v
Evidence
  |
  v
MITRE ATT&CK
  |
  v
Finding
  |
  v
Assessment
  |
  v
Report
```

The project currently focuses on delivering a **functional MVP** rather than a complete offensive-security automation framework.

The architecture is intentionally designed to provide a foundation for future modules and advanced assessment capabilities.

---

# 23. Future Improvements

Potential future capabilities include:

* additional Active Directory enumeration modules;
* privilege and permission analysis;
* attack-path analysis;
* additional credential and authentication assessments;
* richer evidence correlation;
* expanded MITRE ATT&CK coverage;
* advanced task scheduling;
* parallel task execution;
* improved report generation;
* advanced dashboard visualization;
* distributed agents;
* additional target environments;
* stronger secret-management integration.

These capabilities are considered future extensions and are not necessarily part of the current MVP.

---

# 24. Security Disclaimer

RedLab is intended **only for authorized security assessments, academic projects, controlled cybersecurity laboratories, and legitimate security research**.

The platform must only be used against systems for which the operator has explicit authorization.

The operator is responsible for ensuring that:

* the assessment scope is explicitly defined;
* the target infrastructure is authorized;
* offensive operations remain within the approved scope;
* sensitive information is protected;
* collected evidence is handled securely;
* generated reports are stored appropriately.

Unauthorized use of the platform against third-party infrastructure is prohibited.

---

# 25. Project License

See:

```text
LICENSE
```

for the applicable project license and usage conditions.

---

# 26. Maintainers

**Developer / Project Owner**

Haytam Ragueb

**Project**

RedLab

**Repository**

GitHub — `haytamxp/RedLab`

```
```
