import { type FormEvent, useState } from "react";

import type { ApplicationSearchCriteria, ApplicationStatus } from "../types/application";

/**
 * ApplicationSearchFormProps
 *
 * Defines current criteria and search behavior for the applications page.
 */
type ApplicationSearchFormProps = {
  criteria: ApplicationSearchCriteria;
  isSearching: boolean;
  onClear: () => Promise<void>;
  onSearch: (criteria: ApplicationSearchCriteria) => Promise<void>;
};

/**
 * SearchFormState
 *
 * Represents controlled form values for application search.
 */
type SearchFormState = {
  companyName: string;
  source: string;
  status: "" | ApplicationStatus;
  text: string;
};

/**
 * statusOptions
 *
 * Provides readable status options supported by application search.
 */
const statusOptions: { value: "" | ApplicationStatus; label: string }[] = [
  { value: "", label: "Any status" },
  { value: "draft", label: "Draft" },
  { value: "interested", label: "Interested" },
  { value: "applied", label: "Applied" },
  { value: "interviewing", label: "Interviewing" },
  { value: "offer", label: "Offer" },
  { value: "rejected", label: "Rejected" },
  { value: "withdrawn", label: "Withdrawn" },
  { value: "archived", label: "Archived" }
];

/**
 * formStateFromCriteria
 *
 * Converts application search criteria into controlled form state.
 */
function formStateFromCriteria(criteria: ApplicationSearchCriteria): SearchFormState {
  return {
    companyName: criteria.companyName,
    source: criteria.source,
    status: criteria.statuses[0] ?? "",
    text: criteria.text
  };
}

/**
 * criteriaFromFormState
 *
 * Converts controlled form state into application search criteria.
 */
function criteriaFromFormState(values: SearchFormState): ApplicationSearchCriteria {
  return {
    companyName: values.companyName,
    source: values.source,
    statuses: values.status === "" ? [] : [values.status],
    text: values.text
  };
}

/**
 * ApplicationSearchForm
 *
 * Renders search and filter controls for the applications page.
 */
export function ApplicationSearchForm({ criteria, isSearching, onClear, onSearch }: ApplicationSearchFormProps) {
  const [values, setValues] = useState<SearchFormState>(formStateFromCriteria(criteria));

  /**
   * updateValue
   *
   * Updates one controlled search field.
   */
  function updateValue<Field extends keyof SearchFormState>(field: Field, value: SearchFormState[Field]) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  /**
   * handleSubmit
   *
   * Submits search criteria to the parent page.
   */
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    await onSearch(criteriaFromFormState(values));
  }

  /**
   * handleClear
   *
   * Resets search fields and asks the parent page to reload all applications.
   */
  async function handleClear() {
    const clearedValues: SearchFormState = {
      companyName: "",
      source: "",
      status: "",
      text: ""
    };

    setValues(clearedValues);
    await onClear();
  }

  return (
    <section className="form-card" aria-labelledby="application-search-heading">
      <h2 id="application-search-heading">Search applications</h2>

      <form className="search-form" onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="application-search-text">Text</label>
          <input
            id="application-search-text"
            value={values.text}
            onChange={(event) => updateValue("text", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="application-search-company">Company</label>
          <input
            id="application-search-company"
            value={values.companyName}
            onChange={(event) => updateValue("companyName", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="application-search-source">Source</label>
          <input
            id="application-search-source"
            value={values.source}
            onChange={(event) => updateValue("source", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="application-search-status">Status</label>
          <select
            id="application-search-status"
            value={values.status}
            onChange={(event) => updateValue("status", event.target.value as "" | ApplicationStatus)}
          >
            {statusOptions.map((status) => (
              <option key={status.value || "any"} value={status.value}>
                {status.label}
              </option>
            ))}
          </select>
        </div>

        <div className="form-actions form-actions--split">
          <button type="submit" disabled={isSearching}>
            {isSearching ? "Searching..." : "Search"}
          </button>
          <button type="button" className="secondary-button" disabled={isSearching} onClick={() => void handleClear()}>
            Clear search
          </button>
        </div>
      </form>
    </section>
  );
}