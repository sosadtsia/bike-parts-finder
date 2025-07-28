import React from 'react';
import PartCard from './PartCard';

function PartsList({ parts, loading, error }) {
  if (loading) {
    return (
      <div className="col-span-full flex justify-center items-center py-8">
        <div className="loader"></div>
        <p className="ml-4">Loading parts...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="col-span-full bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <p>{error}</p>
      </div>
    );
  }

  if (parts.length === 0) {
    return (
      <div className="col-span-full bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded">
        <p>No parts found matching your search criteria.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {parts.map(part => (
        <PartCard key={part.id} part={part} />
      ))}
    </div>
  );
}

export default PartsList;
