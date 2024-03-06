import React from "react";

import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";

function Pagination({
  previousCursors,
  listData,
  pageSize,
  setPageSize,
  onPreviousClicked,
  onNextClicked
}) {
  
  return (
    <div className="columns">
      <div className="column is-half">
        <span className="select">
          <select class={`input has-text-black`}
            name="pageSize"
            onChange={(e)=>setPageSize(parseInt(e.target.value))}
            value={pageSize}
          >
            {PAGE_SIZE_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </span>
      </div>
      <div className="column is-half has-text-right">
        {previousCursors.length > 0 && (
          <button
            class="button"
            onClick={onPreviousClicked}
          >
            Previous
          </button>
        )}
        {listData.hasNextPage && (
          <>
            <button class="button" onClick={onNextClicked}>
              Next
            </button>
          </>
        )}
      </div>
    </div>
  );
}

export default Pagination;
