import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getPartById } from '../services/api';

const PartDetail = () => {
  const { partId } = useParams();
  const [part, setPart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('details');
  const [selectedRetailer, setSelectedRetailer] = useState(null);

  useEffect(() => {
    const fetchPartDetails = async () => {
      setLoading(true);
      setError(null);

      try {
        const partData = await getPartById(partId);
        setPart(partData);

        if (partData.retailers && partData.retailers.length > 0) {
          setSelectedRetailer(partData.retailers[0]);
        }
      } catch (err) {
        setError('Failed to load part details. Please try again later.');
        console.error('Error fetching part details:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPartDetails();
  }, [partId]);

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
        <p className="mt-2 text-gray-600">Loading part details...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <p className="font-bold">Error</p>
        <p>{error}</p>
        <Link to="/search" className="mt-4 inline-block btn btn-primary">
          Back to Search
        </Link>
      </div>
    );
  }

  if (!part) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold mb-4">Part not found</h2>
        <p className="mb-6">The part you're looking for doesn't exist or has been removed.</p>
        <Link to="/search" className="btn btn-primary">
          Back to Search
        </Link>
      </div>
    );
  }

  return (
    <div>
      {/* Breadcrumbs */}
      <div className="text-sm text-gray-600 mb-6">
        <Link to="/" className="hover:underline">Home</Link>
        <span className="mx-2">/</span>
        <Link to="/search" className="hover:underline">Search</Link>
        <span className="mx-2">/</span>
        <span className="font-medium">{part.name}</span>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-5 gap-8">
        {/* Left Column - Images */}
        <div className="lg:col-span-2">
          <div className="bg-white p-4 rounded-lg shadow-md">
            {part.imageUrl ? (
              <img
                src={part.imageUrl}
                alt={part.name}
                className="w-full h-auto object-contain mb-4"
              />
            ) : (
              <div className="h-80 bg-gray-200 flex items-center justify-center rounded-md">
                <span className="text-gray-500">No image available</span>
              </div>
            )}

            {part.additionalImages && part.additionalImages.length > 0 && (
              <div className="grid grid-cols-4 gap-2 mt-4">
                {part.additionalImages.map((img, index) => (
                  <div
                    key={index}
                    className="aspect-square bg-gray-100 rounded cursor-pointer"
                    onClick={() => {
                      // Handle image gallery functionality
                    }}
                  >
                    <img
                      src={img}
                      alt={`${part.name} view ${index + 1}`}
                      className="w-full h-full object-contain"
                    />
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Right Column - Part Info */}
        <div className="lg:col-span-3">
          <div className="bg-white p-6 rounded-lg shadow-md mb-6">
            <h1 className="text-3xl font-bold mb-2">{part.name}</h1>

            <div className="flex items-center mb-4">
              <span className="bg-primary-100 text-primary-800 px-3 py-1 rounded-full text-sm font-medium mr-2">
                {part.brand}
              </span>
              <span className="text-gray-600">{part.category}</span>
            </div>

            <div className="mb-6">
              <div className="flex items-baseline mb-2">
                <span className="text-3xl font-bold text-primary-600 mr-3">
                  ${part.price?.toFixed(2) || 'Varies by retailer'}
                </span>
                {part.originalPrice && part.originalPrice > part.price && (
                  <span className="text-gray-500 line-through">
                    ${part.originalPrice.toFixed(2)}
                  </span>
                )}
              </div>

              {part.inStock ? (
                <span className="text-green-600 font-medium">In Stock</span>
              ) : (
                <span className="text-red-600 font-medium">Out of Stock</span>
              )}
            </div>

            <div className="space-y-4">
              <div className="border-t border-gray-200 pt-4">
                <h3 className="font-semibold mb-2">Retailers</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {part.retailers && part.retailers.length > 0 ? (
                    part.retailers.map(retailer => (
                      <div
                        key={retailer.id}
                        className={`border rounded-md p-3 cursor-pointer transition-colors ${
                          selectedRetailer && selectedRetailer.id === retailer.id
                            ? 'border-primary-500 bg-primary-50'
                            : 'border-gray-300 hover:border-primary-300'
                        }`}
                        onClick={() => setSelectedRetailer(retailer)}
                      >
                        <div className="flex justify-between items-center">
                          <span className="font-medium">{retailer.name}</span>
                          <span className="font-bold text-primary-600">
                            ${retailer.price.toFixed(2)}
                          </span>
                        </div>
                        <div className="flex justify-between items-center text-sm mt-1">
                          <span className={retailer.inStock ? 'text-green-600' : 'text-red-600'}>
                            {retailer.inStock ? 'In Stock' : 'Out of Stock'}
                          </span>
                          <span className="text-gray-600">
                            {retailer.shippingDays
                              ? `Ships in ${retailer.shippingDays} days`
                              : 'Shipping info unavailable'}
                          </span>
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className="text-gray-600 col-span-2">No retailer information available</p>
                  )}
                </div>
              </div>

              <div className="pt-4">
                {selectedRetailer && (
                  <a
                    href={selectedRetailer.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="btn btn-primary w-full text-center"
                  >
                    Buy Now from {selectedRetailer.name}
                  </a>
                )}
              </div>
            </div>
          </div>

          {/* Tabs for additional info */}
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <div className="flex border-b">
              <button
                className={`flex-1 py-3 font-medium ${
                  activeTab === 'details'
                    ? 'bg-white text-primary-600 border-b-2 border-primary-500'
                    : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                }`}
                onClick={() => setActiveTab('details')}
              >
                Details & Specs
              </button>
              <button
                className={`flex-1 py-3 font-medium ${
                  activeTab === 'compatibility'
                    ? 'bg-white text-primary-600 border-b-2 border-primary-500'
                    : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                }`}
                onClick={() => setActiveTab('compatibility')}
              >
                Compatibility
              </button>
              <button
                className={`flex-1 py-3 font-medium ${
                  activeTab === 'reviews'
                    ? 'bg-white text-primary-600 border-b-2 border-primary-500'
                    : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                }`}
                onClick={() => setActiveTab('reviews')}
              >
                Reviews
              </button>
            </div>

            <div className="p-6">
              {activeTab === 'details' && (
                <div>
                  <p className="mb-4">{part.description || 'No detailed description available.'}</p>

                  {part.specs && Object.keys(part.specs).length > 0 ? (
                    <div className="mt-4">
                      <h3 className="font-semibold mb-2">Specifications</h3>
                      <div className="border rounded-md overflow-hidden">
                        <table className="w-full">
                          <tbody>
                            {Object.entries(part.specs).map(([key, value]) => (
                              <tr key={key} className="border-b last:border-b-0">
                                <td className="py-2 px-4 bg-gray-50 font-medium">{key}</td>
                                <td className="py-2 px-4">{value}</td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  ) : (
                    <p className="text-gray-600">No specifications available.</p>
                  )}
                </div>
              )}

              {activeTab === 'compatibility' && (
                <div>
                  {part.compatibility ? (
                    <div>
                      <h3 className="font-semibold mb-2">Compatible With</h3>
                      <ul className="list-disc pl-5">
                        {part.compatibility.map((item, index) => (
                          <li key={index} className="mb-1">{item}</li>
                        ))}
                      </ul>
                    </div>
                  ) : (
                    <p className="text-gray-600">Compatibility information not available.</p>
                  )}
                </div>
              )}

              {activeTab === 'reviews' && (
                <div>
                  {part.reviews && part.reviews.length > 0 ? (
                    <div>
                      <div className="flex items-center mb-4">
                        <span className="text-xl font-bold mr-2">{part.rating || 0}</span>
                        <div className="flex text-yellow-400">
                          {[1, 2, 3, 4, 5].map((star) => (
                            <svg
                              key={star}
                              xmlns="http://www.w3.org/2000/svg"
                              className={`h-5 w-5 ${star <= (part.rating || 0) ? 'text-yellow-500' : 'text-gray-300'}`}
                              viewBox="0 0 20 20"
                              fill="currentColor"
                            >
                              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                            </svg>
                          ))}
                        </div>
                        <span className="ml-2 text-gray-600">{part.reviews.length} reviews</span>
                      </div>

                      <div className="space-y-4">
                        {part.reviews.map((review) => (
                          <div key={review.id} className="border-b pb-4 last:border-0">
                            <div className="flex justify-between items-center mb-2">
                              <div className="font-medium">{review.user || 'Anonymous'}</div>
                              <div className="text-sm text-gray-500">{new Date(review.date).toLocaleDateString()}</div>
                            </div>
                            <div className="flex text-yellow-400 mb-2">
                              {[1, 2, 3, 4, 5].map((star) => (
                                <svg
                                  key={star}
                                  xmlns="http://www.w3.org/2000/svg"
                                  className={`h-4 w-4 ${star <= review.rating ? 'text-yellow-500' : 'text-gray-300'}`}
                                  viewBox="0 0 20 20"
                                  fill="currentColor"
                                >
                                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                                </svg>
                              ))}
                            </div>
                            <p className="text-gray-800">{review.comment}</p>
                          </div>
                        ))}
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <p className="text-gray-600 mb-2">No reviews yet.</p>
                      <p className="text-gray-600">Be the first to review this product!</p>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PartDetail;
