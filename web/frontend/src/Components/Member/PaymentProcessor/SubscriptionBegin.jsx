import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faFileInvoiceDollar,
  faGauge,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import Scroll from "react-scroll";

import { postCreateStripeSubscriptionCheckoutSessionAPI } from "../../../API/PaymentProcessor";
import { getMemberDetailAPI } from "../../../API/member";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import Footer from "../../Menu/Footer";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import MySubscription from "./MySubscription";
import PricingTable from "./Pricing";

function PaymentProcessorBeginSubscription() {
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
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [datum, setDatum] = useState({});
  const [tabIndex, setTabIndex] = useState(1);

  ////
  //// API.
  ////

  function onSuccess(response) {
    console.log("onSuccess: Starting...");
    console.log("onSuccess: Redirecting at", response.checkoutUrl);

    // Force the user's browser to a different domain address.
    window.location.href = response.checkoutUrl;
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

  function onMemberDetailSuccess(response) {
    console.log("onMemberDetailSuccess: Starting...");

    setCurrentUser(response)
    // console.log(response)
  }

  function onMemberDetailError(apiErr) {
    console.log("onMemberDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onMemberDetailDone() {
    console.log("onMemberDetailDone: Starting...");
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const onClick = (priceID) => {
    // action={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}/api/v1/stripe/create-subscription-checkout-session`}
    setFetching(true);
    postCreateStripeSubscriptionCheckoutSessionAPI(
      priceID,
      onSuccess,
      onError,
      onDone
    );
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      getMemberDetailAPI(
        currentUser.id,
        onMemberDetailSuccess,
        onMemberDetailError,
        onMemberDetailDone
      )
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="container">
      <section className="section">
        {/* Desktop Breadcrumbs */}
        <nav className="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
          <ul>
            <li>
              <Link to="/dashboard" aria-current="page">
                <FontAwesomeIcon className="fas" icon={faGauge} />
                &nbsp;Dashboard
              </Link>
            </li>
            <li className="is-active">
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

        {/*
            RENDER THE FOLLOWING GUI BASED ON WHETHER THE USER ALREADY IS A SUBSCRIBER OR NOT.
            1. IF USER IS A SUBSCRIBER THEN PROVIDE GUI WHICH HELPS THEM SEE WHAT PLAN THEY ARE AND TO CANCEL IF NECESSARY
            2. IF USER DID NOT SUBSCRIBER THEN PROVIDE GUI WHICH HELPS THEM BECOME A SUBSCRIBER.
        */}
        {currentUser.isSubscriber ? (
          <nav class="box">
            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                  &nbsp;My Subscription
                  <br />
                  <br />
                  <MySubscription
                    currentUser={currentUser}
                    setTopAlertMessage={setTopAlertMessage}
                    setTopAlertStatus={setTopAlertStatus}
                    setCurrentUser={setCurrentUser}
                  />
                </h1>
              </div>
            </div>

            <FormErrorBox errors={errors} />
          </nav>
        ) : (
          <nav class="box">
            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                  &nbsp;Pick a Subscription
                </h1>
              </div>
            </div>

            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              // <div className="columns is-vcentered is-multiline">
              //   <div className="column is-6">
              //     <div
              //       className="card has-background-black has-background-gradient is-relative has-text-white has-text-centered p-6"
              //       style={{
              //         backgroundImage: "url('/static/membership-image-1.jpeg')",
              //         backgroundSize: "cover",
              //       }}
              //     >
              //       <form>
              //         <input
              //           type="hidden"
              //           id="priceId"
              //           name="priceId"
              //           value={
              //             process.env
              //               .REACT_APP_WWW_STRIPE_ANNUAL_MEMBERSHIP_PRICE_ID
              //           }
              //         />
              //         <h2 className="title has-text-white is-4 mt-4">
              //           Annual Membership
              //         </h2>
              //         <p className="subtitle is-6 has-text-white">
              //           Get fit with our annual membership plan.
              //         </p>
              //         <p className="is-size-4 has-text-primary">
              //           <b>C$1750/year</b>
              //         </p>
              //         <div className="has-text-centered mt-4">
              //           <button
              //             id="pro-plan-btn"
              //             type="button"
              //             onClick={(e) =>
              //               onClick(
              //                 process.env
              //                   .REACT_APP_WWW_STRIPE_ANNUAL_MEMBERSHIP_PRICE_ID
              //               )
              //             }
              //             className="button is-primary"
              //           >
              //             Get Now
              //           </button>
              //         </div>
              //       </form>
              //     </div>
              //   </div>
              //   <div className="column is-6">
              //     <div
              //       className="card has-background-black has-background-gradient is-relative has-text-white has-text-centered p-6"
              //       style={{
              //         backgroundImage: "url('/static/membership-image-2.jpeg')",
              //         backgroundSize: "cover",
              //       }}
              //     >
              //       <form>
              //         <input
              //           type="hidden"
              //           id="priceId"
              //           name="priceId"
              //           value={
              //             process.env
              //               .REACT_APP_WWW_STRIPE_MONTHLY_MEMBERSHIP_PRICE_ID
              //           }
              //         />
              //         <h2 className="title has-text-white is-4 mt-4">
              //           Monthly Membership
              //         </h2>
              //         <p className="subtitle is-6 has-text-white">
              //           Get started with our flexible monthly membership plan.
              //         </p>
              //         <p className="is-size-4 has-text-primary">
              //           <b>C$160/month</b>
              //         </p>
              //         <div className="has-text-centered mt-4">
              //           <button
              //             id="pro-plan-btn"
              //             type="button"
              //             onClick={(e) =>
              //               onClick(
              //                 process.env
              //                   .REACT_APP_WWW_STRIPE_MONTHLY_MEMBERSHIP_PRICE_ID
              //               )
              //             }
              //             className="button is-primary"
              //           >
              //             Get Now
              //           </button>
              //         </div>
              //       </form>
              //     </div>
              //   </div>
              // </div>
              <PricingTable />
            )}
          </nav>
        )}

        <Footer />
      </section>
    </div>
  );
}

export default PaymentProcessorBeginSubscription;
