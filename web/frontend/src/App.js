import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import HomePage from './pages/HomePage';
import PartDetailsPage from './pages/PartDetailsPage';

function App() {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/parts/:partId" element={<PartDetailsPage />} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
}

export default App;
