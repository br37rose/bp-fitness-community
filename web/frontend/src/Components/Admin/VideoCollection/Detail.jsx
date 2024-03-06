import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrash, faImage, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faVideo, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
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
import DataDisplayRowURL from "../../Reusable/DataDisplayRowURL";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
    VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import {
    EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL,
    EXERCISE_TYPE_SYSTEM
} from "../../../Constants/App";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../Reusable/DataDisplayRowRadio";
import DataDisplayRowSelect from "../../Reusable/DataDisplayRowSelect";


function AdminVideoCollectionDetail() {
    ////
    //// URL Parameters.
    ////

    const { vcid } = useParams()

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
        setTopAlertMessage("Video collection deleted");
        setTimeout(() => {
        console.log(
            "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
            topAlertMessage,
            topAlertStatus
        );
        setTopAlertMessage("");
        }, 2000);

        // Redirect back to the members page.
        setForceURL("/admin/video-collections");
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
                vcid,
                onVideoCollectionDetailSuccess,
                onVideoCollectionDetailError,
                onVideoCollectionDetailDone
            );
        }

        return () => { mounted = false; }
    }, [vcid,]);
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
                            <li class=""><Link to="/admin/video-collections" aria-current="page"><FontAwesomeIcon className="fas" icon={faVideo} />&nbsp;Video Collections</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/video-collections" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to video collections</Link></li>
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
                                    You are about to delete this video collection and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
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

                        {/* Title + Options */}
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faVideo} />&nbsp;Video Collection</p>
                            </div>
                            <div class="column has-text-right">
                                <Link to={`/admin/video-collection/${vcid}/update`} class="button is-warning is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedVideoCollectionForDeletion(datum)}} class="button is-danger is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faTrash} />&nbsp;Delete
                                </Link>
                            </div>
                        </div>}

                        {/* Tab Navigation */}
                        <div class= "tabs is-medium is-size-7-mobile">
                          <ul>
                            <li class="is-active">
                                <Link><strong>Detail</strong></Link>
                            </li>
                            <li>
                                <Link to={`/admin/video-collection/${vcid}/video-contents`}>Contents</Link>
                            </li>
                          </ul>
                        </div>

                        {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

                        {isFetching
                            ?
                            <PageLoadingContent displayMessage={"Please wait..."} />
                            :
                            <>
                                <FormErrorBox errors={errors} />
                                {datum && <div class="container" key={datum.id}>
                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faImage} />&nbsp;Thumbnail</p>
                                    <hr />

                                    <div class="field pb-4">
                                        <label class="label">Preview Image</label>
                                        <div class="control">
                                            {(() => {
                                                switch (datum.thumbnailType) {
                                                    case EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER:
                                                        if (datum.thumbnailObjectUrl !== undefined && datum.thumbnailObjectUrl !== null && datum.thumbnailObjectUrl !== "") {
                                                            return (
                                                                <div className="has-background-black box has-text-white has-text-centered is-size-3" style={{borderRadius: "20px"}}>
                                                                    <img src={datum.thumbnailObjectUrl} alt="Image URL" />
                                                                </div>
                                                            );
                                                        } else {
                                                            return (
                                                                <p>-</p>
                                                            );
                                                        }
                                                    case EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL:
                                                        if (datum.thumbnailUrl !== undefined && datum.thumbnailUrl !== null && datum.thumbnailUrl !== "") {
                                                            return (
                                                                <div className="has-background-black box has-text-white has-text-centered is-size-3" style={{borderRadius: "20px"}}>
                                                                    <img src={datum.thumbnailUrl} alt="Image URL" />
                                                                </div>
                                                            );
                                                        } else {
                                                            return (
                                                                <p>-</p>
                                                            );
                                                        }
                                                    default: return null;
                                                }
                                            })()}
                                        </div>
                                    </div>

                                    <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Information</p>
                                    <hr />

                                    <DataDisplayRowText
                                        label="Name"
                                        value={datum.name}
                                    />

                                    <DataDisplayRowText
                                        label="Summary"
                                        value={datum.summary}
                                    />

                                    <DataDisplayRowText
                                        label="Description"
                                        value={datum.description}
                                    />

                                    <DataDisplayRowURL
                                        label="Category"
                                        urlKey={datum.categoryName}
                                        urlValue={`/admin/video-category/${datum.categoryId}`}
                                        type="external"
                                    />

                                    <DataDisplayRowSelect
                                        label="Type"
                                        selectedValue={datum.type}
                                        options={VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                                    />

                                    <DataDisplayRowSelect
                                        label="Status"
                                        selectedValue={datum.status}
                                        options={VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION}
                                    />

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/video-collections`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to video collections</Link>
                                        </div>
                                        <div class="column is-half has-text-right">
                                            <Link to={`/admin/video-collection/${vcid}/update`} class="button is-warning is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
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
