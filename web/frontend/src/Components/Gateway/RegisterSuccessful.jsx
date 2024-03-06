import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowRight, faKey, faEnvelope } from '@fortawesome/free-solid-svg-icons'
import FormErrorBox from "../Reusable/FormErrorBox";
import useLocalStorage from "../../Hooks/useLocalStorage";
import { postLoginAPI } from "../../API/gateway";

function RegisterSuccessful() {
    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");

    ////
    //// API.
    ////

    // Do nothing.

    ////
    //// Event handling.
    ////

    // Do nothing.

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
                                <div class="columns is-centered p-7">
                                    <div class="is-rounded column is-two-third-tablet">
                                        <article class="message is-primary">
                                          <div class="message-body">
                                                <h1 className="title is-4 has-text-centered has-text-success"><FontAwesomeIcon className="fas" icon={faEnvelope} />&nbsp;Email Sent</h1>
                                                <FormErrorBox errors={errors} />
                                                <p>Thank you for registering - an <b>activation email</b> has bee sent to you. Please be sure to check your social, promotions and spam folders if it does not arrive within 5 minutes.</p>
                                                <p>
                                                    <br />
                                                    <Link to="/"><b>Back to index&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></b></Link>
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
        </>
    );
}

export default RegisterSuccessful;
