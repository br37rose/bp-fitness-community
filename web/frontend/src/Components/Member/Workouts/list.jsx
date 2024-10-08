import {useState, useEffect} from "react";
import {Link} from "react-router-dom";
import Scroll from "react-scroll";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faArrowLeft,
	faPlus,
	faGauge,
	faArrowRight,
	faTable,
	faRefresh,
	faFilter,
	faSearch,
	faFilterCircleXmark,
	faDumbbell,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import {
	topAlertMessageState,
	topAlertStatusState,
	currentUserState,
	videoCollectionsFilterShowState,
	videoCollectionsFilterTemporarySearchTextState,
	videoCollectionsFilterActualSearchTextState,
	videoCollectionsFilterVideoTypeState,
	videoCollectionsFilterSortState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import {VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION} from "../../../Constants/FieldOptions";
import {getWorkoutListApi} from "../../../API/workout";
import MemberListMobile from "./listMobile";
import MemberListDesktop from "./listDesktop";

function MemberWorkoutList() {
	////
	//// Global state.
	////

	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);
	const [currentUser] = useRecoilState(currentUserState);
	const [showFilter, setShowFilter] = useRecoilState(
		videoCollectionsFilterShowState
	); // Filtering + Searching
	const [temporarySearchText, setTemporarySearchText] = useRecoilState(
		videoCollectionsFilterTemporarySearchTextState
	); // Searching - The search field value as your writes their query.
	const [actualSearchText, setActualSearchText] = useRecoilState(
		videoCollectionsFilterActualSearchTextState
	); // Searching - The actual search query value to submit to the API.
	const [status, setStatus] = useState(1);
	const [videoType, setVideoType] = useRecoilState(
		videoCollectionsFilterVideoTypeState
	);
	const [sort, setSort] = useRecoilState(videoCollectionsFilterSortState);

	////
	//// Component states.
	////

	const [errors, setErrors] = useState({});
	const [listData, setListData] = useState("");
	const [isFetching, setFetching] = useState(false);
	const [pageSize, setPageSize] = useState(10); // Pagination
	const [previousCursors, setPreviousCursors] = useState([]); // Pagination
	const [nextCursor, setNextCursor] = useState(""); // Pagination
	const [currentCursor, setCurrentCursor] = useState(""); // Pagination

	////
	//// API.
	////

	function onWorkoutListSuccess(response) {
		if (response.results !== null) {
			setListData(response);
			if (response.hasNextPage) {
				setNextCursor(response.nextCursor); // For pagination purposes.
			}
		} else {
			setListData([]);
			setNextCursor("");
		}
	}

	function onWorkoutlistError(apiErr) {
		setErrors(apiErr);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onWorkoutListDOne() {
		setFetching(false);
	}

	////
	//// Event handling.
	////

	const fetchList = (cur, limit, keywords, st, vt, sbv) => {
		setFetching(true);
		setErrors({});

		let params = new Map();
		params.set("page_size", limit); // Pagination
		params.set("user_id", currentUser.id);

		// DEVELOPERS NOTE: Our `sortByValue` is string with the sort field
		// and sort order combined with a comma seperation. Therefore we
		// need to split as follows.
		if (sbv !== undefined && sbv !== null && sbv !== "") {
			const sortArray = sbv.split(",");
			params.set("sort_field", sortArray[0]); // Sort (1 of 2)
			params.set("sort_order", sortArray[1]); // Sort (2 of 2)
		}

		if (cur !== "") {
			// Pagination
			params.set("cursor", cur);
		}

		params.set("type", vt);

		// Filtering
		if (keywords !== undefined && keywords !== null && keywords !== "") {
			// Searhcing
			params.set("search", keywords);
		}
		if (st !== undefined && st !== null && st !== "") {
			let stArray = new Array();
			stArray.push(st);
			params.set("status_list", stArray);
			if (st === 2) {
				params.set("exclude_archived", false);
			}
		}

		getWorkoutListApi(
			params,
			onWorkoutListSuccess,
			onWorkoutlistError,
			onWorkoutListDOne
		);
	};

	const onNextClicked = (e) => {
		let arr = [...previousCursors];
		arr.push(currentCursor);
		setPreviousCursors(arr);
		setCurrentCursor(nextCursor);
	};

	const onPreviousClicked = (e) => {
		let arr = [...previousCursors];
		const previousCursor = arr.pop();
		setPreviousCursors(arr);
		setCurrentCursor(previousCursor);
	};

	const onSearchButtonClick = (e) => {
		// Searching
		setActualSearchText(temporarySearchText);
	};

	// Function resets the filter state to its default state.
	const onClearFilterClick = (e) => {
		setShowFilter(false);
		setActualSearchText("");
		setTemporarySearchText("");
		setVideoType(0);
		setStatus(0);
		setSort("created,-1");
	};

	////
	//// Misc.
	////

	useEffect(() => {
		let mounted = true;

		if (mounted) {
			window.scrollTo(0, 0); // Start the page at the top of the page.
			fetchList(
				currentCursor,
				pageSize,
				actualSearchText,
				status,
				videoType,
				sort
			);
		}

		return () => {
			mounted = false;
		};
	}, [currentCursor, pageSize, actualSearchText, status, videoType, sort]);

	////
	//// Component rendering.
	////

	return (
		<>
			<div className="container is-fluid">
				<section className="section">
					{/* Desktop Breadcrumbs */}
					<nav className="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
						<ul>
							<li className="">
								<Link to="/dashboard" aria-current="page">
									<FontAwesomeIcon className="fas" icon={faGauge} />
									&nbsp;Dashboard
								</Link>
							</li>
							<li className="is-active">
								<Link aria-current="page">
									<FontAwesomeIcon className="fas" icon={faDumbbell} />
									&nbsp;Workouts
								</Link>
							</li>
						</ul>
					</nav>

					{/* Mobile Breadcrumbs */}
					<nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
						<ul>
							<li class="">
								<Link to="/dashboard" aria-current="page">
									<FontAwesomeIcon className="fas" icon={faArrowLeft} />
									&nbsp;Back to Dashboard
								</Link>
							</li>
						</ul>
					</nav>

					{/* Page */}
					<nav className="box">
						<div className="columns">
							<div className="column">
								<h1 className="title is-4">
									<FontAwesomeIcon className="fas" icon={faDumbbell} />
									&nbsp;Workouts
								</h1>
							</div>
							<div className="column has-text-right">
								<button
									onClick={() =>
										fetchList(
											currentCursor,
											pageSize,
											actualSearchText,
											status,
											videoType,
											sort
										)
									}
									class="is-fullwidth-mobile button is-link is-small"
									type="button">
									<FontAwesomeIcon className="mdi" icon={faRefresh} />
									&nbsp;
									<span class="is-hidden-desktop is-hidden-tablet">
										Refresh
									</span>
								</button>
								&nbsp;
								<button
									onClick={(e) => setShowFilter(!showFilter)}
									class="is-fullwidth-mobile button is-small is-primary"
									type="button">
									<FontAwesomeIcon className="mdi" icon={faFilter} />
									&nbsp;Filter
								</button>
								&nbsp;
								<Link
									to={`/workouts/add`}
									className="is-fullwidth-mobile button is-small is-success"
									type="button">
									<FontAwesomeIcon className="mdi" icon={faPlus} />
									&nbsp;New
								</Link>
							</div>
						</div>

						{/* FILTER */}
						{showFilter && (
							<div
								class="has-background-white-bis"
								style={{borderRadius: "15px", padding: "20px"}}>
								{/* Filter Title + Clear Button */}
								<div class="columns is-mobile">
									<div class="column is-half">
										<strong>
											<u>
												<FontAwesomeIcon className="mdi" icon={faFilter} />
												&nbsp;Filter
											</u>
										</strong>
									</div>
									<div class="column is-half has-text-right">
										<Link onClick={onClearFilterClick}>
											<FontAwesomeIcon
												className="mdi"
												icon={faFilterCircleXmark}
											/>
											&nbsp;Clear Filter
										</Link>
									</div>
								</div>

								{/* Filter Options */}
								<div class="columns">
									<div class="column">
										<FormInputFieldWithButton
											label={"Search"}
											name="temporarySearchText"
											type="text"
											placeholder="Search by name"
											value={temporarySearchText}
											helpText=""
											onChange={(e) => setTemporarySearchText(e.target.value)}
											isRequired={true}
											maxWidth="100%"
											buttonLabel={
												<>
													<FontAwesomeIcon className="fas" icon={faSearch} />
												</>
											}
											onButtonClick={onSearchButtonClick}
										/>
									</div>

									<div class="column">
										<FormSelectField
											label="Status"
											name="status"
											placeholder="Pick"
											selectedValue={status}
											errorText={errors && errors.status}
											helpText=""
											onChange={(e) => setStatus(parseInt(e.target.value))}
											options={
												VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION
											}
										/>
									</div>
								</div>
							</div>
						)}

						{isFetching ? (
							<PageLoadingContent displayMessage={"Please wait..."} />
						) : (
							<>
								<FormErrorBox errors={errors} />
								{listData &&
								listData.results &&
								(listData.results.length > 0 || previousCursors.length > 0) ? (
									<div>
										{/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
										<div class="is-hidden-touch">
											<MemberListDesktop
												listData={listData}
												setPageSize={setPageSize}
												pageSize={pageSize}
												previousCursors={previousCursors}
												onPreviousClicked={onPreviousClicked}
												onNextClicked={onNextClicked}
											/>
										</div>

										{/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
										<div class="is-fullwidth is-hidden-desktop">
											<MemberListMobile
												listData={listData}
												setPageSize={setPageSize}
												pageSize={pageSize}
												previousCursors={previousCursors}
												onPreviousClicked={onPreviousClicked}
												onNextClicked={onNextClicked}
											/>
										</div>
									</div>
								) : (
									<section className="hero is-medium has-background-white-ter">
										<div className="hero-body">
											<p className="title">
												<FontAwesomeIcon className="fas" icon={faTable} />
												&nbsp;No Workouts
											</p>
											<p className="subtitle">
												No class types.{" "}
												<b>
													<Link to="/workouts/add">
														Click here&nbsp;
														<FontAwesomeIcon
															className="mdi"
															icon={faArrowRight}
														/>
													</Link>
												</b>{" "}
												to get started creating your first workout.
											</p>
										</div>
									</section>
								)}
							</>
						)}

						<div class="columns pt-5">
							<div class="column is-half">
								<Link
									class="button is-hidden-touch"
									to={`/dashboard`}>
									<FontAwesomeIcon className="fas" icon={faArrowLeft} />
									&nbsp;Back to Dashboard
								</Link>
								<Link
									class="button is-fullwidth is-hidden-desktop"
									to={`/dashboard`}>
									<FontAwesomeIcon className="fas" icon={faArrowLeft} />
									&nbsp;Back to Dashboard
								</Link>
							</div>
							<div class="column is-half has-text-right">
								<Link
									class="button is-success is-hidden-touch"
									to={`/workouts/add`}>
									<FontAwesomeIcon className="fas" icon={faPlus} />
									&nbsp;New
								</Link>
								<Link
									class="button is-success is-fullwidth is-hidden-desktop"
									to={`/workouts/add`}>
									<FontAwesomeIcon className="fas" icon={faPlus} />
									&nbsp;New
								</Link>
							</div>
						</div>
					</nav>
				</section>
			</div>
		</>
	);
}

export default MemberWorkoutList;
