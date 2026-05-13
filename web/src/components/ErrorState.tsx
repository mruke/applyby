/**
 * ErrorStateProps
 *
 * Defines the title and message shown when a page or section cannot load.
 */
type ErrorStateProps = {
  title?: string;
  message: string;
};

/**
 * ErrorState
 *
 * Displays an accessible error message for failed frontend states.
 */
export function ErrorState({ title = "Something needs attention", message }: ErrorStateProps) {
  return (
    <section className="state-card" role="alert">
      <h2>{title}</h2>
      <p>{message}</p>
    </section>
  );
}