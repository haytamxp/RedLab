RedLab MITRE ATT&CK Integration
1. Purpose

RedLab maps selected Active Directory assessment activities to MITRE ATT&CK techniques.

The technique mapping is returned as structured assessment metadata and persisted with the resulting finding.

2. Current MVP Mapping
RedLab Module	MITRE ID	Technique
AD_USER_ENUMERATION	T1087.002	Domain Account Discovery
AD_COMPUTER_ENUMERATION	T1018	Remote System Discovery
SPN_ENUMERATION	T1558.003	Steal or Forge Kerberos Tickets: Kerberoasting
3. T1087.002 — Domain Account Discovery

RedLab maps:

AD_USER_ENUMERATION

to:

T1087.002

The module enumerates Active Directory domain users.

Example evidence:

{
  "count": 6,
  "users": [
    {
      "sam_account_name": "redlab.user",
      "user_principal_name": "redlab.user@redlab.local"
    }
  ]
}
4. T1018 — Remote System Discovery

RedLab maps:

AD_COMPUTER_ENUMERATION

to:

T1018

The module enumerates Active Directory computer objects.

Example evidence:

{
  "count": 1,
  "computers": [
    {
      "dns_hostname": "DC01.redlab.local",
      "sam_account_name": "DC01$"
    }
  ]
}
5. T1558.003 — Kerberoasting

RedLab maps:

SPN_ENUMERATION

to:

T1558.003

The module enumerates service principal names and identifies accounts associated with SPNs.

Example evidence:

{
  "spn_count": 21,
  "account_count": 2
}

The current MVP performs discovery and evidence collection.

It does not automate credential cracking.

6. Mapping Architecture
Task
  |
  v
Agent Executor
  |
  v
Assessment Module
  |
  v
Evidence
  |
  v
MITRE Mapping
  |
  v
Backend
  |
  v
Finding
  |
  +----> Dashboard
  |
  +----> Report
7. MITRE Metadata

Example:

{
  "mitre_technique": {
    "id": "T1018",
    "name": "Remote System Discovery"
  }
}
8. Finding Structure

A RedLab finding contains information such as:

title
description
severity
technique_id
technique_name
task_id
agent_id
assessment_id
evidence
recommendation
created_at
updated_at
9. Security Interpretation

MITRE ATT&CK mapping describes the activity performed by RedLab.

It does not automatically prove that a vulnerability or compromise exists.

For example:

SPN discovered
     ≠
Kerberos compromise confirmed

The collected evidence must be interpreted by the security analyst.

10. Why MITRE ATT&CK Is Used

MITRE ATT&CK provides a standardized vocabulary for describing adversary behavior and security activities.

In RedLab, it improves:

Assessment traceability
Reporting consistency
Technique identification
Communication of security findings
Mapping between offensive activity and defensive analysis
11. MVP Scope

The RedLab MVP intentionally implements a small technique catalogue instead of the full ATT&CK knowledge base.

The current focus is:

Active Directory discovery
Structured evidence
MITRE mapping
Finding persistence
Reporting

The catalogue can be extended in future versions.


After pasting them into the six files, **do not add the separator lines (`=====`) to the files**. They are only there so you know where each file starts and ends.

Also, one correction from the text you pasted: use local Markdown links such as:

```markdown
[Architecture](docs/ARCHITECTURE.md)
