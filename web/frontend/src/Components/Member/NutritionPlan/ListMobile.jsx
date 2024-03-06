import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faLeaf, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../Constants/FieldOptions";
import DateTimeTextFormatter from "../../Reusable/DateTimeTextFormatter";

/*
Display for both tablet and mobile.
*/
function AdminNutritionPlanListMobile(props) {
    const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onSelectNutritionPlanForDeletion } = props;
    return (
        <>
            {listData && listData.results && listData.results.map(function (datum, i) {
                return <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
                    <strong>Name:</strong>&nbsp;{datum.name}
                    <br />
                    <br />
                    <strong>Status:</strong>&nbsp;{FITNESS_PLAN_STATUS_MAP[datum.status]}
                    <br />
                    <br />
                    <strong>Created At:</strong>&nbsp;<DateTimeTextFormatter value={datum.createdAt} />
                    <br />
                    <br />
                    <strong>Modified At:</strong>&nbsp;<DateTimeTextFormatter value={datum.modifiedAt} />
                    <br />
                    <br />

                    {/* Tablet only */}
                    <div class="is-hidden-mobile" key={`tablet_${datum.id}`}>
                        <div className="buttons is-right">
                            <Link to={`/nutrition-plan/${datum.id}`} class="button is-small is-dark" type="button">
                                <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;View
                            </Link>
                        </div>
                    </div>

                    {/* Mobile only */}
                    <div class="is-hidden-tablet" key={`mobile_${datum.id}`}>
                        <div class="columns is-mobile">
                            <div class="column">
                                <Link to={`/nutrition-plan/${datum.id}`} class="button is-small is-dark is-fullwidth" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;View
                                </Link>
                            </div>
                        </div>
                    </div>

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

export default AdminNutritionPlanListMobile;
