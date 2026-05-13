/**
 * LoadingStateProps
 *
 * Defines the message shown while a page or section is loading.
 */
type LoadingStateProps = {
  message?: string;
};

/**
 * LoadingState
 *
 * Displays an accessible loading message for asynchronous frontend states.
 */
export function LoadingState({ message = "Loading..." }: LoadingStateProps) {
  return (
    <section className="state-card" aria-live="polite">
      <h2>Loading</h2>
      <p>{message}</p>
    </section>
  );
}