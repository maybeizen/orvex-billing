# Orvex Billing

A monorepo containing the billing API and frontend for Orvex.

## Project Structure

```
├── api/          # Go backend API server (Written in Go)
├── frontend/     # Next.js frontend application
├── turbo.json    # Turbo configuration
└── package.json  # Root package configuration
```

## Quick Start

### Prerequisites

- Node.js 18+
- Go 1.24+
- npm/yarn/pnpm

### Installation

```bash
# Install dependencies
npm install

# Start development servers
npm run dev
```

## Available Scripts

- `npm run dev` - Start all development servers
- `npm run build` - Build all packages
- `npm run lint` - Lint all packages
- `npm run test` - Run tests for all packages
- `npm run type-check` - Type check all packages

## Workspaces

### API (`/api`)

Go-based backend API server using Fiber framework.

### Frontend (`/frontend`)

Next.js application with TypeScript and Tailwind CSS.

## Development

This project uses Turbo for fast builds and development. Each workspace can be developed independently or together using the root scripts.
