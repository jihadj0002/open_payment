# Open Payment Gateway — Agent Workflow System

## Overview

This defines the multi-agent workflow system for building the payment gateway. The system has three layers:

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORCHESTRATOR AGENT                            │
│  Role: Project Lead — assigns tasks, tracks progress,           │
│         maintains logbook, ensures compliance                   │
└──────────────┬──────────────────────────────────┬───────────────┘
               │                                  │
    ┌──────────┴──────────┐            ┌──────────┴──────────┐
    │     Task Board      │            │       Logbook       │
    │  (Active tasks)     │            │  (Completed work)   │
    └──────────┬──────────┘            └──────────┬──────────┘
               │                                  │
    ┌──────────┴──────────────────────────────────┴──────────┐
    │                SPECIALIZED AGENTS                       │
    │                                                        │
    │  Backend Eng  │  Frontend Eng  │  DevOps Eng           │
    │  Security Eng │  QA Eng        │  Product Manager      │
    │  Compliance   │  Designer      │                       │
    └────────────────────────────────────────────────────────┘
```

## Protocols (ALL AGENTS MUST COMPLY)

### P1: Task Lifecycle
Every task flows through these states:
```
BACKLOG → TODO → IN_PROGRESS → REVIEW → DONE
                                    ↕
                                 BLOCKED
```

### P2: Task Card Format
Every task card MUST follow this exact format:
```yaml
---
task_id: TASK-XXX
title: Short descriptive title
status: TODO | IN_PROGRESS | REVIEW | DONE | BLOCKED | BACKLOG
priority: CRITICAL | HIGH | MEDIUM | LOW
assignee: agent_role
created: YYYY-MM-DD
started: YYYY-MM-DD
deadline: YYYY-MM-DD
depends_on: [TASK-YYY]
tags: [backend, payment, api]

description: |
  Detailed description of what needs to be done.

acceptance_criteria:
  - Criterion 1
  - Criterion 2

implementation_notes: |
  Notes from the assignee about implementation.

review_notes: |
  Notes from code review.

completed_at: YYYY-MM-DD
---
```

### P3: Logbook Entry Format
Every completed work session MUST be logged:
```yaml
---
date: YYYY-MM-DD
agent: agent_role
tasks_worked: [TASK-XXX, TASK-YYY]
duration_hours: X

work_done: |
  What was accomplished during this session.

decisions_made:
  - Decision 1 with rationale

blockers_found:
  - Blocker 1

next_steps:
  - Step 1
---
```

### P4: Communication Protocol
1. Orchestrator creates tasks with clear acceptance criteria
2. Agent picks task from TODO → moves to IN_PROGRESS
3. Agent works and logs progress daily
4. Agent moves to REVIEW when ready
5. Orchestrator or peer reviews
6. On approval: move to DONE
7. On changes needed: move back to IN_PROGRESS

### P5: Escalation
If a task is BLOCKED for > 24 hours:
- Agent logs the blocker in their logbook entry
- Agent moves task to BLOCKED
- Orchestrator is notified and works to unblock

## File Locations
| Component | Path |
|-----------|------|
| Agent Definitions | `docs/agents/` |
| Task Board | `docs/task-board/` |
| Task Templates | `docs/task-board/templates/` |
| Logbook | `docs/logbook/` |
| OpenCode Skills | `.opencode/agents/` |
| OpenCode Tasks | `.opencode/tasks/` |
| OpenCode Logbook | `.opencode/logbook/` |
