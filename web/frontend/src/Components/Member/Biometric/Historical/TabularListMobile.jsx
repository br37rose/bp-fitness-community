import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faTrophy, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../../Constants/FieldOptions";
import { RANK_POINT_METRIC_TYPE_HEART_RATE, RANK_POINT_METRIC_TYPE_STEP_COUNTER } from "../../../../Constants/App";
import DateTimeTextFormatter from "../../../Reusable/DateTimeTextFormatter";
import GoogleFitDataPointValueToString from "../../../Reusable/SpecificPage/Biometrics/GoogleFitDataPointValueToString";
import GoogleFitDataPointTypeToString from "../../../Reusable/SpecificPage/Biometrics/GoogleFitDataPointTypeToString";

/*
Display for both tablet and mobile.
*/
function MemberHistoricalDataGlobalTabularListMobile(props) {
    const { listRank, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, currentUser } = props;
    return (
        <>
            {listRank && listRank.results && listRank.results.map(function (datum, i) {
                return (
                  <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
                  <strong>Type:</strong>&nbsp;<GoogleFitDataPointTypeToString datapoint={datum} />
                  <br />
                  <strong>Value:</strong>&nbsp;<GoogleFitDataPointValueToString datapoint={datum} />
                  <br />
                  <strong>Start At:</strong>&nbsp;{datum.startAt}
                  <br />
                  <strong>End At:</strong>&nbsp;{datum.endAt}
                  <br />
                  <br />
                  </div>
                );
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
                    {listRank.hasNextPage && <>
                        <button class="button" onClick={onNextClicked}>Next</button>
                    </>}
                </div>
            </div>
        </>
    );
}

export default MemberHistoricalDataGlobalTabularListMobile;
