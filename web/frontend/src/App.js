import React from 'react';

/**
 * This component is kept for backward compatibility.
 * The actual application now uses Next.js and the main entry point is in pages/_app.js
 * This file can be safely removed once all references to it are updated.
 */
function App() {
  return (
    <div className="p-4 bg-yellow-100 border border-yellow-400 text-yellow-700 rounded">
      <h1 className="text-xl font-bold mb-2">App Component (Legacy)</h1>
      <p>
        This is the legacy App component from the Create React App version.
        The application now uses Next.js with pages/_app.js as the main entry point.
      </p>
    </div>
  );
}

export default App;
