import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCalendarMinus, faCalendarPlus, faTrophy, faCalendar, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faCalendarCheck, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, FITNESS_PLAN_STATUS_MAP } from "../../../../../Constants/FieldOptions";
import { RANK_POINT_METRIC_TYPE_HEART_RATE, RANK_POINT_METRIC_TYPE_STEP_COUNTER } from "../../../../../Constants/App";
import DateTimeTextFormatter from "../../../../Reusable/DateTimeTextFormatter";
import LeaderBoardTable from "../../../../Reusable/TableDesigns/Leaderboard/Table";


function MemberLeaderboardPersonalDesktop(props) {
    const { listRank, setPageSize, pageSize, previousCursors, onPreviousClicked, onNextClicked, currentUser } = props;

    const SevenDaysAvgHeader = () => (
        <div className="board_box">
            <div className="board_item">
                <div className="board_content">
                    <h5 className="mt-2 is-size-6 is-size-6-mobile is has-text-centered px-3 py-1 has-background-info has-text-white has has-text-weight-semibold mb-0">
                        7-days-Avg
                    </h5>
                </div>
            </div>
        </div>
    );

    const headers = [
        { title: 'RANK', className: 'is-vcentered' },
        { title: 'LEADERBOARD', className: 'is-vcentered' },
        { title: 'DAILY AVG', className: 'is-vcentered' },
        {
            component: <SevenDaysAvgHeader />,
            className: 'p-0 pb-2'
        },

    ];
    return (
        <LeaderBoardTable
            data={listRank}
            headers={headers}
        />
    );
}

export default MemberLeaderboardPersonalDesktop;
