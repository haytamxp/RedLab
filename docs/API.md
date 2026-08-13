RedLab API Reference

Base URL:

http://127.0.0.1:8080/api/v1
1. Authentication
Register
POST /auth/register
Content-Type: application/json

Example:

{
  "username": "redlabadmin",
  "email": "redlabadmin@example.local",
  "password": "CHANGE_ME",
  "first_name": "RedLab",
  "last_name": "Administrator"
}
Login
POST /auth/login
Content-Type: application/json

Example:

{
  "username": "redlabadmin",
  "password": "CHANGE_ME"
}

Response:

{
  "token": "JWT_TOKEN"
}

Authenticated requests use:

Authorization: Bearer <JWT_TOKEN>
2. Health
GET /health
GET /api/v1/health

Example response:

{
  "success": true,
  "message": "Backend is healthy",
  "data": {
    "service": "RedLab Backend",
    "status": "running",
    "version": "1.0.0"
  }
}
3. Agents
Create Agent
POST /agents
Authorization: Bearer <ADMIN_JWT>
Content-Type: application/json

Example:

{
  "name": "redlab-agent-01",
  "hostname": "WIN-DC01",
  "ip_address": "192.168.128.20",
  "operating_system": "Windows Server 2022",
  "version": "0.1.0"
}

The endpoint returns a registration token.

The token must be stored securely because it is used to authenticate the agent.

List Agents
GET /agents
Authorization: Bearer <ADMIN_JWT>
Get Agent
GET /agents/:id
Authorization: Bearer <ADMIN_JWT>
Heartbeat
POST /agents/:id/heartbeat
Authorization: Bearer <AGENT_TOKEN>
4. Tasks
Create Task
POST /tasks
Authorization: Bearer <ADMIN_JWT>
Content-Type: application/json

Example:

{
  "agent_id": "<AGENT_ID>",
  "type": "AD_USER_ENUMERATION",
  "payload": {
    "base_dn": "DC=redlab,DC=local"
  },
  "priority": 10
}
Poll Next Task
POST /agents/:id/tasks/next
Authorization: Bearer <AGENT_TOKEN>
List Agent Tasks
GET /agents/:id/tasks
Authorization: Bearer <AGENT_TOKEN>
Submit Task Result
POST /agents/:id/tasks/:taskId/result
Authorization: Bearer <AGENT_TOKEN>
Content-Type: application/json
5. Assessments
Create Assessment
POST /assessments
Authorization: Bearer <ADMIN_JWT>
Content-Type: application/json

Example:

{
  "name": "RedLab AD Security Assessment",
  "description": "Authorized Active Directory security assessment.",
  "agent_id": "<AGENT_ID>"
}
List Assessments
GET /assessments
Authorization: Bearer <ADMIN_JWT>
Get Assessment
GET /assessments/:id
Authorization: Bearer <ADMIN_JWT>
Update Assessment Status
PATCH /assessments/:id/status
Authorization: Bearer <ADMIN_JWT>
Content-Type: application/json

Example:

{
  "status": "RUNNING"
}

Supported statuses:

PENDING
RUNNING
COMPLETED
FAILED
Assessment Report
GET /assessments/:id/report
Authorization: Bearer <ADMIN_JWT>

Returns an HTML assessment report.

6. Findings
List Findings
GET /findings
Authorization: Bearer <ADMIN_JWT>
Get Finding
GET /findings/:id
Authorization: Bearer <ADMIN_JWT>
7. MVP Task Types
AD_USER_ENUMERATION
AD_COMPUTER_ENUMERATION
SPN_ENUMERATION
8. MITRE Mapping
Task Type	MITRE ID	Technique
AD_USER_ENUMERATION	T1087.002	Domain Account Discovery
AD_COMPUTER_ENUMERATION	T1018	Remote System Discovery
SPN_ENUMERATION	T1558.003	Steal or Forge Kerberos Tickets: Kerberoasting
9. Authentication Model

Administrator endpoints use JWT authentication.

Agent endpoints use agent registration tokens.

Agent tokens are stored as hashes in PostgreSQL rather than as plaintext.

10. Common HTTP Status Codes
200 OK
201 Created
204 No Content
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
500 Internal Server Error
