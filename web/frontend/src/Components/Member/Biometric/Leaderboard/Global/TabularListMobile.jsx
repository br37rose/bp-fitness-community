import React, {useState, useEffect} from "react";
import {Link} from "react-router-dom";
import Scroll from "react-scroll";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faCalendarMinus,
	faCalendarPlus,
	faTrophy,
	faCalendar,
	faGauge,
	faSearch,
	faEye,
	faPencil,
	faTrashCan,
	faPlus,
	faArrowRight,
	faTable,
	faArrowUpRightFromSquare,
	faFilter,
	faRefresh,
	faCalendarCheck,
	faUsers,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import {
	PAGE_SIZE_OPTIONS,
	FITNESS_PLAN_STATUS_MAP,
} from "../../../../../Constants/FieldOptions";
import {
    RANK_POINT_METRIC_TYPE_HEART_RATE,
	RANK_POINT_METRIC_TYPE_STEP_COUNTER,
	RANK_POINT_METRIC_TYPE_CALORIES_BURNED,
	RANK_POINT_METRIC_TYPE_DISTANCE_DELTA
} from "../../../../../Constants/App";
import DateTimeTextFormatter from "../../../../Reusable/DateTimeTextFormatter";

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
		<main class="main_page">
			{/* <!-- Leaderboard --> */}
			<div class="leaderboard">
				<div class="leader_wrap">
					<Leader listRank={listRank} />
				</div>
			</div>

			<div class="players">
				<table>
					<tbody>
						{/* <!-- row --> */}
						{listRank &&
							listRank.results &&
							listRank.results.map(function (datum, i) {
								return (
									<tr key={`mobile_tablet_${datum.id}`}>
										<td>
											<div class="player_left">
												<div class="player_img">
													{datum.userAvatarObjectUrl ? (
														<figure class="figure-img is-128x128">
															<img src={datum.userAvatarObjectUrl} />
														</figure>
													) : (
														<figure class="figure-img is-128x128">
															<img src="/static/default_user.jpg" />
														</figure>
													)}
												</div>
												<div class="player_left_text">
													<h2 class="text_md f_500 pb_5">
														{datum.userFirstName}
													</h2>
													<p class="text_sm f_300">{`@${datum.userFirstName.toLowerCase()}`}</p>
												</div>
											</div>
										</td>
										<td>
											<div class="player_right">
												<h2 class="text_md f_500 pb_5">
													{datum.value.toFixed()}&nbsp;
													{datum.metricType ===
														RANK_POINT_METRIC_TYPE_HEART_RATE && <>bpm</>}
													{datum.metricType ===
														RANK_POINT_METRIC_TYPE_STEP_COUNTER && <>steps</>}
							                        {datum.metricDataTypeName === RANK_POINT_METRIC_TYPE_CALORIES_BURNED && <>kcal</>}
							                        {datum.metricDataTypeName === RANK_POINT_METRIC_TYPE_DISTANCE_DELTA && <>m</>}
												</h2>
												<img
													src="/static/leaderboard/arrow_up.svg"
													class="rank_arrow"
													alt="arrow_up"
												/>
											</div>
										</td>
									</tr>
								);
							})}
					</tbody>
				</table>
			</div>

			{/* );
            })} */}

			<div class="columns pt-4 is-mobile">
				<div class="column is-half">
					<span class="select">
						<select
							class={`input has-text-grey-light`}
							name="pageSize"
							onChange={(e) => setPageSize(parseInt(e.target.value))}>
							{PAGE_SIZE_OPTIONS.map(function (option, i) {
								return (
									<option
										selected={pageSize === option.value}
										value={option.value}>
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
		</main>
	);
}

export default MemberLeaderboardGlobalTabularListMobile;

const Leader = ({listRank}) => {
	const leaders = [];

	if (listRank && listRank.results) {
		listRank.results.forEach((datum) => {
			const leaderElement = (
				<div
					className={`column is-one-third leader ${
						datum.place === 1 ? "first" : datum.place === 2 ? "second" : "third"
					}`}>
					<div className="leader_img_wrap">
						{/* Leader Image */}
						{datum.userAvatarObjectUrl ? (
							<figure class="figure-img is-128x128">
								<img src={datum.userAvatarObjectUrl} />
							</figure>
						) : (
							<figure class="figure-img is-128x128">
								<img src="/static/default_user.jpg" />
							</figure>
						)}
						{/* Crown for first place */}
						{datum.place === 1 && (
							<img
								src="/static/leaderboard/crown.png"
								className="crown"
								alt="crown"
							/>
						)}
						{/* Badge */}
						<div className="badge_rotate">
							<h5 className="badge">{datum.place}</h5>
						</div>
					</div>
					<div className="leader_texts">
						<h2 className="text_md">{datum.userFirstName}</h2>
						<h4 className="text_lg">{datum.value.toFixed()}</h4>
						<p className="text_sm">@{datum.userFirstName.toLowerCase()}</p>
					</div>
				</div>
			);

			if (datum.place === 1) {
				leaders[1] = leaderElement;
			} else if (datum.place === 2) {
				leaders[0] = leaderElement;
			} else if (datum.place === 3) {
				leaders[2] = leaderElement;
			}
		});
	}

	return <div className="leader_wrap">{leaders}</div>;
};
