import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faGauge,
  faArrowLeft,
  faFileInvoiceDollar,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getAccountDetailAPI } from "../../../API/Account";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import Footer from "../../Menu/Footer";

function PaymentProcessorSubscriptionSuccess() {
  ////
  ////  For debugging purposes only.
  ////

  console.log("REACT_APP_WWW_PROTOCOL:", process.env.REACT_APP_WWW_PROTOCOL);
  console.log("REACT_APP_WWW_DOMAIN:", process.env.REACT_APP_WWW_DOMAIN);
  console.log("REACT_APP_API_PROTOCOL:", process.env.REACT_APP_API_PROTOCOL);
  console.log("REACT_APP_API_DOMAIN:", process.env.REACT_APP_API_DOMAIN);

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

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
      getAccountDetailAPI(onSuccess, onError, onDone);
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
        <nav class="breadcrumb" aria-label="breadcrumbs">
          <ul>
            <li>
              <Link to="/dashboard" aria-current="page">
                <FontAwesomeIcon className="fas" icon={faGauge} />
                &nbsp;Dashboard
              </Link>
            </li>
            <li class="is-active">
              <Link to="/subscriptions" aria-current="page">
                <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                &nbsp;Subscriptions
              </Link>
            </li>
          </ul>
        </nav>

        {/* Mobile Breadcrumbs */}
        <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
          <ul>
            <li class=""><Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link></li>
          </ul>
        </nav>

        <nav class="box">
          <div class="columns">
            <div class="column">
              <h1 class="title is-4">
                <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                &nbsp;Subscriptions
              </h1>
            </div>
          </div>
          <div className="columns is-centered">
            <div className="column is-half">
              <div className="has-text-centered p-6">
                <h1 className="title is-3 custom-title has-text-success">
                  <b>Payment Successful!</b>
                </h1>
                <p className="has-text-dark custom-message">
                  Thank you for your payment! Your subscription is confirmed,
                  and you're all set to enjoy our services.
                </p>
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
                    to={`/subscriptions`}
                    className="button is-primary"
                    type="button"
                  >
                    <span className="icon">
                      <FontAwesomeIcon icon={faArrowLeft} />
                    </span>
                    <span>Back to Subscriptions</span>
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </nav>
        <Footer />
      </section>
    </div>
  );
}

export default PaymentProcessorSubscriptionSuccess;
