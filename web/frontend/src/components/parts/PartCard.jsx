import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

const PartCard = ({ part }) => {
  return (
    <div className="card hover:shadow-lg transition-shadow">
      <div className="flex flex-col md:flex-row gap-4">
        <div className="w-full md:w-36 h-36 bg-gray-200 rounded-md flex items-center justify-center">
          {part.imageUrl ? (
            <img
              src={part.imageUrl}
              alt={part.name}
              className="max-w-full max-h-full object-contain"
            />
          ) : (
            <span className="text-gray-400">No image</span>
          )}
        </div>

        <div className="flex-1">
          <h2 className="text-xl font-semibold mb-2">
            <Link to={`/parts/${part.id}`} className="text-primary-600 hover:underline">
              {part.name}
            </Link>
          </h2>

          <div className="mb-2 flex items-center">
            <span className="bg-gray-200 text-gray-700 px-2 py-1 rounded text-sm mr-2">
              {part.brand}
            </span>
            <span className="text-gray-500">{part.category}</span>
          </div>

          <p className="text-gray-700 mb-4 line-clamp-2">
            {part.description || 'No description available'}
          </p>

          <div className="flex flex-wrap justify-between items-center">
            <div>
              <p className="text-2xl font-bold text-primary-600">
                ${part.price.toFixed(2)}
              </p>
              {part.originalPrice && part.originalPrice > part.price && (
                <p className="text-sm text-gray-500 line-through">
                  ${part.originalPrice.toFixed(2)}
                </p>
              )}
            </div>

            <div className="flex items-center">
              {part.inStock ? (
                <span className="text-green-600 mr-3 text-sm font-medium">In Stock</span>
              ) : (
                <span className="text-red-600 mr-3 text-sm font-medium">Out of Stock</span>
              )}
              <Link to={`/parts/${part.id}`} className="btn btn-primary">
                View Details
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

PartCard.propTypes = {
  part: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    brand: PropTypes.string.isRequired,
    category: PropTypes.string.isRequired,
    description: PropTypes.string,
    price: PropTypes.number.isRequired,
    originalPrice: PropTypes.number,
    imageUrl: PropTypes.string,
    inStock: PropTypes.bool
  }).isRequired
};

export default PartCard;
