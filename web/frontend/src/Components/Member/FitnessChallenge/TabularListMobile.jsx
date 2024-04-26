import React from "react";
import {
  RANK_POINT_METRIC_TYPE_HEART_RATE,
  RANK_POINT_METRIC_TYPE_STEP_COUNTER,
} from "../../../Constants/App";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function MemberLeaderboardGlobalTabularListMobile(props) {
  const {
    listRank,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    currentUser,
  } = props;
  return (
    <>
      {listRank &&
        listRank.results &&
        listRank.results.map(function (datum, i) {
          return (
            <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
              <strong>Place:</strong>&nbsp; #{datum.place}
              <br />
              <br />
              <strong>Picture:</strong>&nbsp;
              {datum.userAvatarObjectUrl ? (
                <figure class="image is-128x128">
                  <img src={datum.userAvatarObjectUrl} />
                </figure>
              ) : (
                <>None</>
              )}
              <br />
              <br />
              <strong>First Name:</strong>&nbsp;
              {datum.userFirstName}
              <br />
              <br />
              <strong>Value:</strong>&nbsp;
              {datum.value}&nbsp;
              {datum.metricType === RANK_POINT_METRIC_TYPE_HEART_RATE && (
                <>bpm</>
              )}
              {datum.metricType === RANK_POINT_METRIC_TYPE_STEP_COUNTER && (
                <>steps</>
              )}
              <br />
              <br />
            </div>
          );
        })}

      <div class="columns pt-4 is-mobile">
        <div class="column is-half">
          <span class="select">
            <select
              class={`input has-text-grey-light`}
              name="pageSize"
              onChange={(e) => setPageSize(parseInt(e.target.value))}
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
            <button class="button" onClick={onPreviousClicked}>
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
    </>
  );
}

export default MemberLeaderboardGlobalTabularListMobile;
