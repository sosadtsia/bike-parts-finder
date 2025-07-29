/**
 * This file is kept for backward compatibility.
 *
 * The actual application now uses Next.js and the main entry point is in:
 * - pages/_app.js (for the application shell)
 * - pages/index.js (for the homepage)
 * - pages/parts/[partId].js (for part details)
 *
 * This file can be safely removed once all references to it are updated.
 */

console.warn(
  'Legacy entry point detected. ' +
  'The application now uses Next.js with pages/_app.js as the main entry point.'
);
// If you want to start measuring performance in this app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
