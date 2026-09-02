# RedLab

## Modular Active Directory Offensive Security Assessment Platform

RedLab is a modular platform designed to **centralize, automate, and document authorized Microsoft Active Directory security assessments**.

It provides a web-based operator interface, a centralized Go backend, a Windows assessment agent, LDAP-based Active Directory enumeration modules, persistent assessment data, MITRE ATT&CK mapping, findings management, and report generation.

> **Important:** RedLab is intended only for authorized security assessments, internal security testing, educational laboratories, and controlled Active Directory environments.

---

## 1. What is RedLab?

Traditional Active Directory assessments often require multiple tools, manual commands, disconnected results, and manually prepared reports.

RedLab introduces a centralized workflow:

```text
                   Operator
                      │
                      ▼
              ┌───────────────┐
              │ RedLab Web UI │
              └───────┬───────┘
                      │ REST / HTTP
                      ▼
              ┌───────────────┐
              │ Go Backend    │
              │ Gin API       │
              └───────┬───────┘
                      │
              ┌───────┼────────┐
              │       │        │
              ▼       ▼        ▼
        PostgreSQL   Tasks   Reporting
                      │
                      ▼
              ┌───────────────┐
              │ Windows Agent │
              └───────┬───────┘
                      │ LDAP
                      ▼
              ┌───────────────┐
              │ Active        │
              │ Directory     │
              └───────────────┘
```

The assessment lifecycle is:

```text
Assessment
    ↓
Task Creation
    ↓
Windows Agent
    ↓
Module Execution
    ↓
Active Directory Enumeration
    ↓
Result Collection
    ↓
MITRE ATT&CK Mapping
    ↓
Finding
    ↓
Assessment Result
    ↓
HTML Report
```

---

# 2. Main Features

### Authentication and authorization

* JWT-based operator authentication
* Role/permission-based authorization
* Protected API endpoints
* Separate authentication mechanism for Windows agents

### Agent management

* Agent registration
* Agent identification using UUID
* Secure registration token generation
* Token stored as a hash in PostgreSQL
* Bearer-token authentication
* Agent heartbeat
* Agent online/offline status
* Task polling

### Active Directory assessment

Current modules include:

| Module                    | Purpose                           |
| ------------------------- | --------------------------------- |
| `DOMAIN_INFO`             | Collect domain information        |
| `AD_USER_ENUMERATION`     | Enumerate domain users            |
| `AD_GROUP_ENUMERATION`    | Enumerate domain groups           |
| `AD_COMPUTER_ENUMERATION` | Enumerate domain computers        |
| `SPN_ENUMERATION`         | Enumerate Service Principal Names |

The modules use LDAP to communicate with Active Directory.

### Security analysis

RedLab associates assessment modules with MITRE ATT&CK techniques.

Examples:

| Module                    | MITRE ATT&CK                         |
| ------------------------- | ------------------------------------ |
| `AD_USER_ENUMERATION`     | T1087.002 — Domain Account Discovery |
| `AD_COMPUTER_ENUMERATION` | T1018 — Remote System Discovery      |
| `SPN_ENUMERATION`         | T1558.003 — Kerberoasting            |

### Findings and reporting

RedLab can store:

* assessment results
* task results
* findings
* severity
* evidence
* affected agents
* MITRE ATT&CK techniques
* recommendations

The platform can generate an HTML assessment report.

---

# 3. Technology Stack

| Component            | Technology                 |
| -------------------- | -------------------------- |
| Backend              | Go                         |
| HTTP Framework       | Gin                        |
| Database             | PostgreSQL                 |
| Directory Services   | Microsoft Active Directory |
| Directory Protocol   | LDAP                       |
| Agent                | Go / Windows               |
| Frontend             | JavaScript / Vite          |
| Authentication       | JWT                        |
| Agent Authentication | Bearer registration token  |
| Password Hashing     | bcrypt                     |
| Security Framework   | MITRE ATT&CK               |
| API                  | REST / JSON                |
| Reporting            | HTML                       |
| Version Control      | Git / GitHub               |

---

# 4. Repository Structure

```text
RedLab/
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── app/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── database/
│   │   ├── dto/
│   │   ├── handlers/
│   │   ├── ldap/
│   │   ├── logger/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── permissions/
│   │   ├── reporting/
│   │   ├── repository/
│   │   ├── response/
│   │   └── services/
│   │
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
│
├── agent/
│   ├── cmd/
│   │   └── agent/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── api/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── heartbeat/
│   │   ├── mitre/
│   │   ├── modules/
│   │   └── tasks/
│   │
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── pages/
│   │   ├── services/
│   │   └── types/
│   │
│   ├── package.json
│   └── vite.config.js
│
├── docs/
│   ├── API.md
│   ├── DEMO.md
│   ├── INSTALLATION.md
│   └── MITRE-ATTACK.md
│
├── reports/
├── lab/
├── scripts/
├── assets/
│
├── README.md
├── SECURITY.md
└── LICENSE
```

---

# 5. Requirements

## Required software

Install the following before starting:

* Git
* Go 1.26+
* PostgreSQL
* Node.js and npm
* Windows 10/11 or Windows Server for the agent
* A Windows Server Active Directory laboratory for LDAP-based modules

Recommended laboratory architecture:

```text
┌──────────────────────────────┐
│ Windows Server 2022          │
│                              │
│ Active Directory Domain      │
│ Controller + DNS + LDAP      │
│ Example IP: 192.168.128.10   │
└──────────────┬───────────────┘
               │ LDAP :389
               │
┌──────────────▼───────────────┐
│ Windows Assessment Machine   │
│                              │
│ RedLab Agent                 │
└──────────────┬───────────────┘
               │ HTTP :8080
               │
┌──────────────▼───────────────┐
│ RedLab Backend               │
│ Go + Gin                     │
└──────────────┬───────────────┘
               │ PostgreSQL :5432
               │
┌──────────────▼───────────────┐
│ PostgreSQL                   │
│ Database: redlab             │
└──────────────────────────────┘
```

---

# 6. Clone the Repository

```powershell
git clone https://github.com/haytamxp/RedLab.git
cd RedLab
```

---

# 7. PostgreSQL Setup

Create the database:

```powershell
psql -h 127.0.0.1 -p 5432 -U postgres -c "CREATE DATABASE redlab;"
```

If PostgreSQL reports:

```text
ERROR: database "redlab" already exists
```

the database is already initialized and you can continue.

Verify the connection:

```powershell
psql -h 127.0.0.1 -p 5432 -U postgres -d redlab
```

Then:

```sql
\q
```

---

# 8. Backend Configuration

Go to the backend directory:

```powershell
cd D:\PFA-RedLab\Project\RedLab\backend
```

Create your local environment file:

```powershell
Copy-Item .\.env.example .\.env
```

Configure the following values inside `.env`:

```text
APP_NAME
APP_VERSION
APP_ENV

SERVER_HOST
SERVER_PORT

SERVER_READ_TIMEOUT
SERVER_WRITE_TIMEOUT
SERVER_IDLE_TIMEOUT

DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
DB_SSLMODE

LDAP_HOST
LDAP_PORT
LDAP_USERNAME
LDAP_PASSWORD
LDAP_BASE_DN
LDAP_USE_TLS

JWT_SECRET
JWT_EXPIRATION

LOG_LEVEL
LOG_FORMAT
LOG_OUTPUT
LOG_DEVELOPMENT
```

Example database configuration:

```text
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=YOUR_POSTGRES_PASSWORD
DB_NAME=redlab
DB_SSLMODE=disable
```

Example LDAP configuration:

```text
LDAP_HOST=192.168.128.10
LDAP_PORT=389
LDAP_USERNAME=YOUR_DOMAIN_ACCOUNT
LDAP_PASSWORD=YOUR_PASSWORD
LDAP_BASE_DN=DC=example,DC=local
LDAP_USE_TLS=false
```

Example JWT configuration:

```text
JWT_SECRET=REPLACE_WITH_A_LONG_RANDOM_SECRET
JWT_EXPIRATION=3600
```

> Never commit the real `.env` file.

---

# 9. Start the Backend

From:

```powershell
cd D:\PFA-RedLab\Project\RedLab\backend
```

Install dependencies:

```powershell
go mod tidy
```

Run tests:

```powershell
go test ./...
```

Start the backend:

```powershell
go run .\cmd\api
```

The backend normally listens on:

```text
http://127.0.0.1:8080
```

---

# 10. Test the Backend

Open another PowerShell window:

```powershell
curl.exe http://127.0.0.1:8080/health
```

Expected response:

```json
{
  "success": true,
  "message": "Backend is healthy"
}
```

The API is versioned under:

```text
/api/v1
```

Health check:

```text
GET /api/v1/health
```

---

# 11. Active Directory Connectivity

Before starting the agent, verify LDAP connectivity from the Windows machine where the agent will run.

Example:

```powershell
Test-NetConnection 192.168.128.10 -Port 389
```

Expected:

```text
TcpTestSucceeded : True
```

The RedLab agent must be able to reach the configured Active Directory LDAP endpoint.

---

# 12. Agent Registration

The Windows agent requires:

```text
REDLAB_SERVER_URL
REDLAB_AGENT_ID
REDLAB_AGENT_TOKEN
```

The recommended workflow is:

```text
Operator logs in
        ↓
Create agent through RedLab
        ↓
Backend generates agent ID
        ↓
Backend generates registration token
        ↓
Token is returned once
        ↓
Store token securely
        ↓
Configure Windows Agent
```

The registration token is not intended to be stored in plaintext in the RedLab database. The backend stores a hash of the token.

---

# 13. Build the Windows Agent

Open another PowerShell window:

```powershell
cd D:\PFA-RedLab\Project\RedLab\agent
```

Install dependencies:

```powershell
go mod tidy
```

Run tests:

```powershell
go test ./...
```

Build the agent:

```powershell
go build -o redlab-agent.exe .\cmd\agent
```

---

# 14. Configure the Windows Agent

Set the environment variables:

```powershell
$env:REDLAB_SERVER_URL = "http://127.0.0.1:8080"
$env:REDLAB_AGENT_ID = "YOUR_AGENT_ID"
$env:REDLAB_AGENT_TOKEN = "YOUR_AGENT_TOKEN"

$env:REDLAB_HEARTBEAT_SECONDS = "30"
$env:REDLAB_POLL_SECONDS = "10"

$env:REDLAB_LDAP_URL = "ldap://192.168.128.10:389"
$env:REDLAB_LDAP_USERNAME = "YOUR_DOMAIN_USERNAME"
$env:REDLAB_LDAP_PASSWORD = "YOUR_DOMAIN_PASSWORD"
$env:REDLAB_LDAP_BASE_DN = "DC=example,DC=local"
```

Then start:

```powershell
.\redlab-agent.exe
```

The agent starts:

```text
Heartbeat Service
        +
Task Poller
```

The heartbeat informs the backend that the agent is alive.

The task poller periodically asks the backend for work.

---

# 15. Agent Communication Flow

The agent communication model is:

```text
Windows Agent
      │
      │ Authorization: Bearer <agent-token>
      ▼
RedLab Backend
      │
      ├── Heartbeat
      │
      ├── Request Next Task
      │
      └── Submit Task Result
```

The current API exposes endpoints such as:

```text
POST /api/v1/agents/:id/heartbeat
POST /api/v1/agents/:id/tasks/next
GET  /api/v1/agents/:id/tasks
POST /api/v1/agents/:id/tasks/:taskId/result
```

Agent authentication is performed independently from human JWT authentication.

---

# 16. Frontend

Open another PowerShell window:

```powershell
cd D:\PFA-RedLab\Project\RedLab\frontend
```

Install dependencies:

```powershell
npm install
```

Start the development server:

```powershell
npm run dev
```

Vite will display the local frontend URL.

Typical example:

```text
http://localhost:5173
```

Production build:

```powershell
npm run build
```

---

# 17. Recommended Startup Order

For a complete laboratory deployment, start the components in this order:

```text
1. PostgreSQL
        ↓
2. Active Directory / Domain Controller
        ↓
3. RedLab Backend
        ↓
4. Windows RedLab Agent
        ↓
5. RedLab Frontend
```

---

# 18. Quick Start — Three Terminals

After the initial configuration, normal development can be done with three terminals.

### Terminal 1 — Backend

```powershell
cd D:\PFA-RedLab\Project\RedLab\backend
go run .\cmd\api
```

### Terminal 2 — Agent

```powershell
cd D:\PFA-RedLab\Project\RedLab\agent
.\redlab-agent.exe
```

### Terminal 3 — Frontend

```powershell
cd D:\PFA-RedLab\Project\RedLab\frontend
npm run dev
```

---

# 19. Complete First Test

Once the backend, frontend, agent, database, and Active Directory are running:

```text
                  ┌────────────┐
                  │   Login    │
                  └─────┬──────┘
                        │
                        ▼
               Create Assessment
                        │
                        ▼
                 Select Agent
                        │
                        ▼
                  Create Task
                        │
                        ▼
              Agent Polls Backend
                        │
                        ▼
              Module Execution
                        │
                        ▼
                  LDAP Query
                        │
                        ▼
              Active Directory
                        │
                        ▼
                   Result
                        │
                        ▼
                Finding / Evidence
                        │
                        ▼
                MITRE ATT&CK
                        │
                        ▼
                 HTML Report
```

Example task:

```text
AD_USER_ENUMERATION
```

The agent resolves the task through the module registry and executes the user enumeration module against LDAP.

---

# 20. Example: AD User Enumeration

The `AD_USER_ENUMERATION` module performs an LDAP search for domain user objects.

Conceptually:

```text
Task:
AD_USER_ENUMERATION
        ↓
Module Registry
        ↓
UsersModule
        ↓
LDAP
        ↓
Active Directory
        ↓
Users
        ↓
Structured JSON result
```

The module collects attributes such as:

```text
sAMAccountName
userPrincipalName
displayName
userAccountControl
```

The result also contains its associated MITRE ATT&CK technique:

```text
T1087.002
Domain Account Discovery
```

---

# 21. Modular Architecture

RedLab uses a common module interface.

Conceptually:

```go
type Module interface {
    Name() string
    Technique() mitre.Technique
    Execute(
        ctx context.Context,
        payload json.RawMessage,
    ) (any, error)
}
```

This allows new assessment capabilities to be added without redesigning the entire platform.

Current registry:

```text
DOMAIN_INFO
AD_USER_ENUMERATION
AD_GROUP_ENUMERATION
AD_COMPUTER_ENUMERATION
SPN_ENUMERATION
```

Future modules can follow the same architecture.

Possible future additions:

```text
AD_GPO_ENUMERATION
AD_OU_ENUMERATION
AD_TRUST_ENUMERATION
AD_DOMAIN_POLICY_ENUMERATION
AD_ACL_ENUMERATION
AD_DELEGATION_ENUMERATION
```

---

# 22. MITRE ATT&CK Integration

RedLab associates technical assessment operations with adversary techniques.

Example:

```text
AD_USER_ENUMERATION
        ↓
T1087.002
Domain Account Discovery
```

```text
AD_COMPUTER_ENUMERATION
        ↓
T1018
Remote System Discovery
```

```text
SPN_ENUMERATION
        ↓
T1558.003
Kerberoasting
```

The goal is to provide security context rather than simply presenting raw LDAP output.

---

# 23. Findings

A finding may contain:

```text
Title
Description
Severity
Affected Target
Evidence
Assessment ID
Task ID
Agent ID
MITRE Technique
Recommendation
Timestamp
```

This transforms technical enumeration output into structured security observations.

---

# 24. Reporting

The reporting workflow is:

```text
Assessment
    ↓
Tasks
    ↓
Results
    ↓
Evidence
    ↓
Findings
    ↓
MITRE ATT&CK
    ↓
HTML Report
```

Assessment reports can contain:

* assessment information
* target information
* agent information
* execution results
* collected evidence
* findings
* severity
* MITRE ATT&CK mappings
* recommendations

---

# 25. REST API

The API is organized under:

```text
/api/v1
```

Main resources include:

```text
/auth
/agents
/tasks
/assessments
/findings
```

Examples:

```text
GET  /health
GET  /api/v1/health

POST /api/v1/auth/register
POST /api/v1/auth/login

POST /api/v1/agents
GET  /api/v1/agents
GET  /api/v1/agents/:id

POST /api/v1/tasks

POST /api/v1/assessments
GET  /api/v1/assessments
GET  /api/v1/assessments/:id
GET  /api/v1/assessments/:id/report

GET /api/v1/findings
GET /api/v1/findings/:id
```

Agent-specific endpoints:

```text
POST /api/v1/agents/:id/heartbeat
GET  /api/v1/agents/:id/tasks
POST /api/v1/agents/:id/tasks/next
POST /api/v1/agents/:id/tasks/:taskId/result
```

---

# 26. Security Architecture

RedLab handles sensitive assessment information, so security is part of the platform design.

### Operator authentication

Human users authenticate through JWT.

```text
Username / Password
        ↓
Authentication
        ↓
JWT
        ↓
Protected API
```

### Agent authentication

Agents authenticate independently:

```text
Agent Token
    ↓
SHA-256 hash
    ↓
Compare with stored token hash
    ↓
Authenticate agent
```

### Secret management

Never commit:

```text
.env
database passwords
LDAP passwords
JWT secrets
agent tokens
private keys
certificates containing private material
production assessment data
```

Use:

```text
.env.example
```

as the public configuration template.

---

# 27. Testing

Backend:

```powershell
cd backend
go test ./...
```

Agent:

```powershell
cd agent
go test ./...
```

Frontend:

```powershell
cd frontend
npm run build
```

Backend health:

```powershell
curl.exe http://127.0.0.1:8080/health
```

LDAP connectivity:

```powershell
Test-NetConnection 192.168.128.10 -Port 389
```

---

# 28. Troubleshooting

## PostgreSQL: database already exists

If you see:

```text
ERROR: database "redlab" already exists
```

do not recreate it.

Test the connection:

```powershell
psql -h 127.0.0.1 -p 5432 -U postgres -d redlab
```

---

## Backend does not compile

Run:

```powershell
cd backend
go mod tidy
go test ./...
```

Do not ignore compiler errors. RedLab depends on several layers:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
PostgreSQL
```

A mismatch between those layers must be resolved before starting the server.

---

## Backend cannot connect to PostgreSQL

Verify PostgreSQL:

```powershell
Test-NetConnection 127.0.0.1 -Port 5432
```

Verify credentials:

```powershell
psql -h 127.0.0.1 -p 5432 -U postgres -d redlab
```

Then check the `.env` database configuration.

---

## LDAP connection fails

Run:

```powershell
Test-NetConnection YOUR_DC_IP -Port 389
```

Verify:

```text
LDAP_HOST
LDAP_PORT
LDAP_USERNAME
LDAP_PASSWORD
LDAP_BASE_DN
LDAP_USE_TLS
```

Also verify that the configured account has the permissions required by the assessment operation.

---

## Agent does not become ONLINE

Check:

```text
REDLAB_SERVER_URL
REDLAB_AGENT_ID
REDLAB_AGENT_TOKEN
```

Then verify backend connectivity:

```powershell
curl.exe http://127.0.0.1:8080/health
```

Check the agent output for:

```text
heartbeat
poll
authentication
```

Also verify the agent token is valid.

---

# 29. Development Principles

RedLab follows several architectural principles.

### Modular by design

Assessment capabilities should be independent modules.

### Orchestration separated from execution

The backend coordinates work while the Windows agent executes assigned assessment tasks.

### Evidence-driven assessment

Findings should be based on collected technical evidence.

### Security by design

Authentication, authorization, validation, secret management, and logging are treated as core requirements.

### Reproducibility

Laboratory configuration and assessment procedures should be documented so results can be reproduced.

### Extendability

New assessment modules should be added through the existing module registry instead of creating a separate execution path.

---

# 30. Project Status

Current MVP capabilities include:

| Capability                     | Status |
| ------------------------------ | ------ |
| JWT authentication             | ✅      |
| Permission-based authorization | ✅      |
| PostgreSQL persistence         | ✅      |
| Agent creation                 | ✅      |
| Agent token authentication     | ✅      |
| Agent heartbeat                | ✅      |
| Task creation                  | ✅      |
| Task polling                   | ✅      |
| Task execution                 | ✅      |
| Result submission              | ✅      |
| Domain information module      | ✅      |
| User enumeration               | ✅      |
| Group enumeration              | ✅      |
| Computer enumeration           | ✅      |
| SPN enumeration                | ✅      |
| MITRE ATT&CK mapping           | ✅      |
| Findings                       | ✅      |
| Assessment management          | ✅      |
| HTML reporting                 | ✅      |
| Web frontend                   | ✅      |

> Feature completeness should still be validated in the intended Active Directory laboratory before being considered production-ready.

---

# 31. Documentation

Additional documentation:

```text
docs/
├── API.md
├── DEMO.md
├── INSTALLATION.md
└── MITRE-ATTACK.md
```

Recommended reading order:

```text
README.md
   ↓
INSTALLATION.md
   ↓
API.md
   ↓
DEMO.md
   ↓
MITRE-ATTACK.md
```

---

# 32. Security and Responsible Use

RedLab is an offensive-security assessment platform.

Use it only against:

* environments you own;
* laboratory environments;
* systems where you have explicit authorization;
* controlled academic or security-testing environments.

Do not use RedLab to access, enumerate, or compromise systems without authorization.

---

# 33. Quick Reference

### Clone

```powershell
git clone https://github.com/haytamxp/RedLab.git
cd RedLab
```

### Backend

```powershell
cd backend
go mod tidy
go test ./...
go run .\cmd\api
```

### Backend health

```powershell
curl.exe http://127.0.0.1:8080/health
```

### Agent

```powershell
cd agent
go mod tidy
go test ./...
go build -o redlab-agent.exe .\cmd\agent
.\redlab-agent.exe
```

### Frontend

```powershell
cd frontend
npm install
npm run dev
```

### LDAP

```powershell
Test-NetConnection YOUR_DC_IP -Port 389
```

---

# 34. End-to-End Architecture

```text
                         ┌───────────────────┐
                         │     Operator      │
                         │    Web Browser    │
                         └─────────┬─────────┘
                                   │
                              HTTP / REST
                                   │
                                   ▼
                    ┌──────────────────────────┐
                    │      RedLab Backend      │
                    │          Go / Gin        │
                    │                          │
                    │ Auth                     │
                    │ Agents                   │
                    │ Tasks                    │
                    │ Assessments              │
                    │ Findings                 │
                    │ Reporting                │
                    └───────┬────────┬─────────┘
                            │        │
                            │        │
                    PostgreSQL       │ Tasks / Results
                            │        │
                            ▼        ▼
                    ┌───────────┐  ┌───────────────┐
                    │ PostgreSQL│  │ Windows Agent │
                    │  redlab   │  │      Go       │
                    └───────────┘  └───────┬───────┘
                                           │
                                      LDAP / LDAPS
                                           │
                                           ▼
                                ┌────────────────────┐
                                │ Active Directory   │
                                │                    │
                                │ Users              │
                                │ Groups             │
                                │ Computers          │
                                │ SPNs               │
                                │ Domain Information │
                                └────────────────────┘
                                           │
                                           ▼
                                  ┌─────────────────┐
                                  │  MITRE ATT&CK   │
                                  │     Mapping     │
                                  └────────┬────────┘
                                           │
                                           ▼
                                  ┌─────────────────┐
                                  │ Findings & HTML │
                                  │     Report      │
                                  └─────────────────┘
```

---

# 35. License

RedLab is distributed under the license included in this repository.

See:

```text
LICENSE
```

---

## RedLab

**Centralize the assessment.
Automate the workflow.
Collect the evidence.
Map the techniques.
Produce the report.**
