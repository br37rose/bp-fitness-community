import React, { useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faUserCircle,
  faPager,
  faArrowRight,
  faGauge,
  faArrowLeft,
  faFileInvoiceDollar,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getAccountDetailAPI } from "../../../API/Account";
import { getCompleteStripeSubscriptionCheckoutSessionAPI } from "../../../API/PaymentProcessor";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import Footer from "../../Menu/Footer";


function PaymentProcessorPurchaseSuccess() {
    ////
    ////  For debugging purposes only.
    ////

    console.log("REACT_APP_WWW_PROTOCOL:", process.env.REACT_APP_WWW_PROTOCOL);
    console.log("REACT_APP_WWW_DOMAIN:", process.env.REACT_APP_WWW_DOMAIN);
    console.log("REACT_APP_API_PROTOCOL:", process.env.REACT_APP_API_PROTOCOL);
    console.log("REACT_APP_API_DOMAIN:", process.env.REACT_APP_API_DOMAIN);

    ////
    //// URL Parameters.
    ////

    const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
    const sessionID = searchParams.get("session_id");

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [checkoutSession, setCheckoutSession] = useState(null);

    ////
    //// API.
    ////

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

    function onCheckoutSessionSuccess(response) {
        console.log("onCheckoutSessionSuccess: Starting...");
        setCheckoutSession(response);
    }

    function onCheckoutSessionError(apiErr) {
        console.log("onCheckoutSessionError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onCheckoutSessionDone() {
        console.log("onCheckoutSessionDone: Starting...");
        setFetching(true);
        setErrors({});
        getAccountDetailAPI(onSuccess, onError, onDone);
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
        window.scrollTo(0, 0); // Start the page at the top of the page.
        setFetching(true);
        setErrors({});
        getCompleteStripeSubscriptionCheckoutSessionAPI(
            sessionID,
            onCheckoutSessionSuccess,
            onCheckoutSessionError,
            onCheckoutSessionDone
        );
    }

    return () => {
      mounted = false;
    };
    }, []);

    ////
    //// Component rendering.
    ////

    return (
        <div className="container">
          <section className="section">
            {/* Desktop Breadcrumbs */}
            <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
              <ul>
                <li className=""><Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                <li className=""><Link to="/account/more" aria-current="page"><FontAwesomeIcon className="fas" icon={faUserCircle} />&nbsp;Account (More)</Link></li>
                <li className="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPager} />&nbsp;Subscribe - Purchase Successful</Link></li>
              </ul>
            </nav>

            {/* Mobile Breadcrumbs */}
            <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
              <ul>
                <li class="">
                    <Link to="/account/more" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Account (More)</Link>
                </li>
              </ul>
            </nav>

            <nav class="box">
              <div class="columns">
                <div class="column">
                  <h1 class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faPager} />
                    &nbsp;Subscribe - Purchase Successful
                  </h1>
                </div>
              </div>
              <div className="columns is-centered">
                <div className="column is-half">
                  <div className="has-text-centered p-6">
                    <h1 className="title is-3 custom-title has-text-success">
                      <b>Payment Successful!</b>
                    </h1>
                    {checkoutSession && <p className="has-text-dark custom-message">
                      Thank you for your payment! Your order of <strong>{checkoutSession.name}</strong> at <i>${checkoutSession.price}&nbsp;{checkoutSession.priceCurrency}</i> is confirmed,
                      and you're all set to enjoy our services.
                    </p>}
                    <div className="custom-image-container">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="200"
                        height="200"
                        viewBox="0 0 200 200"
                      >
                        <circle cx="100" cy="100" r="90" fill="#f2f2f2" />
                        <circle cx="100" cy="100" r="50" fill="#66bb6a" />
                        <path
                          id="tickPath"
                          d="M72 100L90 120L130 80"
                          stroke="#fff"
                          strokeWidth="8"
                          fill="none"
                          strokeLinecap="round"
                          style={{
                            strokeDasharray: 300,
                            strokeDashoffset: 300,
                            animation: "draw 2s ease-out forwards",
                          }}
                        >
                          <animate
                            attributeName="stroke-dashoffset"
                            from="300"
                            to="0"
                            dur="1s"
                            fill="freeze"
                          />
                        </path>
                        <style>
                          {`
                                @keyframes draw {
                                    to {
                                    stroke-dashoffset: 0;
                                    }
                                }
                            `}
                        </style>
                      </svg>
                    </div>

                    <div className="buttons is-centered mt-4">
                      <Link
                        to={`/account/more`}
                        className="button is-primary"
                        type="button">
                        <span>Back to My Account</span>
                        <span className="icon">
                          <FontAwesomeIcon icon={faArrowRight} />
                        </span>
                      </Link>
                    </div>

{/*
                    or

                    <div className="buttons is-centered mt-4">
                      <Link
                        to={`/purchases`}
                        className="button is-black"
                        type="button">
                        <span>Go to Purchases</span>
                        <span className="icon">
                          <FontAwesomeIcon icon={faArrowRight} />
                        </span>
                      </Link>
                    </div>
*/}

                  </div>
                </div>
              </div>
            </nav>
            <Footer />
          </section>
        </div>
    );
}

export default PaymentProcessorPurchaseSuccess;
