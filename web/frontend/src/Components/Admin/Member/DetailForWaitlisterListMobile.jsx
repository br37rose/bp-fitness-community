import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminMemberDetailForWaitlisterListMobile(props) {
    const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onSelectWorkoutSessionForDeletion } = props;
    return (
        <>
            {listData && listData.results && listData.results.map(function(datum, i){
                return <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
                    <strong>Class Type:</strong>&nbsp;
                    <Link to={`/admin/wp-type/${datum.workoutProgramTypeId}`} target="_blank" rel="noreferrer" class="">
                        {datum.workoutProgramTypeName}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                    </Link>
                    <br />
                    <br />
                    <strong>Class:</strong>&nbsp;
                    <Link to={`/admin/branch/${datum.branchId}/class/${datum.workoutProgramId}`} target="_blank" rel="noreferrer" class="">
                        {datum.workoutProgramName}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                    </Link>
                    <br />
                    <br />
                    <strong>Trainer:</strong>&nbsp;
                    <Link to={`/admin/branch/${datum.branchId}/trainer/${datum.trainerId}`} target="_blank" rel="noreferrer" class="">
                        {datum.trainerName}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                    </Link>
                    <br />
                    <br />
                    <strong>Status:</strong>&nbsp;
                    {datum.status === 2
                        ?
                        <span class="tag is-success is-light">Attending</span>
                        :
                        <span class="tag is-danger is-light">Cancelled</span>
                    }
                    <br />
                    <br />

                </div>;
            })}

            <div class="columns pt-4 is-mobile">
                <div class="column is-half">
                    <span class="select">
                        <select class={`input has-text-grey-light`}
                                 name="pageSize"
                             onChange={(e)=>setPageSize(parseInt(e.target.value))}>
                            {PAGE_SIZE_OPTIONS.map(function(option, i){
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

export default AdminMemberDetailForWaitlisterListMobile;
