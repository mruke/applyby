/**
 * createClientId
 *
 * Creates a client-side identity for local create workflows.
 * The fallback exists for browsers or tests that do not expose crypto.randomUUID.
 */
export function createClientId(prefix: string): string {
  if (globalThis.crypto && "randomUUID" in globalThis.crypto) {
    return globalThis.crypto.randomUUID();
  }

  return `${prefix}-${Date.now()}`;
}