import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTasks, faTachometer, faPlus, faTimesCircle, faCheckCircle, faUserCircle, faGauge, faPencil, faDumbbell, faIdCard, faAddressBook, faMessage, faChartPie, faCogs, faEye } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getVideoCategoryDetailAPI, postVideoCategoryCreateAPI } from "../../../API/VideoCategory";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";


function AdminVideoCategoryAdd() {
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
        postVideoCategoryCreateAPI(decamelizedData, onAdminVideoCategoryAddSuccess, onAdminVideoCategoryAddError, onAdminVideoCategoryAddDone);
    }

    function onAdminVideoCategoryAddSuccess(response){
        // For debugging purposes only.
        console.log("onAdminVideoCategoryAddSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Video category created");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAdminVideoCategoryAddSuccess: Delayed for 2 seconds.");
            console.log("onAdminVideoCategoryAddSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/admin/video-category/"+response.id);
    }

    function onAdminVideoCategoryAddError(apiErr) {
        console.log("onAdminVideoCategoryAddError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminVideoCategoryAddError: Delayed for 2 seconds.");
            console.log("onAdminVideoCategoryAddError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAdminVideoCategoryAddDone() {
        console.log("onAdminVideoCategoryAddDone: Starting...");
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
                            <li class=""><Link to="/admin/video-categories" aria-current="page"><FontAwesomeIcon className="fas" icon={faDumbbell} />&nbsp;Video Categories</Link></li>
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
                                    <Link class="button is-medium is-success" to={`/admin/video-categories`}>Yes</Link>
                                    <button class="button is-medium" onClick={(e)=>setShowCancelWarning(false)}>No</button>
                                </footer>
                            </div>
                        </div>

                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New Video Category</p>
                        <FormErrorBox errors={errors} />

                        {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                        {isFetching && <PageLoadingContent displayMessage={"Please wait..."} />}

                        <div class="container">

                            <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</p>
                            <hr />

                            <FormInputField
                                label="Name"
                                name="name"
                                placeholder="Text input"
                                value={name}
                                errorText={errors && errors.name}
                                helpText=""
                                onChange={(e)=>setName(e.target.value)}
                                isRequired={true}
                                maxWidth="380px"
                            />

                            <FormInputField
                                label="No #"
                                name="no"
                                type="number"
                                placeholder="#"
                                value={no}
                                errorText={errors && errors.no}
                                helpText=""
                                onChange={(e)=>setNo(parseInt(e.target.value))}
                                isRequired={true}
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

export default AdminVideoCategoryAdd;
