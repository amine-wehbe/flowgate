// Filter bar — method dropdown, status prefix input, URL keyword search
function FilterBar({ filters, onChange }) {
  return (
    <div className="filter-bar">
      <select
        value={filters.method}
        onChange={(e) => onChange({ ...filters, method: e.target.value })}
      >
        <option value="">All Methods</option>
        {["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"].map((m) => (
          <option key={m} value={m}>{m}</option>
        ))}
      </select>
      <input
        type="text"
        placeholder="Status (e.g. 2, 404)"
        value={filters.status}
        onChange={(e) => onChange({ ...filters, status: e.target.value })}
      />
      <input
        type="text"
        placeholder="Search URL..."
        value={filters.search}
        onChange={(e) => onChange({ ...filters, search: e.target.value })}
      />
      {(filters.method || filters.status || filters.search) && (
        <button onClick={() => onChange({ method: "", status: "", search: "" })}>
          Clear
        </button>
      )}
    </div>
  );
}

export default FilterBar;
