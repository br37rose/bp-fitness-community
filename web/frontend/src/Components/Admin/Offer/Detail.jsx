import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faHandHolding, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faDumbbell, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';

import { getOfferDetailAPI, deleteOfferAPI } from "../../../API/Offer";
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
import DataDisplayRowCheckbox from "../../Reusable/DataDisplayRowCheckbox";
import DataDisplayRowSelect from "../../Reusable/DataDisplayRowSelect";
import FormTextTagRow from "../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../Reusable/FormTextOptionRow";
import {
  OFFER_PAY_FREQUENCY_WITH_EMPTY_OPTIONS,
  BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS,
  OFFER_MEMBERSHIP_RANK_WITH_EMPTY_OPTIONS,
  OFFER_STATUS_WITH_EMPTY_OPTIONS
} from "../../../Constants/FieldOptions";


function AdminOfferDetail() {
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
    const [selectedOfferForDeletion, setSelectedOfferForDeletion] = useState(null);

    ////
    //// Event handling.
    ////

    const onDeleteConfirmButtonClick = () => {
        console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

        deleteOfferAPI(
            selectedOfferForDeletion.id,
            onOfferDeleteSuccess,
            onOfferDeleteError,
            onOfferDeleteDone
        );
        setSelectedOfferForDeletion(null);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onOfferDetailSuccess(response){
        console.log("onOfferDetailSuccess: Starting...");
        setDatum(response);
    }

    function onOfferDetailError(apiErr) {
        console.log("onOfferDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onOfferDetailDone() {
        console.log("onOfferDetailDone: Starting...");
        setFetching(false);
    }

    // --- Delete --- //

    function onOfferDeleteSuccess(response) {
        console.log("onOfferDeleteSuccess: Starting..."); // For debugging purposes only.

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
        setForceURL("/admin/offers");
    }

    function onOfferDeleteError(apiErr) {
        console.log("onOfferDeleteError: Starting..."); // For debugging purposes only.
        setErrors(apiErr);

        // Update notification.
        setTopAlertStatus("danger");
        setTopAlertMessage("Failed deleting");
        setTimeout(() => {
        console.log(
            "onOfferDeleteError: topAlertMessage, topAlertStatus:",
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

    function onOfferDeleteDone() {
        console.log("onOfferDeleteDone: Starting...");
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
            getOfferDetailAPI(
                id,
                onOfferDetailSuccess,
                onOfferDetailError,
                onOfferDetailDone
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
                            <li class=""><Link to="/admin/offers" aria-current="page"><FontAwesomeIcon className="fas" icon={faHandHolding} />&nbsp;Offers</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/offers" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Offers</Link></li>
                        </ul>
                    </nav>

                    {/* Modal */}
                    <nav>
                        {/* Delete modal */}
                        <div class={`modal ${selectedOfferForDeletion !== null ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedOfferForDeletion(null)}></button>
                                </header>
                                <section class="modal-card-body">
                                    You are about to delete this offer and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                                    <button class="button" onClick={(e, ses) => setSelectedOfferForDeletion(null)}>Cancel</button>
                                </footer>
                            </div>
                        </div>
                    </nav>

                    {/* Page */}
                    <nav class="box">
                        {datum && <div class="columns">
                            <div class="column">
                                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faHandHolding} />&nbsp;Offer</p>
                            </div>
                            <div class="column has-text-right">
                                <Link to={`/admin/offer/${id}/update`} class="button is-warning is-fullwidth-mobile" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>&nbsp;
                                <Link onClick={(e,s)=>{setSelectedOfferForDeletion(datum)}} class="button is-danger is-fullwidth-mobile" type="button">
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
                                            <Link to={`/admin/offer/${datum.id}/tags`}>Tags</Link>
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
                                        label="Description"
                                        value={datum.description}
                                    />

                                    <DataDisplayRowText
                                        label="Price"
                                        value={datum.price}
                                    />

                                    <DataDisplayRowText
                                        label="Price Currency"
                                        value={datum.priceCurrency}
                                    />

                                    <DataDisplayRowCheckbox
                                        label="Is Subscription?"
                                        checked={datum.isSubscription}
                                    />

                                    <DataDisplayRowSelect
                                        label="Payment Frequency"
                                        selectedValue={datum.payFrequency}
                                        options={OFFER_PAY_FREQUENCY_WITH_EMPTY_OPTIONS}
                                    />

                                    <DataDisplayRowSelect
                                        label="Business Function"
                                        selectedValue={datum.businessFunction}
                                        options={BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS}
                                        helpText="Select what beneficial business function the user gets granted upon purchasing / enrolling into this offer."
                                    />

                                    <DataDisplayRowSelect
                                        label="Membership Rank"
                                        selectedValue={datum.membershipRank}
                                        options={OFFER_MEMBERSHIP_RANK_WITH_EMPTY_OPTIONS}
                                    />

                                    <DataDisplayRowSelect
                                        label="Status"
                                        selectedValue={datum.status}
                                        options={OFFER_STATUS_WITH_EMPTY_OPTIONS}
                                    />

                                    <div class="columns pt-5">
                                        <div class="column is-half">
                                            <Link class="button is-fullwidth-mobile" to={`/admin/offers`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to offers</Link>
                                        </div>
                                        <div class="column is-half has-text-right">
                                            <Link to={`/admin/offer/${id}/update`} class="button is-warning is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
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

export default AdminOfferDetail;
