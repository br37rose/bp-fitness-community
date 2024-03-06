import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faClock, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faLeaf, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';

import { getNutritionPlanDetailAPI, deleteNutritionPlanAPI } from "../../../API/NutritionPlan";
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
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
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
    MEALS_PER_DAY_MAP,
    CONSUME_FREQUENCY_MAP,
    NUTRITIONAL_GOAL_MAP
} from "../../../Constants/FieldOptions";
import {
    FITNESS_GOAL_STATUS_QUEUED, FITNESS_GOAL_STATUS_ACTIVE,
    GENDER_OTHER, GENDER_MALE, GENDER_FEMALE,
    WORKOUT_INTENSITY_LOW, WORKOUT_INTENSITY_MEDIUM, WORKOUT_INTENSITY_HIGH,
} from "../../../Constants/App";
import Layout from "../../Menu/Layout";


function MemberNutritionPlanSubmissionForm() {
    ////
    //// URL Parameters.
    ////

    const { id } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [datum, setDatum] = useState({});
    const [tabIndex, setTabIndex] = useState(1);
    const [selectedNutritionPlanForDeletion, setSelectedNutritionPlanForDeletion] = useState(null);

    ////
    //// Event handling.
    ////

    const onDeleteConfirmButtonClick = () => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteNutritionPlanAPI(
            selectedNutritionPlanForDeletion.id,
            onNutritionPlanDeleteSuccess,
            onNutritionPlanDeleteError,
            onNutritionPlanDeleteDone
        );
        setSelectedNutritionPlanForDeletion(null);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onNutritionPlanDetailSuccess(response) {
        console.log("onNutritionPlanDetailSuccess: Starting...");
        setDatum(response);
    }

    function onNutritionPlanDetailError(apiErr) {
        console.log("onNutritionPlanDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onNutritionPlanDetailDone() {
        console.log("onNutritionPlanDetailDone: Starting...");
        setFetching(false);
    }

    // --- Delete --- //

    function onNutritionPlanDeleteSuccess(response) {
        console.log("onNutritionPlanDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("Nutrition plan deleted");
        setTimeout(() => {
            console.log(
                "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
                topAlertMessage,
                topAlertStatus
            );
            setTopAlertMessage("");
        }, 2000);

        // Redirect back to the video categories page.
        setForceURL("/nutrition-plans");
    }

    function onNutritionPlanDeleteError(apiErr) {
        console.log("onNutritionPlanDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
            console.log(
                "onNutritionPlanDeleteError: topAlertMessage, topAlertStatus:",
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

    function onNutritionPlanDeleteDone() {
        console.log("onNutritionPlanDeleteDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            setFetching(true);
            getNutritionPlanDetailAPI(
                id,
                onNutritionPlanDetailSuccess,
                onNutritionPlanDetailError,
                onNutritionPlanDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Nutrition Plans', link: '#', icon: faLeaf, isActive: false },
            { text: 'Detail', link: '#', icon: faEye, isActive: true }
        ],
        mobileBackLinkItems: {
            link: '/nutrition-plans',
            text: 'Back to Nutrition Plans',
            icon: faArrowLeft
        }
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            {/* Mobile Breadcrumbs */}
            <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                <ul>
                    <li class=""><Link to="/nutrition-plans" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Nutrition Plans</Link></li>
                </ul>
            </nav>

            {/* Modal */}
            <nav>
                {/* Delete modal */}
                <div class={`modal ${selectedNutritionPlanForDeletion !== null ? 'is-active' : ''}`}>
                    <div class="modal-background"></div>
                    <div class="modal-card">
                        <header class="modal-card-head">
                            <p class="modal-card-title">Are you sure?</p>
                            <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedNutritionPlanForDeletion(null)}></button>
                        </header>
                        <section class="modal-card-body">
                            You are about to delete this nutrition plan and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                        </section>
                        <footer class="modal-card-foot">
                            <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                            <button class="button" onClick={(e, ses) => setSelectedNutritionPlanForDeletion(null)}>Cancel</button>
                        </footer>
                    </div>
                </div>
            </nav>

            {/* Page */}
            <div class="box">
                {datum && <div class="columns">
                    <div class="column">
                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faLeaf} />&nbsp;Nutrition Plan</p>
                    </div>
                    <div class="column has-text-right">

                    </div>
                </div>}
                <FormErrorBox errors={errors} />

                {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

                {isFetching
                    ?
                    <PageLoadingContent displayMessage={"Please wait..."} />
                    :
                    <>
                        {datum && <div class="container" key={datum.id}>
                            {/*
                                      ---------------------------------------------
                                      Queue Status GUI
                                      ---------------------------------------------
                                    */}
                            {datum.status === FITNESS_GOAL_STATUS_QUEUED
                                && <>
                                    <section className="hero is-medium has-background-white-ter">
                                        <div className="hero-body">
                                            <p className="title">
                                                <FontAwesomeIcon className="fas" icon={faClock} />
                                                &nbsp;Nutrition Plan Submitted
                                            </p>
                                            <p className="subtitle">
                                                You have successfully submitted this nutrition plan to our team. The estimated time until our team completes your nutrition plan will take about <b>1 or 2 days</b>. Please check back later.
                                            </p>
                                        </div>
                                    </section>
                                </>
                            }

                            {/*
                                      ---------------------------------------------
                                      Active Status GUI
                                      ---------------------------------------------
                                    */}
                            {datum.status === FITNESS_GOAL_STATUS_ACTIVE
                                && <>
                                    {/* Tab navigation */}

                                    <div class="tabs is-medium is-size-7-mobile">
                                        <ul>
                                            <li>
                                                <Link to={`/nutrition-plan/${datum.id}`}>Detail</Link>
                                            </li>
                                            <li class="is-active">
                                                <Link><strong>Submission Form</strong></Link>
                                            </li>
                                        </ul>
                                    </div>

                                    <p class="title is-6">META</p>
                                    <hr />

                                    <DataDisplayRowText
                                        label="Name"
                                        value={datum.name}
                                    />


                                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;PERSONAL DETAILS</p>
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
                                    {datum.gender === GENDER_OTHER &&
                                        <DataDisplayRowText
                                            label="Gender (Other)"
                                            value={datum.genderOther}
                                        />
                                    }

                                    <DataDisplayRowText
                                        label="What is your ideal weight for your nutrition goal?"
                                        value={`${datum.idealWeight} lbs`}
                                    />

                                    <DataDisplayRowSelectStatic
                                        label="My current level of physical activity is"
                                        selectedValue={datum.physicalActivity}
                                        map={PHYSICAL_ACTIVITY_MAP}
                                    />

                                    <DataDisplayRowRadio
                                        label="My current intensity in my exercise routine is"
                                        value={datum.workoutIntensity}
                                        opt1Value={WORKOUT_INTENSITY_LOW}
                                        opt1Label="Low"
                                        opt2Value={WORKOUT_INTENSITY_MEDIUM}
                                        opt2Label="Medium"
                                        opt3Value={WORKOUT_INTENSITY_HIGH}
                                        opt3Label="High"
                                    />

                                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;ALLERGIES</p>
                                    <hr />

                                    <DataDisplayRowRadio
                                        label="My current intensity in my exercise routine is"
                                        value={datum.hasAllergies}
                                        opt1Value={1}
                                        opt1Label="Yes"
                                        opt2Value={2}
                                        opt2Label="No"
                                    />

                                    {datum.hasAllergies === 1 && <DataDisplayRowText
                                        label="If yes, what are your allergies?"
                                        value={`${datum.allergies}`}
                                    />}

                                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;GOAL(S) FOR NUTRITION PLAN</p>
                                    <hr />

                                    <DataDisplayRowSelectStatic
                                        label="How many meals do you typically eat in a day?"
                                        selectedValue={datum.mealsPerDay}
                                        map={MEALS_PER_DAY_MAP}
                                    />

                                    <DataDisplayRowSelectStatic
                                        label="How often do you consume fast food or junk food?"
                                        selectedValue={datum.consumeJunkFood}
                                        map={CONSUME_FREQUENCY_MAP}
                                    />

                                    <DataDisplayRowSelectStatic
                                        label="How often do you consume fruits and/or vegetables?"
                                        selectedValue={datum.consumeFruitsAndVegetables}
                                        map={CONSUME_FREQUENCY_MAP}
                                    />

                                    <DataDisplayRowMultiSelectStatic
                                        label="Enter your nutrition goals"
                                        selectedValues={datum.goals}
                                        map={NUTRITIONAL_GOAL_MAP}
                                    />

                                    <DataDisplayRowSelectStatic
                                        label="Enter the number of weeks that you would like your training plan to last"
                                        selectedValue={datum.maxWeeks}
                                        map={MAX_WEEK_MAP}
                                    />

                                    <DataDisplayRowRadio
                                        label="Has Intermittent Fasting"
                                        value={datum.hasIntermittentFasting}
                                        opt1Value={1}
                                        opt1Label="Yes"
                                        opt2Value={2}
                                        opt2Label="No"
                                    />
                                </>
                            }

                            <div class="columns pt-5">
                                <div class="column is-half">
                                    <Link class="button is-fullwidth-mobile" to={`/nutrition-plans`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to nutrition plans</Link>
                                </div>
                                <div class="column is-half has-text-right">

                                </div>
                            </div>

                        </div>}
                    </>
                }
            </div>
        </Layout>
    );
}

export default MemberNutritionPlanSubmissionForm;
