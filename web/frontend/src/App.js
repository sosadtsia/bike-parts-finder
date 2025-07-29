import React from 'react';
import Header from './components/Header';
import Footer from './components/Footer';
import HomePage from './pages/HomePage';

// This component is kept for compatibility with existing imports
// The actual app component is now in pages/_app.js
function App() {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <HomePage />
      </main>
      <Footer />
    </div>
  );
}

export default App;
