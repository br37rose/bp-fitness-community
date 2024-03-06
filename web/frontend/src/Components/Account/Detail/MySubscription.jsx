import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCartShopping, faXmark } from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { postSubscriptionCancelAPI } from "../../../API/subscription";
import { getMemberDetailAPI, postMemberCreateAPI } from "../../../API/member";

const MySubscription = ({
  currentUser,
  setTopAlertMessage,
  setTopAlertStatus,
  setCurrentUser,
}) => {
  ////
  //// Global state.
  ////

  ////
  //// Component states.
  ////
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [showModal, setShowModal] = useState(false);

  ////
  //// API.
  ////
  const onMySubscriptionCancelSuccess = (response) => {
    console.log("onMySubscriptionCancelSuccess: Starting...");
    setTopAlertStatus("success");
    setTopAlertMessage("MySubscription Cancelled");

    // Invokes member detail API to update the recoil state
    getMemberDetailUpdated();

    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
  };

  const onMySubscriptionCancelError = (apiErr) => {
    console.log("onMySubscriptionCancelError: Starting...");
    setTopAlertMessage("MySubscription Cancellation is unsuccessful");
    setTopAlertStatus("danger");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
    setErrors(apiErr);
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  };

  const onMySubscriptionCancelDone = () => {
    console.log("onMySubscriptionCancelDone: Starting...");
    setFetching(false);
  };

  const getMemberDetailUpdatedSuccess = (response) => {
    console.log("getMemberDetailUpdatedSuccess: Starting...");
    setCurrentUser(response);
  };

  const getMemberDetailUpdatedError = (apiErr) => {
    console.log("getMemberDetailUpdatedError: Starting...");
    setErrors(apiErr);

    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  };

  const getMemberDetailUpdatedDone = () => {
    console.log("onMySubscriptionCancelDone: Starting...");
    setFetching(false);
  };

  const handleCancelMySubscription = () => {
    postSubscriptionCancelAPI(
      "",
      onMySubscriptionCancelSuccess,
      onMySubscriptionCancelError,
      onMySubscriptionCancelDone
    );
  };

  const getMemberDetailUpdated = () => {
    getMemberDetailAPI(
      currentUser.id,
      getMemberDetailUpdatedSuccess,
      getMemberDetailUpdatedError,
      getMemberDetailUpdatedDone
    );
  };

  const handleCancelMySubscriptionClick = () => {
    setShowModal(!showModal);
  };

  const handleConfirmMySubscriptionCancel = () => {
    handleCancelMySubscription();
    setShowModal(false);
  };

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////
  useEffect(() => {
    let mounted = true;
    if (mounted) {
      window.scrollTo(0, 0);
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

  const { stripeSubscription } = currentUser;

  return (
    <>
      {showModal && (
        <div className={`modal ${showModal ? "is-active" : ""}`}>
          <div className="modal-background"></div>
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="is-size-4 modal-card-title">Are you sure?</p>
              <button
                className="delete"
                aria-label="cancel"
                onClick={handleCancelMySubscriptionClick}
              ></button>
            </header>
            <p className="is-size-5 modal-card-body">
              You are about to <b>cancel</b> your{" "}
              {stripeSubscription !== undefined &&
                stripeSubscription !== null &&
                stripeSubscription !== "" && (
                  <span>{`${stripeSubscription.interval}ly`}</span>
                )}{" "}
              subscription; This action cannot be reversed. Are you sure you
              would like to continue?
            </p>
            <footer className="modal-card-foot">
              <button
                className="button is-danger"
                onClick={handleConfirmMySubscriptionCancel}
              >
                Continue
              </button>
              <button
                className="button"
                onClick={handleCancelMySubscriptionClick}
              >
                Cancel
              </button>
            </footer>
          </div>
        </div>
      )}

      <div className="has-text-centered pb-5">
        <>
          {stripeSubscription !== null && stripeSubscription !== "" && (
            <>
              <p className="title is-size-5 has-text-info pb-1">
                Current Plan: <span>{`${stripeSubscription.interval}`}</span>
              </p>
              <p className="is-size-6 has-text-secondary pb-1">
                You are subscribed to our{" "}
                {`${stripeSubscription.interval === 'year' ? "annual" : "monthly"}`} premium plan.
              </p>
            </>
          )}
          <div className="columns is-centered mt-4">
            <div className="column is-narrow">
              <Link to="/invoices" className="button is-dark">
                <FontAwesomeIcon className="fas" icon={faCartShopping} />
                &nbsp;View Past Invoices
              </Link>
            </div>
            <div className="column is-narrow">
              <button
                className="button is-danger"
                onClick={handleCancelMySubscriptionClick}
              >
                <FontAwesomeIcon className="fas" icon={faXmark} />
                &nbsp;Cancel My Subscription
              </button>
            </div>
          </div>
        </>
      </div>
    </>
  );
};

export default MySubscription;
