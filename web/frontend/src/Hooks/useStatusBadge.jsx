import React from 'react';
import { STATUS_OPTIONS } from '../Constants/FieldOptions';

const useStatusBadge = ({ statusValue }) => {
    console.log(statusValue)
  const statusOption = STATUS_OPTIONS.find((option) => option.value === statusValue);

  const colorClass =
    statusOption &&
    (statusValue === 4 || statusValue === 3
      ? 'is-danger'
      : statusValue === 1
      ? 'is-warning'
      : statusValue === 5
      ? 'is-dark'
      : statusValue === 2
      ? 'is-success'
      : 'is-info');

  return statusOption ? (
    <span className={`tag ${colorClass}`}>{statusOption.label}</span>
  ) : (
    <span className="tag is-info">Unknown</span>
  );
};

export default useStatusBadge;
