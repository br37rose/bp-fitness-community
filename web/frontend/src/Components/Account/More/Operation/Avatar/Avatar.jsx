import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faImage, faFile, faTasks, faTachometer, faPlus, faTimesCircle, faCheckCircle, faUserCircle, faGauge, faPencil, faUsers, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs, faEye, faArrowLeft } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { postAccountAvatarAPI } from "../../../../../API/Account";
import { getAccountDetailAPI } from "../../../../../API/Account";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../../Reusable/FormCountryField";
import FormRegionField from "../../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../../../AppState";
import Layout from "../../../../Menu/Layout";


function AdminAccountAvatarOperation() {
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
    const [selectedFile, setSelectedFile] = useState(null);
    const [currentUser, setCurrentUser] = useState({});

    ////
    //// Event handling.
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

    // --- Avatar Upload --- //

    const onHandleFileChange = (event) => {
        setSelectedFile(event.target.files[0]);
    };

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Starting...");
        console.log("onSubmitClick: user_id:", currentUser.id);
        setFetching(true);
        setErrors({});

        const formData = new FormData();
        formData.append('user_id', currentUser.id);
        formData.append('file', selectedFile);

        postAccountAvatarAPI(
            formData,
            onOperationSuccess,
            onOperationError,
            onOperationDone
        );
        console.log("onSubmitClick: Finished.")
    }

    ////
    //// API.
    ////

    // --- Avatar Operation --- //

    function onOperationSuccess(response) {
        // For debugging purposes only.
        console.log("onOperationSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Photo changed");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onOperationSuccess: Delayed for 2 seconds.");
            console.log("onOperationSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to the user attachments page.
        setForceURL("/account");
    }

    function onOperationError(apiErr) {
        console.log("onOperationError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onOperationError: Delayed for 2 seconds.");
            console.log("onOperationError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onOperationDone() {
        console.log("onOperationDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Account', link: '/account', icon: faUserCircle, isActive: false },
            { text: 'Avatar', link: '#', icon: faImage, isActive: true }
        ],
        mobileBackLinkItems: {
            link: "/account",
            text: "Back to Account",
            icon: faArrowLeft
        }
    }

    // --- Detail --- //

    function onSuccess(response) {
        console.log("onSuccess: Starting...");
        setCurrentUser(response);
    }

    function onError(apiErr) {
        console.log("onError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onDone() {
        console.log("onDone: Starting...");
        setFetching(false);
    }

    const onUnauthorized = () => {
        setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
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
            setFetching(true);
            getAccountDetailAPI(
                onSuccess,
                onError,
                onDone,
                onUnauthorized
            );
        }

        return () => { mounted = false; }
    }, [,]);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            {/* Page Title */}
            <h1 className="title is-2"><FontAwesomeIcon className="fas" icon={faUserCircle} />&nbsp;Account</h1>
            <h4 className="subtitle is-4"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</h4>
            <hr />

            {/* Page */}
            <div className="box">

                {/* Title + Options */}
                <p className="title is-4"><FontAwesomeIcon className="fas" icon={faImage} />&nbsp;Change Photo</p>

                <FormErrorBox errors={errors} />

                {/* <p className="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                {isFetching
                    ?
                    <PageLoadingContent displayMessage={"Submitting..."} />
                    :
                    <>
                        <div className="container">
                            <article className="message is-warning">
                                <div className="message-body">
                                    <strong>Warning:</strong> Submitting with new uploaded file will delete previous upload.
                                </div>
                            </article>

                            {selectedFile !== undefined && selectedFile !== null && selectedFile !== ""
                                ?
                                <>
                                    <article className="message is-success">
                                        <div className="message-body">
                                            <FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;File ready to upload.
                                        </div>
                                    </article>
                                </>
                                :
                                <>
                                    <b>File (Optional)</b>
                                    <br />
                                    <input name="file" type="file" onChange={onHandleFileChange} className="button is-medium" />
                                    <br />
                                    <br />
                                </>
                            }

                            <div className="columns pt-5">
                                <div className="column is-half">
                                    <Link to={`/account`} className="button is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Account</Link>
                                </div>
                                <div className="column is-half has-text-right">
                                    <button className="button is-medium is-success is-fullwidth-mobile" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                </div>
                            </div>

                        </div>
                    </>
                }
            </div>
        </Layout>
    );
}

export default AdminAccountAvatarOperation;
