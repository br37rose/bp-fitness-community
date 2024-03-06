import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
    faEye, faBuilding, faEnvelope, faSquarePhone, faTable, faHome, faLock, faTimesCircle, faArrowRight, faArrowUpRightFromSquare,
    faImage, faEllipsis, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faContactCard, faChartPie, faKey, faUnlock
} from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import QRCode from "qrcode.react";

import { getAccountDetailAPI } from "../../../API/Account";
import { postVertifyOTP } from "../../../API/gateway";
import { postGenerateOTP } from "../../../API/gateway";
import DateTimeTextFormatter from "../../Reusable/EveryPage/DateTimeTextFormatter";
import CheckboxTextFormatter from "../../Reusable/EveryPage/CheckboxTextFormatter";
import SelectTextFormatter from "../../Reusable/EveryPage/SelectTextFormatter";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import { topAlertMessageState, topAlertStatusState, currentUserState, currentOTPResponseState } from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputField from "../../Reusable/FormInputField";
import DataDisplayRowImage from "../../Reusable/DataDisplayRowImage";
import TagsTextFormatter from "../../Reusable/EveryPage/TagsTextFormatter";
import URLTextFormatter from "../../Reusable/EveryPage/URLTextFormatter";
import EmailTextFormatter from "../../Reusable/EveryPage/EmailTextFormatter";
import PhoneTextFormatter from "../../Reusable/EveryPage/PhoneTextFormatter";
import DateTextFormatter from "../../Reusable/EveryPage/DateTextFormatter";
import { COMMERCIAL_CUSTOMER_TYPE_OF_ID } from "../../../Constants/App";
import { ROOT_ROLE_ID, ADMIN_ROLE_ID, TRAINER_ROLE_ID, MEMBER_ROLE_ID } from "../../../Constants/App";
import {
    USER_PHONE_TYPE_OF_OPTIONS_WITH_EMPTY_OPTIONS,
    USER_TYPE_OF_FILTER_OPTIONS,
    USER_ORGANIZATION_TYPE_OPTIONS,
    GENDER_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";


function AccountEnableTwoFactorAuthentication() {
    ////
    ////
    ////

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [otpResponse, setOtpResponse] = useRecoilState(currentOTPResponseState);

    ////
    //// Component states.
    ////

    // Page related states.
    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    // Modal related states.
    const [verificationToken, setVerificationToken] = useState("");


    ////
    //// Event handling.
    ////

    const onVerificationTokenSubmitButtonClick = () => {
        // Remove whitespace characters from verificationToken
        const cleanedVerificationToken = verificationToken.replace(/\s/g, '');

        const payload= {
            verification_token: cleanedVerificationToken,
        }
        postVertifyOTP(
            payload,
            onVerifyOPTSuccess,
            onVerifyOPTError,
            onVerifyOPTDone
        );
    }

    ////
    //// API.
    ////

    // --- Generate OTP --- //

    function onGenerateOPTSuccess(response){
        console.log("onGenerateOPTSuccess: Starting...");
        console.log("response: ", response);
        setOtpResponse(response);
    }

    function onGenerateOPTError(apiErr) {
        console.log("onGenerateOPTError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onGenerateOPTDone() {
        console.log("onGenerateOPTDone: Starting...");
        setFetching(false);
    }

    // --- Enable 2FA --- //

    function onVerifyOPTSuccess(response){
        console.log("onVerifyOPTSuccess: Starting...");
        if (response !== undefined && response !== null && response !== "") {
            console.log("response: ", response);
            if (response.user !== undefined && response.user !== null && response.user !== "") {
                console.log("response.user: ", response.user);

                // Clear errors.
                setErrors({});

                // Save our updated user account.
                setCurrentUser(response.user);

                // Add a temporary banner message in the app and then clear itself after 2 seconds.
                setTopAlertMessage("2FA Enabled");
                setTopAlertStatus("success");
                setTimeout(() => {
                    console.log("onSuccess: Delayed for 2 seconds.");
                    console.log("onSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
                    setTopAlertMessage("");
                }, 2000);

                // Change page.
                setForceURL("/account/2fa");
            }
        }

    }

    function onVerifyOPTError(apiErr) {
        console.log("onVerifyOPTError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed Enabling 2FA");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onSuccess: Delayed for 2 seconds.");
            console.log("onSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVerifyOPTDone() {
        console.log("onVerifyOPTDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////

    const generateBreadcrumbItemLink = (currentUser) => {
        let dashboardLink;
        switch (currentUser.role) {
            case ROOT_ROLE_ID:
                dashboardLink = "/admin/dashboard";
                break;
            case ADMIN_ROLE_ID:
                dashboardLink = "/admin/dashboard";
                break;
            case TRAINER_ROLE_ID:
                dashboardLink = "/admin/dashboard";
                break;
            case MEMBER_ROLE_ID:
                dashboardLink = "/dashboard";
                break;
            default:
                dashboardLink = "/501"; // Default or error handling
                break;
        }
        return dashboardLink;
    };

    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: generateBreadcrumbItemLink(currentUser), isActive: false, icon: faGauge },
            { text: 'Account', link: '/account', icon: faUserCircle, isActive: true }
        ],
        mobileBackLinkItems: {
            link: generateBreadcrumbItemLink(currentUser),
            text: "Back to Dashboard",
            icon: faArrowLeft
        }
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

            postGenerateOTP(
                onGenerateOPTSuccess,
                onGenerateOPTError,
                onGenerateOPTDone
            );

            // if (otpResponse === undefined || otpResponse === null || otpResponse === "") {
            //
            // }
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <>
            <div className="container">
                <section className="section">
                    {/* Desktop Breadcrumbs */}
                    <nav className="breadcrumb has-background-light p-4 is-hidden-touch" aria-label="breadcrumbs">
                        <ul>
                            <li className=""><Link to={generateBreadcrumbItemLink(currentUser)} aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li><Link to={"/account/2fa"} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail (2FA)</Link></li>
                            <li className="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faLock} />&nbsp;Enable 2FA</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav className="breadcrumb has-background-light p-4 is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li className="">
                                <Link to="/account/2fa" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Detail (2FA)</Link>
                            </li>
                        </ul>
                    </nav>

                    {/* Page Title */}
                    <h1 className="title is-2"><FontAwesomeIcon className="fas" icon={faUserCircle} />&nbsp;Profile</h1>
                    <h4 className="subtitle is-4"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail (2FA)</h4>
                    <hr />

                    {/* Page */}
                    <nav className="box">
                        {isFetching
                            ?
                            <PageLoadingContent displayMessage={"Loading..."} />
                            :
                            <>
                                <FormErrorBox errors={errors} />
                                {currentUser && <div className="container">

                                    <div className="content">

                                            <div className="columns">
                                                {/*
                                                    ------
                                                    STEP 1
                                                    ------
                                                */}
                                                <div className="column">

                                                    <h1 className="title is-4">Step 1</h1>
                                                    <hr />
                                                    <p class="has-text-grey">To begin, please download any of the following applications for your mobile device.</p>
                                                    {/* Google Authenticator */}
                                                    <div className="card">
                                                        <div className="card-content">
                                                            <div className="media">
                                                              <div className="media-content">
                                                                <p className="title is-6"><u>Google Authenticator</u></p>
                                                              </div>
                                                            </div>

                                                            <div className="content">
                                                                <p>This 2FA app is created by <b>Google, Inc.</b></p>
                                                                <b>Download for iOS:</b>&nbsp;
                                                                <Link className="" to="https://apps.apple.com/ca/app/google-authenticator/id388497605" target="_blank" rel="noreferrer">
                                                                   Visit App Store&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                                                </Link>
                                                                <br />
                                                                <b>Download for Android:</b>&nbsp;
                                                                <Link className="" to="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&pli=1" target="_blank" rel="noreferrer">
                                                                   Visit Google Play&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                                                </Link>
                                                                <br />
                                                            </div>
                                                        </div>
                                                        {/*}
                                                        <footer className="card-footer">
                                                            <button className="card-footer-item button is-primary is-small">
                                                               Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                                                            </button>
                                                        </footer>
                                                        */}
                                                    </div>

                                                    {/* Authenticator Chrome Extension */}
                                                    <div className="card">
                                                        <div className="card-content">
                                                            <div className="media">

                                                              <div className="media-content">
                                                                <p className="title is-6"><u>Authenticator</u></p>
                                                              </div>
                                                            </div>

                                                            <div className="content">
                                                                <p>This 2FA app is created by <b>authenticator.cc</b></p>
                                                                <b>Download for Chrome:</b>&nbsp;
                                                                <Link className="" to="https://chromewebstore.google.com/detail/authenticator/bhghoamapcdpbohphigoooaddinpkbai?pli=1" target="_blank" rel="noreferrer">
                                                                   Visit Chrome web store&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} />
                                                                </Link>
                                                                <br />
                                                            </div>
                                                        </div>
                                                        {/*}
                                                        <footer className="card-footer">
                                                            <button className="card-footer-item button is-primary is-small">
                                                               Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                                                            </button>
                                                        </footer>
                                                        */}
                                                    </div>
                                                </div>

                                                {/*
                                                    ------
                                                    STEP 2
                                                    ------
                                                */}
                                                <div className="column">
                                                    <h1 className="title is-4">Step 2</h1>
                                                    <hr />
                                                    <p class="has-text-grey">With your 2FA application open, please scan the following QR code with your device and click next when ready.</p>
                                                    <p>&nbsp;</p>
                                                    <div className="columns is-centered is-hidden-mobile">
                                                        <div class="column is-half">
                                                            <figure class="image">
                                                                {otpResponse && <QRCode value={otpResponse.optAuthURL} size={100} className="" />}
                                                                <br />
                                                                <span className="is-centered">Scan with your app</span>
                                                            </figure>
                                                        </div>
                                                    </div>

                                                    <div className="columns is-centered is-hidden-tablet">
                                                        <div class="column is-half">
                                                            <figure class="image">
                                                                {otpResponse && <QRCode value={otpResponse.optAuthURL} size={200} className="" />}
                                                                <br />
                                                                <span className="is-centered">Scan with your app</span>
                                                            </figure>
                                                        </div>
                                                    </div>
                                                </div>
                                                {/*
                                                    ------
                                                    STEP 3
                                                    ------
                                                */}
                                                <div className="column">
                                                    <h1 className="title is-4">Step 3</h1>
                                                    <hr />
                                                    <p class="has-text-grey">Open the two-step verification app on your mobile device to get your verification code.</p>
                                                    <p>&nbsp;</p>

                                                    <FormInputField
                                                        label="Enter your Verification Token"
                                                        name="verificationToken"
                                                        placeholder="See your authenticator app"
                                                        value={verificationToken}
                                                        errorText={errors && errors.verificationToken}
                                                        helpText=""
                                                        onChange={(e)=>setVerificationToken(e.target.value)}
                                                        isRequired={true}
                                                        maxWidth="380px"
                                                    />
                                                    <button className="button is-success is-fullwidth-mobile" onClick={onVerificationTokenSubmitButtonClick}><FontAwesomeIcon className="fas" icon={faLock} />&nbsp;Submit Verification Token</button>
                                                </div>
                                            </div>

                                    </div>

                                    <div className="columns pt-5">
                                        <div className="column is-half">
                                            <Link className="button is-fullwidth-mobile" to={`/account/2fa`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Detail (2FA)</Link>
                                        </div>
                                        <div className="column is-half has-text-right">

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

export default AccountEnableTwoFactorAuthentication;
