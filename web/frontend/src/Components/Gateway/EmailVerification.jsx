import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useSearchParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { postEmailVerificationAPI } from "../../API/gateway";
import FormErrorBox from "../Reusable/FormErrorBox";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowRight, faEnvelope } from '@fortawesome/free-solid-svg-icons'


function EmailVerification() {

    ////
    //// URL Parameters
    ////

    let [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [version, setEmailVerification] = useState("");

    ////
    //// API.
    ////

    function onEmailVerificationSuccess(response){
        console.log("onEmailVerificationSuccess: Starting...");
        setEmailVerification(response);
    }

    function onEmailVerificationError(apiErr) {
        console.log("onEmailVerificationError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onEmailVerificationDone() {
        console.log("onEmailVerificationDone: Starting...");
    }

    ////
    //// Event handling.
    ////

    // (Do nothing)

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            postEmailVerificationAPI(
                searchParams.get("q"), // Extract the verification code from the query parameter in the URL.
                onEmailVerificationSuccess,
                onEmailVerificationError,
                onEmailVerificationDone
            );
        }

        return () => mounted = false;
    }, []);

    ////
    //// Component rendering.
    ////

    return (
        <div class="container column is-4">
            <div class="section">

            <section class="hero is-fullheight">
                <div class="hero-body">
                    <div class="container">
                        <div class="columns is-centered p-7">
                            <div class="is-rounded column is-two-third-tablet">
                                <article class="message is-primary">
                                  <div class="message-body">

                                        <h1 className="title is-4 has-text-centered has-text-success"><FontAwesomeIcon className="fas" icon={faEnvelope} />&nbsp;Email was Verified</h1>
                                        <FormErrorBox errors={errors} />
                                        <p>Thank you for verifying. You may now log into your account now view the login page.</p>
                                        <p>
                                            <br />
                                            <Link to="/login"><b>Go to login&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></b></Link>
                                        </p>
                                  </div>
                                </article>
                            </div>
                            {/* End box */}
                        </div>
                    </div>
                    {/* End container */}
                </div>
                {/* End hero-body */}
            </section>

            </div>
        </div>
      );
}

export default EmailVerification;
