# Orvex Billing

A modern billing platform monorepo built with Express and Next.js, inspired by [Nether Host](https://netherhost.cc) and [Paymenter](https://paymenter.org).

## Tech Stack

- **Backend**: Express.js with TypeScript
- **Frontend**: Next.js with TypeScript and Tailwind CSS
- **Build System**: Turbo for fast builds and caching
- **Package Manager**: pnpm workspaces

## Project Structure

```
├── api/              # Express-based Backend API
│   ├── src/          # Source code
│   ├── package.json  # API dependencies
│   └── .env.example  # Environment template
├── frontend/         # Next.js frontend application
│   ├── src/          # Source code
│   ├── package.json  # Frontend dependencies
│   └── .env.local.example
├── turbo.json        # Turbo configuration
├── package.json      # Root package configuration
└── README.md         # This file
```

## Quick Start

### Prerequisites

- Node.js 18+
- pnpm 8+ (recommended) or npm 9+
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/maybeizen/orvex-billing.git
cd orvex-billing

# Install all dependencies (recommended)
pnpm install

# Or using npm
npm install

# Set up environment variables
cp apps/api/.env.example apps/api/.env

# Start development servers
pnpm run dev

# Or using npm
npm run dev
```

The application will be available at:

- **Frontend**: http://localhost:3000
- **API**: http://localhost:3001

## Available Scripts

### Root Level Commands

- `pnpm run dev` - Start all development servers
- `pnpm run build` - Build all packages
- `pnpm run lint` - Lint all packages
- `pnpm run test` - Run tests for all packages
- `pnpm run type-check` - Type check all packages
- `pnpm run clean` - Clean all build artifacts

_Replace `pnpm` with `npm` if using npm instead of pnpm._

### Workspace-Specific Commands

```bash
# Run commands for specific workspaces (pnpm)
pnpm run dev --filter=api
pnpm run build --filter=frontend
pnpm run lint --filter=api

# Or using npm workspaces
npm run dev --workspace=api
npm run build --workspace=frontend
npm run lint --workspace=api
```

## Workspaces

### API (`/apps/api`)

Express-based backend API server handling:

- User authentication and authorization
- Billing and subscription management
- Payment processing
- Invoice generation

**Key Features:**

- RESTful API design
- TypeScript for type safety
- Environment-based configuration

### Frontend (`/apps/frontend`)

Next.js application providing:

- User dashboard
- Billing management interface
- Responsive design with Tailwind CSS
- Server-side rendering

**Key Features:**

- Modern React with TypeScript
- Tailwind CSS for styling
- Optimized for performance

## Development

This project uses **Turbo** for optimized builds and development workflows. Turbo provides:

- Fast incremental builds
- Smart caching
- Parallel task execution
- Remote caching capabilities

Each workspace can be developed independently or together using the root-level scripts.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting (`pnpm run test && pnpm run lint`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Deployment

### Production Build

```bash
# Build all packages for production (pnpm recommended)
pnpm run build

# Type check before deployment
pnpm run type-check

# Or using npm
npm run build
npm run type-check
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For questions or support, please open an issue in the repository or contact the development team.
