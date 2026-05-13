import type { ReactNode } from "react";
import { NavLink } from "react-router-dom";

/**
 * AppShellProps
 *
 * Defines the content required by the shared application shell.
 */
type AppShellProps = {
  children: ReactNode;
};

/**
 * AppShell
 *
 * Provides the shared page frame for ApplyBy, including brand navigation,
 * primary navigation, and the semantic main content region.
 */
export function AppShell({ children }: AppShellProps) {
  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="app-header__inner">
          <NavLink className="app-brand" to="/">
            ApplyBy
          </NavLink>

          <nav className="app-nav" aria-label="Primary navigation">
            <NavLink to="/">Dashboard</NavLink>
            <NavLink to="/applications">Applications</NavLink>
          </nav>
        </div>
      </header>

      <main className="app-main">{children}</main>
    </div>
  );
}