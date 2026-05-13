/**
 * EmptyStateProps
 *
 * Defines the title and message shown when a page or section has no data.
 */
type EmptyStateProps = {
  title: string;
  message: string;
};

/**
 * EmptyState
 *
 * Displays a clear empty state so the user understands what is missing
 * and what the section will eventually contain.
 */
export function EmptyState({ title, message }: EmptyStateProps) {
  return (
    <section className="state-card">
      <h2>{title}</h2>
      <p>{message}</p>
    </section>
  );
}