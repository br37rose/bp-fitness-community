import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faTrophy, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../../Constants/FieldOptions";
import DateTimeTextFormatter from "../../../Reusable/DateTimeTextFormatter";


function MemberDataPointHistoricalTabularListDesktop(props) {
  const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, currentUser } = props;

  // Defensive Code
  if (currentUser) {
      if (currentUser.primaryHealthTrackingDevice === undefined || currentUser.primaryHealthTrackingDevice === null || currentUser.primaryHealthTrackingDevice === "") {
          console.log("currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId: EMPTY");
          return null;
      }
  }

  console.log("currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId:", currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId);
  console.log("currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId:", currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId);

  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Metric</th>
              <th>Value</th>
              <th>Timestamp</th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                  <tr key={`desktop_${datum.id}`}>
                    <td data-label="Metric">
                        {datum.metricId === currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId &&
                            <>Heart Rate</>
                        }
                        {datum.metricId === currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId &&
                            <>Steps Counter</>
                        }
                    </td>
                    <td data-label="Value">
                        {datum.value}

                        {/* Unit of measure */}
                        {datum.metricId === currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId  &&
                            <>&nbsp;bpm</>
                        }
                        {datum.metricId === currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId  &&
                            <>&nbsp;Steps</>
                        }
                    </td>
                    <td data-label="Timestamp">
                        <DateTimeTextFormatter value={datum.timestamp} />
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

export default MemberDataPointHistoricalTabularListDesktop;
