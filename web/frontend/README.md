# Bike Parts Finder Frontend

This is the React frontend for the Bike Parts Finder application. It allows users to search for bike parts across multiple online retailers.

## Technologies Used

- React 18
- React Router for navigation
- Axios for API requests
- Tailwind CSS for styling

## Getting Started

### Prerequisites

- Node.js 16+ and npm

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm start
```

The application will be available at [http://localhost:3000](http://localhost:3000).

### Building for Production

```bash
npm run build
```

This will create an optimized production build in the `build` folder.

### Using Docker

You can also use Docker for development:

```bash
docker build -f Dockerfile.dev -t bike-parts-finder-frontend .
docker run -p 3000:3000 -v $(pwd):/app bike-parts-finder-frontend
```

## Project Structure

- `/public` - Static assets and HTML template
- `/src` - React source code
  - `/components` - Reusable UI components
  - `/pages` - Page components
  - `/services` - API services
