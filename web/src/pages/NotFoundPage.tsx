/**
 * NotFoundPage
 *
 * Provides a fallback route when the requested frontend page does not exist.
 */
export function NotFoundPage() {
  return (
    <section className="state-card">
      <h1>Page not found</h1>
      <p>The requested page does not exist.</p>
    </section>
  );
}