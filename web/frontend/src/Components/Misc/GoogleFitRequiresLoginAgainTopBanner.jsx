import React, {useState, useEffect} from "react";
import { Link } from "react-router-dom";
import { useRecoilState } from "recoil";
import { useLocation } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheckCircle, faInfoCircle, faArrowUpRightFromSquare, faCircleExclamation } from '@fortawesome/free-solid-svg-icons'
import Scroll from "react-scroll";

import { currentUserState } from "../../AppState";
import {getAccountDetailAPI} from "../../API/Account";
import {
	getGoogleFitRegistrationURLAPI,
	postFitBitAppCreateSimulatorAPI,
} from "../../API/Wearable";


/**
  The purpose of this component is to intercept anonymous users at our
  application URLs which require authorization.
 */
function GoogleFitRequiresLoginAgainTopBanner() {
  ////
  //// Global state.
  ////
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");

  ////
  //// Event handling.
  ////

  const onRegisterClick = (e) => {
      e.preventDefault();
      console.log("onRegisterClick: Clicked");
      setFetching(true);
      setErrors({});
      getGoogleFitRegistrationURLAPI(
          onRegistrationSuccess,
          onRegistrationError,
          onRegistrationDone
      );
  };

  ////
  //// API.
  ////

  // --- Register --- //

  function onRegistrationSuccess(response) {
      console.log("onRegistrationSuccess: Starting...");
      window.location = response.url;
  }

  function onRegistrationError(apiErr) {
      console.log("onRegistrationError: Starting...");
      setErrors(apiErr);

      // The following code will cause the screen to scroll to the top of
      // the page. Please see ``react-scroll`` for more information:
      // https://github.com/fisshy/react-scroll
      var scroll = Scroll.animateScroll;
      scroll.scrollToTop();
  }

  function onRegistrationDone() {
      console.log("onRegistrationDone: Starting...");
      setFetching(false);
  }


  ////
  //// Logic
  ////

  // Get the current location and if we are at specific URL paths then we
  // will not render this component.
  const ignorePathsArr = [
    "/",
    "/register",
    "/register-step-1",
    "/register-step-2",
    "/register-successful",
    "/index",
    "/terms-of-service",
    "/privacy-policy",
    "/login",
    "/login/2fa",
    "/login/2fa/step-1",
    "/login/2fa/step-2",
    "/login/2fa/step-3",
    "/logout",
    "/verify",
    "/forgot-password",
    "/password-reset",
    "/terms",
    "/privacy",
  ];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]); // For debugging purposes only.
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }

  if (currentUser === null) {
    console.log("No current user detected, skipping...");
    return null;
  } else {
    if (currentUser.primaryHealthTrackingDeviceRequiresLoginAgain) {
       console.log("Current user needs to log in again into Google Fit, displaying top banner now...");
        return (
            <div class="notification is-danger is-light">
                <FontAwesomeIcon className="fas" icon={faCircleExclamation} />&nbsp;<b>Authentication Required</b>&nbsp;Your Google Fit fitness tracker has been disconnected by Google and requires you to login
                again.{" "}
                <b>
                    <Link onClick={onRegisterClick}>
                        Click here&nbsp;
                        <FontAwesomeIcon
                            className="mdi"
                            icon={faArrowUpRightFromSquare}
                        />{" "}
                    </Link>
                </b>{" "}
                to login again and meet the requirements of Google.
            </div>
        );
    }
    return null;
  }
}

export default GoogleFitRequiresLoginAgainTopBanner;
