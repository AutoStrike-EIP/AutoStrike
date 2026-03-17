# Contributing to AutoStrike

## Objective

AutoStrike is a breach and attack simulation platform.
All contributions must adhere to a strict framework:
- No actual destructive actions
- Only controlled simulations
- Compliance with security best practices

---

## Setup

Clone the repository:
```bash
git clone https://github.com/{{REPO_NAME}}
cd autostrike
```

---

## Guidelines

- Idiomatic code
- Compliance with standards
- No unnecessary dependencies
- Readable and documented code
- No dangerous logic (no actual malicious execution)

---

## Contribution Workflow

1. Fork the repository
2. Create a branch:
> ```bash
> git checkout -b feature/feature-name
> ```
3. Development
4. Clear commit message:
> ```bash
> # Message foramt: "type(scope): description"
> # Example :
> git commit -m "feat(module/simulation): add new simulation module"
> ```
5. Push:
> ```bash
> git push origin feature/feature-name
> ```
6. Open a Pull Request

---

## Commit conventions

Your commits should follow this format: `type(scope): desciption`

Here is the exhaustive list of the `type`:
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Improvement without functional changes
- `docs`: Documentation
- `chore`: Maintenance

Here is a non-exhaustive list of the `scope`:
- `workflow/ci-cd.yml`
- `simulation/port-scan`
- `README.md`
- *...*

> [!NOTE]
> Use explicit descriptions.

---

## Code Quality

Don't forget to run the linter before any PR:
```bash
*linter-cmd*
```

---

## Scope of Contributions

Accepted:
- new simulation scenarios
- performance improvements
- bug fixes
- documentation improvements

Rejected:
- destructive code
- real-world offensive tools
- Features outside the BAS scope

---

## Questions

Use GitHub Discussions or open an issue.