# Touch Grass Frontend

Next.js-based dashboard for tracking and visualizing WiFi hotspot discoveries.

## Features
- Real-time hotspot discovery notifications
- Interactive map visualization
- User profile management
- Adventure statistics tracking
- Device purchase and management
- WebSocket-based live updates

## Tech Stack
- Framework: Next.js 14
- UI Components: NextUI v2
- Styling: Tailwind CSS
- Maps: React Simple Maps, React Leaflet
- Authentication: Supabase
- Build Tool: PNPM

## Installation

1. Clone the repository
2. Install dependencies:
```bash
pnpm install
```

3. Configure environment variables in `.env.local`:
```bash
NEXT_PUBLIC_SUPABASE_URL=your_supabase_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
NEXT_PUBLIC_API_URL=your_api_url
NEXT_PUBLIC_NEXTJS_URL=your_nextjs_url
NEXT_PUBLIC_GOLANG_URL=your_golang_url
```

## Development
```bash
pnpm dev     # Start development server
pnpm build   # Build for production
pnpm start   # Start production server
```

## Docker Deployment
```bash
docker build --build-arg NEXT_PUBLIC_API_URL=http://your-api-url -t my-nextjs-app .
docker run -p 3000:3000 my-nextjs-app
```

## Project Structure
```
frontend/
├── src/
│   ├── app/           # Next.js app router pages
│   ├── components/    # React components
│   ├── utils/        # Utility functions
│   ├── pages/        # API routes
│   └── data/         # Static data files
├── public/           # Static assets
└── ...config files
```

## Key Components
- Landing Page: Main entry point with parallax effects
- Interactive Map: Australia map visualization
- Discover Page: Location discovery and tracking
- Profile Page: User profile management
- Stats Page: User statistics and achievements

## Authentication
Supabase authentication implementation includes:
- Email/Password authentication
- Email verification
- Password reset functionality
- Session management
- Protected routes

## Styling and UI
- Tailwind CSS for utility-first styling
- NextUI components for consistent UI elements
- Custom glass-effect components
- Mobile-first responsive design

## Performance Optimizations
- Dynamic imports
- Lazy-loaded map components
- Optimized asset delivery
- Client-side navigation
- Server-side rendering where appropriate

## State Management
- React hooks for component state
- Supabase real-time subscriptions
- Server-side data fetching
- Optimistic updates
