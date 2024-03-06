import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArchive, faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";


function AdminMemberDetailForWaitlisterListDesktop(props) {
    const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onSelectWorkoutSessionForDeletion } = props;
    return (
        <div class="b-table">
            <div class="table-wrapper has-mobile-cards">
                <table class="table is-fullwidth is-striped is-hoverable is-fullwidth">
                    <thead>
                        <tr>
                            <th>Class Type</th>
                            <th>Class</th>
                            <th>Trainer</th>
                            <th>Created At</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>

                        {listData && listData.results && listData.results.map(function(datum, i){
                            return <tr key={`desktop_${datum.id}`}>
                                <td data-label="Class Type">
                                    <Link to={`/admin/wp-type/${datum.workoutProgramTypeId}`} target="_blank" rel="noreferrer" class="">
                                        {datum.workoutProgramTypeName}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                    </Link>
                                </td>
                                <td data-label="Class">
                                    {datum.memberEmail
                                        ?
                                        <Link to={`/admin/branch/${datum.branchId}/class/${datum.workoutProgramId}`} target="_blank" rel="noreferrer" class="">
                                        {datum.workoutProgramName}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                        </Link>
                                        :
                                        <>-</>
                                    }
                                </td>
                                <td data-label="Trainer">
                                    <Link to={`/admin/branch/${datum.branchId}/trainer/${datum.trainerId}`} target="_blank" rel="noreferrer" class="">
                                        {datum.trainerName}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                    </Link>
                                </td>
                                <td data-label="Created At">{datum.createdAt}</td>
                                {/*

                                    <td class="is-actions-cell">
                                        <div class="buttons is-right">

                                            <Link to={`/admin/branch/${datum.branchId}/class/${datum.workoutProgramId}/session/${datum.workoutSessionId}/waitlister/${datum.id}`} class="button is-small is-primary" type="button" target="_blank" rel="noreferrer">
                                                <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                            </Link>
                                            <Link to={`/admin/branch/${datum.branchId}/class/${datum.workoutProgramId}/session/${datum.workoutSessionId}/waitlister/${datum.id}/update`} class="button is-small is-warning" type="button" target="_blank" rel="noreferrer">
                                                <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                            </Link>
                                            <button onClick={(e, ses) => onSelectWorkoutSessionForDeletion(e, datum)} class="button is-small is-danger" type="button">
                                                <FontAwesomeIcon className="mdi" icon={faArchive} />&nbsp;Archive
                                            </button>
                                        </div>
                                    </td>
                                */}
                            </tr>;
                        })}

                    </tbody>
                </table>

                <div class="columns">
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

            </div>
        </div>
    );
}

export default AdminMemberDetailForWaitlisterListDesktop;
