import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";


function AdminOfferListDesktop(props) {
  const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onSelectOfferForDeletion } = props;
  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Name</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                  <tr key={`desktop_${datum.id}`}>
                    <td data-label="Name">
                      {datum.name}
                    </td>
                    <td className="is-actions-cell">
                      <div className="buttons is-right">
                        {/*
                                                            <Link to={`/admin/offers/add?datum_id=${datum.id}&datum_name=${datum.name}`} className="button is-small is-success" type="button">
                                                                <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;CPS
                                                            </Link>
                                                        */}
                        <Link
                          to={`/admin/offer/${datum.id}`}
                          className="button is-small is-dark"
                          type="button"
                        >
                          <FontAwesomeIcon
                            className="mdi"
                            icon={faEye}
                          />
                          &nbsp;View
                        </Link>
                        <Link
                          to={`/admin/offer/${datum.id}/update`}
                          className="button is-small is-warning"
                          type="button"
                        >
                          <FontAwesomeIcon
                            className="mdi"
                            icon={faPencil}
                          />
                          &nbsp;Edit
                        </Link>
                        <button
                          onClick={(e, ses) =>
                            onSelectOfferForDeletion(e, datum)
                          }
                          className="button is-small is-danger"
                          type="button"
                        >
                          <FontAwesomeIcon
                            className="mdi"
                            icon={faTrashCan}
                          />
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
                onChange={(e) =>
                  setPageSize(parseInt(e.target.value))
                }
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
      </div>
    </div>
  );
}

export default AdminOfferListDesktop;
