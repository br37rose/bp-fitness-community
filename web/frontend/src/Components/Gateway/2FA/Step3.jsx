import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheckCircle, faArrowLeft, faArrowRight, faEnvelope, faKey, faTriangleExclamation, faArrowUpRightFromSquare } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';
import QRCode from "qrcode.react";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { postVertifyOTP } from "../../../API/gateway";
import { currentOTPResponseState, currentUserState } from "../../../AppState";
import FormInputField from "../../Reusable/FormInputField";
import { ROOT_ROLE_ID, ADMIN_ROLE_ID, TRAINER_ROLE_ID, MEMBER_ROLE_ID } from "../../../Constants/App";


function TwoFactorAuthenticationWizardStep3() {
    ////
    //// Global state.
    ////

    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);
    const [otpResponse] = useRecoilState(currentOTPResponseState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");
    const [verificationToken, setVerificationToken] = useState("");

    ////
    //// API.
    ////

    function onVerifyOPTSuccess(response){
        console.log("onVerifyOPTSuccess: Starting...");
        if (response !== undefined && response !== null && response !== "") {
            console.log("response: ", response);
            if (response.user !== undefined && response.user !== null && response.user !== "") {
                console.log("response.user: ", response.user);

                // Save our updated user account.
                setCurrentUser(response.user);

                switch (response.user.role) {
                    case ROOT_ROLE_ID:
                        setForceURL("/root/tenants");
                        break;
                    case ADMIN_ROLE_ID:
                        setForceURL("/admin/dashboard");
                        break;
                    case TRAINER_ROLE_ID:
                        setForceURL("/admin/dashboard");
                        break;
                    case MEMBER_ROLE_ID:
                        setForceURL("/dashboard");
                        break;
                    default:
                        setForceURL("/501");
                        break;
                }
            }
        }

    }

    function onVerifyOPTError(apiErr) {
        console.log("onVerifyOPTError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onVerifyOPTDone() {
        console.log("onVerifyOPTDone: Starting...");
    }

    ////
    //// Event handling.
    ////

    function onButtonClick(e) {
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
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
        }

        return () => mounted = false;
    }, []);

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
                                            <nav className="box has-background-success-light" >
                                                <p className="subtitle is-5">Step 3 of 3</p>
                                                <progress class="progress is-success" value="100" max="100">100%</progress>
                                            </nav>

                                            {/* Page */}
                                            <form>
                                                <h1 className="title is-2 has-text-centered">Setup Two-Factor Authentication</h1>
                                                <FormErrorBox errors={errors} />
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
                                                <br />
                                            </form>


                                            <nav class="level">
                                                <div class="level-left">
                                                    <div class="level-item">
                                                        <Link class="button is-link is-fullwidth-mobile" to="/login/2fa/step-2"><FontAwesomeIcon icon={faArrowLeft} />&nbsp;Back</Link>
                                                    </div>
                                                </div>
                                                <div class="level-right">
                                                    <div class="level-item">
                                                        <button type="button" class="button is-primary is-fullwidth-mobile" onClick={onButtonClick}>
                                                            <FontAwesomeIcon icon={faCheckCircle} />&nbsp;Subit and Verify
                                                        </button>
                                                    </div>
                                                </div>
                                            </nav>

                                        </div>
                                        {/* End box */}

                                        <div className="has-text-centered">
                                            <p>© 2024 Workery</p>
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

export default TwoFactorAuthenticationWizardStep3;
