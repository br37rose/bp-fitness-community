import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faUserCircle,
  faGauge,
  faArrowLeft,
  faPager,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../AppState";
import Footer from "../../Menu/Footer";

function PaymentProcessorOneTimePurchaseCanceled() {
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

  ////
  //// API.
  ////

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
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
            <li className="is-active">
                <Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPager} />&nbsp;Subscribe - Purchase Cancelled</Link>
            </li>
          </ul>
        </nav>


        {/* Mobile Breadcrumbs */}
        <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
          <ul>
            <li class=""><Link to="/subscriptions" aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Classes &amp; Events</Link></li>
          </ul>
        </nav>

        <nav class="box">
          <div class="columns">
            <div class="column">
              <h1 class="title is-4">
                <FontAwesomeIcon className="fas" icon={faPager} />
                &nbsp;Subscribe - Purchase Cancelled
              </h1>
            </div>
          </div>
          <div className="columns is-centered">
            <div className="column is-half">
              <div className="has-text-centered p-6">
                <h1 className="title is-3 custom-title has-text-danger">
                  <b>Oops!</b>
                </h1>
                <p className="subtitle is-5 custom-subtitle">
                  Your Purchase Has Been Cancelled
                </p>
                <p className="has-text-dark custom-message">
                  Don't worry, we understand that plans can change. Our team is
                  here to help you whenever you're ready to book your next
                  session.
                </p>

                <div className="custom-image-container">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="200"
                    height="200"
                    viewBox="0 0 200 200"
                  >
                    <circle cx="100" cy="100" r="90" fill="#f2f2f2" />
                    <circle cx="100" cy="100" r="50" fill="#ff9999" />
                    <path
                      id="cancelPath1"
                      d="M70 70L130 130"
                      stroke="#fff"
                      strokeWidth="8"
                      fill="none"
                      strokeLinecap="round"
                      style={{
                        strokeDasharray: 200,
                        strokeDashoffset: 200,
                        animation: "draw1 1s ease-out forwards",
                      }}
                    >
                      <animate
                        attributeName="stroke-dashoffset"
                        from="200"
                        to="0"
                        dur="1s"
                        fill="freeze"
                      />
                    </path>
                    <path
                      id="cancelPath2"
                      d="M70 130L130 70"
                      stroke="#fff"
                      strokeWidth="8"
                      fill="none"
                      strokeLinecap="round"
                      style={{
                        strokeDasharray: 200,
                        strokeDashoffset: 200,
                        animation: "draw2 1s ease-out forwards",
                      }}
                    >
                      <animate
                        attributeName="stroke-dashoffset"
                        from="200"
                        to="0"
                        dur="2s"
                        fill="freeze"
                      />
                    </path>
                    <style>
                      {`
                            @keyframes draw1 {
                                to {
                                stroke-dashoffset: 0;
                                }
                            }
                            @keyframes draw2 {
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
                    type="button"
                  >
                    <span className="icon">
                      <FontAwesomeIcon icon={faArrowLeft} />
                    </span>
                    <span>Back to My Account</span>
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

export default PaymentProcessorOneTimePurchaseCanceled;
