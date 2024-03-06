import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faVideo, faVimeo, faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import {
    PAGE_SIZE_OPTIONS,
    VIDEO_COLLECTION_STATUS_MAP,
    VIDEO_CONTENT_VIDEO_TYPE_MAP,
} from "../../../../Constants/FieldOptions";
import {
    EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_VIDEO_TYPE_YOUTUBE,
    EXERCISE_VIDEO_TYPE_VIMEO
} from "../../../../Constants/App";


function MemberVideoContentListDesktop(props) {
  const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onSelectMemberForDeletion } = props;
  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Name</th>
              <th>Category</th>
              <th>Video Type</th>
              <th>Offer</th>
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
                    <td data-label="Category">
                        {datum.categoryName
                            ?
                                <>
                                    <Link to={`/video-category/${datum.categoryId}`} target="_blank" rel="noreferrer">
                                        {datum.categoryName}&nbsp;<FontAwesomeIcon className="mdi" icon={faArrowUpRightFromSquare} />
                                    </Link>
                                </>
                            :
                                <>-</>
                        }
                    </td>
                    <td data-label="Video Type">{VIDEO_CONTENT_VIDEO_TYPE_MAP[datum.type]}</td>
                    <td data-label="Offer">
                        {datum.offerName
                            ?
                                <>
                                    <Link to={`/offer/${datum.offerId}`} target="_blank" rel="noreferrer">
                                        {datum.offerName}&nbsp;<FontAwesomeIcon className="mdi" icon={faArrowUpRightFromSquare} />
                                    </Link>
                                </>
                            :
                                <>-</>
                        }
                    </td>
                    <td data-label="Status">{VIDEO_COLLECTION_STATUS_MAP[datum.status]}</td>
                    <td className="is-actions-cell">
                      <div className="buttons is-right">

                        <Link to={`/video-collection/${datum.collectionId}/video-content/${datum.id}`} className="button is-small is-dark" type="button">
                          <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;View
                        </Link>
                        {/*
                        <Link
                          to={`/member/${datum.id}/update`}
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

export default MemberVideoContentListDesktop;
