import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faHandHolding, faTasks, faTachometer, faPlus, faTimesCircle, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faMessage, faChartPie, faCogs, faEye } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getOfferDetailAPI, putOfferUpdateAPI } from "../../../API/Offer";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormSelectField from "../../Reusable/FormSelectField";
import FormInputField from "../../Reusable/FormInputField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";
import DataDisplayRowCheckbox from "../../Reusable/DataDisplayRowCheckbox";
import DataDisplayRowSelect from "../../Reusable/DataDisplayRowSelect";
import {
  OFFER_PAY_FREQUENCY_WITH_EMPTY_OPTIONS,
  BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS,
  OFFER_MEMBERSHIP_RANK_WITH_EMPTY_OPTIONS,
  OFFER_STATUS_WITH_EMPTY_OPTIONS
} from "../../../Constants/FieldOptions";


function AdminOfferUpdate() {
    ////
    //// URL Parameters.
    ////

    const { id } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [datum, setDatum] = useState({});
    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [businessFunction, setBusinessFunction] = useState(0);
    const [status, setStatus] = useState(0);
    const [showCancelWarning, setShowCancelWarning] = useState(false);
    const [membershipRank, setMembershipRank] = useState(0);

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Beginning...");
        setFetching(true);
        setErrors({});

        // To Snake-case for API from camel-case in React.
        const decamelizedData = {
            id: id,
            business_function: parseInt(businessFunction),
            membership_rank: parseInt(membershipRank),
            status: status,
        };
        console.log("onSubmitClick, decamelizedData:", decamelizedData);
        putOfferUpdateAPI(decamelizedData, onAdminOfferUpdateSuccess, onAdminOfferUpdateError, onAdminOfferUpdateDone);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onOfferDetailSuccess(response){
        console.log("onOfferDetailSuccess: Starting...");
        setBusinessFunction(response.businessFunction);
        setMembershipRank(response.membershipRank);
        setStatus(response.status);
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

    // --- Update --- //

    function onAdminOfferUpdateSuccess(response){
        // For debugging purposes only.
        console.log("onAdminOfferUpdateSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Video category update");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAdminOfferUpdateSuccess: Delayed for 2 seconds.");
            console.log("onAdminOfferUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/admin/offer/"+response.id);
    }

    function onAdminOfferUpdateError(apiErr) {
        console.log("onAdminOfferUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminOfferUpdateError: Delayed for 2 seconds.");
            console.log("onAdminOfferUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAdminOfferUpdateDone() {
        console.log("onAdminOfferUpdateDone: Starting...");
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
    }, [id]);
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
                    <nav class="breadcrumb" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/offers" aria-current="page"><FontAwesomeIcon className="fas" icon={faHandHolding} />&nbsp;Offers</Link></li>
                            <li class=""><Link to={`/admin/offer/${id}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link></li>
                        </ul>
                    </nav>
                    <nav class="box">
                        <div class={`modal ${showCancelWarning ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e)=>setShowCancelWarning(false)}></button>
                                </header>
                                <section class="modal-card-body">
                                    Your record will be cancelled and your work will be lost. This cannot be undone. Do you want to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <Link class="button is-medium is-success" to={`/admin/offers`}>Yes</Link>
                                    <button class="button is-medium" onClick={(e)=>setShowCancelWarning(false)}>No</button>
                                </footer>
                            </div>
                        </div>

                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit Offer</p>
                        <FormErrorBox errors={errors} />

                        {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                        {isFetching && <PageLoadingContent displayMessage={"Please wait..."} />}

                        <div class="container">

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

                           <FormSelectField
                                label="Business Function"
                                name="businessFunction"
                                placeholder="Pick"
                                selectedValue={businessFunction}
                                errorText={errors && errors.businessFunction}
                                helpText="Select what beneficial business function the user gets granted upon purchasing / enrolling into this offer."
                                onChange={(e) => setBusinessFunction(parseInt(e.target.value))}
                                options={BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS}
                            />

                            <FormSelectField
                                label="Membership Rank"
                                name="membershipRank"
                                type="number"
                                placeholder="#"
                                selectedValue={membershipRank}
                                errorText={errors && errors.membershipRank}
                                helpText=""
                                onChange={(e)=>setMembershipRank(parseInt(e.target.value))}
                                isRequired={true}
                                options={OFFER_MEMBERSHIP_RANK_WITH_EMPTY_OPTIONS}
                                maxWidth="80px"
                            />

                            <FormSelectField
                                label="Status"
                                name="status"
                                type="number"
                                placeholder="#"
                                selectedValue={status}
                                errorText={errors && errors.status}
                                helpText={
                                  <ul class="content">
                                    <li>pending - will not show up for members</li>
                                    <li>active - will show up for everyone</li>
                                    <li>archived - will be hidden from everyone</li>
                                  </ul>
                                }
                                onChange={(e)=>setStatus(parseInt(e.target.value))}
                                isRequired={true}
                                options={OFFER_STATUS_WITH_EMPTY_OPTIONS}
                                maxWidth="80px"
                            />

                            <div class="columns pt-5">
                                <div class="column is-half">
                                    <button class="button is-medium is-fullwidth-mobile" onClick={(e)=>setShowCancelWarning(true)}><FontAwesomeIcon className="fas" icon={faTimesCircle} />&nbsp;Cancel</button>
                                </div>
                                <div class="column is-half has-text-right">
                                    <button class="button is-medium is-primary is-fullwidth-mobile" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Submit</button>
                                </div>
                            </div>

                        </div>
                    </nav>
                </section>
            </div>
        </>
    );
}

export default AdminOfferUpdate;
