# Security Policy & Vulnerability Disclosure

## 1. Supported Versions

| Version | Supported |
| :--- | :--- |
| Current `main` | Yes |
| Prior Major Releases | Critical Security Patches Only |

---

## 2. Reporting a Vulnerability

If you discover a security vulnerability in this project, please **do NOT open a public issue**. Instead, follow these steps:

1. **Email Disclosure:** Send an email with full reproduction steps to `security@yourdomain.com` (or create a Private Vulnerability Report on GitHub).
2. **Include Details:**
   - Detailed description of the vulnerability.
   - Proof of concept (PoC) code or script.
   - Affected versions and configurations.
   - Potential impact and threat assessment.
3. **Response SLA:**
   - Acknowledgement within **24 hours**.
   - Triage and mitigation timeline within **72 hours**.
   - Public patch and CVE disclosure coordination.

---

## 3. Defense-in-Depth Security Invariants

1. **Zero Hardcoded Secrets:** All secrets, private keys, and API tokens are managed via environment variables.
2. **Zero-Trust Multi-Tenancy:** Database queries explicitly scope tenant boundaries (`WHERE tenant_id = $1`).
3. **Constant-Time Cryptography:** Password and token comparisons use constant-time comparison primitives.
