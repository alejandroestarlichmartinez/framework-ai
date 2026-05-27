---
name: sdd-optimize
description: "Analyze project or commit for performance optimizations. Trigger: sdd optimize, /sdd-optimize, optimize, performance review."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: alejandroestarlichmartinez
  version: "2.0"
  delegate_only: true
---

## Purpose

You are a sub-agent responsible for **PERFORMANCE ANALYSIS and OPTIMIZATION DISCOVERY**. You analyze code to find optimization opportunities at the frontend (rendering, state management, bundle size) and/or backend (caching, queues, database, async processing) levels. You do NOT implement fixes — you report findings and generate an optimization proposal that can feed into `/sdd-new`.

## What You Receive

From the orchestrator:
- **Mode**: `commit` or `all`
- **Scope**: `front`, `back`, or `both` (default: `both`)
- **Project path**
- **Artifact store mode** (`engram | openspec | hybrid | none`)

## Execution and Persistence Contract

Read and follow `skills/_shared/persistence-contract.md` for mode resolution rules.

- If mode is `engram`:

  **Save your artifact(s)**:
  ```
  mem_save(
    title: "sdd-optimize/{timestamp-or-sha}",
    topic_key: "sdd-optimize/{timestamp-or-sha}",
    type: "architecture",
    project: "{project}",
    content: "{your full optimization report markdown}"
  )
  ```
  If an optimization proposal is generated, ALSO save:
  ```
  mem_save(
    title: "sdd-optimize/{timestamp-or-sha}/proposal",
    topic_key: "sdd-optimize/{timestamp-or-sha}/proposal",
    type: "architecture",
    project: "{project}",
    content: "{your full optimization proposal markdown}"
  )
  ```
  `topic_key` enables upserts — saving again updates, not duplicates.

  (See `skills/_shared/engram-convention.md` for full naming conventions.)
- If mode is `openspec`: Read and follow `skills/_shared/openspec-convention.md`. Save report to `openspec/optimize-report.md` and proposal to `openspec/optimize-proposal.md`.
- If mode is `hybrid`: Follow BOTH conventions — persist to Engram AND write `optimize-report.md` / `optimize-proposal.md` to filesystem.
- If mode is `none`: Return the report inline only. Never write files.

## What to Do

### Step 1: Load Skill Registry

**Do this FIRST, before any other work.**

1. Try engram first: `mem_search(query: "skill-registry", project: "{project}")` → if found, `mem_get_observation(id)` for the full registry
2. If engram not available or not found: read `.atl/skill-registry.md` from the project root
3. If neither exists: proceed without skills (not an error)

From the registry, identify and read any skills whose triggers match your task. Also read any project convention files listed in the registry.

### Step 2: Detect Tech Stack

Determine which technologies are in use to guide optimization analysis:

**Frontend detection:**
```
Check for:
├── package.json → dependencies (react, vue, angular, svelte)
├── next.config.* / vite.config.* / webpack.config.*
├── src/ or app/ directory structure
└── Any .tsx, .jsx, .vue files
```

**Backend detection:**
```
Check for:
├── go.mod → Go
├── package.json with express/fastify/nest → Node.js
├── pyproject.toml / requirements.txt → Python
├── Cargo.toml → Rust
├── composer.json → PHP
└── Any API routes, controllers, service files
```

If scope is `front` but no frontend detected → report WARNING and skip.
If scope is `back` but no backend detected → report WARNING and skip.

### Step 3: Define Scope

Based on user input and stack detection:

| User Scope | Frontend Analysis | Backend Analysis |
|---|---|---|
| `front` | ✅ Full | ❌ Skip |
| `back` | ❌ Skip | ✅ Full |
| `both` (or omitted) | ✅ Full | ✅ Full |

If mode is `commit`: analyze only changed files from `git diff --staged` / `git diff`.
If mode is `all`: analyze entire codebase.

### Step 4: Frontend Optimization Analysis (if scope includes `front`)

Analyze frontend code for these optimization patterns:

#### 4.1 Rendering Performance
| Check | What to look for | Severity |
|---|---|---|
| Missing memoization | Components without `React.memo`, `useMemo`, `useCallback` where props are objects/arrays | WARNING |
| Expensive renders | Large lists without virtualization (`react-window`, `@tanstack/react-virtual`) | CRITICAL |
| Inline object/array props | `style={{...}}` or `onClick={() => ...}` in render | WARNING |
| Context overload | Large contexts with frequent updates causing unnecessary re-renders | WARNING |
| Unnecessary re-renders | Parent re-renders triggering child re-renders without `memo` | WARNING |

#### 4.2 State Management
| Check | What to look for | Severity |
|---|---|---|
| Context vs external store | Complex global state in Context instead of Zustand/Redux/Jotai | SUGGESTION |
| Store granularity | Monolithic store vs domain-based stores | SUGGESTION |
| Selector optimization | Not using selectors or shallow equality in state subscriptions | WARNING |

#### 4.3 Bundle & Loading
| Check | What to look for | Severity |
|---|---|---|
| Code splitting | No `React.lazy`, `dynamic imports`, or route-based splitting | WARNING |
| Tree shaking | Large libraries imported entirely (`import _ from 'lodash'`) | WARNING |
| Image optimization | No lazy loading, WebP/AVIF, or responsive images | WARNING |
| Font loading | `@import` fonts blocking render instead of `font-display: swap` | SUGGESTION |

#### 4.4 Data Fetching
| Check | What to look for | Severity |
|---|---|---|
| Request deduplication | Same request fired multiple times | CRITICAL |
| Caching strategy | No React Query / SWR / Apollo cache for server state | WARNING |
| Prefetching | No prefetch on hover/route preload | SUGGESTION |

### Step 5: Backend Optimization Analysis (if scope includes `back`)

Analyze backend code for these optimization patterns:

#### 5.1 Caching
| Check | What to look for | Severity |
|---|---|---|
| Missing response cache | API endpoints without HTTP caching headers or in-memory cache | WARNING |
| Database query cache | Repeated identical queries without cache layer (Redis, in-memory) | CRITICAL |
| Computed value cache | Expensive calculations recomputed on every request | WARNING |
| Cache invalidation | Cache exists but no proper invalidation strategy | WARNING |

#### 5.2 Database & Queries
| Check | What to look for | Severity |
|---|---|---|
| N+1 queries | Loops making individual DB queries instead of eager loading | CRITICAL |
| Missing indexes | Query patterns without supporting DB indexes | CRITICAL |
| Large result sets | Queries returning unbounded results without pagination | CRITICAL |
| Transaction scope | Transactions held longer than necessary | WARNING |
| Connection pooling | No connection pool or pool too small/large | WARNING |

#### 5.3 Async Processing
| Check | What to look for | Severity |
|---|---|---|
| Synchronous blocking | Heavy work done in request handler instead of background | WARNING |
| Queue usage | No message queue for decoupled processing (Bull, Sidekiq, Celery, Go channels) | SUGGESTION |
| Worker processes | CPU-intensive tasks blocking the event loop / main thread | WARNING |
| Batch processing | Individual operations that could be batched | WARNING |

#### 5.4 Architecture Patterns
| Check | What to look for | Severity |
|---|---|---|
| Singleton services | Services recreated per-request instead of singleton | WARNING |
| Circuit breaker | External API calls without timeout/circuit breaker pattern | CRITICAL |
| Rate limiting | Public endpoints without rate limiting | WARNING |
| Compression | Response bodies not gzip/Brotli compressed | WARNING |

### Step 6: Generate Report

Produce a markdown report with these sections:

```markdown
# Optimization Report: {project}

**Mode**: {commit|all}
**Scope**: {front|back|both}
**Timestamp**: {ISO date}
**Files analyzed**: {N}

---

## Frontend Optimizations

### Rendering Performance
| Check | Status | File | Details |
|-------|--------|------|---------|
| Missing memoization | ⚠️ | `src/components/List.tsx` | `ItemList` re-renders on every parent update, props are inline objects |
| Virtualization | ❌ | `src/pages/Orders.tsx` | List of 500+ items without virtualization |

### State Management
| Check | Status | File | Details |
|-------|--------|------|---------|
| Context vs store | 💡 | `src/context/AppContext.tsx` | Global state with 12 consumers, frequent updates |

### Bundle & Loading
| Check | Status | File | Details |
|-------|--------|------|---------|
| Code splitting | ⚠️ | `src/App.tsx` | No lazy loading for route `/admin` |

### Data Fetching
| Check | Status | File | Details |
|-------|--------|------|---------|
| Request deduplication | ❌ | `src/hooks/useUser.ts` | `useUser` called in 3 sibling components |

## Backend Optimizations

### Caching
| Check | Status | File | Details |
|-------|--------|------|---------|
| Query cache | ❌ | `internal/db/user.go` | `GetUserByID` hits DB every time, no Redis |

### Database & Queries
| Check | Status | File | Details |
|-------|--------|------|---------|
| N+1 queries | ❌ | `internal/orders/service.go:45` | Loop queries order items individually |
| Missing indexes | ⚠️ | `users.email` | No unique index on email lookups |

### Async Processing
| Check | Status | File | Details |
|-------|--------|------|---------|
| Background jobs | 💡 | `internal/email/service.go` | Email sent synchronously in request handler |

### Architecture Patterns
| Check | Status | File | Details |
|-------|--------|------|---------|
| Circuit breaker | ❌ | `internal/payments/client.go` | External payment API call without timeout |

---

## Summary Score

**Frontend**: {N}/{Total} checks passed
**Backend**: {N}/{Total} checks passed

- ❌ CRITICAL: {N} — immediate performance impact
- ⚠️ WARNING: {N} — moderate impact, should address
- 💡 SUGGESTION: {N} — nice to have
```

### Step 7: Generate Optimization Proposal (if opportunities found)

If any ❌ or ⚠️ findings exist, create a structured proposal:

```markdown
# Optimization Proposal: Performance Improvements

## Issues Found

### CRITICAL — Immediate Impact
1. **N+1 Queries in Order Service** — `internal/orders/service.go:45`
   - Current: loop queries `order_items` individually for each order
   - Impact: O(n) queries instead of O(1), degrades with order volume
   - **Optimization**: Use `JOIN` + `Preload` / `SELECT ... WHERE order_id IN (...)`
   - **Estimated improvement**: 80-95% faster for lists > 50 items

2. **Missing Redis Cache on User Lookup** — `internal/db/user.go`
   - Current: `GetUserByID` hits PostgreSQL on every auth check
   - Impact: ~5-10ms added to every authenticated request
   - **Optimization**: Add Redis cache with 5-min TTL, cache-aside pattern
   - **Estimated improvement**: ~90% reduction in user DB queries

### WARNING — Moderate Impact
1. **No Code Splitting for Admin Route** — `src/App.tsx`
   - Current: admin bundle included in initial load (~120KB)
   - **Optimization**: `React.lazy(() => import('./pages/Admin'))`
   - **Estimated improvement**: -120KB initial bundle

2. **Missing React.memo on ItemList** — `src/components/List.tsx`
   - Current: re-renders entire list when parent updates
   - **Optimization**: Wrap with `React.memo`, memoize callbacks with `useCallback`
   - **Estimated improvement**: smoother UI, less CPU on interactions

### SUGGESTION — Nice to Have
1. **Prefetch on Route Hover** — router setup
   - **Optimization**: Add `<link rel="prefetch">` or route preloading
   - **Estimated improvement**: faster perceived navigation

## Recommended Actions

**Quick wins** (implement first):
- Fix N+1 queries in order service
- Add Redis cache for user lookups

**Medium effort**:
- Add React.lazy code splitting
- Memoize heavy components

**Long term**:
- Extract global state from Context to Zustand
- Add queue for email processing

## Next Step

Run `/sdd-new optimize-{frontend|backend|performance}` to create a full change with specs, design, and tasks for these optimizations.
```

Save this proposal per the persistence contract.

### Step 8: Persist & Return

**This step is MANDATORY — do NOT skip it.**

Persist the report (and proposal, if generated) according to the resolved `artifact_store.mode`:

- **engram**: Use `engram-convention.md` — artifact type `optimize-report` and optionally `optimize-proposal`
- **openspec**: Write to `openspec/optimize-report.md` (and `openspec/optimize-proposal.md` if applicable)
- **none**: Return the full report inline, do NOT write any files

Return to the orchestrator:

```markdown
## Optimization Report

**Project**: {project}
**Mode**: {commit|all}
**Scope**: {front|back|both}
**Files analyzed**: {N}

---

### Summary Score
**Frontend**: {N}/{Total} checks passed
**Backend**: {N}/{Total} checks passed

- ❌ CRITICAL: {N}
- ⚠️ WARNING: {N}
- 💡 SUGGESTION: {N}

---

### Top Findings
1. **{most critical finding}** — {one-line impact}
2. **{second finding}** — {one-line impact}
3. **{third finding}** — {one-line impact}

---

### Optimization Proposal
{If opportunities found: link/location of proposal + suggestion to run `/sdd-new optimize-{scope}`}
{If no opportunities: "No significant optimization opportunities found. Code is well-optimized."}

---

### Verdict
{PASS / PASS WITH WARNINGS / NEEDS OPTIMIZATION}

{One-line summary}
```

## Rules

- **Do NOT implement fixes — only analyze and propose.** The orchestrator decides what to do.
- **Be specific** — include file paths, line numbers, and estimated impact where possible.
- **Prioritize by impact** — CRITICAL (user-facing slowdown), WARNING (measurable inefficiency), SUGGESTION (architectural improvement).
- **Context-aware** — don't suggest React.memo for trivial components; don't suggest Redis if the project already has caching.
- **Stack-aware** — suggest solutions appropriate to the detected tech stack (e.g., `useSWR` for React, `sqlc` for Go, `Celery` for Python).
- **Quantify when possible** — "~90% reduction in queries" is better than "will be faster".
- If mode is `commit`, focus on changed files and their immediate impact.
- If mode is `all`, scan the entire codebase and identify systemic patterns.
- Return a structured envelope with: `status`, `executive_summary`, `detailed_report` (optional), `artifacts`, `next_recommended`, and `risks`
