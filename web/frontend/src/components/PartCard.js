import React from 'react';
import { Link } from 'react-router-dom';

function PartCard({ part }) {
  // Get main image or placeholder
  const imageURL = part.images && part.images.length > 0
    ? part.images[0]
    : "https://via.placeholder.com/400x300";

  // Format price
  const price = part.price.toFixed(2);
  const msrp = part.msrp > 0 && part.msrp > part.price
    ? part.msrp.toFixed(2)
    : null;

  return (
    <div className="part-card bg-white rounded-lg shadow-md overflow-hidden">
      <div className="h-48 bg-gray-200">
        <img src={imageURL} alt={`${part.brand} ${part.model}`} className="w-full h-full object-cover" />
      </div>
      <div className="p-4">
        <div className="flex justify-between items-start">
          <div>
            <h3 className="font-bold text-lg">{part.brand} {part.model}</h3>
            <p className="text-sm text-gray-600">{part.category}</p>
          </div>
          <div className="text-right">
            <p className="font-bold text-lg">${price}</p>
            {msrp && <p className="text-sm line-through text-gray-500">${msrp}</p>}
          </div>
        </div>
        <p className="text-sm mt-2">{part.description}</p>
        <div className="mt-4 flex justify-between">
          {part.inStock ? (
            <span className="inline-block bg-green-100 text-green-800 px-2 py-1 rounded text-xs">In Stock</span>
          ) : (
            <span className="inline-block bg-red-100 text-red-800 px-2 py-1 rounded text-xs">Out of Stock</span>
          )}
          <Link to={`/parts/${part.id}`} className="text-blue-600 hover:text-blue-800 text-sm">View Details</Link>
        </div>
      </div>
    </div>
  );
}

export default PartCard;
