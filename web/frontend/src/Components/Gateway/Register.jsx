import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faArrowLeft,
  faCheckCircle,
  faArrowRight,
  faIdCard,
  faContactCard,
  faAddressBook,
  faChartPie,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { postRegisterAPI } from "../../API/gateway";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormSelectField from "../Reusable/FormSelectField";
import FormCheckboxField from "../Reusable/FormCheckboxField";
import FormCountryField from "../Reusable/FormCountryField";
import FormRegionField from "../Reusable/FormRegionField";
import { HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS } from "../../Constants/FieldOptions";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
  onHamburgerClickedState,
  currentUserState,
} from "../../AppState";

function Register() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [onHamburgerClicked, setOnHamburgerClicked] = useRecoilState(
    onHamburgerClickedState
  );
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [email, setEmail] = useState("");
  const [emailRepeated, setEmailRepeated] = useState("");
  const [phone, setPhone] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [password, setPassword] = useState("");
  const [passwordRepeated, setPasswordRepeated] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [addressLine1, setAddressLine1] = useState("");
  const [addressLine2, setAddressLine2] = useState("");
  const [city, setCity] = useState("");
  const [region, setRegion] = useState("");
  const [country, setCountry] = useState("");
  const [agreePromotionsEmail, setHasPromotionalEmail] = useState(true);
  const [agreeTOS, setAgreeTOS] = useState();
  const [howDidYouHearAboutUs, setHowDidYouHearAboutUs] = useState(0);
  const [howDidYouHearAboutUsOther, setHowDidYouHearAboutUsOther] =
    useState("");
  const [hasShippingAddress, setHasShippingAddress] = useState(false);
  const [shippingName, setShippingName] = useState("");
  const [shippingPhone, setShippingPhone] = useState("");
  const [shippingCountry, setShippingCountry] = useState("");
  const [shippingRegion, setShippingRegion] = useState("");
  const [shippingCity, setShippingCity] = useState("");
  const [shippingAddressLine1, setShippingAddressLine1] = useState("");
  const [shippingAddressLine2, setShippingAddressLine2] = useState("");
  const [shippingPostalCode, setShippingPostalCode] = useState("");

  ////
  //// Event handling.
  ////

  function onAgreePromotionsEmailChange(e) {
    setHasPromotionalEmail(!agreePromotionsEmail);
  }

  function onAgreeTOSChange(e) {
    setAgreeTOS(!agreeTOS);
  }

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});
    const decamelizedData = {
      // To Snake-case for API from camel-case in React.
      organization_id: process.env.REACT_APP_API_ORGANIZATION_ID,
      first_name: firstName,
      last_name: lastName,
      email: email,
      email_repeated: emailRepeated,
      phone: phone,
      postal_code: postalCode,
      address_line_1: addressLine1,
      address_line_2: addressLine2,
      city: city,
      region: region,
      country: country,
      status: 1,
      password: password,
      password_repeated: passwordRepeated,
      how_did_you_hear_about_us: howDidYouHearAboutUs,
      how_did_you_hear_about_us_other: howDidYouHearAboutUsOther,
      agree_promotions_email: agreePromotionsEmail,
      agree_tos: agreeTOS,
      has_shipping_address: hasShippingAddress,
      shipping_name: shippingName,
      shipping_phone: shippingPhone,
      shipping_country: shippingCountry,
      shipping_region: shippingRegion,
      shipping_city: shippingCity,
      shipping_address_line1: shippingAddressLine1,
      shipping_address_line2: shippingAddressLine2,
      shipping_postal_code: shippingPostalCode,
    };
    console.log("onSubmitClick, decamelizedData:", decamelizedData);
    postRegisterAPI(
      decamelizedData,
      onRegisterSuccess,
      onRegisterError,
      onRegisterDone
    );
  };

  ////
  //// API.
  ////

  function onRegisterSuccess(response) {
    // For debugging purposes only.
    console.log("onRegisterSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Registration successful");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onRegisterSuccess: Delayed for 2 seconds.");
      console.log(
        "onRegisterSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // Save the data to local storage for persistance in this browser and
    // redirect the user to their respected dahsboard.
    setOnHamburgerClicked(true); // Set to `true` so the side menu loads on startup of app.

    // Store in persistance storage in the browser.
    setCurrentUser(response.user);

    console.log("register", response.user);
    // Redirect the user to a new page.
    if (!response.user.onboardingCompleted && response.user.role === 4) {
      console.log("redrecting to onboarding");
      setForceURL("/onboarding");
    } else {
      setForceURL("/dashboard");
      console.log("redrecting to dashboard");
    }
  }

  function onRegisterError(apiErr) {
    console.log("onRegisterError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onRegisterError: Delayed for 2 seconds.");
      console.log(
        "onRegisterError: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onRegisterDone() {
    console.log("onRegisterDone: Starting...");
    setFetching(false);
  }

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

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faTasks} />
              &nbsp;Register
            </p>
            <FormErrorBox errors={errors} />

            {isFetching && (
              <PageLoadingContent displayMessage={"Please wait..."} />
            )}
            {!isFetching && (
              <div class="container">
                <p class="subtitle is-6">
                  <FontAwesomeIcon className="fas" icon={faIdCard} />
                  &nbsp;Details
                </p>

                <FormInputField
                  label="First Name"
                  name="firstName"
                  placeholder="Text input"
                  value={firstName}
                  errorText={errors && errors.firstName}
                  helpText=""
                  onChange={(e) => setFirstName(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <FormInputField
                  label="Last Name"
                  name="lastName"
                  placeholder="Text input"
                  value={lastName}
                  errorText={errors && errors.lastName}
                  helpText=""
                  onChange={(e) => setLastName(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <FormInputField
                  label="Password"
                  name="password"
                  type="password"
                  placeholder="Password input"
                  value={password}
                  errorText={errors && errors.password}
                  helpText=""
                  onChange={(e) => setPassword(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <FormInputField
                  label="Password Repeated"
                  name="passwordRepeated"
                  type="password"
                  placeholder="Password input"
                  value={passwordRepeated}
                  errorText={errors && errors.passwordRepeated}
                  helpText=""
                  onChange={(e) => setPasswordRepeated(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <p class="subtitle is-6">
                  <FontAwesomeIcon className="fas" icon={faContactCard} />
                  &nbsp;Contact Information
                </p>

                <FormInputField
                  label="Email"
                  name="email"
                  type="email"
                  placeholder="Text input"
                  value={email}
                  errorText={errors && errors.email}
                  helpText=""
                  onChange={(e) => setEmail(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <FormInputField
                  label="Email Repeated"
                  name="emailRepeated"
                  type="emailRepeated"
                  placeholder="Text input"
                  value={emailRepeated}
                  errorText={errors && errors.emailRepeated}
                  helpText="Please re-enter the above email again to confirm the email you entered is correct"
                  onChange={(e) => setEmailRepeated(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <FormInputField
                  label="Phone"
                  name="phone"
                  placeholder="Text input"
                  value={phone}
                  errorText={errors && errors.phone}
                  helpText=""
                  onChange={(e) => setPhone(e.target.value)}
                  isRequired={true}
                  maxWidth="150px"
                />

                <div class="columns">
                  <div class="column">
                    <p class="subtitle is-6">
                      {hasShippingAddress ? (
                        <p class="subtitle is-6">
                          <FontAwesomeIcon
                            className="fas"
                            icon={faAddressBook}
                          />
                          &nbsp;Billing Address
                        </p>
                      ) : (
                        <p class="subtitle is-6">
                          <FontAwesomeIcon
                            className="fas"
                            icon={faAddressBook}
                          />
                          &nbsp;Address
                        </p>
                      )}
                    </p>
                    <FormCountryField
                      priorityOptions={["CA", "US", "MX"]}
                      label="Country"
                      name="country"
                      placeholder="Text input"
                      selectedCountry={country}
                      errorText={errors && errors.country}
                      helpText=""
                      onChange={(value) => setCountry(value)}
                      isRequired={true}
                      maxWidth="160px"
                    />

                    <FormRegionField
                      label="Province/Territory"
                      name="region"
                      placeholder="Text input"
                      selectedCountry={country}
                      selectedRegion={region}
                      errorText={errors && errors.region}
                      helpText=""
                      onChange={(value) => setRegion(value)}
                      isRequired={true}
                      maxWidth="280px"
                    />

                    <FormInputField
                      label="City"
                      name="city"
                      placeholder="Text input"
                      value={city}
                      errorText={errors && errors.city}
                      helpText=""
                      onChange={(e) => setCity(e.target.value)}
                      isRequired={true}
                      maxWidth="380px"
                    />

                    <FormInputField
                      label="Address Line 1"
                      name="addressLine1"
                      placeholder="Text input"
                      value={addressLine1}
                      errorText={errors && errors.addressLine1}
                      helpText=""
                      onChange={(e) => setAddressLine1(e.target.value)}
                      isRequired={true}
                      maxWidth="380px"
                    />

                    <FormInputField
                      label="Address Line 2 (Optional)"
                      name="addressLine2"
                      placeholder="Text input"
                      value={addressLine2}
                      errorText={errors && errors.addressLine2}
                      helpText=""
                      onChange={(e) => setAddressLine2(e.target.value)}
                      isRequired={true}
                      maxWidth="380px"
                    />

                    <FormInputField
                      label="Postal Code"
                      name="postalCode"
                      placeholder="Text input"
                      value={postalCode}
                      errorText={errors && errors.postalCode}
                      helpText=""
                      onChange={(e) => setPostalCode(e.target.value)}
                      isRequired={true}
                      maxWidth="80px"
                    />

                    <FormCheckboxField
                      label="Has shipping address different then above address"
                      name="hasShippingAddress"
                      checked={hasShippingAddress}
                      errorText={errors && errors.hasShippingAddress}
                      onChange={(e) =>
                        setHasShippingAddress(!hasShippingAddress)
                      }
                      maxWidth="180px"
                    />
                  </div>
                  {hasShippingAddress && (
                    <div class="column">
                      <p class="subtitle is-6">
                        <FontAwesomeIcon className="fas" icon={faAddressBook} />
                        &nbsp;Shipping Address
                      </p>

                      <FormInputField
                        label="Name"
                        name="shippingName"
                        placeholder="Text input"
                        value={shippingName}
                        errorText={errors && errors.shippingName}
                        helpText="The name to contact for this shipping address"
                        onChange={(e) => setShippingName(e.target.value)}
                        isRequired={true}
                        maxWidth="350px"
                      />

                      <FormInputField
                        label="Phone"
                        name="shippingPhone"
                        placeholder="Text input"
                        value={shippingPhone}
                        errorText={errors && errors.shippingPhone}
                        helpText="The contact phone number for this shipping address"
                        onChange={(e) => setShippingPhone(e.target.value)}
                        isRequired={true}
                        maxWidth="150px"
                      />

                      <FormCountryField
                        priorityOptions={["CA", "US", "MX"]}
                        label="Country"
                        name="shippingCountry"
                        placeholder="Text input"
                        selectedCountry={shippingCountry}
                        errorText={errors && errors.shippingCountry}
                        helpText=""
                        onChange={(value) => setShippingCountry(value)}
                        isRequired={true}
                        maxWidth="160px"
                      />

                      <FormRegionField
                        label="Province/Territory"
                        name="shippingRegion"
                        placeholder="Text input"
                        selectedCountry={shippingCountry}
                        selectedRegion={shippingRegion}
                        errorText={errors && errors.shippingRegion}
                        helpText=""
                        onChange={(value) => setShippingRegion(value)}
                        isRequired={true}
                        maxWidth="280px"
                      />

                      <FormInputField
                        label="City"
                        name="shippingCity"
                        placeholder="Text input"
                        value={shippingCity}
                        errorText={errors && errors.shippingCity}
                        helpText=""
                        onChange={(e) => setShippingCity(e.target.value)}
                        isRequired={true}
                        maxWidth="380px"
                      />

                      <FormInputField
                        label="Address Line 1"
                        name="shippingAddressLine1"
                        placeholder="Text input"
                        value={shippingAddressLine1}
                        errorText={errors && errors.shippingAddressLine1}
                        helpText=""
                        onChange={(e) =>
                          setShippingAddressLine1(e.target.value)
                        }
                        isRequired={true}
                        maxWidth="380px"
                      />

                      <FormInputField
                        label="Address Line 2 (Optional)"
                        name="shippingAddressLine2"
                        placeholder="Text input"
                        value={shippingAddressLine2}
                        errorText={errors && errors.shippingAddressLine2}
                        helpText=""
                        onChange={(e) =>
                          setShippingAddressLine2(e.target.value)
                        }
                        isRequired={true}
                        maxWidth="380px"
                      />

                      <FormInputField
                        label="Postal Code"
                        name="shippingPostalCode"
                        placeholder="Text input"
                        value={shippingPostalCode}
                        errorText={errors && errors.shippingPostalCode}
                        helpText=""
                        onChange={(e) => setShippingPostalCode(e.target.value)}
                        isRequired={true}
                        maxWidth="80px"
                      />
                    </div>
                  )}
                </div>

                <p class="subtitle is-6">
                  <FontAwesomeIcon className="fas" icon={faChartPie} />
                  &nbsp;Metrics
                </p>

                <FormSelectField
                  label="How did you hear about us?"
                  name="howDidYouHearAboutUs"
                  placeholder="Pick"
                  selectedValue={howDidYouHearAboutUs}
                  errorText={errors && errors.howDidYouHearAboutUs}
                  helpText=""
                  onChange={(e) =>
                    setHowDidYouHearAboutUs(parseInt(e.target.value))
                  }
                  options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                />

                {howDidYouHearAboutUs === 1 && (
                  <FormInputField
                    label="Other (Please specify):"
                    name="howDidYouHearAboutUsOther"
                    placeholder="Text input"
                    value={howDidYouHearAboutUsOther}
                    errorText={errors && errors.howDidYouHearAboutUsOther}
                    helpText=""
                    onChange={(e) =>
                      setHowDidYouHearAboutUsOther(e.target.value)
                    }
                    isRequired={true}
                    maxWidth="380px"
                  />
                )}

                <FormCheckboxField
                  label="I agree to receive electronic updates"
                  name="agreePromotionsEmail"
                  checked={agreePromotionsEmail}
                  errorText={errors && errors.agreePromotionsEmail}
                  onChange={onAgreePromotionsEmailChange}
                  maxWidth="180px"
                />

                <FormCheckboxField
                  label="I agree to terms of service and privacy policy"
                  name="agreeTOS"
                  checked={agreeTOS}
                  errorText={errors && errors.agreeTos}
                  onChange={onAgreeTOSChange}
                  maxWidth="180px"
                />

                <div class="columns">
                  <div class="column is-half">
                    <Link
                      to={`/login`}
                      class="button is-medium is-fullwidth-mobile"
                    >
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back
                    </Link>
                  </div>
                  <div class="column is-half has-text-right">
                    <button
                      class="button is-medium is-primary is-fullwidth-mobile"
                      onClick={onSubmitClick}
                    >
                      <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                      &nbsp;Register
                    </button>
                  </div>
                </div>
              </div>
            )}
          </nav>
          <span className="content is-pulled-right has-text-grey">
            Already have an account?{" "}
            <Link to="/login">
              Click here&nbsp;
              <FontAwesomeIcon className="fas" icon={faArrowRight} />
            </Link>{" "}
            to sign in.
          </span>
        </section>
        <div className="content has-text-centered">
          <br />
          <p>Need help?</p>
          <p>
            <Link to="mailto:mike@bp8fitness.com">mike@bp8fitness.com</Link>
          </p>
          <p>
            <a href="tel:+16479672269">(647) 967-2269</a>
          </p>
        </div>
      </div>
    </>
  );
}

export default Register;
