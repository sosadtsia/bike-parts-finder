import React from 'react';

const About = () => {
  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">About Bike Parts Finder</h1>

      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-2xl font-semibold mb-4">Our Mission</h2>
        <p className="mb-4">
          Bike Parts Finder was created with a simple mission: to help cyclists find the right parts at the best prices.
          Whether you're a casual rider, weekend warrior, or professional cyclist, finding the right components for your
          bike shouldn't be a frustrating experience.
        </p>
        <p className="mb-4">
          By aggregating parts from multiple retailers and providing detailed specifications and compatibility information,
          we aim to simplify the bike part shopping experience and help you make informed decisions.
        </p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-2xl font-semibold mb-4">How It Works</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="text-center p-4">
            <div className="bg-primary-100 text-primary-600 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">Search</h3>
            <p className="text-gray-600">
              Our advanced search engine helps you find the exact parts you need by brand, type, compatibility, and more.
            </p>
          </div>

          <div className="text-center p-4">
            <div className="bg-primary-100 text-primary-600 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">Compare</h3>
            <p className="text-gray-600">
              Compare prices, shipping times, and availability across multiple retailers in one convenient place.
            </p>
          </div>

          <div className="text-center p-4">
            <div className="bg-primary-100 text-primary-600 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">Purchase</h3>
            <p className="text-gray-600">
              Once you've found the right part at the right price, we'll direct you to the retailer to complete your purchase.
            </p>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-2xl font-semibold mb-4">Our Technology</h2>
        <p className="mb-4">
          Bike Parts Finder is built using modern cloud-native technologies, ensuring a fast, reliable experience:
        </p>
        <ul className="list-disc pl-6 mb-4 space-y-2">
          <li>React frontend for a responsive, interactive user interface</li>
          <li>Go backend for high-performance data processing</li>
          <li>PostgreSQL database for reliable data storage</li>
          <li>Redis caching for lightning-fast search results</li>
          <li>Kafka for real-time data streaming</li>
          <li>Deployed on Kubernetes for scalability and reliability</li>
        </ul>
        <p>
          Our web scraping technology continuously monitors bicycle part retailers to ensure prices and availability information
          is always up to date.
        </p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-semibold mb-4">Contact Us</h2>
        <p className="mb-4">
          Have questions, suggestions, or feedback? We'd love to hear from you!
        </p>
        <div className="flex items-center mb-3">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-primary-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          <a href="mailto:contact@bikepartsfinder.com" className="text-primary-600 hover:underline">
            contact@bikepartsfinder.com
          </a>
        </div>
        <div className="flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-primary-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
          </svg>
          <span>(555) 123-4567</span>
        </div>
      </div>
    </div>
  );
};

export default About;
