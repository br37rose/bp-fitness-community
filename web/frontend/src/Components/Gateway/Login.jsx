import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowRight,
  faEnvelope,
  faKey,
  faTriangleExclamation,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import { postLoginAPI } from "../../API/gateway";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { onHamburgerClickedState, currentUserState } from "../../AppState";
import {
  ROOT_ROLE_ID,
  ADMIN_ROLE_ID,
  TRAINER_ROLE_ID,
  MEMBER_ROLE_ID,
} from "../../Constants/App";

function Login() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const isUnauthorized = searchParams.get("unauthorized");

  ////
  //// Global state.
  ////

  const [onHamburgerClicked, setOnHamburgerClicked] = useRecoilState(
    onHamburgerClickedState
  );
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  // const [errors, setErrors] = useState({
  //     "email": "account does not exist",
  //     "password": "invalid password"
  // });
  const [errors, setErrors] = useState({});
  const [validation, setValidation] = useState({
    email: false,
    password: false,
  });
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [forceURL, setForceURL] = useState("");
  const [isFetching, setFetching] = useState(false);

  ////
  //// API.
  ////

  function onLoginSuccess(response) {
    console.log("onLoginSuccess: Starting...");

    if (
      response.user === undefined ||
      response.user === null ||
      response.user === ""
    ) {
      // Defensive code.
      alert(
        "user null error with login endpoint - please contact system administrator."
      );
      return;
    }

    // For debugging purposes only.
    console.log("onLoginSuccess: onHamburgerClicked:", onHamburgerClicked);
    console.log("onLoginSuccess: user:", response.user);

    // Save the data to local storage for persistance in this browser and
    // redirect the user to their respected dahsboard.
    setOnHamburgerClicked(true); // Set to `true` so the side menu loads on startup of app.

    // Store in persistance storage in the browser.
    setCurrentUser(response.user);

    if (
      response.user.otpEnabled === undefined ||
      response.user.otpEnabled === null ||
      response.user.otpEnabled === "" ||
      response.user.otpEnabled === false
    ) {
      // Proceed to redirect the user based on their assigned role.
      switch (response.user.role) {
        case ROOT_ROLE_ID:
          console.log(
            "onLoginSuccess: user is executive (root) admin, redirecting to admin dashboard..."
          );
          // setForceURL("/admin/dashboard");
          alert(
            "unsupported user role detected - executive (`root`) user accounts are not supported at this time."
          );
          break;
        case ADMIN_ROLE_ID:
          console.log(
            "onLoginSuccess: user is admin, redirecting to admin dashboard..."
          );
          setForceURL("/admin/dashboard");
          break;
        case TRAINER_ROLE_ID:
          console.log(
            "onLoginSuccess: user is trainer, redirecting to trainer dashboard..."
          );
          setForceURL("/trainer/dashboard");
          break;
        case MEMBER_ROLE_ID:
          console.log(
            "onLoginSuccess: user is member, redirecting to member dashboard...",
            response.user.onboardingCompleted
          );
          if (!response.user.onboardingCompleted) {
            setForceURL("/onboarding");
          } else {
            setForceURL("/dashboard");
          }
          break;
        default:
          console.log("onLoginSuccess: critical error");
          alert(
            "onLoginSuccess: unsupported user role detected - please contact system administrator."
          );
      }
    } else {
      if (response.user.otpVerified === false) {
        console.log("onLoginSuccess | redirecting to 2fa setup wizard");
        setForceURL("/login/2fa/step-1");
      } else {
        console.log("onLoginSuccess | redirecting to 2fa validation");
        setForceURL("/login/2fa");
      }
    }
  }

  function onLoginError(apiErr) {
    console.log("onLoginError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onLoginDone() {
    console.log("onLoginDone: Starting...");
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  function onEmailChange(e) {
    setEmail(e.target.value);
    validation["email"] = false;
    setValidation(validation);
    // setErrors(errors["email"]="");
  }

  function onPasswordChange(e) {
    setPassword(e.target.value);
    validation["password"] = false;
    setValidation(validation);
    // setErrors(errors["password"]="");
  }

  function onButtonClick(e) {
    var newErrors = {};
    var newValidation = {};
    if (email === undefined || email === null || email === "") {
      newErrors["email"] = "value is missing";
    } else {
      newValidation["email"] = true;
    }
    if (password === undefined || password === null || password === "") {
      newErrors["password"] = "value is missing";
    } else {
      newValidation["password"] = true;
    }

    /// Save to state.
    setErrors(newErrors);
    setValidation(newValidation);

    if (Object.keys(newErrors).length > 0) {
      //
      // Handle errors.
      //

      console.log("failed validation");

      // window.scrollTo(0, 0);  // Start the page at the top of the page.

      // The following code will cause the screen to scroll to the top of
      // the page. Please see ``react-scroll`` for more information:
      // https://github.com/fisshy/react-scroll
      var scroll = Scroll.animateScroll;
      scroll.scrollToTop();
    } else {
      //
      // Submit to server.
      //

      console.log("successful validation, submitting to API server.");

      const postData = {
        email: email,
        password: password,
      };
      setFetching(true);
      postLoginAPI(postData, onLoginSuccess, onLoginError, onLoginDone);
    }
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
    }

    return () => (mounted = false);
  }, []);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="column is-12">
        {isFetching && <PageLoadingContent displayMessage={"Please wait..."} />}
        {!isFetching && (
          <div>
            <section class="hero is-fullheight">
              <div class="hero-body">
                <div class="container is-fluid">
                  <div class="columns is-centered">
                    <div class="column is-one-third-tablet">
                      <div class="box is-rounded">
                        {/* Start Logo */}
                        <nav class="level">
                          <div class="level-item has-text-centered">
                            <figure class="image">
                              <img
                                src="/static/login_header_logo.jpg"
                                style={{ width: "120px" }}
                              />
                            </figure>
                          </div>
                        </nav>
                        {/* End Logo */}
                        <form>
                          <h1 className="title is-4 has-text-centered">
                            Sign In
                          </h1>
                          {isUnauthorized === "true" && (
                            <article class="message is-danger">
                              <div class="message-body">
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faTriangleExclamation}
                                />
                                &nbsp;Your session has ended.
                                <br />
                                Please login again
                              </div>
                            </article>
                          )}
                          <FormErrorBox errors={errors} />

                          <div class="field">
                            <label class="label is-small has-text-grey-light">
                              Email
                            </label>
                            <div class="control has-icons-left has-icons-right">
                              <input
                                class={`input ${
                                  errors && errors.email && "is-danger"
                                } ${
                                  validation && validation.email && "is-success"
                                }`}
                                type="email"
                                placeholder="Email"
                                value={email}
                                onChange={onEmailChange}
                              />
                              <span class="icon is-small is-left">
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faEnvelope}
                                />
                              </span>
                            </div>
                            {errors && errors.email && (
                              <p class="help is-danger">{errors.email}</p>
                            )}
                          </div>

                          <div class="field">
                            <label class="label is-small has-text-grey-light">
                              Password
                            </label>
                            <div class="control has-icons-left has-icons-right">
                              <input
                                class={`input ${
                                  errors && errors.password && "is-danger"
                                } ${
                                  validation &&
                                  validation.password &&
                                  "is-success"
                                }`}
                                type="password"
                                placeholder="Password"
                                value={password}
                                onChange={onPasswordChange}
                              />
                              <span class="icon is-small is-left">
                                <FontAwesomeIcon className="fas" icon={faKey} />
                              </span>
                            </div>
                            {errors && errors.password && (
                              <p class="help is-danger">{errors.password}</p>
                            )}
                          </div>
                          <br />
                          <Link
                            class="button is-fullwidth is-primary"
                            type="button"
                            onClick={onButtonClick}
                          >
                            Login <FontAwesomeIcon icon={faArrowRight} />
                          </Link>
                        </form>
                        <br />
                        <nav class="level">
                          <div class="level-item has-text-centered">
                            <div>
                              <Link
                                to="/forgot-password"
                                className="is-size-7-tablet"
                              >
                                Forgot Password?
                              </Link>
                            </div>
                          </div>
                          <div class="level-item has-text-centered">
                            <div>
                              <Link to="/register" className="is-size-7-tablet">
                                Create an Account
                              </Link>
                            </div>
                          </div>
                        </nav>
                      </div>
                      {/* End box */}

                      <div className="columns">
                        <div className="column has-text-centered">
                          <p>Need help?</p>
                          <p>
                            <Link to="mailto:mike@bp8fitness.com">
                              mike@bp8fitness.com
                            </Link>
                          </p>
                          <p>
                            <a href="tel:+16479672269">(647) 967-2269</a>
                          </p>
                        </div>
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
        )}
      </div>
    </>
  );
}

export default Login;
