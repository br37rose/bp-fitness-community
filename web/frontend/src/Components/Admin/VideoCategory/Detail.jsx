import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faDumbbell, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';

import { getVideoCategoryDetailAPI, deleteVideoCategoryAPI } from "../../../API/VideoCategory";
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
import FormTextTagRow from "../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../Reusable/FormTextOptionRow";


function AdminVideoCategoryDetail() {
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
    const [selectedVideoCategoryForDeletion, setSelectedVideoCategoryForDeletion] = useState(null);

    ////
    //// Event handling.
    ////

    const onDeleteConfirmButtonClick = () => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteVideoCategoryAPI(
            selectedVideoCategoryForDeletion.id,
            onVideoCategoryDeleteSuccess,
            onVideoCategoryDeleteError,
            onVideoCategoryDeleteDone
        );
        setSelectedVideoCategoryForDeletion(null);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onVideoCategoryDetailSuccess(response){
        console.log("onVideoCategoryDetailSuccess: Starting...");
        setDatum(response);
    }

    function onVideoCategoryDetailError(apiErr) {
        console.log("onVideoCategoryDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoCategoryDetailDone() {
        console.log("onVideoCategoryDetailDone: Starting...");
        setFetching(false);
    }

    // --- Delete --- //

    function onVideoCategoryDeleteSuccess(response) {
        console.log("onVideoCategoryDeleteSuccess: Starting..."); // For debugging purposes only.

        // Update notification.
        setTopAlertStatus("success");
        setTopAlertMessage("Video category deleted");
        setTimeout(() => {
        console.log(
            "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
            topAlertMessage,
            topAlertStatus
        );
        setTopAlertMessage("");
        }, 2000);

        // Redirect back to the video categories page.
        setForceURL("/admin/video-categories");
    }

    function onVideoCategoryDeleteError(apiErr) {
        console.log("onVideoCategoryDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
        console.log(
            "onVideoCategoryDeleteError: topAlertMessage, topAlertStatus:",
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

    function onVideoCategoryDeleteDone() {
        console.log("onVideoCategoryDeleteDone: Starting...");
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
            getVideoCategoryDetailAPI(
                id,
                onVideoCategoryDetailSuccess,
                onVideoCategoryDetailError,
                onVideoCategoryDetailDone
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
                            <li class=""><Link to="/admin/video-categories" aria-current="page"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;Video Categories</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/video-categories" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Video Categories</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    <nav>
                        {/* Delete modal */}
                        <div class={`modal ${selectedVideoCategoryForDeletion !== null ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedVideoCategoryForDeletion(null)}></button>
                                </header>
                                <section class="modal-card-body">
                                    You are about to delete this offer and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                                    <button class="button" onClick={(e, ses) => setSelectedVideoCategoryForDeletion(null)}>Cancel</button>
                                </footer>
                            </div>
                        </div>
                    </nav>

                    {/* Page */}
                    <nav class="box">
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;Video Category</p>
                            </div>
                            <div class="column has-text-right">
                                <Link to={`/admin/video-category/${id}/update`} class="button is-warning is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedVideoCategoryForDeletion(datum)}} class="button is-danger is-small is-fullwidth-mobile" type="button">
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

                                    {/* Tab navigation */}
                                    {/*
                                    <div class= "tabs is-medium is-size-7-mobile">
                                      <ul>
                                        <li class="is-active">
                                            <Link><strong>Detail</strong></Link>
                                        </li>
                                        <li>
                                            <Link to={`/admin/video-category/${datum.id}/tags`}>Tags</Link>
                                        </li>
                                      </ul>
                                    </div>
                                    */}

                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</p>
                                    <hr />

                                    <DataDisplayRowText
                                        label="Name"
                                        value={datum.name}
                                    />

                                    <DataDisplayRowText
                                        label="No #"
                                        value={datum.no}
                                    />

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/video-categories`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to video categories</Link>
                                        </div>
                                        <div class="column is-half has-text-right">
                                            <Link to={`/admin/video-category/${id}/update`} class="button is-warning is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
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

export default AdminVideoCategoryDetail;
