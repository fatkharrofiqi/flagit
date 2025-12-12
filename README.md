# Flagit - Feature Flag Management System

A modern feature flag management system built with Go (Fiber) backend and React frontend.

## Tech Stack

### Backend
- **Go** with **Fiber** web framework
- **PostgreSQL** database
- **golang-migrate** for database migrations
- **Server-Sent Events** for real-time updates

### Frontend
- **React 19** with TypeScript
- **TanStack Router** for routing
- **TanStack Query** for server state management
- **Zustand** for client state management
- **Tailwind CSS** + **shadcn/ui** for styling
- **Vite** for build tooling

## Development

### Prerequisites
1. **PostgreSQL** database running
2. **Go** 1.25+ installed
3. **Node.js** 18+ installed
4. **pnpm** package manager

### Setup

1. **Clone and install dependencies**
   ```bash
   git clone <repository-url>
   cd flagits
   pnpm install
   ```

2. **Configure database**
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

3. **Set up database**
   ```bash
   # Run migrations and seed data
   pnpm setup
   ```

4. **Start development servers**
   ```bash
   # Start both frontend and backend concurrently
   pnpm start
   ```

### Available Scripts

#### Root Scripts
- `pnpm setup` - Run migrations and seed database
- `pnpm start` - Start both frontend and backend
- `pnpm migrate` - Run database migrations up
- `pnpm migrate-down` - Rollback migrations
- `pnpm seed` - Seed database with demo data
- `pnpm api` - Start backend only
- `pnpm dev` - Start frontend only (via turbo)

#### Backend Scripts (apps/api)
- `go run cmd/http/main.go` - Start HTTP server
- `go run cmd/seed/main.go` - Seed database

#### Frontend Scripts (apps/admin)
- `pnpm dev` - Start Vite dev server
- `pnpm build` - Build for production

### API Endpoints

#### Projects
- `GET /api/projects` - List all projects
- `POST /api/projects` - Create new project
- `GET /api/projects/:id` - Get project by ID
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

#### Environments
- `GET /api/environments` - List all environments
- `POST /api/environments` - Create new environment
- `GET /api/projects/:projectId/environments` - Get project environments
- `PUT /api/environments/:id` - Update environment
- `DELETE /api/environments/:id` - Delete environment

#### Flags
- `GET /api/projects/:projectId/flags` - Get project flags
- `POST /api/flags` - Create new flag
- `PUT /api/flags/:id` - Update flag
- `DELETE /api/flags/:id` - Delete flag
- `GET /api/flags/:flagId/values` - Get flag values
- `POST /api/flags/values` - Create/update flag value
- `PUT /api/flags/values/:id` - Update flag value

#### Real-time Updates
- `GET /api/events` - SSE endpoint for real-time flag updates

### Architecture

```
├── apps/
│   ├── admin/           # React frontend
│   └── api/            # Go backend (with package.json)
├── package.json          # Root package with turbo scripts
├── turbo.json           # Turbo configuration
└── README.md
```

### Turbo Scripts Structure

- **Root Scripts** - Orchestrates workspace commands using turbo
- **API Package Scripts** (`apps/api/package.json`) - Contains Go-specific commands
- **Admin Package Scripts** (`apps/admin/package.json`) - Contains frontend commands

### Available Scripts

#### From Root (Recommended)
- `pnpm setup` - Run migrations and seed database
- `pnpm start` - Start both frontend and backend
- `pnpm migrate` - Run migrations via turbo
- `pnpm migrate-down` - Rollback migrations
- `pnpm seed` - Seed database with demo data
- `pnpm api` - Start backend only

#### Direct API Commands
- `cd apps/api && pnpm run dev` - Start backend directly
- `cd apps/api && pnpm run migrate` - Run migrations directly
- `cd apps/api && pnpm run seed` - Seed database directly

#### Direct Admin Commands
- `cd apps/admin && pnpm run dev` - Start frontend directly
- `cd apps/admin && pnpm run build` - Build frontend

## Features

- ✅ **Project Management** - Create and manage feature flag projects
- ✅ **Environment Management** - Set up dev, staging, production environments
- ✅ **Feature Flags** - Create boolean, string, number, and JSON flags
- ✅ **Per-Environment Values** - Different flag values per environment
- ✅ **Real-time Updates** - Instant flag changes via Server-Sent Events
- ✅ **Responsive UI** - Works on desktop and mobile
- ✅ **Type Safety** - Full TypeScript support

## Usage

1. **Create a Project** - Start by creating a new project
2. **Set up Environments** - Add dev, staging, production environments
3. **Create Feature Flags** - Add flags with different data types
4. **Configure Values** - Set flag values for each environment
5. **Toggle Flags** - Enable/disable flags in real-time

## Development Tips

- The backend runs on `http://localhost:8080`
- The frontend runs on `http://localhost:3000`
- Database changes require running migrations
- Use `pnpm check` to lint and format code
- Use `pnpm check-types` for TypeScript checking

## Production Deployment

### Backend
```bash
cd apps/api
go build -o flagits-api ./cmd/http
./flagits-api
```

### Frontend
```bash
cd apps/admin
pnpm build
# Serve the dist/ directory
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `pnpm check` to ensure code quality
5. Submit a pull request

## License

MIT License - see LICENSE file for details

### Quick Setup

For a complete one-command setup:

```bash
# Copy environment config and run migrations + seed
pnpm quick-setup
```

This is equivalent to:
```bash
cp .env.example .env
pnpm migrate
pnpm seed
```
