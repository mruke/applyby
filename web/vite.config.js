import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";
/**
 * Defines the Vite configuration for the ApplyBy frontend.
 *
 * This configuration enables React support and configures Vitest to run
 * component tests in a browser-like jsdom environment.
 */
export default defineConfig({
    plugins: [react()],
    test: {
        environment: "jsdom",
        globals: true,
        setupFiles: "./src/test/setup.ts"
    }
});
