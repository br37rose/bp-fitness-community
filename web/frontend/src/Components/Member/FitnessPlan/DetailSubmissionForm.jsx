import React, {useState, useEffect} from "react";
import {Link, Navigate} from "react-router-dom";
import Scroll from "react-scroll";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faClock,
	faRepeat,
	faTasks,
	faTachometer,
	faPlus,
	faArrowLeft,
	faCheckCircle,
	faUserCircle,
	faGauge,
	faPencil,
	faTrophy,
	faEye,
	faIdCard,
	faAddressBook,
	faContactCard,
	faChartPie,
	faCogs,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";
import {useParams} from "react-router-dom";

import {
	getFitnessPlanDetailAPI,
	deleteFitnessPlanAPI,
} from "../../../API/FitnessPlan";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {topAlertMessageState, topAlertStatusState} from "../../../AppState";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../Reusable/DataDisplayRowRadio";
import FormTextTagRow from "../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../Reusable/FormTextOptionRow";
import DataDisplayRowMultiSelect from "../../Reusable/FormTextOptionRow";
import DataDisplayRowMultiSelectStatic from "../../Reusable/DataDisplayRowMultiSelectStatic";
import DataDisplayRowSelectStatic from "../../Reusable/DataDisplayRowSelectStatic";
import {
	HOME_GYM_EQUIPMENT_OPTIONS,
	HOME_GYM_EQUIPMENT_MAP,
	FEET_WITH_EMPTY_OPTIONS,
	INCHES_WITH_EMPTY_OPTIONS,
	GENDER_WITH_EMPTY_OPTIONS,
	PHYSICAL_ACTIVITY_MAP,
	PHYSICAL_ACTIVITY_WITH_EMPTY_OPTIONS,
	WORKOUT_INTENSITY_WITH_EMPTY_OPTIONS,
	DAYS_PER_WEEK_MAP,
	DAYS_PER_WEEK_WITH_EMPTY_OPTIONS,
	TIME_PER_DAY_MAP,
	TIME_PER_DAY_WITH_EMPTY_OPTIONS,
	MAX_WEEK_MAP,
	MAX_WEEK_WITH_EMPTY_OPTIONS,
	FITNESS_GOAL_MAP,
	FITNESS_GOAL_OPTIONS,
	WORKOUT_PREFERENCE_MAP,
	WORKOUT_PREFERENCE_OPTIONS,
} from "../../../Constants/FieldOptions";
import {
	FITNESS_GOAL_STATUS_QUEUED,
	FITNESS_GOAL_STATUS_ACTIVE,
	GENDER_OTHER,
	GENDER_MALE,
	GENDER_FEMALE,
} from "../../../Constants/App";
import Layout from "../../Menu/Layout";

function MemberFitnessPlanSubmissionForm() {
	////
	//// URL Parameters.
	////

	const {id} = useParams();

	////
	//// Global state.
	////

	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);

	////
	//// Component states.
	////

	const [errors, setErrors] = useState({});
	const [isFetching, setFetching] = useState(false);
	const [forceURL, setForceURL] = useState("");
	const [datum, setDatum] = useState({});
	const [tabIndex, setTabIndex] = useState(1);
	const [selectedFitnessPlanForDeletion, setSelectedFitnessPlanForDeletion] =
		useState(null);

	////
	//// Event handling.
	////

	const onDeleteConfirmButtonClick = () => {
		console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

		deleteFitnessPlanAPI(
			selectedFitnessPlanForDeletion.id,
			onFitnessPlanDeleteSuccess,
			onFitnessPlanDeleteError,
			onFitnessPlanDeleteDone
		);
		setSelectedFitnessPlanForDeletion(null);
	};

	////
	//// API.
	////

	// --- Detail --- //

	function onFitnessPlanDetailSuccess(response) {
		console.log("onFitnessPlanDetailSuccess: Starting...");
		setDatum(response);
	}

	function onFitnessPlanDetailError(apiErr) {
		console.log("onFitnessPlanDetailError: Starting...");
		setErrors(apiErr);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onFitnessPlanDetailDone() {
		console.log("onFitnessPlanDetailDone: Starting...");
		setFetching(false);
	}

	// --- Delete --- //

	function onFitnessPlanDeleteSuccess(response) {
		console.log("onFitnessPlanDeleteSuccess: Starting..."); // For debugging purposes only.

		// Update notification.
		setTopAlertStatus("success");
		setTopAlertMessage("Fitness plan deleted");
		setTimeout(() => {
			console.log(
				"onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
				topAlertMessage,
				topAlertStatus
			);
			setTopAlertMessage("");
		}, 2000);

		// Redirect back to the video categories page.
		setForceURL("/fitness-plans");
	}

	function onFitnessPlanDeleteError(apiErr) {
		console.log("onFitnessPlanDeleteError: Starting..."); // For debugging purposes only.
		setErrors(apiErr);

		// Update notification.
		setTopAlertStatus("danger");
		setTopAlertMessage("Failed deleting");
		setTimeout(() => {
			console.log(
				"onFitnessPlanDeleteError: topAlertMessage, topAlertStatus:",
				topAlertMessage,
				topAlertStatus
			);
			setTopAlertMessage("");
		}, 2000);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onFitnessPlanDeleteDone() {
		console.log("onFitnessPlanDeleteDone: Starting...");
		setFetching(false);
	}

	////
	//// BREADCRUMB
	////
	const breadcrumbItems = {
		items: [
			{text: "Dashboard", link: "/dashboard", isActive: false, icon: faGauge},
			{
				text: "Fitness Plans",
				link: "/fitness-plans",
				icon: faTrophy,
				isActive: false,
			},
			{text: "Detail", link: "#", icon: faEye, isActive: true},
		],
		mobileBackLinkItems: {
			link: "/fitness-plans",
			text: "Back to Fitness Plans",
			icon: faArrowLeft,
		},
	};

	////
	//// Misc.
	////

	useEffect(() => {
		let mounted = true;

		if (mounted) {
			window.scrollTo(0, 0); // Start the page at the top of the page.

			setFetching(true);
			getFitnessPlanDetailAPI(
				id,
				onFitnessPlanDetailSuccess,
				onFitnessPlanDetailError,
				onFitnessPlanDetailDone
			);
		}

		return () => {
			mounted = false;
		};
	}, []);
	////
	//// Component rendering.
	////

	if (forceURL !== "") {
		return <Navigate to={forceURL} />;
	}

	return (
		<Layout breadcrumbItems={breadcrumbItems}>
			{/* Modal */}
			<nav>
				{/* Delete modal */}
				<div
					class={`modal ${
						selectedFitnessPlanForDeletion !== null ? "is-active" : ""
					}`}>
					<div class="modal-background"></div>
					<div class="modal-card">
						<header class="modal-card-head">
							<p class="modal-card-title">Are you sure?</p>
							<button
								class="delete"
								aria-label="close"
								onClick={(e, ses) =>
									setSelectedFitnessPlanForDeletion(null)
								}></button>
						</header>
						<section class="modal-card-body">
							You are about to delete this fitness plan and all the data
							associated with it. This action is cannot be undone. Are you sure
							you would like to continue?
						</section>
						<footer class="modal-card-foot">
							<button
								class="button is-success"
								onClick={onDeleteConfirmButtonClick}>
								Confirm
							</button>
							<button
								class="button"
								onClick={(e, ses) => setSelectedFitnessPlanForDeletion(null)}>
								Cancel
							</button>
						</footer>
					</div>
				</div>
			</nav>

			{/* Page */}
			<div class="box">
				{datum && (
					<div class="columns">
						<div class="column">
							<p class="title is-4">
								<FontAwesomeIcon className="fas" icon={faTrophy} />
								&nbsp;Fitness Plan
							</p>
						</div>
						<div class="column has-text-right"></div>
					</div>
				)}
				<FormErrorBox errors={errors} />

				{/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

				{isFetching ? (
					<PageLoadingContent displayMessage={"Please wait..."} />
				) : (
					<>
						{datum && (
							<div key={datum.id}>
								{/*
                                      ---------------------------------------------
                                      Queue Status GUI
                                      ---------------------------------------------
                                    */}
								{datum.status === FITNESS_GOAL_STATUS_QUEUED && (
									<>
										<section className="hero is-medium has-background-white-ter">
											<div className="hero-body">
												<p className="title">
													<FontAwesomeIcon className="fas" icon={faClock} />
													&nbsp;Fitness Plan Submitted
												</p>
												<p className="subtitle">
													You have successfully submitted this fitness plan to
													our team. The estimated time until our team completes
													your fitness plan will take about <b>1 or 2 days</b>.
													Please check back later.
												</p>
											</div>
										</section>
									</>
								)}

								{/*
                                      ---------------------------------------------
                                      Active Status GUI
                                      ---------------------------------------------
                                    */}
								{datum.status === FITNESS_GOAL_STATUS_ACTIVE && (
									<>
										{/* Tab navigation */}

										<div class="tabs is-medium is-size-7-mobile">
											<ul>
												<li>
													<Link to={`/fitness-plan/${datum.id}`}>Detail</Link>
												</li>
												<li class="is-active">
													<Link>
														<strong>Submission Form</strong>
													</Link>
												</li>
											</ul>
										</div>

										<p class="title is-6">META</p>
										<hr />

										<DataDisplayRowText label="Name" value={datum.name} />

										<p class="title is-6 pt-5">
											<FontAwesomeIcon className="fas" icon={faIdCard} />
											&nbsp;EQUIPMENT ACCESS
										</p>
										<hr />

										<DataDisplayRowRadio
											label="What equipment do you have access to"
											value={datum.equipmentAccess}
											opt1Value={1}
											opt1Label="No Equipment (calistanic/outdoor options)"
											opt2Value={2}
											opt2Label="Full Gym Access"
											opt3Value={3}
											opt3Label="Home Gym"
										/>

										<DataDisplayRowMultiSelectStatic
											label="Please select all the home gym equipment that you have (Optional)"
											selectedValues={datum.homeGymEquipment}
											map={HOME_GYM_EQUIPMENT_MAP}
										/>

										<DataDisplayRowRadio
											label="Do you workout at home?"
											value={datum.hasWorkoutsAtHome}
											opt1Value={1}
											opt1Label="Yes"
											opt2Value={2}
											opt2Label="No"
										/>

										<p class="title is-6 pt-5">
											<FontAwesomeIcon className="fas" icon={faIdCard} />
											&nbsp;PERSONAL DETAILS
										</p>
										<hr />

										<DataDisplayRowText
											label="Birthday"
											value={datum.birthday}
											type="date"
										/>

										<DataDisplayRowText
											label="Height"
											value={`${datum.heightFeet}\' ${datum.heightInches}"`}
										/>

										<DataDisplayRowText
											label="Weight"
											value={`${datum.weight} lbs`}
										/>

										<DataDisplayRowRadio
											label="Gender"
											value={datum.gender}
											opt1Value={GENDER_MALE}
											opt1Label="Male"
											opt2Value={GENDER_FEMALE}
											opt2Label="Female"
											opt3Value={GENDER_OTHER}
											opt3Label="Other"
										/>
										{datum.gender === GENDER_OTHER && (
											<DataDisplayRowText
												label="Gender (Other)"
												value={datum.genderOther}
											/>
										)}

										<p class="title is-6 pt-5">
											<FontAwesomeIcon className="fas" icon={faIdCard} />
											&nbsp;CURRENT PHYSICAL ACTIVITY
										</p>
										<hr />

										<DataDisplayRowText
											label="What is your ideal weight for your fitness goal?"
											value={`${datum.idealWeight} lbs`}
										/>

										<DataDisplayRowSelectStatic
											label="My current level of physical activity is"
											selectedValue={datum.physicalActivity}
											map={PHYSICAL_ACTIVITY_MAP}
										/>

										<p class="title is-6 pt-5">
											<FontAwesomeIcon className="fas" icon={faIdCard} />
											&nbsp;GOAL(S) FOR FITNESS PLAN
										</p>
										<hr />

										<DataDisplayRowSelectStatic
											label="Enter the number of days per week that you can train"
											selectedValue={datum.daysPerWeek}
											map={DAYS_PER_WEEK_MAP}
										/>

										<DataDisplayRowSelectStatic
											label="Enter the length of time per day that you can train"
											selectedValue={datum.timePerDay}
											map={TIME_PER_DAY_MAP}
										/>

										<DataDisplayRowSelectStatic
											label="Enter the number of weeks that you would like your training plan to last"
											selectedValue={datum.maxWeeks}
											map={MAX_WEEK_MAP}
										/>

										<DataDisplayRowMultiSelectStatic
											label="Enter your fitness goals"
											selectedValues={datum.goals}
											map={FITNESS_GOAL_MAP}
										/>

										<DataDisplayRowMultiSelectStatic
											label="Enter your workout preferences"
											selectedValues={datum.workoutPreferences}
											map={WORKOUT_PREFERENCE_MAP}
										/>
									</>
								)}

								<div class="columns pt-5">
									<div class="column is-half">
										<Link
											class="button is-hidden-touch"
											to={`/fitness-plans`}>
											<FontAwesomeIcon className="fas" icon={faArrowLeft} />
											&nbsp;Back to fitness plans
										</Link>
										<Link
											class="button is-fullwidth is-hidden-desktop"
											to={`/fitness-plans`}>
											<FontAwesomeIcon className="fas" icon={faArrowLeft} />
											&nbsp;Back to fitness plans
										</Link>
									</div>
									<div class="column is-half has-text-right"></div>
								</div>
							</div>
						)}
					</>
				)}
			</div>
		</Layout>
	);
}

export default MemberFitnessPlanSubmissionForm;
