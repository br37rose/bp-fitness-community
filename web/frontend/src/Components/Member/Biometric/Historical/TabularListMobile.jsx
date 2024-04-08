import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faTrophy, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../../Constants/FieldOptions";
import DateTimeTextFormatter from "../../../Reusable/DateTimeTextFormatter";

/*
Display for both tablet and mobile.
*/
function MemberDataPointHistoricalTabularListMobile(props) {
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
        <>
            {listData && listData.results && listData.results.map(function (datum, i) {
                return <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
                    <strong>Metric:</strong>&nbsp;{datum.name}
                    <br />
                    <br />
                    <strong>Value:</strong>&nbsp;
                    {datum.value}

                    {/* Unit of measure */}
                    {datum.metricId === currentUser.primaryHealthTrackingDevice.heartRateBpmMetricId &&
                        <>&nbsp;bpm</>
                    }
                    {datum.metricId === currentUser.primaryHealthTrackingDevice.stepCountDeltaMetricId &&
                        <>&nbsp;Steps</>
                    }
                    <br />
                    <br />
                    <strong>Timestamp:</strong>&nbsp;
                    <DateTimeTextFormatter value={datum.timestamp} />
                    <br />
                    <br />

                </div>;
            })}

            <div class="columns pt-4 is-mobile">
                <div class="column is-half">
                    <span class="select">
                        <select class={`input has-text-grey-light`}
                            name="pageSize"
                            onChange={(e) => setPageSize(parseInt(e.target.value))}>
                            {PAGE_SIZE_OPTIONS.map(function (option, i) {
                                return <option selected={pageSize === option.value} value={option.value}>{option.label}</option>;
                            })}
                        </select>
                    </span>

                </div>
                <div class="column is-half has-text-right">
                    {previousCursors.length > 0 &&
                        <button class="button" onClick={onPreviousClicked}>Previous</button>
                    }
                    {listData.hasNextPage && <>
                        <button class="button" onClick={onNextClicked}>Next</button>
                    </>}
                </div>
            </div>
        </>
    );
}

export default MemberDataPointHistoricalTabularListMobile;
