import React from "react";

function formatDate(dateString, includeTime = true) {
  if (!dateString) {
    return null;
  }

  const date = new Date(dateString);

  if (isNaN(date)) {
    return null;
  }

  const options = {
    year: "numeric",
    month: "short",
    day: "numeric",
  };

  if (includeTime) {
    options.hour = "numeric";
    options.minute = "numeric";
    options.timeZoneName = "short";
  }

  return date.toLocaleString("en-US", options);
}

function StringToDateConverter({ value, includeTime = true }) {
  const formattedDate = formatDate(value, includeTime);

  if (!formattedDate) {
    return <span style={{ color: "red" }}>Invalid Date</span>;
  }

  return <>{formattedDate}</>;
}

export default StringToDateConverter;
