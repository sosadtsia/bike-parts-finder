import React from 'react';
import { Link } from 'react-router-dom';

const NotFound = () => {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <div className="text-9xl font-bold text-gray-200 mb-4">404</div>
      <h1 className="text-3xl font-bold mb-4">Page Not Found</h1>
      <p className="text-gray-600 mb-8 text-center max-w-md">
        The page you are looking for might have been removed, had its name changed,
        or is temporarily unavailable.
      </p>
      <div className="flex flex-col sm:flex-row gap-4">
        <Link to="/" className="btn btn-primary">
          Go to Home
        </Link>
        <Link to="/search" className="btn btn-outline">
          Browse Bike Parts
        </Link>
      </div>
    </div>
  );
};

export default NotFound;
