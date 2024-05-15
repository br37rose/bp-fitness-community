import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCrown } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../../Components/Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../Constants/FieldOptions";
import { RANK_POINT_METRIC_TYPE_HEART_RATE, RANK_POINT_METRIC_TYPE_STEP_COUNTER } from "../../../Constants/App";
import DateTimeTextFormatter from "../../../Components/Reusable/DateTimeTextFormatter";


function MemberLeaderboardGlobalTabularListDesktop(props) {
  const { listRank, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, currentUser } = props;
  return (
    <div className="container is-fluid">
      <table className="leaderboard-table is-hoverable is-fullwidth">
      <thead>
            <tr>
                <th>Name</th>
                <th>Picture</th>
                <th>Value</th>
            </tr>
          </thead>
        <tbody>
        {listRank &&
              listRank.results &&
              listRank.results.map(function (datum, index) {
                return (
                  <tr key={datum.place} className={`is-${index === 0 ? 'is-highlighted' : 'normal'}`}>
                  <td>{index === 0 ? <FontAwesomeIcon icon={faCrown} /> : null} {datum.userFirstName}</td>
                    <td data-label="Picture">
                        {datum.userAvatarObjectUrl
                            ?
                            <figure class="figure-img is-128x128">
                                <img src={datum.userAvatarObjectUrl} />
                            </figure>
                            :
                            <figure class="figure-img is-128x128">
                                <img src="/static/default_user.jpg" />
                            </figure>
                        }
                    </td>
                    <td data-label="Value">
                        {datum.value}&nbsp;
                        {datum.metricType === RANK_POINT_METRIC_TYPE_HEART_RATE && <>bpm</>}
                        {datum.metricType === RANK_POINT_METRIC_TYPE_STEP_COUNTER && <>steps</>}
                    </td>
                  </tr>
                );
          
})}
        </tbody>
      </table>
      <p className="has-text-centered has-text-light">Last Updated 10th Oct 2020</p>
    </div>
  );
}

export default MemberLeaderboardGlobalTabularListDesktop;
