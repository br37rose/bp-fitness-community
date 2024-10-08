import React, {useState, useEffect, useRef} from "react";
import {Link} from "react-router-dom";
import Scroll from "react-scroll";
import moment from "moment";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faHeart,
	faPersonWalking,
	faRankingStar,
	faChartLine,
	faHeartbeat,
	faRoad,
	faFire,
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
	faUser,
	faShoePrints,
	faMapMarkerAlt,
	faLock,
	faShoppingBag,
	faStar,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import {getMySummaryAPI} from "../../../../API/Biometric";
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
} from "../../../../AppState";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import BarChart from "../../../Reusable/Charts/Bar";
import Layout from "../../../Menu/Layout";

function MemberSummary() {
	////
	//// Global state.
	////

	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);
	const [currentUser] = useRecoilState(currentUserState);

	////
	//// Component states.
	////

	const [errors, setErrors] = useState({});
	const [datum, setDatum] = useState("");
	const [isFetching, setFetching] = useState(false);
	const [isComingSoon, setComingSoon] = useState(true);
	const [showBiometrics, setShowBiometrics] = useState(true); // Toggle state for biometrics

	////
	//// API.
	////

	function onSummarySuccess(response) {
		console.log("onSummarySuccess: Starting...");
		console.log("onSummarySuccess: response:", response);
		setDatum(response);
	}

	function onSummaryError(apiErr) {
		console.log("onSummaryError: Starting...");
		setErrors(apiErr);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onSummaryDone() {
		console.log("onSummaryDone: Starting...");
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
			{text: "My Summary", link: "#", icon: faRankingStar, isActive: true},
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

	const getDatum = (user) => {
		if (user !== undefined && user !== null && user !== "") {
			if (user.primaryHealthTrackingDeviceType !== 0) {
				setFetching(true);
				setErrors({});

				let params = new Map();
				params.set("user_id", user.id);

				getMySummaryAPI(
					params,
					onSummarySuccess,
					onSummaryError,
					onSummaryDone
				);
			} else {
				console.log("user does not have a device, prevented pulling data.");
			}
		}
	};

	const toggleBiometrics = () => setShowBiometrics(!showBiometrics);

	////
	//// Misc.
	////

	useEffect(() => {
		let mounted = true;

		if (mounted) {
			window.scrollTo(0, 0); // Start the page at the top of the page.
			getDatum(currentUser);
		}

		return () => {
			mounted = false;
		};
	}, [currentUser]);

	////
	//// Component rendering.
	////
	const formatToHours = (end) => {
		return moment(end).format("ha"); // Formats the end to hour with am/pm
	};

	const formatToDays = (end) => moment(end).format("Do MMM");

	const formatToWeeks = (end) => `Week ${moment(end).isoWeek()}`;

	const formatToMonths = (end) => moment(end).format("MMM");

	const transformData = (barData, label, text, timeframe, mode) => {
		let formatFunction;

		// Determine the appropriate format function based on the timeframe
		switch (timeframe) {
			case "hours":
				formatFunction = (item) => formatToHours(item.end); // Format to hours
				break;

			case "week":
				formatFunction = (item) => moment(item.end).format("ddd"); // Weekday
				break;
			case "month":
				formatFunction = (item) => moment(item.end).format("MMM D"); // Month and Date
				break;
			case "year":
				formatFunction = (item) => moment(item.end).format("MMM YYYY"); // Month and Year
				break;
			default:
				formatFunction = (item) => formatToHours(item.end); // Default to hours
				break;
		}

		// Determine the chart type and additional properties based on the mode
		let dataset;
		if (mode === 1) {
			// DEVELOPERS NOTE: ???
			dataset = {
				label: ` ${label}`,
				data: barData.map((item) => item.count),
				borderColor: "#E1BD67",
				backgroundColor: "#ffffff",
				type: "line",
				order: 1,
			};
		} else if (mode === 2) {
			// DEVELOPERS NOTE: ???
			dataset = {
				label: ` ${label}`,
				data: barData.map((item) => item.count),
				borderWidth: 0,
				barThickness: 1,
				backgroundColor: ["#5374DF"],
				order: 2,
				borderColor: "rgba(54, 162, 235, 1)",
				borderWidth: 0,
				borderRadius: 5,
				borderSkipped: false,
			};
		} else if (mode === 3) {
			// DEVELOPERS NOTE: `mode=3` is to get `averages` from API results.
			dataset = {
				label: ` ${label}`,
				data: barData.map((item) => item.average),
				borderColor: "#E1BD67",
				backgroundColor: "#ffffff",
				type: "line",
				order: 1,
			};
		} else if (mode === 4) {
			// DEVELOPERS NOTE: `mode=4` is to get `sum` from API results.
			dataset = {
				label: ` ${label}`,
				data: barData.map((item) => item.sum),
				borderColor: "#E1BD67",
				backgroundColor: "#ffffff",
				type: "line",
				order: 1,
			};
		}

		const transformedData = {
			text: text,
			labels: barData.map(formatFunction),
			datasets: [dataset],
		};

		return transformedData;
	};

	try {
		return (
			<Layout breadcrumbItems={breadcrumbItems}>
				<div className="box">
					<div className="columns">
						<div className="column">
							<h1 className="title is-4">
								<FontAwesomeIcon className="fas" icon={faStar} />
								&nbsp;My Summary&nbsp;<span class="tag is-success is-light">Beta</span>
							</h1>
						</div>
						<div className="column has-text-right">
							<button
								onClick={() => getDatum(currentUser)}
								class="is-fullwidth-mobile button is-link is-small"
								type="button">
								<FontAwesomeIcon className="mdi" icon={faRefresh} />
								&nbsp;
								<span class="is-hidden-desktop is-hidden-tablet">Refresh</span>
							</button>
							&nbsp;
							{/* Toggle switch for flipping metrics */}
							{/* <label className="switch is-rounded">
                                        <input
                                            type="checkbox"
                                            checked={showBiometrics}
                                            onChange={toggleBiometrics}
                                        />
                                        <span className="slider round"></span>
                                    </label> */}
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
							{datum !== undefined && datum !== null && datum !== "" ? (
								// Checking if 'datum' is not null or empty
								<>
									{showBiometrics ? (
										// Section 1: Biometrics content
										<div className="section">
											{/* Biometric Summary Row for Heart Rate + Step Count */}
											<div className="columns">
												{/* Heart Rate (Today) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faHeart}
																	/>
																</span>
															</div>
															<div
																title="Heart Rate average for today"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Heart Rate(Today)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.heartRateThisDaySummary &&
																	Math.round(
																		datum.heartRateThisDaySummary.average
																	)
																} bpm`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Heart Rate (Today) */}
												{/* Heart Rate (Week) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faHeart}
																	/>
																</span>
															</div>
															<div
																title="Heart Rate average in the last 7 days"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Heart Rate(Week)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.heartRateThisIsoWeekSummary &&
																	Math.round(
																		datum.heartRateThisIsoWeekSummary.average
																	)
																} bpm`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Heart Rate */}
												{/* Steps (Today) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faShoePrints}
																	/>
																</span>
															</div>
															<div
																title="Steps Count average for today"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Steps Count(Today)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.stepCountDeltaThisDaySummary &&
																	Math.round(
																		datum.stepCountDeltaThisDaySummary.average
																	)
																} counts`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Steps (Today) */}
												{/* Steps (Week) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faShoePrints}
																	/>
																</span>
															</div>
															<div
																title="Steps Count average in the last 7 days"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Steps Count(Week)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.stepCountDeltaThisIsoWeekSummary &&
																	Math.round(
																		datum.stepCountDeltaThisIsoWeekSummary
																			.average
																	)
																} counts`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Steps (Week) */}
											</div>{" "}
											{/* end Biometric Summary Row for Heart Rate + Step Count */}
											{/* Biometric Summary Row for Calories Burned + Distance Delta */}
											<div className="columns">
												{/* Calories Burned (Today) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faFire}
																	/>
																</span>
															</div>
															<div
																title="Heart Rate average for today"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Calories Burned (Today)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.caloriesBurnedThisDaySummary &&
																	Math.round(
																		datum.caloriesBurnedThisDaySummary.average
																	)
																} kcal`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Calories Burned (Today) */}
												{/* Calories Burned (Week) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faFire}
																	/>
																</span>
															</div>
															<div
																title="Heart Rate average in the last 7 days"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Calories Burned (Week)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.caloriesBurnedThisIsoWeekSummary &&
																	Math.round(
																		datum.caloriesBurnedThisIsoWeekSummary
																			.average
																	)
																} kcal`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Calories Burned */}
												{/* Distance Delta (Today) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faRoad}
																	/>
																</span>
															</div>
															<div
																title="Distance average for today"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Distance (Today)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.distanceDeltaThisDaySummary &&
																	Math.round(
																		datum.distanceDeltaThisDaySummary.average
																	)
																} m`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Distance (Today) */}
												{/* Distance (Week) */}
												<div className="column is-one-fourth">
													<div className="box">
														<div className="media">
															<div className="media-left">
																<span className="icon">
																	<FontAwesomeIcon
																		className="fas"
																		icon={faRoad}
																	/>
																</span>
															</div>
															<div
																title="Distance average in the last 7 days"
																className="media-content ">
																<p className="title has-text-weight-semibold is-7">
																	Distance (Week)
																</p>
																<p className="has-text-weight-semibold subtitle is-6">{`${
																	datum &&
																	datum.distanceDeltaThisIsoWeekSummary &&
																	Math.round(
																		datum.distanceDeltaThisIsoWeekSummary
																			.average
																	)
																} m`}</p>
															</div>
														</div>
													</div>
												</div>
												{/* end Distance (Week) */}
											</div>{" "}
											{/* end Biometric Summary Row for Calories Burned + Distance Delta */}
											{/* Biometric Graphs Rows */}
											<div className="columns is-multiline">
												{/* ------ ROW 1 ------ */}
												{/* Heart Rate (Today) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.heartRateThisDayData,
															"Heart Rate",
															"Heart Rate - Today",
															"hours",
															3
														)}
													/>
												</div>
												{/* end Heart Rate (Today) */}
												{/* Heart Rate (Week) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.heartRateThisIsoWeekData,
															"Heart Rate",
															"Heart Rate - Week",
															"week",
															3
														)}
													/>
												</div>
												{/* end Heart Rate (Week) */}
												{/* Heart Rate (Month) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.heartRateThisMonthData,
															"Heart Rate",
															"Heart Rate - Month",
															"month",
															3
														)}
													/>
												</div>
												{/* end Distance Delta (Month) */}
												{/* ------ ROW 2 ------ */}
												{/* Steps Count Delta (Today) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.stepCountDeltaThisDayData,
															"Steps Count",
															"Steps Count - Today",
															"hours",
															4
														)}
													/>
												</div>
												{/* end Steps Count Delta (Today) */}
												{/* Steps Count Delta (Week) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.stepCountDeltaThisIsoWeekData,
															"Steps Count",
															"Steps Count - Week",
															"week",
															4
														)}
													/>
												</div>
												{/* end Steps Count Delta (Week) */}
												{/* Steps Count Delta (Month) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.stepCountDeltaThisMonthData,
															"Steps Count",
															"Steps Count - Month",
															"month",
															4
														)}
													/>
												</div>
												{/* end Steps Count Delta (Month) */}
												{/* ------ ROW 3 ------ */}
												{/* Calories Burned (Today) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.caloriesBurnedThisDayData,
															"Calories Burnede",
															"Calories Burned - Today",
															"hours",
															4
														)}
													/>
												</div>
												{/* end Calories Burned (Today) */}
												{/* Calories Burned (Week) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.caloriesBurnedThisIsoWeekData,
															"Calories Burned",
															"Calories Burned - Week",
															"week",
															4
														)}
													/>
												</div>
												{/* end Calories Burned (Week) */}
												{/* Calories Burned (Month) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.caloriesBurnedThisMonthData,
															"Calories Burned",
															"Calories Burned - Month",
															"month",
															4
														)}
													/>
												</div>
												{/* end Calories Burned (Month) */}
												{/* ------ ROW 4 ------ */}
												{/* Distance Delta (Today) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.heartRateThisDayData,
															"Distance",
															"Distance - Today",
															"hours",
															4
														)}
													/>
												</div>
												{/* end Distance Delta (Today) */}
												{/* Distance Delta (Week) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.heartRateThisIsoWeekData,
															"Distance",
															"Distance - Week",
															"week",
															4
														)}
													/>
												</div>
												{/* end Distance Delta (Week) */}
												{/* Distance Delta (Month) */}
												<div className="column is-one-third">
													<BarChart
														data={transformData(
															datum.heartRateThisMonthData,
															"Distance",
															"Distance - Month",
															"month",
															4
														)}
													/>
												</div>
												{/* end Distance Delta (Month) */}
											</div>
											{/* end Biometric Graphs Rows */}
											{/* ??? */}
											<div className="columns">
												{/* <div className="column is-half">
                                                                <DoughnutChart data={barChartData} options={{ maintainAspectRatio: false }} />
                                                            </div> */}
												{/* <div className="column">
                            <div className="box">
                                <h3 className="title is-4">Top Selling Products</h3>
                                <table className="table is-fullwidth is-striped">
                                    <thead>
                                        <tr>
                                            <th>Product Name</th>
                                            <th>Price</th>
                                            <th>Rating</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {sampleTopSellingData.map((product, index) => (
                                            <tr key={index}>
                                                <td>{product.name}</td>
                                                <td>{formatCurrency(product.price)}</td>
                                                <td>{'★'.repeat(product.rating)}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div> */}
											</div>
											{/* end ??? */}
										</div>
									) : (
										// Section 2: Alternative content
										<div className="alternative-content">
											{/* Your alternative content */}
										</div>
									)}
								</>
							) : (
								<section className="hero is-medium has-background-white-ter">
									<div className="hero-body">
										<p className="title">
											<FontAwesomeIcon className="fas" icon={faTable} />
											&nbsp;No Biometrics
										</p>
										<p className="subtitle">
											You currently have no biometrics data. Please check back
											later!
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
						<div class="column is-half has-text-right">
							{/*
                <Link to={`/account/wearable-tech`} class="button is-success is-fullwidth-mobile">
                    <FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;Register Wearable
                </Link>
                */}
						</div>
					</div>
				</div>
			</Layout>
		);
	} catch (e) {
		console.log("render error:", e);
		return <>GUI FAILED TO LOAD - PLEASE CONTACT SYSTEM ADMINISTRATOR</>;
	}
}

export default MemberSummary;
