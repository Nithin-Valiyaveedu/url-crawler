# URL Crawler Frontend

A React TypeScript frontend for the URL Crawler application with real-time crawling analysis and visualization.

## Features

- Submit URLs for crawling analysis
- Real-time status updates with automatic polling
- Paginated, sortable, and filterable results table
- Detailed analysis view with interactive charts
- Bulk operations (delete, re-run)
- Responsive design with modern UI

## API Integration

This frontend connects to the Go backend API. Make sure your backend is running before starting the frontend.

### Environment Configuration

Create a `.env` file in the frontend directory with the following configuration:

```bash
# Backend API Configuration
VITE_API_BASE_URL=http://localhost:8080
VITE_API_KEY=your-api-key-here

# Environment
VITE_ENV=development
```

### Available Environment Variables

- `VITE_API_BASE_URL`: Backend API base URL (default: http://localhost:8080)
- `VITE_API_KEY`: API key for authentication (default: dev-api-key-12345)
- `VITE_ENV`: Environment mode (development/production)

## Setup and Installation

1. Install dependencies:

```bash
npm install
```

2. Configure environment variables (create `.env` file as shown above)

3. Start the development server:

```bash
npm run dev
```

4. Open [http://localhost:5173](http://localhost:5173) in your browser

## Backend Requirements

Make sure the Go backend is running with the following endpoints available:

- `POST /api/crawl` - Submit URL for crawling
- `GET /api/crawl` - Get paginated crawl results
- `GET /api/crawl/:id` - Get specific crawl result
- `DELETE /api/crawl` - Bulk delete crawl results
- `POST /api/crawl/rerun` - Re-run crawl analysis
- `GET /api/crawl/stats` - Get crawl statistics
- `GET /api/health` - Health check

## API Authentication

The frontend uses Bearer token authentication. Set your API key in the environment variables:

```bash
VITE_API_KEY=your-api-key-here
```

The API key should match the one configured in your Go backend.

## Development

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Technologies Used

- React 18 with TypeScript
- Vite for build tooling
- Tailwind CSS for styling
- Recharts for data visualization
- Lucide React for icons

## Real-time Updates

The application automatically polls the backend every 5 seconds when there are active crawling tasks (running or queued status). This ensures the UI stays up-to-date with the latest crawling progress.
