import React from 'react';
import Link from 'next/link';

function Footer() {
  return (
    <footer className="bg-gray-800 text-white py-8">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between">
          <div className="mb-4 md:mb-0">
            <h3 className="text-lg font-bold mb-2">Bike Parts Finder</h3>
            <p className="text-sm text-gray-400">Find the best deals on bike parts across multiple online retailers</p>
          </div>
          <div>
            <h4 className="font-bold mb-2">Links</h4>
            <ul className="text-sm text-gray-400">
              <li><Link href="/" className="hover:text-white">Home</Link></li>
              <li><Link href="/" className="hover:text-white">Search</Link></li>
              <li><Link href="/" className="hover:text-white">Categories</Link></li>
              <li><Link href="/" className="hover:text-white">About</Link></li>
            </ul>
          </div>
        </div>
        <div className="mt-8 pt-4 border-t border-gray-700 text-sm text-gray-400">
          <p>&copy; {new Date().getFullYear()} Bike Parts Finder. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
