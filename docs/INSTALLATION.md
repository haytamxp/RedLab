RedLab Installation Guide
1. Requirements

Recommended environment:

Windows 10/11 or Windows Server
Go
PostgreSQL
Node.js
Git
VMware or another hypervisor
Windows Server 2022 laboratory domain controller
Active Directory Domain Services
DNS
2. Clone the Repository
git clone https://github.com/haytamxp/RedLab.git
cd RedLab
3. PostgreSQL

Create the database:

CREATE DATABASE redlab;

Verify:

psql -h 127.0.0.1 -p 5432 -U postgres -d redlab
4. Backend Configuration

Go to:

cd D:\PFA-RedLab\Project\RedLab\backend

Create the local environment file:

Copy-Item .\.env.example .\.env

Configure:

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

The real .env must remain local and must not be committed.

5. Backend

Install dependencies:

cd D:\PFA-RedLab\Project\RedLab\backend
go mod tidy

Run tests:

go test ./...

Start:

go run .\cmd\api

Health check:

curl.exe http://127.0.0.1:8080/health

Expected response:

{
  "success": true,
  "message": "Backend is healthy"
}
6. Active Directory

Start the authorized Windows Server 2022 domain controller.

Verify LDAP connectivity:

Test-NetConnection 192.168.128.10 -Port 389

Expected:

TcpTestSucceeded : True

The agent must be able to communicate with the domain controller over the configured LDAP interface.

7. Windows Agent

Go to:

cd D:\PFA-RedLab\Project\RedLab\agent

Install dependencies:

go mod tidy

Run tests:

go test ./...

Build:

go build -o redlab-agent.exe .\cmd\agent

Configure the agent:

$env:REDLAB_SERVER_URL = "http://127.0.0.1:8080"
$env:REDLAB_AGENT_ID = "YOUR_AGENT_ID"
$env:REDLAB_AGENT_TOKEN = "YOUR_AGENT_TOKEN"
$env:REDLAB_HEARTBEAT_SECONDS = "30"
$env:REDLAB_POLL_SECONDS = "10"

Run:

.\redlab-agent.exe

Expected heartbeat:

[heartbeat] ... status=ONLINE
8. Frontend

Go to:

cd D:\PFA-RedLab\Project\RedLab\frontend

Install dependencies:

npm install

Development server:

npm run dev

Production build:

npm run build
9. Startup Order

Recommended order:

1. PostgreSQL
2. Windows Server 2022 / Active Directory
3. RedLab Backend
4. Windows Agent
5. Frontend
10. Verification

Backend:

curl.exe http://127.0.0.1:8080/health

LDAP:

Test-NetConnection 192.168.128.10 -Port 389

Agent:

[heartbeat] ... status=ONLINE

Frontend:

Dashboard loads successfully.
11. Security Notes

Never commit:

.env
database passwords
LDAP passwords
JWT secrets
agent tokens
private keys

Use:

.env.example

as the public configuration template.

For production or exposed environments:

Use HTTPS
Use strong credentials
Rotate secrets
Restrict network exposure
Use a dedicated AD service account
