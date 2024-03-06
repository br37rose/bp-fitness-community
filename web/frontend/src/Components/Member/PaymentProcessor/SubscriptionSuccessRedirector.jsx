import React, { useState, useEffect } from "react";
import { Link, useSearchParams, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faFileInvoice,
  faFileInvoiceDollar,
} from "@fortawesome/free-solid-svg-icons";

import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { getCompleteStripeSubscriptionCheckoutSessionAPI } from "../../../API/PaymentProcessor";
import Footer from "../../Menu/Footer";

// The purpose of this component is to extract the `session_id` from the URL
// parameter (which was set by Stripe, Inc.) and then submit to our API endpoint
// to complete the session and then finally redirect the user.

function PaymentProcessorSubscriptionSuccessRedirector() {
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
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  function onSuccess(response) {
    console.log("onSuccess: Starting...");
    setForceURL("/subscription/completed");
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
      getCompleteStripeSubscriptionCheckoutSessionAPI(
        sessionID,
        onSuccess,
        onError,
        onDone
      );
    }

    return () => {
      mounted = false;
    };
  }, [sessionID]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li class="is-active">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                  &nbsp;Subscription
                </Link>
              </li>
            </ul>
          </nav>
          <nav class="box">
            <div class="columns">
              <div class="column">
                <PageLoadingContent displayMessage={"Loading..."} />
              </div>
            </div>
          </nav>
          <Footer />
        </section>
      </div>
    </>
  );
}

export default PaymentProcessorSubscriptionSuccessRedirector;
