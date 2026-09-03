# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: Actively supported |
| < 1.0   | :x: Best effort, please upgrade to main |

Security fixes are released as patch versions on `main`.

## Reporting a Vulnerability

**Do not open a public issue.** Public disclosure before a fix puts users at risk.

Please use one of these private channels:

1. **Preferred - GitHub Private Advisory:**  
   Go to **Security tab -> Report a vulnerability**  
   https://github.com/moranricardo/cli/security/advisories/new

2. **Email:** moranmaldonadoricardo@gmail.com  
   Subject: `[SECURITY] cli - brief description`  
   Please include steps to reproduce, impact, and if possible a PoC. PGP optional.

### What to expect

- Acknowledgment within 24h
- Initial assessment within 72h
- We follow coordinated disclosure: please give us 90 days to publish a fix before public disclosure.
- If accepted, we will create a GitHub Security Advisory (GHSA), credit you (unless you prefer anonymous), and release a patch.

### Scope

This policy covers the code in this repository (`moranricardo/cli`). Dependencies are out of scope, but we will help upstream if needed.

Thanks for helping keep this project and its users safe.

Related: See [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) for community conduct (not for security reports).
