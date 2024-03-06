import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDownload, faArchive, faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";


function AdminMemberDetailForInvoiceListDesktop(props) {
    const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onDownloadClick } = props;
    return (
        <div class="b-table">
            <div class="table-wrapper has-mobile-cards">
                <table class="table is-fullwidth is-striped is-hoverable is-fullwidth">
                    <thead>
                        <tr>
                            <th>Invoice number</th>
                            <th>Amount</th>
                            <th>Status</th>
                            <th>Payment date</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>

                        {listData && listData.results && listData.results.map(function(datum, i){
                            return (
                              <tr key={`desktop_${datum.id}`}>
                                <td data-label="Invoice number">{datum.number}</td>
                                <td data-label="Amount">${datum.total}</td>
                                <td data-label="Status">
                                    {datum.paid
                                        ?
                                        <span className="tag is-success">Invoice paid</span>
                                        :
                                        <span className="tag is-danger">Invoice not paid</span>
                                    }
                                </td>
                                <td data-label="Payment date">{DateTime.fromISO(  new Date(datum.created * 1000).toISOString()  ).toLocaleString(DateTime.DATETIME_MED)}</td>

                                <td className="is-actions-cell">
                                  <div className="buttons is-right">
                                    <a
                                      target="_blank"
                                      rel="noreferrer"
                                      href={datum.hostedInvoiceUrl}
                                      className="button is-small is-primary"
                                      type="button"
                                    >
                                      <FontAwesomeIcon
                                        className="mdi"
                                        icon={faEye}
                                      />
                                      &nbsp;View&nbsp;
                                      <FontAwesomeIcon
                                        className="fas"
                                        icon={faArrowUpRightFromSquare}
                                      />
                                    </a>
                                    <Link
                                      onClick={(e,l)=>onDownloadClick(e, datum.invoicePdf)}
                                      className="button is-small is-success"
                                      type="button"
                                    >
                                      <FontAwesomeIcon
                                        className="mdi"
                                        icon={faDownload}
                                      />
                                      &nbsp;Download PDF
                                    </Link>

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

export default AdminMemberDetailForInvoiceListDesktop;
