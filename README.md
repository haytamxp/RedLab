

# RedLab

RedLab is a modular Active Directory offensive security assessment platform designed for authorized security assessments and laboratory environments.

## Project Objective

RedLab centralizes authorized Active Directory security assessments through a backend, Windows assessment agent, MITRE ATT&CK mapping, findings, reporting, and a web dashboard.

The platform demonstrates a complete assessment workflow:

Assessment
    |
    v
Task Creation
    |
    v
Windows Agent
    |
    v
AD Enumeration
    |
    v
Evidence
    |
    v
MITRE ATT&CK Mapping
    |
    v
Finding
    |
    v
Report
MVP Scope
JWT authentication
Role-based access control
PostgreSQL persistence
Windows agent registration
Agent token authentication
Agent heartbeat
Task creation and execution
Task polling
Task result submission
Active Directory user enumeration
Active Directory computer enumeration
SPN enumeration
MITRE ATT&CK mapping
Finding persistence
Assessment management
HTML reporting
Web dashboard
Architecture
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
                  +-----------------+-----------------+
                  |                 |                 |
                  v                 v                 v
          +---------------+   +-------------+   +-------------+
          |  PostgreSQL   |   | Task System |   |  Reporting  |
          +---------------+   +------+------+   +-------------+
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
Technology Stack
Component	Technology
Backend	Go, Gin
Database	PostgreSQL
Directory	Active Directory / LDAP
Agent	Go / Windows
Frontend	JavaScript / Vite
Authentication	JWT + Agent Tokens
Framework	MITRE ATT&CK
Supported Assessment Modules
Module	Purpose	MITRE ATT&CK
AD_USER_ENUMERATION	Enumerate domain users	T1087.002
AD_COMPUTER_ENUMERATION	Enumerate domain computers	T1018
SPN_ENUMERATION	Enumerate service principal names	T1558.003
Quick Start
Backend
cd backend
go mod tidy
go test ./...
go run .\cmd\api

Health check:

curl.exe http://127.0.0.1:8080/health
Windows Agent
cd agent
go mod tidy
go test ./...
go build -o redlab-agent.exe .\cmd\agent
.\redlab-agent.exe
Frontend
cd frontend
npm install
npm run dev

Production build:

npm run build
Project Structure
RedLab/
|
+-- backend/
|   +-- cmd/
|   +-- internal/
|   +-- migrations/
|   +-- .env
|   +-- .env.example
|
+-- agent/
|   +-- cmd/
|   +-- internal/
|   +-- scripts/
|
+-- frontend/
|   +-- src/
|   +-- package.json
|
+-- docs/
|   +-- ARCHITECTURE.md
|   +-- INSTALLATION.md
|   +-- API.md
|   +-- DEMO.md
|   +-- MITRE-ATTACK.md
|
+-- README.md
Documentation
Architecture
Installation
API Reference
PFA Demonstration
MITRE ATT&CK Integration
Security

RedLab is intended only for authorized security assessments, academic projects, cybersecurity laboratories, and research environments.

Sensitive configuration must remain in the local .env file.

Never commit:

Database passwords
LDAP passwords
JWT secrets
Agent registration tokens
Private keys

Use .env.example as the configuration template.

MITRE ATT&CK

RedLab maps selected assessment activities to MITRE ATT&CK techniques:

AD_USER_ENUMERATION → T1087.002
AD_COMPUTER_ENUMERATION → T1018
SPN_ENUMERATION → T1558.003

The mapping is stored with the corresponding finding.

See MITRE-ATTACK.md for details.

Reporting

RedLab generates an HTML assessment report containing:

Assessment information
Agent information
Assessment status
Findings
MITRE ATT&CK mappings
Recommendations
Project Status

RedLab currently focuses on a functional MVP rather than a complete offensive-security automation framework.

The MVP demonstrates the complete lifecycle:

Task
  |
  v
Agent
  |
  v
AD Enumeration
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
Security Disclaimer

RedLab must only be used against systems for which the operator has explicit authorization.

The project is intended for:

Academic projects
Cybersecurity laboratories
Authorized penetration testing
Security research environments
