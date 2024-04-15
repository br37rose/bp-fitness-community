import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faTrophy, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../../../Constants/FieldOptions";
import { RANK_POINT_METRIC_TYPE_HEART_RATE, RANK_POINT_METRIC_TYPE_STEP_COUNTER } from "../../../../../Constants/App";
import DateTimeTextFormatter from "../../../../Reusable/DateTimeTextFormatter";


function MemberLeaderboardGlobalTabularListDesktop(props) {
  const { listRank, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, currentUser } = props;
  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
                <th>Place</th>
                <th>Picture</th>
                <th>First Name</th>
                <th>Value</th>
            </tr>
          </thead>
          <tbody>
            {listRank &&
              listRank.results &&
              listRank.results.map(function (datum, i) {
                return (
                  <tr key={`desktop_${datum.id}`}>
                    <td data-label="Timestamp">
                        #{datum.place}
                    </td>
                    <td data-label="Picture">
                        {datum.userAvatarObjectUrl
                            ?
                            <figure class="image is-128x128">
                                <img src={datum.userAvatarObjectUrl} />
                            </figure>
                            :
                            <>None</>
                        }
                    </td>
                    <td data-label="First Name">
                        {datum.userFirstName}
                    </td>
                    <td data-label="Value">
                        {datum.value}&nbsp;
                        {datum.metricType === RANK_POINT_METRIC_TYPE_HEART_RATE && <>bpm</>}
                        {datum.metricType === RANK_POINT_METRIC_TYPE_STEP_COUNTER && <>steps</>}
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
            {listRank.hasNextPage && (
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

export default MemberLeaderboardGlobalTabularListDesktop;
