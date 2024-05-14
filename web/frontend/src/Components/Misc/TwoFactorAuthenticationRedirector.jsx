import React from "react";
import { Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { useLocation } from "react-router-dom";

import { currentUserState } from "../../AppState";

/**
  The purpose of this component is to intercept authenticated users whom have
  two-factor authentication (2FA) enabled and they did not validate their 2FA
  code after logging on.
 */
function TwoFactorAuthenticationRedirector() {
  ////
  //// Global state.
  ////
  const [currentUser] = useRecoilState(currentUserState);

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

  // console.log("TwoFactorAuthenticationRedirector | currentUser:", currentUser)

  if (currentUser === null) {
    console.log("No current user detected, redirecting back to login page.");
    return <Navigate to={`/login?unauthorized=true`} />;
  } else {
    // Only enforce 2FA if the user enabled it.
    if (currentUser.otpEnabled === true) {
      if (currentUser.otpVerified === false) {
        // If the user enabled it but has not verified their OTP then redirect the user to get it setup.
        return <Navigate to={`/login/2fa/step-1`} />;
      } else if (currentUser.otpValidated === false) {
        // If the user did not validate their 2FA code after login then
        // force redirect the user to the 2FA login.
        return <Navigate to={`/login/2fa`} />;
      }
    }
    return null;
  }
}

export default TwoFactorAuthenticationRedirector;
