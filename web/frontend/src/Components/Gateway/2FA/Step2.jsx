import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheckCircle, faArrowLeft, faArrowRight, faEnvelope, faKey, faTriangleExclamation, faArrowUpRightFromSquare } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import QRCode from "qrcode.react";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { postGenerateOTP, postGenerateOTPAndQRCodeImage } from "../../../API/gateway";
import { currentOTPResponseState, currentUserState } from "../../../AppState";
import { ROOT_ROLE_ID, ADMIN_ROLE_ID, TRAINER_ROLE_ID, MEMBER_ROLE_ID } from "../../../Constants/App";


function TwoFactorAuthenticationWizardStep2() {
    ////
    //// Global state.
    ////

    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);
    const [otpResponse, setOtpResponse] = useRecoilState(currentOTPResponseState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");

    ////
    //// API.
    ////

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
    }

    ////
    //// Event handling.
    ////


    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            if (otpResponse === undefined || otpResponse === null || otpResponse === "") {
                postGenerateOTP(
                    onGenerateOPTSuccess,
                    onGenerateOPTError,
                    onGenerateOPTDone
                );
            }
        }

        return () => mounted = false;
    }, [otpResponse]);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL}  />
    }

    return (
        <>
            <div className="container column is-12">
                <div className="section">


                    <section className="hero is-fullheight">
                        <div className="hero-body">
                            <div className="container">
                                <div className="columns is-centered">
                                    <div className="column is-half-tablet">
                                        <div className="box is-rounded">
                                            {/* Progress Wizard */}
                                            <nav className="box has-background-light" >
                                                <p className="subtitle is-5">Step 2 of 3</p>
                                                <progress class="progress is-success" value="66" max="100">66%</progress>
                                            </nav>

                                            {/* Page */}
                                            <form>
                                                <h1 className="title is-2 has-text-centered">Setup Two-Factor Authentication</h1>
                                                <FormErrorBox errors={errors} />
                                                <p class="has-text-grey">With your 2FA application open, please scan the following QR code with your device and click next when ready.</p>
                                                <p>&nbsp;</p>
                                                <div className="columns is-centered is-hidden-mobile">
                                                    <div class="column is-half">
                                                        <figure class="image">
                                                            {otpResponse &&
                                                                <QRCode value={otpResponse.optAuthURL} size={250} className="" />}
                                                            Scan with your app
                                                        </figure>
                                                    </div>
                                                </div>

                                            </form>
                                            <br />

                                            <nav class="level">
                                                <div class="level-left">
                                                    <div class="level-item">
                                                        <Link class="button is-link is-fullwidth-mobile" to="/login/2fa/step-1"><FontAwesomeIcon icon={faArrowLeft} />&nbsp;Back</Link>
                                                    </div>
                                                </div>
                                                <div class="level-right">
                                                    <div class="level-item">
                                                        <Link type="button" class="button is-primary is-fullwidth-mobile" to="/login/2fa/step-3">
                                                            Next&nbsp;<FontAwesomeIcon icon={faArrowRight} />
                                                        </Link>
                                                    </div>
                                                </div>
                                            </nav>

                                        </div>
                                        {/* End box */}

                                        <div className="has-text-centered">
                                            <p>© 2024 BP8 Fitness Community</p>
                                        </div>
                                        {/* End suppoert text. */}

                                    </div>
                                    {/* End Column */}
                                </div>
                            </div>
                            {/* End container */}
                        </div>
                        {/* End hero-body */}
                    </section>

                </div>
            </div>
        </>
    );
}

export default TwoFactorAuthenticationWizardStep2;
