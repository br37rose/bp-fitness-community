import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsis, faArrowUpRightFromSquare, faHeartPulse, faArrowRight, faTable, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faContactCard, faChartPie, faKey } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import RedirectURL from "../../../Hooks/RedirectURL";
import { getAccountDetailAPI } from "../../../API/Account";
import { getGoogleFitRegistrationURLAPI, postFitBitAppCreateSimulatorAPI } from "../../../API/Wearable";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";
import { SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS, SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS } from "../../../Constants/FieldOptions";
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
import FormTextRow from "../../Reusable/FormTextRow";
import FormTextTagRow from "../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../Reusable/FormTextOptionRow";
import Layout from "../../Menu/Layout";


function AccountWearableTechLaunchpad() {
    ////
    ////
    ////

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");

    ////
    //// Event handling.
    ////

    const onRegisterClick = (e) => {
        setFetching(true);
        setErrors({});
        getGoogleFitRegistrationURLAPI(
            onRegistrationSuccess,
            onRegistrationError,
            onRegistrationDone
        );
    }

    const onCreateSimulator = (e) => {
        setFetching(true);
        setErrors({});
        postFitBitAppCreateSimulatorAPI(
            currentUser.id,
            "random",
            onCreateSimulatorSuccess,
            onCreateSimulatorError,
            onCreateSimulatorDone
        );
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onAccountDetailSuccess(response) {
        console.log("onAccountDetailSuccess: Starting...");
        setCurrentUser(response);
    }

    function onAccountDetailError(apiErr) {
        console.log("onAccountDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountDetailDone() {
        console.log("onAccountDetailDone: Starting...");
        setFetching(false);
    }

    // --- Simulator --- //

    function onCreateSimulatorSuccess(response) {
        console.log("onCreateSimulatorSuccess: Starting...");
        window.location.reload();
    }

    function onCreateSimulatorError(apiErr) {
        console.log("onCreateSimulatorError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onCreateSimulatorDone() {
        console.log("onCreateSimulatorDone: Starting...");
        setFetching(false);
    }

    // --- Register --- //

    function onRegistrationSuccess(response) {
        console.log("onRegistrationSuccess: Starting...");
        setForceURL(response.url);
    }

    function onRegistrationError(apiErr) {
        console.log("onRegistrationError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onRegistrationDone() {
        console.log("onRegistrationDone: Starting...");
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
            setErrors({});
            getAccountDetailAPI(
                onAccountDetailSuccess,
                onAccountDetailError,
                onAccountDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <RedirectURL url={forceURL} />
    }

    return (
        <div>

            <div className="columns">
                <div className="column">
                    {/* Subtitle */}
                    <p class="title is-6"><FontAwesomeIcon className="fas" icon={faHeartPulse} />&nbsp;Google Fit Fitness Tracker</p>
                    <hr />

                    {/* Empty list */}
                    {currentUser.primaryHealthTrackingDeviceType === 0 && <section className="hero has-background-white-ter">
                        <div className="hero-body">
                            <p className="title">
                                <FontAwesomeIcon className="fas" icon={faHeartPulse} />
                                &nbsp;No Connection
                            </p>
                            <p className="subtitle">
                                Your Google Fit fitness tracker is not connected with us.{" "}
                                <b>
                                    <Link click={onRegisterClick}>Click here&nbsp;<FontAwesomeIcon className="mdi" icon={faArrowUpRightFromSquare} /> </Link>
                                </b>{" "}
                                to get started by registering your device and let us read the latest your biometrics. We currently we will extract the following data from your device:
                                <div className="content">
                                    <ul>
                                        <li>Activity and exercise</li>
                                        <li>Heart rate</li>
                                    </ul>
                                    <p> However, we ask for numerous other biometrics so when new features come around, we will support those device types. Please accept those as well.</p>
                                    <br />
                                    <p><i>DEVELOPERS ONLY: <b><Link onClick={onCreateSimulator}>Click here&nbsp;<FontAwesomeIcon className="mdi" icon={faArrowRight} /> </Link></b>{" "}attach a simulator fitbit with fake data.</i></p>
                                </div>
                            </p>
                        </div>
                    </section>}

                    {/* Google Fit */}
                    {currentUser.primaryHealthTrackingDeviceType === 1 && <section className="hero has-background-white-ter">
                        <div className="hero-body">
                            <p className="title">
                                <FontAwesomeIcon className="fas" icon={faHeartPulse} />
                                &nbsp;Google Fit Connected
                            </p>
                            <p className="subtitle">
                                Your Google Fit fitness tracker is connected with us.{" "}
                            </p>
                        </div>
                    </section>}
                </div>
            </div>

            <div class="columns pt-5">
                <div class="column is-half">
                    <Link class="button is-medium is-fullwidth-mobile" to={"/dashboard"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                </div>
                <div class="column is-half has-text-right">
                    {currentUser.primaryHealthTrackingDeviceType === 0 && <button onClick={onRegisterClick} class="button is-medium is-success is-fullwidth-mobile" type="button">
                        Register My Google Fit&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                    </button>}
                </div>
            </div>
        </div>
    );
}

export default AccountWearableTechLaunchpad;
