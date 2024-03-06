import React from "react";
import PropTypes from "prop-types";

import Pagination from "./Pagination";
import useStatusBadge from "../../../Hooks/useStatusBadge";
import StringToDateConverter from "../StringToDateConverter";

const TableBody = ({
  listData,
  columns,
  emptyMessage,
  enablePagination,
  previousCursors,
  onNextClicked,
  onPreviousClicked,
  pageSize,
  setPageSize,
}) => {
  console.log(listData)
  const renderTableHeader = () => (
    <thead>
      <tr>
        {columns.map((column) => (
          <th key={column.key}>{column.title}</th>
        ))}
      </tr>
    </thead>
  );

  const formatDuration = (value) => {
    return value.toString().concat(" ", "minutes");
  };

  const renderTableCell = (item, column) => {
    const value = item[column.key];
    console.log(value)
    if (column.key === "durationInMinutes") {
      return <td key={column.key}>{formatDuration(value)}</td>;
    } else if (column.key === "actions") {
      return <td key={column.key}>{column.renderCell(item)}</td>;
    } else if (column.key === "status") {
      return (
        <td key={column.key}>
          <RenderStatusCell statusValue={value} />
        </td>
      );
    } else if (column.key === "paymentDate") {
      return (
        <td key={column.key}>
          <StringToDateConverter value={value} />
        </td>
      );
    } else {
      return <td key={column.key}>{value}</td>;
    }
  };

  const renderTableBody = () => {
    if (
      listData &&
      listData.results &&
      (listData.results.length > 0 || previousCursors.length > 0)
    ) {
      return (
        <tbody>
          {listData.results.map((item) => (
            <tr key={item.id}>
              {columns.map((column) => renderTableCell(item, column))}
            </tr>
          ))}
        </tbody>
      );
    } else {
      return (
        <tbody>
          <tr>
            <td colSpan={columns.length} className="has-text-centered">
              {emptyMessage}
            </td>
          </tr>
        </tbody>
      );
    }
  };

  return (
    <div className="table-container">
      <div className="b-table">
        <div className="table-wrapper has-mobile-cards">
          <table className="table is-fullwidth is-striped is-hoverable">
            {renderTableHeader()}
            {renderTableBody()}
          </table>
          {enablePagination && (
            <Pagination
              previousCursors={previousCursors}
              listData={listData}
              pageSize={pageSize}
              setPageSize={setPageSize}
              onPreviousClicked={onPreviousClicked}
              onNextClicked={onNextClicked}
            />
          )}
        </div>
      </div>
    </div>
  );
};

TableBody.propTypes = {
  listData: PropTypes.shape({
    results: PropTypes.arrayOf(PropTypes.object),
  }).isRequired,
  columns: PropTypes.arrayOf(
    PropTypes.shape({
      key: PropTypes.string.isRequired,
      title: PropTypes.string.isRequired,
      renderCell: PropTypes.func,
    })
  ).isRequired,
  emptyMessage: PropTypes.string.isRequired,
  tableTitle: PropTypes.string,
  enablePagination: PropTypes.bool,
};

TableBody.defaultProps = {
  enablePagination: false,
};

export default TableBody;

const RenderStatusCell = (statusValue) => {
  const statusBadge = useStatusBadge(statusValue);
  return statusBadge;
};
