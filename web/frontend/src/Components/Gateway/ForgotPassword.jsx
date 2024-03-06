import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowRight, faArrowLeft, faEnvelope, faKey, faTriangleExclamation, faCheckCircle } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../Reusable/FormErrorBox";
import useLocalStorage from "../../Hooks/useLocalStorage";
import { postForgotPasswordAPI } from "../../API/gateway";
import { topAlertMessageState, topAlertStatusState } from "../../AppState";


function ForgotPassword() {
    ////
    //// URL Parameters.
    ////

    const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
    const isUnauthorized = searchParams.get("unauthorized");

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [validation, setValidation] = useState({
        "email": false,
    });
    const [email, setEmail] = useState("");
    const [forceURL, setForceURL] = useState("");
    const [wasEmailSent, setWasEmailSent] = useState(false);

    ////
    //// API.
    ////

    function onForgotPasswordSuccess(){
        console.log("onForgotPasswordSuccess: Starting...");

        setTopAlertMessage("Email sent");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onOrganizationUpdateSuccess: Delayed for 2 seconds.");
            console.log("onOrganizationUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        setWasEmailSent(true);
    }

    function onForgotPasswordError(apiErr) {
        console.log("onForgotPasswordError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onForgotPasswordDone() {
        console.log("onForgotPasswordDone: Starting...");
    }

    ////
    //// Event handling.
    ////

    function onEmailChange(e) {
        setEmail(e.target.value);
        validation["email"]=false
        setValidation(validation);
        // setErrors(errors["email"]="");
    }

    function onButtonClick(e) {
        var newErrors = {};
        var newValidation = {};
        if (email === undefined || email === null || email === "") {
            newErrors["email"] = "value is missing";
        } else {
            newValidation["email"] = true
        }

        /// Save to state.
        setErrors(newErrors);
        setValidation(newValidation);

        if (Object.keys(newErrors).length > 0) {
            //
            // Handle errors.
            //

            console.log("failed validation");

            // window.scrollTo(0, 0);  // Start the page at the top of the page.

            // The following code will cause the screen to scroll to the top of
            // the page. Please see ``react-scroll`` for more information:
            // https://github.com/fisshy/react-scroll
            var scroll = Scroll.animateScroll;
            scroll.scrollToTop()
        } else {
            //
            // Submit to server.
            //

            console.log("successful validation, submitting to API server.");
            postForgotPasswordAPI(
                email,
                onForgotPasswordSuccess,
                onForgotPasswordError,
                onForgotPasswordDone
            );
        }
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

    if (forceURL !== "") {
        return <Navigate to={forceURL}  />
    }

    ////
    //// Component rendering.
    ////

    return (
        <>
            <div class="container column is-12">
                <div class="section">

                    <section class="hero is-fullheight">
                        <div class="hero-body">
                            <div class="container">
                                <div class="columns is-centered">
                                    <div class="column is-one-third-tablet">
                                        <div class="box is-rounded">
                                            {/* Start Logo */}
                                            <nav class="level">
                                                <div class="level-item has-text-centered">
                                                    <figure class='image'>
                                                        <img src='/static/logo.jpeg' style={{width:"256px"}} />
                                                    </figure>
                                                </div>
                                            </nav>
                                            {/* End Logo */}
                                            {!wasEmailSent
                                                ?
                                                <>
                                                    <form>
                                                        <h1 className="title is-4 has-text-centered">Forgot Password</h1>
                                                        <p className="pb-5 has-text-grey">Please enter your email and we will send you a password reset email.</p>
                                                        {isUnauthorized === "true" &&
                                                            <article class="message is-danger">
                                                              <div class="message-body"><FontAwesomeIcon className="fas" icon={faTriangleExclamation} />&nbsp;Your session has ended.<br/>Please login again</div>
                                                            </article>
                                                        }
                                                        <FormErrorBox errors={errors} />

                                                        <div class="field">
                                                            <label class="label is-small has-text-grey-light">Email</label>
                                                            <div class="control has-icons-left has-icons-right">
                                                                <input class={`input ${errors && errors.email && 'is-danger'} ${validation && validation.email && 'is-success'}`}
                                                                        type="email"
                                                                 placeholder="Email"
                                                                       value={email}
                                                                    onChange={onEmailChange}/>
                                                                <span class="icon is-small is-left">
                                                                    <FontAwesomeIcon className="fas" icon={faEnvelope} />
                                                                </span>
                                                            </div>
                                                            {errors && errors.email &&
                                                                <p class="help is-danger">{errors.email}</p>
                                                            }
                                                        </div>

                                                        <br />
                                                        <button class="button is-medium is-block is-fullwidth is-primary" type="button" onClick={onButtonClick} style={{backgroundColor:"#FF0000"}}>
                                                            ForgotPassword <FontAwesomeIcon icon={faArrowRight} />
                                                        </button>
                                                    </form>
                                                    <br />
                                                </>
                                                :
                                                <article class="message is-success">
                                                  <div class="message-body">
                                                    <h1 className="is-size-4"><FontAwesomeIcon icon={faCheckCircle} />&nbsp;<b>Email Sent</b></h1>
                                                    <p>The password reset email has been sent to your inbox. Please check and follow the instructions in the email.</p>
                                                    <br />
                                                    <p>Didn't receive the email? <a onClick={(e)=>onButtonClick()}>Click here</a> to resend again</p>
                                                  </div>
                                                </article>
                                            }
                                            <nav class="level">
                                                <div class="level-item has-text-centered">
                                                    <div>
                                                        <Link to="/login" className="is-size-7-tablet"><FontAwesomeIcon icon={faArrowLeft} />&nbsp;Back</Link>
                                                    </div>
                                                </div>
                                            </nav>
                                        </div>
                                        {/* End box */}

                                        <div className="has-text-centered">
                                            <p>Need help?</p>
                                            <p><Link to="Support@cpscapsule.com">Support@cpscapsule.com</Link></p>
                                            <p><a href="tel:+15199142685">(519) 914-2685</a></p>
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

export default ForgotPassword;
