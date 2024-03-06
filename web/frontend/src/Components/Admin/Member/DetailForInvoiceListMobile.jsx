import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDownload, faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminMemberDetailForInvoiceListMobile(props) {
    const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onDownloadClick } = props;
    return (
        <>
            {listData && listData.results && listData.results.map(function(datum, i){
                return <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
                    <strong>Invoice number:</strong>&nbsp;
                    {datum.number}
                    <br />
                    <br />
                    <strong>Amount:</strong>&nbsp;
                    ${datum.total}
                    <br />
                    <br />
                    <strong>Status:</strong>&nbsp;
                    {datum.paid
                        ?
                        <span className="tag is-success">Invoice paid</span>
                        :
                        <span className="tag is-danger">Invoice not paid</span>
                    }
                    <br />
                    <br />

                    {/* Tablet only */}
                    <div class="is-hidden-mobile" key={`tablet_${datum.id}`}>
                        <div className="buttons is-right">
                            <Link to={datum.hostedInvoiceUrl} class="button is-small is-primary" type="button" target="_blank" rel="noreferrer">
                                <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                            </Link>
                            <Link onClick={(e,l)=>onDownloadClick(e, datum.invoicePdf)} class="button is-small is-warning" type="button">
                                <FontAwesomeIcon className="mdi" icon={faDownload} />&nbsp;Download PDF
                            </Link>
                        </div>
                    </div>

                    {/* Mobile only */}
                    <div class="is-hidden-tablet" key={`mobile_${datum.id}`}>
                        <div class="columns is-mobile">
                            <div class="column">
                                <Link to={datum.hostedInvoiceUrl} class="button is-small is-primary is-fullwidth" type="button" target="_blank" rel="noreferrer">
                                    <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;View&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                </Link>
                            </div>
                            <div class="column">
                                <Link onClick={(e,l)=>onDownloadClick(e, datum.invoicePdf)} class="button is-small is-warning is-fullwidth">
                                    <FontAwesomeIcon className="mdi" icon={faDownload} />&nbsp;Download PDF
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

export default AdminMemberDetailForInvoiceListMobile;
