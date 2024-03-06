import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrash, faVideo, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faDumbbell, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';
import Vimeo from '@u-wave/react-vimeo';

import { getExerciseDetailAPI, deleteExerciseAPI } from "../../../API/Exercise";
import FormErrorBox from "../../Reusable/FormErrorBox";
import DataDisplayRowOffer from "../../Reusable/DataDisplayRowOffer";
import DataDisplayRowCheckbox from "../../Reusable/DataDisplayRowCheckbox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
    EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_TYPE_WITH_EMPTY_OPTIONS,
    EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    TIMED_LOCK_DURATION_WITH_EMPTY_OPTIONS
} from "../../../Constants/FieldOptions";
import {
    EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_VIDEO_TYPE_YOUTUBE,
    EXERCISE_VIDEO_TYPE_VIMEO,
    EXERCISE_TYPE_SYSTEM
} from "../../../Constants/App";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../Reusable/DataDisplayRowRadio";
import DataDisplayRowSelect from "../../Reusable/DataDisplayRowSelect";
import DataDisplayRowImage from "../../Reusable/DataDisplayRowImage";


function AdminExerciseDetail() {
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
    const [selectedExerciseForDeletion, setSelectedExerciseForDeletion] = useState(null);

    ////
    //// Event handling.
    ////

    const onDeleteConfirmButtonClick = () => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteExerciseAPI(
            selectedExerciseForDeletion.id,
            onExerciseDeleteSuccess,
            onExerciseDeleteError,
            onExerciseDeleteDone
        );
        setSelectedExerciseForDeletion(null);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onExerciseDetailSuccess(response){
        console.log("onExerciseDetailSuccess: Starting...");
        setDatum(response);
    }

    function onExerciseDetailError(apiErr) {
        console.log("onExerciseDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onExerciseDetailDone() {
        console.log("onExerciseDetailDone: Starting...");
        setFetching(false);
    }

    // --- Delete --- //

    function onExerciseDeleteSuccess(response) {
        console.log("onExerciseDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("Exercise deleted");
        setTimeout(() => {
        console.log(
            "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
            topAlertMessage,
            topAlertStatus
        );
        setTopAlertMessage("");
        }, 2000);

        // Redirect back to the members page.
        setForceURL("/admin/exercises");
    }

    function onExerciseDeleteError(apiErr) {
        console.log("onExerciseDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
        console.log(
            "onExerciseDeleteError: topAlertMessage, topAlertStatus:",
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

    function onExerciseDeleteDone() {
        console.log("onExerciseDeleteDone: Starting...");
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
            getExerciseDetailAPI(
                id,
                onExerciseDetailSuccess,
                onExerciseDetailError,
                onExerciseDetailDone
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
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/exercises" aria-current="page"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;Exercises</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/exercises" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Exercises</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    <nav>
                        {/* Delete modal */}
                        <div class={`modal ${selectedExerciseForDeletion !== null ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedExerciseForDeletion(null)}></button>
                                </header>
                                <section class="modal-card-body">
                                    You are about to delete this exercise and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                                    <button class="button" onClick={(e, ses) => setSelectedExerciseForDeletion(null)}>Cancel</button>
                                </footer>
                            </div>
                        </div>
                    </nav>

                    {/* Page */}
                    <nav class="box">
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;Exercise</p>
                            </div>
                            <div class="column has-text-right">
                                <Link to={`/admin/exercise/${id}/update`} class="button is-small is-warning is-fullwidth-mobile" type="button" disabled={datum.type === EXERCISE_TYPE_SYSTEM}>
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedExerciseForDeletion(datum)}} class="button is-small is-danger is-fullwidth-mobile" type="button" disabled={datum.type === EXERCISE_TYPE_SYSTEM}>
                                    <FontAwesomeIcon className="mdi" icon={faTrash} />&nbsp;Delete
                                </Link>
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
                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faVideo} />&nbsp;Video</p>
                                    <hr />

                                    {(() => {
                                        switch (datum.videoType) {
                                            case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER: return (
                                                <>
                                                    <video style={{width:'100%', height:'100%'}} controls>
                                                        <source src={datum.videoObjectUrl}
                                                            type="video/mp4"
                                                        />
                                                    </video>
                                                </>
                                            );
                                            case EXERCISE_VIDEO_TYPE_YOUTUBE: return (
                                                <>yt</>
                                            );
                                            case EXERCISE_VIDEO_TYPE_VIMEO: return (
                                                <Vimeo
                                                  video={`${datum.videoUrl}`}
                                                  autoplay
                                                />
                                            );
                                            default: return null;
                                        }
                                    })()}

                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Information</p>
                                    <hr />

                                    <DataDisplayRowImage label="Thumbnail" src={datum.thumbnailObjectUrl} alt="Thumbnail" />                                            

                                    <DataDisplayRowText
                                        label="Name"
                                        value={datum.name}
                                    />

                                    <DataDisplayRowText
                                        label="Alternate Name"
                                        value={datum.alternateName}
                                    />

                                    <DataDisplayRowText
                                        label="Description"
                                        value={datum.description}
                                    />

                                    <DataDisplayRowRadio
                                        label="Gender"
                                        value={datum.gender}
                                        opt1Value="Male"
                                        opt1Label="Male"
                                        opt2Value="Female"
                                        opt2Label="Female"
                                        opt3Value="Other"
                                        opt3Label="Other"
                                    />

                                    <DataDisplayRowSelect
                                        label="Movement Type"
                                        selectedValue={datum.movementType}
                                        options={EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                                    />

                                    <DataDisplayRowSelect
                                        label="Category"
                                        selectedValue={datum.category}
                                        options={EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION}
                                    />

                                    <DataDisplayRowSelect
                                        label="Type"
                                        selectedValue={datum.type}
                                        options={EXERCISE_TYPE_WITH_EMPTY_OPTIONS}
                                    />

                                    <DataDisplayRowSelect
                                        label="Status"
                                        selectedValue={datum.status}
                                        options={EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION}
                                    />

                                    <DataDisplayRowCheckbox
                                        label="Has monetization?"
                                        checked={datum.hasMonetization}
                                    />

                                    {datum.hasMonetization && <>
                                        <DataDisplayRowOffer
                                            label="Offer"
                                            offerID={datum.offerId}
                                            helpText=""
                                        />

                                        <DataDisplayRowCheckbox
                                            label="Has Timed Lock?"
                                            checked={datum.hasTimedLock}
                                        />

                                        {datum.hasTimedLock
                                            ? <>
                                                <DataDisplayRowSelect
                                                    label="Timed Lock"
                                                    selectedValue={datum.timedLock}
                                                    options={TIMED_LOCK_DURATION_WITH_EMPTY_OPTIONS}
                                                />
                                            </>
                                            : <>
                                            </>
                                        }
                                    </>}

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/exercises`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Exercises</Link>
                                        </div>
                                        <div class="column is-half has-text-right">
                                            <Link to={`/admin/exercise/${id}/update`} class="button is-warning is-fullwidth-mobile" disabled={datum.type === EXERCISE_TYPE_SYSTEM}><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
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

export default AdminExerciseDetail;
