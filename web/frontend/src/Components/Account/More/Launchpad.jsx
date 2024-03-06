import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCreditCard, faKey, faUserCircle, faBuildingUser, faImage, faPaperclip, faAddressCard, faSquarePhone, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faHardHat, faGauge, faPencil, faUsers, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faBuilding, faEllipsis, faArchive, faBoxOpen, faTrashCan, faHomeUser } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getAccountDetailAPI } from "../../../API/Account";
import AlertBanner from "../../Reusable/EveryPage/AlertBanner";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import BubbleLink from "../../Reusable/EveryPage/BubbleLink";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";
import Layout from "../../Menu/Layout";


function AccountMoreLaunchpad() {
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
    const [currentUser, setCurrentUser] = useState({});

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
        <div>
            <div className="columns is-vcentered is-multiline">
                {currentUser.status === 1 && <div className="column">
                    <BubbleLink
                        title={`Photo`}
                        subtitle={`Upload a photo of the currentUser`}
                        faIcon={faImage}
                        url={`/account/more/avatar`}
                        bgColour={`has-background-danger-dark`}
                    />
                </div>}
                {currentUser.status === 1 && currentUser.isSubscriber === false && <div className="column">
                    <BubbleLink
                        title={`Subscribe`}
                        subtitle={`Get a membership with us`}
                        faIcon={faCreditCard}
                        url={`/account/more/subscribe`}
                        bgColour={`has-background-success-dark`}
                    />
                </div>}
            </div>


            {/* Bottom Navigation */}
            < div className="columns pt-5" >
                <div className="column is-half">
                    <Link className="button is-medium is-fullwidth-mobile" to={`/dashboard`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                </div>
                <div className="column is-half has-text-right">

                </div>
            </div >
        </div >
    );
}

export default AccountMoreLaunchpad;
