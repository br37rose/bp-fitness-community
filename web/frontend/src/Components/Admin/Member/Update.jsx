import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { aTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faGauge, faPencil, faUsers, faEye, faIdCard, faAddressBook, faContactCard, faChartPie, faCogs } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';
import { useParams } from 'react-router-dom';

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { getMemberDetailAPI, putMemberUpdateAPI } from "../../../API/member";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS } from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";


function AdminMemberUpdate() {
    ////
    //// URL Parameters.
    ////

    const { bid, id } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [organizationID, setOrganizationID] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [email, setEmail] = useState("");
    const [phone, setPhone] = useState("");
    const [postalCode, setPostalCode] = useState("");
    const [addressLine1, setAddressLine1] = useState("");
    const [addressLine2, setAddressLine2] = useState("");
    const [city, setCity] = useState("");
    const [region, setRegion] = useState("");
    const [country, setCountry] = useState("");
    const [status, setStatus] = useState(0);
    const [agreePromotionsEmail, setHasPromotionalEmail] = useState(true);
    const [howDidYouHearAboutUs, setHowDidYouHearAboutUs] = useState(0);
    const [howDidYouHearAboutUsOther, setHowDidYouHearAboutUsOther] = useState("");
    const [password, setPassword] = useState("");
    const [passwordRepeated, setPasswordRepeated] = useState("");

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
        setErrors({});
        const decamelizedData = {
            id: id,
            organization_id: currentUser.organizationID,
            first_name: firstName,
            last_name: lastName,
            email: email,
            phone: phone,
            postal_code: postalCode,
            address_line_1: addressLine1,
            address_line_2: addressLine2,
            city: city,
            region: region,
            country: country,
            status: status,
            password: password,
            password_repeated: passwordRepeated,
            how_did_you_hear_about_us: howDidYouHearAboutUs,
            how_did_you_hear_about_us_other: howDidYouHearAboutUsOther,
            agree_promotions_email: agreePromotionsEmail,
        };
        console.log("onSubmitClick, decamelizedData:", decamelizedData);
        putMemberUpdateAPI(decamelizedData, onAdminMemberUpdateSuccess, onAdminMemberUpdateError, onAdminMemberUpdateDone);
    }

    function onMemberDetailSuccess(response){
        console.log("onMemberDetailSuccess: Starting...");
        setOrganizationID(response.organizationId);
        setFirstName(response.firstName);
        setLastName(response.lastName);
        setEmail(response.email);
        setPhone(response.phone);
        setPostalCode(response.postalCode);
        setAddressLine1(response.addressLine1);
        setAddressLine2(response.addressLine2);
        setCity(response.city);
        setRegion(response.region);
        setCountry(response.country);
        setStatus(response.status);
        setHowDidYouHearAboutUs(response.howDidYouHearAboutUs);
        setHowDidYouHearAboutUsOther(response.howDidYouHearAboutUsOther);
        setHasPromotionalEmail(response.agreePromotionsEmail);
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

    function onAdminMemberUpdateSuccess(response){
        // For debugging purposes only.
        console.log("onAdminMemberUpdateSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Member updated");
        setTopAlertStatus("Workout Member");
        setTimeout(() => {
            console.log("onAdminMemberUpdateSuccess: Delayed for 2 seconds.");
            console.log("onAdminMemberUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/admin/member/"+response.id);
    }

    function onAdminMemberUpdateError(apiErr) {
        console.log("onAdminMemberUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminMemberUpdateError: Delayed for 2 seconds.");
            console.log("onAdminMemberUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAdminMemberUpdateDone() {
        console.log("onAdminMemberUpdateDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            setFetching(true);
            getMemberDetailAPI(
                id,
                onMemberDetailSuccess,
                onMemberDetailError,
                onMemberDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL}  />
    }

    return (
        <>
            <div class="container">
                <section class="section">
                    {/* Desktop Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/members" aria-current="page"><FontAwesomeIcon className="fas" icon={faUsers} />&nbsp;Members</Link></li>
                            <li class=""><Link to={`/admin/member/${id}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Update</Link></li>
                        </ul>
                    </nav>

                    {/* Mobile Breadcrumbs */}
                    <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to={`/admin/member/${id}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Detail</Link></li>
                        </ul>
                    </nav>

                    {/* Page */}
                    <nav class="box">
                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faUsers} />&nbsp;Members</p>
                        <FormErrorBox errors={errors} />

                        {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

                        {isFetching
                            ?
                            <PageLoadingContent displayMessage={"Please wait..."} />
                            :
                            <div class="container" key={id}>

                                <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;Office Information</p>
                                <hr />

                                <FormInputField
                                    label="First Name"
                                    name="firstName"
                                    placeholder="Text input"
                                    value={firstName}
                                    errorText={errors && errors.firstName}
                                    helpText=""
                                    onChange={(e)=>setFirstName(e.target.value)}
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
                                    onChange={(e)=>setLastName(e.target.value)}
                                    isRequired={true}
                                    maxWidth="380px"
                                />

                                <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faContactCard} />&nbsp;Contact Information</p>
                                <hr />

                                <FormInputField
                                    label="Email"
                                    name="email"
                                    placeholder="Text input"
                                    value={email}
                                    errorText={errors && errors.email}
                                    helpText=""
                                    onChange={(e)=>setEmail(e.target.value)}
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
                                    onChange={(e)=>setPhone(e.target.value)}
                                    isRequired={true}
                                    maxWidth="150px"
                                />

                                <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faAddressBook} />&nbsp;Address</p>
                                <hr />

                                <FormCountryField
                                    priorityOptions={["CA","US","MX"]}
                                    label="Country"
                                    name="country"
                                    placeholder="Text input"
                                    selectedCountry={country}
                                    errorText={errors && errors.country}
                                    helpText=""
                                    onChange={(value)=>setCountry(value)}
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
                                    onChange={(value)=>setRegion(value)}
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
                                    onChange={(e)=>setCity(e.target.value)}
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
                                    onChange={(e)=>setAddressLine1(e.target.value)}
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
                                    onChange={(e)=>setAddressLine2(e.target.value)}
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
                                    onChange={(e)=>setPostalCode(e.target.value)}
                                    isRequired={true}
                                    maxWidth="80px"
                                />

                                <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faChartPie} />&nbsp;Metrics</p>
                                <hr />

                                <FormSelectField
                                    label="How did you hear about us?"
                                    name="howDidYouHearAboutUs"
                                    placeholder="Pick"
                                    selectedValue={howDidYouHearAboutUs}
                                    errorText={errors && errors.howDidYouHearAboutUs}
                                    helpText=""
                                    onChange={(e)=>setHowDidYouHearAboutUs(parseInt(e.target.value))}
                                    options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                                />

                                {howDidYouHearAboutUs === 1 && <FormInputField
                                    label="Other (Please specify):"
                                    name="howDidYouHearAboutUsOther"
                                    placeholder="Text input"
                                    value={howDidYouHearAboutUsOther}
                                    errorText={(e)=>setHowDidYouHearAboutUsOther(e.target.value)}
                                    helpText=""
                                    onChange={(e)=>setHowDidYouHearAboutUsOther(e.target.value)}
                                    isRequired={true}
                                    maxWidth="380px"
                                />}

                                <FormCheckboxField
                                    label="I agree to receive electronic updates from my local gym"
                                    name="agreePromotionsEmail"
                                    checked={agreePromotionsEmail}
                                    errorText={errors && errors.agreePromotionsEmail}
                                    onChange={onAgreePromotionsEmailChange}
                                    maxWidth="180px"
                                />

                                <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faCogs} />&nbsp;Settings</p>
                                <hr />

                                <FormRadioField
                                    label="Status"
                                    name="status"
                                    placeholder="Pick"
                                    value={status}
                                    opt2Value={1}
                                    opt2Label="Active"
                                    opt4Value={100}
                                    opt4Label="Archived"
                                    errorText={errors && errors.status}
                                    onChange={(e)=>setStatus(parseInt(e.target.value))}
                                    maxWidth="180px"
                                    disabled={false}
                                />

                                <FormInputField
                                    label="Password (Optional)"
                                    name="password"
                                    type="password"
                                    placeholder="Text input"
                                    value={password}
                                    errorText={errors && errors.password}
                                    helpText=""
                                    onChange={(e)=>setPassword(e.target.value)}
                                    isRequired={true}
                                    maxWidth="380px"
                                />

                                <FormInputField
                                    label="Password Repeated (Optional)"
                                    name="passwordRepeated"
                                    type="password"
                                    placeholder="Text input"
                                    value={passwordRepeated}
                                    errorText={errors && errors.passwordRepeated}
                                    helpText=""
                                    onChange={(e)=>setPasswordRepeated(e.target.value)}
                                    isRequired={true}
                                    maxWidth="380px"
                                />

                                <div class="columns pt-5">
                                    <div class="column is-half">
                                        <Link class="button is-hidden-touch" to={`/admin/member/${id}`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back</Link>
                                        <Link class="button is-fullwidth is-hidden-desktop" to={`/admin/member/${id}`}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back</Link>
                                    </div>
                                    <div class="column is-half has-text-right">
                                        <button class="button is-primary is-hidden-touch" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                        <button class="button is-primary is-fullwidth is-hidden-desktop" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                    </div>
                                </div>

                            </div>
                        }
                    </nav>
                </section>
            </div>
        </>
    );
}

export default AdminMemberUpdate;
