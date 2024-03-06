import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faImage, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faVideo, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';
import Vimeo from '@u-wave/react-vimeo';

import { getVideoContentDetailAPI } from "../../../../API/VideoContent";
import DataDisplayRowOffer from "../../../Reusable/DataDisplayRowOffer";
import DataDisplayRowCheckbox from "../../../Reusable/DataDisplayRowCheckbox";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../Reusable/FormCountryField";
import FormRegionField from "../../../Reusable/FormRegionField";
import DataDisplayRowURL from "../../../Reusable/DataDisplayRowURL";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../../AppState";
import {
    VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION,
    VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    TIMED_LOCK_DURATION_WITH_EMPTY_OPTIONS
} from "../../../../Constants/FieldOptions";
import {
    EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL,
    EXERCISE_TYPE_SYSTEM,
    EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_VIDEO_TYPE_YOUTUBE,
    EXERCISE_VIDEO_TYPE_VIMEO
} from "../../../../Constants/App";
import DataDisplayRowText from "../../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../../Reusable/DataDisplayRowRadio";
import DataDisplayRowSelect from "../../../Reusable/DataDisplayRowSelect";


function MemberVideoContentDetail() {
    ////
    //// URL Parameters.
    ////

    const { vcid, vconid } = useParams()

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

    ////
    //// Event handling.
    ////


    ////
    //// API.
    ////

    // --- Detail --- //

    function onVideoContentDetailSuccess(response){
        console.log("onVideoContentDetailSuccess: Starting...");
        setDatum(response);
    }

    function onVideoContentDetailError(apiErr) {
        console.log("onVideoContentDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVideoContentDetailDone() {
        console.log("onVideoContentDetailDone: Starting...");
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
            getVideoContentDetailAPI(
                vconid,
                onVideoContentDetailSuccess,
                onVideoContentDetailError,
                onVideoContentDetailDone
            );
        }

        return () => { mounted = false; }
    }, [vconid,]);
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
                            <li class=""><Link to="/video-collections" aria-current="page"><FontAwesomeIcon className="fas" icon={faVideo} />&nbsp;Video Collections</Link></li>
                            <li class=""><Link to={`/video-collection/${vcid}/video-contents`} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail&nbsp;(Video Content)</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Video Content</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to={`/video-collection/${vcid}/video-contents`} aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Details (Video Content)</Link></li>
                        </ul>
                    </nav>

                    {/* Page */}
                    <nav class="box">

                        {/* Title + Options */}
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faVideo} />&nbsp;Video Content</p>
                            </div>
                            <div class="column has-text-right">
                                {/*
                                <Link to={`/video-collection/${datum.collectionId}/video-content/${datum.id}/update`}class="button is-warning is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedVideoContentForDeletion(datum)}} class="button is-danger is-small is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Delete
                                </Link>
                                */}
                            </div>
                        </div>}

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
                                                <>YouTube (TODO)</>
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
                                    />

                                    <DataDisplayRowText
                                        label="No #"
                                        value={datum.no}
                                    />

                                    <DataDisplayRowText
                                        label="Description"
                                        value={datum.description}
                                    />

                                    <DataDisplayRowText
                                        label="Author Name"
                                        value={datum.authorName}
                                    />

                                    <DataDisplayRowText
                                        label="Author URL Address"
                                        value={datum.authorUrl}
                                    />

                                    <DataDisplayRowText
                                        label="Duration"
                                        value={datum.duration}
                                    />

                                    <DataDisplayRowURL
                                        label="Category"
                                        urlKey={datum.categoryName}
                                        urlValue={`/video-category/${datum.categoryId}`}
                                        type="external"
                                    />

                                    <DataDisplayRowText
                                        label="Collection"
                                        value={datum.collectionName}
                                    />

                                    <DataDisplayRowURL
                                        label="Collection"
                                        urlKey={datum.collectionName}
                                        urlValue={`/video-collection/${datum.collectionId}`}
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

                                    <DataDisplayRowCheckbox
                                        label="Has monetization?"
                                        checked={datum.hasMonetization}
                                    />

                                    {datum.hasMonetization && <>
                                        <DataDisplayRowURL
                                            label="Offer"
                                            urlKey={datum.offerName}
                                            urlValue={`/offer/${datum.offerId}`}
                                            type="external"
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
                                            <Link class="button is-fullwidth-mobile" to={`/video-collection/${datum.collectionId}/video-contents`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to video contents</Link>
                                        </div>
                                        <div class="column is-half has-text-right">

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

export default MemberVideoContentDetail;
