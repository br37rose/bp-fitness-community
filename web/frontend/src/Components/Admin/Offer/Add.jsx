import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowUpRightFromSquare, faInfoCircle, faArrowLeft, faHandHolding, faTasks, faTachometer, faPlus, faTimesCircle, faCheckCircle, faUserCircle, faGauge, faPencil, faDumbbell, faIdCard, faAddressBook, faMessage, faChartPie, faCogs, faEye } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getOfferDetailAPI, postOfferCreateAPI } from "../../../API/Offer";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";


function AdminOfferAdd() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [name, setName] = useState("");
    const [no, setNo] = useState("");
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    ////
    //// Event handling.
    ////

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Beginning...");
        setFetching(true);
        setErrors({});

        // To Snake-case for API from camel-case in React.
        const decamelizedData = {
            no: parseInt(no),
            name: name,
        };
        console.log("onSubmitClick, decamelizedData:", decamelizedData);
        postOfferCreateAPI(decamelizedData, onAdminOfferAddSuccess, onAdminOfferAddError, onAdminOfferAddDone);
    }

    function onAdminOfferAddSuccess(response){
        // For debugging purposes only.
        console.log("onAdminOfferAddSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Video category created");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAdminOfferAddSuccess: Delayed for 2 seconds.");
            console.log("onAdminOfferAddSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/admin/offer/"+response.id);
    }

    function onAdminOfferAddError(apiErr) {
        console.log("onAdminOfferAddError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminOfferAddError: Delayed for 2 seconds.");
            console.log("onAdminOfferAddError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAdminOfferAddDone() {
        console.log("onAdminOfferAddDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
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
                    <nav class="breadcrumb" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/offers" aria-current="page"><FontAwesomeIcon className="fas" icon={faHandHolding} />&nbsp;Offers</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New</Link></li>
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

                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New Offer</p>
                        <FormErrorBox errors={errors} />

                        {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                        {isFetching && <PageLoadingContent displayMessage={"Please wait..."} />}

                        <div class="container content">

                            <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faInfoCircle} />&nbsp;How to Add an Offer via Stripe</p>
                            <hr />
                            <p>Please follow the following steps if you would like to add an offer into the system.</p>
                            <ol>
                                <li>Log into <Link to="https://stripe.com"  target="_blank" rel="noreferrer">Stripe, Inc.&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link></li>
                                <li>Go to "Products"</li>
                                <li>Click "Add Product"</li>
                                <li>Fill in the form and please include the following:
                                    <ul>
                                        <li>Name</li>
                                        <li>Description</li>
                                        <li>Image (only upload one)</li>
                                        <li>Metadata (key="OrganizationID" value="648763d3f6fbead15f5bd4d2")</li>
                                        <li>Recurring: Yes</li>
                                        <li>Monthly billing period</li>
                                        <li>Price</li>
                                    </ul>
                                </li>
                                <li>Click "Save product"</li>
                                <li>Inside this web-app, go back to the offers list, reload and you should see the new offer automatically populated.</li>
                            </ol>



                            <div class="columns pt-5">
                                <div class="column is-half">
                                    <Link to="/admin/offers" class="button is-medium is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back</Link>
                                </div>
                                <div class="column is-half has-text-right">
                                    {/*
                                    <Link to="/admin/offers" class="button is-medium is-primary is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;I understand</Link>
                                    */}
                                </div>
                            </div>

                        </div>
                    </nav>
                </section>
            </div>
        </>
    );
}

export default AdminOfferAdd;
