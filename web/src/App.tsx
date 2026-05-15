import { Route, Routes } from "react-router-dom";

import { AppShell } from "./components/AppShell";
import { ApplicationDetailPage } from "./pages/ApplicationDetailPage";
import { ApplicationEditPage } from "./pages/ApplicationEditPage";
import { ApplicationsPage } from "./pages/ApplicationsPage";
import { ContactEditPage } from "./pages/ContactEditPage";
import { DocumentEditPage } from "./pages/DocumentEditPage";
import { DashboardPage } from "./pages/DashboardPage";
import { NotFoundPage } from "./pages/NotFoundPage";

/**
 * App
 *
 * Defines the top-level frontend route structure for ApplyBy.
 * It keeps routing centralized and wraps all pages in the shared app shell.
 */
export default function App() {
  return (
    <AppShell>
      <Routes>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/applications" element={<ApplicationsPage />} />
        <Route path="/applications/:applicationId/edit" element={<ApplicationEditPage />} />
        <Route path="/applications/:applicationId/contacts/:contactId/edit" element={<ContactEditPage />} />
        <Route path="/applications/:applicationId/documents/:documentId/edit" element={<DocumentEditPage />} />
        <Route path="/applications/:applicationId" element={<ApplicationDetailPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </AppShell>
  );
}
