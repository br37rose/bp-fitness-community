import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faImage, faEllipsis, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faContactCard, faChartPie, faKey } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getAccountDetailAPI } from "../../API/Account";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormRadioField from "../Reusable/FormRadioField";
import FormMultiSelectField from "../Reusable/FormMultiSelectField";
import FormSelectField from "../Reusable/FormSelectField";
import FormCheckboxField from "../Reusable/FormCheckboxField";
import FormCountryField from "../Reusable/FormCountryField";
import FormRegionField from "../Reusable/FormRegionField";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS, SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS } from "../../Constants/FieldOptions";
import FormTextRow from "../Reusable/FormTextRow";
import FormTextTagRow from "../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../Reusable/FormTextOptionRow";
import DataDisplayRowImage from "../Reusable/DataDisplayRowImage";
import Layout from "../Menu/Layout";


function AccountDetail() {
    ////
    ////
    ////

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
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    ////
    //// Event handling.
    ////

    //

    ////
    //// API.
    ////

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

    ////
    //// BREADCRUMB
    ////

    const generateBreadcrumbItemLink = (currentUser) => {
        let dashboardLink;
        switch (currentUser.role) {
            case 1:
                dashboardLink = "/root/dashboard";
                break;
            case 2:
                dashboardLink = "/admin/dashboard";
                break;
            case 3:
                dashboardLink = "/trainer/dashboard";
                break;
            case 4:
                dashboardLink = "/dashboard";
                break;
            default:
                dashboardLink = "/"; // Default or error handling
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
        return <Navigate to={forceURL} />
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div class="box">

                {/* Title + Options */}
                {currentUser && <div class="columns">
                    <div class="column">
                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faUserCircle} />&nbsp;Account</p>
                    </div>
                    <div class="column has-text-right">
                        {/* Mobile Specific */}
                        <Link to={`/account/change-password`} class="button is-medium is-success is-fullwidth is-hidden-desktop" type="button">
                            <FontAwesomeIcon className="mdi" icon={faKey} />&nbsp;Change Password
                        </Link>
                        {/* Desktop Specific */}
                        <Link to={`/account/change-password`} class="button is-medium is-success is-hidden-touch" type="button">
                            <FontAwesomeIcon className="mdi" icon={faKey} />&nbsp;Change Password
                        </Link>
                    </div>
                </div>}
                <FormErrorBox errors={errors} />

                {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

                {isFetching
                    ?
                    <PageLoadingContent displayMessage={"Please wait..."} />
                    : <>
                        {currentUser && <div class="container">

                            {/* Tab Navigation */}
                            <div class="tabs is-medium is-size-7-mobile">
                                <ul>
                                    <li class="is-active">
                                        <Link><strong>Detail</strong></Link>
                                    </li>
                                    <li>
                                        <Link to={`/account/tags`}>Tags</Link>
                                    </li>
                                    <li>
                                        <Link to={`/account/friends`}>Friends</Link>
                                    </li>
                                    <li>
                                        <Link to={`/account/wearable-tech`}>Wearable Tech</Link>
                                    </li>
                                    <li>
                                        <Link to={`/account/subscription`}>Subscription</Link>
                                    </li>
                                    <li>
                                        <Link to={`/account/2fa`}>2FA</Link>
                                    </li>
                                    <li>
                                        <Link to={`/account/more`}>More&nbsp;<FontAwesomeIcon className="fas" icon={faEllipsis} /></Link>
                                    </li>
                                </ul>
                            </div>

                            {currentUser.avatarObjectUrl !== undefined && currentUser.avatarObjectUrl !== null && currentUser.avatarObjectUrl !== "" && <>
                                <p class="title is-6"><FontAwesomeIcon className="fas" icon={faImage} />&nbsp;Photo</p>
                                <hr />
                                <DataDisplayRowImage
                                    label={`Photo`}
                                    alt={`Download File`}
                                    src={currentUser.avatarObjectUrl}
                                    helpText={``}
                                    maxWidth="380px"
                                />
                            </>}

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;Full Name</p>
                            <hr />

                            <FormTextRow
                                label="First Name"
                                value={currentUser.firstName}
                            />

                            <FormTextRow
                                label="Last Name"
                                value={currentUser.lastName}
                            />

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faContactCard} />&nbsp;Contact Information</p>
                            <hr />

                            <FormTextRow
                                label="Email"
                                value={currentUser.email}
                            />

                            <FormTextRow
                                label="Phone"
                                value={currentUser.phone}
                            />

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faAddressBook} />&nbsp;Address</p>
                            <hr />

                            <FormTextRow
                                label="Country"
                                value={currentUser.country}
                            />

                            <FormTextRow
                                label="Province/Territory"
                                value={currentUser.region}
                            />

                            <FormTextRow
                                label="City"
                                value={currentUser.city}
                            />

                            <FormTextRow
                                label="Address Line 1"
                                value={currentUser.addressLine1}
                            />

                            <FormTextRow
                                label="Address Line 2"
                                value={currentUser.addressLine2}
                            />

                            <FormTextRow
                                label="Postal Code"
                                value={currentUser.postalCode}
                            />

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faChartPie} />&nbsp;Metrics</p>
                            <hr />

                            <FormTextYesNoRow
                                label="I agree to receive electronic updates from my local gym"
                                value={currentUser.agreePromotionsEmail}
                            />

                            {currentUser !== undefined && currentUser !== null && currentUser !== "" && <>
                                {currentUser.stripeSubscription !== undefined && currentUser.stripeSubscription !== null && currentUser.stripeSubscription !== "" && <>
                                    <p class="title is-6"><FontAwesomeIcon className="fas" icon={faRepeat} />&nbsp;Subscription</p>
                                    <hr />
                                    <FormTextOptionRow
                                        label="Interval"
                                        selectedValue={currentUser.stripeSubscription.interval}
                                        options={SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS}
                                    />
                                    <FormTextOptionRow
                                        label="Status"
                                        selectedValue={currentUser.stripeSubscription.status}
                                        options={SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS}
                                    />
                                </>}
                            </>}

                            <div class="columns pt-5">
                                <div class="column is-half">
                                    <Link class="button is-medium is-fullwidth-mobile" to={"/dashboard"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                                </div>
                                <div class="column is-half has-text-right">
                                    <Link to={"/account/update"} class="button is-medium is-primary is-hidden-touch"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
                                    <Link to={"/account/update"} class="button is-medium is-primary is-fullwidth is-hidden-desktop"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
                                </div>
                            </div>

                        </div>}
                    </>
                }
            </div>
        </Layout>
    );
}

export default AccountDetail;
