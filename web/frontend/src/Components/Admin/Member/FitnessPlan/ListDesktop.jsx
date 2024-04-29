import React from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEye, faTrashCan } from "@fortawesome/free-solid-svg-icons";

import {
  PAGE_SIZE_OPTIONS,
  FITNESS_PLAN_STATUS_MAP,
} from "../../../../Constants/FieldOptions";
import DateTimeTextFormatter from "../../../Reusable/DateTimeTextFormatter";

function AdminFitnessPlanListDesktop(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectFitnessPlanForDeletion,
  } = props;
  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Name</th>
              <th>Status</th>
              <th>Created At</th>
              <th>Modified At</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                  <tr key={`desktop_${datum.id}`}>
                    <td data-label="Name">{datum.name}</td>
                    <td data-label="Status">
                      {FITNESS_PLAN_STATUS_MAP[datum.status]}
                    </td>
                    <td data-label="Created At">
                      <DateTimeTextFormatter value={datum.createdAt} />
                    </td>
                    <td data-label="Modified At">
                      <DateTimeTextFormatter value={datum.modifiedAt} />
                    </td>
                    <td className="is-actions-cell">
                      <div className="buttons is-right">
                        {/*
                                                            <Link to={`/fitness-plans/add?datum_id=${datum.id}&datum_name=${datum.name}`} className="button is-small is-success" type="button">
                                                                <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;CPS
                                                            </Link>
                                                        */}
                        <Link
                          to={`/fitness-plan/${datum.id}`}
                          className="button is-small is-dark"
                          type="button"
                        >
                          <FontAwesomeIcon className="mdi" icon={faEye} />
                          &nbsp;View
                        </Link>
                        <button
                          onClick={(e, ses) =>
                            onSelectFitnessPlanForDeletion(e, datum)
                          }
                          className="button is-small is-danger"
                          type="button"
                        >
                          <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                          &nbsp;Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>

        <div class="columns">
          <div class="column is-half">
            <span class="select">
              <select
                class={`input has-text-grey-light`}
                name="pageSize"
                onChange={(e) => setPageSize(parseInt(e.target.value))}
              >
                {PAGE_SIZE_OPTIONS.map(function (option, i) {
                  return (
                    <option
                      selected={pageSize === option.value}
                      value={option.value}
                    >
                      {option.label}
                    </option>
                  );
                })}
              </select>
            </span>
          </div>
          <div class="column is-half has-text-right">
            {previousCursors.length > 0 && (
              <button class="button" onClick={onPreviousClicked}>
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
      </div>
    </div>
  );
}

export default AdminFitnessPlanListDesktop;
