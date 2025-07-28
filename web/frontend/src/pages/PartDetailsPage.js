import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getPartById } from '../services/api';

function PartDetailsPage() {
  const { partId } = useParams();
  const [part, setPart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPartDetails = async () => {
      setLoading(true);
      setError(null);

      try {
        const data = await getPartById(partId);
        setPart(data);
      } catch (err) {
        setError(err.message || 'Failed to fetch part details');
      } finally {
        setLoading(false);
      }
    };

    fetchPartDetails();
  }, [partId]);

  if (loading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="loader"></div>
        <p className="ml-4">Loading part details...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <p>{error}</p>
        <Link to="/" className="text-blue-600 hover:text-blue-800 mt-2 inline-block">
          Back to search
        </Link>
      </div>
    );
  }

  if (!part) {
    return (
      <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded">
        <p>Part not found</p>
        <Link to="/" className="text-blue-600 hover:text-blue-800 mt-2 inline-block">
          Back to search
        </Link>
      </div>
    );
  }

  // Build image carousel
  const renderImages = () => {
    if (part.images && part.images.length > 0) {
      return (
        <div className="flex overflow-x-auto space-x-4 mb-4 p-2">
          {part.images.map((img, index) => (
            <div key={index} className="flex-shrink-0 w-64 h-64">
              <img
                src={img}
                alt={`${part.brand} ${part.model} - view ${index + 1}`}
                className="w-full h-full object-contain"
              />
            </div>
          ))}
        </div>
      );
    }

    return (
      <div className="w-full h-64 bg-gray-200 mb-4">
        <img
          src="https://via.placeholder.com/640x480"
          alt={`${part.brand} ${part.model}`}
          className="w-full h-full object-contain"
        />
      </div>
    );
  };

  // Stock status badge
  const stockStatus = part.inStock ? (
    <span className="inline-block bg-green-100 text-green-800 px-2 py-1 rounded text-xs">
      In Stock
    </span>
  ) : (
    <span className="inline-block bg-red-100 text-red-800 px-2 py-1 rounded text-xs">
      Out of Stock
    </span>
  );

  // Price information
  const renderPriceInfo = () => {
    let priceInfo = <span className="text-2xl font-bold">${part.price.toFixed(2)}</span>;

    if (part.msrp > 0 && part.msrp > part.price) {
      const discountPercentage = Math.round(((part.msrp - part.price) / part.msrp) * 100);

      priceInfo = (
        <div>
          {priceInfo}
          <div className="flex items-center">
            <span className="line-through text-gray-500 mr-2">${part.msrp.toFixed(2)}</span>
            <span className="inline-block bg-green-100 text-green-800 px-2 py-1 rounded text-xs">
              Save {discountPercentage}%
            </span>
          </div>
        </div>
      );
    }

    return priceInfo;
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      <div className="p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold">{part.brand} {part.model}</h2>
          <Link to="/" className="text-blue-600 hover:text-blue-800">
            Back to results
          </Link>
        </div>

        <div className="flex flex-col md:flex-row">
          <div className="md:w-1/2 mb-4 md:mb-0 md:pr-6">
            {renderImages()}
          </div>

          <div className="md:w-1/2">
            <div className="mb-4">
              <p className="text-sm text-gray-600">{part.category} &gt; {part.subCategory}</p>
              <div className="flex items-center mt-2">
                {stockStatus}
                <div className="ml-4">{renderPriceInfo()}</div>
              </div>
            </div>

            <hr className="my-4" />

            <div className="mb-4">
              <h3 className="font-bold mb-2">Description</h3>
              <p>{part.description}</p>
            </div>

            <div className="mb-4">
              <a
                href={part.url}
                target="_blank"
                rel="noopener noreferrer"
                className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded inline-block"
              >
                View on retailer site
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default PartDetailsPage;
