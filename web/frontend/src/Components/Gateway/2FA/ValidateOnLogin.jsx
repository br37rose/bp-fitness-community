import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCheckCircle,
  faArrowLeft,
  faArrowRight,
  faEnvelope,
  faKey,
  faTriangleExclamation,
  faArrowUpRightFromSquare,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import QRCode from "qrcode.react";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { postValidateOTP } from "../../../API/gateway";
import { currentOTPResponseState, currentUserState } from "../../../AppState";
import FormInputField from "../../Reusable/FormInputField";
import {
  ROOT_ROLE_ID,
  ADMIN_ROLE_ID,
  TRAINER_ROLE_ID,
  MEMBER_ROLE_ID,
} from "../../../Constants/App";

function TwoFactorAuthenticationValidateOnLogin() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const paramToken = searchParams.get("token");

  ////
  //// Global state.
  ////

  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);
  const [otpResponse] = useRecoilState(currentOTPResponseState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");
  const [token, setToken] = useState("");
  const [submittedParamToken, setSubmittedParamToken] = useState(false);

  ////
  //// API.
  ////

  function onVerifyOPTSuccess(response) {
    console.log("onVerifyOPTSuccess: Starting...");
    if (response !== undefined && response !== null && response !== "") {
      console.log("response: ", response);
      if (
        response.user !== undefined &&
        response.user !== null &&
        response.user !== ""
      ) {
        console.log("response.user: ", response.user);

        // Save our updated user account.
        setCurrentUser(response.user);

        switch (response.user.role) {
          case ROOT_ROLE_ID:
            setForceURL("/root/tenants");
            break;
          case ADMIN_ROLE_ID:
            setForceURL("/admin/dashboard");
            break;
          case TRAINER_ROLE_ID:
            setForceURL("/admin/dashboard");
            break;
          case MEMBER_ROLE_ID:
            console.log("memeber onboarding after login otp");
            setForceURL("/onboarding");
            break;
          default:
            setForceURL("/501");
            break;
        }
      }
    }
  }

  function onVerifyOPTError(apiErr) {
    console.log("onVerifyOPTError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onVerifyOPTDone() {
    console.log("onVerifyOPTDone: Starting...");
  }

  ////
  //// Event handling.
  ////

  function onButtonClick(e) {
    // Remove whitespace characters from token
    const cleanedToken = token.replace(/\s/g, "");

    const payload = {
      token: cleanedToken,
    };
    postValidateOTP(
      payload,
      onVerifyOPTSuccess,
      onVerifyOPTError,
      onVerifyOPTDone
    );
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      // DEVELOPERS NOTE:
      // It appears that `Apple Verification` service submits a `token` url
      // parameter to the page with the uniquely generated 2FA code; as a result,
      // the following code will check to see if this `token` url parameter
      // exists and whether it was submitted or not and if it wasn't submitted
      // then we submit for OTP verification and proceed.
      if (
        submittedParamToken === false &&
        paramToken !== undefined &&
        paramToken !== null &&
        paramToken !== ""
      ) {
        const payload = {
          token: paramToken,
        };
        postValidateOTP(
          payload,
          onVerifyOPTSuccess,
          onVerifyOPTError,
          onVerifyOPTDone
        );
        setSubmittedParamToken(true);
        setToken(paramToken);
      }
    }

    return () => (mounted = false);
  }, [paramToken, submittedParamToken]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div className="container column is-12">
        <div className="section">
          <section className="hero is-fullheight">
            <div className="hero-body">
              <div className="container">
                <div className="columns is-centered">
                  <div className="column is-half-tablet">
                    <div className="box is-rounded">
                      {/* Start Logo */}
                      <nav className="level">
                        <div className="level-item has-text-centered">
                          <figure className="image">
                            <Link to="/">
                              <img
                                src="/static/logo.png"
                                style={{ width: "256px" }}
                              />
                            </Link>
                          </figure>
                        </div>
                      </nav>
                      {/* End Logo */}

                      {/* Page */}
                      <form>
                        <h1 className="title is-3 has-text-centered">
                          Two-Factor Authentication
                        </h1>
                        <FormErrorBox errors={errors} />
                        <p class="has-text-grey">
                          Open the two-step verification app on your mobile
                          device, get your token and input here to finish your
                          login.
                        </p>
                        <p>&nbsp;</p>
                        <FormInputField
                          label="Enter your Token"
                          name="token"
                          placeholder="See your authenticator app"
                          value={token}
                          errorText={errors && errors.token}
                          helpText=""
                          onChange={(e) => setToken(e.target.value)}
                          isRequired={true}
                          maxWidth="380px"
                        />
                        <br />
                      </form>

                      <nav class="level">
                        <div class="level-left">
                          <div class="level-item">
                            <Link
                              class="button is-link is-fullwidth-mobile"
                              to="/login"
                            >
                              <FontAwesomeIcon icon={faArrowLeft} />
                              &nbsp;Back to Login
                            </Link>
                          </div>
                        </div>
                        <div class="level-right">
                          <div class="level-item">
                            <button
                              type="button"
                              class="button is-primary is-fullwidth-mobile"
                              onClick={onButtonClick}
                            >
                              <FontAwesomeIcon icon={faCheckCircle} />
                              &nbsp;Subit and Login
                            </button>
                          </div>
                        </div>
                      </nav>
                    </div>
                    {/* End box */}

                    <div className="has-text-centered">
                      <p>© 2024 BP8 Fitness Community</p>
                    </div>
                    {/* End suppoert text. */}
                  </div>
                  {/* End Column */}
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

export default TwoFactorAuthenticationValidateOnLogin;
