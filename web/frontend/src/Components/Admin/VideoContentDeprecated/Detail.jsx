import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faVideo, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faDumbbell, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';
import Vimeo from '@u-wave/react-vimeo';

import { getVideoCollectionDetailAPI, deleteVideoCollectionAPI } from "../../../API/VideoCollection";
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
import {
    EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_TYPE_WITH_EMPTY_OPTIONS,
    EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION
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


function AdminVideoCollectionDetail() {
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
    const [selectedVideoCollectionForDeletion, setSelectedVideoCollectionForDeletion] = useState(null);

    ////
    //// Event handling.
    ////

    const onDeleteConfirmButtonClick = () => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteVideoCollectionAPI(
            selectedVideoCollectionForDeletion.id,
            onVideoCollectionDeleteSuccess,
            onVideoCollectionDeleteError,
            onVideoCollectionDeleteDone
        );
        setSelectedVideoCollectionForDeletion(null);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onVideoCollectionDetailSuccess(response){
        console.log("onVideoCollectionDetailSuccess: Starting...");
        setDatum(response);
    }

    function onVideoCollectionDetailError(apiErr) {
        console.log("onVideoCollectionDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoCollectionDetailDone() {
        console.log("onVideoCollectionDetailDone: Starting...");
        setFetching(false);
    }

    // --- Delete --- //

    function onVideoCollectionDeleteSuccess(response) {
        console.log("onVideoCollectionDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("VideoCollection deleted");
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

    function onVideoCollectionDeleteError(apiErr) {
        console.log("onVideoCollectionDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
        console.log(
            "onVideoCollectionDeleteError: topAlertMessage, topAlertStatus:",
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

    function onVideoCollectionDeleteDone() {
        console.log("onVideoCollectionDeleteDone: Starting...");
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
            getVideoCollectionDetailAPI(
                id,
                onVideoCollectionDetailSuccess,
                onVideoCollectionDetailError,
                onVideoCollectionDetailDone
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
                            <li class=""><Link to="/admin/exercises" aria-current="page"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;VideoCollections</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/exercises" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to VideoCollections</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    <nav>
                        {/* Delete modal */}
                        <div class={`modal ${selectedVideoCollectionForDeletion !== null ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedVideoCollectionForDeletion(null)}></button>
                                </header>
                                <section class="modal-card-body">
                                    You are about to delete this member and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                                    <button class="button" onClick={(e, ses) => setSelectedVideoCollectionForDeletion(null)}>Cancel</button>
                                </footer>
                            </div>
                        </div>
                    </nav>

                    {/* Page */}
                    <nav class="box">
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;VideoCollection</p>
                            </div>
                            <div class="column has-text-right">
                                <Link to={`/admin/exercise/${id}/update`} class="button is-medium is-warning is-fullwidth-mobile" type="button" disabled={datum.type === EXERCISE_TYPE_SYSTEM}>
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedVideoCollectionForDeletion(datum)}} class="button is-medium is-danger is-fullwidth-mobile" type="button" disabled={datum.type === EXERCISE_TYPE_SYSTEM}>
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Delete
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

                                    <DataDisplayRowText
                                        label="Name"
                                        value={datum.name}
                                        autoplay={false}
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

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/exercises`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to VideoCollections</Link>
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

export default AdminVideoCollectionDetail;
