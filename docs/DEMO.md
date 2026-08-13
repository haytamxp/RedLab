RedLab MVP Demonstration
1. Objective

The demonstration shows the complete lifecycle of an authorized Active Directory security assessment.

Login
  |
  v
Create Assessment
  |
  v
Create Task
  |
  v
Agent Polls
  |
  v
AD Module Executes
  |
  v
Evidence Returned
  |
  v
MITRE ATT&CK Mapping
  |
  v
Finding Created
  |
  v
Assessment Completed
  |
  v
HTML Report
  |
  v
Dashboard
2. Start PostgreSQL

Ensure PostgreSQL is running and the RedLab database exists.

Database:

redlab
3. Start Active Directory

Start the Windows Server 2022 laboratory domain controller.

Verify LDAP connectivity:

Test-NetConnection 192.168.128.10 -Port 389

Expected:

TcpTestSucceeded : True
4. Start the Backend
cd D:\PFA-RedLab\Project\RedLab\backend
go run .\cmd\api

Verify:

curl.exe http://127.0.0.1:8080/health
5. Start the Windows Agent

Configure the agent environment variables and run:

cd D:\PFA-RedLab\Project\RedLab\agent
.\redlab-agent.exe

Expected:

[heartbeat] ... status=ONLINE
6. Login

PowerShell:

cd D:\PFA-RedLab\Project\RedLab\backend

$loginBody = @{
    username = "redlabadmin"
    password = "YOUR_PASSWORD"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod `
    -Method POST `
    -Uri "http://127.0.0.1:8080/api/v1/auth/login" `
    -ContentType "application/json" `
    -Body $loginBody

$adminJwt = $loginResponse.token
7. Create Assessment
$assessmentBody = @{
    name = "RedLab AD Security Assessment"
    description = "Authorized Active Directory offensive security assessment."
    agent_id = "YOUR_AGENT_ID"
} | ConvertTo-Json

$assessmentResponse = Invoke-RestMethod `
    -Method POST `
    -Uri "http://127.0.0.1:8080/api/v1/assessments" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    } `
    -ContentType "application/json" `
    -Body $assessmentBody

$assessmentId = $assessmentResponse.data.id

Write-Host "Assessment ID:" $assessmentId
8. Start the Assessment
$statusBody = @{
    status = "RUNNING"
} | ConvertTo-Json

Invoke-RestMethod `
    -Method PATCH `
    -Uri "http://127.0.0.1:8080/api/v1/assessments/$assessmentId/status" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    } `
    -ContentType "application/json" `
    -Body $statusBody
9. AD User Enumeration

Create:

$taskBody = @{
    agent_id = "YOUR_AGENT_ID"
    type = "AD_USER_ENUMERATION"
    payload = @{
        assessment_id = $assessmentId
        base_dn = "DC=redlab,DC=local"
    }
    priority = 10
} | ConvertTo-Json -Depth 10

$taskResponse = Invoke-RestMethod `
    -Method POST `
    -Uri "http://127.0.0.1:8080/api/v1/tasks" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    } `
    -ContentType "application/json" `
    -Body $taskBody

$taskResponse | ConvertTo-Json -Depth 10

Expected MITRE mapping:

T1087.002
Domain Account Discovery
10. AD Computer Enumeration

Create:

$taskBody = @{
    agent_id = "YOUR_AGENT_ID"
    type = "AD_COMPUTER_ENUMERATION"
    payload = @{
        assessment_id = $assessmentId
        base_dn = "DC=redlab,DC=local"
    }
    priority = 10
} | ConvertTo-Json -Depth 10

Expected mapping:

T1018
Remote System Discovery
11. SPN Enumeration

Create:

$taskBody = @{
    agent_id = "YOUR_AGENT_ID"
    type = "SPN_ENUMERATION"
    payload = @{
        assessment_id = $assessmentId
        base_dn = "DC=redlab,DC=local"
    }
    priority = 10
} | ConvertTo-Json -Depth 10

$taskResponse = Invoke-RestMethod `
    -Method POST `
    -Uri "http://127.0.0.1:8080/api/v1/tasks" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    } `
    -ContentType "application/json" `
    -Body $taskBody

Expected mapping:

T1558.003
Steal or Forge Kerberos Tickets: Kerberoasting

The MVP performs SPN discovery and evidence collection. It does not automate password cracking.

12. Verify Findings
$findings = Invoke-RestMethod `
    -Method GET `
    -Uri "http://127.0.0.1:8080/api/v1/findings" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    }

$findings | ConvertTo-Json -Depth 10
13. Complete the Assessment
$completeBody = @{
    status = "COMPLETED"
} | ConvertTo-Json

Invoke-RestMethod `
    -Method PATCH `
    -Uri "http://127.0.0.1:8080/api/v1/assessments/$assessmentId/status" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    } `
    -ContentType "application/json" `
    -Body $completeBody
14. Generate the HTML Report
$report = Invoke-WebRequest `
    -UseBasicParsing `
    -Method GET `
    -Uri "http://127.0.0.1:8080/api/v1/assessments/$assessmentId/report" `
    -Headers @{
        Authorization = "Bearer $adminJwt"
    }

$report.Content | Set-Content `
    -Path "D:\PFA-RedLab\Project\RedLab\RedLab-Assessment.html" `
    -Encoding UTF8

Start-Process `
    "D:\PFA-RedLab\Project\RedLab\RedLab-Assessment.html"
15. Frontend Demonstration

Start:

cd D:\PFA-RedLab\Project\RedLab\frontend
npm run dev

The dashboard should display:

Backend status
Registered agents
Assessments
Findings
MITRE ATT&CK techniques
Report access
16. Final Demonstration

The final PFA demonstration should show:

Administrator
     |
     v
Assessment
     |
     v
Task
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
MITRE ATT&CK
     |
     v
Finding
     |
     v
Report
     |
     v
Dashboard
