# Bike Parts Finder Frontend

This is the React frontend for the Bike Parts Finder application, built with Next.js. It allows users to search for bike parts across multiple online retailers.

## Technologies Used

- React 19
- Next.js 15.4 (build system and routing)
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
npm run dev
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
docker-compose up
```

Or build and run the production container:

```bash
npm run docker:prod
```

## Project Structure

- `/public` - Static assets (images, favicon, manifest)
- `/src` - React application source code
  - `/components` - Reusable UI components
  - `/pages` - Page components (Next.js routing)
  - `/services` - API services
