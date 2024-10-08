import React, {useState, useEffect} from "react";
import {Link, useNavigate, Navigate} from "react-router-dom";
import Scroll from "react-scroll";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faHeart,
	faPersonWalking,
	faRankingStar,
	faChartLine,
	faHeartbeat,
	faFilterCircleXmark,
	faArrowLeft,
	faUsers,
	faEye,
	faPencil,
	faTrashCan,
	faPlus,
	faGauge,
	faArrowRight,
	faTable,
	faArrowUpRightFromSquare,
	faRefresh,
	faFilter,
	faSearch,
	faPercent,
	faAdd,
	faFire,
	faMapMarkerAlt,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import {getLeaderboardListAPI} from "../../../../../API/Biometric";
import {
	topAlertMessageState,
	topAlertStatusState,
	currentUserState,
	dataPointFilterShowState,
	dataPointFilterTemporarySearchTextState,
	dataPointFilterActualSearchTextState,
	dataPointFilterSortState,
	dataPointFilterStatusState,
	dataPointFilterIsHeartRateState,
	dataPointFilterIsStepCountDeltaState,
} from "../../../../../AppState";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../../Reusable/FormInputFieldWithButton";
import {PAGE_SIZE_OPTIONS} from "../../../../../Constants/FieldOptions";
import {
	RANK_POINT_PERIOD_DAY,
	RANK_POINT_PERIOD_WEEK,
	RANK_POINT_PERIOD_MONTH,
	RANK_POINT_PERIOD_YEAR,
	RANK_POINT_FUNCTION_AVERAGE,
	RANK_POINT_FUNCTION_SUM,
} from "../../../../../Constants/App";
import MemberLeaderboardGlobalTabularListDesktop from "./TabularListDesktop";
import MemberLeaderboardGlobalTabularListMobile from "./TabularListMobile";
import {
	RANK_POINT_METRIC_TYPE_HEART_RATE,
	RANK_POINT_METRIC_TYPE_STEP_COUNTER,
	RANK_POINT_METRIC_TYPE_CALORIES_BURNED,
	RANK_POINT_METRIC_TYPE_DISTANCE_DELTA
} from "../../../../../Constants/App";
import Layout from "../../../../Menu/Layout";
import MobileLeaderboard from "../MobileView/MobileLeaderboard";

function MemberLeaderboardGlobalTabularList() {
	let navigate = useNavigate();

	////
	//// Global state.
	////

	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);
	const [currentUser] = useRecoilState(currentUserState);
	const [showFilter, setShowFilter] = useRecoilState(dataPointFilterShowState); // Filtering + Searching
	const [sort, setSort] = useRecoilState(dataPointFilterSortState); // Sorting
	const [temporarySearchText, setTemporarySearchText] = useRecoilState(
		dataPointFilterTemporarySearchTextState
	); // Searching - The search field value as your writes their query.
	const [actualSearchText, setActualSearchText] = useRecoilState(
		dataPointFilterActualSearchTextState
	); // Searching - The actual search query value to submit to the API.
	const [status, setStatus] = useRecoilState(dataPointFilterStatusState);
	// const [isHeartRate, setIsHeartRate] = useRecoilState(dataPointFilterIsHeartRateState);
	// const [isStepsCounter, setIsStepsCounter] = useRecoilState(dataPointFilterIsStepCountDeltaState);

	////
	//// Component states.
	////

	const [errors, setErrors] = useState({});
	const [listRank, setListRank] = useState("");
	const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] =
		useState("");
	const [isFetching, setFetching] = useState(false);
	const [pageSize, setPageSize] = useState(2); // Pagination
	const [previousCursors, setPreviousCursors] = useState([]); // Pagination
	const [nextCursor, setNextCursor] = useState(""); // Pagination
	const [currentCursor, setCurrentCursor] = useState(""); // Pagination
	const [isHeartRate, setIsHeartRate] = useState(true);
	const [isStepsCounter, setIsStepsCounter] = useState(false);
	const [isCaloriesBurned, setIsCaloriesBurned] = useState(false);
	const [isDistanceDelta, setIsDistanceDelta] = useState(false);
	const [period, setPeriod] = useState(RANK_POINT_PERIOD_DAY);
	const [calcFunction, setCalcFunction] = useState(RANK_POINT_FUNCTION_AVERAGE);

	////
	//// API.
	////

	function onRankPointistSuccess(response) {
		console.log("onRankPointistSuccess: Starting...");
		console.log("onRankPointistSuccess: Results:", response);
		if (response.results !== null) {
			setListRank(response);
			if (response.hasNextPage) {
				setNextCursor(response.nextCursor); // For pagination purposes.
			}
		} else {
			setListRank([]);
			setNextCursor("");
		}
	}

	function onRankPointistError(apiErr) {
		console.log("onRankPointistError: Starting...");
		setErrors(apiErr);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onRankPointistDone() {
		console.log("onRankPointistDone: Starting...");
		setFetching(false);
	}

	////
	//// BREADCRUMB
	////
	const breadcrumbItems = {
		items: [
			{text: "Dashboard", link: "/dashboard", isActive: false, icon: faGauge},
			{
				text: "Biometrics",
				link: "/biometrics",
				icon: faHeartbeat,
				isActive: false,
			},
			{text: "Leaderboard", link: "#", icon: faRankingStar, isActive: true},
		],
		mobileBackLinkItems: {
			link: "/biometrics",
			text: "Back to Biometrics",
			icon: faArrowLeft,
		},
	};

	////
	//// Event handling.
	////

	const handleNavigateToAccount = (e) => {
		e.preventDefault();
		navigate("/account", {state: {activeTabProp: "wearableTech"}});
	};

	const onHeartRateButtonClick = (e) => {
		e.preventDefault(); // Do not remove this line!
		setIsStepsCounter(false);
		setIsHeartRate(true);
		setIsCaloriesBurned(false);
		setIsDistanceDelta(false);
		setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
	};

	const onStepCounterButtonClick = (e) => {
		e.preventDefault(); // Do not remove this line!
		setIsStepsCounter(true);
		setIsHeartRate(false);
		setIsCaloriesBurned(false);
		setIsDistanceDelta(false);
		setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
	};

	const onCaloriesBurnedButtonClick = (e) => {
		e.preventDefault(); // Do not remove this line!
		setIsStepsCounter(false);
		setIsHeartRate(false);
		setIsCaloriesBurned(true);
		setIsDistanceDelta(false);
		setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
	};

	const onDistanceDeltaButtonClick = (e) => {
		e.preventDefault(); // Do not remove this line!
		setIsStepsCounter(false);
		setIsHeartRate(false);
		setIsCaloriesBurned(false);
		setIsDistanceDelta(true);
		setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
	};



	const fetchList = (
		user,
		cur,
		limit,
		keywords,
		stat,
		sbv,
		isHeartRate,
		isStepsCounter,
		isCaloriesBurned,
		isDistanceDelta,
		p,
		cf
	) => {
		setFetching(true);
		setErrors({});

		let metricType;
		if (isHeartRate) {
			metricType = RANK_POINT_METRIC_TYPE_HEART_RATE;
		} else if (isStepsCounter) {
			metricType = RANK_POINT_METRIC_TYPE_STEP_COUNTER;
		} else if (isCaloriesBurned) {
			metricType = RANK_POINT_METRIC_TYPE_CALORIES_BURNED;
		} else if (isDistanceDelta) {
			metricType = RANK_POINT_METRIC_TYPE_DISTANCE_DELTA;
		}

		let params = new Map();
		params.set("page_size", limit);
		params.set("page_size", "place");
		params.set("sort_order", "ASC");
		params.set("metric_data_type_names", metricType);
		params.set("period", p);
		params.set("function", cf);
		params.set("user_id", user.id);

		console.log("params:", params);

		// Make the submission to the API backend.
		getLeaderboardListAPI(
			params,
			onRankPointistSuccess,
			onRankPointistError,
			onRankPointistDone
		);
	};

	const onNextClicked = (e) => {
		console.log("onNextClicked");
		let arr = [...previousCursors];
		arr.push(currentCursor);
		setPreviousCursors(arr);
		setCurrentCursor(nextCursor);
	};

	const onPreviousClicked = (e) => {
		console.log("onPreviousClicked");
		let arr = [...previousCursors];
		const previousCursor = arr.pop();
		setPreviousCursors(arr);
		setCurrentCursor(previousCursor);
	};

	const onSearchButtonClick = (e) => {
		// Searching
		console.log("Search button clicked...");
		setActualSearchText(temporarySearchText);
	};

	const onSelectFitnessPlanForDeletion = (e, datum) => {
		console.log("onSelectFitnessPlanForDeletion", datum);
		setSelectedFitnessPlanForDeletion(datum);
	};

	const onDeselectFitnessPlanForDeletion = (e) => {
		console.log("onDeselectFitnessPlanForDeletion");
		setSelectedFitnessPlanForDeletion("");
	};

	// Function resets the filter state to its default state.
	const onClearFilterClick = (e) => {
		setShowFilter(false);
		setActualSearchText("");
		setTemporarySearchText("");
		setSort("timestamp,DESC");
		setStatus(0);
		setIsHeartRate(true);
		setIsStepsCounter(false);
		setIsCaloriesBurned(false);
		setIsDistanceDelta(false);
	};

	////
	//// Misc.
	////

	useEffect(() => {
		let mounted = true;

		if (mounted) {
			window.scrollTo(0, 0); // Start the page at the top of the page.
			fetchList(
				currentUser,
				currentCursor,
				pageSize,
				actualSearchText,
				status,
				sort,
				isHeartRate,
				isStepsCounter,
				isCaloriesBurned,
				isDistanceDelta,
				period,
				calcFunction
			);
		}

		return () => {
			mounted = false;
		};
	}, [
		currentUser,
		currentCursor,
		pageSize,
		actualSearchText,
		status,
		sort,
		isHeartRate,
		isStepsCounter,
		isCaloriesBurned,
		isDistanceDelta,
		period,
		calcFunction,
	]);

	////
	//// Component rendering.
	////

	console.log("--->", listRank);

	return (
		<Layout breadcrumbItems={breadcrumbItems}>
			<div className="box">
				<div className="columns">
					<div className="column">
						<h1 className="title is-4">
							<FontAwesomeIcon className="fas" icon={faRankingStar} />
							&nbsp;Leaderboard - Global&nbsp;<span class="tag is-success is-light">Beta</span>
						</h1>
					</div>
					<div className="column has-text-right">
						<button
							onClick={() =>
								fetchList(
									currentUser,
									currentCursor,
									pageSize,
									actualSearchText,
									status,
									sort,
									isHeartRate,
									isStepsCounter,
									isCaloriesBurned,
									isDistanceDelta,
									period,
									calcFunction
								)
							}
							class="is-fullwidth-mobile button is-link is-small"
							type="button">
							<FontAwesomeIcon className="mdi" icon={faRefresh} />
							&nbsp;
							<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
						</button>
						&nbsp;
						{/*
                      <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                          <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                      </button>

                      &nbsp;
                      <Link to={`/account/wearable-tech`} className="is-fullwidth-mobile button is-small is-success" type="button">
                          <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;Register Wearable
                      </Link>
                      */}
					</div>
				</div>

				{isFetching ? (
					<PageLoadingContent displayMessage={"Please wait..."} />
				) : (
					<>
						<FormErrorBox errors={errors} />

						{/* Section for selecting `metric type` */}
						<div class="column has-text-right">
							<button
								class={`button is-small ${isHeartRate && `is-info`}`}
								type="button"
								onClick={(e) => {
									onHeartRateButtonClick(e);
								}}>
								<FontAwesomeIcon className="mdi" icon={faHeart} />
								&nbsp;Heart Rate
							</button>
							<Link
								class={`button is-small ${isStepsCounter && `is-info`}`}
								type="button"
								onClick={(e) => {
									onStepCounterButtonClick(e);
								}}>
								<FontAwesomeIcon className="mdi" icon={faPersonWalking} />
								&nbsp;Steps Count
							</Link>
							&nbsp;
							<Link
								class={`button is-small ${isCaloriesBurned && `is-info`}`}
								type="button"
								onClick={(e) => {
									onCaloriesBurnedButtonClick(e);
								}}>
								<FontAwesomeIcon className="mdi" icon={faFire} />
								&nbsp;Calories Burned
							</Link>
							&nbsp;
							<Link
								class={`button is-small ${isDistanceDelta && `is-info`}`}
								type="button"
								onClick={(e) => {
									onDistanceDeltaButtonClick(e);
								}}>
								<FontAwesomeIcon className="mdi" icon={faMapMarkerAlt} />
								&nbsp;Distance
							</Link>
							&nbsp;
							{/*
                                DEVELOPERS NOTE:
                                - As we add more sensors, add your new sensors here...
                            */}
						</div>
						{/* Section for selecting `function` */}
						<div class="column has-text-right">
							{/*
                                DEVELOPERS NOTE:
                                - Some functions are not available for some metric types because it makes sense. Why would you keep a summation of heart rate? This
                                is the reason for code restrictions below.
                                - The only functions you should use are `RANK_POINT_FUNCTION_AVERAGE` and `RANK_POINT_FUNCTION_SUM`.
                          */}
							<button
								class={`button is-small ${
									calcFunction === RANK_POINT_FUNCTION_AVERAGE && `is-info`
								}`}
								type="button"
								onClick={(e) => {
									setCalcFunction(RANK_POINT_FUNCTION_AVERAGE);
								}}>
								<FontAwesomeIcon className="mdi" icon={faPercent} />
								&nbsp;Average
							</button>
							{!isHeartRate && (
								<Link
									class={`button is-small ${
										calcFunction === RANK_POINT_FUNCTION_SUM && `is-info`
									}`}
									type="button"
									onClick={(e) => {
										setCalcFunction(RANK_POINT_FUNCTION_SUM);
									}}>
									<FontAwesomeIcon className="mdi" icon={faAdd} />
									&nbsp;Sum
								</Link>
							)}
							&nbsp;
						</div>

						{/* Section for selecting `period` */}
						<div class="column has-text-right">
							{/*
                        DEVELOPERS NOTE:
                        - Period refers to the period of time the ranking is between. For example `day` would mean ranking for today.
                        - Week is ISO week, meaning it the week starts on Sunday and ends on Saturday.
                        - Month or year ranking are only ment for THIS month or year.
                    */}
							<button
								class={`button is-small ${
									period === RANK_POINT_PERIOD_DAY && `is-info`
								}`}
								type="button"
								onClick={(e) => {
									e.preventDefault();
									setPeriod(RANK_POINT_PERIOD_DAY);
								}}>
								Today
							</button>
							<Link
								class={`button is-small ${
									period === RANK_POINT_PERIOD_WEEK && `is-info`
								}`}
								type="button"
								onClick={(e) => {
									e.preventDefault();
									setPeriod(RANK_POINT_PERIOD_WEEK);
								}}>
								Week
							</Link>
							<Link
								class={`button is-small ${
									period === RANK_POINT_PERIOD_MONTH && `is-info`
								}`}
								type="button"
								onClick={(e) => {
									e.preventDefault();
									setPeriod(RANK_POINT_PERIOD_MONTH);
								}}>
								Month
							</Link>
							<Link
								class={`button is-small ${
									period === RANK_POINT_PERIOD_YEAR && `is-info`
								}`}
								type="button"
								onClick={(e) => {
									e.preventDefault();
									setPeriod(RANK_POINT_PERIOD_YEAR);
								}}>
								Year
							</Link>
							&nbsp;
						</div>

						{listRank &&
						listRank.results &&
						(listRank.results.length > 0 || previousCursors.length > 0) ? (
							<div>
								{/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
								<div class="is-hidden-touch">
									<MemberLeaderboardGlobalTabularListDesktop
										listRank={listRank}
										setPageSize={setPageSize}
										pageSize={pageSize}
										previousCursors={previousCursors}
										onPreviousClicked={onPreviousClicked}
										onNextClicked={onNextClicked}
										currentUser={currentUser}
									/>
								</div>

								{/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
								<div class="is-hidden-desktop">
									<MemberLeaderboardGlobalTabularListMobile
										listRank={listRank}
										setPageSize={setPageSize}
										pageSize={pageSize}
										previousCursors={previousCursors}
										onPreviousClicked={onPreviousClicked}
										onNextClicked={onNextClicked}
										currentUser={currentUser}
									/>
									{/* <MobileLeaderboard
									listRank={listRank}
									setPageSize={setPageSize}
									pageSize={pageSize}
									previousCursors={previousCursors}
									onPreviousClicked={onPreviousClicked}
									onNextClicked={onNextClicked}
									currentUser={currentUser}
									/> */}
								</div>
							</div>
						) : (
							<section className="hero is-medium has-background-white-ter">
								<div className="hero-body">
									<p className="title">
										<FontAwesomeIcon className="fas" icon={faTable} />
										&nbsp;No Biometrics
									</p>
									<p className="subtitle">
										You currently have no biometrics data.{" "}
										<b>
											<Link onClick={(e) => handleNavigateToAccount(e)}>
												Click here&nbsp;
												<FontAwesomeIcon className="mdi" icon={faArrowRight} />
											</Link>
										</b>{" "}
										to get started registering your wearable tech!
									</p>
								</div>
							</section>
						)}
					</>
				)}

				<div class="columns pt-5">
					<div class="column is-half">
						<Link class="button is-hidden-touch" to={`/biometrics`}>
							<FontAwesomeIcon className="fas" icon={faArrowLeft} />
							&nbsp;Back to Biometrics
						</Link>
						<Link
							class="button is-fullwidth is-hidden-desktop"
							to={`/biometrics`}>
							<FontAwesomeIcon className="fas" icon={faArrowLeft} />
							&nbsp;Back to Biometrics
						</Link>
					</div>
				</div>
			</div>
		</Layout>
	);
}

export default MemberLeaderboardGlobalTabularList;
