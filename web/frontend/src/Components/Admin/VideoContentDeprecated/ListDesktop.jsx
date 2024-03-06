import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faVideo, faVimeo, faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import {
    PAGE_SIZE_OPTIONS,
    EXERCISE_CATEGORY_MAP,
    EXERCISE_MOMENT_TYPE_MAP,
    EXERCISE_TYPE_MAP,
    EXERCISE_STATUS_MAP
} from "../../../Constants/FieldOptions";
import {
    EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_VIDEO_TYPE_YOUTUBE,
    EXERCISE_VIDEO_TYPE_VIMEO
} from "../../../Constants/App";


function AdminVideoCollectionListDesktop(props) {
  const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onSelectMemberForDeletion } = props;
  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Name</th>
              <th>Gender</th>
              <th>Movement Type</th>
              <th>Category</th>
              <th>Type</th>
              <th>Video Type</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                    <tr key={`desktop_${datum.id}`}>
                    <td data-label="Name">{datum.name}{datum.alternateName && <span>&nbsp;{datum.alternateName}</span>}</td>
                    <td data-label="Gender"> {datum.gender}</td>
                    <td data-label="Movement Type">{EXERCISE_MOMENT_TYPE_MAP[datum.movementType]}</td>
                    <td data-label="Category">{EXERCISE_CATEGORY_MAP[datum.category]}</td>
                    <td data-label="Type">{EXERCISE_TYPE_MAP[datum.type]}</td>
                    <td data-label="Video Type">
                        {(() => {
                            switch (datum.videoType) {
                                case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER: return (
                                    <>
                                        S3
                                    </>
                                );
                                case EXERCISE_VIDEO_TYPE_YOUTUBE: return (
                                    <>
                                        YouTube
                                    </>
                                );
                                case EXERCISE_VIDEO_TYPE_VIMEO: return (
                                    <>
                                        Vimeo
                                    </>
                                );
                                default: return null;
                            }
                        })()}
                    </td>
                    <td data-label="Status">{EXERCISE_STATUS_MAP[datum.status]}</td>
                    <td className="is-actions-cell">
                      <div className="buttons is-right">
                        {/*
                                                            <Link to={`/admin/members/add?datum_id=${datum.id}&datum_name=${datum.name}`} className="button is-small is-success" type="button">
                                                                <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;CPS
                                                            </Link>
                                                        */}
                        <Link
                          to={`/admin/exercise/${datum.id}`}
                          className="button is-small is-dark"
                          type="button"
                        >
                          <FontAwesomeIcon
                            className="mdi"
                            icon={faEye}
                          />
                          &nbsp;View
                        </Link>
                        {/*
                        <Link
                          to={`/admin/member/${datum.id}/update`}
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
                            onSelectMemberForDeletion(e, datum)
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
                        */}
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

export default AdminVideoCollectionListDesktop;
