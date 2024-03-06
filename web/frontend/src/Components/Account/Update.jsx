import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faContactCard, faChartPie } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import useLocalStorage from "../../Hooks/useLocalStorage";
import { getAccountDetailAPI, putAccountUpdateAPI } from "../../API/Account";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormRadioField from "../Reusable/FormRadioField";
import FormMultiSelectField from "../Reusable/FormMultiSelectField";
import FormSelectField from "../Reusable/FormSelectField";
import FormCheckboxField from "../Reusable/FormCheckboxField";
import FormCountryField from "../Reusable/FormCountryField";
import FormRegionField from "../Reusable/FormRegionField";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import Layout from "../Menu/Layout";


function AccountUpdate() {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [email, setEmail] = useState("");
    const [phone, setPhone] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [password, setPassword] = useState("");
    const [passwordRepeated, setPasswordRepeated] = useState("");
    const [companyName, setCompanyName] = useState("");
    const [postalCode, setPostalCode] = useState("");
    const [addressLine1, setAddressLine1] = useState("");
    const [addressLine2, setAddressLine2] = useState("");
    const [city, setCity] = useState("");
    const [region, setRegion] = useState("");
    const [country, setCountry] = useState("");
    const [agreePromotionsEmail, setHasPromotionalEmail] = useState(true);

    ////
    //// Event handling.
    ////

    function onAgreePromotionsEmailChange(e) {
        setHasPromotionalEmail(!agreePromotionsEmail);
    }

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Beginning...");
        setFetching(true);

        const submission = {
            Email: email,
            Phone: phone,
            FirstName: firstName,
            LastName: lastName,
            Password: password,
            PasswordRepeated: passwordRepeated,
            CompanyName: companyName,
            PostalCode: postalCode,
            AddressLine1: addressLine1,
            AddressLine2: addressLine2,
            City: city,
            Region: region,
            Country: country,
            AgreePromotionsEmail: agreePromotionsEmail,
        };
        console.log("onSubmitClick, submission:", submission);
        putAccountUpdateAPI(submission, onAccountUpdateSuccess, onAccountUpdateError, onAccountUpdateDone);
    }

    function onAccountDetailSuccess(response) {
        console.log("onAccountDetailSuccess: Starting...");
        setEmail(response.email);
        setPhone(response.phone);
        setFirstName(response.firstName);
        setLastName(response.lastName);
        setCompanyName(response.companyName);
        setPostalCode(response.postalCode);
        setAddressLine1(response.addressLine1);
        setAddressLine2(response.addressLine2);
        setCity(response.city);
        setRegion(response.region);
        setCountry(response.country);
        setHasPromotionalEmail(response.agreePromotionsEmail);
        setCurrentUser(response);
    }

    function onAccountDetailError(apiErr) {
        console.log("onAccountDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountDetailDone() {
        console.log("onAccountDetailDone: Starting...");
        setFetching(false);
    }

    function onAccountUpdateSuccess(response) {
        // For debugging purposes only.
        console.log("onAccountUpdateSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Account updated");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAccountUpdateSuccess: Delayed for 2 seconds.");
            console.log("onAccountUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/account");
    }

    function onAccountUpdateError(apiErr) {
        console.log("onAccountUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAccountUpdateError: Delayed for 2 seconds.");
            console.log("onAccountUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountUpdateDone() {
        console.log("onAccountUpdateDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////

    const generateBreadcrumbItemLink = (currentUser) => {
        let dashboardLink;
        switch (currentUser.role) {
            case 1:
                dashboardLink = "/root/dashboard";
                break;
            case 2:
                dashboardLink = "/admin/dashboard";
                break;
            case 3:
                dashboardLink = "/trainer/dashboard";
                break;
            case 4:
                dashboardLink = "/dashboard";
                break;
            default:
                dashboardLink = "/"; // Default or error handling
                break;
        }
        return dashboardLink;
    };

    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: generateBreadcrumbItemLink(currentUser), isActive: false, icon: faGauge },
            { text: 'Account', link: '/account', icon: faUserCircle, isActive: false },
            { text: 'Edit', link: '#', icon: faPencil, isActive: true }
        ],
        mobileBackLinkItems: {
            link: '/account',
            text: "Back to Account",
            icon: faArrowLeft
        }
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            setFetching(true);
            getAccountDetailAPI(
                onAccountDetailSuccess,
                onAccountDetailError,
                onAccountDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div class="box">
                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faUserCircle} />&nbsp;Account</p>
                <FormErrorBox errors={errors} />

                {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

                {isFetching
                    ?
                    <PageLoadingContent displayMessage={"Please wait..."} />
                    : <>
                        <div class="container">

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;Full Name</p>
                            <hr />

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

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faContactCard} />&nbsp;Contact Information</p>
                            <hr />

                            <FormInputField
                                label="Email"
                                name="email"
                                placeholder="Text input"
                                value={email}
                                errorText={errors && errors.email}
                                helpText=""
                                onChange={(e) => setEmail(e.target.value)}
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
                                maxWidth="380px"
                            />

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faAddressBook} />&nbsp;Address</p>
                            <hr />

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
                                label="Address Line 2"
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
                                maxWidth="380px"
                            />

                            <p class="title is-6"><FontAwesomeIcon className="fas" icon={faChartPie} />&nbsp;Metrics</p>
                            <hr />

                            <FormCheckboxField
                                label="I agree to receive electronic updates from my local branch and/or BP8 Fitness."
                                name="agreePromotionsEmail"
                                checked={agreePromotionsEmail}
                                errorText={errors && errors.agreePromotionsEmail}
                                onChange={onAgreePromotionsEmailChange}
                                maxWidth="180px"
                            />

                            <div class="columns pt-5">
                                <div class="column is-half">
                                    <Link class="button is-hidden-touch" to={"/account"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back</Link>
                                    <Link class="button is-fullwidth is-hidden-desktop" to={"/account"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back</Link>
                                </div>
                                <div class="column is-half has-text-right">
                                    <button class="button is-primary is-hidden-touch" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                    <button class="button is-primary is-fullwidth is-hidden-desktop" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                </div>
                            </div>

                        </div>
                    </>
                }
            </div>
        </Layout>
    );
}

export default AccountUpdate;
