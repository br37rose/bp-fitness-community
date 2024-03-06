import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrash, faClock, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faLeaf, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';

import { getNutritionPlanDetailAPI, deleteNutritionPlanAPI } from "../../../../API/NutritionPlan";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../Reusable/FormCountryField";
import FormRegionField from "../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../../AppState";
import DataDisplayRowText from "../../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../../Reusable/DataDisplayRowRadio";
import FormTextTagRow from "../../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../../Reusable/FormTextOptionRow";
import DataDisplayRowMultiSelect from "../../../Reusable/FormTextOptionRow";
import DataDisplayRowMultiSelectStatic from "../../../Reusable/DataDisplayRowMultiSelectStatic";
import DataDisplayRowSelectStatic from "../../../Reusable/DataDisplayRowSelectStatic";
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
   WORKOUT_PREFERENCE_OPTIONS
} from "../../../../Constants/FieldOptions";
import {
   FITNESS_GOAL_STATUS_QUEUED, FITNESS_GOAL_STATUS_ACTIVE,
   GENDER_OTHER, GENDER_MALE, GENDER_FEMALE
} from "../../../../Constants/App";


function MemberNutritionPlanDetail() {
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
    const [selectedAdminNutritionPlanForDeletion, setSelectedAdminNutritionPlanForDeletion] = useState(null);

    ////
    //// Event handling.
    ////

    const onDeleteConfirmButtonClick = () => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteNutritionPlanAPI(
            selectedAdminNutritionPlanForDeletion.id,
            onAdminNutritionPlanDeleteSuccess,
            onAdminNutritionPlanDeleteError,
            onAdminNutritionPlanDeleteDone
        );
        setSelectedAdminNutritionPlanForDeletion(null);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onNutritionPlanDetailSuccess(response){
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

    function onAdminNutritionPlanDeleteSuccess(response) {
        console.log("onAdminNutritionPlanDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("AdminNutrition plan deleted");
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

    function onAdminNutritionPlanDeleteError(apiErr) {
        console.log("onAdminNutritionPlanDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
        console.log(
            "onAdminNutritionPlanDeleteError: topAlertMessage, topAlertStatus:",
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

    function onAdminNutritionPlanDeleteDone() {
        console.log("onAdminNutritionPlanDeleteDone: Starting...");
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
        return <Navigate to={forceURL}  />
    }

    return (
        <>
            <div class="container">
                <section class="section">

                    {/* Desktop Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/nutrition-plans" aria-current="page"><FontAwesomeIcon className="fas" icon={faLeaf} />&nbsp;AdminNutrition Plans</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/nutrition-plans" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to AdminNutrition Plans</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    <nav>
                        {/* Delete modal */}
                        <div class={`modal ${selectedAdminNutritionPlanForDeletion !== null ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedAdminNutritionPlanForDeletion(null)}></button>
                                </header>
                                <section class="modal-card-body">
                                    You are about to delete this nutrition plan and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                                    <button class="button" onClick={(e, ses) => setSelectedAdminNutritionPlanForDeletion(null)}>Cancel</button>
                                </footer>
                            </div>
                        </div>
                    </nav>

                    {/* Page */}
                    <nav class="box">
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faLeaf} />&nbsp;AdminNutrition Plan</p>
                            </div>
                            {datum.status === FITNESS_GOAL_STATUS_ACTIVE && <div class="column has-text-right">
                                <Link to={`/nutrition-plan/${id}/update`} class="button is-warning is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit & Re-request
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedAdminNutritionPlanForDeletion(datum)}} class="button is-danger is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faTrash} />&nbsp;Delete
                                </Link>
                            </div>}
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
                                              &nbsp;AdminNutrition Plan Submitted
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

                                        <div class= "tabs is-medium is-size-7-mobile">
                                          <ul>
                                            <li class="is-active">
                                                <Link><strong>Detail</strong></Link>
                                            </li>
                                            <li>
                                                <Link to={`/nutrition-plan/${datum.id}/submission-form`}>Submission Form</Link>
                                            </li>
                                          </ul>
                                        </div>

                                        <p class="title is-6">META</p>
                                        <hr />

                                        <DataDisplayRowText
                                            label="Name"
                                            value={datum.name}
                                        />

                                        <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;DETAIL</p>
                                        <hr />

                                        <DataDisplayRowText
                                            label="Instructions"
                                            value={datum.instructions}
                                            type="text_with_linebreaks"
                                        />
                                      </>
                                    }

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/nutrition-plans`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to nutrition plans</Link>
                                        </div>
                                        <div class="column is-half has-text-right">
                                            {datum.status === FITNESS_GOAL_STATUS_ACTIVE &&
                                                <Link to={`/nutrition-plan/${id}/update`} class="button is-warning is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit & Re-request</Link>
                                            }
                                        </div>
                                    </div>

                                </div>}
                            </>
                        }
                    </nav>
                </section>
            </div>
        </>
    );
}

export default MemberNutritionPlanDetail;
