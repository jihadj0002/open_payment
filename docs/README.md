# Open Payment Gateway — Documentation Index

> **Last Updated:** 2026-07-16
> **Status Key:** ✅ Accurate | ⚠️ Partial/Spec (some items not implemented) | ❌ Not Implemented

## Project Overview (`01-project-overview/`)
| File | Status | Description |
|------|--------|-------------|
| `01-introduction.md` | ✅ | Project introduction and vision |
| `02-business-requirements.md` | ✅ | Business requirements |
| `03-goals-and-success-metrics.md` | ✅ | Goals and success metrics |
| `04-scope-and-limitations.md` | ✅ | Scope and limitations |

## Requirements (`02-requirements/`)
| File | Status | Description |
|------|--------|-------------|
| `01-functional-requirements.md` | ✅ | Functional requirements |
| `02-non-functional-requirements.md` | ✅ | Non-functional requirements |
| `03-user-stories.md` | ✅ | User stories |
| `04-compliance-requirements.md` | ⚠️ Partially updated | PCI DSS partial, GDPR/PSD2 TODO |

## System Architecture (`03-system-architecture/`)
| File | Status | Description |
|------|--------|-------------|
| `01-high-level-architecture.md` | ✅ | High-level architecture |
| `02-microservices-boundaries.md` | ⚠️ | Modular monolith (not microservices) |
| `03-data-flow-diagrams.md` | ✅ | Data flow diagrams |
| `04-payment-state-machine.md` | ✅ | Payment state machine |
| `05-event-driven-architecture.md` | ✅ | Event-driven architecture |
| `06-database-architecture.md` | ✅ | Database architecture |
| `07-caching-strategy.md` | ✅ | Caching strategy (Redis) |
| `08-deployment-architecture.md` | ⚠️ | Deployment architecture (Railway, not AWS) |

## API Design (`04-api-design/`)
| File | Status | Description |
|------|--------|-------------|
| `01-api-overview.md` | ✅ Updated | API overview, auth, pagination, rate limiting |
| `02-merchant-rest-api.md` | ✅ Updated | Merchant REST API spec |
| `03-admin-rest-api.md` | ✅ Updated | Admin REST API spec |
| `04-webhooks.md` | ✅ Updated | Webhook events, retry, signature verification |
| `05-internal-grpc-api.md` | ✅ | Internal gRPC API (not implemented) |
| `06-error-handling.md` | ✅ Updated | Error taxonomy, codes, idempotency |

## Database Design (`05-database-design/`)
| File | Status | Description |
|------|--------|-------------|
| `01-er-diagrams.md` | ✅ Updated | ER diagrams with implementation status |
| `02-merchant-schema.md` | ✅ Updated | Merchant tables (migrations 001, 003, 010, 016) |
| `03-payment-schema.md` | ✅ Updated | Payment tables (migrations 001, 004, 008, 009, 011, 014) |
| `04-ledger-schema.md` | ⚠️ Updated | Ledger schema (simplified implementation) |
| `05-customer-schema.md` | ⚠️ Updated | Customer schema (migrations 001, 004, 014) |
| `06-security-schema.md` | ✅ Updated | Security tables (migrations 007, 013) |
| `07-indexing-and-partitioning.md` | ⚠️ Updated | Indexing strategy (partitioning not implemented) |

## Frontend Spec (`06-frontend-spec/`)
| File | Status | Description |
|------|--------|-------------|
| `01-design-system.md` | ✅ Updated | UI component library spec |
| `02-merchant-dashboard.md` | ✅ Updated | Merchant dashboard pages |
| `03-admin-dashboard.md` | ❌ TODO | Admin dashboard (not built) |
| `04-developer-portal.md` | ❌ TODO | Developer portal (not built) |
| `05-state-and-data-flow.md` | ✅ | State management with React Query + Zustand |
| `06-routing-and-navigation.md` | ✅ Updated | App router structure |
| `07-mobile-app-spec.md` | ❌ TODO | Mobile app (not built) |

## Security (`07-security/`)
| File | Status | Description |
|------|--------|-------------|
| `01-encryption-model.md` | ✅ | Encryption model |
| `02-authentication.md` | ✅ | Authentication |
| `03-authorization.md` | ✅ | Authorization |
| `04-request-security.md` | ✅ | Request security |
| `05-infrastructure-security.md` | ✅ | Infrastructure security |
| `06-audit-and-threat-model.md` | ✅ | Audit and threat model |

## Infrastructure (`08-infrastructure/`)
| File | Status | Description |
|------|--------|-------------|
| `01-kubernetes-setup.md` | ⚠️ | K8s setup (Railway, not self-managed K8s) |
| `02-terraform-modules.md` | ⚠️ | Terraform modules |
| `03-ci-cd-pipeline.md` | ✅ Updated | CI/CD pipeline (GitHub Actions) |
| `04-monitoring-stack.md` | ⚠️ | Monitoring stack |
| `05-disaster-recovery.md` | ✅ | Disaster recovery plan |
| `06-scaling-policies.md` | ✅ | Scaling policies |

## Implementation (`09-implementation/`)
| File | Status | Description |
|------|--------|-------------|
| `01-golang-modular-monolith.md` | ✅ | Go modular monolith architecture |
| `02-service-implementation.md` | ✅ | Service implementation details |
| `03-frontend-implementation.md` | ✅ | Frontend implementation |
| `04-testing-strategy.md` | ✅ | Testing strategy |
| `05-code-quality.md` | ✅ | Code quality standards |
| `06-performance-targets.md` | ✅ | Performance targets |

## Project Management (`10-project-management/`)
| File | Status | Description |
|------|--------|-------------|
| `01-sprint-roadmap.md` | ⚠️ | Sprint roadmap (may need update) |
| `02-detailed-task-list.md` | ⚠️ | Detailed task list (may need update) |
| `03-team-requirements.md` | ✅ | Team requirements |
| `04-risk-register.md` | ✅ | Risk register |
| `05-delivery-milestones.md` | ✅ | Delivery milestones |

## Workflow (`11-workflow/`)
| File | Status | Description |
|------|--------|-------------|
| `01-step-by-step-build-process.md` | ✅ | Build process |
| `02-first-pr-to-production.md` | ✅ | PR to production workflow |
| `03-operational-runbook.md` | ✅ | Operational runbook |

## Task Board
| File | Status | Description |
|------|--------|-------------|
| `00-task-board.md` | ✅ Updated | Main task board with all phases |
| `01-backlog.md` | ✅ | Backlog items |
| `02-completed.md` | ✅ | Completed tasks (historical) |

## Deployment
| File | Status | Description |
|------|--------|-------------|
| `local-development.md` | ✅ Updated | Local development guide |
| `production-deployment.md` | ⚠️ Updated | Production deployment (Railway) |
| `README.md` | ✅ | Deployment overview |

## Other Docs
| File | Status | Description |
|------|--------|-------------|
| `api/openapi.yml` | ✅ | OpenAPI 3.0 spec |
| `core_principles.md` | ✅ | Core principles |
| `project-capabilities.md` | ✅ | Project capabilities |
| `ideas.md` | ✅ | Ideas and roadmap |
| `work_process.md` | ✅ | Work process |
| `payment_gateway_processor_and_security_explained.md` | ✅ | Processor/security explainer |

## Agents
| File | Status | Description |
|------|--------|-------------|
| `agents/README.md` | ✅ | Agent system overview |
| `agents/00-orchestrator.md` | ✅ | Orchestrator agent |
| `agents/01-backend-engineer.md` | ✅ | Backend engineer agent |
| `agents/02-frontend-engineer.md` | ✅ | Frontend engineer agent |
| `agents/03-devops-engineer.md` | ✅ | DevOps engineer agent |
| `agents/04-security-engineer.md` | ✅ | Security engineer agent |
| `agents/05-qa-engineer.md` | ✅ | QA engineer agent |
| `agents/06-product-manager.md` | ✅ | Product manager agent |
| `agents/07-compliance-officer.md` | ✅ | Compliance officer agent |

## Runbooks
| File | Status | Description |
|------|--------|-------------|
| `runbooks/database-recovery.md` | ✅ | Database recovery |
| `runbooks/processor-downtime.md` | ✅ | Processor downtime |
| `runbooks/webhook-backlog.md` | ✅ | Webhook backlog |
| `runbooks/rate-limit-breach.md` | ✅ | Rate limit breach |
