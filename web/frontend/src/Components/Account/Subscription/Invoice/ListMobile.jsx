import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faDumbbell, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import { DateTime } from "luxon";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS } from "../../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AccountInvoiceListMobile(props) {
  const { listData, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, onDownloadClick } = props;
  return (
    <>
      {listData && listData.results && listData.results.map(function (datum, i) {
        return <div class="pb-2 subtitle mb-4" key={`mobile_tablet_${datum.id}`}>
          <hr />
          <strong>Invoice number:</strong>&nbsp;{datum.number}
          <br />
          <br />
          <strong>Amount:</strong>&nbsp;${datum.total}
          <br />
          <br />
          <strong>Payment date:</strong>&nbsp;
          {DateTime.fromISO(new Date(datum.created * 1000).toISOString()).toLocaleString(DateTime.DATETIME_MED)}
          <br />
          <br />
          {/*
                        ######
                        TABLET
                        ######
                    */}
          <div class="is-hidden-mobile" key={`tablet_${datum.id}`}>
            <div className="buttons is-right">
              <a
                target="_blank"
                rel="noreferrer"
                href={datum.hostedInvoiceUrl}
                className="button is-small is-dark"
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
              <Link onClick={(e, l) => onDownloadClick(e, datum.invoicePdf)} className="button is-small is-success" type="button">
                <FontAwesomeIcon className="mdi" icon={faEye} />&nbsp;Download PDF
              </Link>

            </div>
          </div>

          {/*
                        ######
                        MOBILE
                        ######
                    */}
          <div class="is-hidden-tablet pb-4" key={`mobile_${datum.id}`}>
            <div className="buttons is-right">
              <a
                target="_blank"
                rel="noreferrer"
                href={datum.hostedInvoiceUrl}
                className="button is-small is-primary is-fullwidth"
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
                onClick={(e, l) => onDownloadClick(e, datum.invoicePdf)}
                className="button is-small is-success is-fullwidth"
                type="button"
              >
                <FontAwesomeIcon
                  className="mdi"
                  icon={faEye}
                />
                &nbsp;Download PDF
              </Link>

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

export default AccountInvoiceListMobile;
